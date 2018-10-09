package main

import (
	"fmt"
	"log"
	"strings"
	"time"

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

	// define a queue
	q, err := ch.QueueDeclare(
		"Hello queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to create a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
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
			count := strings.Count(body, ".")
			time.Sleep(time.Duration(count) * time.Second)
			fmt.Println("done!")
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
