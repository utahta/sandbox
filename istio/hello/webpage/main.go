package main

import (
	"bytes"
	"log"
	"net/http"
)

var version string

func main() {
	http.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var body bytes.Buffer
		body.WriteString("hello world")

		w.WriteHeader(http.StatusOK)
		w.Write(body.Bytes())

		log.Printf("[%v] %s", version, body.String())
	}))

	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
