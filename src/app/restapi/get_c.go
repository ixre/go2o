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
	"github.com/labstack/echo"
	"go2o/src/core/infrastructure/gen"
	"go2o/src/core/service/dps"
	"io"
	"net/url"
	"strconv"
)

type getC struct {
}

// 下载邀请二维码
func (this *getC) Invite_qr(ctx *echo.Context) error {
	domain := ctx.Query("domain")                       //域名
	memberId, _ := strconv.Atoi(ctx.Query("member_id")) //会员编号
	targetUrl := ctx.Query("target_url")                //目标跳转地址
	if len(targetUrl) == 0 {
		targetUrl = "/main/app"
	}
	m := dps.MemberService.GetMember(memberId)
	if m != nil {
		var url string = domain + "/change_device?device=3&return_url=/main/t/" + m.InvitationCode +
			url.QueryEscape("?return_url="+targetUrl)
		qrBytes := gen.BuildQrCodeForUrl(url)
		ctx.Response().Header().Add("Content-Type", "Image/Jpeg")
		ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=tgcode_%d.jpg", memberId))
		ctx.Response().Write(qrBytes)
	}
	return nil
}

// 创建二维码
func (this *getC) GenQr(ctx *echo.Context) error {
	link := ctx.Query("url")
	qrBytes := gen.BuildQrCodeForUrl(link)
	t := sha1.New()
	io.WriteString(t, link)
	hash := string(t.Sum(nil))
	ctx.Response().Header().Add("Content-Type", "Image/Jpeg")
	ctx.Response().Header().Set("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s.jpg", hash))
	ctx.Response().Write(qrBytes)
	return nil
}
