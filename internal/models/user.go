package models

type User struct {
	UserID      string `json:"userId" dynamodbav:"UserID"`           // ユーザーID
	DisplayName string `json:"displayName" dynamodbav:"DisplayName"` // 表示名
	Email       string `json:"email" dynamodbav:"Email"`             // メールアドレス
	Password    string `json:"password" dynamodbav:"Password"`       // パスワード
	AccessLevel string `json:"accessLevel" dynamodbav:"AccessLevel"` // 権限（OWNER/EDITOR/VIEWER）
}
