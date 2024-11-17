package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/computerdane/bop/bop"
	"github.com/computerdane/bop/lib"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cfgFile          string
	shouldSaveConfig bool

	host          string
	mpvArgs       string
	port          int
	shouldShuffle bool
)

var rootCmd = &cobra.Command{
	Use:   "bop [search]",
	Short: "Bop your songs",

	Run: func(cmd *cobra.Command, args []string) {
		// connect to grpc
		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			Crash("did not connect: ", err)
		}
		defer conn.Close()
		client := bop.NewBopClient(conn)

		// list urls based on search
		request := bop.ListRequest{}
		if len(args) > 0 {
			search := strings.Join(args, " ")
			request.Search = &search
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := client.List(ctx, &request)
		if err != nil {
			Crash(err)
		}
		names := reply.GetName()
		if len(names) == 0 {
			Crash("no results")
		}

		if shouldShuffle {
			rand.Shuffle(len(names), func(i, j int) {
				names[i], names[j] = names[j], names[i]
			})
		} else {
			// group songs by directory, keeping ranking of dirs
			groups := make(map[string]([]string))
			dirs := []string{}
			for _, name := range names {
				dir := path.Dir(name)
				if _, ok := groups[dir]; ok {
					groups[dir] = append(groups[dir], name)
				} else {
					dirs = append(dirs, dir)
					groups[dir] = []string{name}
				}
			}

			// sort songs within their directories (great for playing albums with numbered songs)
			i := 0
			for _, dir := range dirs {
				group := groups[dir]
				sort.Strings(group)
				for _, name := range group {
					names[i] = name
					i++
				}
			}
		}

		// launch mpv with urls
		mpvArgsArray := append(strings.Split(mpvArgs, " "), names...)
		if err := exec.Command("mpv", mpvArgsArray...).Start(); err != nil {
			Crash(err)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/bop/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&shouldSaveConfig, "save-config", false, "save to the config file with the provided flags")

	lib.AddOption(rootCmd, lib.Option{P: &host, Name: "host", Shorthand: "H", Value: "localhost", Usage: "api host without port"})
	lib.AddOption(rootCmd, lib.Option{P: &mpvArgs, Name: "mpv-args", Shorthand: "", Value: "--force-window --title=${filename}", Usage: "args to pass to mpv"})
	lib.AddOption(rootCmd, lib.Option{P: &port, Name: "port", Shorthand: "P", Value: 8085, Usage: "api port"})
	lib.AddOption(rootCmd, lib.Option{P: &shouldShuffle, Name: "shuffle", Shorthand: "s", Value: false, Usage: "shuffle songs"})
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			Crash(err)
		}
		cfgFile = home + "/.config/bop/config.yaml"
	}
	viper.SetConfigFile(cfgFile)

	if _, err := os.Stat(cfgFile); err != nil {
		genConfig()
	}

	if err := viper.ReadInConfig(); err == nil {
		lib.LoadOptions()
	}

	if shouldSaveConfig {
		genConfig()
	}
}

func genConfig() {
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
