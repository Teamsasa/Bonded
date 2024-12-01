package main

import (
	"bonded/internal/handler"
	"bonded/internal/infra/db"
	"bonded/internal/repository"
	"bonded/internal/usecase"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	dynamoClient := db.DynamoDBClientRequest()
	calendarRepo := repository.CalendarRepositoryRequest(dynamoClient)
	appUsecase := usecase.CalendarUsecaseRequest(calendarRepo)
	h := handler.HandlerRequest(calendarRepo, appUsecase)

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/hello":
			if request.HTTPMethod == "GET" {
				return h.HelloHandler(ctx, request)
			}
		case "/dynamodb-test":
			if request.HTTPMethod == "GET" {
				return h.DynamoDBTestHandler(ctx, request)
			}
		case "/calendar/list":
			if request.HTTPMethod == "GET" {
				return h.HandleGetCalendars(ctx, request)
			}
		case "/calendar/create":
			if request.HTTPMethod == "POST" {
				return h.HandleCreateCalendar(ctx, request)
			}
		case "/calendar/edit/" + request.PathParameters["id"]:
			if request.HTTPMethod == "PUT" {
				return h.HandlePutCalendarEdit(ctx, request)
			}
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	})
}
