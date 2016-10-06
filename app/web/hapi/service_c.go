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
	"github.com/jsix/gof"
	"go2o/core/service/dps"
	"go2o/x/echox"
	"net/http"
	"strconv"
)

type serviceC struct {
	gof.App
}

func (m *serviceC) Favorite(ctx *echox.Context) error {
	memberId := getMemberId(ctx)
	if memberId <= 0 {
		return requestLogin(ctx)
	}
	result := gof.Message{}

	favType := ctx.Query("type")
	id, _ := strconv.Atoi(ctx.Query("id"))
	if id <= 0 || favType == "" {
		result.Message = "收藏失败"
	} else {
		var err error
		ms := dps.MemberService
		if favType == "shop" {
			err = ms.FavoriteShop(memberId, id)
		} else {
			err = ms.FavoriteGoods(memberId, id)
		}
		result.Error(err)
	}
	return ctx.JSONP(http.StatusOK, ctx.Query("callback"), result)
}
