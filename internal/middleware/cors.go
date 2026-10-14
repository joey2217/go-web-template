package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/example/go-web-template/internal/config"
)

// CORS 返回跨域中间件；当白名单包含 "*" 时放行所有来源。
func CORS(cfg config.CORS) gin.HandlerFunc {
	origins := cfg.AllowOrigins
	allowAll := false
	for _, o := range origins {
		if o == "*" {
			allowAll = true
			break
		}
	}
	if len(origins) == 0 {
		allowAll = true
	}

	corsCfg := cors.Config{
		AllowAllOrigins:  allowAll,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Id"},
		ExposeHeaders:    []string{"X-Request-Id"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	if !allowAll {
		corsCfg.AllowOrigins = origins
	}
	return cors.New(corsCfg)
}
