package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8181"

type Config struct {
	httpClient *http.Client
	Rabbitmq   *amqp.Connection

	// TODO: Learn and implement the circuitBreaker patte
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

	// conect to rabbitmq
	rabbitConn, err := connectMq()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to rabbitMQ")
	app := Config{
		Rabbitmq: rabbitConn,
	}

	// define the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: cfg.Routes(),
	}

	// start the server
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// connect to rabbit RabbitMQ

func connectMq() (*amqp.Connection, error) {
	var count int64
	var backOffTime = 2 * time.Second
	maxBackOffTime := 30 * time.Second

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Printf("RabbitMQ not yet ready... %v", err)
			count++
		} else {
			log.Println("connected to rabbitMq")
			return conn, nil
		}

		if count > 5 {
			return nil, fmt.Errorf("failed to connect with rabbitmq after multiple attempts: %w", err)
		}

		// add some jitterance for better connection
		jitter := time.Duration(int64((backOffTime)))
		backOffTime = time.Duration(math.Min(float64(maxBackOffTime), float64(backOffTime*2)))
		log.Printf("backing off for %v seconds (with jitter)...", backOffTime+jitter)
		time.Sleep(backOffTime + jitter)

	}
}
