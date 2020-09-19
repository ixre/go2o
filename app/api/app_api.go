package api

import (
	"context"
	"github.com/ixre/gof/api"
	"go2o/core/service"
	"go2o/core/service/proto"
	"time"
)

const appVersion = "app_version"
const appAndroidVersion = "app_android_version"
const appIOSVersion = "app_ios_version"
const appReleaseInfo = "app_release_info"
const appApkFileUrl = "app_apk_file_url"
const appIOSFileUrl = "app_ios_file_url"

var _ api.Handler = new(AppApi)

type AppApi struct {
}

func NewAppApi() *AppApi {
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
	trans, cli, err := service.RegistryServeClient()
	if err == nil {
		defer trans.Close()
		keys := []string{appVersion, appAndroidVersion, appIOSVersion,
			appReleaseInfo, appApkFileUrl, appIOSFileUrl}
		mp, _ := cli.GetRegistries(context.TODO(), &proto.StringArray{Value: keys})
		version := ""
		url := ""
		if prodType == "android" {
			version = mp.Value[appAndroidVersion]
			url = mp.Value[appApkFileUrl]
		} else if prodType == "ios" {
			version = mp.Value[appIOSVersion]
			url = mp.Value[appIOSFileUrl]
		} else {
			version = mp.Value[appVersion]
			url = ""
		}
		info := mp.Value[appReleaseInfo]
		isLatest := api.CompareVersion(prodVersion, version) >= 0
		data := map[string]interface{}{
			"version": version,
			"latest":  isLatest,
			"force":   true,
			"url":     url,
			"info":    info,
		}
		return data
	}
	return api.ResponseWithCode(-1, "无法检测版本信息")
}

func (a *AppApi) init() *AppApi {
	time.Sleep(time.Second * 5) // 等待RPC服务启动5秒
	trans, cli, err := service.RegistryServeClient()
	if err == nil {
		defer trans.Close()
		cli.CreateRegistry(context.TODO(),
			&proto.RegistryCreateRequest{
				Key:          appVersion,
				DefaultValue: "1.0.0",
				Description:  "APP版本号",
			})
		cli.CreateRegistry(context.TODO(),
			&proto.RegistryCreateRequest{
				Key:          appAndroidVersion,
				DefaultValue: "1.0.0",
				Description:  "安卓APP版本号",
			})
		cli.CreateRegistry(context.TODO(),
			&proto.RegistryCreateRequest{
				Key:          appIOSVersion,
				DefaultValue: "1.0.0",
				Description:  "苹果APP版本号",
			})
		cli.CreateRegistry(context.TODO(),
			&proto.RegistryCreateRequest{
				Key:          appReleaseInfo,
				DefaultValue: "修复已知BUG\n界面调整",
				Description:  "版本发布日志",
			})
		cli.CreateRegistry(context.TODO(),
			&proto.RegistryCreateRequest{
				Key:          appApkFileUrl,
				DefaultValue: "",
				Description:  "安卓APK文件下载地址",
			})
		cli.CreateRegistry(context.TODO(),
			&proto.RegistryCreateRequest{
				Key:          appIOSFileUrl,
				DefaultValue: "",
				Description:  "苹果APP文件下载地址",
			})
	} else {
		println("init app api err:", err.Error())
	}
	return a
}
