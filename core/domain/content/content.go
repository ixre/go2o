/**
 * Copyright 2015 @ 56x.net.
 * name : content
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

import "github.com/ixre/go2o/core/domain/interface/content"

var _ content.IContentAggregateRoot = new(Content)

type Content struct {
	artRepo        content.IArticleRepo
	catRepo        content.IArticleCategoryRepo
	userId         int64
	articleManager content.IArticleManager
	pageManager    content.IPageManager
	pageRepo       content.IPageRepo
}

func NewContent(userId int64, rep content.IArticleRepo,
	catRepo content.IArticleCategoryRepo,
	pageRepo content.IPageRepo) content.IContentAggregateRoot {
	return &Content{
		artRepo:  rep,
		pageRepo: pageRepo,
		userId:   userId,
		catRepo:  catRepo,
	}
}

// GetAggregateRootId 获取聚合根编号
func (c *Content) GetAggregateRootId() int {
	return int(c.userId)
}

// ArticleManager 文章服务
func (c *Content) ArticleManager() content.IArticleManager {
	if c.articleManager == nil {
		c.articleManager = newArticleManagerImpl(c.userId, c.catRepo, c.artRepo)
	}
	return c.articleManager
}

// PageManager implements content.IContentAggregateRoot.
func (c *Content) PageManager() content.IPageManager {
	if c.pageManager == nil {
		c.pageManager = &pageManagerImpl{
			zoneId:   int(c.userId),
			pageRepo: c.pageRepo,
		}
	}
	return c.pageManager
}
