package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/utahta/sandbox/istio/cb/helloworld"
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

	addr := strings.TrimSpace(os.Getenv("ECHO_ADDR"))
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

	numClients, _ := strconv.Atoi(os.Getenv("NUM_CLIENTS"))
	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		id := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					func() {
						ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
						defer cancel()

						var wg sync.WaitGroup
						wg.Add(1)
						go func() {
							defer wg.Done()
							r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: fmt.Sprintf("world: %d", id)})
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
		}()
	}

	wg.Wait()
}
