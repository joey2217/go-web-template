package model

// Post 帖子实体。
type Post struct {
	Base
	Title   string `gorm:"type:varchar(255);not null;index" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
	Status  int16  `gorm:"type:smallint;not null;default:0;index" json:"status"`
	Author  string `gorm:"type:varchar(64);not null;index" json:"author"`
}

// 帖子状态取值。
const (
	PostStatusDraft     int16 = 0 // 草稿
	PostStatusPublished int16 = 1 // 已发布
)

// TableName 指定表名。
func (Post) TableName() string { return "posts" }
