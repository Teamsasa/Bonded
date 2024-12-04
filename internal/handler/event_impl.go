package handler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"bonded/internal/models"
)

func (h *Handler) HandleCreateEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarID := request.PathParameters["calendarId"]
	var event models.Event
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error unmarshalling request: " + err.Error(),
		}, nil
	}

	calendar, err := h.CalendarUsecase.FindCalendar(ctx, calendarID) // カレンダーで返される
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error finding calendar: " + err.Error(),
		}, nil
	}
	if calendar == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Calendar not found",
		}, nil
	}

	err = h.EventUsecase.CreateEvent(ctx, calendar, &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating event: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       `{"message":"Event created successfully."}`,
	}, nil
}
