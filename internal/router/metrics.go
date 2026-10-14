package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// registerMetrics 挂载 Prometheus 抓取端点。
// OpenMetrics 格式开启——Prometheus 2.0+ 原生支持（带 _created 序列等）。
//
// 同时注册 GET 和 HEAD：Prometheus 偶尔用 HEAD 探活，纯 GET 会返回 405/404。
//
// 路径默认 /metrics，可通过 config.metrics.path 改。生产建议：
//   - 内网暴露 / 由 Prometheus 直接抓
//   - 公网暴露时加 basic auth / IP 白名单（这里简化为内网约定）
func registerMetrics(r *gin.Engine, reg *prometheus.Registry, path string) {
	handler := gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		Registry:          reg,
		EnableOpenMetrics: true,
	}))
	r.GET(path, handler)
	r.HEAD(path, handler)
}
