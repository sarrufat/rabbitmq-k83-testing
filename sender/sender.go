package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	hostName := "localhost"
	if len(os.Args) == 2 {
		hostName = os.Args[1]
	}
	connStr := fmt.Sprintf("amqp://guest:guest@%s:5672/", hostName)
	conn, err := amqp.Dial(connStr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"test01_queue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")
	forever := make(chan bool)

	body := bodyFrom()
	go func() {
		for true {
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         body,
				})
			failOnError(err, "Failed to publish a message")
			log.Printf(" [x] Sent")
			time.Sleep(50 * time.Millisecond)
		}
	}()
	<-forever
}

func bodyFrom() []byte {

	dat, err := ioutil.ReadFile("generated.json")
	if err != nil {
		log.Fatal(err)
	}
	return dat
}
