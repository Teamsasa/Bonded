package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"

	"bonded/internal/models"
)

// 一旦、個別取得のみ。userIDから一覧取得とかもあった方がいいかも
func (h *Handler) HandleGetEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventID := request.PathParameters["eventId"]
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
		Body:       `{"message":"Event created successfully.", "eventId":"` + event.EventID + `"}`,
	}, nil
}

func (h *Handler) HandleGetEvents(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarID := request.PathParameters["calendarId"]

	// カレンダーの存在確認
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

	// イベントを取得
	eventList, err := h.EventUsecase.FindEvents(ctx, calendarID) // 変数名を変更
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error finding events: " + err.Error(),
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
