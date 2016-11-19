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
	"github.com/jsix/goex/echox"
	"github.com/jsix/gof"
	"github.com/labstack/echo"
	"go2o/app/cache"
	autil "go2o/app/util"
	"go2o/core/domain/interface/member"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/dps"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/variable"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 会员登录后才能调用接口
type MemberC struct {
}

// 登录
func (mc *MemberC) Login(c echo.Context) error {
	var result dto.MemberLoginResult
	r := c.Request()
	usr := strings.TrimSpace(r.FormValue("usr"))
	pwd := strings.TrimSpace(r.FormValue("pwd"))
	if len(usr) == 0 || len(pwd) == 0 {
		result.Message = "会员不存在"
	} else {
		encodePwd := domain.MemberSha1Pwd(pwd)
		mp, _ := dps.MemberService.Login(usr, encodePwd, true)
		id := mp["Id"]
		rst := mp["Result"]
		if id > 0 {
			m, _ := dps.MemberService.GetMember(id)
			// 登录成功，生成令牌
			token := autil.SetMemberApiToken(sto, id, encodePwd)
			result.Member = &dto.LoginMember{
				Id:         int(id),
				Token:      token,
				UpdateTime: m.UpdateTime,
			}
			result.Result = true
		} else {
			switch rst {
			case -1:
				result.Message = member.ErrNoSuchMember.Error()
			case -2:
				result.Message = member.ErrCredential.Error()
			case -3:
				result.Message = member.ErrDisabled.Error()
			default:
				result.Message = "登陆失败"
			}
		}
	}
	return c.JSON(http.StatusOK, result)
}

// 注册
func (mc *MemberC) Register(c echo.Context) error {
	r := c.Request()
	result := gof.Message{}
	mchId := getMerchantId(c)
	usr := r.FormValue("usr")
	pwd := r.FormValue("pwd")
	phone := r.FormValue("phone")
	registerFrom := r.FormValue("reg_from")          // 注册来源
	invitationCode := r.FormValue("invitation_code") // 邀请码
	var regIp string
	if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
		regIp = r.RemoteAddr[:i]
	}
	m := &define.Member{}
	pro := &define.Profile{}
	m.Usr = usr
	m.Pwd = domain.MemberSha1Pwd(pwd)
	m.RegIp = regIp
	m.RegFrom = registerFrom
	pro.Phone = phone
	pro.Name = m.Usr
	_, err := dps.MemberService.RegisterMember(mchId,
		m, pro, "", invitationCode)
	return c.JSON(http.StatusOK, result.Error(err))
}

func (mc *MemberC) Ping(c echo.Context) error {
	//log.Println("---", ctx.Request.FormValue("member_id"), ctx.Request.FormValue("member_token"))
	return c.String(http.StatusOK, "PONG")
}

// 同步
func (mc *MemberC) Async(c echo.Context) error {
	var rlt AsyncResult
	var form = url.Values(c.Request().Form)
	var mut, aut, kvMut, kvAut int
	memberId := int32(GetMemberId(c))
	mut, _ = strconv.Atoi(form.Get("member_update_time"))
	aut, _ = strconv.Atoi(form.Get("account_update_time"))
	mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, memberId)
	sto.Get(mutKey, &kvMut)
	autKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
	sto.Get(autKey, &kvAut)
	if kvMut == 0 {
		m, _ := dps.MemberService.GetMember(memberId)
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
	return c.JSON(http.StatusOK, rlt)
}

// 获取最新的会员信息
func (mc *MemberC) Get(c echo.Context) error {
	memberId := int32(GetMemberId(c))
	m, _ := dps.MemberService.GetMember(memberId)
	m.DynamicToken, _ = autil.GetMemberApiToken(sto, memberId)
	return c.JSON(http.StatusOK, m)
}

// 汇总信息
func (mc *MemberC) Summary(c echo.Context) error {
	memberId := GetMemberId(c)
	var updateTime int64 = dps.MemberService.GetMemberLatestUpdateTime(memberId)
	var v *dto.MemberSummary = new(dto.MemberSummary)
	var key = fmt.Sprintf("cac:mm:summary:%d", memberId)
	if cache.GetKVS().Get(key, &v) != nil || v.UpdateTime < updateTime {
		v = dps.MemberService.GetMemberSummary(memberId)
		cache.GetKVS().SetExpire(key, v, 3600*48) // cache 48 hours
	}
	return c.JSON(http.StatusOK, v)
}

// 获取最新的会员账户信息
func (mc *MemberC) Account(c echo.Context) error {
	memberId := GetMemberId(c)
	m := dps.MemberService.GetAccount(memberId)
	return c.JSON(http.StatusOK, m)
}

// 断开
func (mc *MemberC) Disconnect(c *echox.Context) error {
	var result gof.Message
	if autil.MemberHttpSessionDisconnect(c) {
		result.Result = true
	} else {
		result.Message = "disconnect fail"
	}
	return c.JSON(http.StatusOK, result)
}
