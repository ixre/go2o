/**
 * Copyright 2015 @ z3q.net.
 * name : content
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

import "go2o/core/domain/interface/content"

var _ content.IContent = new(Content)

type Content struct {
	contentRep     content.IContentRep
	userId         int64
	articleManager content.IArticleManager
}

func NewContent(userId int64, rep content.IContentRep) content.IContent {
	return &Content{
		contentRep: rep,
		userId:     userId,
	}
}

// 获取聚合根编号
func (c *Content) GetAggregateRootId() int64 {
	return c.userId
}

// 文章服务
func (c *Content) ArticleManager() content.IArticleManager {
	if c.articleManager == nil {
		c.articleManager = newArticleManagerImpl(c.userId, c.contentRep)
	}
	return c.articleManager
}

// 创建页面
func (c *Content) CreatePage(v *content.Page) content.IPage {
	return newPage(c.GetAggregateRootId(), c.contentRep, v)
}

// 获取页面
func (c *Content) GetPage(id int64) content.IPage {
	v := c.contentRep.GetPageById(c.GetAggregateRootId(), id)
	if v != nil {
		return c.CreatePage(v)
	}
	return nil
}

// 根据字符串标识获取页面
func (c *Content) GetPageByStringIndent(indent string) content.IPage {
	v := c.contentRep.GetPageByStringIndent(c.GetAggregateRootId(), indent)
	if v != nil {
		return c.CreatePage(v)
	}
	return nil
}

// 删除页面
func (c *Content) DeletePage(id int64) error {
	return c.contentRep.DeletePage(c.GetAggregateRootId(), id)
}
