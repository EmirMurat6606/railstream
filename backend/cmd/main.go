// Package main is the entry point of the backend program
package main

import (
	pub "github.com/EmirMurat6606/railstream/internal/mqtt"
)

func check(err error){
	if (err != nil){
		panic(err)
	}
}

func main() {
	publisher, err := pub.NewMqttPublisher("49970051230a4dacb5b34fe2e2447647.s1.eu.hivemq.cloud", 8883, "railstreamer", ".env")
	check(err)
	publisher.Print()
	var loop bool = true

	for (loop) {
	
	}
}
