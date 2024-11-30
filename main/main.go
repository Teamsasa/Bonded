package main

import (
	"context"
	"encoding/json"
	"log" // ログパッケージをインポート
	"net/http"
	"os"

	"bonded/internal/infra/db"
	"bonded/internal/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Handler struct {
	Repo db.CalendarRepository
}

func NewHandler() *Handler {
	log.Println("Initializing Handler") // 初期化ログ
	return &Handler{
		Repo: db.NewCalendarRepository(),
	}
}

func (h *Handler) ListCalendars(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("ListCalendars called with request:", request) // 関数呼び出しログ
	userID := request.QueryStringParameters["userId"]
	if userID == "" {
		log.Println("userId is missing") // エラーログ
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"userId is required"}`,
		}, nil
	}

	calendars, err := h.Repo.FindByUserID(ctx, userID)
	if err != nil {
		log.Println("Error finding calendars:", err) // エラーログ
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"Internal server error"}`,
		}, nil
	}
	body, err := json.Marshal(calendars)
	if err != nil {
		log.Println("Error marshalling calendars:", err) // エラーログ
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"Failed to marshal response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func (h *Handler) CreateCalendar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var calendar models.Calendar
	err := json.Unmarshal([]byte(request.Body), &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"Invalid request body"}`,
		}, nil
	}

	if calendar.ID == "" || calendar.UserID == "" || calendar.Name == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"id, userId, and name are required"}`,
		}, nil
	}

	err = h.Repo.Save(ctx, &calendar)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"Failed to create calendar"}`,
		}, nil
	}

	response := map[string]string{
		"message": "Calendar created successfully.",
	}
	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

// handlerProxy 関数を定義し、正しいシグネチャで lambda.Start に渡す
func handlerProxy(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	handler := NewHandler()
	switch event.Path {
	case "/calendar/list":
		return handler.ListCalendars(ctx, event)
	case "/calendar/create":
		return handler.CreateCalendar(ctx, event)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error":"Not found"}`,
		}, nil
	}
}

func main() {
	// 環境変数が設定されていない場合はデフォルト値を使用
	if os.Getenv("DYNAMODB_ENDPOINT") == "" {
		os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8000")
	}

	// handlerProxy 関数を lambda.Start に渡す
	lambda.Start(handlerProxy)
}
