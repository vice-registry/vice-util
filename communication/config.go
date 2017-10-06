package communication

import (
	"log"

	"github.com/streadway/amqp"
)

var rabbitmqCredentials = struct {
	Location string
	Username string
	Password string
	Channel  *amqp.Channel
}{}

// SetRabbitmqCredentials set the login credentials to Couchbase cluster
func SetRabbitmqCredentials(location string, username string, password string) error {
	rabbitmqCredentials.Location = location
	rabbitmqCredentials.Username = username
	rabbitmqCredentials.Password = password

	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + location + "/")
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		return err
	}

	rabbitmqCredentials.Channel = channel

	return nil
}

// CloseConnection closes all open connections to RabbitMQ
func CloseConnection() {
	rabbitmqCredentials.Channel.Close()
}
