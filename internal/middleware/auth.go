package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/example/go-web-template/internal/config"
)

// ContextKeyUserID 是写入 gin.Context 的身份键。
const ContextKeyUserID = "userID"

// Auth 是鉴权中间件占位实现。
//
// 当 auth.enabled=false（默认）时直接放行。
// 当开启后，应在此校验凭证、解析身份并写入上下文，例如：
//
//	c.Set(ContextKeyUserID, userID)
//
// 校验失败时返回 401：
//
//	response.Fail(c, apperror.Unauthorized("未认证"))
//
// 这样后续接入 JWT / Session / API Key 只需修改本文件，handler 与 service 无需改动。
func Auth(cfg config.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Enabled {
			c.Next()
			return
		}
		// TODO: 在此实现真实鉴权逻辑。
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
