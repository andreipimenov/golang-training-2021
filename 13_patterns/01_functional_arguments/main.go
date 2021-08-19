package main

import (
	"log"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK"}`))
}

func main() {
	s := NewServer(
		WithAddr(":8000"),
		WithRoute("/health", health),
	)

	log.Println("Server is listening on 8000")
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", health)
	mux.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"status": "Not implemented"}`))
	})
	return mux
}
