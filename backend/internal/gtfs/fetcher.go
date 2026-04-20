// Package gtfs is responsible for fetching realtime train data from the Belgian Mobility API
package gtfs


func getStationID(name string) string {
	stationMu.RLock()
	defer stationMu.RUnlock()
	return stationToId[name]
}
