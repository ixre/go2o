/**
 * Copyright 2015 @ to2.net.
 * name : article
 * author : jarryliu
 * date : 2016-06-27 15:57
 * description :
 * history :
 */
package content

import (
	"go2o/core/domain/interface/content"
	"strconv"
	"time"
)

var _ content.ICategory = new(categoryImpl)

type categoryImpl struct {
	contentRepo content.IContentRepo
	value       *content.ArticleCategory
	manager     *articleManagerImpl
}

func NewCategory(v *content.ArticleCategory, m *articleManagerImpl,
	rep content.IContentRepo) content.ICategory {
	return &categoryImpl{
		contentRepo: rep,
		value:       v,
		manager:     m,
	}
}

// 获取领域编号
func (c *categoryImpl) GetDomainId() int32 {
	return c.value.ID
}

// 获取文章数量
func (c *categoryImpl) ArticleNum() int {
	return c.contentRepo.GetArticleNumByCategory(c.GetDomainId())
}

// 获取值
func (c *categoryImpl) GetValue() content.ArticleCategory {
	return *c.value
}

// 设置值
func (c *categoryImpl) SetValue(v *content.ArticleCategory) error {
	if c.contentRepo.CategoryExists(c.value.Alias, c.GetDomainId()) {
		return content.ErrCategoryAliasExists
	}
	v.ID = c.GetDomainId()
	c.value.Name = v.Name
	c.value.Alias = v.Alias
	c.value.Location = v.Location
	c.value.Title = v.Title
	c.value.SortNum = v.SortNum
	c.value.ParentId = v.ParentId
	c.value.Title = v.Title
	c.value.Keywords = v.Keywords
	c.value.Description = v.Description
	// 设置访问权限
	if v.PermFlag > 0 {
		c.value.PermFlag = v.PermFlag
	}
	if c.value.PermFlag <= 0 {
		c.value.PermFlag = content.PermAll
	}
	return nil
}

// 保存
func (c *categoryImpl) Save() (int32, error) {
	c.value.UpdateTime = time.Now().Unix()
	id, err := c.contentRepo.SaveCategory(c.value)
	if err == nil {
		c.value.ID = id
	}
	return id, err
}

var _ content.IArticle = new(articleImpl)

type articleImpl struct {
	_rep      content.IContentRepo
	_value    *content.Article
	_category content.ICategory
	_manager  content.IArticleManager
}

func NewArticle(v *content.Article, m content.IArticleManager,
	rep content.IContentRepo) content.IArticle {
	return &articleImpl{
		_rep:     rep,
		_value:   v,
		_manager: m,
	}
}

// 获取领域编号
func (a *articleImpl) GetDomainId() int32 {
	return a._value.ID
}

// 获取值
func (a *articleImpl) GetValue() content.Article {
	return *a._value
}

// 设置值
func (a *articleImpl) SetValue(v *content.Article) error {
	a._value.Title = v.Title
	a._value.SmallTitle = v.SmallTitle
	a._value.SortNum = v.SortNum
	a._value.Location = v.Location
	a._value.Content = v.Content
	a._value.Thumbnail = v.Thumbnail
	a._value.CatId = v.CatId
	a._value.AccessKey = v.AccessKey
	a._value.Priority = v.Priority
	a._value.UpdateTime = time.Now().Unix()

	if a._value.CreateTime == 0 {
		a._value.CreateTime = a._value.UpdateTime
	}
	if a._value.PublisherId <= 0 {
		a._value.PublisherId = v.PublisherId
	}
	return nil
}

// 栏目
func (a *articleImpl) Category() content.ICategory {
	if a._category == nil {
		a._category = a._manager.GetCategory(a._value.CatId)
	}
	return a._category
}

// 保存文章
func (a *articleImpl) Save() (int32, error) {
	if a.Category() == nil {
		return a.GetDomainId(), content.NotSetCategory
	}
	id, err := a._rep.SaveArticle(a._value)
	a._value.ID = id
	return id, err
}

