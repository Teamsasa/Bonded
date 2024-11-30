package models

type Calendar struct {
	ID     string `json:"id" dynamodbav:"ID"`
	UserID string `json:"userId" dynamodbav:"UserID"`
	Name   string `json:"name" dynamodbav:"Name"`
}
