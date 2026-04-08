package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

// parseStaticGTFS extracts static feed data
//
// This function extracts the gtfs.zip, gets all relevant data and returns a staticFeed object or an error
func parseStaticGTFS() (*StaticFeed, error) {

	gtfs, err := zip.OpenReader("./gtfs/gtfs.zip")

	if err != nil {
		return nil, err
	}

	defer gtfs.Close()

	// Construct the static feed
	feed := &StaticFeed{
		Stops:     make(map[string]*Stop),
		Trips:     make(map[string]*Trip),
		StopTimes: make(map[string][]*StopTime),
		Services:  make(map[string]*Service),
	}

	for _, f := range gtfs.File {
		switch f.Name {

		case "stops.txt":
			if err := parseStops(f, feed); err != nil {
				return nil, err
			}
		case "trips.txt":
			if err := parseTrips(f, feed); err != nil {
				return nil, err
			}
		case "stop_times.txt":
			if err := parseStopTimes(f, feed); err != nil {
				return nil, err
			}
		case "calendar.txt":
			if err := parseCalendar(f, feed); err != nil {
				return nil, err
			}
		default:
			// Nothing to do
		}
	}

	return nil, nil
}

func parseStops(f *zip.File, feed *StaticFeed) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}

	defer rc.Close()

	reader := csv.NewReader(rc)

	header, err := reader.Read()
	fmt.Print(header)
	if err != nil {
		return err
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
