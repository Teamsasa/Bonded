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
	item["SortKey"] = &dynamodb.AttributeValue{S: aws.String("CALENDAR")}
	item["UserID"] = &dynamodb.AttributeValue{S: aws.String(calendar.OwnerUserID)}
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
	item["SortKey"] = &dynamodb.AttributeValue{S: aws.String("CALENDAR")}
	item["UserID"] = &dynamodb.AttributeValue{S: aws.String(calendar.OwnerUserID)}
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
			"SortKey":    {S: aws.String("CALENDAR")}, // SortKeyを追加
		},
	}
	_, err := r.dynamoDB.DeleteItemWithContext(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func (r *calendarRepository) FindByCalendarID(ctx context.Context, calendarID string) (*models.Calendar, error) {
	// カレンダー情報を取得　（カレンダーとイベント、ユーザー情報を取得。カレンダー情報だけにするべき？）
    input := &dynamodb.GetItemInput{
        TableName: aws.String(r.tableName),
        Key: map[string]*dynamodb.AttributeValue{
            "CalendarID": {S: aws.String(calendarID)},
            "SortKey":    {S: aws.String("CALENDAR")},
        },
    }
    result, err := r.dynamoDB.GetItemWithContext(ctx, input)
    if err != nil {
        return nil, err
    }
    if result.Item == nil {
        return nil, fmt.Errorf("calendar with CalendarID %s not found", calendarID)
    }
    var calendar models.Calendar
    err = dynamodbattribute.UnmarshalMap(result.Item, &calendar)
    if err != nil {
        return nil, err
    }

    // 関連するイベントを取得
    eventInput := &dynamodb.QueryInput{
        TableName: aws.String(r.tableName),
        KeyConditionExpression: aws.String("CalendarID = :cid AND begins_with(SortKey, :sk)"),
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":cid": {S: aws.String(calendarID)},
            ":sk":  {S: aws.String("EVENT#")},
        },
    }
    eventResult, err := r.dynamoDB.QueryWithContext(ctx, eventInput)
    if err != nil {
        return nil, err
    }

    var events []models.Event
    err = dynamodbattribute.UnmarshalListOfMaps(eventResult.Items, &events)
    if err != nil {
        return nil, err
    }

    calendar.Events = events

    // 関連するユーザーを取得
    userInput := &dynamodb.QueryInput{
        TableName: aws.String(r.tableName),
        KeyConditionExpression: aws.String("CalendarID = :cid AND begins_with(SortKey, :sk)"),
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":cid": {S: aws.String(calendarID)},
            ":sk":  {S: aws.String("USER#")},
        },
    }
    userResult, err := r.dynamoDB.QueryWithContext(ctx, userInput)
    if err != nil {
        return nil, err
    }
    var users []models.User
    err = dynamodbattribute.UnmarshalListOfMaps(userResult.Items, &users)
    if err != nil {
        return nil, err
    }
    calendar.Users = users

    return &calendar, nil
}

func (r *calendarRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Calendar, error) {
    // GSIを使用してユーザーが所属するカレンダーを取得
    input := &dynamodb.QueryInput{
        TableName: aws.String(r.tableName),
        IndexName: aws.String("UserID-index"),
        KeyConditionExpression: aws.String("UserID = :uid"),
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":uid": {S: aws.String(userID)},
        },
    }
    result, err := r.dynamoDB.QueryWithContext(ctx, input)
    if err != nil {
        return nil, err
    }

    // カレンダー情報を取得
    var calendars []*models.Calendar
    calendarIDSet := make(map[string]struct{})
    for _, item := range result.Items {
        calendarID := *item["CalendarID"].S
        if _, exists := calendarIDSet[calendarID]; !exists {
            calendar, err := r.FindByCalendarID(ctx, calendarID)
            if err != nil {
                return nil, err
            }
            calendars = append(calendars, calendar)
            calendarIDSet[calendarID] = struct{}{}
        }
    }

    return calendars, nil
}
