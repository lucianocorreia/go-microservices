package event

import "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp091.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareRandomQueue(ch *amqp091.Channel) (amqp091.Queue, error) {
	return ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
}
