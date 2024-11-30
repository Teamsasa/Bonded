package db

import (
	"context"
	"fmt"
	"os"

	"bonded/internal/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CalendarRepository interface {
	Save(ctx context.Context, calendar *models.Calendar) error
	Update(ctx context.Context, calendar *models.Calendar) error
	Delete(ctx context.Context, calendarID string) error
	FindByID(ctx context.Context, calendarID string) (*models.Calendar, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.Calendar, error)
}

type calendarRepository struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

func NewCalendarRepository() CalendarRepository {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000" // デフォルトエンドポイント
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String(endpoint),
	}))
	return &calendarRepository{
		dynamoDB:  dynamodb.New(sess),
		tableName: "Calendars",
	}
}

func (r *calendarRepository) Save(ctx context.Context, calendar *models.Calendar) error {
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

func (r *calendarRepository) Update(ctx context.Context, calendar *models.Calendar) error {
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

func (r *calendarRepository) Delete(ctx context.Context, calendarID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(calendarID)},
		},
	}
	_, err := r.dynamoDB.DeleteItemWithContext(ctx, input)
	return err
}

func (r *calendarRepository) FindByID(ctx context.Context, calendarID string) (*models.Calendar, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(calendarID)},
		},
	}
	result, err := r.dynamoDB.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, fmt.Errorf("Calendar with ID %s not found", calendarID)
	}
	var calendar models.Calendar
	err = dynamodbattribute.UnmarshalMap(result.Item, &calendar)
	return &calendar, err
}

func (r *calendarRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Calendar, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		IndexName: aws.String("UserID-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"UserID": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{S: aws.String(userID)},
				},
			},
		},
	}
	result, err := r.dynamoDB.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	var calendars []*models.Calendar
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &calendars)
	return calendars, err
}
