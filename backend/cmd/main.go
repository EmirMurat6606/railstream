// Package main is the entry point of the backend program
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	mqtt "github.com/EmirMurat6606/railstream/internal/mqtt"
	gtfs "github.com/EmirMurat6606/railstream/internal/gtfs"
)

// waitGroup for all go routines
var waitgroup = sync.WaitGroup{}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	pub, err := mqtt.NewMqttPublisher(
		"49970051230a4dacb5b34fe2e2447647.s1.eu.hivemq.cloud",
		8883,
		"railstream_cluster",
		".env",
	)
	check(err)

	waitgroup.Go(func() {
		err := gtfs.StartLoader(ctx)
		check(err)
	})

	waitgroup.Go(func() {
		err := gtfs.StartRealtimeFeed(ctx, pub)
		check(err)
	})

	waitgroup.Wait()
}
