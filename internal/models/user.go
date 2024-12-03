package models

type User struct {
	UserId string `json:"userId" dynamodbav:"UserID"`
	Name   string `json:"name" dynamodbav:"Name"`
}
