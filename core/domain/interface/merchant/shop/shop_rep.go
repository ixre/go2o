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

		// 获取站点配置
		GetSiteConf(int) *ShopSiteConf

		SaveSiteConf(merchantId int, v *ShopSiteConf) error

		SaveShop(*Shop) (int, error)

		GetValueShop(merchantId, shopId int) *Shop

		GetShopsOfMerchant(merchantId int) []*Shop

		DeleteShop(merchantId, shopId int) error
	}
)
