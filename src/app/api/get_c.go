/**
 * Copyright 2015 @ z3q.net.
 * name : get_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api

import (
	"github.com/jsix/gof/web"
	"go2o/src/core/infrastructure/gen"
	"go2o/src/core/service/dps"
	"net/url"
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
		var url string = domain + "/main/change_device?device=3&return_url=/main/t/" + m.InvitationCode +
			url.QueryEscape("?return_url=/main/app")
		var qrBytes []byte = gen.BuildQrCodeForUrl(url)
		ctx.Response.Header().Add("Content-Type", "Image/Jpeg")
		ctx.Response.Write(qrBytes)
	}
}
