# go-web-template

一个开箱即用的 **Go RESTful API 服务模板**：Gin + GORM + PostgreSQL，内置鉴权中间件（占位放行）、Redis 分布式限流、Swagger 文档、统一响应/错误处理、结构化日志、配置管理、数据库迁移与 Docker 一键启动。

---

## 1. 项目定位

| 项 | 决策 |
| --- | --- |
| 类型 | 纯 RESTful API 服务（只输出 JSON，无前端页面） |
| 使用方式 | 作为后端 / 微服务模板复制后按业务扩展 |
| 语言 | **Go 1.26** |
| 鉴权 | **预留 `Auth` 中间件，默认放行**（无配置项，代码内 TODO 待实现） |
| 限流 | **基于 Redis 的分布式限流** |
| API 文档 | **Swagger（swaggo 注解生成 OpenAPI + Swagger UI）** |
| 配置 | Viper 加载 YAML + 环境变量覆盖，多环境 |
| 日志 | Zap 结构化日志 |
| 数据库 | PostgreSQL + GORM，golang-migrate 做版本化迁移 |
| 交付 | Dockerfile + docker-compose（app + pgsql + redis 一键起） |
| 示例业务 | `posts`（title / content / status / author） |

**设计目标**：分层清晰、依赖单向、易测试、易扩展。

---

## 2. 技术选型

| 类别 | 选型 | 说明 |
| --- | --- | --- |
| Web 框架 | `github.com/gin-gonic/gin` | 生态成熟、中间件丰富 |
| ORM | `gorm.io/gorm` + `gorm.io/driver/postgres` | |
| 配置 | `github.com/spf13/viper` | YAML + 环境变量自动覆盖 |
| 日志 | `go.uber.org/zap` | 结构化、高性能；含 GORM logger 适配 |
| 参数校验 | `github.com/go-playground/validator/v10` | 与 Gin binding 集成，含中文错误翻译 |
| UUID | `github.com/google/uuid` | 主键（v7 优先） |
| Redis 客户端 | `github.com/redis/go-redis/v9` | 限流 / 缓存底层 |
| 限流 | `github.com/go-redis/redis_rate/v10` | 基于 Redis 的 GCRA 分布式限流 |
| 数据库迁移 | `github.com/golang-migrate/migrate/v4` | 版本化 SQL 迁移 |
| CORS | `github.com/gin-contrib/cors` | 可配置白名单 |
| Request ID | `github.com/gin-contrib/requestid` | 贯穿日志与响应头 |
| 日志中间件 | `github.com/gin-contrib/zap` | 访问日志 |
| API 文档 | `github.com/swaggo/swag` + `github.com/swaggo/gin-swagger` | 注解生成 OpenAPI + UI |
| 测试 | `github.com/stretchr/testify` | 断言 + mock |

**依赖方向（单向）**：

```
handler  →  service  →  repository  →  model
   ↘                        ↙
        dto / pkg / apperror
```

- 上层依赖下层，下层**不反向依赖**上层。
- `service` **不 import gin**（保持框架无关，便于复用与测试）。
- `handler` **不直接访问数据库**。
- `pkg/*` 为可复用通用包，不依赖 `internal/*`。

---

## 3. 目录结构

