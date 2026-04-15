package gtfs

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"sync"

	gtfs "github.com/jamespfennell/gtfs"
	env "github.com/joho/godotenv"
)

// Global WaitGroup for all go routines
var Waitgroup = sync.WaitGroup{}

// StaticFeed stores information from the static GTFS feed
type StaticFeed struct {
    Stops     []*Stop
    stopIndex map[string]*Stop // interne
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

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("bmc-partner-key", apiKey)

	resp, _ := client.Do(req)


	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	defer resp.Body.Close()

	// Use the gtfs library
	b, _ := io.ReadAll(resp.Body)
	staticData, err := gtfs.ParseStatic(b, gtfs.ParseStaticOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("The Belgian train has %d routes and %d stations\n", len(staticData.Routes), len(staticData.Stops))


	req, err = http.NewRequest("GET", "https://api-management-opendata-production.azure-api.net/api/gtfs/feed/nmbssncb/rt/trip-update/?format=protobuf", nil)
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("bmc-partner-key", apiKey)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	
	defer resp.Body.Close()


	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	realtimeData, err := gtfs.ParseRealtime(b, &gtfs.ParseRealtimeOptions{})
	if err != nil {
		return err
	}

	var bestTrip *gtfs.Trip
	var bestTime time.Time

	temseID := "gs:nmbssncb:8894672"
	berchemID := "gs:nmbssncb:8821121"

	for i := range realtimeData.Trips {
		trip := &realtimeData.Trips[i]

		var temseTime *time.Time
		var temseSeq, berchemSeq int

		for _, stu := range trip.StopTimeUpdates {

			if *stu.StopID == temseID {
				temseSeq = int(*stu.StopSequence)
				if stu.Departure != nil && stu.Departure.Time != nil {
					temseTime = stu.Departure.Time
				}
			}

			if *stu.StopID == berchemID {
				berchemSeq = int(*stu.StopSequence)
			}
		}

		// ✔ juiste richting check
		if temseTime != nil && temseSeq > 0 && berchemSeq > 0 && temseSeq < berchemSeq {

			if bestTrip == nil || temseTime.Before(bestTime) {
				bestTrip = trip
				bestTime = *temseTime
			}
		}
	}

	if bestTrip != nil {
    fmt.Println("Volgende trip gevonden om:", bestTime)
		for _, stu := range bestTrip.StopTimeUpdates {
			if *stu.StopID == temseID && stu.Departure != nil && stu.Departure.Delay != nil {
				fmt.Println("Delay in Temse:", *stu.Departure.Delay, "seconden")
			}
		}
	} else {
		fmt.Println("Geen trip gevonden")
	}
	return nil
}

// StartLoader executes a goroutine that periodically updates the static feed
//
// It loads the static feed immediately and then schedules updates every 24 hours at 05:30 Brussels time
func StartLoader() error {

	loc, err := time.LoadLocation("Europe/Brussels")

	if err != nil {
		return err
	}

	// Start goroutine for periodic fetching

	Waitgroup.Add(1)
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
