package models

type Calendar struct {
	CalendarID  string  `json:"calendarId,omitempty" dynamodbav:"CalendarID"`   // カレンダーのID
	SortKey     string  `json:"sortKey,omitempty" dynamodbav:"SortKey"`         // ソートキー
	Name        string  `json:"name" dynamodbav:"Name"`                         // カレンダー名
	IsPublic    *bool    `json:"isPublic" dynamodbav:"IsPublic"`                 // 公開フラグ
	OwnerUserID string  `json:"ownerUserId,omitempty" dynamodbav:"OwnerUserID"` // オーナーのユーザーID
	Users       []User  `json:"users,omitempty" dynamodbav:"Users"`             // 共有ユーザーのIDリスト
	Events      []Event `json:"events,omitempty"`                               // カレンダー内のイベント
}
