package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, error: %v", msg, err)
	}
}

func randomString() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func fibRPC(n int) int {
	// create connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "failed to start server")
	defer conn.Close()

	// create channel
	ch, err := conn.Channel()
	failOnError(err, "failed to create channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // queue name
		false, // durable
		false, // auto del
		true,  // exclusive
		false, // no wait
		nil,   // args
	)
	failOnError(err, "failed to create a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // arguments
	)
	failOnError(err, "failed to register a consumer")

	corrID := randomString()
	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			CorrelationId: corrID,
			ReplyTo:       q.Name,
			ContentType:   "text/plain",
			Body:          []byte(strconv.Itoa(n)),
		},
	)
	failOnError(err, "failed to publish")

	res := 0
	for m := range msgs {
		if corrID == m.CorrelationId {
			res, err = strconv.Atoi(string(m.Body))
			failOnError(err, "failed to convert to integer")
			break
		}
	}
	return res
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	n := bodyForm(os.Args)
	log.Printf(" [x] Requesting fib(%d)", n)

	start := time.Now()

	// use Ticker for timeouts
	go func() {
		select {
		case <-time.Tick(6 * time.Millisecond):
			failOnError(fmt.Errorf("timeout"), "could not process the request")
		}
	}()

	val := fibRPC(n)
	log.Printf(" [.] Got %d, time:%fms", val, time.Since(start).Seconds()*1000)
}

func bodyForm(args []string) int {
	s := "30"
	if len(args) == 2 && args[1] != "" {
		s = args[1]
	}
	n, err := strconv.Atoi(s)
	failOnError(err, "failed to convert to integer")
	return n
}
