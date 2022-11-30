/**
 * Copyright 2015 @ 56x.net.
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
		GetShop(shopId int64) IShop
		// 获取自营店铺
		QuerySelfSupportShops() []Shop
		// 获取门店
		GetStore(storeId int64) IShop

		SaveShop(*Shop) (int64, error)
		// 检查商户商城是否存在(创建)
		CheckShopExists(vendorId int64) bool
		// 获取店铺数量
		ShopCount(vendorId int64, shopType int32) int
		// 商店别名是否存在
		ShopAliasExists(alias string, shopId int) bool
		// 获取商店值
		GetValueShop(shopId int64) *Shop
		// 获取线下商店
		GetOfflineShop(shopId int) *OfflineShop

		// 获取商户的店铺
		GetOnlineShopOfMerchant(vendorId int) *OnlineShop

		// 获取商户所有商店
		GetShopId(mchId int64) []Shop

		// 删除线上商店
		DeleteOnlineShop(mchId, shopId int64) error

		// 删除线下门店
		DeleteOfflineShop(mchId, shopId int64) error

		// 保存线上商店
		SaveOnlineShop(v *OnlineShop) (int64, error)

		// 保存线下商店
		SaveOfflineShop(v *OfflineShop, create bool) error

		// GetShopIdByAlias 根据alias获取店铺编号
		GetShopIdByAlias(alias string) int64
	}
)
