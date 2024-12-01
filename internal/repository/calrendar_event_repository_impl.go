package repository

import (
	"bonded/internal/models"
	"context"
)

type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	Edit(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, eventID string) error
	FindByID(ctx context.Context, eventID string) (*models.Event, error)
	FindByCalendarID(ctx context.Context, calendarID string) ([]*models.Event, error)
}

// type eventRepository struct {
// 	repo *EventRepository
// }

// func (r *eventRepository) Cretae(ctx context.Context, event *models.Event) error {
// 	// イベントを保存する処理
// 	return nil
// }

// func (r *eventRepository) Edit(ctx context.Context, event *models.Event) error {
// 	// イベントを更新する処理
// 	return nil
// }

// func (r *eventRepository) Delete(ctx context.Context, eventID string) error {
// 	// イベントを削除する処理
// 	return nil
// }

// func (r *eventRepository) FindByID(ctx context.Context, eventID string) (*models.Event, error) {
// 	// イベントをIDで検索する処理
// 	return nil, nil
// }

// func (r *eventRepository) FindByCalendarID(ctx context.Context, calendarID string) ([]*models.Event, error) {
// 	// カレンダーIDでイベントを検索する処理
// 	return nil, nil
// }
