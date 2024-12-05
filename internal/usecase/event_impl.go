package usecase

import (
	"bonded/internal/models"
	"context"
	"github.com/google/uuid"
)

func (u *eventUsecase) CreateEvent(ctx context.Context, calendar *models.Calendar, event *models.Event) error {
	event.EventID = uuid.New().String()
	return u.eventRepo.CreateEvent(ctx, calendar, event)
}

func (u *eventUsecase) FindEvents(ctx context.Context, calendarID string) ([]*models.Event, error) {
	return u.eventRepo.FindEvents(ctx, calendarID)
}
