package events

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,  // durable?
		false, //auto-deletion?
		false, //internal?
		false, //no-wait?
		nil,   // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",
		false, // durable
		false, // delete when used?
		true,  //exclusive
		false, //no wait?
		nil,   // args

	)
}
