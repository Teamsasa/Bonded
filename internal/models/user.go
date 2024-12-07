package models

type User struct {
	UserID      string `json:"userId" dynamodbav:"UserID"`           // ユーザーID
	DisplayName string `json:"displayName" dynamodbav:"DisplayName"` // 表示名
	AccessLevel string `json:"accessLevel" dynamodbav:"AccessLevel"` // 権限（OWNER/EDITOR/VIEWER）
}
