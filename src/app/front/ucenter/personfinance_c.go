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
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type personFinanceRiseC struct {
}

func (this *personFinanceRiseC) chkOpenState(memberId int, c *echox.Context) (
	bool, personfinance.RiseInfoValue, error) {
	e, err := dps.PersonFinanceService.GetRiseInfo(memberId)
	b := err == nil
	if !b { //如果未开通,则跳转到开通界面
		return false, e, c.Redirect(http.StatusFound, "openService?return_url="+
			c.Request().URL.String())
	}
	return true, e, nil
}

func (this *personFinanceRiseC) OpenService(c *echox.Context) error {
	if c.Request().Method == "POST" {
		return this.openService_post(c)
	}
	d := c.NewData()
	d.Map["AliasRise"] = variable.AliasRisePersonFinance
	return c.RenderOK("rise_openservice.html", d)
}

func (this *personFinanceRiseC) openService_post(c *echox.Context) error {
	msg := gof.Message{}
	memberId := GetSessionMemberId(c)
	if err := dps.PersonFinanceService.OpenRiseService(memberId); err != nil {
		if strings.Index(err.Error(), "exists") != -1 {
			msg.Message = "您已经开通服务!"
		} else {
			msg.Message = err.Error()
		}
	} else {
		msg.Result = true
	}
	return c.JSON(http.StatusOK, msg)
}

func (this *personFinanceRiseC) TransferIn(c *echox.Context) error {
	if c.Request().Method == "POST" {
		return this.transferIn_post(c)
	}
	memberId := GetSessionMemberId(c)
	if b, _, err := this.chkOpenState(memberId, c); !b {
		return err
	}
	d := c.NewData()
	d.Map["MinTranfer"] = format.FormatFloat(personfinance.RiseMinTransferInAmount)
	d.Map["RiseDate"] = time.Now().AddDate(0, 0, personfinance.RiseSettleTValue).Format("2006-01-02")
	d.Map["AliasPresentBalance"] = variable.AliasPresentAccount
	d.Map["AliasBalance"] = variable.AliasBalanceAccount
	d.Map["Account"] = dps.MemberService.GetAccount(memberId)
	return c.Render(http.StatusOK, "rise_transferin.html", d)
}

func (this *personFinanceRiseC) transferIn_post(c *echox.Context) error {
	msg := gof.Message{}
	memberId := GetSessionMemberId(c)
	amount, err := strconv.ParseFloat(c.Form("Amount"), 32)
	transferFrom, _ := strconv.Atoi(c.Form("TransferWith"))
	if err == nil {
		err = dps.PersonFinanceService.RiseTransferIn(memberId,
			personfinance.TransferWith(transferFrom), float32(amount))
	}
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	return c.JSON(http.StatusOK, msg)
}

func (this *personFinanceRiseC) TransferOut(c *echox.Context) error {
	if c.Request().Method == "POST" {
		return this.transferOut_post(c)
	}
	memberId := GetSessionMemberId(c)
	if b, e, err := this.chkOpenState(memberId, c); !b {
		return err
	} else {
		d := c.NewData()
		d.Map["AliasBalance"] = variable.AliasBalanceAccount
		d.Map["MinTranfer"] = format.FormatFloat(personfinance.RiseMinTransferOutAmount)
		d.Map["MaxTransferAmount"] = format.FormatFloat(e.Balance)
		return c.Render(http.StatusOK, "rise_transferout.html", d)
	}
}

func (this *personFinanceRiseC) transferOut_post(c *echox.Context) error {
	msg := gof.Message{}
	memberId := GetSessionMemberId(c)
	amount, err := strconv.ParseFloat(c.Form("Amount"), 32)
	transferTo, _ := strconv.Atoi(c.Form("TransferWith"))
	if err == nil {
		err = dps.PersonFinanceService.RiseTransferOut(memberId,
			personfinance.TransferWith(transferTo), float32(amount))
	}
	if err != nil {
		msg.Message = err.Error()
	} else {
		msg.Result = true
	}
	return c.JSON(http.StatusOK, msg)
}
