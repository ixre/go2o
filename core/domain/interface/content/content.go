/**
 * Copyright 2015 @ 56x.net.
 * name : 9.content.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

import (
	"github.com/ixre/go2o/core/infrastructure/domain"
	d "github.com/ixre/go2o/core/infrastructure/domain"
)

var (
	ErrCategoryContainArchive = domain.NewError(
		"err_category_contain_archive", "栏目包含文章,不允许删除")

	ErrCategoryAliasExists = domain.NewError(
		"err_category_alias_exists", "已存在相同标识的栏目")

	ErrAliasHasExists = domain.NewError(
		"err_content_alias_exists", "页面别名已存在")

	NotSetCategory = domain.NewError(
		"err_not_set_category", "请选择分类")

	ErrUserNotMatch = domain.NewError(
		"err_content_user_not_match", "用户不匹配")

	ErrNoSuchPage = domain.NewError(
		"err_no_such_page", "页面不存在")

	ErrInternalPage = domain.NewError(
		"err_internal_page", "不允许操作内置页面")
)

const (
	// FlagInternal 是否内置
	FCategoryInternal = 1 << iota
	// FCategorySupportPost 是否支持投稿
	FCategorySupportPost = 2
	// FCategoryOpen 是否对会员开放
	FCategoryOpen = 4
)

const ()

type (
	IContentAggregateRoot interface {
		d.IAggregateRoot
		// ArticleManager 文章服务
		ArticleManager() IArticleManager
		// PageManager 页面服务
		PageManager() IPageManager
	}
)
