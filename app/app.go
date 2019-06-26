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
	"github.com/ixre/gof"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/shell"
	"time"
)

const (
	// 网页服务
	FlagWebServe = 1 << iota
	// 常驻程序
	FlagDaemon
	// 远程调用服务
	FlagRpcServe
	// Tcp服务
	FlagTcpServe
	// 单元测试
	FlagTesting
	// 二进制程序
	FlagBin
	// 系统后台
	FlagMasterServe
	// 商户系统
	FlagMchServe
	// 用户中心
	FlagUCenterServe
	// 通行证
	FlagPassportServe
	// 商铺系统
	FlagShopServe
)

const (
	FsPortal = iota
	FsPortalMobile
	FsPassport
	FsPassportMobile
	FsUCenter
	FsUCenterMobile
	FsShop
	FsShopMobile
	FsMch
	FsWholesale
	fsLast
)

var (
	FlagWebApp = FlagWebServe | FlagMchServe | FlagUCenterServe |
		FlagPassportServe | FlagShopServe

	// 网站模板文件监视
	webFs = make(map[int]bool)
)

type CustomConfig func(gof.App, int) error



// 自动安装包
func AutoInstall() {
	execInstall()
	d := time.Second * 15
	t := time.NewTimer(d)
	for {
		select {
		case <-t.C:
			if err := execInstall(); err == nil {
				t.Reset(d)
			} else {
				break
			}
		}
	}
}

func execInstall() error {
	_, _, err := shell.Run("go install .")
	if err != nil {
		log.Println("[ Go2o][ Install]:", err)
	}
	return err
}
