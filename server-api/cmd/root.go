package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/computerdane/bop/bop"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	dir     string
	baseUrl string
	port    int
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	bop.UnimplementedBopServer
}

func (s *server) List(_ context.Context, in *bop.ListRequest) (*bop.ListReply, error) {
	var names []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		path = strings.Replace(path, dir+"/", "", 1)
		if err == nil && !d.IsDir() {
			names = append(names, path)
		}
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "walking directory tree failed")
	}
	matches := fuzzy.RankFindNormalizedFold(strings.ToLower(in.GetSearch()), names)
	sort.Reverse(matches)
	urls := make([]string, len(matches))
	for i, match := range matches {
		urls[i] = baseUrl + match.Target
	}
	return &bop.ListReply{Name: urls}, nil
}

var rootCmd = &cobra.Command{
	Use:   "bop-api listen [dir] [baseUrl]",
	Short: "API server for bop",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "listen" {
			log.Fatalf("Usage: %s", cmd.Use)
		}
		dir = args[1]
		baseUrl = args[2]
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		bop.RegisterBopServer(s, &server{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8085, "port to listen on")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
