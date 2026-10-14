package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPMetrics 记录每个请求的 计数 / 延迟 / 在飞 三个维度。
//
// 关键点：
//   - c.FullPath() 返回路由模板（如 /api/v1/users/:id），不是真实 URL——
//     防止被 /random/xxx 这种扫描把 Prometheus 标签基数撑爆
//   - FullPath() == "" 表示 404 路由（未匹配），归到 __not_found__ 桶（保留扫描可见性）
//   - 记录代码放在 defer 里：即使 c.Next() 内部 panic、Recovery 兜住写 500，
//     defer 仍能拿到最终状态码并完成计数
//   - 顺序：本中间件应放在最外层（Recovery 之前），让 HTTPRequestsInFlight 包含
//     Recovery 兜底时间、duration 反映完整请求耗时
func HTTPMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		HTTPRequestsInFlight.Inc()
		defer func() {
			HTTPRequestsInFlight.Dec()

			route := c.FullPath()
			if route == "" {
				route = "__not_found__"
			}
			status := strconv.Itoa(c.Writer.Status())
			method := c.Request.Method

			HTTPRequestsTotal.WithLabelValues(method, route, status).Inc()
			HTTPRequestDuration.WithLabelValues(method, route, status).Observe(time.Since(start).Seconds())
		}()
		c.Next()
	}
}
