/**
 * Copyright 2014 @ ops.
 * name :
 * author : jarryliu
 * date : 2013-11-10 09:44
 * description :
 * history :
 */

package variable

var (
	// 零售门户前缀
	DOMAIN_PREFIX_PORTAL = "www."
	// 批发门户域名前缀
	DOMAIN_PREFIX_WHOLESALE_PORTAL = "whs."
	// 零售门户手机端域名前缀
	DOMAIN_PREFIX_PORTAL_MOBILE = "m."
	// 会员中心域名前缀
	DOMAIN_PREFIX_MEMBER = "u."
	// 商户系统域名前缀
	DOMAIN_PREFIX_MERCHANT = "mch."
	// 通行证域名前缀
	DOMAIN_PREFIX_PASSPORT = "passport."

	// 通行证域名协议,默认为http,可以使用https安全加密
	DOMAIN_PASSPORT_PROTO = "http"
	// API系统
	DOMAIN_PREFIX_HApi = "hapi."
	// 静态服务器前缀
	DOMAIN_PREFIX_STATIC = "static."
	// 图片服务器前缀
	DOMAIN_PREFIX_IMAGE = "img."
	// 批发中心移动端
	DOMAIN_PREFIX_M_WHOLESALE = "mwhs."
	// 会员中心域名前缀(移动端)
	DOMAIN_PREFIX_M_MEMBER = "mu."
	// 通行证域名前缀(移动端)
	DOMAIN_PREFIX_M_PASSPORT = "mpp."
)

const (
	DEnabledSSL            = "D_EnabledSSL"
	DStaticPathr          = "D_StaticPathr"
	DImageServer           = "D_ImageServer"
	DUrlHash               = "D_Hash"
	DRetailPortal          = "D_RetailPortal"
	DWholesalePortal       = "D_WholesalePortal"
	DUCenter               = "D_UCenter"
	DPassport              = "D_Passport"
	DMerchant              = "D_Merchant"
	DHApi                  = "D_HApi"
	DRetailMobilePortal    = "D_RetailMobilePortal"
	DWholesaleMobilePortal = "D_WholesaleMobilePortal"
	DMobilePassport        = "D_MobilePassport"
	DMobileUCenter         = "D_MobileUCenter"
)
