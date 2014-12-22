/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package partner

type IPartner interface {
	GetAggregateRootId() int

	GetValue() ValuePartner

	// 获取销售配置
	GetSaleConf() SaleConf

	// 保存销售配置
	SaveSaleConf(*SaleConf)error

	// 获取站点配置
	GetSiteConf()SiteConf

	// 保存站点配置
	SaveSiteConf(*SiteConf)error
}
