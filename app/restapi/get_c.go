/**
 * Copyright 2015 @ z3q.net.
 * name : get_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"crypto/sha1"
	"fmt"
	"go2o/app/util"
	"go2o/core/infrastructure/gen"
	"go2o/core/service/dps"
	"gopkg.in/labstack/echo.v1"
	"io"
	"strconv"
)

type getC struct {
}

// 下载邀请二维码
func (this *getC) Invite_qr(ctx *echo.Context) error {
	domain := ctx.Query("domain")                       //域名
	memberId, _ := strconv.Atoi(ctx.Query("member_id")) //会员编号
	targetUrl := ctx.Query("target_url")                //目标跳转地址
	if len(domain) == 0 {
		domain = "http://" + ctx.Request().Host
	}
	if len(targetUrl) == 0 {
		targetUrl = dps.BaseService.GetRegisterPerm().CallBackUrl
	}
	m := dps.MemberService.GetMember(memberId)
	if m != nil {
		ctx.Response().Header().Add("Content-Type", "Image/Jpeg")
		ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=tgcode_%s.jpg", m.InvitationCode))
		ctx.Response().Write(util.GenerateInvitationQr(domain, m.InvitationCode, targetUrl))
	}
	return nil
}

// 创建二维码
func (this *getC) GenQr(ctx *echo.Context) error {
	link := ctx.Query("url")
	qrBytes := gen.BuildQrCodeForUrl(link, 10)
	t := sha1.New()
	io.WriteString(t, link)
	hash := string(t.Sum(nil))
	ctx.Response().Header().Add("Content-Type", "Image/Jpeg")
	ctx.Response().Header().Set("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s.jpg", hash))
	ctx.Response().Write(qrBytes)
	return nil
}
