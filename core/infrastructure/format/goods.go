/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-11 21:15
 * description :
 * history :
 */
package format

import (
	"fmt"
	"github.com/ixre/gof"
	"go2o/core/infrastructure"
	"go2o/core/variable"
	"strconv"
	"strings"
)

var (
	imageServe   string
	noPicUrl     string
	_noPicUrl    string
	picCfgLoaded bool
)

//todo: not used
// 格式化商品编号，不足位用０补齐
func FormatGoodsNo(d int) string {
	const l int = 6
	s := strconv.Itoa(d)
	sl := len(s)
	if sl >= 6 {
		return s
	}
	return strings.Repeat("0", l-sl) + s
}

// 获取无图片地址
func GetNoPicPath() string {
	if len(_noPicUrl) == 0 {
		ctx := infrastructure.GetApp()
		_noPicUrl = ctx.Config().GetString(variable.NoPicPath)
	}
	return _noPicUrl
}

func containProto(s string) bool {
	return strings.HasPrefix(s, "//") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://")
}

func getImageServe() string {
	app := gof.CurrentApp
	cfg := app.Config()
	s := cfg.GetString(variable.ImageServer)
	if containProto(s) {
		return s
	}
	return fmt.Sprintf("//www.%s%s",
		cfg.GetString(variable.ServerDomain), s)
}

// 获取资源前缀
func GetResUrlPrefix() string {
	if len(imageServe) == 0 {
		imageServe = getImageServe()
	}
	return imageServe
}

// 获取商品图片地址
func GetGoodsImageUrl(image string) string {
	if !picCfgLoaded {
		ctx := infrastructure.GetApp()
		if len(imageServe) == 0 {
			imageServe = getImageServe()
		}

		if len(noPicUrl) == 0 {
			noPicUrl = imageServe + "/" + ctx.Config().GetString(variable.NoPicPath)
		}
		picCfgLoaded = true
	}

	if len(image) == 0 {
		return noPicUrl
	}

	if containProto(image) {
		return image
	}
	return imageServe + "/" + image
}

// 获取资源地址
func GetResUrl(image string) string {
	if !picCfgLoaded {
		ctx := infrastructure.GetApp()
		if len(imageServe) == 0 {
			imageServe = getImageServe()
		}

		if len(noPicUrl) == 0 {
			noPicUrl = imageServe + "/" + ctx.Config().GetString(variable.NoPicPath)
		}
		picCfgLoaded = true
	}

	if len(image) == 0 {
		return noPicUrl
	}

	if containProto(image) {
		return image
	}
	return imageServe + "/" + image
}

// 获取URL/路径的名称
func GetName(url string) string {
	if url != "" {
		arr := strings.Split(url, "/")
		if l := len(arr); l > 0 {
			return arr[l-1]
		}
	}
	return ""
}
