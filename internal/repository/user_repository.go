package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/example/go-web-template/internal/dto"
	"github.com/example/go-web-template/internal/model"
	"github.com/example/go-web-template/internal/repository/gen"
)

// UserRepository 定义用户数据访问的消费方接口。
//
// 底层实现由 GORM CLI 生成的代码提供（见 ./gen 包）；
// Create 通过 GORM 原生调用完成，以保留时间戳自动填充与 BeforeCreate 钩子。
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	List(ctx context.Context, opts dto.ListOptions) ([]model.User, int64, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// NewUserRepository 创建用户仓储实现。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return gen.User[model.User](r.db).GetByID(ctx, id)
}

func (r *userRepository) List(ctx context.Context, opts dto.ListOptions) ([]model.User, int64, error) {
	q := gen.User[model.User](r.db)
	total, err := q.CountUsers(ctx, opts)
	if err != nil {
		return nil, 0, err
	}
	users, err := q.ListUsers(ctx, opts)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// GORM CLI 生成的 UpdateByID 不检查受影响行数，这里预先校验存在性以保留 404 语义。
	if _, err := r.GetByID(ctx, user.ID); err != nil {
		return err
	}
	return gen.User[model.User](r.db).UpdateByID(ctx, user)
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.GetByID(ctx, id); err != nil {
		return err
	}
	return gen.User[model.User](r.db).DeleteByID(ctx, id)
}
