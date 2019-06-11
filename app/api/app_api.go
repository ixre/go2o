package api

import (
	"github.com/ixre/gof/api"
	"go2o/core/service/thrift"
	"time"
)

const appVersion = "app_version"
const appReleaseInfo = "app_release_info"
const appApkFileUrl = "app_apk_file_url"
const appIosFileUrl = "app_ios_file_url"

var _ api.Handler = new(AppApi)

type AppApi struct {
}

func NewAppApi()*AppApi {
	r := &AppApi{}
	go r.init()
	return r
}

func (a AppApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"check": a.check,
	})
}

// 检测版本
func (a AppApi) check(ctx api.Context) interface{} {
	form := ctx.Form()
	prodVersion := form.GetString("prod_version")
	prodType := form.GetString("prod_type")
	trans,cli,err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		keys := []string{appVersion, appReleaseInfo, appApkFileUrl, appIosFileUrl}
		mp, _ := cli.GetRegistries(thrift.Context, keys)
		info := mp[keys[1]]
		version := mp[keys[0]]
		isLatest := api.CompareVersion(prodVersion, version) >= 0
		data := map[string]interface{}{
			"version": version,
			"latest":  isLatest,
			"force":   true,
			"url":     "",
			"info":    info,
		}
		if prodType == "ios" {
			data["url"] = mp[keys[3]]
			return data
		}
		if prodType == "android" {
			data["url"] = mp[keys[2]]
			return data
		}
	}
	return api.ResponseWithCode(-1,"无法检测版本信息")
}

func (a *AppApi) init()*AppApi {
	time.Sleep(time.Second * 5) // 等待RPC服务启动5秒
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		cli.CreateUserRegistry(thrift.Context, appVersion, "1.0.0", "APP版本号")
		cli.CreateUserRegistry(thrift.Context, appReleaseInfo, "修复已知BUG\n界面调整", "版本发布日志")
		cli.CreateUserRegistry(thrift.Context, appApkFileUrl, "", "安卓APK文件下载地址")
		cli.CreateUserRegistry(thrift.Context, appIosFileUrl, "", "苹果APP文件下载地址")
	} else {
		println("init app api err:", err.Error())
	}
	return a
}
