/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package product

import (
	"go2o/core/domain/interface/pro_model"
	"go2o/core/infrastructure/domain"
	"sort"
)

var (
	ErrReadonlyCategory = domain.NewError(
		"err_readonly_category", "无权修改系统分类",
	)
	ErrNoSuchCategory = domain.NewError(
		"err_category_not_exist", "分类不存在",
	)

	ErrCategoryCycleReference = domain.NewError(
		"err_category_cycle_reference", "分类上级循环引用")

	ErrHasChildCategories = domain.NewError(
		"err_has_child_categories", "分类包含子分类,无法删除",
	)
	ErrCategoryContainGoods = domain.NewError(
		"err_category_contain_goods", "分类包含商品,无法删除",
	)

	ErrIncorrectCategoryType = domain.NewError(
		"err_category_incorrect_type", "分类类型不允许修改")

	ErrVirtualCatNoUrl = domain.NewError(
		"err_category_virtual_no_url", "虚拟分类必须设置跳转链接")

	ErrCategoryFloorShow = domain.NewError(
		"err_category_floor_show", "非一级分类不能设置首页显示")
)

type (
	ICategory interface {
		// 获取领域编号
		GetDomainId() int
		// 获取值
		GetValue() *Category
		// 设置值
		SetValue(*Category) error
		//todo: 做成界面,同时可后台管理项
		// 获取扩展数据
		GetOption() domain.IOptionStore
		// 保存
		Save() (int, error)
		// 获取子栏目的编号
		GetChildes() []int
	}
	//分类
	Category struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 上级分类
		ParentId int `db:"parent_id"`
		// 产品模型
		ModelId int `db:"prod_model"`
		// 优先级
		Priority int `db:"priority"`
		// 分类名称
		Name string `db:"name"`
		// 是否为虚拟分类
		VirtualCat int `db:"virtual_cat"`
		// 分类链接地址
		CatUrl string `db:"cat_url"`
		// 虚拟分类跳转地址
		RedirectUrl string `db:"redirect_url"`
		// 图标
		Icon string `db:"icon"`
		// 图标坐标
		IconPoint string `db:"icon_xy"`
		// 分类层级
		Level int `db:"level"`
		// 序号
		SortNum int `db:"sort_num"`
		// 是否楼层显示
		FloorShow int `db:"floor_show"`
		// 是否启用
		Enabled int `db:"enabled"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 子分类
		Children []*Category `db:"-"`
	}
	ICategoryRepo interface {
		// 获取系统的栏目服务
		GlobCatService() IGlobCatService

		// 保存分类
		SaveCategory(*Category) (int, error)

		// 检查分类是否关联商品
		CheckContainGoods(vendorId int64, catId int) bool

		// 删除分类
		DeleteCategory(vendorId int64, catId int) error

		// 获取分类
		GetCategory(mchId, id int) *Category

		// 获取所有分类
		GetCategories(mchId int) []*Category
		// 获取关联的品牌
		GetRelationBrands(idArr []int) []*promodel.ProductBrand
	}

	// 公共分类服务
	IGlobCatService interface {
		// 是否只读,当商户共享系统的分类时,
		// 没有修改的权限,即只读!
		ReadOnly() bool
		// 创建分类
		CreateCategory(*Category) ICategory
		// 获取分类
		GetCategory(id int) ICategory
		// 获取所有分类
		GetCategories() []ICategory
		// 删除分类
		DeleteCategory(id int) error
		// 递归获取下级分类
		CategoryTree(parentId int) *Category
		// 获取分类关联的品牌
		RelationBrands(catId int) []*promodel.ProductBrand
	}
)

var (
	C_OptionViewName = "viewName" //显示的视图名称
	C_OptionDescribe = "describe" //描述
)

var _ sort.Interface = new(CategoryList)

type CategoryList []*Category

func (c CategoryList) Len() int {
	return len(c)
}

func (c CategoryList) Less(i, j int) bool {
	return c[i].SortNum < c[j].SortNum ||
		// 如果序号相同,则判断ID
		(c[i].SortNum == c[j].SortNum && c[i].Id < c[j].Id)
}

func (c CategoryList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
