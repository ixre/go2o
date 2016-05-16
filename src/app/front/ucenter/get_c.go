/**
 * Copyright 2015 @ z3q.net.
 * name : basic_c
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ucenter

import (
	"errors"
	"fmt"
	"go2o/src/core/infrastructure/gen"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
)

type getC struct {
}

func (this *getC) GetQR(ctx *echox.Context) error {
	var domain string
	var code string
	//var w int
	code = ctx.P(0)
	if len(code) == 0 {
		return errors.New("Code error")
	}
	//w,_ =  strconv.Atoi(ctx.P(1))
	//if w < 40 || w > 200{
	//	w = 60
	//}
	domain = ctx.Query("domain")
	if len(domain) == 0 {
		id := dps.MemberService.GetMemberIdByInvitationCode(code)
		if id <= 0 {
			return errors.New("Code error")
		}
		rl := dps.MemberService.GetRelation(id)
		if rl == nil {
			return errors.New("Code error")
		}
		merchantId := rl.RegisterMerchantId
		pt, _ := dps.PartnerService.GetMerchant(merchantId)
		if pt == nil {
			return errors.New("Except member")
		}
		domain = fmt.Sprintf("http://%s.%s", pt.Usr,
			ctx.App.Config().GetString(variable.ServerDomain))
	}

	var url string = domain + "/change_device?device=3&return_url=/main/t/" + code
	qrBytes := gen.BuildQrCodeForUrl(url, 10)
	ctx.Response().Header().Add("Content-Type", "Image/Jpeg")
	ctx.HttpResponse().Header().Set("Content-Disposition",
		fmt.Sprintf("attachment;filename=tgcode_%s.jpg", code))
	ctx.HttpResponse().Write(qrBytes)
	return nil
}
