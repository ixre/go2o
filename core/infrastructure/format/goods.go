/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2014-02-11 21:15
 * description :
 * history :
 */
package format

import (
	"strings"
)

var (
	imageServe   string
	noPicUrl     string
	_noPicUrl    string
	picCfgLoaded bool
)

// 静态文件服务器路径
var GlobalFileServerPath = ""

// 获取无图片地址
func GetNoPicPath() string {
	if len(_noPicUrl) == 0 {
		_noPicUrl = "res/nopic.gif"
	}
	return _noPicUrl
}

func containProto(s string) bool {
	return strings.HasPrefix(s, "//") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://")
}

// GetFileFullUrl 获取完整的文件地址
func GetFileFullUrl(url string) string {
	if len(url) == 0 {
		return ""
	}
	if containProto(url) {
		return url
	}
	return GlobalFileServerPath + url
}

// 获取商品图片地址
func GetGoodsImageUrl(image string) string {
	return GetFileFullUrl(image)
	if !picCfgLoaded {
		if len(imageServe) == 0 {
			imageServe = GlobalFileServerPath
		}

		if len(noPicUrl) == 0 {
			noPicUrl = imageServe + "/res/nopic.gif"
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
