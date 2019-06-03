/**
 * Copyright 2015 @ to2.net.
 * name : shop_repo.go
 * author : jarryliu
 * date : 2016-05-28 13:10
 * description :
 * history :
 */
package shop

type (
	IShopRepo interface {
		// 获取商店
		GetShop(shopId int32) IShop
		SaveShop(*Shop) (int32, error)
		// 获取店铺数量
		ShopCount(vendorId int32, shopType int32) int
		// 商店别名是否存在
		ShopAliasExists(alias string, shopId int32) bool
		// 获取商店值
		GetValueShop(shopId int32) *Shop
		// 获取线上商店
		GetOnlineShop(shopId int32) *OnlineShop
		// 获取线下商店
		GetOfflineShop(shopId int32) *OfflineShop

		// 获取商户所有商店
		GetShopsOfMerchant(mchId int32) []Shop

		// 删除线上商店
		DeleteOnlineShop(mchId, shopId int32) error

		// 删除线下门店
		DeleteOfflineShop(mchId, shopId int32) error

		// 保存线上商店
		SaveOnlineShop(v *OnlineShop, create bool) error

		// 保存线下商店
		SaveOfflineShop(v *OfflineShop, create bool) error
	}
)
