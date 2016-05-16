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
	partnerId := GetMerchantId(ctx)
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

func (t *jsonC) getMultiParams(s string) (p string, size, begin int) {
	arr := strings.Split(s, "*")
	l := len(arr)
	if l == 1 {
		p = s
		size = 10 //返回默认10条
	} else {
		p = arr[0]
		size, _ = strconv.Atoi(arr[1])
		if l > 2 {
			begin, _ = strconv.Atoi(arr[2])
		}
	}
	return p, size, begin
}

// 商品
func (this *jsonC) Simple_goods(ctx *echox.Context) error {
	typeParams := strings.TrimSpace(ctx.Form("params"))
	types := strings.Split(typeParams, "|")
	partnerId := GetMerchantId(ctx)
	result := make(map[string]interface{}, len(types))

	key := fmt.Sprint("go2o:front:sg:%d-%s", partnerId, getMd5(typeParams))
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil {
		//从缓存中读取
		log.Println(err)
		ss := dps.SaleService
		for _, t := range types {
			p, size, begin := this.getMultiParams(t)
			switch p {
			case "new-goods":
				_, result[p] = ss.GetPagedOnShelvesGoods(partnerId,
					-1, begin, begin+size, "gs_goods.id DESC")
			case "hot-sales":
				_, result[p] = ss.GetPagedOnShelvesGoods(partnerId,
					-1, begin, begin+size, "gs_goods.sale_num DESC")
			}
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return ctx.Debug(ctx.JSON(http.StatusOK, result))
}

// 获取销售标签获取商品
func (this *jsonC) Saletag_goods(ctx *echox.Context) error {
	codeParams := strings.TrimSpace(ctx.Form("params"))
	codes := strings.Split(codeParams, "|")
	partnerId := GetMerchantId(ctx)
	result := make(map[string]interface{}, len(codes))

	key := fmt.Sprint("go2o:front:stg:%d--%s", partnerId, getMd5(codeParams))
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil { //从缓存中读取
		log.Println(err)
		for _, param := range codes {
			code, size, begin := this.getMultiParams(param)
			list := dps.SaleService.GetValueGoodsBySaleTag(
				partnerId, code, "", begin, begin+size)
			result[code] = list
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return ctx.Debug(ctx.JSON(http.StatusOK, result))
}
