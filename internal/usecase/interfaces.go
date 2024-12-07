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
	userRepo     repository.UserRepository
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
	CreateCalendar(ctx context.Context, calendar *models.CreateCalendar) error
	EditCalendar(ctx context.Context, calendar *models.Calendar, input *models.Calendar) error
	DeleteCalendar(ctx context.Context, calendarID string) error
	FindPublicCalendars(ctx context.Context) ([]*models.Calendar, error)
	FindCalendars(ctx context.Context, userID string) ([]*models.Calendar, error)
	FindCalendar(ctx context.Context, calendarID string) (*models.Calendar, error)
	FollowCalendar(ctx context.Context, calendar *models.Calendar, userID string) error
	UnfollowCalendar(ctx context.Context, calendar *models.Calendar, userID string) error
}

type EventUsecase interface {
	CreateEvent(ctx context.Context, calendar *models.Calendar, event *models.Event) error
	FindEvents(ctx context.Context, calendarID string) ([]*models.Event, error)
}
