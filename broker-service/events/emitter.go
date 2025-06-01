package events

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	return declareExchange(channel)
}

// Push() -> pushes messages to the queue

func (e *Emitter) Push(event string, sev string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	log.Printf("Pushing event to %v channel", channel)

	err = channel.Publish(
		"logs_event",
		sev,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewEmitter constructor -> to init the Emitter type
func NewEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}
	return emitter, nil
}
