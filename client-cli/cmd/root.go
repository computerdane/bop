package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"github.com/computerdane/bop/bop"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cfgFile string
	addr    string
)

var rootCmd = &cobra.Command{
	Use:   "bop [dir]",
	Short: "Bop your songs",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			Crash("did not connect: ", err)
		}
		defer conn.Close()
		c := bop.NewBopClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.List(ctx, &bop.ListRequest{})
		if err != nil {
			Crash("could not greet: ", err)
		}

		mpvPath, err := exec.LookPath("mpv")
		if err != nil {
			Crash(err)
		}

		if err := syscall.Exec(mpvPath, r.GetName(), os.Environ()); err != nil {
			Crash(err)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/bop/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&addr, "addr", "localhost:8085", "addr of api")

	viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			Crash(err)
		}
		cfgFile = home + "/.config/bop/config.yaml"
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err == nil {
		addr = viper.GetString("addr")
	}

	// try to generate config file
	cfgFileDir := path.Dir(cfgFile)
	if err := os.MkdirAll(cfgFileDir, os.ModePerm); err != nil {
		Warn("failed to make config directory: ", err)
	}
	if _, err := os.OpenFile(cfgFile, os.O_CREATE|os.O_RDONLY, 0600); err != nil {
		Warn("failed to create config file: ", err)
	}
	if err := viper.WriteConfig(); err != nil {
		Warn("failed to generate config: ", err)
	}
}

var (
	red    = color.New(color.FgRed).FprintlnFunc()
	yellow = color.New(color.FgYellow).FprintlnFunc()
)

func Warn(a ...any) {
	yellow(os.Stderr, a...)
}

func Crash(a ...any) {
	if len(a) == 0 {
		red(os.Stderr, "unknown error!")
	} else {
		red(os.Stderr, a...)
	}
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
