package models

import "time"

type Event struct {
	ID          string
	Title       string
	Description string
	CalendarID  string
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	AllDay      bool
}
