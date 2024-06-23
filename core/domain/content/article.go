/**
 * Copyright 2015 @ 56x.net.
 * name : article
 * author : jarryliu
 * date : 2016-06-27 15:57
 * description :
 * history :
 */
package content

import (
	"strconv"
	"time"

	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/infrastructure/util/collections"
)

var _ content.IArticle = new(articleImpl)

type articleImpl struct {
	_repo     content.IArticleRepo
	_value    *content.Article
	_category *content.Category
	_manager  content.IArticleManager
}

func NewArticle(v *content.Article, m content.IArticleManager,
	rep content.IArticleRepo) content.IArticle {
	return &articleImpl{
		_repo:    rep,
		_value:   v,
		_manager: m,
	}
}

// 获取领域编号
func (a *articleImpl) GetDomainId() int {
	return a._value.Id
}

// 获取值
func (a *articleImpl) GetValue() content.Article {
	return *a._value
}

// SetValue 设置值
func (a *articleImpl) SetValue(v *content.Article) error {
	a._value.Title = v.Title
	a._value.ShortTitle = v.ShortTitle
	a._value.SortNum = v.SortNum
	a._value.Location = v.Location
	a._value.Content = v.Content
	a._value.Thumbnail = v.Thumbnail
	a._value.CatId = v.CatId
	a._value.Flag = v.Flag
	a._value.MchId = v.MchId
	a._value.AccessToken = v.AccessToken
	a._value.Priority = v.Priority
	a._value.UpdateTime = int(time.Now().Unix())

	if a._value.CreateTime == 0 {
		a._value.CreateTime = a._value.UpdateTime
	}
	if a._value.PublisherId <= 0 {
		a._value.PublisherId = v.PublisherId
	}
	return nil
}

// Category 栏目
func (a *articleImpl) Category() *content.Category {
	if a._category == nil {
		a._category = a._manager.GetCategory(int(a._value.CatId))
	}
	return a._category
}

// Save 保存文章
func (a *articleImpl) Save() error {
	if a.Category() == nil {
		return content.NotSetCategory
	}
	_, err := a._repo.Save(a._value)
	return err
}

// Dislike implements content.IArticle.
func (a *articleImpl) Dislike(memberId int) error {
	//todo: 记录会员的点赞记录
	a._value.DislikeCount += 1
	return a.Save()
}

// IncreaseViewCount implements content.IArticle.
func (a *articleImpl) IncreaseViewCount(memberId int, count int) error {
	a._value.ViewCount += count
	return a.Save()
}

// Like implements content.IArticle.
func (a *articleImpl) Like(memberId int) error {
	a._value.LikeCount += 1
	return a.Save()
}

// DeleteArticle 删除文章
func (a *articleImpl) Destory() error {
	return a._repo.Delete(&content.Article{Id: a.GetDomainId()})
}

var _ content.IArticleManager = new(articleManagerImpl)

type articleManagerImpl struct {
	_rep         content.IArticleCategoryRepo
	_artRepo     content.IArticleRepo
	_userId      int64
	categoryList []*content.Category
}

func newArticleManagerImpl(userId int64, rep content.IArticleCategoryRepo, artRepo content.IArticleRepo) content.IArticleManager {
	return &articleManagerImpl{
		_rep:     rep,
		_userId:  userId,
		_artRepo: artRepo,
	}
}

// GetAllCategory 获取所有的栏目
func (a *articleManagerImpl) GetAllCategory() []content.Category {
	if a.categoryList == nil {
		list := a._rep.FindList("")
		l := len(list)
		//如果没有分类,则为系统初始化数据
		if l == 0 && a._userId <= 0 {
			a.initSystemCategory()
			list = a._rep.FindList("")
			l = len(list)
		}
		catList := make([]*content.Category, l)
		for i, v := range list {
			catList[i] = v
		}
		a.categoryList = catList
	}
	return collections.MapList(a.categoryList, func(a *content.Category) content.Category {
		return *a
	})

}

// 初始化系统数据
func (a *articleManagerImpl) initSystemCategory() {
	list := []*content.Category{
		{

			Alias: "about",
			Name:  "关于商城",
		},
		{
			Alias: "wholesale",
			Name:  "批发公告",
		},
		{
			Alias:    "member",
			Name:     "会员公告",
			PermFlag: content.FCategoryOpen,
		},
		{
			Alias: "merchant",
			Name:  "商户公告",
		},
		{
			Alias: "service",
			Name:  "客户服务",
		},
		{
			Alias: "help",
			Name:  "帮助中心",
		},
		{
			Alias:    "news",
			Name:     "新闻资讯",
			PermFlag: content.FCategorySupportPost,
		},
		{
			Alias: "notification",
			Name:  "系统通知",
		},
	}
	for _, v := range list {
		v.PermFlag = v.PermFlag | content.FCategoryInternal
		a.SaveCategory(v)
	}
}

// 获取栏目
func (a *articleManagerImpl) GetCategory(id int) *content.Category {
	v := collections.FindArray(a.GetAllCategory(), func(v content.Category) bool {
		return v.Id == id
	})
	if v.Id > 0 {
		return &v
	}
	return nil
}

// GetCategoryByAlias 根据标识获取文章栏目
func (a *articleManagerImpl) GetCategoryByAlias(alias string) *content.Category {
	v := collections.FindArray(a.GetAllCategory(), func(v content.Category) bool {
		return v.Alias == alias || strconv.Itoa(int(v.Id)) == alias
	})
	if v.Id > 0 {
		return &v
	}
	return nil
}

func (a *articleManagerImpl) SaveCategory(v *content.Category) error {
	exit := a._rep.FindBy("alias = ? and id <> ?", v.Alias, v.Id)
	if exit != nil {
		return content.ErrCategoryAliasExists
	}
	var r *content.Category
	if v.Id > 0 {
		r = a.GetCategory(v.Id)
	} else {
		r = &content.Category{}
	}
	r.Name = v.Name
	r.Alias = v.Alias
	r.Location = v.Location
	r.Title = v.Title
	r.SortNum = v.SortNum
	r.ParentId = v.ParentId
	r.Title = v.Title
	r.Keywords = v.Keywords
	r.Description = v.Description

	// 设置访问权限
	if v.PermFlag > 0 {
		r.PermFlag = v.PermFlag
	}
	r.UpdateTime = time.Now().Unix()
	_, err := a._rep.Save(r)
	if err == nil {
		a.categoryList = nil
	}
	return err
}

// DeleteCategory 删除栏目
func (a *articleManagerImpl) DeleteCategory(id int) error {
	c := a.GetCategory(id)
	if c != nil {
		n := a._artRepo.GetArticleNumByCategory(id)
		if n > 0 {
			return content.ErrCategoryContainArchive
		}
		return a._rep.Delete(&content.Category{Id: id})
	}
	return nil

}

// CreateArticle 创建文章
func (a *articleManagerImpl) CreateArticle(v *content.Article) content.IArticle {
	return NewArticle(v, a, a._artRepo)
}

// GetArticle 获取文章
func (a *articleManagerImpl) GetArticle(id int) content.IArticle {
	v := a._artRepo.Get(id)
	if v != nil {
		return NewArticle(v, a, a._artRepo)
	}
	return nil
}
