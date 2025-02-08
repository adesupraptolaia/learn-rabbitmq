package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	var (
		username = "guest"
		password = "guest"
	)

	// connect to RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@localhost:5672/", username, password))
	if err != nil {
		log.Panicln("Failed to connect to RabbitMQ, err: ", err.Error())
	}
	defer conn.Close()

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Panicln("Failed to open a channel, err: ", err.Error())
	}
	defer ch.Close()

	// open queue
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicln("Failed to declare a queue, err: ", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// send 20 messages
	for i := 1; i <= 20; i++ {
		body := fmt.Sprintf("Hello World %d! %s", i, time.Now().Format(time.RFC3339Nano))

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			// "test no queue and mandatory=true", // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Panicf("\nFailed to publish a message i: %d, err: %s\n", i, err.Error())
		}
		log.Printf(" [x] Sent %s\n", body)
	}
}
