/**
 * Copyright 2015 @ z3q.net.
 * name : json_c.go
 * author : jarryliu
 * date : 2016-04-25 23:09
 * description :
 * history :
 */
package ols

import (
	"github.com/jsix/gof"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"net/http"
	"strconv"
	"strings"
)

type jsonC struct {
}

// 广告
func (t *jsonC) Ad(ctx *echox.Context) error {
	names := strings.Split(ctx.Query("names"), "|")
	partnerId := GetPartnerId(ctx.Request(), ctx.Session)
	as := dps.AdvertisementService
	result := make(map[string]map[string]interface{}, len(names))
	for _, n := range names { //分别绑定广告
		ad, data := as.GetAdvertisementAndDataByName(partnerId, n)
		if ad == nil {
			result[n] = nil
			continue
		}
		result[n] = map[string]interface{}{
			"ad":   ad,
			"data": data,
		}
	}
	return ctx.JSON(http.StatusOK, result)
}

// 商品
func (t *jsonC) Simple_goods(ctx *echox.Context) error {
	types := strings.Split(ctx.Form("types"), "|")
	size, err := strconv.Atoi(ctx.Form("size"))
	if err != nil {
		msg := &gof.Message{}
		return ctx.JSON(http.StatusNotFound, msg.Error(err))
	}
	partnerId := GetPartnerId(ctx.Request(), ctx.Session)
	data := make(map[string]interface{}, len(types))
	ss := dps.SaleService
	for _, t := range types {
		switch t {
		case "new-goods":
			_, data[t] = ss.GetPagedOnShelvesGoods(partnerId,
				-1, 0, size, "gs_goods.id DESC")
		case "hot-sales":
			_, data[t] = ss.GetPagedOnShelvesGoods(partnerId,
				-1, 0, size, "gs_goods.sale_num DESC")
		}
	}
	return ctx.Debug(ctx.JSON(http.StatusOK, data))
}
