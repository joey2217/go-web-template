APP_NAME := go-web-template
MAIN     := ./cmd/server
BIN_DIR  := bin
MODULE   := $(shell go list -m)

DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/app?sslmode=disable

.PHONY: help deps run build test fmt fmt-check lint clean gen swagger \
        migrate-up migrate-down migrate-create migrate-force \
        docker-up docker-down docker-logs

help: ## 显示可用命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

deps: ## 整理依赖
	go mod tidy

run: ## 本地运行
	go run $(MAIN)

build: ## 编译二进制到 bin/
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/$(APP_NAME) $(MAIN)

test: ## 单元测试
	go test -race -cover ./...

fmt: ## 格式化代码并整理 import 分组（gofmt + goimports）
	go tool goimports -w -local $(MODULE) .

fmt-check: ## 校验格式与 import 分组（CI 用，不修改文件）
	@out=$$(go tool goimports -l -local $(MODULE) .); \
	if [ -n "$$out" ]; then \
		echo "以下文件格式或 import 不符合规范，请运行 make fmt："; \
		echo "$$out"; \
		exit 1; \
	fi

lint: ## 静态检查
	golangci-lint run ./...

clean: ## 清理产物
	rm -rf $(BIN_DIR)

swagger: ## 生成 Swagger 文档
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

gen: ## 生成 GORM CLI 数据访问代码（生成后自动整理 import）
	rm -rf internal/repository/gen
	gorm gen -i internal/repository/dao -o internal/repository/gen
	go tool goimports -w -local $(MODULE) internal/repository/gen

migrate-up: ## 执行迁移
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down: ## 回滚一步迁移
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create: ## 新建迁移，用法：make migrate-create name=add_xxx
	migrate create -ext sql -dir migrations -seq $(name)

migrate-force: ## 强制设置版本，用法：make migrate-force version=1
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

docker-up: ## 启动依赖与应用
	docker compose -f deployments/docker-compose.yml up -d --build

docker-down: ## 停止并移除容器
	docker compose -f deployments/docker-compose.yml down

docker-logs: ## 查看应用日志
	docker compose -f deployments/docker-compose.yml logs -f app
