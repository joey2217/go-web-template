package dto

import (
	"time"

	"github.com/example/go-web-template/internal/model"
)

// CreateUserRequest 创建用户入参。
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64" example:"alice"`
	Email    string `json:"email" binding:"required,email,max=128" example:"alice@example.com"`
	Nickname string `json:"nickname" binding:"omitempty,max=64" example:"Alice"`
	Status   *int16 `json:"status" binding:"omitempty,oneof=0 1" example:"0"`
}

// UpdateUserRequest 更新用户入参（局部更新，nil 字段不修改）。
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,max=64" example:"Alice"`
	Email    *string `json:"email" binding:"omitempty,email,max=128" example:"alice@example.com"`
	Status   *int16  `json:"status" binding:"omitempty,oneof=0 1" example:"0"`
}

// ListUsersQuery 用户列表查询入参。
type ListUsersQuery struct {
	PageQuery
	Keyword string `form:"keyword" binding:"omitempty,max=128"`
	Status  *int16 `form:"status" binding:"omitempty,oneof=0 1"`
}

// UserResponse 用户响应体。
type UserResponse struct {
	ID        string    `json:"id" example:"0191f0a0-0000-7000-8000-000000000000"`
	Username  string    `json:"username" example:"alice"`
	Email     string    `json:"email" example:"alice@example.com"`
	Nickname  string    `json:"nickname" example:"Alice"`
	Status    int16     `json:"status" example:"0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserResponse 将实体转换为响应体。
func NewUserResponse(u *model.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// NewUserResponseList 批量转换实体为响应体。
func NewUserResponseList(users []model.User) []UserResponse {
	list := make([]UserResponse, 0, len(users))
	for i := range users {
		list = append(list, *NewUserResponse(&users[i]))
	}
	return list
}
