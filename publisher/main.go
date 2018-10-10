package main

import (
	"log"
	"os"
	"strings"

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

	msg := bodyFrom(os.Args)
	route := topicLogLevel(os.Args)
	// route := logLevel(os.Args)

	// the same function can be used for `direct` as well as `topic` exchange
	direct(ch, "logs-topic", "topic", route, msg)

	// fanout(ch, msg)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, error: %+v", msg, err)
	}
}

func bodyFrom(args []string) string {
	if len(args) < 3 {
		return "hello"
	}
	return strings.Join(os.Args[2:], " ")
}

func logLevel(args []string) string {
	if len(args) < 2 {
		return "info"
	}
	return args[1]
}

// Lets us assume routing keys of logs will have two words `<facility>.<severity>`.
func topicLogLevel(args []string) string {
	if len(args) < 2 || args[1] == "" {
		return "anonymous.info"
	}
	return args[1]
}
