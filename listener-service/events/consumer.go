package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// constructor to init the consumer

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup consumer
func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// loop over the topics and read it
	for _, topic := range topics {
		ch.QueueBind(
			q.Name,
			topic,
			"log_topic",
			false,
			nil,
		)
		// if err != nil {
		// 	return err

		// }
	}
	// consume things
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	foreverConsume := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			// another go routine to handle the payloads
			go handlePayload(payload)
		}
	}()
	fmt.Printf("Waiting for message [Exchange,Queue] [logs_topic, %s]\n", q.Name)

	// reciver infinetely
	<-foreverConsume
	return nil
}

// handlePayload() handles the payload and do some stuffs
func handlePayload(data Payload) {
	switch data.Name {
	case "log", "event":
		err := logEvent(data)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		//authentication data
	default:
		err := logEvent(data)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil

}
