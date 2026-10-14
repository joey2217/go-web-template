// Package model 定义数据库实体（GORM 模型）。
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base 是所有实体共用的基础字段：UUID 主键 + 时间戳 + 软删除。
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 在插入前生成 UUID v7（按时间有序），失败时回退到 v4。
func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.ID != uuid.Nil {
		return nil
	}
	id, err := uuid.NewV7()
	if err != nil {
		id = uuid.New()
	}
	b.ID = id
	return nil
}
