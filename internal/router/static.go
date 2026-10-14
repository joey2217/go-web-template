package router

import "github.com/gin-gonic/gin"

// registerStatic 挂载静态资源目录与 favicon。
//
// 目录约定：
//
//	./public/             静态资源根目录（CSS / JS / 图片等）
//	./public/favicon.ico  浏览器自动请求的根路径图标
//
// 访问路径：
//
//	GET /static/<file>    任意静态文件
//	GET /favicon.ico      根路径 favicon（浏览器 / 爬虫默认请求）
func registerStatic(r *gin.Engine) {
	r.Static("/static", "./public")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
}
