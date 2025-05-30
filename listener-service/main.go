package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/abhiraj-ku/go_micro-practise/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// conect to rabbitmq
	rabbitConn, err := connectMq()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to rabbitMQ")

	// start listening from messages
	log.Println("Listening for messages from rabbitmq")

	// create consumers for the messages
	consumer, err := events.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

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
