/**
 * Copyright 2015 @ z3q.net.
 * name : item
 * author : jarryliu
 * date : 2016-06-29 09:31
 * description :
 * history :
 */
package item

import (
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
)

const (
	// 已下架
	ShelvesDown = 1
	// 已上架
	ShelvesOn = 2
	// 已拒绝上架 (不允许上架)
	ShelvesIncorrect = 3
)

var (
	ErrVendor *domain.DomainError = domain.NewDomainError(
		"err_not_be_review", "商品供应商不正确")

	ErrNotBeReview *domain.DomainError = domain.NewDomainError(
		"err_not_be_review", "商品还未通过审核")

	ErrItemNameLength *domain.DomainError = domain.NewDomainError(
		"err_item_name_length", "商品标题至少10个字")

	ErrItemIncorrect *domain.DomainError = domain.NewDomainError(
		"err_item_incorrect", "商品已被违规下架")

	ErrNotUploadImage *domain.DomainError = domain.NewDomainError(
		"err_goods_not_upload_image", "请上传商品图片")

	ErrDescribeLength *domain.DomainError = domain.NewDomainError(
		"err_item_describe_length", "商品描述至少20个字符")

	ErrNilRejectRemark *domain.DomainError = domain.NewDomainError(
		"err_item_nil_reject_remark", "原因不能为空")
)

type (
	IItemRep interface {
		// 获取货品
		GetValueItem(itemId int32) *Item

		// 根据id获取货品
		GetItemByIds(ids ...int64) ([]*Item, error)

		SaveValueItem(*Item) (int32, error)

		// 获取在货架上的商品
		GetPagedOnShelvesItem(supplierId int32, catIds []int64, start, end int) (total int, goods []*Item)

		// 获取货品销售总数
		GetItemSaleNum(supplierId int32, id int32) int

		// 删除货品
		DeleteItem(supplierId, goodsId int32) error
	}

	// 商品值
	Item struct {
		// 编号
		Id int32 `db:"id" auto:"yes" pk:"yes"`
		// 分类
		CategoryId int32 `db:"category_id"`
		// 名称
		Name string `db:"name"`
		//供应商编号(暂时同mch_id)
		VendorId int32 `db:"supplier_id"`
		// 货号
		GoodsNo string `db:"goods_no"`
		// 小标题
		SmallTitle string `db:"small_title"`
		// 图片
		Image string `db:"img"`
		// 成本价
		Cost float32 `db:"cost"`
		// 单件重量,单位:千克(kg)
		Weight float32 `db:"weight"`
		//定价
		Price float32 `db:"price"`
		//参考销售价
		SalePrice float32 `db:"sale_price"`
		// 运费模板编号
		ExpressTplId int32 `db:"express_tid"`
		// 描述
		Description string `db:"description"`
		// 上架状态
		ShelveState int `db:"shelve_state"`
		// 审核状态
		ReviewState int `db:"review_state"`
		// 备注
		Remark string `db:"remark"`
		// 状态
		State int `db:"state"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)

// 转换包含部分数据的产品值对象
func ParseToPartialValueItem(v *valueobject.Goods) *Item {
	return &Item{
		Id:         v.Item_Id,
		CategoryId: v.CategoryId,
		Name:       v.Name,
		GoodsNo:    v.GoodsNo,
		Image:      v.Image,
		Price:      v.Price,
		SalePrice:  v.SalePrice,
	}
}
