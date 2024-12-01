package usecase

import (
	"bonded/internal/models"
	"bonded/internal/repository"
	"context"
)

func EventUsecaseRequest(eventRepo repository.EventRepository) EventUsecase {
	return &eventUsecase{
		eventRepo: eventRepo,
	}
}

type EventUsecase interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	EditEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEvents(ctx context.Context, calendarID string) ([]*models.Event, error)
}

type eventUsecase struct {
	eventRepo repository.EventRepository
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
