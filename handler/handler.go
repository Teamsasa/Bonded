package handler

import (
	"bonded/internal/infra/db"
	"bonded/internal/models"
	"bonded/usecase"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Handler 構造体はリポジトリを保持します。
type Handler struct {
	Repo db.CalendarRepository
}

// HandleGetCalendars は GET /calendar/list のハンドラーです。
func (h *Handler) HandleGetCalendars(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := request.QueryStringParameters["userId"]
	calendars, err := h.Repo.FindByUserID(ctx, userID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error finding calendars: " + err.Error(),
		}, nil
	}
	body, err := json.Marshal(calendars)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error marshalling response: " + err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

// HandleCreateCalendar は POST /calendar/create のハンドラーです。
func (h *Handler) HandleCreateCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var calendar models.Calendar
	err := json.Unmarshal([]byte(request.Body), &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload",
		}, nil
	}
	err = h.Repo.Save(ctx, &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error saving calendar: " + err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message":"Calendar created successfully."}`,
	}, nil
}

// HandlePutCalendarUpdate は PUT /calendar/update/{id} のハンドラーです。
func (h *Handler) HandlePutCalendarUpdate(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	var input struct {
		Name string `json:"name"`
	}
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload",
		}, nil
	}

	calendar, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Calendar not found",
		}, nil
	}

	calendar.Name = input.Name
	err = h.Repo.Update(ctx, calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to update calendar",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message":"Calendar updated successfully."}`,
	}, nil
}

// HelloHandler は GET /hello のハンドラーです。
func (h *Handler) HelloHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	greeting := usecase.GetGreeting(ctx, request.RequestContext.Identity.SourceIP)
	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

// DynamoDBTestHandler は GET /dynamodb-test のハンドラーです。
func (h *Handler) DynamoDBTestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
