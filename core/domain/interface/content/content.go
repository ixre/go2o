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
	d	"github.com/ixre/go2o/core/domain"
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
		// CreatePage 创建页面
		CreatePage(*Page) IPage
		// GetPage 获取页面
		GetPage(id int32) IPage
		// GetPageByCode 根据字符串标识获取页面
		GetPageByCode(indent string) IPage
		// DeletePage 删除页面
		DeletePage(id int32) error
	}

	IArchiveRepo interface {
		// GetContent 获取内容
		GetContent(userId int64) IContentAggregateRoot
		// GetPageById 根据编号获取页面
		GetPageById(userId, id int32) *Page
		// GetPageByCode 根据标识获取页面
		GetPageByCode(userId int, code string) *Page
		// DeletePage 删除页面
		DeletePage(userId, id int32) error
		// SavePage 保存页面
		SavePage(userId int32, v *Page) (int32, error)
		// GetAllArticleCategory 获取所有栏目
		GetAllArticleCategory() []*Category
		// GetArticleNumByCategory 获取文章数量
		GetArticleNumByCategory(categoryId int32) int
		// SaveCategory 保存栏目
		SaveCategory(v *Category) (int32, error)
		// CategoryExists 判断栏目是否存在
		CategoryExists(alias string, id int32) bool
		// DeleteCategory 删除栏目
		DeleteCategory(id int32) error
		// 获取文章
		GetArticleById(id int32) *Article
		// 获取文章列表
		GetArticleList(categoryId int32, begin int, end int) []*Article
		// 保存文章
		SaveArticle(v *Article) (int32, error)
		// 删除文章
		DeleteArticle(id int32) error
	}
)
