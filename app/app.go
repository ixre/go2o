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
	"github.com/jsix/gof/util"
	"go2o/core/variable"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
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
	err := c(app, tag)
	if tag&TagWebServe == TagWebServe {
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
			newBytes := []byte(fmt.Sprintf("var domain='%s';var hapi='%s://%s'+domain;",
				variable.Domain,
				util.BoolExt.TString(variable.DOMAIN_PREFIX_SSL, "https", "http"),
				variable.DOMAIN_PREFIX_HAPI,
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

// 重设MAC OX下的文件监视更改
func resetFsOnDarwin() {
	webFs[FsMain] = false
	webFs[FsMainMobile] = !false
	webFs[FsPassport] = false
	webFs[FsPassportMobile] = false
	webFs[FsUCenter] = false
	webFs[FsUCenterMobile] = false
	webFs[FsShop] = false
	webFs[FsShopMobile] = false
	webFs[FsMch] = !false
}

// 获取模板是否监视更改
func GetFs(i int) bool {
	return webFs[i]
}
