/**
 * Copyright 2013 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-03 23:18
 * description :
 * history :
 */
package ols

import (
	"bytes"
	"fmt"
	"github.com/atnet/gof"
	gutil "github.com/atnet/gof/util"
	"github.com/atnet/gof/web"
	"go2o/src/app/front/shop"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/service/dps"
	"html/template"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// 处理自定义错误
func HandleCustomError(w http.ResponseWriter, ctx *web.Context, err error) {
	if err != nil {
		ctx.App.Template().Execute(w, gof.TemplateDataMap{
			"error":  err.Error(),
			"statck": template.HTML(strings.Replace(string(debug.Stack()), "\n", "<br />", -1)),
		},
			strings.Replace("views/shop/{device}/error.html", "{device}", ctx.Items["device_view_dir"].(string), -1))
	}
}

func GetShops(c gof.App, partnerId int) []byte {
	//分店
	var buf *bytes.Buffer = bytes.NewBufferString("")

	shops := dps.PartnerService.GetShopsOfPartner(partnerId)
	if len(shops) == 0 {
		return []byte("<div class=\"nodata noshop\">还未添加分店</div>")
	}
	buf.WriteString("<ul class=\"shops\">")
	for i, v := range shops {
		buf.WriteString(fmt.Sprintf(`<li class="s%d">
			<div class="name"><span><strong>%s</strong></div>
			<span class="shop-state shopstate%d">%s</span>
			<div class="phone">%s</div>
			<div class="address">%s</div>
			</li>`, i+1, v.Name, v.State, enum.GetFrontShopStateName(v.State), v.Phone, v.Address))
	}
	buf.WriteString("</ul>")
	return buf.Bytes()
}

func GetCategories(c gof.App, partnerId int, secret string) []byte {
	var buf *bytes.Buffer = bytes.NewBufferString("")
	categories := dps.SaleService.GetCategories(partnerId)

	buf.WriteString(`<ul class="categories">
		<li class="s0 current" val="0">
			<div class="name"><span><strong>全部</strong></div>
		</li>
	`)
	for i, v := range categories {
		buf.WriteString(fmt.Sprintf(`<li class="s%d" val="%d">
			<div class="name"><span><strong>%s</strong></div>
			</li>`, i+1, v.Id, v.Name))
	}
	buf.WriteString("</ul>")
	return buf.Bytes()
}

const (
	clientDeviceTypeCookieId string = "client_device_type"
)

// 获取浏览设备
func GetBrownerDevice(ctx *web.Context) string {
	ck, err := ctx.Request.Cookie(clientDeviceTypeCookieId)
	if err == nil && ck != nil {
		switch ck.Value {
		case "1":
			return shop.DevicePC
		case "2":
			return shop.DeviceMobile
		case "3":
			return shop.DeviceTouchPad
		}
	}
	if gutil.IsMobileAgent(ctx.Request.UserAgent()) {
		return shop.DeviceMobile
	}
	return shop.DevicePC
}

// 设置浏览设备
func SetBrownerDevice(ctx *web.Context, deviceType string) {
	ck, err := ctx.Request.Cookie(clientDeviceTypeCookieId)
	isDefaultDevice := deviceType == "" || deviceType == "1"
	if err == nil && ck != nil {
		if isDefaultDevice {
			ck.Value = deviceType
			ck.Expires = time.Now().Add(-time.Hour * 48)
		} else {
			ck.Value = deviceType
			ck.Expires = time.Now().Add(time.Hour * 24)
		}
	} else if !isDefaultDevice {
		ck = &http.Cookie{
			Name:    clientDeviceTypeCookieId,
			Value:   deviceType,
			Expires: time.Now().Add(time.Hour * 24),
		}
	}

	if ck != nil {
		ck.HttpOnly = false
		ck.Path = "/"
		http.SetCookie(ctx.ResponseWriter, ck)
	}
}
