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
func (this *Content) GetAggregateRootId() int {
	return this._userId
}

// 文章服务
func (this *Content) ArticleManager() content.IArticleManager {
	if this._articleManager == nil {
		this._articleManager = newArticleManagerImpl(this._userId, this._contentRep)
	}
	return this._articleManager
}

// 创建页面
func (this *Content) CreatePage(v *content.Page) content.IPage {
	return newPage(this.GetAggregateRootId(), this._contentRep, v)
}

// 获取页面
func (this *Content) GetPage(id int) content.IPage {
	v := this._contentRep.GetPageById(this.GetAggregateRootId(), id)
	if v != nil {
		return this.CreatePage(v)
	}
	return nil
}

// 根据字符串标识获取页面
func (this *Content) GetPageByStringIndent(indent string) content.IPage {
	v := this._contentRep.GetPageByStringIndent(this.GetAggregateRootId(), indent)
	if v != nil {
		return this.CreatePage(v)
	}
	return nil
}

// 删除页面
func (this *Content) DeletePage(id int) error {
	return this._contentRep.DeletePage(this.GetAggregateRootId(), id)
}
