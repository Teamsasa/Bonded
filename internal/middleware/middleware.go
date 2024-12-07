package middleware

import (
	"bonded/internal/contextKey"
	"bonded/internal/usecase"
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type IAuthMiddleware interface {
	AuthMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type authMiddleware struct {
	authUsecase usecase.IAuthUsecase
	publicPaths map[string]*regexp.Regexp
}

func NewAuthMiddleware(authUsecase usecase.IAuthUsecase) IAuthMiddleware {
	publicPaths := map[string]*regexp.Regexp{
		"/hello":                regexp.MustCompile(`^/hello$`),
		"/calendar/list/public": regexp.MustCompile(`^/calendar/list/public$`),
		"/calendar/":            regexp.MustCompile(`^/calendar/[0-9a-fA-F-]+$`),
	}

	return &authMiddleware{
		authUsecase: authUsecase,
		publicPaths: publicPaths,
	}
}

func (am *authMiddleware) AuthMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		for _, re := range am.publicPaths {
			if re.MatchString(request.Path) {
				return next(ctx, request)
			}
		}

		authHeader, ok := request.Headers["Authorization"]
		if !ok || !strings.HasPrefix(authHeader, "Bearer ") {
			return unauthorizedResponse("Missing or invalid Authorization header")
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		jwtData, err := am.authUsecase.ValidateJWT(accessToken)
		if err != nil {
			return unauthorizedResponse(err.Error())
		}

		// idTokenHeader, ok := request.Headers["X-Id-Token"]
		// if !ok || !strings.HasPrefix(idTokenHeader, "Bearer ") {
		// 	return unauthorizedResponse("Missing or invalid ID token header")
		// }
		// idToken := strings.TrimPrefix(idTokenHeader, "Bearer ")
		// jwtData, err := am.authUsecase.ValidateJWT(idToken)
		// if err != nil {
		// 	return unauthorizedResponse(err.Error())
		// }

		ctx = context.WithValue(ctx, contextKey.JwtDataKey, jwtData)

		return next(ctx, request)
	}
}

func unauthorizedResponse(message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 401,
		Body:       message,
	}, nil
}
