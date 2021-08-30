package main

import (
	"log"
	"net/http"

	"github.com/andreipimenov/golang-training-2021/13_patterns/02_method_chaining/server"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": "PONG"}`))
}

func main() {
	s := server.New().Addr(":8000").Route("/ping", ping)

	log.Println("Server is listening on 8000")

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
