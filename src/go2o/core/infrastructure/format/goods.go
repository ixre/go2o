/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-11 21:15
 * description :
 * history :
 */
package format

import (
	"go2o/core/infrastructure"
	"go2o/share/variable"
	"strconv"
	"strings"
)

var (
	imageServe string
	noPicUrl   string
)

func init() {
	ctx := infrastructure.GetContext()
	if len(imageServe) == 0 {
		imageServe = ctx.Config().GetString(variable.ImageServer)
	}
	if len(noPicUrl) == 0 {
		noPicUrl = ctx.Config().GetString(variable.StaticServer) +
			"/" + ctx.Config().GetString(variable.NoPicPath)
	}
}

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

// 获取商品图片地址
func GetGoodsImageUrl(image string) string {
	if len(image) == 0 {
		return noPicUrl
	}
	return imageServe + "/" + image
}
