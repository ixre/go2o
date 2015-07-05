/**
 * Copyright 2015 @ S1N1 Team.
 * name : get_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api

import (
	"github.com/atnet/gof/web"
	"go2o/src/core/infrastructure/gen"
	"go2o/src/core/service/dps"
	"strconv"
)

type getC struct {
	*BaseC
}

// 下载邀请二维码
func (this *getC) Invite_qr(ctx *web.Context) {
	form := ctx.Request.URL.Query()
	domain := form.Get("domain")
	memberId, _ := strconv.Atoi(form.Get("member_id"))
	m := dps.MemberService.GetMember(memberId)
	if m != nil {
		var url string = domain + "/main/t/" + m.InvitationCode
		var qrBytes []byte = gen.BuildQrCodeForUrl(url)
		ctx.Response.Header().Add("Content-Type", "Image/Jpeg")
		ctx.Response.Write(qrBytes)
	}
}
