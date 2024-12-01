package models

import "time"

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string
	CalendarID  string `json:"calendar_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string
	AllDay      bool
}
