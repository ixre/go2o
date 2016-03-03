/**
 * Copyright 2015 @ z3q.net.
 * name : device
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util

import (
	gutil "github.com/jsix/gof/util"
<<<<<<< HEAD
	"net/http"
=======
	"github.com/jsix/gof/web"
	"net/http"
	"net/url"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"time"
)

const (
	clientDeviceTypeCookieId string = "client_device_type"
<<<<<<< HEAD
=======

>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	// PC设备
	DevicePC string = "1"
	// 手持设备
	DeviceMobile string = "2"
	// 触摸设备
	DeviceTouchPad string = "3"
	// APP内嵌网页
	DeviceAppEmbed string = "4"
)

// 获取浏览设备
<<<<<<< HEAD
func GetBrownerDevice(r *http.Request) string {
	ck, err := r.Cookie(clientDeviceTypeCookieId)
=======
func GetBrownerDevice(ctx *web.Context) string {
	ck, err := ctx.Request.Cookie(clientDeviceTypeCookieId)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	if err == nil && ck != nil {
		switch ck.Value {
		case "1":
			return DevicePC
		case "2":
			return DeviceMobile
		case "3":
			return DeviceTouchPad
		case "4":
			return DeviceAppEmbed
		}
	}
<<<<<<< HEAD
	if gutil.IsMobileAgent(r.UserAgent()) {
=======
	if gutil.IsMobileAgent(ctx.Request.UserAgent()) {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		return DeviceMobile
	}
	return DevicePC
}

// 设置浏览设备
<<<<<<< HEAD
func SetBrownerDevice(w http.ResponseWriter, r *http.Request, deviceType string) {
	ck, err := r.Cookie(clientDeviceTypeCookieId)
=======
func SetBrownerDevice(ctx *web.Context, deviceType string) {
	ck, err := ctx.Request.Cookie(clientDeviceTypeCookieId)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
	if ck != nil {
		ck.HttpOnly = false
		ck.Path = "/"
		http.SetCookie(w, ck)
	}
}

func SetDeviceByUrlQuery(w http.ResponseWriter, r *http.Request) bool {
	dvType := r.URL.Query().Get("device")
	if len(dvType) != 0 {
		SetBrownerDevice(w, r, dvType)
=======

	if ck != nil {
		ck.HttpOnly = false
		ck.Path = "/"
		http.SetCookie(ctx.Response, ck)
	}
}

func SetDeviceByUrlQuery(ctx *web.Context, form *url.Values) bool {
	dvType := form.Get("device")
	if len(dvType) != 0 {
		SetBrownerDevice(ctx, dvType)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		return true
	}
	return false
}
