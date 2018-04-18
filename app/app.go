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
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/log"
	"github.com/jsix/gof/shell"
	"go2o/core/variable"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
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

// 自定义配置应用系统
func Configure(c CustomConfig, app gof.App, tag int) error {
	err := c(app, tag)
	if tag&FlagWebServe == FlagWebServe {
		defer flushJsGlob()
	}
	return err
}

// 将变量输出到JS中
func flushJsGlob() {
	filePath := "static/assets/js/base.js"
	fi, err := os.Open(filePath)
	if err == nil {
		defer fi.Close()
		data, err := ioutil.ReadAll(fi)
		if err == nil {
			newBytes := []byte(fmt.Sprintf("var domain='%s';var hapi='//%s'+domain;",
				variable.Domain, variable.DOMAIN_PREFIX_HApi,
			))
			txt := string(data)
			delimer := "/*~*/"
			i := strings.Index(txt, delimer)
			if i == -1 {
				newBytes = append(newBytes, delimer...)
				newBytes = append(newBytes, txt...)
			} else {
				newBytes = append(newBytes, txt[i:]...)
			}
			ioutil.WriteFile(filePath, newBytes, os.ModePerm)
		}
	} else {
		log.Println("[ Flush][ JS][ Error]:", err)
	}
}

// 初始化，如果不调用此操作。则默认全部不监视。
func FsInit(debug bool) {
	for i := 0; i < fsLast; i++ {
		webFs[i] = debug
	}
	// MAC OX 连接数限制,在这里指定部分应用更新模板
	if debug && runtime.GOOS == "darwin" {
		resetFsOnDarwin()
	}
}

// 获取模板是否监视更改
func GetFs(i int) bool {
	return webFs[i]
}

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

// 重设MAC OX下的文件监视更改
func resetFsOnDarwin() {
	webFs[FsPortal] = false
	webFs[FsPortalMobile] = !false
	webFs[FsPassport] = false
	webFs[FsPassportMobile] = false
	webFs[FsUCenter] = false
	webFs[FsUCenterMobile] = false
	webFs[FsShop] = false
	webFs[FsShopMobile] = false
	webFs[FsMch] = false
	webFs[FsWholesale] = false
}
