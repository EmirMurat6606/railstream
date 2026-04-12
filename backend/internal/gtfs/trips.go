package gtfs


type Stop struct {
    ID   string `json:"stop_id"`
    Name string `json:"stop_name"`
    StopTimes []StopTime `json:"stop_times"`
}

type Trip struct {
    ID string `json:"trip_id"`
    Date string `json:"date"`
}

type StopTime struct {
    Trip Trip `json:"trip"`
    ArrivalTime string `json:"arrival_time"`
    DepartureTime string `json:"departure_time"`
    StopSequence uint8 `json:"stop_sequence"`
}