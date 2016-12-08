/**
 * Copyright 2015 @ z3q.net.
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
	value      *content.ArticleCategory
	manager    *articleManagerImpl
}

func NewCategory(v *content.ArticleCategory, m *articleManagerImpl,
	rep content.IContentRepo) content.ICategory {
	return &categoryImpl{
		contentRepo: rep,
		value:      v,
		manager:    m,
	}
}

// 获取领域编号
func (c *categoryImpl) GetDomainId() int32 {
	return c.value.Id
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
	v.Id = c.GetDomainId()
	c.value = v
	return nil
}

// 保存
func (c *categoryImpl) Save() (int32, error) {
	c.value.UpdateTime = time.Now().Unix()
	id, err := c.contentRepo.SaveCategory(c.value)
	if err == nil {
		c.manager._categories = nil
		c.manager._categoryMap = nil
	}
	c.value.Id = id
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
	return a._value.Id
}

// 获取值
func (a *articleImpl) GetValue() content.Article {
	return *a._value
}

// 设置值
func (a *articleImpl) SetValue(v *content.Article) error {
	v.Id = a.GetDomainId()
	a._value = v
	return nil
}

// 栏目
func (a *articleImpl) Category() content.ICategory {
	if a._category == nil {
		a._category = a._manager.GetCategory(a._value.CategoryId)
	}
	return a._category
}

// 保存文章
func (a *articleImpl) Save() (int32, error) {
	if a.Category() == nil {
		return a.GetDomainId(), content.NotSetCategory
	}
	unix := time.Now().Unix()
	a._value.UpdateTime = unix
	if a._value.CreateTime == 0 {
		a._value.CreateTime = unix
	}
	id, err := a._rep.SaveArticle(a._value)
	a._value.Id = id
	return id, err
}

var _ content.IArticleManager = new(articleManagerImpl)

type articleManagerImpl struct {
	_rep         content.IContentRepo
	_categories  []content.ICategory
	_categoryMap map[int32]content.ICategory
	_userId      int32
}

func newArticleManagerImpl(userId int32, rep content.IContentRepo) content.IArticleManager {
	return &articleManagerImpl{
		_rep:         rep,
		_userId:      userId,
		_categoryMap: map[int32]content.ICategory{},
	}
}

// 获取所有的栏目
func (a *articleManagerImpl) GetAllCategory() []content.ICategory {
	if a._categories == nil {
		list := a._rep.GetAllArticleCategory()
		l := len(list)

		//如果没有分类,则为系统初始化数据
		if l == 0 && a._userId <= 0 {
			a.initSystemCategory()
			return []content.ICategory{}
		}

		a._categories = make([]content.ICategory, l)
		a._categoryMap = make(map[int32]content.ICategory)
		for i, v := range list {
			a._categories[i] = NewCategory(v, a, a._rep)
			a._categoryMap[v.Id] = a._categories[i]
		}
	}

	return a._categories
}

// 初始化系统数据
func (a *articleManagerImpl) initSystemCategory() {
	list := []*content.ArticleCategory{
		{
			Id:    0,
			Alias: "mch-notice",
			Name:  "商户公告",
		},
		{
			Id:    0,
			Alias: "mall-activity",
			Name:  "商城活动",
		},
		{
			Id:    0,
			Alias: "mem-notice",
			Name:  "会员公告",
		},
	}

	// 因为保存会清除categoryMap 和categories
	catList := make([]content.ICategory, len(list))
	catMap := make(map[int32]content.ICategory)

	for i, v := range list {
		c := NewCategory(v, a, a._rep)
		if c.GetDomainId() == 0 {
			c.Save() //如果为新建的分类,则保存
		}
		catList[i] = c
		catMap[v.Id] = c
	}
	//赋值
	a._categories = catList
	a._categoryMap = catMap
}

// 获取栏目
func (a *articleManagerImpl) GetCategory(id int32) content.ICategory {
	a.GetAllCategory()
	return a._categoryMap[id]
}

// 创建栏目
func (a *articleManagerImpl) CreateCategory(v *content.ArticleCategory) content.ICategory {
	return NewCategory(v, a, a._rep)
}

// 根据标识获取文章栏目
func (a *articleManagerImpl) GetCategoryByAlias(alias string) content.ICategory {
	a.GetAllCategory()
	for _, v := range a._categories {
		if v2 := v.GetValue(); v2.Alias == alias ||
			strconv.Itoa(int(v2.Id)) == alias {
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
		err := a._rep.DeleteCategory(id)
		if err == nil {
			a._categories = nil
			a._categoryMap = nil
		}
		return err
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
