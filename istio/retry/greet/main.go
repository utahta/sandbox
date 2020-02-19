package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/utahta/sandbox/istio/hello/helloworld"
	"google.golang.org/grpc"
)

const (
	monitorPort = ":18000"
)

var version string

func main() {
	go func() {
		http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		http.ListenAndServe(monitorPort, nil)
	}()

	var wg sync.WaitGroup
	addrs := strings.Split(os.Getenv("ECHO_ADDR"), ",")
	for _, addr := range addrs {
		addr := strings.TrimSpace(addr)
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()

			conn, err := grpc.Dial(
				addr,
				grpc.WithInsecure(),
				grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
			)
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := helloworld.NewGreeterClient(conn)

			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					func() {
						ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
						defer cancel()

						var wg sync.WaitGroup
						wg.Add(1)
						go func() {
							defer wg.Done()
							r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: "world"})
							if err != nil {
								log.Printf("[%s][ERROR] could not SayHello: %v\n", version, err)
								return
							}
							log.Printf("[%s] SayHello: %s\n", version, r.Message)
						}()

						wg.Add(1)
						go func() {
							defer wg.Done()
							r, err := c.SayMorning(ctx, &helloworld.MorningRequest{Name: "world"})
							if err != nil {
								log.Printf("[%s][ERROR] could not SayMorning: %v\n", version, err)
								return
							}
							log.Printf("[%s] SayMorning: %s\n", version, r.Message)
						}()

						wg.Wait()
					}()
				}
			}
		}(addr)
	}
	wg.Wait()
}
