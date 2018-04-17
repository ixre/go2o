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
	"go2o/core/variable"
	"net/http"
	"strings"
	"time"
)

const (
	clientDeviceTypeCookieId = "client_device_type"
	// PC设备
	DevicePC = "1"
	// 手持设备
	DeviceMobile = "2"
	// 触摸设备
	DeviceTouchPad = "3"
	// APP内嵌网页
	DeviceAppEmbed = "4"
)

// 获取浏览设备
func GetBrownerDevice(r *http.Request) string {
	//return DeviceMobile
	//return getDevice(r)
	ck, err := r.Cookie(clientDeviceTypeCookieId)
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
	if gutil.IsMobileAgent(r.UserAgent()) {
		return DeviceMobile
	}
	return DevicePC
}
func getDevice(r *http.Request) string {
	if gutil.IsMobileAgent(r.UserAgent()) {
		return DeviceMobile
	}
	return DevicePC
}

// 设置浏览设备
func SetBrownerDevice(w http.ResponseWriter, r *http.Request, deviceType string) {
	ck, err := r.Cookie(clientDeviceTypeCookieId)
	isDefaultDevice := deviceType == "" || deviceType == "1"
	domain := getCookieDomain(r.Host)
	if domain == "" {
		domain = variable.Domain
	}
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
		ck.Domain = domain
		http.SetCookie(w, ck)
	}
}

// 根据host获取cookie域名
func getCookieDomain(host string) string {
	i := strings.LastIndex(host, ":")
	if i > 0 {
		host = host[:i]
	}
	i = strings.LastIndex(host, ".")
	if i > 0 {
		j := strings.LastIndex(host[:i], ".")
		if j >= 0 {
			return host[j+1:]
		}
	}
	return host
}

func SetDeviceByUrlQuery(w http.ResponseWriter, r *http.Request) bool {
	dvType := r.URL.Query().Get("device")
	if len(dvType) != 0 {
		SetBrownerDevice(w, r, dvType)
		return true
	}
	return false
}

// 是否为移动设备请求
func MobileRequest(r *http.Request) bool {
	//判断是否来自移动请求
	referrer := r.Referer()
	//todo: 判断是否为手机
	i := strings.Index(referrer, variable.DOMAIN_PREFIX_PORTAL_MOBILE)
	if i != -1 {
		preChar := referrer[i-1]
		if preChar == '.' || preChar == '/' {
			return true
		}
	}
	return false
}
