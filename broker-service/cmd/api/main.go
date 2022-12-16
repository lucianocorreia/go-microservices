package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	webport = "80"
)

type Config struct {
	Rabbit *amqp091.Connection
}

func main() {
	// try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	fmt.Printf("starting web broker on port %s\n", webport)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp091.Connection, error) {
	var count int64
	var backOff = 1 * time.Second
	var conn *amqp091.Connection

	for {
		c, err := amqp091.Dial("amqp://guest:guest@rabbitmq")
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
