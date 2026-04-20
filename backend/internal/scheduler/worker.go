// Package scheduler provides utilities for running periodic tasks
// within a defined daily time window, with support for context-based
// cancellation and graceful shutdown.
//
// It is designed for long-running background workers such as data
// fetchers, ETL jobs, or streaming pipelines that must only execute
// during specific hours of the day.
package scheduler

import (
	"context"
	"log"
	"time"
)

// Task is a function that can be scheduled
type Task func(ctx context.Context) error

// Window defines a daily active time window
type Window struct {
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
	Location    *time.Location
}

// RunPeriodic runs a task every interval while inside the window.
// It blocks until context is cancelled or task returns a fatal error.
func RunPeriodic(
	ctx context.Context,
	interval time.Duration,
	window Window,
	task Task,
) error {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		now := time.Now().In(window.Location)

		start := time.Date(
			now.Year(), now.Month(), now.Day(),
			window.StartHour, window.StartMinute,
			0, 0, window.Location,
		)

		end := time.Date(
			now.Year(), now.Month(), now.Day(),
			window.EndHour, window.EndMinute,
			0, 0, window.Location,
		)

		if now.Before(start) {
			sleepUntil(ctx, start)
			continue
		}

		if now.After(end) {
			sleepUntil(ctx, start.Add(24*time.Hour))
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			if err := task(ctx); err != nil {
				log.Println("scheduled task error:", err)
			}
		}
	}
}

func sleepUntil(ctx context.Context, t time.Time) {
	d := time.Until(t)
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}