package usecase

import (
	"bonded/internal/models"
	"bonded/internal/repository"
	"context"
)

func CalendarUsecaseRequest(calendarRepo repository.CalendarRepository, eventRepo repository.EventRepository) Usecase {
	return &usecase{
		calendarUsecase: &calendarUsecase{
			calendarRepo: calendarRepo,
		},
		eventUsecase: &eventUsecase{
			eventRepo: eventRepo,
		},
	}
}

type usecase struct {
	calendarUsecase CalendarUsecase
	eventUsecase    EventUsecase
}

type calendarUsecase struct {
	calendarRepo repository.CalendarRepository
}

type eventUsecase struct {
	eventRepo repository.EventRepository
}

type Usecase interface {
	Calendar() CalendarUsecase
	Event() EventUsecase
}

func (u *usecase) Calendar() CalendarUsecase {
	return u.calendarUsecase
}

func (u *usecase) Event() EventUsecase {
	return u.eventUsecase
}

type CalendarUsecase interface {
	CreateCalendar(ctx context.Context, calendar *models.Calendar) error
	EditCalendar(ctx context.Context, calendar *models.Calendar) error
	DeleteCalendar(ctx context.Context, calendarID string) error
	FindCalendars(ctx context.Context, userID string) ([]*models.Calendar, error)
}

type EventUsecase interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	FindEvent(ctx context.Context, eventID string) (*models.Event, error)
}
