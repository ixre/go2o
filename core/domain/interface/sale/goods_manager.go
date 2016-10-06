/**
 * Copyright 2015 @ z3q.net.
 * name : sale_goods.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import (
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/valueobject"
)

type (
	// 商品
	IGoods interface {
		// 获取领域对象编号
		GetDomainId() int

		// 商品快照
		SnapshotManager() goods.ISnapshotManager

		// 获取货品
		GetItem() IItem

		// 设置值
		GetValue() *goods.ValueGoods

		// 获取包装过的商品信息
		GetPackedValue() *valueobject.Goods

		// 获取促销信息
		GetPromotions() []promotion.IPromotion

		// 获取促销价
		GetPromotionPrice(level int) float32

		// 获取会员价销价,返回是否有会原价及价格
		GetLevelPrice(level int) (bool, float32)

		// 获取促销描述
		GetPromotionDescribe() map[string]string

		// 获取会员价
		GetLevelPrices() []*goods.MemberPrice

		// 保存会员价
		SaveLevelPrice(*goods.MemberPrice) (int, error)

		// 设置值
		SetValue(*goods.ValueGoods) error

		// 保存
		Save() (int, error)

		// 更新销售数量,扣减库存
		AddSaleNum(quantity int) error

		// 取消销售
		CancelSale(quantity int, orderNo string) error

		//// 生成快照
		//GenerateSnapshot() (int, error)
		//
		//// 获取最新的快照
		//GetLatestSnapshot() *goods.GoodsSnapshot
	}

	// 商品服务
	IGoodsManager interface {
		// 创建商品
		CreateGoodsByItem(IItem, *goods.ValueGoods) IGoods

		// 创建商品
		CreateGoods(*goods.ValueGoods) IGoods

		// 根据产品编号获取商品
		GetGoods(int) IGoods

		// 根据产品SKU获取商品
		GetGoodsBySku(itemId, sku int) IGoods

		// 删除商品
		DeleteGoods(int) error

		//// 获取指定的商品快照
		//GetSaleSnapshot(id int) *goods.GoodsSnapshot
		//
		//// 根据Key获取商品快照
		//GetSaleSnapshotByKey(key string) *goods.GoodsSnapshot

		// 获取指定数量已上架的商品
		GetOnShelvesGoods(start, end int, sortBy string) []*valueobject.Goods
	}

	// 简单商品信息
	SimpleGoods struct {
		GoodsId    int    `json:"id"`
		GoodsImage string `json:"img"`
		Name       string `json:"name"`
		Quantity   string `json:"qty"`
	}
)
