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
	ErrInvalidCategory = domain.NewError(
		"err_invalid_category", "文章分类不正确")
	ErrDisallowPostArticle = domain.NewError(
		"err_disallow_post_article", "分类不允许投稿")
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
	// FCategoryPost 是否支持投稿
	FCategoryPost = 2
	// FCategoryOpen 是否对会员显示
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
