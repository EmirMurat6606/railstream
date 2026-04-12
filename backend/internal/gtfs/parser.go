package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"strconv"
	"strings"
)

// parseStaticGTFS extracts static feed data
//
// This function extracts the gtfs.zip, gets all relevant data and returns a staticFeed object or an error
func ParseStaticGTFS() (*StaticFeed, error) {

	gtfs, err := zip.OpenReader("gtfs.zip")

	if err != nil {
		return nil, err
	}

	defer gtfs.Close()

	// Create a new StaticFeed object to store the parsed data
	feed := &StaticFeed{
		Stops:     []*Stop{},
		stopIndex: make(map[string]*Stop),
	}

	for _, f := range gtfs.File {
		switch f.Name {

		case "gtfs/stops.txt":
			for _, f2 := range gtfs.File {
				switch f2.Name {
				case "gtfs/stop_times.txt":
					if err := parseStops(f, f2, feed); err != nil {
						return nil, err
					}
				}
			}

		case "gtfs/trips.txt":
			// code
		default:
			// Nothing to do
		}
	}

	return feed, nil
}

func parseStops(stops *zip.File, stopTimes *zip.File, feed *StaticFeed) error {
	rc, err := stops.Open()

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

	for index, field := range header {
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

	// First, get all the stop IDs

	for _, row := range stopData {
		switch row[stopNameIndex] {
		case "Tamise", "Anvers-Berchem", "Puurs", "Saint-Nicolas", "Malines":
			if row[stopDescIndex] == "NMBSSNCB  STATION" {
				stop := &Stop{
					ID:        normalizeStopID(row[stopIDIndex]),
					Name:      row[stopNameIndex],
					StopTimes: []StopTime{},
				}

				feed.Stops = append(feed.Stops, stop)
				feed.stopIndex[stop.ID] = stop // internal pointer mapping for O(1) lookup
			}
		default:
			// nothing to do
		}
	}

	// Parse the stop times and sequences and update the feed
	if err = parseStopTimes(stopTimes, feed); err != nil {
		return err
	}

	return nil
}

func parseStopTimes(f *zip.File, feed *StaticFeed) error {

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

	var at, dt, sId, sSeq, tID uint8

	for index, field := range header {
		switch field {
		case "arrival_time":
			at = uint8(index)
		case "departure_time":
			dt = uint8(index)
		case "stop_id":
			sId = uint8(index)
		case "stop_sequence":
			sSeq = uint8(index)
		case "trip_id":
			tID = uint8(index)
		default:
			// Nothing to do
		}
	}

	data, err := reader.ReadAll()

	if err != nil {
		return err
	}

	for _, row := range data {
		stop, ok := feed.stopIndex[normalizeStopID(row[sId])]
		if !ok {
			continue
		}

		stopSeq, err := strconv.Atoi(row[sSeq])

		if err != nil {
			return err
		}

		tripId := row[tID]
		date := trimDate(tripId)

		stop.StopTimes = append(stop.StopTimes, StopTime{
			Trip: Trip{
				ID: tripId,
				Date: date,
			},
			ArrivalTime:   row[at],
			DepartureTime: row[dt],
			StopSequence:  uint8(stopSeq),
		})
	}

	return nil
}

func parseCalendar(f *zip.File, feed *StaticFeed) error {
	return nil
}

func normalizeStopID(id string) string {
	parts := strings.Split(id, ":")
	last := parts[len(parts)-1]

	// verwijder enkel de 'S' prefix
	last = strings.TrimPrefix(last, "S")

	return parts[0] + ":" + parts[1] + ":" + last
}

func trimDate(tripId string) string {
	parts := strings.Split(tripId, ":")
	var idx uint8
	if len(parts[len(parts) - 1]) <= 2 {
		idx = uint8(len(parts) - 2)
	}else{
		idx = uint8(len(parts) - 1)
	}
	return parts[idx]
}
