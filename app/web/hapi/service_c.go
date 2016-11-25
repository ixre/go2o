/**
 * Copyright 2015 @ z3q.net.
 * name : service_c.go
 * author : jarryliu
 * date : 2016-09-09 23:45
 * description :
 * history :
 */
package hapi

import (
	"github.com/jsix/goex/echox"
	"github.com/jsix/gof"
	"github.com/jsix/gof/util"
	"go2o/core/service/rsi"
	"net/http"
	"strconv"
)

type serviceC struct {
	gof.App
}

func (m *serviceC) Favorite(c *echox.Context) error {
	memberId := getMemberId(c)
	if memberId <= 0 {
		return requestLogin(c)
	}
	result := gof.Message{}

	favType := c.QueryParam("type")
	id, _ := util.I32Err(strconv.Atoi(c.QueryParam("id")))
	if id <= 0 || favType == "" {
		result.Message = "收藏失败"
	} else {
		var err error
		ms := rsi.MemberService
		if favType == "shop" {
			err = ms.FavoriteShop(memberId, id)
		} else {
			err = ms.FavoriteGoods(memberId, id)
		}
		result.Error(err)
	}
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), result)
}
