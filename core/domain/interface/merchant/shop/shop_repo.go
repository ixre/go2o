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
		// 创建电子商城
		CreateShop(shop *OnlineShop) IShop
		// 获取商店
		GetShop(shopId int) IShop

		SaveShop(*Shop) (int32, error)
		// 检查商户商城是否存在(创建)
		CheckShopExists(vendorId int) bool
		// 获取店铺数量
		ShopCount(vendorId int32, shopType int32) int
		// 商店别名是否存在
		ShopAliasExists(alias string, shopId int) bool
		// 获取商店值
		GetValueShop(shopId int) *Shop
		// 获取线上商店
		GetOnlineShop(shopId int) *OnlineShop
		// 获取线下商店
		GetOfflineShop(shopId int) *OfflineShop

		// 获取商户的店铺
		GetOnlineShopOfMerchant(vendorId int) *OnlineShop

		// 获取商户所有商店
		GetShopsOfMerchant(mchId int32) []Shop

		// 删除线上商店
		DeleteOnlineShop(mchId, shopId int) error

		// 删除线下门店
		DeleteOfflineShop(mchId, shopId int) error

		// 保存线上商店
		SaveOnlineShop(v *OnlineShop) (int, error)

		// 保存线下商店
		SaveOfflineShop(v *OfflineShop, create bool) error
	}
)
