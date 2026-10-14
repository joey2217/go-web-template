package router

import (
	"github.com/gin-gonic/gin"

	"github.com/example/go-web-template/internal/middleware"
)

// registerAPI 注册 /api/v1 下的业务路由。
func registerAPI(r *gin.Engine, deps Deps) {
	v1 := r.Group("/api/v1")

	// 受保护分组：先限流，再鉴权（鉴权默认放行）。
	protected := v1.Group("")
	protected.Use(middleware.RateLimit(deps.Config.RateLimit, deps.Redis, deps.Logger))
	protected.Use(middleware.Auth())

	h := deps.Handlers

	posts := protected.Group("/posts")
	{
		posts.POST("", h.Post.Create)
		posts.GET("", h.Post.List)
		posts.GET("/:id", h.Post.Get)
		posts.PUT("/:id", h.Post.Update)
		posts.DELETE("/:id", h.Post.Delete)
	}

	users := protected.Group("/users")
	{
		users.POST("", h.User.Create)
		users.GET("", h.User.List)
		users.GET("/:id", h.User.Get)
		users.PUT("/:id", h.User.Update)
		users.DELETE("/:id", h.User.Delete)
	}
}
