// Package mqtt is responsible for publishing (train) trip data to the MQTT broker
//
// HiveMQ is used to host a broker that can ease communication between
// both publishers and subscribers.
package mqtt

import (
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	env "github.com/joho/godotenv"
)

// this callback triggers when a message is received, it then prints the message (in the payload) and topic
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

// upon connection to the client, this is called
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

// this is called when the connection to the client is lost, it prints "Connection lost" and the corresponding error
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

type credentials struct {
	Url      string
	Username string
	Password string
}

func loadCredentials() (*credentials, error) {
	url := os.Getenv("MQTT_URL")
	username := os.Getenv("MQTT_USERNAME")
	password := os.Getenv("MQTT_PASSWORD")

	if url == "" || username == "" || password == "" {
		return nil, errors.New("missing MQTT environment variables")
	}

	return &credentials{
		Url:      url,
		Username: username,
		Password: password,
	}, nil
}

type Publisher struct {
	client mqtt.Client
}

// NewMqttPublisher creates a new mqttPublisher object
func NewMqttPublisher(port uint16, clientName string) (*Publisher, error) {

	var url, username, password string

	if len(credentialsPath) != 0 {
		credentials, err := loadCredentials()

		if err != nil {
			return nil, err
		}

		url, username, password = credentials.Url, credentials.Username, credentials.Password

	} else {
		return nil, errors.New("No credentials present")
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", url, port))
	opts.SetClientID(clientName)

	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	var publisher = Publisher{
		client: mqtt.NewClient(opts),
	}

	// return an error if the connection isn't successfull
	if token := publisher.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &publisher, nil
}

func (p *Publisher) Publish(data string, topic string) error {
	token := p.client.Publish(topic, 1, true, data)

	token.Wait()

	return token.Error()
}
