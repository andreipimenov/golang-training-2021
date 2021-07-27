package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Working")

	r := chi.NewRouter()

	r.Get("/price/{ticker}/stat", stockHandler)

	err := http.ListenAndServe(":80", r)
	if err != nil {
		log.Fatal(err)
	}
}