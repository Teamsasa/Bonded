package models

import "time"

type Event struct {
	EventID     string    `json:"event_id,omitempty" dynamodbav:"EventID"`
	Title       string    `json:"title" dynamodbav:"Title"`
	Description string    `json:"description" dynamodbav:"Description"`
	StartTime   time.Time `json:"start_time" dynamodbav:"StartTime"`
	EndTime     time.Time `json:"end_time" dynamodbav:"EndTime"`
	Location    string    `json:"location" dynamodbav:"Location"`
	AllDay      bool      `json:"all_day" dynamodbav:"AllDay"`
}
