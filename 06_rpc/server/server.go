package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	var (
		username = "guest"
		password = "guest"
	)

	// handle notify signals
	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

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

	// create a queue
	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Panicln("Failed to declare a queue, err: ", err.Error())
	}

	// set Qos
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Panicln("Failed to set Qos, err: ", err.Error())
	}

	// listen to queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicln("Failed to register a consumer, err: ", err.Error())
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var err error

		for msg := range msgs {
			msgBody := string(msg.Body)
			fmt.Println("Received a message: ", msgBody)
			fmt.Println("Processing...")

			time.Sleep(10 * time.Second)

			response := fmt.Sprintf("response of msg %s and correlation id %s", msgBody, msg.CorrelationId)

			// send response
			err = ch.PublishWithContext(ctx,
				"",          // exchange
				msg.ReplyTo, // routing key
				false,       // mandatory
				false,       // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: msg.CorrelationId,
					Body:          []byte(response),
				})
			if err != nil {
				log.Panicln("Failed to publish a message, err: ", err.Error())
			}

			fmt.Println("success send response to correlation id: ", msg.CorrelationId)

			msg.Ack(false)
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-quit

	fmt.Println(" [*] Closing Channel and Connection")
}
