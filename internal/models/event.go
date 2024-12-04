package models

import "time"

type Event struct {
    EventID     string    `json:"eventId" dynamodbav:"EventID"`       // イベントID
    Title       string    `json:"title" dynamodbav:"Title"`           // イベント名
    Description string    `json:"description" dynamodbav:"Description"` // 詳細
    StartTime   time.Time `json:"startTime" dynamodbav:"StartTime"`   // 開始時間
    EndTime     time.Time `json:"endTime" dynamodbav:"EndTime"`       // 終了時間
    Location    string    `json:"location" dynamodbav:"Location"`     // 場所
    AllDay      bool      `json:"allDay" dynamodbav:"AllDay"`         // 終日フラグ
}
