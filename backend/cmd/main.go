package main

import (
	"context"
	"os"
	"os/signal"

	"sync"
	"syscall"

	gtfs "github.com/EmirMurat6606/railstream/internal/gtfs"
	"github.com/EmirMurat6606/railstream/internal/mqtt"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// waitGroup for all go routines
var waitGroup = sync.WaitGroup{}

func main() {
	// This is the actual main function
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	pub, err := mqtt.NewMqttPublisher(
		".env",
		8883,
		"railstream_ec2_publisher",
	)

	check(err)

	waitGroup.Go(func() {
		err := gtfs.StartLoader(ctx)
		check(err)
	})

	waitGroup.Go(func() {
		err := gtfs.StartRealtimeFetcher(ctx, pub)
		check(err)
	})

	waitGroup.Wait()
}
