// Package main 是服务入口。
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/go-web-template/internal/app"
	"github.com/example/go-web-template/internal/config"
	"github.com/example/go-web-template/pkg/validator"

	_ "github.com/example/go-web-template/docs" // Swagger 生成产物
)

// @title			go-web-template API
// @version		1.0
// @description	Go Web 模板：Gin + GORM + PostgreSQL
// @BasePath		/
// @schemes		http
func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "", "配置文件路径，默认读取 ./config/config.yaml")
	flag.Parse()

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	validator.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	if err := application.Run(ctx); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
