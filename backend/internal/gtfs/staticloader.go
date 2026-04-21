package gtfs

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	gtfs "github.com/jamespfennell/gtfs"
	env "github.com/joho/godotenv"
)

// data stores the environment variables
var data map[string]string


func init() {
	var err error
	data, err = env.Read(".env")
	if err != nil{
		panic(err)
	}
}

// StartLoader executes a goroutine that periodically updates the static feed
//
// It loads the static feed immediately and then schedules updates every 24 hours at 05:30 Brussels time
func StartLoader(ctx context.Context) error {

	loc, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		return err
	}

	if err := updateStatic(); err != nil {
		log.Println("Initial static update failed:", err)
	}

	for {
		now := time.Now().In(loc)

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

		timer := time.NewTimer(duration)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()

		case <-timer.C:
			if err := updateStatic(); err != nil {
				log.Println("Static update failed:", err)
			}
		}
	}
}

func updateStatic() error {
	requestURL := "https://api-management-opendata-production.azure-api.net/api/gtfs/feed/nmbssncb/static/"

	client := &http.Client{}

	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("bmc-partner-key", data["API_KEY_BEL_MOBILITY"])

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Use the gtfs library
	b, _ := io.ReadAll(resp.Body)
	staticData, err := gtfs.ParseStatic(b, gtfs.ParseStaticOptions{})
	if err != nil {
		return err
	}

	for _, stop := range staticData.Stops {
		if stop.Description == "NMBSSNCB  STATION" {
			switch stop.Name {

			case "Tamise", "Saint-Nicolas", "Puurs", "Malines":
				stationMu.Lock()
				s := stationRegistry[stop.Name]
				s.ExtID = extractNumericID(stop.Id)
				stationRegistry[stop.Name] = s  
				stationMu.Unlock()
			default:
				// nothing to do
			}
		}
	}

	return nil
}

func extractNumericID(fullID string) string {
	idx := strings.LastIndex(fullID, ":")
	if idx == -1 {
		return ""
	}

	s := fullID[idx+1:]

	// remove all after _
	if u := strings.Index(s, "_"); u != -1 {
		s = s[:u]
	}

	// remove S prefix if present
	s = strings.TrimPrefix(s, "S")

	return s
}
