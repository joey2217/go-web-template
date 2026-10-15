// Package router 负责路由树与全局中间件装配。
package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/example/go-web-template/internal/config"
	"github.com/example/go-web-template/internal/handler"
	"github.com/example/go-web-template/internal/middleware"
	"github.com/example/go-web-template/pkg/apperror"
	"github.com/example/go-web-template/pkg/metrics"
	"github.com/example/go-web-template/pkg/response"
)

// Deps 是路由装配所需的依赖。
type Deps struct {
	Config   *config.Config
	Logger   *zap.Logger
	DB       *gorm.DB
	Redis    *goredis.Client
	Handlers *handler.Handlers
	Metrics  *prometheus.Registry // 传 nil 表示不挂 /metrics 端点
}

// New 构建 Gin 引擎并装配全局中间件与路由。
func New(deps Deps) *gin.Engine {
	if deps.Config.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	} else if !deps.Config.App.Debug {
		gin.SetMode(gin.TestMode)
	}

	r := gin.New()

	// 中间件顺序：
	//   1. HTTPMetrics 最外层——记录整个请求的真实耗时（含 Recovery 兜底时间），
	//      即使 handler panic，defer 仍能拿到 Recovery 写回的 500 状态码
	//   2. Recovery 兜住后续所有 panic
	//   3. RequestID 在 Recovery 之前调用，其 defer 后于 Recovery 的 defer 执行（LIFO），
	//      因此 panic 日志能携带 request_id
	r.Use(metrics.HTTPMetrics())
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

	if deps.Metrics != nil {
		registerMetrics(
			r,
			deps.Metrics,
			deps.Config.Metrics.Path,
			deps.Config.Metrics.Username,
			deps.Config.Metrics.Password,
		)
	}

	registerAPI(r, deps)

	// 未匹配路由走统一 ErrorBody 信封（与 handler 层资源 404 一致）。
	// 资源级的 404（GET /users/:id 找不到）由 handler 抛 apperror，不走这里。
	//
	// gin v1 的 trigger404 不会合并 engine.Handlers（r.Use 注册的全局中间件），
	// NoRoute 存的就是原始 handlers，不含 Recovery / RequestID / AccessLog / HTTPMetrics。
	// 这里手动 append(r.Handlers, ...) 让 X-Request-Id 响应头、访问日志、指标、panic 恢复
	// 对未匹配路由同样生效——否则 /typo 这类请求连 request_id 都没有，指标也不计数。
	r.NoRoute(append(r.Handlers, func(c *gin.Context) {
		response.Fail(c, apperror.NotFound(
			fmt.Sprintf("路由不存在: %s %s", c.Request.Method, c.Request.URL.Path),
		))
	})...)

	return r
}
