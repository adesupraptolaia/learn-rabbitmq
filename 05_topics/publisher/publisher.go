package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	err = ch.ExchangeDeclare(
		"exchange_topic_name", // exchange name
		amqp.ExchangeTopic,    // exchange type
		false,                 // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		log.Panicf("Failed to declare an exchange, err: %s\n", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	routingKey := "routing_key"
	if len(os.Args) > 1 {
		routingKey = os.Args[1]
	}

	body := fmt.Sprintf("Hello World: routing_key %s", routingKey)

	err = ch.PublishWithContext(ctx,
		"exchange_topic_name", // exchange
		routingKey,            // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		log.Panicf("\nFailed to publish a message_body: %s err: %s\n", body, err.Error())
	}
	log.Printf(" [x] Sent %s\n", body)
}
