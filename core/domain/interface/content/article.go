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
		// GetDomainId 获取领域编号
		GetDomainId() int32
		// GetValue 获取值
		GetValue() Article
		// SetValue 设置值
		SetValue(*Article) error
		// Category 栏目
		Category() *Category
		// Save 保存文章
		Save() (int32, error)
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
		// DelCategory 删除栏目
		DelCategory(id int) error
		// CreateArticle 创建文章
		CreateArticle(*Article) IArticle
		// GetArticle 获取文章
		GetArticle(id int32) IArticle
		// DeleteArticle 删除文章
		DeleteArticle(id int32) error
	}

	// IArticleCategoryRepo 文章栏目仓储
	IArticleCategoryRepo interface {
		fw.Repository[Category]
	}

	// IArticleManager 文章服务
	IArticleRepo interface {
		// GetContent 获取内容
		GetContent(userId int64) IContentAggregateRoot
		// GetArticleNumByCategory 获取文章数量
		GetArticleNumByCategory(categoryId int) int

		// 获取文章
		GetArticleById(id int32) *Article
		// 保存文章
		SaveArticle(v *Article) (int32, error)
		// 删除文章
		DeleteArticle(id int32) error
	}

	// Category 栏目
	Category struct {
		domain.IValueObject
		//编号
		ID int `db:"id" pk:"yes" auto:"yes"`
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
		ID int32 `db:"id" auto:"yes" pk:"yes"`
		// 栏目编号
		CatId int32 `db:"cat_id"`
		// 标题
		Title string `db:"title"`
		// 小标题
		SmallTitle string `db:"small_title"`
		// 文章附图
		Thumbnail string `db:"thumbnail"`
		// 重定向URL
		Location string `db:"location"`
		// 优先级,优先级越高，则置顶
		Priority int `db:"priority"`
		// 浏览钥匙
		AccessKey string `db:"access_key"`
		// 作者
		PublisherId int `db:"publisher_id"`
		// 文档内容
		Content string `db:"content"`
		// 标签（关键词）
		Tags string `db:"tags"`
		// 显示次数
		ViewCount int `db:"view_count"`
		// 排序序号
		SortNum int `db:"sort_num"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 最后修改时间
		UpdateTime int64 `db:"update_time"`
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
