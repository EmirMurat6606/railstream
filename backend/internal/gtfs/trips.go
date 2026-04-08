package gtfs

import (
	"time"
)

// Trip stores information about the relevant train trips
type Trip struct {
    TripID    string
    RouteID   string
    ServiceID string
}

// StopTime stores the arrival and departure time of a stop (station)
type StopTime struct {
    StopID       string
    ArrivalTime  time.Duration
    DepartureTime time.Duration
    Sequence     int
}

type Stop struct {
    ID   string
    Name string
}

type Service struct {
    ID        string
    Monday    bool
    Tuesday   bool
    Wednesday bool
    Thursday  bool
    Friday    bool
    Saturday  bool
    Sunday    bool
}