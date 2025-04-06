package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8181"

type Config struct{}

func main() {
	cfg := Config{}
	log.Printf("Starting broker service on port %s\n", webPort)

	// define the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: cfg.routes(),
	}

	// start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
