package dto

import (
	"time"

	"github.com/example/go-web-template/internal/model"
)

// CreatePostRequest 创建帖子入参。
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=255" example:"Hello World"`
	Content string `json:"content" binding:"required" example:"post content"`
	Author  string `json:"author" binding:"required,max=64" example:"alice"`
	Status  int16  `json:"status" binding:"omitempty,oneof=0 1" example:"0"`
}

// UpdatePostRequest 更新帖子入参（全量更新）。
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,max=255" example:"Hello World"`
	Content string `json:"content" binding:"required" example:"post content"`
	Author  string `json:"author" binding:"required,max=64" example:"alice"`
	Status  int16  `json:"status" binding:"omitempty,oneof=0 1" example:"1"`
}

// ListPostsQuery 帖子列表查询入参。
type ListPostsQuery struct {
	PageQuery
	Keyword string `form:"keyword" binding:"omitempty,max=255"`
	Status  *int16 `form:"status" binding:"omitempty,oneof=0 1"`
}

// ListOptions 是传递给数据访问层的列表查询选项（与 GORM CLI 模板共用的结构体参数）。
// 独立于 ListPostsQuery 以避免在生成代码与业务 DTO 之间形成导入循环。
type ListOptions struct {
	Keyword string
	Status  *int16
	Offset  int
	Limit   int
}

// PostResponse 帖子响应体。
type PostResponse struct {
	ID        string    `json:"id" example:"0191f0a0-0000-7000-8000-000000000000"`
	Title     string    `json:"title" example:"Hello World"`
	Content   string    `json:"content" example:"post content"`
	Status    int16     `json:"status" example:"0"`
	Author    string    `json:"author" example:"alice"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPostResponse 将实体转换为响应体。
func NewPostResponse(p *model.Post) *PostResponse {
	if p == nil {
		return nil
	}
	return &PostResponse{
		ID:        p.ID.String(),
		Title:     p.Title,
		Content:   p.Content,
		Status:    p.Status,
		Author:    p.Author,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// NewPostResponseList 批量转换实体为响应体。
func NewPostResponseList(posts []model.Post) []PostResponse {
	list := make([]PostResponse, 0, len(posts))
	for i := range posts {
		list = append(list, *NewPostResponse(&posts[i]))
	}
	return list
}
