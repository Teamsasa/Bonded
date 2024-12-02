package usecase

import (
	"bonded/internal/models"
	"context"
	"github.com/google/uuid"
)

func (u *calendarUsecase) CreateCalendar(ctx context.Context, calendar *models.Calendar) error {
	calendar.CalendarID = uuid.New().String()
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
