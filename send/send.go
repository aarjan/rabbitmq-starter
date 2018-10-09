package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	// Create connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "failed to start server")
	defer conn.Close()

	// Create channel
	ch, err := conn.Channel()
	failOnError(err, "failed to create channel")
	defer ch.Close()

	/**
	* Exchanges receives messages from the producer and pushes it to queues
	* Exchanges must know what to do with the messages it recieves, based on the `exchange type`.
	* fanout: publishes to all queues
	* topic:
	* headers:
	* direct:
	 */
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // kind
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
	failOnError(err, "failed to create an exchange")

	/**
	* A Queue is a buffer used to store messages.
	* Giving queue a name is important when you want to share messages between a publisher and consumer, but, if you want to publish messages to all the queues, `name` is not mandatory.
	* Here, RabbitMQ auto generates random queue name for this instance
	 */
	q, err := ch.QueueDeclare(
		"",    // queue name
		true,  // durable
		false, // delete when consumer disconnects
		true,  // exclusive
		false, // noWait
		nil,   // args
	)

	/**
	* Now, that we have created a `fanout` exchange and a queue, we need to bind the exchange to our queues.
	 */
	err = ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)
	failOnError(err, "failed to bind exchange to the queue")

	// Add Qos on channel
	failOnError(err, "failed to create a queue")
	err = ch.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	)
	failOnError(err, "failed to add Qos")

	body := os.Args[1]

	/*
	* Publish message to the channel
	* If the `exchange` param is not specified, messages are routed to the queue with the name specified in the `routing_key` param, if exists.
	* The `routing key` is used to send message to specific routes, i.e queues, but, since, we are using `fan out` exchange model, we emit messages to all our keys.
	*/
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		}, // publishing protocol
	)
	failOnError(err, "failed to publish message")

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, error: %+v", msg, err)
	}
}
