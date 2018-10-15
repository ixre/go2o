/**
 * Copyright 2015 @ z3q.net.
 * name : finance_c.go
 * author : jarryliu
 * date : 2016-01-07 21:36
 * description :
 * history :
 */
package partner

import (
	"github.com/jsix/gof"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
)

type financeC struct {
}

func (t *financeC) Balance_info(c *echox.Context) error {
	memberId, _ := strconv.Atoi(c.Query("member_id"))
	m := dps.MemberService.GetMember(memberId)
	d := c.NewData()
	d.Map = map[string]interface{}{
		"memberId": memberId,
		"member":   m,
	}
	return c.RenderOK("finance.balance_info.html", d)
}

func (t *financeC) New_balance_ticket(c *echox.Context) error {
	if c.Request().Method == "POST" {
		return t.new_balance_ticket_post(c)
	}
	memberId, _ := strconv.Atoi(c.Query("member_id"))
	m := dps.MemberService.GetMember(memberId)
	d := c.NewData()
	d.Map = map[string]interface{}{
		"member": m,
		"bpName": variable.AliasPresentAccount,
		"bpKt":   member.KindBalancePresent,
		"bKt":    member.KindBalanceCharge,
		"bName":  variable.AliasBalanceAccount,
	}
	return c.RenderOK("finance.member_newticket.html", d)
}

func (t *financeC) new_balance_ticket_post(c *echox.Context) error {
	var msg = gof.Result{ErrCode: 0}
	partnerId := getPartnerId(c)
	memberId, _ := strconv.Atoi(c.Form("member_id"))
	//kt := strings.Split(c.Form("kt"), "-")
	//if len(kt) < 2 {
	//	return c.JSON(http.StatusOK, gof.Result{Message: "参数错误"})
	//}
	//kind, _ := strconv.Atoi(kt[0])
	//ktype, _ := strconv.Atoi(kt[1])

	kind, _ := strconv.Atoi(c.Form("kt"))
	oper := c.Form("oper")
	remark := c.Form("remark")
	amtStr := c.Form("amount")
	amount, err := strconv.ParseFloat(oper+amtStr, 32)
	if err == nil {
		_, err = dps.MemberService.NewBalanceTicket(partnerId, memberId, kind, remark, float32(amount))
	}
	if err != nil {
		msg.ErrMsg = err.Error()
		msg.ErrCode = 1
	}
	return c.JSON(http.StatusOK, msg)
}
