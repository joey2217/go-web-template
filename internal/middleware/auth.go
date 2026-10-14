package middleware

import (
	"github.com/gin-gonic/gin"
)

// ContextKeyUserID 是写入 gin.Context 的身份键。
const ContextKeyUserID = "userID"

// Auth 是鉴权中间件占位实现，当前默认放行所有请求。
//
// TODO: 在此实现真实鉴权逻辑：
//   - 校验凭证（如 Authorization 头：JWT / Session / API Key）
//   - 解析身份并写入上下文：c.Set(ContextKeyUserID, userID)
//   - 校验失败时中断请求并返回 401：
//     c.AbortWithStatusJSON(http.StatusUnauthorized, ...)
//     或 c.Error(apperror.Unauthorized("未认证"))
//
// 接入后仅需修改本文件，handler 与 service 通过 UserID(c) 读取身份，无需改动。
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现鉴权校验；当前直接放行。
		c.Next()
	}
}

// UserID 从上下文中读取当前用户身份。
func UserID(c *gin.Context) (string, bool) {
	v, ok := c.Get(ContextKeyUserID)
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}
