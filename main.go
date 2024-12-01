package main

import (
	"bonded/handler"
	"bonded/internal/infra/db"
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	repo := db.NewCalendarRepository()
	h := &handler.Handler{
		Repo: repo,
	}

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		path := request.Path
		method := request.HTTPMethod

		switch {
		case path == "/hello" && method == "GET":
			return h.HelloHandler(ctx, request)
		case path == "/dynamodb-test" && method == "GET":
			return h.DynamoDBTestHandler(ctx, request)
		case path == "/calendar/list" && method == "GET":
			return h.HandleGetCalendars(ctx, request)
		case path == "/calendar/create" && method == "POST":
			return h.HandleCreateCalendar(ctx, request)
		case strings.HasPrefix(path, "/calendar/update/") && method == "PUT":
			return h.HandlePutCalendarUpdate(ctx, request)
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       "Not Found",
			}, nil
		}
	})
}
