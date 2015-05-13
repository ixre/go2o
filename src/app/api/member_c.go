/**
 * Copyright 2015 @ S1N1 Team.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api

import (
	"errors"
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/crypto"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/dto"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"strconv"
	"strings"
)

var _ mvc.Filter = new(memberC)

type memberC struct {
	*baseC
}

// 会员登陆后才能调用接口
func (this *memberC) Requesting(ctx *web.Context) bool {
	if this.baseC != nil && this.baseC.Requesting(ctx) {
		r := ctx.Request
		memberId, _ := strconv.Atoi(r.FormValue("member_id"))
		token := r.FormValue("token")
		if chkMemberToken(ctx.App.Storage(), memberId, token) {
			return true
		}
		this.errorOutput(ctx, "invalid request!")
	}
	return false
}

const offset string = "%$^&@#"

func chkStorage(sto gof.Storage) {
	if sto == nil {
		panic(errors.New("[ Api] - api token storage is null !"))
	}
}

func getMemberTokenKey(memberId int) string {
	return fmt.Sprintf("api:member:token:%d", memberId)
}

// 设置令牌，并返回
func setMemberToken(sto gof.Storage, memberId int, pwd string) string {
	chkStorage(sto)
	cyp := crypto.NewUnixCrypto(pwd+offset, offset)
	var token string = string(cyp.Encode())
	var key string = getMemberTokenKey(memberId)

	sto.Set(key, token)      // 存储令牌
	sto.Set(key+"base", pwd) // 存储令牌凭据

	return token
}

// 校验令牌
func chkMemberToken(sto gof.Storage, memberId int, token string) bool {
	chkStorage(sto)

	if len(token) == 0 {
		return false
	}

	var key = getMemberTokenKey(memberId)
	var srcToken, tokenBase string

	sto.Get(key, &srcToken)
	sto.Get(key+"base", &tokenBase)

	if len(srcToken) == 0 || len(tokenBase) == 0 {
		return false
	}

	cyp := crypto.NewUnixCrypto(tokenBase+offset, offset)
	b, _, _ := cyp.Compare(token)
	return b
}

// 处理请求
func (this *memberC) handle(ctx *web.Context) {
	mvc.Handle(this, ctx, false)
}

// 登陆
func (this *memberC) login(ctx *web.Context) {
	if this.baseC.Requesting(ctx) {

		r := ctx.Request
		var usr, pwd string = r.FormValue("usr"), r.FormValue("pwd")
		var result dto.MemberLoginResult

		if len(usr) == 0 || len(pwd) == 0 {
			result.Message = "会员不存在"
		} else {
			b, e, err := dps.MemberService.Login(usr, pwd)
			result.Result = b

			if b {
				// 生成令牌
				e.DynamicToken = setMemberToken(ctx.App.Storage(), e.Id, e.Pwd)
				result.Member = e
			}
			if err != nil {
				result.Message = err.Error()
			}
		}

		this.jsonOutput(ctx, result)
	}
}

// 注册
func (this *memberC) register(ctx *web.Context) {
	if this.baseC.Requesting(ctx) {
		r := ctx.Request
		var result dto.MessageResult

		var partnerId int = this.GetPartnerId(ctx)
		var invMemberId int // 邀请人
		var usr string = r.FormValue("usr")
		var pwd string = r.FormValue("pwd")
		var registerFrom string = r.FormValue("reg_from")          // 注册来源
		var invitationCode string = r.FormValue("invitation_code") // 推荐码
		var regIp string
		if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
			regIp = r.RemoteAddr[:i]
		}

		fmt.Println(usr, pwd, "REGFROM:", registerFrom, "INVICODE:", invitationCode)

		// 检验
		if len(invitationCode) != 0 {
			invMemberId = dps.MemberService.GetMemberIdByInvitationCode(invitationCode)
			if invMemberId == 0 {
				result.Message = "推荐/邀请人不存在！"
				this.jsonOutput(ctx, result)
				return
			}
		}

		var member member.ValueMember
		member.Usr = usr
		member.Pwd = domain.EncodeMemberPwd(usr, pwd)
		member.RegIp = regIp
		member.RegFrom = registerFrom

		memberId, err := dps.MemberService.SaveMember(&member)
		if err == nil {
			result.Result = true
			err = dps.MemberService.SaveRelation(memberId, "-", invMemberId, partnerId)
		}

		if err != nil {
			result.Message = err.Error()
		}

		this.jsonOutput(ctx, result)
	}
}
