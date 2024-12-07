package handler

import (
	"bonded/internal/contextKey"
	"bonded/internal/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v4"
)

func (h *Handler) HandleGetCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendarID := request.PathParameters["calendarId"]
	calendar, err := h.CalendarUsecase.FindCalendar(ctx, calendarID)
	if err != nil || calendar == nil {
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
	calendars, err := h.CalendarUsecase.FindCalendars(ctx)
	if err != nil || calendars == nil {
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

func (h *Handler) HandleGetPublicCalendars(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	calendars, err := h.CalendarUsecase.FindPublicCalendars(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error finding public calendars: " + err.Error(),
		}, nil
	}
	if calendars == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "No public calendars found",
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

func (h *Handler) HandleUnfollowCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody struct {
		CalendarID string `json:"calendarId"`
	}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload: " + err.Error(),
		}, nil
	}
	if requestBody.CalendarID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing required fields: calendarId",
		}, nil
	}

	calendar, err := h.CalendarUsecase.FindCalendar(ctx, requestBody.CalendarID)
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

	err = h.CalendarUsecase.UnfollowCalendar(ctx, calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error unfollowing calendar: " + err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message":"Calendar unfollowed successfully."}`,
	}, nil
}

func (h *Handler) HandleCreateCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var calendar models.CreateCalendar
	err := json.Unmarshal([]byte(request.Body), &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload: " + err.Error(),
		}, nil
	}

	if calendar.Name == "" || calendar.IsPublic == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing required fields: name or isPublic",
		}, nil
	}

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
	var input models.Calendar
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload" + fmt.Sprint(err),
		}, nil
	}
	calendarId := request.PathParameters["calendarId"]
	input.CalendarID = calendarId

	calendar, err := h.CalendarUsecase.FindCalendar(ctx, input.CalendarID)
	if err != nil || calendar == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Calendar not found",
		}, nil
	}

	err = h.CalendarUsecase.EditCalendar(ctx, calendar, &input)
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

func (h *Handler) HandleFollowCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody struct {
		CalendarID string `json:"calendarId"`
	}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload: " + err.Error(),
		}, nil
	}
	if requestBody.CalendarID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing required fields: calendarId",
		}, nil
	}

	calendar, err := h.CalendarUsecase.FindCalendar(ctx, requestBody.CalendarID)
	if err != nil || calendar == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error finding calendar: " + err.Error(),
		}, nil
	}
	if !*calendar.IsPublic {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Calendar is not public",
		}, nil
	}

	err = h.CalendarUsecase.FollowCalendar(ctx, calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error following calendar: " + err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message":"Calendar followed successfully."}`,
	}, nil
}

func (h *Handler) HandleInviteUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody struct {
		InviteUserID string `json:"inviteUserId"`
		CalendarID   string `json:"calendarId"`
		AccessLevel  string `json:"accessLevel"`
	}

	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request payload: " + err.Error(),
		}, nil
	}

	fmt.Println("Inviting user", requestBody.InviteUserID, "to calendar", requestBody.CalendarID, "with access level", requestBody.AccessLevel)

	// アクセスレベルの検証
	if requestBody.AccessLevel != "EDITOR" && requestBody.AccessLevel != "VIEWER" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid access level. Must be either 'EDITOR' or 'VIEWER'",
		}, nil
	}

	// トークンからユーザーIDを取得
	jwtData := ctx.Value(contextKey.JwtDataKey).(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	err = h.CalendarUsecase.InviteUser(ctx, requestBody.CalendarID, userID, requestBody.InviteUserID, requestBody.AccessLevel)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error inviting user: " + err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message":"User invited successfully"}`,
	}, nil
}
