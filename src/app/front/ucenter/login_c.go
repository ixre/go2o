/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package ucenter

import (
	"github.com/atnet/gof/web"
	"go2o/src/app/util"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/core/service/goclient"
	"strconv"
	"github.com/atnet/gof"
	"encoding/json"
)

type loginC struct {
}

//登陆
func (this *loginC) Index(ctx *web.Context) {
	executeTemplate(ctx,nil, nil, "views/ucenter/{device}/login.html")
}
func (this *loginC) Index_post(ctx *web.Context) {
	r := ctx.Request
	r.ParseForm()
var result gof.Message
	usr, pwd := r.Form.Get("usr"), r.Form.Get("pwd")
	b, m, err := dps.MemberService.Login(-1, usr, pwd)
	if b {
		ctx.Session().Set("member", m)
		ctx.Session().Save()
		result.Result = true
	}else{
		if err != nil{
			result.Message = err.Error()
		}else{
			result.Message = "登陆失败"
		}
	}
	js,_ := json.Marshal(result)
	ctx.ResponseWriter.Write(js)

}

//从partner登录过来的信息
func (this *loginC) Partner_connect(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	sessionId := r.URL.Query().Get("sessionId")
	var m *member.ValueMember
	var err error

	if sessionId == "" {
		// 第三方连接，传入memberId 和 token
		memberId, err := strconv.Atoi(r.URL.Query().Get("mid"))
		token := r.URL.Query().Get("token")
		if err == nil && token != "" {
			m, err = goclient.Member.GetMember(memberId, token)
			ctx.Session().Set("member", m)
		}
	} else {
		// 从统一平台连接过来（标准版商户PC前端)
		ctx.Session().UseInstead(sessionId)
		m = ctx.Session().Get("member").(*member.ValueMember)
	}

	// 设置访问设备
	util.SetBrownerDevice(ctx, ctx.Request.URL.Query().Get("device"))

	if err == nil || m != nil {
		rl := dps.MemberService.GetRelation(m.Id)
		if rl.RegisterPartnerId > 0 {
			ctx.Session().Set("member:rel_partner", rl.RegisterPartnerId)
			ctx.Session().Save()
			w.Write([]byte("<script>location.replace('/')</script>"))
			return
		}
	}
	w.Write([]byte("<script>location.replace('/login')</script>"))
}

//从partner端退出
func (this *loginC) Partner_disconnect(ctx *web.Context) {
	ctx.Session().Destroy()
	ctx.ResponseWriter.Write([]byte("{state:1}"))
}
