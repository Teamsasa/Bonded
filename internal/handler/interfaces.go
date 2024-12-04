package handler

import (
	"bonded/internal/repository"
	"bonded/internal/usecase"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	Repo            repository.CalendarRepository
	CalendarUsecase usecase.CalendarUsecase
	EventUsecase    usecase.EventUsecase
}

func HandlerRequest(repo repository.CalendarRepository, usecase usecase.Usecase) *Handler {
	return &Handler{
		Repo:            repo,
		CalendarUsecase: usecase.Calendar(),
		EventUsecase:    usecase.Event(),
	}
}

func (h *Handler) HelloHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	greeting := usecase.GetGreeting(ctx, request.RequestContext.Identity.SourceIP)
	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

