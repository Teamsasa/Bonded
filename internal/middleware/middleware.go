package middleware

import (
	"bonded/internal/usecase"
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type contextKey string

const jwtDataKey contextKey = "jwtData"

type IAuthMiddleware interface {
	AuthMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type authMiddleware struct {
	authUsecase usecase.IAuthUsecase
	publicPaths map[string]bool
}

func NewAuthMiddleware(authUsecase usecase.IAuthUsecase) IAuthMiddleware {
	publicPaths := map[string]bool{
		"/hello":                true,
		"/calendar/list/public": true,
	}

	return &authMiddleware{
		authUsecase: authUsecase,
		publicPaths: publicPaths,
	}
}

func (am *authMiddleware) AuthMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if am.publicPaths[request.Path] {
			return next(ctx, request)
		}

		authHeader, ok := request.Headers["Authorization"]
		if !ok || !strings.HasPrefix(authHeader, "Bearer ") {
			return unauthorizedResponse("Missing or invalid Authorization header")
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := am.authUsecase.ValidateJWT(accessToken)
		if err != nil {
			return unauthorizedResponse(err.Error())
		}

		idTokenHeader, ok := request.Headers["X-Id-Token"]
		if !ok || !strings.HasPrefix(idTokenHeader, "Bearer ") {
			return unauthorizedResponse("Missing or invalid ID token header")
		}
		idToken := strings.TrimPrefix(idTokenHeader, "Bearer ")
		jwtData, err := am.authUsecase.ValidateJWT(idToken)
		if err != nil {
			return unauthorizedResponse(err.Error())
		}

		ctx = context.WithValue(ctx, jwtDataKey, jwtData)

		return next(ctx, request)
	}
}

func unauthorizedResponse(message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 401,
		Body:       message,
	}, nil
}
