package model

// User 用户实体。
type User struct {
	Base
	Username string `gorm:"type:varchar(64);not null;uniqueIndex:idx_users_username" json:"username"`
	Email    string `gorm:"type:varchar(128);not null;uniqueIndex:idx_users_email" json:"email"`
	Nickname string `gorm:"type:varchar(64)" json:"nickname"`
	Status   int16  `gorm:"type:smallint;not null;default:0" json:"status"`
}

// 用户状态取值。
const (
	UserStatusActive   int16 = 0 // 正常
	UserStatusDisabled int16 = 1 // 禁用
)

// TableName 指定表名。
func (User) TableName() string { return "users" }
