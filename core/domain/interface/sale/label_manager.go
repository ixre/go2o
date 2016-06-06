/**
 * Copyright 2015 @ z3q.net.
 * name : label_manager
 * author : jarryliu
 * date : 2016-06-06 21:37
 * description :
 * history :
 */
package sale

import "go2o/core/domain/interface/valueobject"

type (
	// 销售标签
	SaleLabel struct {
		Id int `db:"id" auto:"yes" pk:"yes"`

		// 商户编号
		MerchantId int `db:"mch_id"`

		// 标签代码
		TagCode string `db:"tag_code"`

		// 标签名
		TagName string `db:"tag_name"`

		// 商品的遮盖图
		LabelImage string `db:"label_image"`

		// 是否启用
		Enabled int `db:"enabled"`
	}

	// 销售标签接口
	ISaleLabel interface {
		GetDomainId() int

		// 获取值
		GetValue() *SaleLabel

		// 设置值
		SetValue(v *SaleLabel) error

		// 保存
		Save() (int, error)

		// 是否为系统内置
		System() bool

		// 获取标签下的商品
		GetValueGoods(sortBy string, begin, end int) []*valueobject.Goods

		// 获取标签下的分页商品
		GetPagedValueGoods(sortBy string, begin, end int) (total int, goods []*valueobject.Goods)
	}

	ISaleLabelRep interface {
		// 创建销售标签
		CreateSaleLabel(v *SaleLabel) ISaleLabel

		// 获取所有的销售标签
		GetAllValueSaleLabels(merchantId int) []*SaleLabel

		// 获取销售标签值
		GetValueSaleLabel(merchantId int, tagId int) *SaleLabel

		// 根据Code获取销售标签
		GetSaleLabelByCode(merchantId int, code string) *SaleLabel

		// 删除销售标签
		DeleteSaleLabel(merchantId int, id int) error

		// 获取销售标签
		GetSaleLabel(merchantId int, tagId int) ISaleLabel

		// 保存销售标签
		SaveSaleLabel(merchantId int, v *SaleLabel) (int, error)

		// 获取商品
		GetValueGoodsBySaleLabel(merchantId, tagId int, sortBy string,
			begin, end int) []*valueobject.Goods

		// 获取分页商品
		GetPagedValueGoodsBySaleLabel(merchantId, tagId int, sortBy string,
			begin, end int) (int, []*valueobject.Goods)

		// 获取商品的销售标签
		GetItemSaleLabels(itemId int) []*SaleLabel

		// 清理商品的销售标签
		CleanItemSaleLabels(itemId int) error

		// 保存商品的销售标签
		SaveItemSaleLabels(itemId int, tagIds []int) error
	}
)
