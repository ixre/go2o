/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-12 16:52
 * description :
 * history :
 */

package partner

type IPartnerRep interface {

	CreatePartner(*ValuePartner) IPartner

	GetPartner(int) IPartner

	// 获取销售配置
	GetSaleConf(int) *SaleConf

	SaveSaleConf(*SaleConf) error

	// 获取站点配置
	GetSiteConf(int) *SiteConf

	SaveSiteConf(*SiteConf) error
}
