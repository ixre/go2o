/**
 * Copyright 2015 @ z3q.net.
 * name : content.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

import "go2o/core/infrastructure/domain"

var (
	ErrCategoryContainArchive *domain.DomainError = domain.NewDomainError(
		"err_category_contain_archive", "栏目包含文章,不允许删除")

	ErrCategoryAliasExists *domain.DomainError = domain.NewDomainError(
		"err_category_alias_exists", "已存在相同标识的栏目")

	NotSetCategory *domain.DomainError = domain.NewDomainError(
		"err_not_set_category", "请选择分类")

	ErrUserNotMatch *domain.DomainError = domain.NewDomainError(
		"err_content_user_not_match", "用户不匹配")
)

type (
	IContent interface {
		// 获取聚合根编号
		GetAggregateRootId() int

		// 文章服务
		ArticleManager() IArticleManager

		// 创建页面
		CreatePage(*Page) IPage

		// 获取页面
		GetPage(id int) IPage

		// 根据字符串标识获取页面
		GetPageByStringIndent(indent string) IPage

		// 删除页面
		DeletePage(id int) error
	}

	IPage interface {
		// 获取领域编号
		GetDomainId() int

		// 获取值
		GetValue() *Page

		// 设置值
		SetValue(*Page) error

		// 保存
		Save() (int, error)
	}

	IContentRep interface {
		// 获取内容
		GetContent(userId int) IContent

		// 根据编号获取页面
		GetPageById(userId, id int) *Page

		// 根据标识获取页面
		GetPageByStringIndent(userId int, indent string) *Page

		// 删除页面
		DeletePage(userId, id int) error

		// 保存页面
		SavePage(userId int, v *Page) (int, error)

		// 获取所有栏目
		GetAllArticleCategory() []*ArticleCategory

		// 获取文章数量
		GetArticleNumByCategory(categoryId int) int

		// 保存栏目
		SaveCategory(v *ArticleCategory) (int, error)

		// 判断栏目是否存在
		CategoryExists(indent string, id int) bool

		// 删除栏目
		DeleteCategory(id int) error

		// 获取文章
		GetArticleById(id int) *Article

		// 获取文章列表
		GetArticleList(categoryId int, begin int, end int) []*Article

		// 保存文章
		SaveArticle(v *Article) (int, error)

		// 删除文章
		DeleteArticle(id int) error
	}

	Page struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`

		// 商户编号
		UserId int `db:"mch_id"`

		// 标题
		Title string `db:"title"`

		// 字符标识
		StrIndent string `db:"str_indent"`

		// 关键词
		KeyWord string `db:"keyword"`

		// 描述
		Description string `db:"description"`

		// 样式表地址
		CssPath string `db:"css_path"`

		// 内容
		Body string `db:"body"`

		// 修改时间
		UpdateTime int64 `db:"update_time"`

		// 是否启用
		Enabled int `db:"enabled"`
	}
)
