/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

import (
	"go2o/core/infrastructure/domain"
	"sort"
)

type (
	ICategory interface {
		GetDomainId() int

		GetValue() *Category

		GetOption() domain.IOptionStore

		SetValue(*Category) error

		Save() (int, error)

		// 获取子栏目的编号
		GetChildId() []int
	}
	//分类
	Category struct {
		Id int `db:"id" auto:"yes" pk:"yes"`
		//父分类
		ParentId int `db:"parent_id"`
		//供应商编号
		MerchantId int `db:"merchant_id"`
		//名称
		Name       string `db:"name"`
		SortNumber int    `db:"sort_number"`
		//层级,用于判断2个分类是否为同一级
		Level       int         `db:"level"`
		Icon        string      `db:"icon"`
		Url         string      `db:"url"`
		CreateTime  int64       `db:"create_time"`
		Enabled     int         `db:"enabled"`
		Description string      `db:"description"`
		Child       []*Category `db:"-"`
	}
	ICategoryRep interface {
		// 保存分类
		SaveCategory(*Category) (int, error)

		DeleteCategory(mchId, id int) error

		GetCategory(mchId, id int) *Category

		GetCategories(mchId int) CategoryList

		// 获取与栏目相关的栏目
		GetRelationCategories(mchId, categoryId int) CategoryList

		// 获取子栏目
		GetChildCategories(mchId, categoryId int) CategoryList
	}

	// 分类服务
	ICategoryManager interface {
		// 是否只读,当商户共享系统的分类时,
		// 没有修改的权限,即只读!
		ReadOnly() bool

		// 创建分类
		CreateCategory(*Category) ICategory

		// 获取分类
		GetCategory(int) ICategory

		// 获取所有分类
		GetCategories() []ICategory

		// 删除分类
		DeleteCategory(int) error
	}
)

var (
	C_OptionViewName string = "viewName" //显示的视图名称
	C_OptionDescribe string = "describe" //描述
)

var _ sort.Interface = new(CategoryList)

type CategoryList []*Category

func (c CategoryList) Len() int {
	return len(c)
}

func (c CategoryList) Less(i, j int) bool {
	return c[i].SortNumber < c[j].SortNumber ||
		// 如果序号相同,则判断ID
		(c[i].SortNumber == c[j].SortNumber && c[i].Id < c[j].Id)
}

func (c CategoryList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
