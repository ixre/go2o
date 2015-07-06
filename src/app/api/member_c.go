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
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/dto"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"strings"
	"go2o/src/cache"
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

// 登陆
func (this *MemberC) Login(ctx *web.Context) {
	if this.BaseC.Requesting(ctx) {

		r := ctx.Request
		var usr, pwd string = r.FormValue("usr"), r.FormValue("pwd")
		partnerId := this.GetPartnerId(ctx)
		var result dto.MemberLoginResult

		if len(usr) == 0 || len(pwd) == 0 {
			result.Message = "会员不存在"
		} else {
			b, e, err := dps.MemberService.Login(partnerId, usr, pwd)
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
		ctx.Response.JsonOutput(result)
	}
}

// 注册
func (this *MemberC) Register(ctx *web.Context) {
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
				ctx.Response.JsonOutput(result)
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

		ctx.Response.JsonOutput(result)
	}
}

// 断开
func (this *MemberC) Disconnect(ctx *web.Context) {
	var result gof.Message
	if util.MemberHttpSessionDisconnect(ctx) {
		result.Result = true
	} else {
		result.Message = "disconnect fail"
	}
	ctx.Response.JsonOutput(result)
}

// 汇总信息
func (this *MemberC) Summary(ctx *web.Context){
	memberId := this.GetMemberId(ctx)
	var updateTime int64 = dps.MemberService.GetMemberLatestUpdateTime(memberId)
	var v *dto.MemberSummary = new(dto.MemberSummary)
	var key = fmt.Sprintf("cache:member:summary:%d",memberId)
	if cache.GetKVS().Get(key,&v) != nil || v.UpdateTime < updateTime {
		v = dps.MemberService.GetMemberSummary(memberId)
		cache.GetKVS().SetExpire(key, v, 3600*48) // cache 48 hours
	}
	ctx.Response.JsonOutput(v)
}