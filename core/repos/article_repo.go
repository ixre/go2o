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
	"errors"

	impl "github.com/ixre/go2o/core/domain/content"
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

var _ content.IArticleRepo = new(contentRepo)
var _ content.IPageRepo = new(pageRepoImpl)

type contentRepo struct {
	fw.Repository[content.Article]
	db.Connector
	o        orm.Orm
	pageRepo content.IPageRepo
	catRepo  content.IArticleCategoryRepo
}

// 内容仓储
func NewArticleRepo(o orm.Orm,
	catRepo content.IArticleCategoryRepo,
	pageRepo content.IPageRepo) content.IArticleRepo {
	return &contentRepo{
		Connector: o.Connector(),
		pageRepo:  pageRepo,
		catRepo:   catRepo,
		o:         o,
	}
}

// 获取内容
func (c *contentRepo) GetContent(userId int64) content.IContentAggregateRoot {
	return impl.NewContent(userId, c, c.catRepo, c.pageRepo)
}

// 获取文章数量
func (c *contentRepo) GetArticleNumByCategory(categoryId int) int {
	num := 0
	c.Connector.ExecScalar("SELECT COUNT(1) FROM article_list WHERE cat_id= $1",
		&num, categoryId)
	return num
}

var _ content.IArticleCategoryRepo = new(articleCategoryRepo)

type articleCategoryRepo struct {
	fw.BaseRepository[content.Category]
}

func NewArticleCategoryRepo(o fw.ORM) content.IArticleCategoryRepo {
	s := &articleCategoryRepo{}
	s.ORM = o
	return s
}

var _ content.IPageRepo = new(pageRepoImpl)

type pageRepoImpl struct {
	fw.BaseRepository[content.Page]
}

func NewPageRepo(o fw.ORM) content.IPageRepo {
	s := &pageRepoImpl{}
	s.ORM = o
	return s
}

// 根据编号获取页面
func (p *pageRepoImpl) GetPageById(zoneId, id int) content.IPage {
	v := p.Get(id)
	if v != nil && v.UserId == zoneId {
		return impl.NewPage(zoneId, p, v)
	}
	return nil
}

// 根据标识获取页面
func (p *pageRepoImpl) GetPageByCode(zoneId int, code string) content.IPage {
	v := p.FindBy("user_id = ? and code = ?", zoneId, code)
	if v != nil && v.UserId == zoneId {
		return impl.NewPage(zoneId, p, v)
	}
	return nil
}

// 删除页面
func (p *pageRepoImpl) DeletePage(zoneId, id int) error {
	_, err := p.DeleteBy("user_id= ? AND id= ?", zoneId, id)
	return err
}

// 保存页面
func (p *pageRepoImpl) SavePage(zoneId int, v *content.Page) error {
	if v.UserId != zoneId {
		return errors.New("zone id not match")
	}
	_, err := p.Save(v)
	return err
}
