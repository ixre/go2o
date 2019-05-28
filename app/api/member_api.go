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
}

func (m MemberApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"login":   m.login,
		"account": m.account,
		"profile": m.profile,
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
			memberId, _ := strconv.Atoi(r.Data["MemberId"])
			token, _ := cli.GetToken(thrift.Context, int64(memberId), false)
			result.Member = &dto.LoginMember{
				ID:         memberId,
				Token:      token,
				UpdateTime: time.Now().Unix(),
			}
		}
	}
	return result
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
