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
	_contentRep content.IContentRep
	_merchantId int
}

func NewContent(merchantId int, rep content.IContentRep) content.IContent {
	return &Content{
		_contentRep: rep,
		_merchantId: merchantId,
	}
}

// 获取聚合根编号
func (this *Content) GetAggregateRootId() int {
	return this._merchantId
}

// 创建页面
func (this *Content) CreatePage(v *content.ValuePage) content.IPage {
	return NewPage(this.GetAggregateRootId(), this._contentRep, v)
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

// 创建文章
func (this *Content) CreateArticle(*content.ValuePage) content.IArticle {
	//todo:
	return nil
}

// 获取文章
func (this *Content) GetArticle(id int) content.IArticle {
	//todo:
	return nil
}

// 获取文章列表
func (this *Content) GetArticleList(categoryId int, start, over int) []content.IArticle {
	//todo:
	return nil
}

// 删除文章
func (this *Content) DeleteArticle(id int) error {
	//todo:
	return nil
}
