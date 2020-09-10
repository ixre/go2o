/**
 * Copyright 2015 @ to2.net.
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
	contentRepo    content.IContentRepo
	userId         int64
	articleManager content.IArticleManager
}

func NewContent(userId int64, rep content.IContentRepo) content.IContent {
	return &Content{
		contentRepo: rep,
		userId:      userId,
	}
}

// 获取聚合根编号
func (c *Content) GetAggregateRootId() int {
	return int(c.userId)
}

// 文章服务
func (c *Content) ArticleManager() content.IArticleManager {
	if c.articleManager == nil {
		c.articleManager = newArticleManagerImpl(c.userId, c.contentRepo)
	}
	return c.articleManager
}

// 创建页面
func (c *Content) CreatePage(v *content.Page) content.IPage {
	return newPage(int32(c.GetAggregateRootId()), c.contentRepo, v)
}

// 获取页面
func (c *Content) GetPage(id int32) content.IPage {
	v := c.contentRepo.GetPageById(int32(c.GetAggregateRootId()), id)
	if v != nil {
		return c.CreatePage(v)
	}
	return nil
}

// 根据字符串标识获取页面
func (c *Content) GetPageByStringIndent(indent string) content.IPage {
	v := c.contentRepo.GetPageByStringIndent(int32(c.GetAggregateRootId()), indent)
	if v != nil {
		return c.CreatePage(v)
	}
	return nil
}

// 删除页面
func (c *Content) DeletePage(id int32) error {
	return c.contentRepo.DeletePage(int32(c.GetAggregateRootId()), id)
}
