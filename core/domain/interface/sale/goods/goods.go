/**
 * Copyright 2015 @ z3q.net.
 * name : goods
 * author : jarryliu
 * date : 2016-06-28 23:54
 * description :
 * history :
 */
package goods

import "go2o/core/domain/interface/valueobject"

type (
	// 商品仓储
	IGoodsRep interface {
		// 获取商品
		GetValueGoods(itemId int, sku int) *ValueGoods

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
