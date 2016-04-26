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
	"encoding/gob"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/crypto"
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/service/dps"
	"go2o/src/x/echox"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	//todo: ??? 设置为可配置
	maxSeconds int64 = 10
)

func init() {
	gob.Register(map[string]map[string]interface{}{})
	gob.Register(ad.ValueGallery{})
	gob.Register(ad.ValueAdvertisement{})
	gob.Register([]*valueobject.Goods{})
	gob.Register(valueobject.Goods{})
}

type jsonC struct {
}

func getMd5(s string) string {
	return crypto.Md5([]byte(s))[8:16]
}

// 广告
func (t *jsonC) Ad(ctx *echox.Context) error {
	namesParams := strings.TrimSpace(ctx.Query("names"))
	names := strings.Split(namesParams, "|")
	partnerId := GetPartnerId(ctx)
	as := dps.AdvertisementService

	result := make(map[string]map[string]interface{}, len(names))
	key := fmt.Sprint("go2o:front:ad:%d-%s", partnerId, getMd5(namesParams))
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil { //从缓存中读取
		log.Println(err)
		for _, n := range names {
			//分别绑定广告
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
		sto.SetExpire(key, result, maxSeconds)
	}
	return ctx.JSON(http.StatusOK, result)
}

// 商品
func (t *jsonC) Simple_goods(ctx *echox.Context) error {
	typeParams := strings.TrimSpace(ctx.Form("types"))
	types := strings.Split(typeParams, "|")
	size, err := strconv.Atoi(ctx.Form("size"))
	if err != nil {
		msg := &gof.Message{}
		return ctx.JSON(http.StatusNotFound, msg.Error(err))
	}
	partnerId := GetPartnerId(ctx)
	result := make(map[string]interface{}, len(types))

	key := fmt.Sprint("go2o:front:sg:%d-%d-%s", partnerId, size, getMd5(typeParams))
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil { //从缓存中读取
		log.Println(err)
		ss := dps.SaleService
		for _, t := range types {
			switch t {
			case "new-goods":
				_, result[t] = ss.GetPagedOnShelvesGoods(partnerId,
					-1, 0, size, "gs_goods.id DESC")
			case "hot-sales":
				_, result[t] = ss.GetPagedOnShelvesGoods(partnerId,
					-1, 0, size, "gs_goods.sale_num DESC")
			}
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return ctx.Debug(ctx.JSON(http.StatusOK, result))
}

// 获取销售标签获取商品
func (t *jsonC) Saletag_goods(ctx *echox.Context) error {
	codeParams := strings.TrimSpace(ctx.Form("codes"))
	codes := strings.Split(codeParams, "|")
	begin, _ := strconv.Atoi(ctx.Form("begin"))
	size, err := strconv.Atoi(ctx.Form("size"))
	if err != nil {
		msg := &gof.Message{}
		return ctx.JSON(http.StatusNotFound, msg.Error(err))
	}
	partnerId := GetPartnerId(ctx)
	result := make(map[string]interface{}, len(codes))

	key := fmt.Sprint("go2o:front:stg:%d-%d-%d-%s", partnerId, begin, size, getMd5(codeParams))
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil { //从缓存中读取
		log.Println(err)
		for _, code := range codes {
			list := dps.SaleService.GetValueGoodsBySaleTag(
				partnerId, code, "", begin, begin+size)
			result[code] = list
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return ctx.Debug(ctx.JSON(http.StatusOK, result))
}
