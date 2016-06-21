/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-11-22 20:01
 * description :
 * history :
 */

package shop

import "fmt"

type (
	IShopManager interface {
		// 新建商店
		CreateShop(*Shop) IShop

		// 获取所有商店
		GetShops() []IShop

		// 获取营业中的商店
		GetBusinessInShops() []IShop

		// 获取商店
		GetShop(int) IShop

		// 根据名称获取商店
		GetShopByName(name string) IShop

		// 删除门店
		DeleteShop(shopId int) error

		// 重新加载数据
		Reload()
	}
)

//位置(经度+"/"+纬度)
func (this OfflineShop) Location() string {
	return fmt.Sprintf("%f/%f", this.Lng, this.Lat)
}
