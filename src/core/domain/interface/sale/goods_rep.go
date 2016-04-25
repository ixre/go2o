/**
 * Copyright 2015 @ z3q.net.
 * name : goods_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import "go2o/src/core/domain/interface/valueobject"

// 商品仓储
type IGoodsRep interface {

	// 获取商品
	GetValueGoods(itemId int, sku int) *ValueGoods

	// 获取商品
	GetValueGoodsById(goodsId int) *ValueGoods

	// 根据SKU获取商品
	GetValueGoodsBySku(itemId, sku int) *ValueGoods

	// 保存商品
	SaveValueGoods(*ValueGoods) (int, error)

	// 获取在货架上的商品
	GetOnShelvesGoods(partnerId int, start, end int,
		sortBy string) []*valueobject.Goods

	// 获取在货架上的商品
	GetPagedOnShelvesGoods(partnerId int, catIds []int, start, end int,
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
