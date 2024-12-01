package usecase

import (
	"bonded/internal/models"
	"bonded/internal/repository"
	"context"
	"github.com/google/uuid"
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

type Usecase interface {
	Calendar() CalendarUsecase
	Event() EventUsecase
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

type usecase struct {
	calendarUsecase CalendarUsecase
	eventUsecase    EventUsecase
}

func (u *usecase) Calendar() CalendarUsecase {
	return u.calendarUsecase
}

func (u *usecase) Event() EventUsecase {
	return u.eventUsecase
}

type calendarUsecase struct {
	calendarRepo repository.CalendarRepository
}

func (u *calendarUsecase) CreateCalendar(ctx context.Context, calendar *models.Calendar) error {
	calendar.ID = uuid.New().String()
	return u.calendarRepo.Create(ctx, calendar)
}

func (u *calendarUsecase) EditCalendar(ctx context.Context, calendar *models.Calendar) error {
	return u.calendarRepo.Edit(ctx, calendar)
}

func (u *calendarUsecase) DeleteCalendar(ctx context.Context, calendarID string) error {
	return u.calendarRepo.Delete(ctx, calendarID)
}

func (u *calendarUsecase) FindCalendars(ctx context.Context, userID string) ([]*models.Calendar, error) {
	return u.calendarRepo.FindByUserID(ctx, userID)
}

type eventUsecase struct {
	eventRepo repository.EventRepository
}

func (u *eventUsecase) CreateEvent(ctx context.Context, event *models.Event) error {
	event.ID = uuid.New().String()
	return u.eventRepo.Create(ctx, event)
}

func (u *eventUsecase) FindEvent(ctx context.Context, eventID string) (*models.Event, error) {
	return u.eventRepo.FindByEventID(ctx, eventID)
}
