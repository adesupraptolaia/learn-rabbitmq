package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	workerName := "worker-X"
	if len(os.Args) > 0 {
		workerName = os.Args[1]
	}
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
		"worker_queue_name", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Panicln("Failed to declare a queue, err: ", err.Error())
	}

	// make fair distribution, set prefetch count = 1
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Panicln("Failed to set Qos, err: ", err.Error())
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

	// regx := regexp.MustCompile(`\d+`)

	// consume messages in other goroutines
	go func() {
		for msg := range msgs {
			// false = only ack this message
			// true = ack this message and all prior unacknowledged
			msgBody := string(msg.Body)
			fmt.Println("Received a message: ", msgBody, "from go routine", workerName)

			num, err := strconv.Atoi(regexp.MustCompile(`\d+`).FindString(msgBody))
			if err != nil {
				log.Panicln("Failed to convert string to int, err: ", err.Error())
			}

			if num%2 == 0 {
				time.Sleep(5 * time.Second)
			}

			msg.Ack(false)
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-quit

	fmt.Println(" [*] Closing Channel and Connection")

}
