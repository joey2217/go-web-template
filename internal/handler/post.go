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

// PostHandler 处理帖子相关 HTTP 请求。
type PostHandler struct {
	svc service.PostService
}

// NewPostHandler 创建帖子处理器。
func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// Create 创建帖子
//
//	@Summary	创建帖子
//	@Tags		posts
//	@Accept		json
//	@Produce	json
//	@Param		body	body		dto.CreatePostRequest	true	"帖子信息"
//	@Success	201		{object}	dto.PostResponse
//	@Failure	400		{object}	response.ErrorBody
//	@Failure	500		{object}	response.ErrorBody
//	@Router		/api/v1/posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	log := logger.FromGin(c)
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("bind post json failed", zap.Error(err))
		c.Error(apperror.BadRequest(validator.Translate(err)))
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		log.Error("create post failed", zap.Error(err))
		c.Error(err)
		return
	}
	log.Info("post created", zap.String("post_id", resp.ID), zap.String("title", resp.Title))
	response.Created(c, resp)
}

// List 分页查询帖子
//
//	@Summary	分页查询帖子
//	@Tags		posts
//	@Produce	json
//	@Param		page		query		int		false	"页码，默认 1"
//	@Param		page_size	query		int		false	"每页数量，默认 20，最大 100"
//	@Param		keyword		query		string	false	"标题/内容关键字"
//	@Param		status		query		int		false	"状态：0 草稿，1 已发布"
//	@Success	200			{object}	response.PageData
//	@Failure	400			{object}	response.ErrorBody
//	@Failure	500			{object}	response.ErrorBody
//	@Router		/api/v1/posts [get]
func (h *PostHandler) List(c *gin.Context) {
	var query dto.ListPostsQuery
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

// Get 查询帖子详情
//
//	@Summary	查询帖子详情
//	@Tags		posts
//	@Produce	json
//	@Param		id	path		string	true	"帖子 ID"
//	@Success	200	{object}	dto.PostResponse
//	@Failure	404	{object}	response.ErrorBody
//	@Failure	500	{object}	response.ErrorBody
//	@Router		/api/v1/posts/{id} [get]
func (h *PostHandler) Get(c *gin.Context) {
	resp, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, resp)
}

// Update 更新帖子
//
//	@Summary	更新帖子
//	@Tags		posts
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"帖子 ID"
//	@Param		body	body		dto.UpdatePostRequest	true	"帖子信息"
//	@Success	200		{object}	dto.PostResponse
//	@Failure	400		{object}	response.ErrorBody
//	@Failure	404		{object}	response.ErrorBody
//	@Failure	500		{object}	response.ErrorBody
//	@Router		/api/v1/posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	var req dto.UpdatePostRequest
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

// Delete 删除帖子
//
//	@Summary	删除帖子
//	@Tags		posts
//	@Produce	json
//	@Param		id	path	string	true	"帖子 ID"
//	@Success	204	"No Content"
//	@Failure	404	{object}	response.ErrorBody
//	@Failure	500	{object}	response.ErrorBody
//	@Router		/api/v1/posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.Error(err)
		return
	}
	response.NoContent(c)
}
