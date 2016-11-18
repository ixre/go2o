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
	"github.com/jsix/gof/db/orm"
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
func (c *contentRep) GetContent(userId int64) content.IContent {
	return contentImpl.NewContent(userId, c)
}

// 根据编号获取页面
func (c *contentRep) GetPageById(mchId, id int64) *content.Page {
	var e content.Page
	if err := c.Connector.GetOrm().Get(id, &e); err == nil && e.UserId == mchId {
		return &e
	}
	return nil
}

// 根据标识获取页面
func (c *contentRep) GetPageByStringIndent(userId int64, indent string) *content.Page {
	var e content.Page
	if err := c.Connector.GetOrm().GetBy(&e, "user_id=? and str_indent=?", userId, indent); err == nil {
		return &e
	}
	return nil
}

// 删除页面
func (c *contentRep) DeletePage(userId, id int64) error {
	_, err := c.Connector.GetOrm().Delete(content.Page{}, "user_id=? AND id=?", userId, id)
	return err
}

// 保存页面
func (c *contentRep) SavePage(userId int64, v *content.Page) (int64, error) {
	return orm.Save(c.GetOrm(), v, v.Id)
}

// 获取文章数量
func (c *contentRep) GetArticleNumByCategory(categoryId int64) int {
	num := 0
	c.Connector.ExecScalar("SELECT COUNT(0) FROM con_article WHERE cat_id=?",
		&num, categoryId)
	return num
}

// 获取栏目
func (c *contentRep) GetAllArticleCategory() []*content.ArticleCategory {
	list := []*content.ArticleCategory{}
	c.Connector.GetOrm().Select(&list, "")
	return list
}

// 判断栏目是否存在
func (c *contentRep) CategoryExists(indent string, id int64) bool {
	num := 0
	c.Connector.ExecScalar("SELECT COUNT(0) FROM con_category WHERE indent=? and id<>id",
		&num, indent, id)
	return num > 0
}

// 保存栏目
func (c *contentRep) SaveCategory(v *content.ArticleCategory) (int64, error) {
	return orm.Save(c.GetOrm(), v, v.Id)
}

// 删除栏目
func (c *contentRep) DeleteCategory(id int64) error {
	return c.Connector.GetOrm().DeleteByPk(&content.ArticleCategory{}, id)
}

// 获取文章
func (c *contentRep) GetArticleById(id int64) *content.Article {
	e := content.Article{}
	if c.Connector.GetOrm().Get(id, &e) == nil {
		return &e
	}
	return nil
}

// 获取文章列表
func (c *contentRep) GetArticleList(categoryId int64, begin int, end int) []*content.Article {
	list := []*content.Article{}
	c.Connector.GetOrm().SelectByQuery(&content.Article{},
		"cat_id=? LIMIT ?,?", categoryId, begin, end-begin)
	return list
}

// 保存文章
func (c *contentRep) SaveArticle(v *content.Article) (int64, error) {
	return orm.Save(c.GetOrm(), v, v.Id)
}

// 删除文章
func (c *contentRep) DeleteArticle(id int64) error {
	return c.Connector.GetOrm().DeleteByPk(&content.Article{}, id)
}
