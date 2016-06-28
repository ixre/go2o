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
	"go2o/core/domain/interface/content"
	"go2o/core/domain/interface/merchant"
	"go2o/core/query"
)

type contentService struct {
	_contentRep content.IContentRep
	_query      *query.ContentQuery
	_sysContent content.IContent
}

func NewContentService(rep content.IContentRep, q *query.ContentQuery) *contentService {
	return &contentService{
		_contentRep: rep,
		_query:      q,
		_sysContent: rep.GetContent(0),
	}
}

// 获取页面
func (this *contentService) GetPage(merchantId, id int) *content.Page {
	c := this._contentRep.GetContent(merchantId)
	page := c.GetPage(id)
	if page != nil {
		return page.GetValue()
	}
	return nil
}

// 根据标识获取页面
func (this *contentService) GetPageByIndent(merchantId int, indent string) *content.Page {
	c := this._contentRep.GetContent(merchantId)
	page := c.GetPageByStringIndent(indent)
	if page != nil {
		return page.GetValue()
	}
	return nil
}

// 保存页面
func (this *contentService) SavePage(merchantId int, v *content.Page) (int, error) {
	c := this._contentRep.GetContent(merchantId)
	var page content.IPage

	if v.UserId != merchantId {
		return -1, merchant.ErrMerchantNotMatch
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

// 获取所有栏目
func (this *contentService) GetArticleCategories() []content.ArticleCategory {
	list := this._sysContent.ArticleManager().GetAllCategory()
	arr := make([]content.ArticleCategory, len(list))
	for i, v := range list {
		arr[i] = v.GetValue()
	}
	return arr
}

// 获取文章栏目
func (this *contentService) GetArticleCategory(id int) content.ArticleCategory {
	m := this._sysContent.ArticleManager().GetCategory(id)
	if m != nil {
		return m.GetValue()
	}
	return content.ArticleCategory{}
}

// 根据标识获取文章栏目
func (this *contentService) GetArticleCategoryByAlias(alias string) content.ArticleCategory {
	m := this._sysContent.ArticleManager().GetCategoryByAlias(alias)
	if m != nil {
		return m.GetValue()
	}
	return content.ArticleCategory{}
}

// 保存文章栏目
func (this *contentService) SaveArticleCategory(v *content.ArticleCategory) (int, error) {
	m := this._sysContent.ArticleManager()
	c := m.GetCategory(v.Id)
	if c == nil {
		c = m.CreateCategory(v)
	}
	err := c.SetValue(v)
	if err == nil {
		return c.Save()
	}
	return -1, err
}

// 删除文章分类
func (this *contentService) DeleteArticleCategory(categoryId int) error {
	return this._sysContent.ArticleManager().DelCategory(categoryId)
}

// 获取文章
func (this *contentService) GetArticle(id int) content.Article {
	a := this._sysContent.ArticleManager().GetArticle(id)
	if a != nil {
		return a.GetValue()
	}
	return content.Article{}
}

// 删除文章
func (this *contentService) DeleteArticle(id int) error {
	return this._sysContent.ArticleManager().DeleteArticle(id)
}

// 保存文章
func (this *contentService) SaveArticle(e *content.Article) (int, error) {
	m := this._sysContent.ArticleManager()
	a := m.GetArticle(e.Id)
	if a == nil {
		a = m.CreateArticle(e)
	}
	err := a.SetValue(e)
	if err == nil {
		return a.Save()
	}
	return -1, err
}
