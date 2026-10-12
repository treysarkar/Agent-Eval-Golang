package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           NewServer(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Listening at 8080")
	log.Fatal(server.ListenAndServe())
}