var _ content.IArticleManager = new(articleManagerImpl)

type articleManagerImpl struct {
	_rep    content.IContentRepo
	_userId int32
}

func newArticleManagerImpl(userId int32, rep content.IContentRepo) content.IArticleManager {
	return &articleManagerImpl{
		_rep:    rep,
		_userId: userId,
	}
}

// 获取所有的栏目
func (a *articleManagerImpl) GetAllCategory() []content.ICategory {
	list := a._rep.GetAllArticleCategory()
	l := len(list)
	//如果没有分类,则为系统初始化数据
	if l == 0 && a._userId <= 0 {
		a.initSystemCategory()
		list = a._rep.GetAllArticleCategory()
		l = len(list)
	}
	catList := make([]content.ICategory, l)
	for i, v := range list {
		catList[i] = NewCategory(v, a, a._rep)
	}
	return catList
}

// 初始化系统数据
func (a *articleManagerImpl) initSystemCategory() {
	list := []*content.ArticleCategory{
		{
			ID:       0,
			Alias:    "news",
			Name:     "商城公告",
			PermFlag: content.PermAll,
		},
		{
			ID:       0,
			Alias:    "wholesale",
			Name:     "批发中心公告",
			PermFlag: content.PermAll,
		},
		{
			ID:       0,
			Alias:    "member",
			Name:     "会员公告",
			PermFlag: content.PermMember,
		},
		{
			ID:       0,
			Alias:    "merchant",
			Name:     "商户公告",
			PermFlag: content.PermVendor,
		},
		{
			ID:       0,
			Alias:    "service",
			Name:     "客户服务",
			PermFlag: content.PermAll,
		},
		{
			ID:       0,
			Alias:    "help",
			Name:     "帮助中心",
			PermFlag: content.PermAll,
		},
		{
			ID:       0,
			Alias:    "about",
			Name:     "关于商城",
			PermFlag: content.PermAll,
		},
	}
	for _, v := range list {
		c := NewCategory(v, a, a._rep)
		if c.GetDomainId() == 0 {
			c.Save() //如果为新建的分类,则保存
		}
	}
}

// 获取栏目
func (a *articleManagerImpl) GetCategory(id int32) content.ICategory {
	list := a.GetAllCategory()
	for _, v := range list {
		if v.GetDomainId() == id {
			return v
		}
	}
	return nil
}

// 创建栏目
func (a *articleManagerImpl) CreateCategory(v *content.ArticleCategory) content.ICategory {
	return NewCategory(v, a, a._rep)
}

// 根据标识获取文章栏目
func (a *articleManagerImpl) GetCategoryByAlias(alias string) content.ICategory {
	list := a.GetAllCategory()
	for _, v := range list {
		if v2 := v.GetValue(); v2.Alias == alias ||
			strconv.Itoa(int(v2.ID)) == alias {
			return v
		}
	}
	return nil
}

// 删除栏目
func (a *articleManagerImpl) DelCategory(id int32) error {
	c := a.GetCategory(id)
	if c != nil {
		if c.ArticleNum() > 0 {
			return content.ErrCategoryContainArchive
		}
		return a._rep.DeleteCategory(id)
	}
	return nil

}

// 创建文章
func (a *articleManagerImpl) CreateArticle(v *content.Article) content.IArticle {
	return NewArticle(v, a, a._rep)
}

// 获取文章
func (a *articleManagerImpl) GetArticle(id int32) content.IArticle {
	v := a._rep.GetArticleById(id)
	if v != nil {
		return NewArticle(v, a, a._rep)
	}
	return nil
}

// 获取文章列表
func (a *articleManagerImpl) GetArticleList(categoryId int32,
	begin, end int) []*content.Article {
	return a._rep.GetArticleList(categoryId, begin, end)
}

// 删除文章
func (a *articleManagerImpl) DeleteArticle(id int32) error {
	return a._rep.DeleteArticle(id)
}
