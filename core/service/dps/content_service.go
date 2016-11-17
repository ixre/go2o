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
//todo: 取消merchantId
func (cs *contentService) GetPage(merchantId, id int) *content.Page {
	c := cs._contentRep.GetContent(merchantId)
	page := c.GetPage(id)
	if page != nil {
		return page.GetValue()
	}
	return nil
}

// 根据标识获取页面
//todo: 取消merchantId
func (cs *contentService) GetPageByIndent(userId int, indent string) *content.Page {
	c := cs._contentRep.GetContent(userId)
	page := c.GetPageByStringIndent(indent)
	if page != nil {
		return page.GetValue()
	}
	return nil
}

// 保存页面
func (cs *contentService) SavePage(merchantId int, v *content.Page) (int64, error) {
	c := cs._contentRep.GetContent(merchantId)
	var page content.IPage
	var err error
	if v.UserId != merchantId {
		return -1, merchant.ErrMerchantNotMatch
	}

	if v.Id > 0 {
		page = c.GetPage(v.Id)
	} else {
		page = c.CreatePage(v)
	}
	err = page.SetValue(v)
	if err != nil {
		return 0, err
	}
	return page.Save()
}

// 删除页面
func (cs *contentService) DeletePage(merchantId int, pageId int) error {
	c := cs._contentRep.GetContent(merchantId)
	return c.DeletePage(pageId)
}

// 获取所有栏目
func (cs *contentService) GetArticleCategories() []*content.ArticleCategory {
	list := cs._sysContent.ArticleManager().GetAllCategory()
	arr := make([]*content.ArticleCategory, len(list))
	for i, v := range list {
		val := v.GetValue()
		arr[i] = &val
	}
	return arr
}

// 获取文章栏目
func (cs *contentService) GetArticleCategory(id int) content.ArticleCategory {
	m := cs._sysContent.ArticleManager().GetCategory(id)
	if m != nil {
		return m.GetValue()
	}
	return content.ArticleCategory{}
}

// 根据标识获取文章栏目
func (cs *contentService) GetArticleCategoryByAlias(alias string) content.ArticleCategory {
	m := cs._sysContent.ArticleManager().GetCategoryByAlias(alias)
	if m != nil {
		return m.GetValue()
	}
	return content.ArticleCategory{}
}

// 保存文章栏目
func (cs *contentService) SaveArticleCategory(v *content.ArticleCategory) (int64, error) {
	m := cs._sysContent.ArticleManager()
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
func (cs *contentService) DeleteArticleCategory(categoryId int) error {
	return cs._sysContent.ArticleManager().DelCategory(categoryId)
}

// 获取文章
func (cs *contentService) GetArticle(id int64) *content.Article {
	a := cs._sysContent.ArticleManager().GetArticle(id)
	if a != nil {
		v := a.GetValue()
		return &v
	}
	return nil
}

// 删除文章
func (cs *contentService) DeleteArticle(id int64) error {
	return cs._sysContent.ArticleManager().DeleteArticle(id)
}

// 保存文章
func (cs *contentService) SaveArticle(e *content.Article) (int64, error) {
	m := cs._sysContent.ArticleManager()
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

func (cs *contentService) PagedArticleList(catId, begin, size int,
	where string) (int, []*content.Article) {
	return cs._query.PagedArticleList(catId, begin, size, where)
}
