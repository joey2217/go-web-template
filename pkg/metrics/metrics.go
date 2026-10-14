// Package metrics 提供 Prometheus 指标注册与采集。
//
// 三层结构：
//   - HTTP:    中间件自动采集（请求数 / 延迟 / 在飞）
//   - Cache:   装饰器 cache.NewInstrumented() 采集（操作数 / 命中率）
//   - Runtime: 自动注册（go_* / process_*）+ DB 连接池（db_pool_*）
//
// 关键约束：
//   - 用独立 prometheus.Registry 而非 DefaultRegisterer：避免测试/多实例污染
//   - HTTP 标签里的 route 用 c.FullPath()（/users/:id）而非 URL.Path，防基数爆炸
package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// HTTPRequestsTotal 按 method/route/status 统计请求数。
var HTTPRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests processed.",
	},
	[]string{"method", "route", "status"},
)

// HTTPRequestDuration 按 method/route/status 统计请求延迟（秒）。
// Buckets 覆盖 5ms~10s，足够 Web API 场景。
var HTTPRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "route", "status"},
)

// HTTPRequestsInFlight 当前在飞请求数（gauge）。
var HTTPRequestsInFlight = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Number of HTTP requests currently being processed.",
	},
)

// CacheOps 按 op/result 统计缓存操作。
//   - op:     get | set | delete | delete_pattern
//   - result: hit | miss | ok | error
var CacheOps = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "cache_operations_total",
		Help: "Total number of cache operations by op and result.",
	},
	[]string{"op", "result"},
)

// NewRegistry 构造独立的 Prometheus 注册表，注册所有自定义 + 标准运行时指标。
func NewRegistry() *prometheus.Registry {
	r := prometheus.NewRegistry()
	r.MustRegister(
		// 自定义
		HTTPRequestsTotal,
		HTTPRequestDuration,
		HTTPRequestsInFlight,
		CacheOps,
		// 标准：Go runtime（goroutine / GC / memstats）+ 进程（CPU / RSS / FDs）
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	return r
}

// DBPoolCollector 把 sql.DBStats 暴露为 Prometheus 指标。
// 实现 prometheus.Collector 接口：每次 scrape 拉取最新快照（零开销、无后台 goroutine）。
type DBPoolCollector struct {
	db         *sql.DB
	openDesc   *prometheus.Desc
	inUseDesc  *prometheus.Desc
	idleDesc   *prometheus.Desc
	waitDesc   *prometheus.Desc
	waitDurSec *prometheus.Desc
}

// NewDBPoolCollector 创建 DB 连接池采集器。db=nil 时返回 nil。
func NewDBPoolCollector(db *sql.DB) *DBPoolCollector {
	if db == nil {
		return nil
	}
	return &DBPoolCollector{
		db:         db,
		openDesc:   prometheus.NewDesc("db_pool_open_connections", "Open connections (in_use + idle).", nil, nil),
		inUseDesc:  prometheus.NewDesc("db_pool_in_use_connections", "Connections currently in use.", nil, nil),
		idleDesc:   prometheus.NewDesc("db_pool_idle_connections", "Idle connections.", nil, nil),
		waitDesc:   prometheus.NewDesc("db_pool_wait_count_total", "Total connections waited for.", nil, nil),
		waitDurSec: prometheus.NewDesc("db_pool_wait_duration_seconds_total", "Total time blocked waiting for a new connection.", nil, nil),
	}
}

func (c *DBPoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.openDesc
	ch <- c.inUseDesc
	ch <- c.idleDesc
	ch <- c.waitDesc
	ch <- c.waitDurSec
}

func (c *DBPoolCollector) Collect(ch chan<- prometheus.Metric) {
	s := c.db.Stats()
	ch <- prometheus.MustNewConstMetric(c.openDesc, prometheus.GaugeValue, float64(s.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.inUseDesc, prometheus.GaugeValue, float64(s.InUse))
	ch <- prometheus.MustNewConstMetric(c.idleDesc, prometheus.GaugeValue, float64(s.Idle))
	ch <- prometheus.MustNewConstMetric(c.waitDesc, prometheus.CounterValue, float64(s.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.waitDurSec, prometheus.CounterValue, s.WaitDuration.Seconds())
}
