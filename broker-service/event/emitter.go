package event

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp091.Connection
}

func NewEventEmitter(conn amqp091.Connection) (Emitter, error) {
	emitter := Emitter{
		conn: &conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}

func (e *Emitter) setup() error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch)
}

func (e *Emitter) Push(event string, severity string) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("pushing to channel...")

	err = ch.Publish("logs_topic", severity, false, false, amqp091.Publishing{
		ContentType: "text/plainl",
		Body:        []byte(event),
	})
	if err != nil {
		return err
	}

	return nil
}
