/**
* Direct Exchange directs the message to particular queues.
* The routing algorithm behind the `direct` exchange directs the messages to queues whose `binding key` exactly matches the `routing key` of the message
 */
package main

import (
	"log"

	"github.com/streadway/amqp"
)

func direct(ch *amqp.Channel, exchangeName, exchangeType, route, msg string) {

	/**
	* Declare a `direct` exchange
	 */
	err := ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // kind
		true,         // durable
		false,        // autoDelete
		false,        // internal
		false,        // noWait
		nil,          // args
	)
	failOnError(err, "failed to create an exchange")

	/*
	* Publish message to the channel
	* If the `exchange` param is not specified, messages are routed to the queue with the name specified in the `routing_key` param, if exists.
	* The `routing key` is used to send message to specific routes, i.e queues.
	 */
	err = ch.Publish(
		exchangeName, // exchange
		route,        // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // required?
			ContentType:  "text/plain",
			Body:         []byte(msg),
		}, // publishing protocol
	)
	failOnError(err, "failed to publish message")

	log.Printf(" [x] Send msg: %q to route: %s using exchange type: %s on exchange: %s", msg, route, exchangeType, exchangeName)
}
