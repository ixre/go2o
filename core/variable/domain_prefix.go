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
	// 主站域名前缀
	DOMAIN_PREFIX_PORTAL = "www."
	// 手机端域名前缀
	DOMAIN_PREFIX_MOBILE = "m."
	// 系统管理域名前缀
	DOMAIN_PREFIX_WEBMASTER = "webmaster."
	// 会员中心域名前缀
	DOMAIN_PREFIX_MEMBER = "u."

	// 商户系统域名前缀
	DOMAIN_PREFIX_MERCHANT = "mch."
	// 通行证域名前缀
	DOMAIN_PREFIX_PASSPORT = "passport."

	// 通行证域名协议,默认为http,可以使用https安全加密
	DOMAIN_PASSPORT_PROTO = "http"
	// 是否启用安全连接
	DOMAIN_PREFIX_SSL = false
	// HTTP API应用前缀
	DOMAIN_PREFIX_HAPI = "hapi."
	// 静态服务器前缀
	DOMAIN_PREFIX_STATIC = "static."
	// 图片服务器前缀
	DOMAIN_PREFIX_IMAGE = "img."

	// 会员中心域名前缀(移动端)
	DOMAIN_PREFIX_M_MEMBER = DOMAIN_PREFIX_MEMBER + DOMAIN_PREFIX_MOBILE
	// 通行证域名前缀(移动端)
	DOMAIN_PREFIX_M_PASSPORT = DOMAIN_PREFIX_PASSPORT + DOMAIN_PREFIX_MOBILE
)
