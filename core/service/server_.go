package service

import (
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/gof"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : server_.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-10-09 17:33
 * description :
 * history :
 */

// RPC服务初始化
func prepareRpcServer(ctx gof.App) {
	gf := ctx.Config().GetString

	//ssl := gf("ssl_enabled")
	//prefix := "http://"
	//if ssl == "true" || ssl == "1" {
	//	prefix = "https://"
	//}
	repo := inject.GetRegistryRepo()
	update := repo.UpdateValue
	update(registry.ApiRequireVersion, gf("api_require_version"))

	// 更新静态服务器的地址(解偶合)
	//prefix := repo.Get(registry.DomainPrefixImage)
	//format.GlobalImageServer = fmt.Sprintf("%s://%s%s", protocol, prefix, domain)

	//hash := crypto.Md5([]byte(strconv.Itoa(int(time.Now().Unix()))))[8:14]

	//mp[variable.DEnabledSSL] = gf("ssl_enabled")
	//mp[consts.DStaticPath] = gf("static_server")
	//mp[variable.DImageServer] = gf("image_server")
	//mp[variable.DUrlHash] = hash
	//mp[variable.DRetailPortal] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_PORTAL, domain}, "")
	//mp[variable.DWholesalePortal] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_WHOLESALE_PORTAL, domain}, "")
	//mp[variable.DUCenter] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_MEMBER, domain}, "")
	//mp[variable.DPassport] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_PASSPORT, domain}, "")
	//mp[variable.DMerchant] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_MERCHANT, domain}, "")
	//mp[variable.DHApi] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_HApi, domain}, "")
	//
	//mp[variable.DRetailMobilePortal] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_PORTAL_MOBILE, domain}, "")
	//mp[variable.DWholesaleMobilePortal] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_M_WHOLESALE, domain}, "")
	//mp[variable.DMobilePassport] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_M_PASSPORT, domain}, "")
	//mp[variable.DMobileUCenter] = strings.Join([]string{prefix,
	//	variable.DOMAIN_PREFIX_M_MEMBER, domain}, "")

}
