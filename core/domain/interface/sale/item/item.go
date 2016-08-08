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

var (
	ErrVendor *domain.DomainError = domain.NewDomainError(
		"err_not_be_review", "商品供应商不正确")

	ErrNotBeReview *domain.DomainError = domain.NewDomainError(
		"err_not_be_review", "商品还未通过审核")

	ErrNotUploadImage *domain.DomainError = domain.NewDomainError(
		"err_goods_not_upload_image", "请上传商品图片")
)

type (
	IItemRep interface {
		// 获取货品
		GetValueItem(itemId int) *Item

		// 根据id获取货品
		GetItemByIds(ids ...int) ([]*Item, error)

		SaveValueItem(*Item) (int, error)

		// 获取在货架上的商品
		GetPagedOnShelvesItem(supplierId int, catIds []int, start, end int) (total int, goods []*Item)

		// 获取货品销售总数
		GetItemSaleNum(supplierId int, id int) int

		// 删除货品
		DeleteItem(supplierId, goodsId int) error
	}

	// 商品值
	Item struct {
		// 编号
		Id int `db:"id" auto:"yes" pk:"yes"`
		// 分类
		CategoryId int `db:"category_id"`
		// 名称
		Name string `db:"name"`
		//供应商编号(暂时同mch_id)
		VendorId int `db:"supplier_id"`
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
		ExpressTplId int `db:"express_tid"`
		// 供应门店 //todo: 去掉
		ApplySubs string `db:"apply_subs"`
		//简单备注,如:(限时促销),todo: 去掉
		Remark string `db:"remark"`
		// 描述
		Description string `db:"description"`
		// 是否上架,1为上架
		OnShelves int `db:"on_shelves"`
		// 是否审核
		HasReview int `db:"has_review"`
		// 是否审核通过
		ReviewPass int `db:"review_pass"`
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
