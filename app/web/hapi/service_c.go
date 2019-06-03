/**
 * Copyright 2015 @ to2.net.
 * name : service_c.go
 * author : jarryliu
 * date : 2016-09-09 23:45
 * description :
 * history :
 */
package hapi

import (
	"fmt"
	"github.com/ixre/goex/echox"
	"github.com/ixre/gof"
	"github.com/ixre/gof/util"
	ut "go2o/app/util"
	"go2o/core/infrastructure/gen"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"go2o/core/variable"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type serviceC struct {
	gof.App
}

// 切换设备
func (m *serviceC) Device(c *echox.Context) error {
	device := c.QueryParam("device")
	app := c.QueryParam("app")
	if device != "" {
		ut.SetBrownerDevice(c.Response(), c.Request(), device)
	}
	if app != "" {
		//todo::
	}
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), "ok")
}

// 登陆状态
func (s *serviceC) LoginState(c *echox.Context) error {
	mp := make(map[string]interface{})
	mobileReq := ut.MobileRequest(c.Request())
	mPrefix := util.BoolExt.TString(mobileReq, variable.DOMAIN_PREFIX_M_PASSPORT,
		variable.DOMAIN_PREFIX_PASSPORT)
	pstUrl := fmt.Sprintf("//%s%s", mPrefix, variable.Domain)
	memberId := getMemberId(c)
	if memberId <= 0 {
		registry, _ := rsi.FoundationService.GetRegistries(thrift.Context,
			[]string{"PlatformName"})
		mp["PtName"] = registry["PlatformName"]
		mp["LoginUrl"] = pstUrl + "/auth/login"
		mp["RegisterUrl"] = pstUrl + "/register"
		mp["Login"] = 0
	} else {
		mmUrl := fmt.Sprintf("//%s%s",
			variable.DOMAIN_PREFIX_MEMBER, variable.Domain)
		m, _ := rsi.MemberService.GetProfile(thrift.Context, int64(memberId))
		mp["MMName"] = m.Name
		mp["LogoutUrl"] = pstUrl + "/auth/logout"
		mp["MMUrl"] = mmUrl
		mp["Login"] = 1
	}
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), mp)
}

// 收藏
func (m *serviceC) Favorite(c *echox.Context) error {
	memberId := getMemberId(c)
	if memberId <= 0 {
		return requestLogin(c)
	}
	result := &gof.Result{}

	favType := c.QueryParam("type")
	id, _ := util.I32Err(strconv.Atoi(c.QueryParam("id")))
	if id <= 0 || favType == "" {
		result.ErrCode = 1
		result.ErrMsg = "收藏失败"
	} else {
		var err error
		ms := rsi.MemberService
		if favType == "shop" {
			err = ms.FavoriteShop(int64(memberId), id)
		} else {
			err = ms.FavoriteGoods(int64(memberId), id)
		}
		result.Error(err)
	}
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), result)
}

// 二维码
func (s *serviceC) QrCode(c *echox.Context) error {
	qurl := c.QueryParam("url")
	u, err := url.Parse(qurl)
	if err != nil || !strings.HasSuffix(u.Host, variable.Domain) {
		return c.StringOK("not service")
	}
	data := gen.BuildQrCodeForUrl(qurl, 20)
	c.Response().Header().Set("Content-Type", "Image/Jpeg")
	c.Response().Header().Add("Content-Disposition", "attachment;filename=qr.jpg")
	c.Response().Write(data)
	return nil
}
