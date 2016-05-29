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
		// 获取站点配置
		GetSiteConf() ShopSiteConf

		// 保存站点配置
		SaveSiteConf(*ShopSiteConf) error

		// 新建商店
		CreateShop(*Shop) IShop

		// 获取所有商店
		GetShops() []IShop

		// 获取营业中的商店
		GetBusinessInShops() []IShop

		// 获取商店
		GetShop(int) IShop

		// 删除门店
		DeleteShop(shopId int) error
	}

	//合作商店网站配置
	ShopSiteConf struct {
		//合作商编号
		MerchantId int `db:"merchant_id" auto:"no" pk:"yes"`

		//主机
		Host string `db:"host"`

		//前台Logo
		Logo string `db:"logo"`

		//首页标题
		IndexTitle string `db:"index_title"`

		//子页面标题
		SubTitle string `db:"sub_title"`

		//状态: 0:暂停  1：正常
		State     int    `db:"state"`
		StateHtml string `db:"state_html"`
	}
)

//位置(经度+"/"+纬度)
func (this OfflineShop) Location() string {
	return fmt.Sprintf("%f/%f", this.Lng, this.Lat)
}
