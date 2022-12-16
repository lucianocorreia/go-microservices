package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rConn.Close()

	// start listening to messages
	log.Println("listening and cosuming RabbitMQ messages")

	// create a consumer
	consumer, err := event.NewConsumer(rConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and cosume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var count int64
	var backOff = 1 * time.Second
	var conn *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("rabbitmq is not ready yet...")
			count++
		} else {
			conn = c
			break
		}

		if count > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return conn, nil
}
