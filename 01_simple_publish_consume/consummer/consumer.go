package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// get flags from command line
	engine := flag.String("engine", "engine-X", "RabbitMQ username")
	flag.Parse()

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

	var quit = make(chan os.Signal, 1)
	// handle notify signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// consume messages in other goroutines
	go func() {
		for msg := range msgs {
			// false = only ack this message
			// true = ack this message and all prior unacknowledged
			msg.Ack(false)
			fmt.Println("Received a message: ", string(msg.Body), "from go routine", *engine)
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-quit

	fmt.Println(" [*] Closing Channel and Connection")
}
