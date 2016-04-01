/**
 * Copyright 2015 @ z3q.net.
 * name : personfinance_c
 * author : jarryliu
 * date : 2016-04-01 10:15
 * description :
 * history :
 */
package ucenter

import (
	"github.com/jsix/gof"
	"go2o/src/core/domain/interface/personfinance"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
    "go2o/src/core/variable"
    "strings"
)

type personFinanceRiseC struct {
}


func (this *personFinanceRiseC) chkOpenState(memberId int,c *echox.Context)(bool,error){
   _,err := dps.PersonFinanceService.GetRiseInfo(memberId)
    b := err == nil
    if !b{ //如果未开通,则跳转到开通界面
        return false,c.Redirect(http.StatusFound,"openService?return_url="+
            c.Request().URL.String())
    }
    return true,nil
}

func (this *personFinanceRiseC) OpenService(c *echox.Context) error {
	if c.Request().Method == "POST" {
		return this.openService_post(c)
	}
	d := c.NewData()
    d.Map["AliasRise"] = variable.AliasRisePersonFinance
	return c.Render(http.StatusOK, "rise_openservice.html", d)
}

func (this *personFinanceRiseC) openService_post(c *echox.Context) error {
	msg := gof.Message{}
	memberId := GetSessionMemberId(c)
	if err := dps.PersonFinanceService.OpenRiseService(memberId); err != nil {
        if strings.Index(err.Error(),"exists") != -1{
            msg.Message ="您已经开通服务!"
        }else{
            msg.Message = err.Error()
        }
	} else {
		msg.Result = true
	}
	return c.JSON(http.StatusOK, msg)
}

func (this *personFinanceRiseC) TransferIn(c *echox.Context) error {
    memberId := GetSessionMemberId(c)
    if b,err := this.chkOpenState(memberId,c);!b{
        return err
    }
	d := c.NewData()
	d.Map["min_tranfer"] = format.FormatFloat(personfinance.MinRiseTransferInAmount)
	return c.Render(http.StatusOK, "rise_transferin.html", d)
}

func (this *personFinanceRiseC) TransferOut(c *echox.Context) error {
	if c.Request().Method == "POST" {
		return this.transferOut_post(c)
	}
    memberId := GetSessionMemberId(c)
    if b,err := this.chkOpenState(memberId,c);!b{
        return err
    }
	d := c.NewData()
	d.Map["min_tranfer"] = format.FormatFloat(personfinance.MinRiseTransferOutAmount)
	return c.Render(http.StatusOK, "rise_transferout.html", d)
}

func (this *personFinanceRiseC) transferOut_post(c *echox.Context) error {
	msg := gof.Message{}
	memberId := GetSessionMemberId(c)
	amount, err := strconv.ParseFloat(c.Form("amount"), 32)
	if err == nil {
		err = dps.PersonFinanceService.RiseTransferOut(memberId, float32(amount))
	}
	if err != nil {
		msg.Message = msg.Message
	} else {
		msg.Result = true
	}
	return c.JSON(http.StatusOK, msg)
}
