package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// create connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "failed to start server")
	defer conn.Close()

	// create channel
	ch, err := conn.Channel()
	failOnError(err, "failed to create channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to create an exchange")

	// define a queue
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "failed to create a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)
	failOnError(err, "failed to bind exchange to queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to consume message")

	forever := make(chan struct{})

	go func() {
		for m := range msgs {
			body := string(m.Body)
			fmt.Println("Received a message, ", body)
			m.Ack(false)
		}
	}()

	<-forever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, error: %+v", msg, err)
	}
}
