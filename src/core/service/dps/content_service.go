/**
 * Copyright 2015 @ z3q.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dps

import (
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/query"
)

type contentService struct {
	_contentRep content.IContentRep
	_query      *query.ContentQuery
}

func NewContentService(rep content.IContentRep, q *query.ContentQuery) *contentService {
	return &contentService{
		_contentRep: rep,
		_query:      q,
	}
}

// 获取页面
func (this *contentService) GetPage(merchantId, id int) *content.ValuePage {
	c := this._contentRep.GetContent(merchantId)
	page := c.GetPage(id)
	if page != nil {
		return page.GetValue()
	}
	return nil
}

// 根据标识获取页面
func (this *contentService) GetPageByIndent(merchantId int, indent string) *content.ValuePage {
	c := this._contentRep.GetContent(merchantId)
	page := c.GetPageByStringIndent(indent)
	if page != nil {
		return page.GetValue()
	}
	return nil
}

// 保存页面
func (this *contentService) SavePage(merchantId int, v *content.ValuePage) (int, error) {
	c := this._contentRep.GetContent(merchantId)
	var page content.IPage

	if v.MerchantId != merchantId {
		return -1, merchant.ErrPartnerNotMatch
	}

	if v.Id > 0 {
		page = c.GetPage(v.Id)
		page.SetValue(v)
	} else {
		page = c.CreatePage(v)
	}

	return page.Save()
}

// 删除页面
func (this *contentService) DeletePage(merchantId int, pageId int) error {
	c := this._contentRep.GetContent(merchantId)
	return c.DeletePage(pageId)
}
