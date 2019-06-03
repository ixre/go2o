/**
 * Copyright 2015 @ to2.net.
 * name : item
 * author : jarryliu
 * date : 2016-06-29 09:31
 * description :
 * history :
 */
package product

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrNoSuchProduct *domain.DomainError = domain.NewError(
		"err_product_no_such_product", "产品不存在")
	ErrNoSuchAttr *domain.DomainError = domain.NewError(
		"err_product_no_such_attr", "产品属性不存在")
	ErrNoBrand *domain.DomainError = domain.NewError(
		"err_product_no_brand", "未设置商品品牌")

	ErrVendor *domain.DomainError = domain.NewError(
		"err_not_be_review", "商品供应商不正确")

	ErrNotBeReview *domain.DomainError = domain.NewError(
		"err_not_be_review", "商品还未通过审核")

	ErrItemNameLength *domain.DomainError = domain.NewError(
		"err_item_name_length", "商品标题至少10个字")

	ErrItemIncorrect *domain.DomainError = domain.NewError(
		"err_item_incorrect", "商品已被违规下架")

	ErrNotUploadImage *domain.DomainError = domain.NewError(
		"err_goods_not_upload_image", "请上传商品图片")

	ErrDescribeLength *domain.DomainError = domain.NewError(
		"err_item_describe_length", "商品描述至少20个字符")

	ErrNilRejectRemark *domain.DomainError = domain.NewError(
		"err_item_nil_reject_remark", "原因不能为空")
)

type (
	IProduct interface {
		// 获取聚合根编号
		GetAggregateRootId() int64
		// 获取商品的值
		GetValue() Product
		// 设置产品的值
		SetValue(v *Product) error
		// 设置产品属性
		SetAttr(attr []*Attr) error
		// 获取属性
		Attr() []*Attr
		// 获取销售标签
		//GetSaleLabels() []*Label

		// 保存销售标签
		//SaveSaleLabels([]int) error

		// 设置商品描述
		SetDescribe(describe string) error
		// 保存
		Save() (int64, error)

		// 销毁产品
		Destroy() error
	}

	IProductRepo interface {
		// 创建产品
		CreateProduct(*Product) IProduct
		// 根据产品编号获取货品
		GetProduct(id int64) IProduct
		// 获取货品
		GetProductValue(itemId int64) *Product
		// 根据id获取货品
		GetProductsById(ids ...int32) ([]*Product, error)
		SaveProduct(*Product) (int, error)
		//todo:  到商品
		// 获取在货架上的商品
		GetPagedOnShelvesProduct(supplierId int32, catIds []int32, start, end int) (total int, goods []*Product)
		//todo:  到商品
		// 获取货品销售总数
		GetProductSaleNum(productId int64) int
		// 删除货品
		DeleteProduct(productId int64) error
		// Get Attr
		GetAttr(primary interface{}) *Attr
		// Select Attr
		SelectAttr(where string, v ...interface{}) []*Attr
		// Save Attr
		SaveAttr(v *Attr) (int, error)
		// Delete Attr
		DeleteAttr(primary interface{}) error
		// Batch Delete Attr
		BatchDeleteAttr(where string, v ...interface{}) (int64, error)
	}
)

type (
	// 产品
	Product struct {
		// 编号
		Id int64 `db:"id" auto:"yes" pk:"yes"`
		// 分类
		CatId int32 `db:"cat_id"`
		// 名称
		Name string `db:"name"`
		//供应商编号(暂时同mch_id)
		VendorId int32 `db:"supplier_id"`
		// 品牌编号
		BrandId int32 `db:"brand_id"`
		// 商家编码
		Code string `db:"code"`
		// 图片
		Image string `db:"img"`
		// 描述
		Description string `db:"description"`
		// 上架状态
		//ShelveState int32 `db:"-"`
		// 审核状态
		//ReviewState int32 `db:"-"`
		// 备注
		Remark string `db:"remark"`
		// 状态
		State int32 `db:"state"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 排序编号
		SortNum int32 `db:"sort_num"`
		// 产品属性
		Attr []*Attr `db:"-"`
	}

	// 产品属性
	Attr struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 产品编号
		ProductId int64 `db:"product_id"`
		// 属性编号
		AttrId int32 `db:"attr_id"`
		// 属性值
		AttrData string `db:"attr_data"`
		// 属性文本
		AttrWord string `db:"attr_word"`
	}
)
