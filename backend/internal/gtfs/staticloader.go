package gtfs

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	env "github.com/joho/godotenv"
)

type StaticFeed struct {
	Dummy string
}

func updateStatic() error {
	requestURL := "https://api-management-opendata-production.azure-api.net/api/gtfs/feed/nmbssncb/static/"

	client := &http.Client{}

	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		return err
	}

	data, err := env.Read(".env")

	if err != nil {
		return err
	}

	apiKey := data["API_KEY_BEL_MOBILITY"]

	req.Header.Add("Cache-Contorl", "no-cache")
	req.Header.Add("bmc-partner-key", apiKey)

	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the zip file
	out, err := os.Create("gtfs.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// StartLoader executes a goroutine that periodically updates the static feed
func StartLoader() error {

	loc, err := time.LoadLocation("Europe/Brussels")

	if err != nil {
		return err
	}

	// Start goroutine for periodic fetching
	go func() {

		// Load the static data initially
		if err := updateStatic(); err != nil {
			log.Println("Initial static update failed:", err)
		}

		for {
			now := time.Now().In(loc)

			// Vandaag 05:30
			nextRun := time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				5, 30, 0, 0,
				loc,
			)

			if now.After(nextRun) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			duration := time.Until(nextRun)
			log.Println("Next static GTFS update at:", nextRun)

			time.Sleep(duration)

			err := updateStatic()
			if err != nil {
				log.Println("Static update failed:", err)
			}
		}
	}()

	return nil
}
