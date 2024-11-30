package handler

import (
	"bonded/internal/infra/db"
	"bonded/internal/models"
	"bonded/internal/usecase"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

var (
	calendarRepo    db.CalendarRepository
	calendarUsecase usecase.CalendarUsecase
)

func init() {
	calendarRepo = db.NewCalendarRepository()
	calendarUsecase = usecase.NewCalendarUsecase(calendarRepo)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.Path {
	case "/hello":
		return helloHandler(ctx, request)
	case "/dynamodb-test":
		return dynamoDBTestHandler(ctx, request)
	case "/calendar/create":
		return createCalendarHandler(ctx, request)
	case "/calendar/edit":
		return editCalendarHandler(ctx, request)
	case "/calendar/delete":
		return deleteCalendarHandler(ctx, request)
	case "/calendar/list":
		return listCalendarsHandler(ctx, request)
	case "/event/create":
		return createEventHandler(ctx, request)
	case "/event/edit":
		return editEventHandler(ctx, request)
	case "/event/delete":
		return deleteEventHandler(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			Body:       "Not Found",
			StatusCode: 404,
		}, nil
	}
}

func helloHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	greeting := usecase.GetGreeting(ctx, request.RequestContext.Identity.SourceIP)
	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

func dynamoDBTestHandler(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoRepo := db.NewDynamoDB()
	dynamoUsecase := usecase.NewDynamoUsecase(dynamoRepo)
	err := dynamoUsecase.DynamoDBTest(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to access DynamoDB table: " + err.Error(),
			StatusCode: 500,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       "Successfully accessed DynamoDB table",
		StatusCode: 200,
	}, nil
}

func createCalendarHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var calendar models.Calendar
	err := json.Unmarshal([]byte(request.Body), &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request payload: " + err.Error(),
			StatusCode: 400,
		}, nil
	}

	// ユーザーIDを設定（実際の環境では認証情報から取得）
	if calendar.UserID == "" {
		calendar.UserID = "user1"
	}

	err = calendarUsecase.CreateCalendar(ctx, &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to create calendar: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	responseBody, _ := json.Marshal(calendar)
	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 201,
	}, nil
}

func editCalendarHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var calendar models.Calendar
	err := json.Unmarshal([]byte(request.Body), &calendar)
	if err != nil || calendar.ID == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request payload",
			StatusCode: 400,
		}, nil
	}

	err = calendarUsecase.EditCalendar(ctx, &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to edit calendar: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Calendar updated successfully",
		StatusCode: 200,
	}, nil
}

func deleteCalendarHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarID := request.QueryStringParameters["id"]
	if calendarID == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Calendar ID is required",
			StatusCode: 400,
		}, nil
	}

	err := calendarUsecase.DeleteCalendar(ctx, calendarID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to delete calendar: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Calendar deleted successfully",
		StatusCode: 200,
	}, nil
}

func createEventHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event models.Event
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request payload: " + err.Error(),
			StatusCode: 400,
		}, nil
	}

	eventRepo := db.NewEventRepository()
	eventUsecase := usecase.NewEventUsecase(eventRepo)
	err = eventUsecase.CreateEvent(ctx, &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to create event: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	responseBody, _ := json.Marshal(event)
	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 201,
	}, nil
}

func editEventHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event models.Event
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil || event.ID == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request payload",
			StatusCode: 400,
		}, nil
	}

	eventRepo := db.NewEventRepository()
	eventUsecase := usecase.NewEventUsecase(eventRepo)
	err = eventUsecase.EditEvent(ctx, &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to edit event: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Event updated successfully",
		StatusCode: 200,
	}, nil
}

func deleteEventHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventID := request.QueryStringParameters["id"]
	if eventID == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Event ID is required",
			StatusCode: 400,
		}, nil
	}

	eventRepo := db.NewEventRepository()
	eventUsecase := usecase.NewEventUsecase(eventRepo)
	err := eventUsecase.DeleteEvent(ctx, eventID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to delete event: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Event deleted successfully",
		StatusCode: 200,
	}, nil
}

// カレンダー一覧を取得するハンドラー
func listCalendarsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 仮のユーザーIDを設定（実際の環境では認証情報から取得）
	userID := "user1"

	calendars, err := calendarUsecase.GetCalendars(ctx, userID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to get calendars: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	responseBody, err := json.Marshal(calendars)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to marshal response: " + err.Error(),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}