```
go-web-template
├── cmd
│   └── server
│       └── main.go                 # 入口：解析参数 → 加载配置 → 装配 → 优雅启停
├── internal
│   ├── app
│   │   └── app.go                  # 依赖装配(DI)、App 结构体、Run/Shutdown
│   ├── config
│   │   ├── config.go               # 配置结构体 + Viper 加载
│   │   └── validate.go             # 必填项校验
│   ├── router
│   │   ├── router.go               # 全局中间件 + 路由树装配
│   │   ├── api.go                  # /api/v1 路由分组
│   │   ├── swagger.go              # Swagger UI 挂载（非 prod）
│   │   └── health.go               # /healthz /readyz
│   ├── middleware
│   │   ├── requestid.go            # 请求 ID
│   │   ├── accesslog.go            # 访问日志(zap)
│   │   ├── recovery.go             # panic 恢复
│   │   ├── cors.go                 # 跨域
│   │   ├── ratelimit.go            # 限流（Redis）
│   │   ├── auth.go                 # 鉴权占位（默认放行，TODO 待实现）
│   │   └── error.go                # 统一错误渲染
│   ├── handler                     # HTTP 层 (Controller)
│   │   └── post.go
│   ├── service                     # 业务逻辑层
│   │   └── post.go
│   ├── repository                  # 数据访问层（接口 + GORM 实现）
│   │   └── post.go
│   ├── model                       # GORM 实体
│   │   ├── base.go
│   │   └── post.go
│   └── dto                         # 请求/响应结构体 + 校验标签
│       ├── post.go
│       └── common.go               # 分页等
├── pkg                             # 可复用通用包
│   ├── logger/logger.go            # Zap 封装 + GORM logger 适配
│   ├── database/postgres.go        # GORM 初始化、连接池
│   ├── redis/redis.go              # Redis 客户端初始化
│   ├── apperror/error.go           # 业务错误类型 + 错误码
│   ├── response/response.go        # 统一响应封装
│   └── validator/validator.go      # 校验器 + 错误翻译
├── docs                            # Swagger 生成产物（swag init 输出）
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── config
│   ├── config.yaml                 # 默认配置（可提交）
│   ├── config.local.yaml           # 本地覆盖（.gitignore）
│   └── config.prod.yaml            # 生产覆盖
├── migrations                      # golang-migrate SQL
│   ├── 000001_create_posts.up.sql
│   └── 000001_create_posts.down.sql
├── deployments
│   ├── Dockerfile
│   └── docker-compose.yml
├── scripts                         # 辅助脚本
├── .env.example
├── .gitignore
├── .golangci.yml
├── Makefile
├── go.mod
└── README.md
```

**为什么用 `internal/`**：Go 编译器强制 `internal/` 只能被本项目引用，防止模板被复制后内部实现被外部误依赖。

---

## 4. 分层职责

| 层 | 职责 | 禁止 |
| --- | --- | --- |
| `handler` | 解析/绑定入参、调用 service、写响应 | 写业务逻辑、直接访问 DB |
| `service` | 业务编排、事务边界、权限判断 | 依赖 gin / HTTP 细节 |
| `repository` | 数据访问，屏蔽 GORM 细节 | 写业务规则 |
| `model` | 表结构映射（GORM 实体） | 承载 DTO 逻辑 |
| `dto` | 对外契约（请求/响应结构 + 校验） | 直接映射数据库表 |
| `pkg` | 通用能力（日志/DB/Redis/错误…） | 依赖 `internal` |

**接口放置**：`repository` 包内定义接口（如 `PostRepository`）+ GORM 实现，`service` 依赖接口，便于单测 mock。

---

## 5. 请求生命周期

```
Client
  → RequestID            注入 X-Request-Id
  → AccessLog            开始计时、绑定 logger 到 ctx
  → Recovery             捕获 panic
  → CORS                 跨域处理
  → RateLimit            限流（Redis，敏感路由）
  → Auth                 鉴权占位（默认放行，预留身份注入 ctx）
  → Handler              绑定/校验 dto → 调用 service
       → Service         业务逻辑 → 调用 repository
            → Repository → GORM → PostgreSQL
  → ErrorHandler         统一渲染错误（若有 c.Errors）
  → Response             统一 JSON 信封
```

---

## 6. 配置设计

**优先级**：命令行 flag > 环境变量 > `config.{env}.yaml` > `config.yaml` > 代码默认值。

```yaml
app:
  name: go-web-template
  env: dev                 # dev | test | prod
  debug: true
server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 15s
  write_timeout: 15s
  shutdown_timeout: 10s
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: app
  sslmode: disable
  max_open_conns: 50
  max_idle_conns: 10
  conn_max_lifetime: 1h
  log_level: warn           # silent | error | warn | info
redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 10
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
ratelimit:
  enabled: true
  key_prefix: rl            # Redis key 前缀
  rate: 20                  # 每 window 允许请求数
  burst: 40                 # 突发上限
  window: 1s                # 统计窗口
  fail_open: true           # Redis 不可用时是否放行
swagger:
  enabled: true             # 生产建议置 false
log:
  level: info               # debug | info | warn | error
  format: json              # json | console
  output: stdout            # stdout | 文件路径
cors:
  allow_origins: ["*"]
```

- 环境变量映射：`APP_SERVER_PORT` → `server.port`（Viper `SetEnvPrefix("APP")` + `SetEnvKeyReplacer`）。
- 敏感项（DB / Redis 密码）通过环境变量注入，**不提交仓库**。
- 用 `-c ./config/config.local.yaml` 或 `CONFIG_FILE` 指定配置路径。

---

## 7. 统一响应与错误处理

**成功响应**（HTTP 200）：

```json
{ "code": 0, "message": "ok", "data": { "...": "..." }, "request_id": "..." }
```

**失败响应**（HTTP 状态码语义化，body 带业务码）：

