// Package router 负责路由树与全局中间件装配。
package router

import (
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/example/go-web-template/internal/config"
	"github.com/example/go-web-template/internal/handler"
	"github.com/example/go-web-template/internal/middleware"
)

// Deps 是路由装配所需的依赖。
type Deps struct {
	Config   *config.Config
	Logger   *zap.Logger
	DB       *gorm.DB
	Redis    *goredis.Client
	Handlers *handler.Handlers
}

// New 构建 Gin 引擎并装配全局中间件与路由。
func New(deps Deps) *gin.Engine {
	if deps.Config.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	} else if !deps.Config.App.Debug {
		gin.SetMode(gin.TestMode)
	}

	r := gin.New()

	// 顺序即执行顺序：Recovery 在最外层以保证能兜住后续所有 panic。
	// RequestID 先于 Recovery 的 defer 执行，因此 panic 日志也能携带 request_id。
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID(deps.Logger))
	r.Use(middleware.AccessLog())
	// r.Use(middleware.CORS(deps.Config.CORS))
	r.Use(middleware.ErrorHandler())

	registerHealth(r, deps)
	registerStatic(r)

	if deps.Config.Swagger.Enabled {
		registerSwagger(r)
	}

	registerAPI(r, deps)

	return r
}
