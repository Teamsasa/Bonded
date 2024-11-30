package handler

import (
	"bonded/internal/infra/db"
	"bonded/internal/usecase"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.Path {
	case "/hello":
		return helloHandler(ctx, request)
	case "/dynamodb-test":
		return dynamoDBTestHandler(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			Body:       "Not Found",
			StatusCode: 404,
		}, nil
	}
}

func helloHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	greeting := usecase.GetGreeting(ctx, request.RequestContext.Identity.SourceIP)
	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

func dynamoDBTestHandler(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoRepo := db.NewDynamoDB()
	dynamoUsecase := usecase.NewDynamoDBUsecase(dynamoRepo)
	err := dynamoUsecase.DynamoDBTest(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to access DynamoDB table: " + err.Error(),
			StatusCode: 500,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       "Successfully accessed DynamoDB table",
		StatusCode: 200,
	}, nil
}
