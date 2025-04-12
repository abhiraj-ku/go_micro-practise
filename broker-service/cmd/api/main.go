package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const webPort = "8181"

type Config struct {
	httpClient *http.Client

	// TODO: Learn and implement the circuitBreaker pattern for resilliency
	// circuitBreaker CircuitBreaker
}

func main() {
	cfg := Config{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 25,
				IdleConnTimeout:     90 * time.Second,
			},
		}}
	log.Printf("Starting broker service on port %s\n", webPort)

	// define the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: cfg.Routes(),
	}

	// start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
