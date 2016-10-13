/**
 * Copyright 2015 @ at3.net.
 * name : app.go
 * author : jarryliu
 * date : 2016-10-14 00:21
 * description :
 * history :
 */
package app

import "github.com/jsix/gof"

const (
	// 网页服务
	TagWebServe = 1 << iota
	// 常驻程序
	TagDaemon
	// 系统后台
	TagMasterServe
	// Tcp服务
	TagTcpServe
	// 商户系统
	TagMchServe
	// 用户中心
	TagUCenterServe
	// 通行证
	TagPassportServe
	// 商铺系统
	TagShopServe
)

var (
	TagWeb = TagWebServe | TagMchServe | TagUCenterServe |
		TagPassportServe | TagShopServe
)

type CustomConfig func(gof.App, int) error

// 自定义配置应用系统
func Configure(c CustomConfig, app gof.App, tag int) error {
	return c(app, tag)
}
