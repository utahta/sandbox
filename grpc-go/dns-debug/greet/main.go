package main

import (
	"context"
	"log"
	"time"

	"github.com/utahta/sandbox/grpc-go/dns-debug/greet/dns"
	"github.com/utahta/sandbox/grpc-go/dns-debug/helloworld"
	"google.golang.org/grpc"
)

const (
	//address = "mydns:///echo-server-hs.echo-server.svc.cluster.local:5000"
	address = "dns:///echo-server-hs.echo-server.svc.cluster.local:5000"
)

var version string

func main() {
	// Register gRPC resolver for client-side load balancing.
	dns.Register()

	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithBalancerName("round_robin"),
		grpc.WithDisableServiceConfig(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := helloworld.NewGreeterClient(conn)

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()

				r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: "world"})
				if err != nil {
					log.Printf("[ERROR] could not greet: %v\n", err)
					return
				}
				log.Printf("Greeting[%s]: %s\n", version, r.Message)
			}()
		}
	}
}
