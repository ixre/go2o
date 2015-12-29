/**
 * Copyright 2015 @ z3q.net.
 * name : ucenter
 * author : jarryliu
 * date : 2015-12-13 11:39
 * description :
 * history :
 */
package ucenter

import (
	"github.com/jsix/gof/web/session"
	"github.com/labstack/echo"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/url"
	"strings"
	"sync"
)

var (
	mux sync.Mutex
)

// 会员登陆检查
func memberLogonCheck(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		path := ctx.Request().URL.Path

		const ignore string = "/login|/partner_connect|/msc|/msd|/partner_disconnect|"
		if strings.Index(ignore, path+"|") != -1 || strings.HasPrefix(path, "/static/") {
			return h(ctx)
		}
		session := session.Default(ctx.Response(), ctx.Request())
		if m := session.Get("member"); m != nil {
			return h(ctx)
		}
		ctx.Response().Header().Set("Location", "/login?return_url="+
			url.QueryEscape(ctx.Request().URL.String()))
		ctx.Response().WriteHeader(302)
		return nil
	}
}

// 获取会员编号
func GetSessionMemberId(ctx *echox.Context) int {
	if m := ctx.Session.Get("member"); m != nil {
		return m.(*member.ValueMember).Id
	}
	return 0
}

// 获取会员
func getMember(ctx *echox.Context) *member.ValueMember {
	s := ctx.Session.Get("member")
	if s != nil {
		return s.(*member.ValueMember)
	}
	return nil
}

// 重新缓存会员
func reCacheMember(ctx *echox.Context, memberId int) *member.ValueMember {
	v := dps.MemberService.GetMember(memberId)
	ctx.Session.Set("member", v)
	ctx.Session.Save()
	return v
}

// 获取商户
func getPartner(ctx *echox.Context) *partner.ValuePartner {
	val := ctx.Session.Get("member:rel_partner")
	if val != nil {
		return cache.GetValuePartnerCache(val.(int))
	} else {
		m := getMember(ctx)
		if m != nil {
			rel := dps.MemberService.GetRelation(m.Id)
			ctx.Session.Set("member:rel_partner", rel.RegisterPartnerId)
			ctx.Session.Save()
			return cache.GetValuePartnerCache(rel.RegisterPartnerId)
		}
	}
	return nil
}

// 获取商户的站点设置
func getSiteConf(partnerId int) *partner.SiteConf {
	return cache.GetPartnerSiteConf(partnerId)
}
