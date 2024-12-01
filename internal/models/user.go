package models

type User struct {
	UniqueID    string // 固有ID
	DisplayName string // 表示名
	Email       string
	Password    string // ハッシュ化されたパスワード
}
