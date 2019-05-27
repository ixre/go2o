package api

import "github.com/ixre/gof/api"


var _ api.Handler = new(AppApi)

type AppApi struct {
}

func (a AppApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"check": a.check,
	})
}

// 检测版本
func (a AppApi) check(ctx api.Context)interface{} {
	form := ctx.Form()
	version := "1.0.1"
	prodVersion := form.GetString("prod_version")
	prodType := form.GetString("prod_type")
	info := "修复已知BUG\n界面更新修复\n主要业务逻辑调整"
	isLatest := api.CompareVersion(prodVersion,version) >= 0
	data := map[string]interface{}{
		"version":version,
		"latest":isLatest,
		"force":true,
		"url":"http://s.to2.net/apk",
		"info":info,
	};
	if prodType == "ios"{
		data["url"] = "http://s.to2.net/ios"
		return data
	}
	if prodType == "android"{
		return data
	}
	return data
}