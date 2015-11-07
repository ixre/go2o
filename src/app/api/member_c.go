/**
 * Copyright 2015 @ z3q.net.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api

import (
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/web"
	"github.com/jsix/gof/web/mvc"
	"go2o/src/app/util"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/dto"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"strconv"
	"strings"
)

var _ mvc.Filter = new(MemberC)

type MemberC struct {
	*BaseC
}

// 会员登陆后才能调用接口
func (this *MemberC) Requesting(ctx *web.Context) bool {
	if this.BaseC == nil || !this.BaseC.Requesting(ctx) {
		return false
	}
	rlt := this.BaseC.CheckMemberToken(ctx)
	//fmt.Printf("%#v\n",ctx.Request.Form)
	return rlt
}

// 登陆
func (this *MemberC) Login(ctx *web.Context) {
	if this.BaseC.Requesting(ctx) {

		r := ctx.Request
		var usr, pwd string = r.FormValue("usr"), r.FormValue("pwd")
		partnerId := this.GetPartnerId(ctx)
		var result dto.MemberLoginResult

		pwd = strings.TrimSpace(pwd)

		if len(usr) == 0 || len(pwd) == 0 {
			result.Message = "会员不存在"
		} else {
			encodePwd := domain.MemberSha1Pwd(pwd)
			b, e, err := dps.MemberService.Login(partnerId, usr, encodePwd)
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
		var err error

		var partnerId int = this.GetPartnerId(ctx)
		var invMemberId int // 邀请人
		var usr string = r.FormValue("usr")
		var pwd string = r.FormValue("pwd")
		var phone string = r.FormValue("phone")
		var registerFrom string = r.FormValue("reg_from")          // 注册来源
		var invitationCode string = r.FormValue("invitation_code") // 推荐码
		var regIp string
		if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
			regIp = r.RemoteAddr[:i]
		}

		if err = dps.PartnerService.CheckRegisterMode(partnerId, invitationCode); err != nil {
			result.Message = err.Error()
			ctx.Response.JsonOutput(result)
			return
		}

		//fmt.Println(usr, pwd, "REGFROM:", registerFrom, "INVICODE:", invitationCode)

		// 检验
		if len(invitationCode) != 0 {
			invMemberId = dps.MemberService.GetMemberIdByInvitationCode(invitationCode)
			if invMemberId == 0 {
				result.Message = "1011:推荐码错误"
				ctx.Response.JsonOutput(result)
				return
			}
		}

		var member member.ValueMember
		member.Usr = usr
		member.Pwd = domain.MemberSha1Pwd(pwd)
		member.RegIp = regIp
		member.Phone = phone
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

func (this *MemberC) Ping(ctx *web.Context) {
	//log.Println("---", ctx.Request.FormValue("member_id"), ctx.Request.FormValue("member_token"))
	ctx.Response.Write([]byte("pang"))
}

// 同步
func (this *MemberC) Async(ctx *web.Context) {
	var rlt AsyncResult
	var form = ctx.Request.Form
	var mut, aut, kvMut, kvAut int
	memberId := this.GetMemberId(ctx)
	mut, _ = strconv.Atoi(form.Get("member_update_time"))
	aut, _ = strconv.Atoi(form.Get("account_update_time"))
	mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, memberId)
	ctx.App.Storage().Get(mutKey, &kvMut)
	autKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
	ctx.App.Storage().Get(autKey, &kvAut)
	if kvMut == 0 {
		m := dps.MemberService.GetMember(memberId)
		kvMut = int(m.UpdateTime)
		ctx.App.Storage().Set(mutKey, kvMut)
	}
	//kvAut = 0
	if kvAut == 0 {
		acc := dps.MemberService.GetAccount(memberId)
		kvAut = int(acc.UpdateTime)
		ctx.App.Storage().Set(autKey, kvAut)
	}
	rlt.MemberId = memberId
	rlt.MemberUpdated = kvMut != mut
	rlt.AccountUpdated = kvAut != aut
	ctx.Response.JsonOutput(rlt)
}

// 获取最新的会员信息
func (this *MemberC) Get(ctx *web.Context) {
	memberId := this.GetMemberId(ctx)
	m := dps.MemberService.GetMember(memberId)
	m.DynamicToken, _ = util.GetMemberApiToken(ctx.App.Storage(), memberId)
	ctx.Response.JsonOutput(m)
}

// 汇总信息
func (this *MemberC) Summary(ctx *web.Context) {
	memberId := this.GetMemberId(ctx)
	var updateTime int64 = dps.MemberService.GetMemberLatestUpdateTime(memberId)
	var v *dto.MemberSummary = new(dto.MemberSummary)
	var key = fmt.Sprintf("cache:member:summary:%d", memberId)
	if cache.GetKVS().Get(key, &v) != nil || v.UpdateTime < updateTime {
		v = dps.MemberService.GetMemberSummary(memberId)
		cache.GetKVS().SetExpire(key, v, 3600*48) // cache 48 hours
	}
	ctx.Response.JsonOutput(v)
}

// 获取最新的会员账户信息
func (this *MemberC) Account(ctx *web.Context) {
	memberId := this.GetMemberId(ctx)
	m := dps.MemberService.GetAccount(memberId)
	ctx.Response.JsonOutput(m)
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
