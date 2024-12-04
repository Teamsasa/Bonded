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

	item, err := dynamodbattribute.MarshalMap(calendar)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.dynamoDB.PutItemWithContext(ctx, input)
	return err
}
