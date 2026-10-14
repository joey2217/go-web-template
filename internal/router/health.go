package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// registerHealth 注册存活与就绪探针。
func registerHealth(r *gin.Engine, deps Deps) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		checks := gin.H{}
		ready := true

		if sqlDB, err := deps.DB.DB(); err != nil {
			checks["postgres"] = err.Error()
			ready = false
		} else if err := sqlDB.PingContext(ctx); err != nil {
			checks["postgres"] = err.Error()
			ready = false
		} else {
			checks["postgres"] = "ok"
		}

		if deps.Redis != nil {
			if err := deps.Redis.Ping(ctx).Err(); err != nil {
				checks["redis"] = err.Error()
				ready = false
			} else {
				checks["redis"] = "ok"
			}
		} else {
			checks["redis"] = "disabled"
		}

		status := http.StatusOK
		state := "ok"
		if !ready {
			status = http.StatusServiceUnavailable
			state = "unavailable"
		}
		c.JSON(status, gin.H{"status": state, "checks": checks})
	})
}
