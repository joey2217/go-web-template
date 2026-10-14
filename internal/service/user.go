// Package service 实现业务逻辑层。
package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/example/go-web-template/internal/dto"
	"github.com/example/go-web-template/internal/model"
	"github.com/example/go-web-template/internal/repository"
	"github.com/example/go-web-template/pkg/apperror"
	"github.com/example/go-web-template/pkg/response"
)

// UserService 定义用户业务接口。
type UserService interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	Get(ctx context.Context, id string) (*dto.UserResponse, error)
	List(ctx context.Context, query *dto.ListUsersQuery) (*response.PageData, error)
	Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService 创建用户业务实现。
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	status := int16(model.UserStatusActive)
	if req.Status != nil {
		status = *req.Status
	}
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   status,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("用户名或邮箱已存在").WithCause(err)
		}
		return nil, apperror.Internal("创建用户失败").WithCause(err)
	}
	return dto.NewUserResponse(user), nil
}

func (s *userService) Get(ctx context.Context, id string) (*dto.UserResponse, error) {
	uid, err := parseID(id)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, mapNotFound(err, "用户不存在")
	}
	return dto.NewUserResponse(user), nil
}

func (s *userService) List(ctx context.Context, query *dto.ListUsersQuery) (*response.PageData, error) {
	query.Normalize()
	users, total, err := s.repo.List(ctx, dto.ListOptions{
		Keyword: query.Keyword,
		Status:  query.Status,
		Offset:  query.Offset(),
		Limit:   query.PageSize,
	})
	if err != nil {
		return nil, apperror.Internal("查询用户列表失败").WithCause(err)
	}
	return &response.PageData{
		List:     dto.NewUserResponseList(users),
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func (s *userService) Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	uid, err := parseID(id)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, mapNotFound(err, "用户不存在")
	}

	// 局部更新：仅当字段非 nil 时覆盖。
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.repo.Update(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("邮箱已被其他用户使用").WithCause(err)
		}
		return nil, mapNotFound(err, "用户不存在")
	}

	// 回查以获取数据库生成的最新时间戳。
	updated, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, mapNotFound(err, "用户不存在")
	}
	return dto.NewUserResponse(updated), nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	uid, err := parseID(id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, uid); err != nil {
		return mapNotFound(err, "用户不存在")
	}
	return nil
}
