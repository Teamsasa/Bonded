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
	// 1. 関連アイテムの作成
	relatedItem := map[string]*dynamodb.AttributeValue{
		"CalendarID": {
			S: aws.String(calendar.CalendarID),
		},
		"SortKey": {
			S: aws.String(fmt.Sprintf("CAL#%s#%s", calendar.CalendarID, calendar.OwnerUserID)),
		},
		"UserID": {
			S: aws.String(calendar.OwnerUserID),
		},
	}

	relatedInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      relatedItem,
	}
	_, err := r.dynamoDB.PutItemWithContext(ctx, relatedInput)
	if err != nil {
		return err
	}

	// 2. メインカレンダーアイテムの作成
	mainItem := map[string]*dynamodb.AttributeValue{
		"SortKey": {
			S: aws.String("CALENDAR"),
		},
		"OwnerUserID": {
			S: aws.String(calendar.OwnerUserID),
		},
		"CalendarID": {
			S: aws.String(calendar.CalendarID),
		},
		"IsPublic": {
			BOOL: aws.Bool(*calendar.IsPublic),
		},
		"Name": {
			S: aws.String(calendar.Name),
		},
	}

	mainInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      mainItem,
	}

	_, err = r.dynamoDB.PutItemWithContext(ctx, mainInput)
	if err != nil {
		return err
	}

	// 3. ユーザーアイテムの作成
	userItem := map[string]*dynamodb.AttributeValue{
		"Email": {
			S: aws.String(calendar.Users[0].Email),
		},
		"CalendarID": {
			S: aws.String(calendar.CalendarID),
		},
		"UserID": {
			S: aws.String(calendar.Users[0].UserID),
		},
		"DisplayName": {
			S: aws.String(calendar.Users[0].DisplayName),
		},
		"SortKey": {
			S: aws.String(fmt.Sprintf("USER#%s", calendar.Users[0].UserID)),
		},
		"AccessLevel": {
			S: aws.String(calendar.Users[0].AccessLevel),
		},
		"Password": {
			S: aws.String(calendar.Users[0].Password),
		},
	}

	userInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      userItem,
	}
	_, err = r.dynamoDB.PutItemWithContext(ctx, userInput)
	if err != nil {
		return err
	}

	return nil
}

func (r *calendarRepository) Edit(ctx context.Context, calendarID *models.Calendar, input *models.Calendar) error {
	calendar, err := r.FindByCalendarID(ctx, calendarID.CalendarID)
	if err != nil {
		return err
	}

	if input.Name != "" {
		calendar.Name = input.Name
	}
	if input.IsPublic != nil {
		calendar.IsPublic = input.IsPublic
	}
	if input.OwnerUserID != "" {
		calendar.OwnerUserID = input.OwnerUserID
	}

	item, err := dynamodbattribute.MarshalMap(calendar)
	if err != nil {
		return err
	}
	item["SortKey"] = &dynamodb.AttributeValue{S: aws.String("CALENDAR")}
	item["UserID"] = &dynamodb.AttributeValue{S: aws.String(calendar.OwnerUserID)}
	item["IsPublic"] = &dynamodb.AttributeValue{BOOL: calendar.IsPublic}
	updateInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}
	_, err = r.dynamoDB.PutItemWithContext(ctx, updateInput)
	return err
}

func (r *calendarRepository) Delete(ctx context.Context, calendarID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"CalendarID": {S: aws.String(calendarID)},
			"SortKey":    {S: aws.String("CALENDAR")},
		},
	}
	_, err := r.dynamoDB.DeleteItemWithContext(ctx, input)
	return err
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
		TableName:              aws.String(r.tableName),
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
		TableName:              aws.String(r.tableName),
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
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("UserID-index"),
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

