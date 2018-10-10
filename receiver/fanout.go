package main

import (
	"github.com/streadway/amqp"
)

func fanout(ch *amqp.Channel) <-chan amqp.Delivery {

	err := ch.ExchangeDeclare(
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

	/**
	* Consume returns a channel that immediately delivering queued messages
	* Consumers range over the channel to ensure all deliveries are recieved.
	* Unrecived deliveries will block all methods on same connection
	* All deliveries in AMQP must be acknowledged.
	* If the consumer is cancelled or the channel or connection is closed any unacknowledged deliveries will be requeued at the end of the same queue.
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
