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
	dynamoClient, err := db.DynamoDBClientRequest()
	if err != nil {
		panic(err)
	}
	calendarRepo := repository.CalendarRepositoryRequest(dynamoClient)
	eventRepo := repository.EventRepositoryRequest(dynamoClient)
	appUsecase := usecase.CalendarUsecaseRequest(calendarRepo, eventRepo)
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
		case "/calendar/" + request.PathParameters["calendarId"]: // ok
			if request.HTTPMethod == "GET" {
				return h.HandleGetCalendar(ctx, request)
			}
		case "/calendar/list/" + request.PathParameters["userId"]: // ok
			if request.HTTPMethod == "GET" {
				return h.HandleGetCalendars(ctx, request)
			}
		case "/calendar/create/" + request.PathParameters["userId"]: // ok
			if request.HTTPMethod == "POST" {
				return h.HandleCreateCalendar(ctx, request)
			}
		case "/calendar/edit/" + request.PathParameters["calendarId"]: // ok 今は誰でも編集できる状態になっているので、呼び出す時にEDITORかどうかを見たい
			if request.HTTPMethod == "PUT" {
				return h.HandleEditCalendar(ctx, request)
			}
		case "/calendar/delete/" + request.PathParameters["calendarId"]: // ok
			if request.HTTPMethod == "DELETE" {
				return h.HandleDeleteCalendar(ctx, request)
			}
		case "/event/create/" + request.PathParameters["calendarId"]:
			if request.HTTPMethod == "POST" {
				return h.HandleCreateEvent(ctx, request)
			}
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	})
}
