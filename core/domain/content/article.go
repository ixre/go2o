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
	_contentRep content.IContentRep
	_value      *content.ArticleCategory
	_manager    *articleManagerImpl
}

func NewCategory(v *content.ArticleCategory, m *articleManagerImpl,
	rep content.IContentRep) content.ICategory {
	return &categoryImpl{
		_contentRep: rep,
		_value:      v,
		_manager:    m,
	}
}

// 获取领域编号
func (this *categoryImpl) GetDomainId() int {
	return this._value.Id
}

// 获取文章数量
func (this *categoryImpl) ArticleNum() int {
	return this._contentRep.GetArticleNumByCategory(this.GetDomainId())
}

// 获取值
func (this *categoryImpl) GetValue() content.ArticleCategory {
	return *this._value
}

// 设置值
func (this *categoryImpl) SetValue(v *content.ArticleCategory) error {
	if this._contentRep.CategoryExists(this._value.Alias, this.GetDomainId()) {
		return content.ErrCategoryAliasExists
	}
	v.Id = this.GetDomainId()
	this._value = v
	return nil
}

// 保存
func (this *categoryImpl) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix()
	id, err := this._contentRep.SaveCategory(this._value)
	if err == nil {
		this._manager._categories = nil
		this._manager._categoryMap = nil
	}
	this._value.Id = id
	return id, err
}

var _ content.IArticle = new(articleImpl)

type articleImpl struct {
	_rep      content.IContentRep
	_value    *content.Article
	_category content.ICategory
	_manager  content.IArticleManager
}

func NewArticle(v *content.Article, m content.IArticleManager,
	rep content.IContentRep) content.IArticle {
	return &articleImpl{
		_rep:     rep,
		_value:   v,
		_manager: m,
	}
}

// 获取领域编号
func (this *articleImpl) GetDomainId() int {
	return this._value.Id
}

// 获取值
func (this *articleImpl) GetValue() content.Article {
	return *this._value
}

// 设置值
func (this *articleImpl) SetValue(v *content.Article) error {
	v.Id = this.GetDomainId()
	this._value = v
	return nil
}

// 栏目
func (this *articleImpl) Category() content.ICategory {
	if this._category == nil {
		this._category = this._manager.GetCategory(this._value.CategoryId)
	}
	return this._category
}

// 保存文章
func (this *articleImpl) Save() (int, error) {
	if this.Category() == nil {
		return this.GetDomainId(), content.NotSetCategory
	}
	unix := time.Now().Unix()
	this._value.UpdateTime = unix
	if this._value.CreateTime == 0 {
		this._value.CreateTime = unix
	}
	id, err := this._rep.SaveArticle(this._value)
	this._value.Id = id
	return id, err
}

var _ content.IArticleManager = new(articleManagerImpl)

type articleManagerImpl struct {
	_rep         content.IContentRep
	_categories  []content.ICategory
	_categoryMap map[int]content.ICategory
	_userId      int
}

func newArticleManagerImpl(userId int, rep content.IContentRep) content.IArticleManager {
	return &articleManagerImpl{
		_rep:         rep,
		_userId:      userId,
		_categoryMap: map[int]content.ICategory{},
	}
}

// 获取所有的栏目
func (this *articleManagerImpl) GetAllCategory() []content.ICategory {
	if this._categories == nil {
		list := this._rep.GetAllArticleCategory()
		l := len(list)

		//如果没有分类,则为系统初始化数据
		if l == 0 && this._userId <= 0 {
			this.initSystemCategory()
			return []content.ICategory{}
		}

		this._categories = make([]content.ICategory, l)
		this._categoryMap = make(map[int]content.ICategory)
		for i, v := range list {
			this._categories[i] = NewCategory(v, this, this._rep)
			this._categoryMap[v.Id] = this._categories[i]
		}
	}

	return this._categories
}

// 初始化系统数据
func (this *articleManagerImpl) initSystemCategory() {
	list := []*content.ArticleCategory{
		&content.ArticleCategory{
			Id:    0,
			Alias: "mch-notice",
			Name:  "商户公告",
		},
		&content.ArticleCategory{
			Id:    0,
			Alias: "mall-activity",
			Name:  "商城活动",
		},
		&content.ArticleCategory{
			Id:    0,
			Alias: "mem-notice",
			Name:  "会员公告",
		},
	}

	// 因为保存会清除categoryMap 和categories
	catList := make([]content.ICategory, len(list))
	catMap := make(map[int]content.ICategory)

	for i, v := range list {
		c := NewCategory(v, this, this._rep)
		if c.GetDomainId() == 0 {
			c.Save() //如果为新建的分类,则保存
		}
		catList[i] = c
		catMap[v.Id] = c
	}
	//赋值
	this._categories = catList
	this._categoryMap = catMap
}

// 获取栏目
func (this *articleManagerImpl) GetCategory(id int) content.ICategory {
	this.GetAllCategory()
	return this._categoryMap[id]
}

// 创建栏目
func (this *articleManagerImpl) CreateCategory(v *content.ArticleCategory) content.ICategory {
	return NewCategory(v, this, this._rep)
}

// 根据标识获取文章栏目
func (this *articleManagerImpl) GetCategoryByAlias(alias string) content.ICategory {
	this.GetAllCategory()
	for _, v := range this._categories {
		if v2 := v.GetValue(); v2.Alias == alias || strconv.Itoa(v2.Id) == alias {
			return v
		}
	}
	return nil
}

// 删除栏目
func (this *articleManagerImpl) DelCategory(id int) error {
	c := this.GetCategory(id)
	if c != nil {
		if c.ArticleNum() > 0 {
			return content.ErrCategoryContainArchive
		}
		err := this._rep.DeleteCategory(id)
		if err == nil {
			this._categories = nil
			this._categoryMap = nil
		}
		return err
	}
	return nil

}

// 创建文章
func (this *articleManagerImpl) CreateArticle(v *content.Article) content.IArticle {
	return NewArticle(v, this, this._rep)
}

// 获取文章
func (this *articleManagerImpl) GetArticle(id int) content.IArticle {
	v := this._rep.GetArticleById(id)
	if v != nil {
		return NewArticle(v, this, this._rep)
	}
	return nil
}

// 获取文章列表
func (this *articleManagerImpl) GetArticleList(categoryId int,
	begin, end int) []*content.Article {
	return this._rep.GetArticleList(categoryId, begin, end)
}

// 删除文章
func (this *articleManagerImpl) DeleteArticle(id int) error {
	return this._rep.DeleteArticle(id)
}
