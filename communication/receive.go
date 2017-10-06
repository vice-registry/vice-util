package communication

import (
	"log"

	"github.com/streadway/amqp"
)

// NewConsumer register a new consumer on a queue
func NewConsumer(queueName string) (<-chan amqp.Delivery, error) {

	queue, err := rabbitmqCredentials.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		return nil, err
	}

	err = rabbitmqCredentials.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		return nil, err
	}

	msgs, err := rabbitmqCredentials.Channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		return nil, err
	}

	return msgs, nil

}
