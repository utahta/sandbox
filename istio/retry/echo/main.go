package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/utahta/sandbox/istio/hello/helloworld"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port        = ":5000"
	monitorPort = ":18000"
)

var version string

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	log.Printf("[%s] Received SayHello: %v\n", version, in.Name)

	if err := failureErr(); err != nil {
		return nil, err
	}
	return &helloworld.HelloResponse{Message: fmt.Sprintf("Hello %s from %s", in.Name, os.Getenv("HOSTNAME"))}, nil
}

func (s *server) SayMorning(ctx context.Context, in *helloworld.MorningRequest) (*helloworld.MorningResponse, error) {
	log.Printf("[%s] Received SayMorning: %v\n", version, in.Name)

	if err := failureErr(); err != nil {
		return nil, err
	}
	return &helloworld.MorningResponse{Message: fmt.Sprintf("Good morning %s from %s", in.Name, os.Getenv("HOSTNAME"))}, nil
}

func failureErr() error {
	if s := os.Getenv("FAILURE_RATE"); s != "" {
		rate, _ := strconv.Atoi(s) // [0, 100]
		if rate > rand.Intn(100) {
			code := codes.Unavailable
			for i := codes.OK; i <= codes.Unauthenticated; i++ {
				switch strings.ToLower(os.Getenv("RESPONSE_ERROR_CODE")) {
				case strings.ToLower(i.String()):
					code = i
				}
			}
			return  status.Error(code, fmt.Sprintf("error from %s", os.Getenv("HOSTNAME")))
		}
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
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
