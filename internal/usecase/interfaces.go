package usecase

import (
	"bonded/internal/models"
	"bonded/internal/repository"
	"context"
)

func CalendarUsecaseRequest(calendarRepo repository.CalendarRepository, eventRepo repository.EventRepository, userRepo repository.UserRepository) Usecase {
	return &usecase{
		calendarUsecase: &calendarUsecase{
			calendarRepo: calendarRepo,
			userRepo:     userRepo,
		},
		eventUsecase: &eventUsecase{
			eventRepo:    eventRepo,
			calendarRepo: calendarRepo,
		},
	}
}

type usecase struct {
	calendarUsecase CalendarUsecase
	eventUsecase    EventUsecase
}

type calendarUsecase struct {
	calendarRepo repository.CalendarRepository
	userRepo     repository.UserRepository
}

type eventUsecase struct {
	eventRepo    repository.EventRepository
	calendarRepo repository.CalendarRepository
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
	CreateCalendar(ctx context.Context, calendar *models.CreateCalendar) error
	EditCalendar(ctx context.Context, calendar *models.Calendar, input *models.Calendar) error
	DeleteCalendar(ctx context.Context, calendarID string) error
	FindPublicCalendars(ctx context.Context) ([]*models.Calendar, error)
	FindCalendars(ctx context.Context) ([]*models.Calendar, error)
	FindCalendar(ctx context.Context, calendarID string) (*models.Calendar, error)
	FollowCalendar(ctx context.Context, calendar *models.Calendar) error
	UnfollowCalendar(ctx context.Context, calendar *models.Calendar) error
}

type EventUsecase interface {
	CreateEvent(ctx context.Context, calendar *models.Calendar, event *models.Event) error
	FindEvents(ctx context.Context, calendarID string) ([]*models.Event, error)
	EditEvent(ctx context.Context, calendarID string, event *models.Event) (*models.Event, error)
	DeleteEvent(ctx context.Context, calendarID string, eventID string) error
}
