package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

// parseStaticGTFS extracts static feed data
//
// This function extracts the gtfs.zip, gets all relevant data and returns a staticFeed object or an error
func ParseStaticGTFS() (*StaticFeed, error) {

	gtfs, err := zip.OpenReader("gtfs.zip")
	fmt.Println("1")


	if err != nil {
		return nil, err
	}

	defer gtfs.Close()

	// Create a new StaticFeed object to store the parsed data
	feed := &StaticFeed{
		Stops: []Stop{},
	}

	for _, f := range gtfs.File {
		switch f.Name {

		case "gtfs/stops.txt":
			fmt.Println("2")
			if err := parseStops(f, feed); err != nil {
				return nil, err
			}
		default:
			// Nothing to do
		}
	}

	return feed, nil
}

func parseStops(f *zip.File, feed *StaticFeed) error {
	rc, err := f.Open()

	if err != nil {
		return err
	}
	defer rc.Close()

	reader := csv.NewReader(rc)

	header, err := reader.Read()

	if err != nil {
		return err
	}

	var stopIDIndex, stopNameIndex, stopDescIndex uint8

	for index, field := range header{
		switch field {
		case "stop_id":
			stopIDIndex = uint8(index)
		case "stop_name":
			stopNameIndex = uint8(index)
		case "stop_desc":
			stopDescIndex = uint8(index)
		default:
			// Nothing to do
		}
	}

	stopData, err := reader.ReadAll()

	if err != nil {
		return err
	}

	for _, row := range stopData{
		switch row[stopNameIndex] {
		case "Tamise", "Anvers-Berchem", "Puurs", "Saint-Nicolas", "Malines":
			if row[stopDescIndex] == "NMBSSNCB  STATION" {
				feed.Stops = append(feed.Stops, Stop{
					ID: row[stopIDIndex],
					Name: row[stopNameIndex],
				})
			}
		default:
			// nothing to do
		}
	}

	return nil
}

func parseTrips(f *zip.File, feed *StaticFeed) error {
	return nil
}
func parseStopTimes(f *zip.File, feed *StaticFeed) error {
	return nil
}

func parseCalendar(f *zip.File, feed *StaticFeed) error {
	return nil
}
