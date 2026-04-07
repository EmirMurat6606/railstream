// Package mqtt is responsible for publishing (train) trip data to the MQTT broker
//
// HiveMQ is used to host a broker that can ease communication between
// both publishers and subscribers.
package mqtt


import(
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

type credentials struct{
	Username string
	Password string
}

func newCredentials(path string) *credentials{
	// Read credentials from the path
	return nil
}

type MqttPublisher struct{

	Broker string
	Credentials credentials

}