```json
{
  "code": 12001,
  "message": "帖子不存在",
  "request_id": "..."
}
```

**错误码分段**：

| 段 | 含义 | 示例 |
| --- | --- | --- |
| `0` | 成功 | |
| `10xxx` | 通用/参数 | 10001 参数校验失败 |
| `11xxx` | 认证鉴权 | 11001 未认证、11002 无权限 |
| `12xxx` | 业务 | 12001 资源不存在、12002 资源冲突 |
| `13xxx` | 限流 | 13001 请求过于频繁 |
| `50xxx` | 系统 | 50000 内部错误 |

设计要点：

- `pkg/apperror` 定义 `AppError{ Code int, HTTPStatus int, Message string, Cause error }` 与构造函数（`BadRequest/Unauthorized/NotFound/Conflict/TooManyRequests/Internal`）。
- Handler 出错时调用 `c.Error(err)` 并 `return`；`middleware/error.go` 在链路末端读取 `c.Errors`，将 `*AppError` 映射为对应 HTTP 状态 + 业务码，未知错误统一 500 且**不泄露内部细节**。
- `message` 面向用户/前端，`Cause` 仅进日志。

---

## 8. 鉴权设计（预留）

v1 **不实现具体鉴权逻辑**，只保留中间件扩展点，默认放行；**无任何配置项**。

**`middleware/auth.go`**：

- `Auth()` 中间件挂载在 `/api/v1` 受保护分组上，当前实现为直接 `c.Next()` 放行。
- 真实鉴权逻辑以 `TODO` 标注在函数体内，接入时直接实现，无需新增配置。

**预留约定**（后续接入时保持接口不变）：

- 身份信息统一写入 `gin.Context`（`ctx.Set(middleware.ContextKeyUserID, ...)`），业务层通过 `middleware.UserID(c)` 读取，不感知鉴权实现细节。
- 未通过校验统一返回 401（业务码 `11001`）。
- 后续可平滑接入 JWT / Session / API Key，均只需改 `middleware/auth.go`，handler/service 无需改动。
- 限流维度可据此升级：鉴权关闭时按客户端 IP，鉴权开启后可按用户身份。

**v1 不包含**：注册/登录接口、token 签发与刷新、密码哈希与校验。

---

## 9. 数据库与迁移

**实体（示例）**：

- `posts`：`id uuid pk`、`title varchar`、`content text`、`status smallint`、`author varchar`、`created_at`、`updated_at`、`deleted_at`（软删）。

**约定**：

- 主键 UUID，时间统一 **UTC**，软删除用 `gorm.DeletedAt`。
- 迁移用 **golang-migrate**（`migrations/*.up.sql` / `*.down.sql`），生产禁止 `AutoMigrate`；仅本地 dev 可通过开关启用 `AutoMigrate` 便利开发。
- GORM 配置连接池（`MaxOpenConns/MaxIdleConns/ConnMaxLifetime`），日志接入 Zap。

