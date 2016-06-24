/**
 * Copyright 2015 @ z3q.net.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"fmt"
	"github.com/jsix/gof"
	"go2o/app/cache"
	"go2o/app/util"
	"go2o/core/domain/interface/member"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/dps"
	"go2o/core/variable"
	"go2o/x/echox"
	"gopkg.in/labstack/echo.v1"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 会员登陆后才能调用接口
type MemberC struct {
}

// 登陆
func (this *MemberC) Login(ctx *echo.Context) error {
	r := ctx.Request()
	var usr, pwd string = r.FormValue("usr"), r.FormValue("pwd")
	//merchantId := getMerchantId(ctx)
	var result dto.MemberLoginResult

	pwd = strings.TrimSpace(pwd)

	if len(usr) == 0 || len(pwd) == 0 {
		result.Message = "会员不存在"
	} else {
		encodePwd := domain.MemberSha1Pwd(pwd)
		e, err := dps.MemberService.TryLogin(usr, encodePwd, true)

		if err == nil {
			// 生成令牌
			e.DynamicToken = util.SetMemberApiToken(sto, e.Id, e.Pwd)
			mm := dps.MemberService.GetMemberSummary(e.Id)
			result.Member = mm
			result.Result = true
		} else {
			result.Message = err.Error()
			result.Result = false
		}
	}
	return ctx.JSON(http.StatusOK, result)

}

// 注册
func (this *MemberC) Register(ctx *echo.Context) error {
	r := ctx.Request()
	var result dto.MessageResult
	var err error
	var merchantId int = getMerchantId(ctx)
	var usr string = r.FormValue("usr")
	var pwd string = r.FormValue("pwd")
	var phone string = r.FormValue("phone")
	var registerFrom string = r.FormValue("reg_from")          // 注册来源
	var invitationCode string = r.FormValue("invitation_code") // 邀请码
	var regIp string
	if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
		regIp = r.RemoteAddr[:i]
	}

	m := &member.Member{}
	pro := &member.Profile{}
	m.Usr = usr
	m.Pwd = domain.MemberSha1Pwd(pwd)
	m.RegIp = regIp
	m.RegFrom = registerFrom

	pro.Phone = phone
	pro.Name = m.Usr

	_, err = dps.MemberService.RegisterMember(merchantId, m, pro, "", invitationCode)
	if err == nil {
		result.Result = true
	} else {
		result.Message = err.Error()
	}
	return ctx.JSON(http.StatusOK, result)
}

func (this *MemberC) Ping(ctx *echo.Context) error {
	//log.Println("---", ctx.Request.FormValue("member_id"), ctx.Request.FormValue("member_token"))
	return ctx.String(http.StatusOK, "PONG")
}

// 同步
func (this *MemberC) Async(ctx *echo.Context) error {
	var rlt AsyncResult
	var form = url.Values(ctx.Request().Form)
	var mut, aut, kvMut, kvAut int
	memberId := GetMemberId(ctx)
	mut, _ = strconv.Atoi(form.Get("member_update_time"))
	aut, _ = strconv.Atoi(form.Get("account_update_time"))
	mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, memberId)
	sto.Get(mutKey, &kvMut)
	autKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
	sto.Get(autKey, &kvAut)
	if kvMut == 0 {
		m := dps.MemberService.GetMember(memberId)
		kvMut = int(m.UpdateTime)
		sto.Set(mutKey, kvMut)
	}
	//kvAut = 0
	if kvAut == 0 {
		acc := dps.MemberService.GetAccount(memberId)
		kvAut = int(acc.UpdateTime)
		sto.Set(autKey, kvAut)
	}
	rlt.MemberId = memberId
	rlt.MemberUpdated = kvMut != mut
	rlt.AccountUpdated = kvAut != aut
	return ctx.JSON(http.StatusOK, rlt)
}

// 获取最新的会员信息
func (this *MemberC) Get(ctx *echo.Context) error {
	memberId := GetMemberId(ctx)
	m := dps.MemberService.GetMember(memberId)
	m.DynamicToken, _ = util.GetMemberApiToken(sto, memberId)
	return ctx.JSON(http.StatusOK, m)
}

// 汇总信息
func (this *MemberC) Summary(ctx *echo.Context) error {
	memberId := GetMemberId(ctx)
	var updateTime int64 = dps.MemberService.GetMemberLatestUpdateTime(memberId)
	var v *dto.MemberSummary = new(dto.MemberSummary)
	var key = fmt.Sprintf("cac:mm:summary:%d", memberId)
	if cache.GetKVS().Get(key, &v) != nil || v.UpdateTime < updateTime {
		v = dps.MemberService.GetMemberSummary(memberId)
		cache.GetKVS().SetExpire(key, v, 3600*48) // cache 48 hours
	}
	return ctx.JSON(http.StatusOK, v)
}

// 获取最新的会员账户信息
func (this *MemberC) Account(ctx *echo.Context) error {
	memberId := GetMemberId(ctx)
	m := dps.MemberService.GetAccount(memberId)
	return ctx.JSON(http.StatusOK, m)
}

// 断开
func (this *MemberC) Disconnect(ctx *echox.Context) error {
	var result gof.Message
	if util.MemberHttpSessionDisconnect(ctx) {
		result.Result = true
	} else {
		result.Message = "disconnect fail"
	}
	return ctx.JSON(http.StatusOK, result)
}
