/**
* In RPC, a client sends request message and the server responds with response message.
* In order to recieve the response we need to send a 'callback' queue address with the request.
* AMQP protocol predefines 14 properties that go with message.
* Some of the important ones are:
* persistent: Marks a message persistent.
* content_type: Describe the mime-type of encoding.
* reply_to: Commonly used to name the callback queue.
* correlation_id: Used to correlate RPC responses with requests.

* In general, our RPC will work like this:
* For a RPC request, the client sends a message with two properties:
* 	1. reply_to: which is the callback queue
*	2. correlation_id: unique value for each request
* The request is sent to `rpc_queue` queue.
* The RPC worker/server is waiting for the request on that queue. It checks the `correlation_id` property.
* If matched, it returns the response on the callback queue.
 */

package main

import (
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, error: %+v", msg, err)
	}
}

// We use a simple fibonacci generator as method that processes the input and returns a result
func fib(n int) int {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		x, y = y, x+y
	}
	return x
}

func main() {
	// Create connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "failed to start server")
	defer conn.Close()

	// Create channel
	ch, err := conn.Channel()
	failOnError(err, "failed to create channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "failed to set Qos")

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // arguments
	)
	failOnError(err, "failed to establish consumer")

	forever := make(chan struct{})

	go func() {
		for m := range msgs {
			n, err := strconv.Atoi(string(m.Body))
			failOnError(err, "invalid integer")

			response := fib(n)
			log.Printf(" [.] fib(%d): %d", n, response)

			err = ch.Publish(
				"",        // exchange
				m.ReplyTo, // routing key
				true,      // mandatory
				false,     // immediate
				amqp.Publishing{
					CorrelationId: m.CorrelationId,
					ContentType:   "text/plain",
					Body:          []byte(strconv.Itoa(response)),
				},
			)
			failOnError(err, "failed to publish message")

			m.Ack(false)
		}
	}()

	<-forever
}
