/**
 * Copyright 2015 @ 56x.net.
 * name : article
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

import (
	"reflect"

	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

// 文章
type (
	// IArticle 文章
	IArticle interface {
		domain.IDomain
		// GetValue 获取值
		GetValue() Article
		// SetValue 设置值
		SetValue(*Article) error
		// Category 栏目
		Category() *Category
		// Save 保存文章
		Save() error
		// Destory 删除文章
		Destory() error
		// 增加浏览次数
		IncreaseViewCount(memberId int, count int) error
		// 喜欢
		Like(memberId int) error
		// 不喜欢
		Dislike(memberId int) error
	}

	// IArticleManager 文章管理器
	IArticleManager interface {
		// GetCategory 获取栏目
		GetCategory(id int) *Category
		// GetCategoryByAlias 根据标识获取文章栏目
		GetCategoryByAlias(alias string) *Category
		// GetAllCategory 获取所有的栏目
		GetAllCategory() []Category
		// SaveCategory 保存栏目
		SaveCategory(v *Category) error
		// DeleteCategory 删除栏目
		DeleteCategory(id int) error
		// CreateArticle 创建文章
		CreateArticle(*Article) IArticle
		// GetArticle 获取文章
		GetArticle(id int) IArticle
	}

	// IArticleCategoryRepo 文章栏目仓储
	IArticleCategoryRepo interface {
		fw.Repository[Category]
	}

	// IArticleManager 文章服务
	IArticleRepo interface {
		fw.Repository[Article]
		// GetContent 获取内容
		GetContent(userId int64) IContentAggregateRoot
		// GetArticleNumByCategory 获取文章数量
		GetArticleNumByCategory(categoryId int) int
	}

	// Category 栏目
	Category struct {
		domain.IValueObject
		//编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		//父类编号,如为一级栏目则为0
		ParentId int `db:"parent_id"`
		// 浏览权限
		PermFlag int `db:"perm_flag"`
		// 名称(唯一)
		Name string `db:"name"`
		// 别名
		Alias string `db:"cat_alias"`
		// 排序编号
		SortNum int `db:"sort_num"`
		// 定位路径（打开栏目页定位到的路径）
		Location string `db:"location"`
		// 页面标题
		Title string `db:"title"`
		// 关键字
		Keywords string `db:"keywords"`
		// 描述
		Description string `db:"describe"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// Article 文章
	Article struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes" json:"id" bson:"id"`
		// 分类编号
		CatId int `db:"cat_id" json:"catId" bson:"catId"`
		// 标题
		Title string `db:"title" json:"title" bson:"title"`
		// ShortTitle
		ShortTitle string `db:"short_title" json:"shortTitle" bson:"shortTitle"`
		// 标志
		Flag int `db:"flag" json:"flag" bson:"flag"`
		// 缩略图
		Thumbnail string `db:"thumbnail" json:"thumbnail" bson:"thumbnail"`
		// 访问密钥
		AccessToken string `db:"access_token" json:"accessToken" bson:"accessToken"`
		// 作者
		PublisherId int `db:"publisher_id" json:"publisherId" bson:"publisherId"`
		// 地址
		Location string `db:"location" json:"location" bson:"location"`
		// 优先级
		Priority int `db:"priority" json:"priority" bson:"priority"`
		// 短标题
		MchId int `db:"mch_id" json:"mchId" bson:"mchId"`
		// 内容
		Content string `db:"content" json:"content" bson:"content"`
		// 标签
		Tags string `db:"tags" json:"tags" bson:"tags"`
		// 喜欢的数量
		LikeCount int `db:"like_count" json:"likeCount" bson:"likeCount"`
		// 不喜欢的数量
		DislikeCount int `db:"dislike_count" json:"dislikeCount" bson:"dislikeCount"`
		// 浏览次数
		ViewCount int `db:"view_count" json:"viewCount" bson:"viewCount"`
		// 排列序号
		SortNum int `db:"sort_num" json:"sortNum" bson:"sortNum"`
		// 创建时间
		CreateTime int `db:"create_time" json:"createTime" bson:"createTime"`
		// 更新时间
		UpdateTime int `db:"update_time" json:"updateTime" bson:"updateTime"`
	}
)

func (c Category) TableName() string {
	return "arc_category"
}

var _ domain.IValueObject = new(Category)

func (c Category) Equal(v interface{}) bool {
	// 判断两个对象值是否相等
	return reflect.DeepEqual(c, v)
}

func (c Article) TableName() string {
	return "arc_article"
}
