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
	Label struct {
		Id int64 `db:"id" auto:"yes" pk:"yes"`

		// 商户编号
		MerchantId int64 `db:"mch_id"`

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
		GetDomainId() int64

		// 获取值
		GetValue() *Label

		// 设置值
		SetValue(v *Label) error

		// 保存
		Save() (int64, error)

		// 是否为系统内置
		System() bool

		// 获取标签下的商品
		GetValueGoods(sortBy string, begin, end int) []*valueobject.Goods

		// 获取标签下的分页商品
		GetPagedValueGoods(sortBy string, begin, end int) (total int,
			goods []*valueobject.Goods)
	}

	// 标签服务
	ILabelManager interface {
		// 获取所有的销售标签
		GetAllSaleLabels() []ISaleLabel

		// 初始化销售标签
		InitSaleLabels() error

		// 获取销售标签
		GetSaleLabel(id int64) ISaleLabel

		// 根据Code获取销售标签
		GetSaleLabelByCode(code string) ISaleLabel

		// 创建销售标签
		CreateSaleLabel(v *Label) ISaleLabel

		// 删除销售标签
		DeleteSaleLabel(id int64) error
	}

	ISaleLabelRep interface {
		// 创建销售标签
		CreateSaleLabel(v *Label) ISaleLabel

		// 获取所有的销售标签
		GetAllValueSaleLabels(mchId int64) []*Label

		// 获取销售标签值
		GetValueSaleLabel(mchId int64, tagId int64) *Label

		// 根据Code获取销售标签
		GetSaleLabelByCode(mchId int64, code string) *Label

		// 删除销售标签
		DeleteSaleLabel(mchId int64, id int64) error

		// 获取销售标签
		GetSaleLabel(mchId int64, tagId int64) ISaleLabel

		// 保存销售标签
		SaveSaleLabel(mchId int64, v *Label) (int64, error)

		// 获取商品
		GetValueGoodsBySaleLabel(mchId int64, tagId int64, sortBy string,
			begin, end int) []*valueobject.Goods

		// 获取分页商品
		GetPagedValueGoodsBySaleLabel(mchId int64, tagId int64, sortBy string,
			begin, end int) (int, []*valueobject.Goods)

		// 获取商品的销售标签
		GetItemSaleLabels(itemId int64) []*Label

		// 清理商品的销售标签
		CleanItemSaleLabels(itemId int64) error

		// 保存商品的销售标签
		SaveItemSaleLabels(itemId int64, tagIds []int) error
	}
)
