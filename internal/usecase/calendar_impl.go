package usecase

import (
	"bonded/internal/models"
	"context"

	"github.com/google/uuid"
)

func (u *calendarUsecase) FindCalendar(ctx context.Context, calendarID string) (*models.Calendar, error) {
	return u.calendarRepo.FindByCalendarID(ctx, calendarID)
}

func (u *calendarUsecase) CreateCalendar(ctx context.Context, calendar *models.CreateCalendar) error {
	if calendar.OwnerName == "" {
		user, err := u.userRepo.FindByUserID(ctx, calendar.OwnerUserID)
		if err != nil {
			return err
		}
		calendar.OwnerName = user.DisplayName
	}
	user := models.User{
		UserID:      calendar.OwnerUserID,
		DisplayName: calendar.OwnerName,
		Email:       calendar.OwnerUserID + "@example.com",
		Password:    "password",
		AccessLevel: "OWNER",
	}
	calendar.Users = []models.User{user}

	calendar.CalendarID = uuid.New().String()

	// CreateCalendarのフィールドをCalendarに変換
	calendarReq := models.Calendar{
		CalendarID:  calendar.CalendarID,
		SortKey:     "CALENDAR",
		Name:        calendar.Name,
		IsPublic:    calendar.IsPublic,
		OwnerUserID: calendar.OwnerUserID,
		Users:       calendar.Users,
		Events:      calendar.Events,
	}

	return u.calendarRepo.Create(ctx, &calendarReq)
}

func (u *calendarUsecase) EditCalendar(ctx context.Context, calendar *models.Calendar, input *models.Calendar) error {
	return u.calendarRepo.Edit(ctx, calendar, input)
}

func (u *calendarUsecase) DeleteCalendar(ctx context.Context, calendarID string) error {
	return u.calendarRepo.Delete(ctx, calendarID)
}

func (u *calendarUsecase) FindCalendars(ctx context.Context, userID string) ([]*models.Calendar, error) {
	return u.calendarRepo.FindByUserID(ctx, userID)
}
