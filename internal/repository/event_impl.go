package repository

import (
	"bonded/internal/models"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (r *eventRepository) Create(ctx context.Context, event *models.Event) error {
	item, err := dynamodbattribute.MarshalMap(event)
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

func (r *eventRepository) FindByEventID(ctx context.Context, eventID string) (*models.Event, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"EventID": {
				S: aws.String(eventID),
			},
		},
	}
	result, err := r.dynamoDB.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	var event models.Event
	err = dynamodbattribute.UnmarshalMap(result.Item, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
