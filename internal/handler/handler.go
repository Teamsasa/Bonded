package handler

import (
	"bonded/internal/infra/db"
	"bonded/internal/models"
	"bonded/internal/repository"
	"bonded/internal/usecase"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	Repo             repository.CalendarRepository
	CalendearUsecase usecase.CalendarUsecase
	EventUsecase     usecase.EventUsecase
}

func HandlerRequest(repo repository.CalendarRepository, usecase usecase.Usecase) *Handler {
	return &Handler{
		Repo:             repo,
		CalendearUsecase: usecase.Calendar(),
		EventUsecase:     usecase.Event(),
	}
}

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

func (h *Handler) HandleCreateCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var calendar models.Calendar
	err := json.Unmarshal([]byte(request.Body), &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload",
		}, nil
	}

	err = h.CalendearUsecase.CreateCalendar(ctx, &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error saving calendar: " + err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       `{"message":"Calendar created successfully."}`,
	}, nil
}

func (h *Handler) HandlePutCalendarEdit(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarId := request.PathParameters["calendarId"]

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

	calendar, err := h.Repo.FindByCalendarID(ctx, calendarId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Calendar not found",
		}, nil
	}

	calendar.Name = input.Name
	err = h.Repo.Edit(ctx, calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to edit calendar",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message":"Calendar edited successfully."}`,
	}, nil
}

func (h *Handler) HandleDeleteCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarId := request.PathParameters["calendarId"]

	err := h.Repo.Delete(ctx, calendarId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to delete calendar",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message":"Calendar deleted successfully."}`,
	}, nil
}

func (h *Handler) HelloHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	greeting := usecase.GetGreeting(ctx, request.RequestContext.Identity.SourceIP)
	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

func (h *Handler) DynamoDBTestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoRepo := db.NewDynamoDB()
	dynamoUsecase := usecase.DynamoUsecaseRequest(dynamoRepo)
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

// 一旦、個別取得のみ。userIDから一覧取得とかもあった方がいいかも
func (h *Handler) HandleGetEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventID := request.PathParameters["id"]
	event, err := h.EventUsecase.FindEvent(ctx, eventID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error finding event: " + err.Error(),
		}, nil
	}

	body, err := json.Marshal(event)
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

func (h *Handler) HandleCreateEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event models.Event
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload",
		}, nil
	}

	err = h.EventUsecase.CreateEvent(ctx, &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error saving event: " + err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       `{"message":"イベントが正常に作成されました "}` + event.ID,
	}, nil
}
