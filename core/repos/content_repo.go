/**
 * Copyright 2015 @ z3q.net.
 * name : content_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repos

import (
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	contentImpl "go2o/core/domain/content"
	"go2o/core/domain/interface/content"
)

var _ content.IContentRepo = new(contentRepo)

type contentRepo struct {
	db.Connector
}

// 内容仓储
func NewContentRepo(c db.Connector) content.IContentRepo {
	return &contentRepo{
		Connector: c,
	}
}

// 获取内容
func (c *contentRepo) GetContent(userId int32) content.IContent {
	return contentImpl.NewContent(userId, c)
}

// 根据编号获取页面
func (c *contentRepo) GetPageById(mchId, id int32) *content.Page {
	var e content.Page
	if err := c.Connector.GetOrm().Get(id, &e); err == nil && e.UserId == mchId {
		return &e
	}
	return nil
}

// 根据标识获取页面
func (c *contentRepo) GetPageByStringIndent(userId int32, indent string) *content.Page {
	var e content.Page
	if err := c.Connector.GetOrm().GetBy(&e, "user_id=? and str_indent=?", userId, indent); err == nil {
		return &e
	}
	return nil
}

// 删除页面
func (c *contentRepo) DeletePage(userId, id int32) error {
	_, err := c.Connector.GetOrm().Delete(content.Page{}, "user_id=? AND id=?", userId, id)
	return err
}

// 保存页面
func (c *contentRepo) SavePage(userId int32, v *content.Page) (int32, error) {
	return orm.I32(orm.Save(c.GetOrm(), v, int(v.Id)))
}

// 获取文章数量
func (c *contentRepo) GetArticleNumByCategory(categoryId int32) int {
	num := 0
	c.Connector.ExecScalar("SELECT COUNT(0) FROM article_list WHERE cat_id=?",
		&num, categoryId)
	return num
}

// 获取栏目
func (c *contentRepo) GetAllArticleCategory() []*content.ArticleCategory {
	list := []*content.ArticleCategory{}
	c.Connector.GetOrm().Select(&list, "")
	return list
}

// 判断栏目是否存在
func (c *contentRepo) CategoryExists(alias string, id int32) bool {
	num := 0
	c.Connector.ExecScalar("SELECT COUNT(0) FROM article_category WHERE cat_alias=? and id<>id",
		&num, alias, id)
	return num > 0
}

// 保存栏目
func (c *contentRepo) SaveCategory(v *content.ArticleCategory) (int32, error) {
	return orm.I32(orm.Save(c.GetOrm(), v, int(v.ID)))
}

// 删除栏目
func (c *contentRepo) DeleteCategory(id int32) error {
	return c.Connector.GetOrm().DeleteByPk(&content.ArticleCategory{}, id)
}

// 获取文章
func (c *contentRepo) GetArticleById(id int32) *content.Article {
	e := content.Article{}
	if c.Connector.GetOrm().Get(id, &e) == nil {
		return &e
	}
	return nil
}

// 获取文章列表
func (c *contentRepo) GetArticleList(categoryId int32, begin int, end int) []*content.Article {
	list := []*content.Article{}
	c.Connector.GetOrm().SelectByQuery(&content.Article{},
		"cat_id=? LIMIT ?,?", categoryId, begin, end-begin)
	return list
}

// 保存文章
func (c *contentRepo) SaveArticle(v *content.Article) (int32, error) {
	return orm.I32(orm.Save(c.GetOrm(), v, int(v.ID)))
}

// 删除文章
func (c *contentRepo) DeleteArticle(id int32) error {
	return c.Connector.GetOrm().DeleteByPk(&content.Article{}, id)
}
