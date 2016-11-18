/**
 * Copyright 2015 @ z3q.net.
 * name : shop_rep.go
 * author : jarryliu
 * date : 2016-05-28 13:10
 * description :
 * history :
 */
package shop

type (
	IShopRep interface {
		SaveShop(*Shop) (int32, error)

		// 商店别名是否存在
		ShopAliasExists(alias string, shopId int32) bool

		// 获取商店值
		GetValueShop(mchId, shopId int32) *Shop

		// 获取商户所有商店
		GetShopsOfMerchant(mchId int32) []Shop

		// 删除线上商店
		DeleteOnlineShop(mchId, shopId int32) error

		// 删除线下门店
		DeleteOfflineShop(mchId, shopId int32) error

		// 获取线上商店
		GetOnlineShop(shopId int32) *OnlineShop

		// 保存线上商店
		SaveOnlineShop(v *OnlineShop, create bool) error

		// 获取线下商店
		GetOfflineShop(shopId int32) *OfflineShop

		// 保存线下商店
		SaveOfflineShop(v *OfflineShop, create bool) error
	}
)
