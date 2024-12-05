package main

import (
	"bonded/internal/handler"
	"bonded/internal/infra/db"
	"bonded/internal/middleware"
	"bonded/internal/repository"
	"bonded/internal/usecase"
	"context"
	"fmt"
	"os"

	"github.com/MicahParks/keyfunc"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var jwks *keyfunc.JWKS

func init() {
	var err error
	jwks, err = keyfunc.Get(os.Getenv("COGNITO_JWKS_URL"), keyfunc.Options{})
	if err != nil {
		panic(fmt.Sprintf("Failed to get JWKS: %v", err))
	}
}

func main() {
	clientID := os.Getenv("COGNITO_CLIENT_ID")
	cognitoIssuer := os.Getenv("COGNITO_ISSUER")
	dynamoClient := db.DynamoDBClientRequest()
	calendarRepo := repository.CalendarRepositoryRequest(dynamoClient)
	eventRepo := repository.EventRepositoryRequest(dynamoClient)
	userRepo := repository.UserRepositoryRequest(dynamoClient)
	caledarUsecase := usecase.CalendarUsecaseRequest(calendarRepo, eventRepo, userRepo)
	authUsecase := usecase.NewAuthUsecase(jwks, clientID, cognitoIssuer)
	middleware := middleware.NewAuthMiddleware(authUsecase)
	h := handler.HandlerRequest(caledarUsecase)

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		authenticatedHandler := middleware.AuthMiddleware(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			switch request.Path {
			case "/hello":
				if request.HTTPMethod == "GET" {
					return h.HelloHandler(ctx, request)
				}
			case "/calendar/" + request.PathParameters["calendarId"]: // ok
				if request.HTTPMethod == "GET" {
					return h.HandleGetCalendar(ctx, request)
				}
			case "/calendar/list/" + request.PathParameters["userId"]: // ok
				if request.HTTPMethod == "GET" {
					return h.HandleGetCalendars(ctx, request)
				}
			case "/calendar/list/public":
				if request.HTTPMethod == "GET" {
					return h.HandleGetPublicCalendars(ctx, request)
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
			case "/event/create/" + request.PathParameters["calendarId"]: //
				if request.HTTPMethod == "POST" {
					return h.HandleCreateEvent(ctx, request)
				}
			}
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       "Not Found",
			}, nil
		})
		return authenticatedHandler(ctx, request)
	})
}
