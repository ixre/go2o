/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

import "go2o/core/domain/interface/sale/item"

type (
	// 物品
	IItem interface {
		GetDomainId() int64

		// 获取商品的值
		GetValue() item.Item

		// 是否上架
		IsOnShelves() bool

		// 获取销售标签
		GetSaleLabels() []*Label

		// 保存销售标签
		SaveSaleLabels([]int) error

		// 设置商品值
		SetValue(*item.Item) error

		// 设置商品描述
		SetDescribe(describe string) error

		// 设置上架
		SetShelve(state int, remark string) error

		// 审核
		Review(pass bool, remark string) error

		// 标记为违规
		Incorrect(remark string) error

		// 保存
		Save() (int64, error)
	}

	// 货品服务
	IItemManager interface {
		// 创建产品
		CreateItem(*item.Item) IItem

		// 根据产品编号获取货品
		GetItem(id int64) IItem

		// 删除货品
		DeleteItem(id int64) error
	}
)