**示例接口（posts）**：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/v1/posts` | 创建 |
| GET | `/api/v1/posts` | 分页列表 |
| GET | `/api/v1/posts/:id` | 详情 |
| PUT | `/api/v1/posts/:id` | 更新 |
| DELETE | `/api/v1/posts/:id` | 删除（软删） |

---

## 10. API 文档（Swagger）

- 使用 **swaggo** 注解生成 OpenAPI，产物输出到 `docs/`（`swag init`）。
- UI 通过 `gin-swagger` 挂载在 `/swagger/index.html`，仅在 `swagger.enabled=true` 时注册（生产建议关闭）。
- Handler 上以注释声明接口：

```go
// CreatePost 创建帖子
// @Summary      创建帖子
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreatePostRequest  true  "帖子信息"
// @Success      200   {object}  response.Body{data=dto.PostResponse}
// @Failure      400   {object}  response.Body
// @Router       /api/v1/posts [post]
```

- 全局信息（title/version/host/basePath）在 `cmd/server/main.go` 顶部注释声明。
- `make swagger` 重新生成；CI 中校验 `docs/` 与代码一致（`swag init` 后 git diff 为空）。

---

## 11. 日志与可观测性

- Zap 双模式：`json`（生产）/ `console`（开发）。
- 每条日志携带 `request_id`，与响应头 `X-Request-Id` 一致，便于全链路追踪。
- 访问日志字段：`method`、`path`、`status`、`latency`、`client_ip`、`request_id`。
- 健康检查：`/healthz`（存活）、`/readyz`（含 PostgreSQL `Ping` + Redis `Ping`）。
- 可选：Prometheus `/metrics`、Gin pprof（仅 dev）。

---

## 12. 中间件清单

| 中间件 | 作用 | 作用域 |
| --- | --- | --- |
| RequestID | 生成/透传请求 ID | 全局 |
| AccessLog | 结构化访问日志 | 全局 |
| Recovery | panic → 500 + 堆栈日志 | 全局 |
| CORS | 跨域白名单 | 全局 |
| ErrorHandler | 统一错误渲染 | 全局 |
| RateLimit | Redis 分布式限流（超限 429 / 13001） | 敏感路由 |
| Auth | 鉴权占位（无配置，默认放行） | `/api/v1` 受保护分组 |

---

## 13. 编码约定

- 包名小写单词；文件名 `snake_case`；导出标识符必须有注释。
- Handler 方法挂在结构体上（`PostHandler`），构造函数注入依赖：`handler.NewPostHandler(svc)`。
- 统一用 **构造函数注入** 手工装配（`internal/app`），暂不引入 wire；规模变大可平滑替换。
- 分页约定：入参 `page`(默认 1)、`page_size`(默认 20，上限 100)；响应 `{ list, total, page, page_size }`。
- 时间返回 RFC3339 UTC；`Content-Type: application/json; charset=utf-8`。
- 所有错误必须显式处理，禁止裸 `panic`；`context.Context` 从 handler 透传到 repository。
- 提交前运行 `make fmt`（gofmt + goimports，本地包单独成组）与 `make lint`；CI 用 `make fmt-check` 卡格式。
- 工具版本由 `go.mod` 的 `tool` 指令固定（`go tool goimports`），无需手动 `go install`。

---

## 14. Docker 与部署

- `deployments/Dockerfile`：多阶段构建（`golang:1.26-alpine` 编译 → `gcr.io/distroless/static` 或 `alpine` 运行），产出最小镜像。
- `deployments/docker-compose.yml`：`app` + `postgres` + `redis` 三服务，含健康检查、数据卷、`.env` 注入；`app` 依赖 pgsql / redis 健康后启动。
- 镜像内包含 `migrations/`，通过 `make migrate-up` 或 entrypoint 执行迁移。

---

## 15. 常用命令（Makefile）

```
make run            # 本地运行
make build          # 编译二进制
make test           # 单元测试
make fmt            # 格式化代码 + 整理 import 分组（gofmt + goimports）
make fmt-check      # 校验格式（CI 用，不修改文件）
make lint           # golangci-lint（formatter 含 gofmt/goimports）
make swagger        # 生成 Swagger 文档（swag init）
make migrate-up     # 执行迁移
make migrate-down   # 回滚一步
make migrate-create name=xxx   # 新建迁移
make docker-up      # docker-compose up -d
make docker-down
```

---

## 16. 测试策略

- **service 层**：用 `testify/mock` mock repository，覆盖业务分支（重点）。
- **handler 层**：`httptest` + mock service，验证状态码/响应体/绑定校验。
- **repository 层**：可选 `testcontainers-go` 起真实 PostgreSQL 做集成测试。
- **限流中间件**：用 `miniredis` 做无外部依赖的 Redis 单测。
- **pkg 层**：apperror、response、validator 纯函数单测。
- CI 门禁：`go vet` + `golangci-lint` + `go test -race ./...`。

---

## 17. 后续可选增强（不在 v1 范围）

- 鉴权落地：JWT / Session / API Key（填实 `middleware/auth.go`）。
- 权限模型：RBAC（角色/权限表 + Casbin）。
- 缓存：业务数据缓存（复用 `pkg/redis`）。
- 可观测：OpenTelemetry 链路追踪 + Prometheus 指标。
- 消息队列、定时任务、文件上传（对象存储）。
- 代码生成脚手架（entity → repo/service/handler）。

---

## 18. 待确认事项

1. **Swagger 生产策略**：确认生产 `swagger.enabled=false`（仅非 prod 暴露 UI）？
2. **限流维度**：鉴权默认关闭，v1 限流按**客户端 IP**；后续接入鉴权再按用户维度，是否认可？
3. **Redis 不可用策略**：`ratelimit.fail_open=true`（放行，保可用性）是否符合预期？
4. **`posts.status` 语义**：如 `0=草稿 / 1=已发布`，是否需要？

---

> 评审确认后，下一步按此设计搭建骨架代码（目录 + go.mod + 配置 + main.go + 基础中间件 + Redis 限流 + Swagger + posts 示例 CRUD）。
