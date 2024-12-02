package usecase

import (
	"bonded/internal/models"
	"context"
	"github.com/google/uuid"
)

func (u *eventUsecase) CreateEvent(ctx context.Context, calendar *models.Calendar, event *models.Event) error {
	event.EventID = uuid.New().String()
	err := u.eventRepo.CreateEvent(ctx, calendar, event)
	if err != nil {
		return err
	}
	return nil
}
