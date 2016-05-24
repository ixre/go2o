/**
 * Copyright 2015 @ z3q.net.
 * name : article
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

// 文章
type IArticle interface {
	// 获取领域编号
	GetDomainId() int

	// 获取值
	GetValue() *ValueArticle

	// 设置值
	SetValue(*ValueArticle) error

	// 保存
	Save() (int, error)
}
