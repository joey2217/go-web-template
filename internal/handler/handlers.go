package handler

import (
	"gorm.io/gorm"

	"github.com/example/go-web-template/internal/repository"
	"github.com/example/go-web-template/internal/service"
	"github.com/example/go-web-template/pkg/cache"
)

// Handlers 聚合所有 HTTP 处理器。
//
// 由 NewHandlers 统一构造（内部完成 repository → service → handler 串联），
// 路由层只需持有一个 Handlers 引用即可访问所有端点，
// 避免 router.Deps 字段随业务增长而膨胀，也消除 app 层显式串联依赖链的样板。
//
// 单元测试时仍可直接使用各 handler 的 NewXxxHandler(svc) 注入 mock service。
type Handlers struct {
	Post *PostHandler
	User *UserHandler
}

// NewHandlers 一站式构造所有 HTTP 处理器及其底层依赖链。
// c 可以是 cache.Noop{}（未启用缓存时），service 侧无需判 nil。
func NewHandlers(db *gorm.DB, c cache.Cache) *Handlers {
	return &Handlers{
		Post: newPostHandler(db, c),
		User: newUserHandler(db, c),
	}
}

func newPostHandler(db *gorm.DB, c cache.Cache) *PostHandler {
	repo := repository.NewPostRepository(db)
	svc := service.NewPostService(repo, c)
	return NewPostHandler(svc)
}

func newUserHandler(db *gorm.DB, c cache.Cache) *UserHandler {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, c)
	return NewUserHandler(svc)
}
