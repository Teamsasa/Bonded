package usecase

import (
	"bonded/internal/models"
	"context"
	"github.com/google/uuid"
)

func (u *eventUsecase) CreateEvent(ctx context.Context, event *models.Event) error {
	event.EventID = uuid.New().String()
	return u.eventRepo.Create(ctx, event)
}

func (u *eventUsecase) FindEvent(ctx context.Context, eventID string) (*models.Event, error) {
	return u.eventRepo.FindByEventID(ctx, eventID)
}
