/**
 * Copyright 2015 @ z3q.net.
 * name : person_finance.go
 * author : jarryliu
 * date : 2016-04-05 10:07
 * description :
 * history :
 */
package merchant

import (
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
)

type personFinanceC struct {
}

func (this *personFinanceC) Earnings_account(c *echox.Context) error {
	d := c.NewData()
	d.Map["GrowAlias"] = variable.AliasGrowAccount
	return c.Render(http.StatusOK, "pf.earnings_account.html", d)
}

func (this *personFinanceC) Earnings_log(c *echox.Context) error {
	d := c.NewData()

	personId, err := strconv.Atoi(c.Query("person_id"))
	if err == nil {
		if m := dps.MemberService.GetMember(personId); m != nil {
			d.Map["PersonName"] = m.Name
		}
	}
	d.Map["PersonId"] = personId
	d.Map["GrowAlias"] = variable.AliasGrowAccount
	return c.Render(http.StatusOK, "pf.earnings_log.html", d)
}

func (this *personFinanceC) Earnings_transferlog(c *echox.Context) error {
	d := c.NewData()
	d.Map["GrowAlias"] = variable.AliasGrowAccount
	return c.Render(http.StatusOK, "pf.earnings_transferlog.html", d)
}
