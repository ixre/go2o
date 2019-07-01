/**
 * Copyright 2015 @ to2.net.
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
	"go2o/core/infrastructure/gen"
	"io"
)

type getC struct {
}

// 下载邀请二维码
func (g *getC) Invite_qr(c echo.Context) error {
	panic("not implement")
	//domain := c.QueryParam("domain")                                    //域名
	//memberId,_ := util.I64Err(strconv.Atoi(c.QueryParam("member_id"))) //会员编号
	//targetUrl := c.QueryParam("target_url")                             //目标跳转地址
	//if len(domain) == 0 {
	//	domain = "http://" + c.Request().Host
	//}
	//if len(targetUrl) == 0 {
	//	targetUrl = rsi.FoundationService.GetRegisterPerm().CallBackUrl
	//}
	//m, _ := rsi.MemberService.GetMember(thrift.Context, memberId)
	//if m != nil {
	//	query := "return_url=" + targetUrl
	//	c.Response().Header().Add("Content-Type", "Image/Jpeg")
	//	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=tgcode_%s.jpg", m.InvitationCode))
	//	c.Response().Write(autil.GenerateInvitationQr(domain, m.InvitationCode, query))
	//}
	//return nil
}

// 创建二维码
func (g *getC) GenQr(c echo.Context) error {
	link := c.QueryParam("url")
	qrBytes := gen.BuildQrCodeForUrl(link, 10)
	t := sha1.New()
	io.WriteString(t, link)
	hash := string(t.Sum(nil))
	c.Response().Header().Add("Content-Type", "Image/Jpeg")
	c.Response().Header().Set("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s.jpg", hash))
	c.Response().Write(qrBytes)
	return nil
}
