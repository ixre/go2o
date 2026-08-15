/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-11-22 20:01
 * description :
 * history :
 */

package shop

type (
	IShopManager interface {
		// CreateOnlineShop 创建线上店铺
		CreateOnlineShop(o *OnlineShop) (IShop, error)
		//todo: will be removed
		// 新建商店
		CreateShop(*Shop) IShop
		// 获取店铺
		GetOnlineShop() IShop
		// 获取门店
		GetStore(id int) IShop

		// 获取所有商店
		GetShops() []IShop

		// 获取营业中的商店
		GetBusinessInShops() []IShop

		// 根据名称获取商店
		GetShopByName(name string) IShop

		// 删除门店
		DeleteShop(shopId int32) error
	}
)
