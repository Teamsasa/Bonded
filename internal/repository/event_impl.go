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

func (r *eventRepository) FindEvents(ctx context.Context, calendarID string) ([]*models.Event, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("CalendarID = :calendarID AND begins_with(SortKey, :sortPrefix)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":calendarID": {S: aws.String(calendarID)},
			":sortPrefix": {S: aws.String("EVENT#")},
		},
	}

	result, err := r.dynamoDB.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	events := make([]*models.Event, 0, len(result.Items))
	for _, item := range result.Items {
		var event models.Event
		err = dynamodbattribute.UnmarshalMap(item, &event)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *eventRepository) EventExists(ctx context.Context, calendarID string, eventID string) bool {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"CalendarID": {S: aws.String(calendarID)},
			"SortKey":    {S: aws.String("EVENT#" + eventID)},
		},
	}

	result, err := r.dynamoDB.GetItemWithContext(ctx, input)
	if err != nil {
		return false
	}

	if result.Item == nil {
		return false
	}

	return true
}

func (r *eventRepository) EditEvent(ctx context.Context, calendarID string, event *models.Event) (*models.Event, error) {
	updateExpression, attributeNames, attributeValues := buildUpdateExpression(event)

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"CalendarID": {S: aws.String(calendarID)},
			"SortKey":    {S: aws.String("EVENT#" + event.EventID)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  attributeNames,
		ExpressionAttributeValues: attributeValues,
		ReturnValues:              aws.String("ALL_NEW"),
	}

	result, err := r.dynamoDB.UpdateItemWithContext(ctx, updateInput)
	if err != nil {
		return nil, err
	}

	var updatedEvent models.Event
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedEvent)
	if err != nil {
		return nil, err
	}

	return &updatedEvent, nil
}

func buildUpdateExpression(event *models.Event) (string, map[string]*string, map[string]*dynamodb.AttributeValue) {
	return "SET Title = :title, Description = :desc, StartTime = :startTime, EndTime = :endTime, #location = :location, AllDay = :allDay",
		map[string]*string{
			"#location": aws.String("Location"),
		},
		map[string]*dynamodb.AttributeValue{
			":title":     {S: aws.String(event.Title)},
			":desc":      {S: aws.String(event.Description)},
			":startTime": {S: aws.String(event.StartTime)},
			":endTime":   {S: aws.String(event.EndTime)},
			":location":  {S: aws.String(event.Location)},
			":allDay":    {BOOL: aws.Bool(event.AllDay)},
		}
}
