package handler

import (
	"gorm.io/gorm"

	"github.com/example/go-web-template/internal/repository"
	"github.com/example/go-web-template/internal/service"
)

// Handlers 聚合所有 HTTP 处理器。
//
// 由 NewHandlers 统一构造（内部完成 repository → service → handler 装配），
// 路由层只需持有一个 Handlers 引用即可访问所有端点，
// 避免 router.Deps 字段随业务增长而膨胀，也消除 app 层显式串联依赖链的样板。
//
// 单元测试时仍可直接使用各 handler 的 NewXxxHandler(svc) 注入 mock service。
type Handlers struct {
	Post *PostHandler
	User *UserHandler
}

// NewHandlers 一站式构造所有 HTTP 处理器及其底层依赖链。
func NewHandlers(db *gorm.DB) *Handlers {
	return &Handlers{
		Post: newPostHandler(db),
		User: newUserHandler(db),
	}
}

func newPostHandler(db *gorm.DB) *PostHandler {
	repo := repository.NewPostRepository(db)
	svc := service.NewPostService(repo)
	return NewPostHandler(svc)
}

func newUserHandler(db *gorm.DB) *UserHandler {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	return NewUserHandler(svc)
}
