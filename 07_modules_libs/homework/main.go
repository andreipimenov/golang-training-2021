package main

import (
	"github.com/joho/godotenv"
	"github.com/rabbit72/golang-training-2021/07_modules_libs/homework/server"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	err := server.Serve(":8888")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
