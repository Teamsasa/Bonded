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
	appUsecase := usecase.CalendarUsecaseRequest(calendarRepo)
	authUsecase := usecase.NewAuthUsecase(jwks, clientID, cognitoIssuer)
	middleware := middleware.NewAuthMiddleware(authUsecase)
	h := handler.HandlerRequest(calendarRepo, appUsecase)

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		authenticatedHandler := middleware.AuthMiddleware(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
			}
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       "Not Found",
			}, nil
		})
		return authenticatedHandler(ctx, request)
	})
}
