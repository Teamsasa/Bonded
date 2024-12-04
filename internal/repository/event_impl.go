package repository

import (
	"bonded/internal/models"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (r *eventRepository) CreateEvent(ctx context.Context, calendar *models.Calendar, event *models.Event) error {
	calendar.Events = append(calendar.Events, *event)

	item, err := dynamodbattribute.MarshalMap(event)
	if err != nil {
		return err
	}

	item["CalendarID"] = &dynamodb.AttributeValue{S: aws.String(calendar.CalendarID)}
	item["SortKey"] = &dynamodb.AttributeValue{S: aws.String("EVENT#" + event.EventID)}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.dynamoDB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	gsiItem := map[string]*dynamodb.AttributeValue{
		"CalendarID": {S: aws.String(calendar.CalendarID)},
		"SortKey":    {S: aws.String("CAL#" + calendar.CalendarID + "#" + event.EventID)},
		"UserID":     {S: aws.String(calendar.OwnerUserID)},
	}
	gsiInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      gsiItem,
	}

	_, err = r.dynamoDB.PutItemWithContext(ctx, gsiInput)
	return err
}
