/**
 * Copyright 2015 @ z3q.net.
 * name : sale_goods.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package item

import (
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
)

type (
	// 商品
	IGoodsItem interface {
		// 获取聚合根编号
		GetAggregateRootId() int32

		// 商品快照
		SnapshotManager() ISnapshotManager

		// 获取货品
		GetItem() product.IProduct

		// 设置值
		GetValue() *GoodsItem

		// 获取包装过的商品信息
		GetPackedValue() *valueobject.Goods

		// 获取促销信息
		GetPromotions() []promotion.IPromotion

		// 获取促销价
		GetPromotionPrice(level int32) float32

		// 获取会员价销价,返回是否有会原价及价格
		GetLevelPrice(level int32) (bool, float32)

		// 获取促销描述
		GetPromotionDescribe() map[string]string

		// 获取会员价
		GetLevelPrices() []*MemberPrice

		// 保存会员价
		SaveLevelPrice(*MemberPrice) (int32, error)

		// 设置值
		SetValue(*GoodsItem) error

		// 保存
		Save() (int32, error)

		// 更新销售数量,扣减库存
		AddSalesNum(quantity int32) error

		// 取消销售
		CancelSale(quantity int32, orderNo string) error

		// 占用库存
		TakeStock(quantity int32) error

		// 释放库存
		FreeStock(quantity int32) error

		//// 生成快照
		//GenerateSnapshot() (int64, error)
		//
		//// 获取最新的快照
		//GetLatestSnapshot() *goods.GoodsSnapshot

		// 删除商品
		Destroy() error
	}

	// 简单商品信息
	SimpleGoods struct {
		GoodsId    int32  `json:"id"`
		GoodsImage string `json:"img"`
		Name       string `json:"name"`
		Quantity   string `json:"qty"`
	}
)
