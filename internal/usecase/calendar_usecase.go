package usecase

import (
	"bonded/internal/models"
	"bonded/internal/repository"
	"context"
)

func CalendarUsecaseRequest(calendarRepo repository.CalendarRepository) CalendarUsecase {
	return &calendarUsecase{
		calendarRepo: calendarRepo,
	}
}

type CalendarUsecase interface {
	CreateCalendar(ctx context.Context, calendar *models.Calendar) error
	EditCalendar(ctx context.Context, calendar *models.Calendar) error
	DeleteCalendar(ctx context.Context, calendarID string) error
	GetCalendars(ctx context.Context, userID string) ([]*models.Calendar, error)
}

type calendarUsecase struct {
	calendarRepo repository.CalendarRepository
}

func NewCalendarUsecase(calendarRepo repository.CalendarRepository) CalendarUsecase {
	return &calendarUsecase{
		calendarRepo: calendarRepo,
	}
}

func (u *calendarUsecase) CreateCalendar(ctx context.Context, calendar *models.Calendar) error {
	return u.calendarRepo.Save(ctx, calendar)
}

func (u *calendarUsecase) EditCalendar(ctx context.Context, calendar *models.Calendar) error {
	return u.calendarRepo.Update(ctx, calendar)
}

func (u *calendarUsecase) DeleteCalendar(ctx context.Context, calendarID string) error {
	return u.calendarRepo.Delete(ctx, calendarID)
}

func (u *calendarUsecase) GetCalendars(ctx context.Context, userID string) ([]*models.Calendar, error) {
	return u.calendarRepo.FindByUserID(ctx, userID)
}
