package communication

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/vice-registry/vice-util/models"
)

// SendMessage publishes a message to a queue (creates queue if not exist)
func SendMessage(queueName string, message string) error {
	importQueue, err := rabbitmqCredentials.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Failed to declare a queue %s: %s", queueName, err)
		return err
	}

	body := message
	err = rabbitmqCredentials.Channel.Publish(
		"",               // exchange
		importQueue.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		log.Printf("Failed to publish a message to queue %s: %s", queueName, err)
		return err
	}
	return nil
}

// NewImportAction creates a new action for import request in the message queue
func NewImportAction(image *models.Image) error {
	return SendMessage("import", image.ID)
}

// NewExportAction creates a new action for export request in the message queue
func NewExportAction(deployment *models.Deployment) error {
	return SendMessage("export", deployment.ID)
}
