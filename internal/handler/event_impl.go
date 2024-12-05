package handler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"bonded/internal/models"
)

func (h *Handler) HandleCreateEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event models.Event
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error unmarshalling request: " + err.Error(),
		}, nil
	}
	calendarID := request.PathParameters["calendarId"]

	calendar, err := h.CalendarUsecase.FindCalendar(ctx, calendarID)
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

func (h *Handler) HandleEditEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event models.Event
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error unmarshalling request: " + err.Error(),
		}, nil
	}

	calendarID := request.PathParameters["calendarId"]
	updatedEvent, err := h.EventUsecase.EditEvent(ctx, calendarID, &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error editing event: " + err.Error(),
		}, nil
	}

	updatedEventJSON, err := json.Marshal(updatedEvent)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error marshalling response: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(updatedEventJSON),
	}, nil
}

func (h *Handler) HandleGetEventList(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarID := request.PathParameters["calendarId"]
	eventList, err := h.EventUsecase.FindEvents(ctx, calendarID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error finding calendar: " + err.Error(),
		}, nil
	}

	body, err := json.Marshal(eventList)
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
