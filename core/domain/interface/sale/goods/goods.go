/**
 * Copyright 2015 @ z3q.net.
 * name : goods
 * author : jarryliu
 * date : 2016-06-28 23:54
 * description :
 * history :
 */
package goods

import (
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
)

var (
	ErrNoSuchGoods *domain.DomainError = domain.NewDomainError(
		"no_such_goods", "商品不存在")

	ErrLatestSnapshot *domain.DomainError = domain.NewDomainError(
		"latest_snapshot", "已经是最新的快照")

	ErrNoSuchSnapshot *domain.DomainError = domain.NewDomainError(
		"no_such_snapshot", "商品快照不存在")

	ErrNotOnShelves *domain.DomainError = domain.NewDomainError(
		"not_on_shelves", "商品未上架")
)

type (
	// 商品仓储
	IGoodsRep interface {
		// 获取商品
		GetValueGoods(itemId int, sku int) *ValueGoods

		// 根据SKU-ID获取商品,SKU-ID为商品ID
		//todo: 循环引有,故为interface{}
		GetGoodsBySKuId(skuId int) interface{}

		// 获取商品
		GetValueGoodsById(goodsId int) *ValueGoods

		// 根据SKU获取商品
		GetValueGoodsBySku(itemId, sku int) *ValueGoods

		// 保存商品
		SaveValueGoods(*ValueGoods) (int, error)

		// 获取在货架上的商品
		GetOnShelvesGoods(merchantId int, start, end int,
			sortBy string) []*valueobject.Goods

		// 获取在货架上的商品
		GetPagedOnShelvesGoods(merchantId int, catIds []int, start, end int,
			where, orderBy string) (total int, goods []*valueobject.Goods)

		// 根据编号获取商品
		GetGoodsByIds(ids ...int) ([]*valueobject.Goods, error)

		// 获取会员价
		GetGoodsLevelPrice(goodsId int) []*MemberPrice

		// 保存会员价
		SaveGoodsLevelPrice(*MemberPrice) (int, error)

		// 移除会员价
		RemoveGoodsLevelPrice(id int) error

		// 保存快照
		SaveSnapshot(*Snapshot) (int, error)

		// 获取最新的商品快照
		GetLatestSnapshot(skuId int) *Snapshot

		// 获取指定的商品快照
		GetSaleSnapshot(id int) *GoodsSnapshot

		// 根据Key获取商品快照
		GetSaleSnapshotByKey(key string) *GoodsSnapshot

		// 根据指定商品快照
		GetSnapshots(skuIdArr []int) []Snapshot
	}

	// 商品
	ValueGoods struct {
		Id int `db:"id" pk:"yes" auto:"yes"`

		// 货品编号
		ItemId int `db:"item_id"`

		// 是否为赠品
		IsPresent int `db:"is_present"`

		// 规格
		SkuId int `db:"sku_id"`

		// 促销标志
		PromotionFlag int `db:"prom_flag"`

		// 库存
		StockNum int `db:"stock_num"`

		// 已售件数
		SaleNum int `db:"sale_num"`

		// 销售价
		SalePrice float32 `db:"-"`

		// 促销价
		PromPrice float32 `db:"-"`

		// 实际价
		Price float32 `db:"-"`
	}

	// 会员价
	MemberPrice struct {
		Id      int     `db:"id" pk:"yes" auto:"yes"`
		GoodsId int     `db:"goods_id"`
		Level   int     `db:"level"`
		Price   float32 `db:"price"`
		Enabled int     `db:"enabled"`
	}
)

// 转换为商品值对象
func ParseToValueGoods(v *valueobject.Goods) *ValueGoods {
	return &ValueGoods{
		Id:            v.GoodsId,
		ItemId:        v.Item_Id,
		IsPresent:     v.IsPresent,
		SkuId:         v.SkuId,
		PromotionFlag: v.PromotionFlag,
		StockNum:      v.StockNum,
		SaleNum:       v.SaleNum,
		SalePrice:     v.SalePrice,
		PromPrice:     v.PromPrice,
		Price:         v.Price,
	}
}
