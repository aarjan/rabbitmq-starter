/**
* Recieve from a `topic` exchange
* Here, we create a new binding for different log level `<facility>.<severity>`
 */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func topic(ch *amqp.Channel) <-chan amqp.Delivery {

	err := ch.ExchangeDeclare(
		"logs-topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "failed to create an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // auto-deleted
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "failed to create a queue")

	if len(os.Args) < 2 {
		fmt.Printf("Usage %s: [binding key]", os.Args[0])
		os.Exit(0)
	}

	for _, key := range os.Args[1:] {

		log.Printf("Binding queue: %s to exchange: %s using routing key: %s",
			q.Name, "logs-topic", key)

		err = ch.QueueBind(
			q.Name,       // queue name
			key,          // bind key
			"logs-topic", // exchange name
			false,        // noWait
			nil,          // args
		)
		failOnError(err, "failed to bind the queue to the exchange")
	}

	/**

	 */
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer name
		true,   // auto ack
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	failOnError(err, "failed to register a consumer")

	return msgs
}
