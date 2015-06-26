/**
 * Copyright 2015 @ S1N1 Team.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dps
import (
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/query"
	"go2o/src/core/domain/interface/partner"
)


type contentService struct {
	_contentRep content.IContentRep
	_query     *query.ContentQuery
}

func NewContentService(rep content.IContentRep, q *query.ContentQuery) *contentService {
	return &contentService{
		_contentRep: rep,
		_query:     q,
	}
}

// 获取页面
func (this *contentService) GetPage(partnerId,id int)*content.ValuePage{
	c := this._contentRep.GetContent(partnerId)
	page := c.GetPage(id)
	if page != nil{
		return page.GetValue()
	}
	return nil
}

// 保存页面
func (this *contentService)  SavePage(partnerId int, v *content.ValuePage)(int,error){
	c := this._contentRep.GetContent(partnerId)
	var page content.IPage

	if v.PartnerId != partnerId{
		return -1, partner.ErrNotMatch
	}

	if v.Id > 0 {
		page = c.GetPage(v.Id)
		page.SetValue(v)
	}else{
		page = c.CreatePage(v)
	}

	return page.Save()
}