// Package dto 定义对外请求/响应结构体与校验规则。
package dto

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// PageQuery 是通用分页入参。
type PageQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// Normalize 填充分页默认值并做上限保护。
func (p *PageQuery) Normalize() {
	if p.Page <= 0 {
		p.Page = defaultPage
	}
	if p.PageSize <= 0 {
		p.PageSize = defaultPageSize
	}
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
}

// Offset 返回 SQL OFFSET。
func (p *PageQuery) Offset() int { return (p.Page - 1) * p.PageSize }
