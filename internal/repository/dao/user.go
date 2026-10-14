package dao

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/go-web-template/internal/dto"
	"github.com/example/go-web-template/internal/model"
)

// User 用户数据访问接口（GORM CLI）。
type User interface {
	// GetByID 根据主键查询单个用户（排除软删）。
	//
	//	SELECT * FROM @@table WHERE id = @id AND deleted_at IS NULL
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// UpdateByID 更新用户。
	//
	//	UPDATE @@table
	//	{{set}}
	//	  updated_at = NOW(),
	//	  username = @u.Username,
	//	  email = @u.Email,
	//	  nickname = @u.Nickname,
	//	  status = @u.Status
	//	{{end}}
	//	WHERE id = @u.ID AND deleted_at IS NULL
	UpdateByID(ctx context.Context, u *model.User) error

	// DeleteByID 软删用户。
	//
	//	UPDATE @@table SET deleted_at = NOW() WHERE id = @id AND deleted_at IS NULL
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// ListUsers 分页查询用户。
	//
	//	SELECT * FROM @@table
	//	{{where}}
	//	  deleted_at IS NULL
	//	  {{if opts.Keyword != ""}} AND (username ILIKE @opts.Keyword OR email ILIKE @opts.Keyword OR nickname ILIKE @opts.Keyword) {{end}}
	//	  {{if opts.Status != nil}} AND status = @opts.Status {{end}}
	//	{{end}}
	//	ORDER BY created_at DESC
	//	LIMIT @opts.Limit OFFSET @opts.Offset
	ListUsers(ctx context.Context, opts dto.ListOptions) ([]model.User, error)

	// CountUsers 统计符合条件的用户数量。
	//
	//	SELECT COUNT(*) FROM @@table
	//	{{where}}
	//	  deleted_at IS NULL
	//	  {{if opts.Keyword != ""}} AND (username ILIKE @opts.Keyword OR email ILIKE @opts.Keyword OR nickname ILIKE @opts.Keyword) {{end}}
	//	  {{if opts.Status != nil}} AND status = @opts.Status {{end}}
	//	{{end}}
	CountUsers(ctx context.Context, opts dto.ListOptions) (int64, error)
}
