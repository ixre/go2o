/**
 * Copyright 2015 @ z3q.net.
 * name : person_finance.go
 * author : jarryliu
 * date : 2016-04-05 10:07
 * description :
 * history :
 */
package partner

import (
    "go2o/src/x/echox"
    "net/http"
)

type personFinanceC struct{
}

func (this *personFinanceC) Earnings_accounts(c *echox.Context)error{
    d := c.NewData()
    return c.Render(http.StatusOK,"pf.earnings_accounts",d)
}
