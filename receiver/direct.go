/**
* Here, we create a new binding for different log level i.e. `info`, `warn`, `error`
 */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func direct(ch *amqp.Channel) <-chan amqp.Delivery {

	// Why `durable` true in exchange but not in queue?
	err := ch.ExchangeDeclare(
		"logs-direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "failed to create an exchange")

	// Why `exclusive` true in queue but not in `exchange`?
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
		fmt.Printf("Usage %s: [info] [warn] [error]", os.Args[0])
		os.Exit(0)
	}

	for _, key := range os.Args[1:] {

		log.Printf("Binding queue: %s to exchange: %s using routing key: %s",
			q.Name, "logs-direct", key)

		err = ch.QueueBind(
			q.Name,        // queue name
			key,           // bind key
			"logs-direct", // exchange name
			false,         // noWait
			nil,           // args
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
	failOnError(err, "failed to consume message")

	return msgs
}
