package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
		"broadcast",         // exchange name
		amqp.ExchangeFanout, // exchange type
		false,               // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Panicf("Failed to declare an exchange, err: %s\n", err.Error())
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Panicln("Failed to declare a queue, err: ", err.Error())
	}

	// bind the queue to exchange "broadcast"
	ch.QueueBind(
		q.Name,      // queue name
		"",          // routing key
		"broadcast", // exchange
		false,       // no-wait
		nil,         // arguments
	)

	// create a consumer
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack = false, I want to ack manually
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicln("Failed to register a consumer, err: ", err.Error())
	}

	// handle quit signals
	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// consume messages in other goroutines
	go func() {
		for msg := range msgs {
			// false = only ack this message
			// true = ack this message and all prior unacknowledged
			msgBody := string(msg.Body)
			fmt.Println("Received a message: ", msgBody)

			msg.Ack(false)
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-quit

	fmt.Println(" [*] Closing Channel and Connection")

}
