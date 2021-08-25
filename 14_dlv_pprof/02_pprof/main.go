package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"
)

const (
	port       = 8080
	difficulty = 5
	timeout    = time.Second
)

func main() {
	http.HandleFunc("/mining", MiningHandler(difficulty, timeout))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Server is listening on :%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