func (r *calendarRepository) FollowCalendar(ctx context.Context, calendar *models.Calendar, user *models.User) error {
	item, err := dynamodbattribute.MarshalMap(calendar)
	if err != nil {
		return err
	}

	// メインカレンダーアイテムの更新
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.dynamoDB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	// 関連アイテムの作成
	relatedItem := map[string]*dynamodb.AttributeValue{
		"CalendarID": {
			S: aws.String(calendar.CalendarID),
		},
		"SortKey": {
			S: aws.String(fmt.Sprintf("CAL#%s#%s", calendar.CalendarID, user.UserID)),
		},
		"UserID": {
			S: aws.String(user.UserID),
		},
	}

	relatedInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      relatedItem,
	}
	_, err = r.dynamoDB.PutItemWithContext(ctx, relatedInput)
	if err != nil {
		return err
	}

	// ユーザー情報の作成
	userItem := map[string]*dynamodb.AttributeValue{
		"Email": {
			S: aws.String(user.Email),
		},
		"CalendarID": {
			S: aws.String(calendar.CalendarID),
		},
		"UserID": {
			S: aws.String(user.UserID),
		},
		"DisplayName": {
			S: aws.String(user.DisplayName),
		},
		"SortKey": {
			S: aws.String(fmt.Sprintf("USER#%s", user.UserID)),
		},
		"AccessLevel": {
			S: aws.String("VIEWER"),
		},
		"Password": {
			S: aws.String(user.Password),
		},
	}

	userInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      userItem,
	}
	_, err = r.dynamoDB.PutItemWithContext(ctx, userInput)
	if err != nil {
		return err
	}
	return nil
}

func (r *calendarRepository) FindAllCalendars(ctx context.Context) ([]*models.Calendar, error) {
	// カレンダーのメイン情報を全件取得
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("SortKey = :sk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sk": {S: aws.String("CALENDAR")},
		},
	}
	result, err := r.dynamoDB.ScanWithContext(ctx, input)
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

func (r *calendarRepository) UnfollowCalendar(ctx context.Context, calendar *models.Calendar, user *models.User) error {
	// メインカレンダーアイテムの更新
	item, err := dynamodbattribute.MarshalMap(calendar)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}
	_, err = r.dynamoDB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	// 関連アイテムの削除
	relatedInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"CalendarID": {S: aws.String(calendar.CalendarID)},
			"SortKey":    {S: aws.String(fmt.Sprintf("CAL#%s#%s", calendar.CalendarID, user.UserID))},
		},
	}
	_, err = r.dynamoDB.DeleteItemWithContext(ctx, relatedInput)
	if err != nil {
		return err
	}

	// ユーザーアイテムの削除
	userInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"CalendarID": {S: aws.String(calendar.CalendarID)},
			"SortKey":    {S: aws.String(fmt.Sprintf("USER#%s", user.UserID))},
		},
	}
	_, err = r.dynamoDB.DeleteItemWithContext(ctx, userInput)
	if err != nil {
		return err
	}
	return nil
}

func (r *calendarRepository) InviteUser(ctx context.Context, calendar *models.Calendar, user *models.User) error {
	fmt.Println("InviteUser")
	fmt.Println(calendar)
	fmt.Println(user)
	// 関連アイテムの作成
	relatedItem := map[string]*dynamodb.AttributeValue{
		"CalendarID": {
			S: aws.String(calendar.CalendarID),
		},
		"SortKey": {
			S: aws.String(fmt.Sprintf("CAL#%s#%s", calendar.CalendarID, user.UserID)),
		},
		"UserID": {
			S: aws.String(user.UserID),
		},
	}

	relatedInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      relatedItem,
	}
	_, err := r.dynamoDB.PutItemWithContext(ctx, relatedInput)
	if err != nil {
		return err
	}

	// ユーザー情報の作成
	userItem := map[string]*dynamodb.AttributeValue{
		"Email": {
			S: aws.String(user.Email),
		},
		"CalendarID": {
			S: aws.String(calendar.CalendarID),
		},
		"UserID": {
			S: aws.String(user.UserID),
		},
		"DisplayName": {
			S: aws.String(user.DisplayName),
		},
		"SortKey": {
			S: aws.String(fmt.Sprintf("USER#%s", user.UserID)),
		},
		"AccessLevel": {
			S: aws.String(user.AccessLevel),
		},
		"Password": {
			S: aws.String(user.Password),
		},
	}

	userInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      userItem,
	}
	_, err = r.dynamoDB.PutItemWithContext(ctx, userInput)
	return err
}
