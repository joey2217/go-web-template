// Package handler 实现 HTTP 层（Controller）。
package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/example/go-web-template/internal/dto"
	"github.com/example/go-web-template/internal/service"
	"github.com/example/go-web-template/pkg/apperror"
	"github.com/example/go-web-template/pkg/logger"
	"github.com/example/go-web-template/pkg/response"
	"github.com/example/go-web-template/pkg/validator"
)

// UserHandler 处理用户相关 HTTP 请求。
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler 创建用户处理器。
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Create 创建用户
//
//	@Summary	创建用户
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		body	body		dto.CreateUserRequest	true	"用户信息"
//	@Success	201		{object}	response.Body{data=dto.UserResponse}
//	@Failure	400		{object}	response.Body
//	@Failure	409		{object}	response.Body
//	@Failure	500		{object}	response.Body
//	@Router		/api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	log := logger.FromGin(c)
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("bind user json failed", zap.Error(err))
		c.Error(apperror.BadRequest(validator.Translate(err)))
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		log.Error("create user failed", zap.Error(err))
		c.Error(err)
		return
	}
	log.Info("user created", zap.String("user_id", resp.ID), zap.String("username", resp.Username))
	response.Created(c, resp)
}

// List 分页查询用户
//
//	@Summary	分页查询用户
//	@Tags		users
//	@Produce	json
//	@Param		page		query		int		false	"页码，默认 1"
//	@Param		page_size	query		int		false	"每页数量，默认 20，最大 100"
//	@Param		keyword		query		string	false	"用户名/邮箱/昵称关键字"
//	@Param		status		query		int		false	"状态：0 正常，1 禁用"
//	@Success	200			{object}	response.Body{data=response.PageData}
//	@Failure	400			{object}	response.Body
//	@Failure	500			{object}	response.Body
//	@Router		/api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	var query dto.ListUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(apperror.BadRequest(validator.Translate(err)))
		return
	}
	data, err := h.svc.List(c.Request.Context(), &query)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, data)
}

// Get 查询用户详情
//
//	@Summary	查询用户详情
//	@Tags		users
//	@Produce	json
//	@Param		id	path		string	true	"用户 ID"
//	@Success	200	{object}	response.Body{data=dto.UserResponse}
//	@Failure	404	{object}	response.Body
//	@Failure	500	{object}	response.Body
//	@Router		/api/v1/users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	resp, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, resp)
}

// Update 更新用户（局部更新）
//
//	@Summary	更新用户
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"用户 ID"
//	@Param		body	body		dto.UpdateUserRequest	true	"用户信息（仅传需要修改的字段）"
//	@Success	200		{object}	response.Body{data=dto.UserResponse}
//	@Failure	400		{object}	response.Body
//	@Failure	404		{object}	response.Body
//	@Failure	409		{object}	response.Body
//	@Failure	500		{object}	response.Body
//	@Router		/api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest(validator.Translate(err)))
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, resp)
}

// Delete 删除用户
//
//	@Summary	删除用户
//	@Tags		users
//	@Produce	json
//	@Param		id	path	string	true	"用户 ID"
//	@Success	204	"No Content"
//	@Failure	404	{object}	response.Body
//	@Failure	500	{object}	response.Body
//	@Router		/api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.Error(err)
		return
	}
	response.NoContent(c)
}
