package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/utahta/sandbox/istio/hello/helloworld"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const (
	port        = ":5000"
	monitorPort = ":18000"
)

var version string

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("[2][%s] Received: %v\n", version, in.Name)
	return &helloworld.HelloReply{Message: "Hello2 " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})

	var eg errgroup.Group
	eg.Go(func() error {
		return s.Serve(lis)
	})

	eg.Go(func() error {
		http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		return http.ListenAndServe(monitorPort, nil)
	})

	if err := eg.Wait(); err != nil {
		log.Fatalf("[%v] %v", version, err)
	}
}
