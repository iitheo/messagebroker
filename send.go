package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

//func failOnError(err error, msg string) {
//	if err != nil {
//		log.Fatalf("%s: %s", msg, err)
//	}
//}

func main() {
	//connect to RabbitMQ server
	//The connection abstracts the socket connection, and takes care of protocol version negotiation
	//and authentication and so on for us.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	//Next we create a channel, which is where most of the API for getting things done resides:
	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
	}

	defer ch.Close()

	//To send, we must declare a queue for us to send to; then we can publish a message to the queue:
	//Declaring a queue is idempotent -
	//it will only be created if it doesn't exist already.
	//The message content is a byte array, so you can encode whatever you like there.
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}


	body := "Hello World!"+time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST")
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)

	if err != nil {
		failOnError(err, "Failed to publish a message")
	}

}

//We also need a helper function to check the return value for each amqp call:
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}