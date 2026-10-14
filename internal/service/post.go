// Package service 实现业务逻辑层。
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/example/go-web-template/internal/dto"
	"github.com/example/go-web-template/internal/model"
	"github.com/example/go-web-template/internal/repository"
	"github.com/example/go-web-template/pkg/apperror"
	"github.com/example/go-web-template/pkg/cache"
	"github.com/example/go-web-template/pkg/response"
)

// PostService 定义帖子业务接口。
type PostService interface {
	Create(ctx context.Context, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	Get(ctx context.Context, id string) (*dto.PostResponse, error)
	List(ctx context.Context, query *dto.ListPostsQuery) (*response.PageData, error)
	Update(ctx context.Context, id string, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	Delete(ctx context.Context, id string) error
}

type postService struct {
	repo  repository.PostRepository
	cache cache.Cache
}

// NewPostService 创建帖子业务实现。
// c 可以是 cache.Noop{}（未启用缓存时），方法内无需判 nil。
func NewPostService(repo repository.PostRepository, c cache.Cache) PostService {
	return &postService{repo: repo, cache: c}
}

const (
	postCachePrefix = "post:"
	postCacheTTL    = 5 * time.Minute
)

func postCacheKey(uid uuid.UUID) string { return postCachePrefix + uid.String() }

func (s *postService) Create(ctx context.Context, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	post := &model.Post{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
		Author:  req.Author,
	}
	if err := s.repo.Create(ctx, post); err != nil {
		return nil, apperror.Internal("创建帖子失败").WithCause(err)
	}
	return dto.NewPostResponse(post), nil
}

// Get 按 ID 查询帖子，cache-aside：先查缓存，未命中回源并回填。
// 缓存失败 fail-open，降级到 DB，不影响业务可用性。
func (s *postService) Get(ctx context.Context, id string) (*dto.PostResponse, error) {
	uid, err := parseID(id)
	if err != nil {
		return nil, err
	}
	key := postCacheKey(uid)

	if cached, ok, err := cache.Get[dto.PostResponse](ctx, s.cache, key); err == nil && ok {
		return &cached, nil
	}

	post, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, mapNotFound(err, "帖子不存在")
	}
	resp := dto.NewPostResponse(post)
	_ = cache.Set[dto.PostResponse](ctx, s.cache, key, *resp, postCacheTTL)
	return resp, nil
}

func (s *postService) List(ctx context.Context, query *dto.ListPostsQuery) (*response.PageData, error) {
	query.Normalize()
	posts, total, err := s.repo.List(ctx, dto.ListOptions{
		Keyword: query.Keyword,
		Status:  query.Status,
		Offset:  query.Offset(),
		Limit:   query.PageSize,
	})
	if err != nil {
		return nil, apperror.Internal("查询帖子列表失败").WithCause(err)
	}
	return &response.PageData{
		List:     dto.NewPostResponseList(posts),
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func (s *postService) Update(ctx context.Context, id string, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	uid, err := parseID(id)
	if err != nil {
		return nil, err
	}
	post, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, mapNotFound(err, "帖子不存在")
	}

	post.Title = req.Title
	post.Content = req.Content
	post.Status = req.Status
	post.Author = req.Author
	if err := s.repo.Update(ctx, post); err != nil {
		return nil, mapNotFound(err, "帖子不存在")
	}

	// 回查以获取数据库生成的最新时间戳。
	updated, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, mapNotFound(err, "帖子不存在")
	}

	// 写后失效缓存（best-effort）。
	_ = s.cache.Delete(ctx, postCacheKey(uid))
	return dto.NewPostResponse(updated), nil
}

func (s *postService) Delete(ctx context.Context, id string) error {
	uid, err := parseID(id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, uid); err != nil {
		return mapNotFound(err, "帖子不存在")
	}
	_ = s.cache.Delete(ctx, postCacheKey(uid))
	return nil
}

// parseID 校验并解析路径参数中的 UUID。
func parseID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperror.BadRequest("无效的 ID")
	}
	return uid, nil
}

// mapNotFound 将 GORM 的记录不存在错误映射为业务错误，其余归为内部错误。
func mapNotFound(err error, message string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NotFound(message)
	}
	return apperror.Internal("数据库操作失败").WithCause(err)
}
