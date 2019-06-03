/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-11-22 20:01
 * description :
 * history :
 */

package shop

type (
	IShopManager interface {
		// 新建商店
		CreateShop(*Shop) IShop

		// 获取所有商店
		GetShops() []IShop

		// 获取营业中的商店
		GetBusinessInShops() []IShop

		// 获取商铺
		GetOnlineShop() IShop

		// 获取商店
		GetShop(id int32) IShop

		// 根据名称获取商店
		GetShopByName(name string) IShop

		// 删除门店
		DeleteShop(shopId int32) error
	}
)
