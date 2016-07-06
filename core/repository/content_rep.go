/**
 * Copyright 2015 @ z3q.net.
 * name : content_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	contentImpl "go2o/core/domain/content"
	"go2o/core/domain/interface/content"
)

var _ content.IContentRep = new(contentRep)

type contentRep struct {
	db.Connector
}

// 内容仓储
func NewContentRep(c db.Connector) content.IContentRep {
	return &contentRep{
		Connector: c,
	}
}

// 获取内容
func (this *contentRep) GetContent(merchantId int) content.IContent {
	return contentImpl.NewContent(merchantId, this)
}

// 根据编号获取页面
func (this *contentRep) GetPageById(merchantId, id int) *content.Page {
	var e content.Page
	if err := this.Connector.GetOrm().Get(id, &e); err == nil && e.UserId == merchantId {
		return &e
	}
	return nil
}

// 根据标识获取页面
func (this *contentRep) GetPageByStringIndent(merchantId int, indent string) *content.Page {
	var e content.Page
	if err := this.Connector.GetOrm().GetBy(&e, "mch_id=? and str_indent=?", merchantId, indent); err == nil {
		return &e
	}
	return nil
}

// 删除页面
func (this *contentRep) DeletePage(merchantId, id int) error {
	_, err := this.Connector.GetOrm().Delete(content.Page{}, "mch_id=? AND id=?", merchantId, id)
	return err
}

// 保存页面
func (this *contentRep) SavePage(merchantId int, v *content.Page) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		var id2 int64
		_, id2, err = orm.Save(nil, v)
		v.Id = int(id2)
	}
	return v.Id, err
}

// 获取文章数量
func (this *contentRep) GetArticleNumByCategory(categoryId int) int {
	num := 0
	this.Connector.ExecScalar("SELECT COUNT(0) FROM con_article WHERE cat_id=?",
		&num, categoryId)
	return num
}

// 获取栏目
func (this *contentRep) GetAllArticleCategory() []*content.ArticleCategory {
	list := []*content.ArticleCategory{}
	this.Connector.GetOrm().Select(&list, "")
	return list
}

// 判断栏目是否存在
func (this *contentRep) CategoryExists(indent string, id int) bool {
	num := 0
	this.Connector.ExecScalar("SELECT COUNT(0) FROM con_category WHERE indent=? and id<>id",
		&num, indent, id)
	return num > 0
}

// 保存栏目
func (this *contentRep) SaveCategory(v *content.ArticleCategory) (id int, err error) {
	orm := this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		var id2 int64
		_, id2, err = orm.Save(nil, v)
		v.Id = int(id2)
	}
	return v.Id, err
}

// 删除栏目
func (this *contentRep) DeleteCategory(id int) error {
	return this.Connector.GetOrm().DeleteByPk(&content.ArticleCategory{}, id)
}

// 获取文章
func (this *contentRep) GetArticleById(id int) *content.Article {
	e := content.Article{}
	if this.Connector.GetOrm().Get(id, &e) == nil {
		return &e
	}
	return nil
}

// 获取文章列表
func (this *contentRep) GetArticleList(categoryId int, begin int, end int) []*content.Article {
	list := []*content.Article{}
	this.Connector.GetOrm().SelectByQuery(&content.Article{},
		"cat_id=? LIMIT ?,?", categoryId, begin, end-begin)
	return list
}

// 保存文章
func (this *contentRep) SaveArticle(v *content.Article) (i int, err error) {
	orm := this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		var id2 int64
		_, id2, err = orm.Save(nil, v)
		v.Id = int(id2)
	}
	return v.Id, err
}

// 删除文章
func (this *contentRep) DeleteArticle(id int) error {
	return this.Connector.GetOrm().DeleteByPk(&content.Article{}, id)
}
