package main

import (
	"github.com/streadway/amqp"
)

func runQueue(conn *amqp.Connection, config Configuration) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		config.RunQueue, // name
		true,            // durable
		false,           // delete when usused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	for d := range msgs {
		runBenchmark(d.Body)
	}
	<-forever

}
