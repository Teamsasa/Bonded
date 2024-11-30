package usecase

import (
	"bonded/internal/infra/db"
	"bonded/internal/models"
	"context"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	EditEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEvents(ctx context.Context, calendarID string) ([]*models.Event, error)
}

type eventUsecase struct {
	eventRepo db.EventRepository
}

func NewEventUsecase(eventRepo db.EventRepository) EventUsecase {
	return &eventUsecase{
		eventRepo: eventRepo,
	}
}

func (u *eventUsecase) CreateEvent(ctx context.Context, event *models.Event) error {
	return u.eventRepo.Save(ctx, event)
}

func (u *eventUsecase) EditEvent(ctx context.Context, event *models.Event) error {
	return u.eventRepo.Update(ctx, event)
}

func (u *eventUsecase) DeleteEvent(ctx context.Context, eventID string) error {
	return u.eventRepo.Delete(ctx, eventID)
}

func (u *eventUsecase) GetEvents(ctx context.Context, calendarID string) ([]*models.Event, error) {
	return u.eventRepo.FindByCalendarID(ctx, calendarID)
}
