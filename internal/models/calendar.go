package models

type Calendar struct {
	CalendarID string  `json:"calendarId" dynamodbav:"CalendarID"`
	UserID     string  `json:"userId" dynamodbav:"UserID"`
	Name       string  `json:"name" dynamodbav:"Name"`
	Event      []Event `json:"event" dynamodbav:"Event"`
}
