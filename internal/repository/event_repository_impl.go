package repository

import (
	"context"
	"bonded/internal/infra/db"
	"bonded/internal/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type eventRepository struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	FindByEventID(ctx context.Context, eventID string) (*models.Event, error)
	// ...必要に応じて他のメソッド...
}

func EventRepositoryRequest(dynamoClient *db.DynamoDBClient) EventRepository {
	return &eventRepository{
		dynamoDB:  dynamoClient.Client,
		tableName: "Events",
	}
}

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
			"ID": {
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