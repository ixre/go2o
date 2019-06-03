/**
 * Copyright 2015 @ to2.net.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"fmt"
	"github.com/ixre/gof"
	"github.com/labstack/echo"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
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
	user := strings.TrimSpace(r.FormValue("user"))
	pwd := strings.TrimSpace(r.FormValue("pwd"))
	if len(user) == 0 || len(pwd) == 0 {
		result.ErrMsg = "会员不存在"
		return c.JSON(http.StatusOK, result)
	}
	trans, cli, err := thrift.MemberServeClient()
	if err != nil {
		result.ErrMsg = "网络连接失败"
	} else {
		defer trans.Close()
		encPwd := domain.MemberSha1Pwd(pwd)
		r, _ := cli.CheckLogin(thrift.Context, user, encPwd, true)
		result.ErrMsg = r.ErrMsg
		result.ErrCode = int(r.ErrCode)
		if r.ErrCode == 0 {
			memberId, _ := strconv.Atoi(r.Data["member_id"])
			token, _ := cli.GetToken(thrift.Context, int64(memberId), false)
			result.Member = &dto.LoginMember{
				ID:         memberId,
				Token:      token,
				UpdateTime: time.Now().Unix(),
			}
		}
	}
	return c.JSON(http.StatusOK, result)
}

func (mc *MemberC) Ping(c echo.Context) error {
	//log.Println("---", ctx.Request.FormValue("member_id"), ctx.Request.FormValue("member_token"))
	return c.String(http.StatusOK, "PONG")
}

// 同步 todo: 过时
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
		//m, _ := rsi.MemberService.GetMember(thrift.Context, memberId)
		//kvMut = int(m.UpdateTime)
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
