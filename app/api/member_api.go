package api

import (
	"github.com/ixre/gof/api"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/thrift"
	"strconv"
	"strings"
	"time"
)

var _ api.Handler = new(MemberApi)

type MemberApi struct {
	*apiUtil
}

func (m MemberApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"login":   m.login,
		"get":m.getMember,
		"account": m.account,
		"profile": m.profile,
		"checkToken":m.checkToken,
	})
}

// 登录
func (m MemberApi) login(ctx api.Context) interface{} {
	form := ctx.Form()
	user := strings.TrimSpace(form.GetString("user"))
	pwd := strings.TrimSpace(form.GetString("pwd"))
	if len(user) == 0 {
		return dto.MemberLoginResult{
			ErrCode: 2,
			ErrMsg:  "缺少参数:user",
			Member:  nil,
		}
	}
	if len(pwd) == 0 {
		return dto.MemberLoginResult{
			ErrCode: 2,
			ErrMsg:  "缺少参数:pwd",
			Member:  nil,
		}
	}
	var result dto.MemberLoginResult
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
				Code:       r.Data["member_code"],
				Token:      token,
				UpdateTime: time.Now().Unix(),
			}
		}
	}
	return result
}

// 注册
func (m MemberApi) Register(ctx api.Context) interface{} {
	user := ctx.Form().GetString("user")
	pwd := ctx.Form().GetString("pwd")
	phone := ctx.Form().GetString("phone")
	regFrom := ctx.Form().GetString("reg_from")       // 注册来源
	inviteCode := ctx.Form().GetString("invite_code") // 邀请码
	regIp := ctx.Form().GetString("$user_ip_addr")    // IP地址
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mp := map[string]string{
			"reg_ip":      regIp,
			"reg_from":    regFrom,
			"invite_code": inviteCode,
		}
		r, _ := cli.RegisterMemberV2(thrift.Context, user, pwd, 0, "", phone, "", "", mp)
		return r
	}
	return m.SResult(err)
}

// 账号信息
func (m MemberApi) account(ctx api.Context) interface{} {
	memberId := ctx.Form().GetInt("memberId")
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		r, err1 := cli.GetAccount(thrift.Context, int64(memberId))
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}

// 账号信息
func (m MemberApi) profile(ctx api.Context) interface{} {
	memberId := ctx.Form().GetInt("memberId")
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		r, err1 := cli.GetProfile(thrift.Context, int64(memberId))
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}
// 账号信息
func (m MemberApi) checkToken(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	if len(code) == 0 || len(token) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId,_ := cli.GetMemberId(thrift.Context,code)
		r, err1 := cli.CheckToken(thrift.Context, memberId, token)
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}

// 获取会员信息
func (m MemberApi) getMember(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0{
		return api.NewErrorResponse("missing params: code")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId,_ := cli.GetMemberId(thrift.Context,code)
		if memberId <= 0{
			return api.NewErrorResponse("no such member")
		}
		r, _ := cli.GetMember(thrift.Context, memberId)
		return r
	}
	return api.NewErrorResponse(err.Error())
}