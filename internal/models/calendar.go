package models

type Calendar struct {
	CalendarID  string  `json:"calendarId" dynamodbav:"CalendarID"`   // カレンダーのID
	Name        string  `json:"name" dynamodbav:"Name"`               // カレンダー名
	IsPublic    bool    `json:"isPublic" dynamodbav:"IsPublic"`       // 公開フラグ
	OwnerUserID string  `json:"ownerUserId" dynamodbav:"OwnerUserID"` // オーナーのユーザーID
	Users       []User  `json:"users" dynamodbav:"Users"`             // 共有ユーザーのIDリスト
	Events      []Event `json:"events,omitempty"`                     // カレンダー内のイベント
}
