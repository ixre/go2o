/**
 * Copyright 2015 @ 56x.net.
 * name : content_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repos

import (
	contentImpl "github.com/ixre/go2o/core/domain/content"
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

var _ content.IContentRepo = new(contentRepo)

type contentRepo struct {
	db.Connector
	o orm.Orm
}

// 内容仓储
func NewContentRepo(o orm.Orm) content.IContentRepo {
	return &contentRepo{
		Connector: o.Connector(),
		o:         o,
	}
}

// 获取内容
func (c *contentRepo) GetContent(userId int64) content.IContent {
	return contentImpl.NewContent(userId, c)
}

// 根据编号获取页面
func (c *contentRepo) GetPageById(mchId, id int32) *content.Page {
	var e content.Page
	if err := c.o.Get(id, &e); err == nil && e.UserId == int64(mchId) {
		return &e
	}
	return nil
}

// 根据标识获取页面
func (c *contentRepo) GetPageByStringIndent(userId int32, indent string) *content.Page {
	var e content.Page
	if err := c.o.GetBy(&e, "user_id= $1 and str_indent= $2", userId, indent); err == nil {
		return &e
	}
	return nil
}

// 删除页面
func (c *contentRepo) DeletePage(userId, id int32) error {
	_, err := c.o.Delete(content.Page{}, "user_id= $1 AND id= $2", userId, id)
	return err
}

// 保存页面
func (c *contentRepo) SavePage(userId int32, v *content.Page) (int32, error) {
	return orm.I32(orm.Save(c.o, v, int(v.Id)))
}

// 获取文章数量
func (c *contentRepo) GetArticleNumByCategory(categoryId int32) int {
	num := 0
	c.Connector.ExecScalar("SELECT COUNT(0) FROM article_list WHERE cat_id= $1",
		&num, categoryId)
	return num
}

// 获取栏目
func (c *contentRepo) GetAllArticleCategory() []*content.ArticleCategory {
	var list []*content.ArticleCategory
	c.o.Select(&list, "")
	return list
}

// 判断栏目是否存在
func (c *contentRepo) CategoryExists(alias string, id int32) bool {
	num := 0
	c.Connector.ExecScalar("SELECT COUNT(0) FROM article_category WHERE cat_alias= $1 and id <> $2",
		&num, alias, id)
	return num > 0
}

// 保存栏目
func (c *contentRepo) SaveCategory(v *content.ArticleCategory) (int32, error) {
	return orm.I32(orm.Save(c.o, v, int(v.ID)))
}

// 删除栏目
func (c *contentRepo) DeleteCategory(id int32) error {
	return c.o.DeleteByPk(&content.ArticleCategory{}, id)
}

// 获取文章
func (c *contentRepo) GetArticleById(id int32) *content.Article {
	e := content.Article{}
	if c.o.Get(id, &e) == nil {
		return &e
	}
	return nil
}

// 获取文章列表
func (c *contentRepo) GetArticleList(categoryId int32, begin int, end int) []*content.Article {
	list := []*content.Article{}
	c.o.SelectByQuery(&content.Article{},
		"cat_id= $1 LIMIT $3 OFFSET $2", categoryId, begin, end-begin)
	return list
}

// 保存文章
func (c *contentRepo) SaveArticle(v *content.Article) (int32, error) {
	return orm.I32(orm.Save(c.o, v, int(v.ID)))
}

// 删除文章
func (c *contentRepo) DeleteArticle(id int32) error {
	return c.o.DeleteByPk(&content.Article{}, id)
}
