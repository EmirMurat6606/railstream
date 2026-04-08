// Package mqtt is responsible for publishing (train) trip data to the MQTT broker
//
// HiveMQ is used to host a broker that can ease communication between
// both publishers and subscribers.
package mqtt

import (
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
	Username string
	Password string
}

func loadCredentials(path string) (*credentials, error) {
	// Read credentials from the path
	data, err := env.Read(path)

	if err != nil {
		return nil, err
	}

	var credentials = credentials{
		Username: data["MQTT_USERNAME"],
		Password: data["MQTT_PASSWORD"],
	}

	return &credentials, nil
}

type publisher struct {
	client mqtt.Client
}

// NewMqttPublisher creates a new mqttPublisher object
func NewMqttPublisher(broker string, port uint16, clientName string, credentialsPath string) (*publisher, error) {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker, port))
	opts.SetClientID(clientName)

	if len(credentialsPath) != 0 {
		credentials, err := loadCredentials(credentialsPath)

		if err != nil {
			return nil, err
		}

		opts.SetUsername(credentials.Username)
		opts.SetPassword(credentials.Password)
	}

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	var publisher = publisher{
		client: mqtt.NewClient(opts),
	}

	// return an error if the connection isn't successfull
	if token := publisher.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &publisher, nil
}

func (p *publisher) Publish(data byte, topic string) error {
	token := p.client.Publish(topic, 1, true, data)

	if token.Error() != nil {
		return token.Error()
	}

	return nil
}
