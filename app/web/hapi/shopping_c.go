package hapi

import (
	"github.com/jsix/goex/echox"
	"github.com/jsix/gof"
	"go2o/core/service/rsi"
	"net/http"
)

type shoppingC struct {
	app gof.App
}

func (s *shoppingC) preReq(c *echox.Context) error {
	//if getMemberId(c)
	return nil
}

// 收货地址列表
func (s *shoppingC) AddressList(c *echox.Context) error {
	memberId := getMemberId(c)
	if memberId <= 0 {
		return requestLogin(c)
	}
	address := rsi.MemberService.GetAddressList(memberId)
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), address)
}
