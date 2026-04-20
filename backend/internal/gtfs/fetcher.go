// Package gtfs is responsible for fetching realtime train data from the Belgian Mobility API
package gtfs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	scheduler "github.com/EmirMurat6606/railstream/internal/scheduler"
	mqtt "github.com/EmirMurat6606/railstream/internal/mqtt"
	gtfs "github.com/jamespfennell/gtfs"
)

// StartRealtimeFeed fetches the Belgian Mobility API periodically for live trip updates
//
// text...
func StartRealtimeFeed(ctx context.Context, publisher *mqtt.Publisher) error {

	loc, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		return err
	}

	window := scheduler.Window{
		StartHour:   6,
		StartMinute: 0,
		EndHour:     12,
		EndMinute:   30,
		Location:    loc,
	}

	return scheduler.RunPeriodic(
		ctx,
		30*time.Second,
		window,
		func(ctx context.Context) error {
			return fetchRealTimeData()
		},
	)
}

func fetchRealTimeData(publisher *mqtt.Publisher) error {

	requestUrl := "https://api-management-opendata-production.azure-api.net/api/gtfs/feed/nmbssncb/rt/trip-update/?format=protobuf"

	client := &http.Client{}

	req, err := http.NewRequest("GET", requestUrl, nil)
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("bmc-partner-key", data["API_KEY_BEL_MOBILITY"])

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	realtimeData, err := gtfs.ParseRealtime(b, &gtfs.ParseRealtimeOptions{})
	if err != nil {
		return err
	}

	var tamiseSeq, stnicolasSeq, puursSeq, malinesSeq uint16

	for _, i := range realtimeData.Trips {
		tamiseSeq, stnicolasSeq, puursSeq, malinesSeq = 0, 0, 0, 0
		for _, j := range i.StopTimeUpdates {

			id := extractNumericID(*j.StopID)

			switch id {

			case getStationID("Tamise"):
				tamiseSeq = uint16(*j.StopSequence)
			case getStationID("Saint-Nicolas"):
				stnicolasSeq = uint16(*j.StopSequence)
			case getStationID("Puurs"):
				puursSeq = uint16(*j.StopSequence)
			case getStationID("Malines"):
				malinesSeq = uint16(*j.StopSequence)
			default:
				// nothing to do
			}
		}

		if tamiseSeq > 0 {
			if stnicolasSeq > 0 && tamiseSeq < stnicolasSeq {
				fmt.Print("On trip ")
				fmt.Println(i.ID)
				fmt.Println("Found: Temse -> St-Niklaas")
				fmt.Println(i.StopTimeUpdates[len(i.StopTimeUpdates)-1].Arrival.Time)
				fmt.Println("---------------------------------")
			}
			if puursSeq > 0 && tamiseSeq < puursSeq {
				fmt.Print("On trip ")
				fmt.Println(i.ID)
				fmt.Print("     ")
				fmt.Println("Found: Temse -> Puurs")
				fmt.Println("---------------------------------")

			}
			if malinesSeq > 0 && tamiseSeq < malinesSeq {
				fmt.Print("On trip ")
				fmt.Println(i.ID)
				fmt.Println("Found: Temse -> Mechelen")
				fmt.Println("---------------------------------")

			}
		}
	}

	return nil

}

func getStationID(name string) string {
	stationMu.RLock()
	defer stationMu.RUnlock()
	return stationToId[name]
}
