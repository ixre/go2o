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
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"github.com/labstack/echo"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"go2o/core/variable"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
		return c.JSON(http.StatusOK, result)
	}
	trans, cli, err := thrift.MemberServeClient()
	if err != nil {
		result.Message = "网络连接失败"
	} else {
		defer trans.Close()
		encPwd := domain.MemberSha1Pwd(pwd)
		r, _ := cli.CheckLogin(thrift.Context, usr, encPwd, true)
		result.Message = r.ErrMsg
		result.Result = r.ErrCode == 0
		if r.ErrCode == 0 {
			memberId, _ := strconv.Atoi(r.Data["MemberId"])
			token, _ := cli.GetToken(thrift.Context, int64(memberId), false)
			result.Member = &dto.LoginMember{
				Id:         memberId,
				Token:      token,
				UpdateTime: time.Now().Unix(),
			}
		}
	}
	return c.JSON(http.StatusOK, result)
}

// 注册
func (mc *MemberC) Register(c echo.Context) error {
	result := &gof.Result{}
	return c.JSON(http.StatusOK,
		result.Error(errors.New("注册暂停，请通过微信或其他方式注册!")))

	r := c.Request()
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
	m := &member_service.SMember{}
	pro := &member_service.SProfile{}
	m.Usr = usr
	m.Pwd = domain.MemberSha1Pwd(pwd)
	m.RegIp = regIp
	m.RegFrom = registerFrom
	pro.Phone = phone
	pro.Name = m.Usr
	_, err := rsi.MemberService.RegisterMember(mchId,
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
	memberId := GetMemberId(c)
	mut, _ = strconv.Atoi(form.Get("member_update_time"))
	aut, _ = strconv.Atoi(form.Get("account_update_time"))
	mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, memberId)
	sto.Get(mutKey, &kvMut)
	autKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
	sto.Get(autKey, &kvAut)
	if kvMut == 0 {
		m, _ := rsi.MemberService.GetMember(thrift.Context, memberId)
		kvMut = int(m.UpdateTime)
		sto.Set(mutKey, kvMut)
	}
	//kvAut = 0
	if kvAut == 0 {
		acc, _ := rsi.MemberService.GetAccount(thrift.Context, memberId)
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
	memberId := GetMemberId(c)
	m, _ := rsi.MemberService.GetMember(thrift.Context, memberId)
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		m.DynamicToken, _ = cli.GetToken(thrift.Context, memberId, false)
	}
	return c.JSON(http.StatusOK, m)
}

// 汇总信息
func (mc *MemberC) Summary(c echo.Context) error {
	memberId := GetMemberId(c)
	v, _ := rsi.MemberService.Complex(thrift.Context, memberId)
	return c.JSON(http.StatusOK, v)
}

// 获取最新的会员账户信息
func (mc *MemberC) Account(c echo.Context) error {
	memberId := GetMemberId(c)
	m, _ := rsi.MemberService.GetAccount(thrift.Context, memberId)
	return c.JSON(http.StatusOK, m)
}

// 断开
// todo: token不允许删除，只能自动过期
func (mc *MemberC) Disconnect(c echo.Context) error {
	result := &gof.Result{}
	return c.JSON(http.StatusOK, result)
	//mStr := c.QueryParam("member_id")
	//memberId, err := util.I32Err(strconv.Atoi(mStr))
	//token := c.QueryParam("token")
	//trans,cli, err := thrift.MemberClient()
	//if err == nil {
	//	defer trans.Close()
	//	if b, _ := cli.CheckToken(thrift.Context,memberId, token); b {
	//		cli.RemoveToken(memberId)
	//	} else {
	//		err = errors.New("error credential")
	//	}
	//}
	//return c.JSON(http.StatusOK, result.Error(err))
}
