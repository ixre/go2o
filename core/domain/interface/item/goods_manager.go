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
	IGoods interface {
		// 获取领域对象编号
		GetDomainId() int32

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
		AddSalesNum(quantity int) error

		// 取消销售
		CancelSale(quantity int, orderNo string) error

		// 占用库存
		TakeStock(quantity int) error

		// 释放库存
		FreeStock(quantity int) error

		//// 生成快照
		//GenerateSnapshot() (int64, error)
		//
		//// 获取最新的快照
		//GetLatestSnapshot() *goods.GoodsSnapshot
	}

	// 商品服务
	IGoodsManager interface {
		// 创建商品
		CreateGoodsByItem(product.IProduct, *GoodsItem) IGoods

		// 创建商品
		CreateGoods(*GoodsItem) IGoods

		// 根据产品编号获取商品
		GetGoods(id int32) IGoods

		// 根据产品SKU获取商品
		GetGoodsBySku(itemId, skuId int32) IGoods

		// 删除商品
		DeleteGoods(id int32) error

		//// 获取指定的商品快照
		//GetSaleSnapshot(id int32) *goods.GoodsSnapshot
		//
		//// 根据Key获取商品快照
		//GetSaleSnapshotByKey(key string) *goods.GoodsSnapshot

		// 获取指定数量已上架的商品
		GetOnShelvesGoods(start, end int, sortBy string) []*valueobject.Goods
	}

	// 简单商品信息
	SimpleGoods struct {
		GoodsId    int32  `json:"id"`
		GoodsImage string `json:"img"`
		Name       string `json:"name"`
		Quantity   string `json:"qty"`
	}
)
