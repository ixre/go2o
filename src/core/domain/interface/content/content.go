/**
 * Copyright 2015 @ S1N1 Team.
 * name : content.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

type IContent interface {
	// 获取聚合根编号
	GetAggregateRootId() int

	// 创建页面
	CreatePage(*ValuePage) IPage

	// 获取页面
	GetPage(id int) IPage

	// 根据字符串标识获取页面
	GetPageByStringIndent(indent string) IPage

	// 删除页面
	DeletePage(id int) error

	// 创建文章
	CreateArticle(*ValuePage) IArticle

	// 获取文章
	GetArticle(id int) IArticle

	// 获取文章列表
	GetArticleList(categoryId int, start, over int) []IArticle

	// 删除文章
	DeleteArticle(id int) error
}
