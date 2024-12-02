package handler

import (
	"bonded/internal/infra/db"
	"bonded/internal/repository"
	"bonded/internal/usecase"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	Repo             repository.CalendarRepository
	CalendearUsecase usecase.CalendarUsecase
	EventUsecase     usecase.EventUsecase
}

func HandlerRequest(repo repository.CalendarRepository, usecase usecase.Usecase) *Handler {
	return &Handler{
		Repo:             repo,
		CalendearUsecase: usecase.Calendar(),
		EventUsecase:     usecase.Event(),
	}
}

func (h *Handler) HelloHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	greeting := usecase.GetGreeting(ctx, request.RequestContext.Identity.SourceIP)
	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

func (h *Handler) DynamoDBTestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dynamoRepo := db.NewDynamoDB()
	dynamoUsecase := usecase.DynamoUsecaseRequest(dynamoRepo)
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
