/**
 * Copyright 2015 @ at3.net.
 * name : app.go
 * author : jarryliu
 * date : 2016-10-14 00:21
 * description :
 * history :
 */
package app

import (
	"github.com/jsix/gof"
	"runtime"
)

const (
	// 网页服务
	TagWebServe = 1 << iota
	// 常驻程序
	TagDaemon
	// 单元测试
	TagTesting
	// 二进制程序
	TagBin
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

const (
	FsMain = iota
	FsMainMobile
	FsPassport
	FsPassportMobile
	FsUCenter
	FsUCenterMobile
	FsShop
	FsShopMobile
	FsMch
	fsLast
)

var (
	TagWeb = TagWebServe | TagMchServe | TagUCenterServe |
		TagPassportServe | TagShopServe

	// 网站模板文件监视
	webFs map[int]bool = make(map[int]bool)
)

type CustomConfig func(gof.App, int) error

// 自定义配置应用系统
func Configure(c CustomConfig, app gof.App, tag int) error {
	return c(app, tag)
}

// 初始化，如果不调用此操作。则默认全部不监视。
func FsInit(debug bool) {
	for i := 0; i < fsLast; i++ {
		webFs[i] = debug
	}
	// MAC OX 连接数限制,在这里指定部分应用更新模板
	if debug && runtime.GOOS == "darwin" {
		resetFsOnDarwin(webFs)
	}
}

// 重设MAC OX下的文件监视更改
func resetFsOnDarwin(fs map[int]bool) {
	fs[FsMain] = false
	fs[FsMainMobile] = false
	fs[FsPassport] = !false
	fs[FsPassportMobile] = !false
	fs[FsUCenter] = !false
	fs[FsUCenterMobile] = !false
	fs[FsShop] = false
	fs[FsShopMobile] = false
	fs[FsMch] = false
}

// 获取模板是否监视更改
func GetFs(i int) bool {
	return webFs[i]
}
