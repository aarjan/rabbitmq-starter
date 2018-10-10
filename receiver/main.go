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

	forever := make(chan struct{})

	// type of exchange to choose!
	msgs := topic(ch)
	// msgs := direct(ch)
	// msgs := fanout(ch)

	go func() {
		for m := range msgs {
			body := string(m.Body)
			fmt.Println("Received a message, ", body)
		}
	}()

	log.Printf(" [x] Waiting for logs. To exit press Ctrl+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, error: %+v", msg, err)
	}
}
