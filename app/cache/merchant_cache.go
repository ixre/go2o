/**
 * Copyright 2015 @ to2.net.
 * name : partner_cache
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package cache

import (
	"bytes"
	"context"
	"fmt"
	"go2o/core/domain/interface/merchant"
	"go2o/core/service"
	"go2o/core/service/impl"
	"go2o/core/service/proto"
	"sort"
	"strconv"
	"strings"
)

// 获取商户信息缓存
func GetValueMerchantCache(mchId int) *proto.SMerchant {
	var v proto.SMerchant
	var sto = GetKVS()
	var key = GetValueMerchantCacheCK(mchId)
	if sto.Get(key, &v) != nil {
		v2, _ := impl.MerchantService.GetMerchant(context.TODO(), &proto.Int64{Value: int64(mchId)})
		if v2 != nil {
			sto.SetExpire(key, *v2, DefaultMaxSeconds)
			return v2
		}
	}
	return &v

}

// 设置商户信息缓存
func GetValueMerchantCacheCK(mchId int) string {
	return fmt.Sprintf("cache:partner:value:%d", mchId)
}

// 设置商户站点配置
func GetMerchantSiteConfCK(mchId int) string {
	return fmt.Sprintf("cache:partner:siteconf:%d", mchId)
}

func DelMerchantCache(mchId int) {
	kvs := GetKVS()
	kvs.Del(GetValueMerchantCacheCK(mchId))
	kvs.Del(GetMerchantSiteConfCK(mchId))
}

// 根据主机头识别会员编号
func GetMerchantIdByHost(host string) int {
	key := "cache:host-for:" + host
	sto := GetKVS()
	id, err := sto.GetInt(key)
	mchId := id
	if err != nil || mchId <= 0 {
		mchId = int(impl.MerchantService.GetMerchantIdByHost(host))
		if mchId > 0 {
			sto.SetExpire(key, mchId, DefaultMaxSeconds)
		}
	}
	return mchId
}

// 根据API ID获取商户ID
func GetMerchantIdByApiId(apiId string) int64 {
	var mchId int64
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:id-%s", apiId)
	kvs.Get(key, &mchId)
	if mchId == 0 {
		mchId = impl.MerchantService.GetMerchantIdByApiId(apiId)
		if mchId != 0 {
			kvs.Set(key, mchId)
		}
	}
	return mchId
}

// 获取API 信息
func GetMerchantApiInfo(mchId int64) *merchant.ApiInfo {
	var d = new(merchant.ApiInfo)
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:info-%d", mchId)
	err := kvs.Get(key, &d)
	if err != nil {
		if d = impl.MerchantService.GetApiInfo(int(mchId)); d != nil {
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
	trans,cli,_ := service.ExpressServiceClient()
	defer trans.Close()
	list,_ := cli.GetProviders(context.TODO(),&proto.Empty{})
	iMap := make(map[string][]*proto.SExpressProvider, 0)
	var letArr []string
	for _, v := range list.Value {
		for _, g := range strings.Split(v.GroupFlag, ",") {
			if g == "" {
				continue
			}
			arr, ok := iMap[g]
			if !ok {
				arr = []*proto.SExpressProvider{}
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
