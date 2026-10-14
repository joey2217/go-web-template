package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/example/go-web-template/internal/dto"
	"github.com/example/go-web-template/internal/model"
	"github.com/example/go-web-template/internal/repository/gen"
)

// PostRepository 定义帖子数据访问的消费方接口。
//
// 底层实现由 GORM CLI 生成的代码提供（见 ./gen 包）；
// Create 通过 GORM 原生调用完成，以保留时间戳自动填充与 BeforeCreate 钩子。
type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
	List(ctx context.Context, opts dto.ListOptions) ([]model.Post, int64, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// NewPostRepository 创建帖子仓储实现。
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

type postRepository struct {
	db *gorm.DB
}

func (r *postRepository) Create(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	return gen.Post[model.Post](r.db).GetByID(ctx, id)
}

func (r *postRepository) List(ctx context.Context, opts dto.ListOptions) ([]model.Post, int64, error) {
	q := gen.Post[model.Post](r.db)
	total, err := q.CountPosts(ctx, opts)
	if err != nil {
		return nil, 0, err
	}
	posts, err := q.ListPosts(ctx, opts)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (r *postRepository) Update(ctx context.Context, post *model.Post) error {
	// GORM CLI 生成的 UpdateByID 不检查受影响行数，这里预先校验存在性以保留 404 语义。
	if _, err := r.GetByID(ctx, post.ID); err != nil {
		return err
	}
	return gen.Post[model.Post](r.db).UpdateByID(ctx, post)
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.GetByID(ctx, id); err != nil {
		return err
	}
	return gen.Post[model.Post](r.db).DeleteByID(ctx, id)
}
