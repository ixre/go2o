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
		SaveShop(*Shop) (int, error)

		// 商店别名是否存在
		ShopAliasExists(alias string, shopId int) bool

		GetValueShop(merchantId, shopId int) *Shop

		GetShopsOfMerchant(merchantId int) []*Shop

		// 删除线上商店
		DeleteOnlineShop(mchId, shopId int) error

		// 删除线下门店
		DeleteOfflineShop(mchId, shopId int) error

		// 获取线上商店
		GetOnlineShop(shopId int) *OnlineShop

		// 保存线上商店
		SaveOnlineShop(v *OnlineShop, create bool) error

		// 获取线下商店
		GetOfflineShop(shopId int) *OfflineShop

		// 保存线下商店
		SaveOfflineShop(v *OfflineShop, create bool) error
	}
)
