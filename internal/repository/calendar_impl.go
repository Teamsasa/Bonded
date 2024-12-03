package repository

import (
	"bonded/internal/models"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (r *calendarRepository) Create(ctx context.Context, calendar *models.Calendar) error {
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

func (r *calendarRepository) Edit(ctx context.Context, calendar *models.Calendar) error {
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
			"CalendarID": {S: aws.String(calendarID)},
		},
	}
	result, err := r.dynamoDB.DeleteItemWithContext(ctx, input)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("calendar with calendarID %s not found", calendarID)
	}
	return nil
}

func (r *calendarRepository) FindByCalendarID(ctx context.Context, calendarID string) (*models.Calendar, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"CalendarID": {S: aws.String(calendarID)},
		},
	}
	result, err := r.dynamoDB.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, fmt.Errorf("calendar with calendarID %s not found", calendarID)
	}
	var calendar models.Calendar
	err = dynamodbattribute.UnmarshalMap(result.Item, &calendar)
	return &calendar, err
}

func (r *calendarRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Calendar, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditionExpression: aws.String("UserID = :uid"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":uid": {S: aws.String(userID)},
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
