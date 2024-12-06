package usecase

import (
	"bonded/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (u *eventUsecase) CreateEvent(ctx context.Context, calendar *models.Calendar, event *models.Event) error {
	event.EventID = uuid.New().String()
	return u.eventRepo.CreateEvent(ctx, calendar, event)
}

func (u *eventUsecase) FindEvents(ctx context.Context, calendarID string) ([]*models.Event, error) {
	return u.eventRepo.FindEvents(ctx, calendarID)
}

func (u *eventUsecase) EditEvent(ctx context.Context, calendarID string, event *models.Event) (*models.Event, error) {
	if event.EventID == "" {
		return nil, errors.New("eventID is required")
	}

	res, err := u.calendarRepo.FindByCalendarID(ctx, calendarID)
	if err != nil || res == nil {
		return nil, err
	}

	exists := u.eventRepo.EventExists(ctx, calendarID, event.EventID)
	if !exists {
		return nil, errors.New("event not found")
	}

	return u.eventRepo.EditEvent(ctx, calendarID, event)
}
