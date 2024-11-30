package usecase

import (
	"context"
)

type IGreetungUsecase interface {
	GetGreeting(ctx context.Context, sourceIP string) string
}

type GreetingUsecase struct {
}

func GetGreeting(ctx context.Context, sourceIP string) string {
	if sourceIP == "" {
		return "Hello, world!\n"
	}
	return "Hello, " + sourceIP + "!\n"
}
