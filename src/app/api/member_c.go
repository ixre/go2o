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
	"fmt"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/dto"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"strings"
)

var _ mvc.Filter = new(MemberC)

type MemberC struct {
	*BaseC
}

// 会员登陆后才能调用接口
func (this *MemberC) Requesting(ctx *web.Context) bool {
	return this.BaseC != nil && this.BaseC.Requesting(ctx) &&
	this.BaseC.CheckMemberToken(ctx)
}

// 处理请求
func (this *MemberC) handle(ctx *web.Context) {
	mvc.Handle(this, ctx, false)
}

// 登陆
func (this *MemberC) login(ctx *web.Context) {
	if this.BaseC.Requesting(ctx) {

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
				e.DynamicToken = util.SetMemberApiToken(ctx.App.Storage(), e.Id, e.Pwd)
				result.Member = e
			}
			if err != nil {
				result.Message = err.Error()
			}
		}

		this.JsonOutput(ctx, result)
	}
}

// 注册
func (this *MemberC) register(ctx *web.Context) {
	if this.BaseC.Requesting(ctx) {
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
				this.JsonOutput(ctx, result)
				return
			}
		}

		var member member.ValueMember
		member.Usr = usr
		member.Pwd = domain.Md5MemberPwd(usr, pwd)
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

		this.JsonOutput(ctx, result)
	}
}
