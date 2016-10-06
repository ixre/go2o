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
	_contentRep     content.IContentRep
	_userId         int
	_articleManager content.IArticleManager
}

func NewContent(userId int, rep content.IContentRep) content.IContent {
	return &Content{
		_contentRep: rep,
		_userId:     userId,
	}
}

// 获取聚合根编号
func (c *Content) GetAggregateRootId() int {
	return c._userId
}

// 文章服务
func (c *Content) ArticleManager() content.IArticleManager {
	if c._articleManager == nil {
		c._articleManager = newArticleManagerImpl(c._userId, c._contentRep)
	}
	return c._articleManager
}

// 创建页面
func (c *Content) CreatePage(v *content.Page) content.IPage {
	return newPage(c.GetAggregateRootId(), c._contentRep, v)
}

// 获取页面
func (c *Content) GetPage(id int) content.IPage {
	v := c._contentRep.GetPageById(c.GetAggregateRootId(), id)
	if v != nil {
		return c.CreatePage(v)
	}
	return nil
}

// 根据字符串标识获取页面
func (c *Content) GetPageByStringIndent(indent string) content.IPage {
	v := c._contentRep.GetPageByStringIndent(c.GetAggregateRootId(), indent)
	if v != nil {
		return c.CreatePage(v)
	}
	return nil
}

// 删除页面
func (c *Content) DeletePage(id int) error {
	return c._contentRep.DeletePage(c.GetAggregateRootId(), id)
}
