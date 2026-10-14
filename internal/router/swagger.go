package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// registerSwagger 挂载 Swagger UI，访问路径为 /swagger/index.html。
//
// 依赖 cmd/server/main.go 中通过 swag init 生成的 docs 包（由 main 包匿名导入）。
func registerSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
