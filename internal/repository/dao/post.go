// Package dao 集中存放 GORM CLI 数据访问接口（仅用于代码生成，不在运行时被业务层直接引用）。
//
// 业务层的仓储接口与实现在 internal/repository 包（本包的父包）中。
// 这样隔离是因为 GORM CLI 工具要求被处理的接口方法必须带 SQL 模板注释，
// 而消费方仓储接口（PostRepository / UserRepository）的方法没有这类注释，
// 混在一起会导致生成失败。
//
// 运行 `make gen` 重新生成实现到 ../gen 包。
package dao

import (
	"context"

	"github.com/google/uuid"

	"github.com/example/go-web-template/internal/dto"
	"github.com/example/go-web-template/internal/model"
)

// Post 帖子数据访问接口（GORM CLI）。
type Post interface {
	// GetByID 根据主键查询单个帖子（排除软删）。
	//
	//	SELECT * FROM @@table WHERE id = @id AND deleted_at IS NULL
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)

	// UpdateByID 更新帖子。
	//
	//	UPDATE @@table
	//	{{set}}
	//	  updated_at = NOW(),
	//	  title = @p.Title,
	//	  content = @p.Content,
	//	  status = @p.Status,
	//	  author = @p.Author
	//	{{end}}
	//	WHERE id = @p.ID AND deleted_at IS NULL
	UpdateByID(ctx context.Context, p *model.Post) error

	// DeleteByID 软删帖子。
	//
	//	UPDATE @@table SET deleted_at = NOW() WHERE id = @id AND deleted_at IS NULL
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// ListPosts 分页查询帖子。
	//
	//	SELECT * FROM @@table
	//	{{where}}
	//	  deleted_at IS NULL
	//	  {{if opts.Keyword != ""}} AND (title ILIKE @opts.Keyword OR content ILIKE @opts.Keyword) {{end}}
	//	  {{if opts.Status != nil}} AND status = @opts.Status {{end}}
	//	{{end}}
	//	ORDER BY created_at DESC
	//	LIMIT @opts.Limit OFFSET @opts.Offset
	ListPosts(ctx context.Context, opts dto.ListOptions) ([]model.Post, error)

	// CountPosts 统计符合条件的帖子数量。
	//
	//	SELECT COUNT(*) FROM @@table
	//	{{where}}
	//	  deleted_at IS NULL
	//	  {{if opts.Keyword != ""}} AND (title ILIKE @opts.Keyword OR content ILIKE @opts.Keyword) {{end}}
	//	  {{if opts.Status != nil}} AND status = @opts.Status {{end}}
	//	{{end}}
	CountPosts(ctx context.Context, opts dto.ListOptions) (int64, error)
}
