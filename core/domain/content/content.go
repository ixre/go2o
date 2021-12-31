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

var _ content.IContent = new(Content)

type Content struct {
	contentRepo    content.IArchiveRepo
	userId         int64
	articleManager content.IArticleManager
}

func NewContent(userId int64, rep content.IArchiveRepo) content.IContent {
	return &Content{
		contentRepo: rep,
		userId:      userId,
	}
}

// GetAggregateRootId 获取聚合根编号
func (c *Content) GetAggregateRootId() int {
	return int(c.userId)
}

// ArticleManager 文章服务
func (c *Content) ArticleManager() content.IArticleManager {
	if c.articleManager == nil {
		c.articleManager = newArticleManagerImpl(c.userId, c.contentRepo)
	}
	return c.articleManager
}

// CreatePage 创建页面
func (c *Content) CreatePage(v *content.Page) content.IPage {
	return newPage(int32(c.GetAggregateRootId()), c.contentRepo, v)
}

// GetPage 获取页面
func (c *Content) GetPage(id int32) content.IPage {
	v := c.contentRepo.GetPageById(int32(c.GetAggregateRootId()), id)
	if v != nil {
		return c.CreatePage(v)
	}
	return nil
}

// GetPageByCode 根据字符串标识获取页面
func (c *Content) GetPageByCode(indent string) content.IPage {
	v := c.contentRepo.GetPageByCode(c.GetAggregateRootId(), indent)
	if v != nil {
		return c.CreatePage(v)
	}
	return nil
}

// DeletePage 删除页面
func (c *Content) DeletePage(id int32) error {
	ip := c.GetPage(id)
	if ip == nil{
		return content.ErrNoSuchPage
	}
	if ip.GetValue().Flag & content.FlagInternal == content.FlagInternal{
		return content.ErrInternalPage
	}
	return c.contentRepo.DeletePage(int32(c.GetAggregateRootId()), id)
}
