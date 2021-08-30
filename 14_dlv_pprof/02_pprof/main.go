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
	timeout    = 30 * time.Second
)

func main() {
	http.HandleFunc("/mining", MiningHandler(difficulty, timeout))

	log.Printf("Server is listening on :%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
