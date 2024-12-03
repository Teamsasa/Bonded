package handler

import (
	"bonded/internal/models"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func (h *Handler) HandleGetCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarID := request.PathParameters["calendarId"]
	calendar, err := h.CalendarUsecase.FindCalendar(ctx, calendarID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error finding calendar: " + err.Error(),
		}, nil
	}

	body, err := json.Marshal(calendar)
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

func (h *Handler) HandleGetCalendars(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := request.PathParameters["userId"]
	calendars, err := h.CalendarUsecase.FindCalendars(ctx, userID)
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
	userID := request.PathParameters["userId"]
	var calendar models.Calendar
	err := json.Unmarshal([]byte(request.Body), &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload",
			}, nil
		}
	calendar.UserID = userID

	err = h.CalendarUsecase.CreateCalendar(ctx, &calendar)
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

func (h *Handler) HandleEditCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	calendar, err := h.CalendarUsecase.FindCalendar(ctx, calendarId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Calendar not found",
		}, nil
	}

	calendar.Name = input.Name
	err = h.CalendarUsecase.EditCalendar(ctx, calendar)
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
	err := h.CalendarUsecase.DeleteCalendar(ctx, calendarId)
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
