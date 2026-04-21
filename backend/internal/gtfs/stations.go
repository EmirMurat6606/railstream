package gtfs

import (
	"net/url"
	"sync"
)

type station struct {
	Name      string
	ExtID     string
	reference string // raw Hafas reference string (not URL-encoded)
}

// stationRegistry contains all stations with their Hafas-reference
//
// The key values are in French because the static gtfs data provides French station names
// which makes implementation slightly easier
var (
	stationRegistry = map[string]station{
		"Tamise": {
			Name:      "Temse",
			ExtID:     "8894672",
			reference: "A=1@O=Temse@X=4221352@Y=51126086@U=80@L=8894672@B=1@p=1776635490@",
		},
		"Saint-Nicolas": {
			Name:      "Sint-Niklaas",
			ExtID:     "8894508",
			reference: "A=1@O=Sint-Niklaas@X=4142966@Y=51171472@U=80@L=8894508@B=1@p=1776549195@",
		},
		"Puurs": {
			Name:      "Puurs",
			ExtID:     "8822715",
			reference: "A=1@O=Puurs@X=4282704@Y=51077220@U=80@L=8822715@B=1@p=1776549195@",
		},
		"Malines": {
			Name:      "Mechelen",
			ExtID:     "8822004",
			reference: "A=1@O=Mechelen@X=4482786@Y=51017649@U=80@L=8822004@B=1@p=1776549195@",
		},
	}

	stationMu sync.RWMutex
)

// getStation gets a station from the register
func getStation(name string) (station, bool) {
	s, ok := stationRegistry[name]
	return s, ok
}

// encodedReference returns URL-encoded Hafas reference
func (s station) encodedReference() string {
	return url.QueryEscape(s.reference)
}
