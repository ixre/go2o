/**
 * Copyright 2015 @ z3q.net.
 * name : partner_cache
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package cache

import (
	"bytes"
	"fmt"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/merchant"
	"go2o/core/service/rsi"
	"sort"
	"strconv"
	"strings"
)

// 获取商户信息缓存
func GetValueMerchantCache(mchId int32) *merchant.Merchant {
	var v merchant.Merchant
	var sto storage.Interface = GetKVS()
	var key string = GetValueMerchantCacheCK(mchId)
	if sto.Get(key, &v) != nil {
		v2 := rsi.MerchantService.GetMerchant(mchId)
		if v2 != nil {
			sto.SetExpire(key, *v2, DefaultMaxSeconds)
			return v2
		}
	}
	return &v

}

// 设置商户信息缓存
func GetValueMerchantCacheCK(mchId int32) string {
	return fmt.Sprintf("cache:partner:value:%d", mchId)
}

// 设置商户站点配置
func GetMerchantSiteConfCK(mchId int32) string {
	return fmt.Sprintf("cache:partner:siteconf:%d", mchId)
}

func DelMerchantCache(mchId int32) {
	kvs := GetKVS()
	kvs.Del(GetValueMerchantCacheCK(mchId))
	kvs.Del(GetMerchantSiteConfCK(mchId))
}

// 根据主机头识别会员编号
func GetMerchantIdByHost(host string) int32 {
	key := "cache:host-for:" + host
	sto := GetKVS()
	id, err := sto.GetInt(key)
	mchId := int32(id)
	if err != nil || mchId <= 0 {
		mchId = rsi.MerchantService.GetMerchantIdByHost(host)
		if mchId > 0 {
			sto.SetExpire(key, mchId, DefaultMaxSeconds)
		}
	}
	return mchId
}

// 根据API ID获取商户ID
func GetMerchantIdByApiId(apiId string) int32 {
	var mchId int32
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:id-%s", apiId)
	kvs.Get(key, &mchId)
	if mchId == 0 {
		mchId = rsi.MerchantService.GetMerchantIdByApiId(apiId)
		if mchId != 0 {
			kvs.Set(key, mchId)
		}
	}
	return mchId
}

// 获取API 信息
func GetMerchantApiInfo(mchId int32) *merchant.ApiInfo {
	var d *merchant.ApiInfo = new(merchant.ApiInfo)
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:info-%d", mchId)
	err := kvs.Get(key, &d)
	if err != nil {
		if d = rsi.MerchantService.GetApiInfo(mchId); d != nil {
			kvs.Set(key, d)
		}
	}
	return d
}

var (
	expressCacheKey = "go2o:repo:express:ship-tab"
)

// 获取发货的快递选项卡
func GetShipExpressTab() string {
	sto := GetKVS()
	html, err := sto.GetString(expressCacheKey)
	if err != nil {
		html = getRealShipExpressTab()
		sto.Set(expressCacheKey, html)
	}
	return html
}

func getRealShipExpressTab() string {
	buf := bytes.NewBuffer(nil)
	list := rsi.ExpressService.GetEnabledProviders()
	iMap := make(map[string][]*express.ExpressProvider, 0)
	letArr := []string{}
	for _, v := range list {
		for _, g := range strings.Split(v.GroupFlag, ",") {
			if g == "" {
				continue
			}
			arr, ok := iMap[g]
			if !ok {
				arr = []*express.ExpressProvider{}
				letArr = append(letArr, g)
			}
			arr = append(arr, v)
			iMap[g] = arr
		}

	}
	sort.Strings(letArr)
	l := len(letArr)
	if letArr[l-1] == "常用" {
		letArr = append(letArr[l-1:], letArr[:l-1]...)
	}

	buf.WriteString(`<div class="gra-tabs" id="express-tab"><ul class="tabs">`)
	for _, v := range letArr {
		buf.WriteString(`<li title="`)
		buf.WriteString(v)
		buf.WriteString(`" href="`)
		buf.WriteString(v)
		buf.WriteString(`"><span class="tab-title">`)
		buf.WriteString(v)
		buf.WriteString(`</span></li>`)
	}
	buf.WriteString("</ul>")
	buf.WriteString(`<div class="frames">`)
	i := 0
	for _, l := range letArr {
		buf.WriteString(`<div class="frame"><ul class="list">`)
		for _, v := range iMap[l] {
			i++
			buf.WriteString("<li><input type=\"radio\" name=\"ProviderId\" field=\"ProviderId\" value=\"")
			buf.WriteString(strconv.Itoa(int(v.Id)))
			buf.WriteString(`" id="provider_`)
			buf.WriteString(strconv.Itoa(i))
			buf.WriteString(`"/><label for="provider_`)
			buf.WriteString(strconv.Itoa(i))
			buf.WriteString(`">`)
			buf.WriteString(v.Name)
			buf.WriteString("</label></li>")
		}
		buf.WriteString("</ul></div>")
	}
	buf.WriteString("</div></div>")
	return buf.String()
}
