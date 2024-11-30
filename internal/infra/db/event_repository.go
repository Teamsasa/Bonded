package db

import (
	"bonded/internal/models"
	"context"
)

type EventRepository interface {
	Save(ctx context.Context, event *models.Event) error
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, eventID string) error
	FindByID(ctx context.Context, eventID string) (*models.Event, error)
	FindByCalendarID(ctx context.Context, calendarID string) ([]*models.Event, error)
}

// 実装例
type eventRepository struct {
	// データベース接続情報など
}

func NewEventRepository() EventRepository {
	return &eventRepository{
		// 初期化処理
	}
}

func (r *eventRepository) Save(ctx context.Context, event *models.Event) error {
	// イベントを保存する処理
	return nil
}

func (r *eventRepository) Update(ctx context.Context, event *models.Event) error {
	// イベントを更新する処理
	return nil
}

func (r *eventRepository) Delete(ctx context.Context, eventID string) error {
	// イベントを削除する処理
	return nil
}

func (r *eventRepository) FindByID(ctx context.Context, eventID string) (*models.Event, error) {
	// イベントをIDで検索する処理
	return nil, nil
}

func (r *eventRepository) FindByCalendarID(ctx context.Context, calendarID string) ([]*models.Event, error) {
	// カレンダーIDでイベントを検索する処理
	return nil, nil
}
