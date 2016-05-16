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
	"fmt"
	"github.com/jsix/gof"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
)

// 获取商户信息缓存
func GetValuePartnerCache(merchantId int) *merchant.MerchantValue {
	var v merchant.MerchantValue
	var sto gof.Storage = GetKVS()
	var key string = GetValuePartnerCacheCK(merchantId)
	if sto.Get(key, &v) != nil {
		v2, err := dps.PartnerService.GetMerchant(merchantId)
		if v2 != nil && err == nil {
			sto.SetExpire(key, *v2, DefaultMaxSeconds)
			return v2
		}
	}
	return &v

}

// 设置商户信息缓存
func GetValuePartnerCacheCK(merchantId int) string {
	return fmt.Sprintf("cache:partner:value:%d", merchantId)
}

// 设置商户站点配置
func GetPartnerSiteConfCK(merchantId int) string {
	return fmt.Sprintf("cache:partner:siteconf:%d", merchantId)
}

func DelPartnerCache(merchantId int) {
	kvs := GetKVS()
	kvs.Del(GetValuePartnerCacheCK(merchantId))
	kvs.Del(GetPartnerSiteConfCK(merchantId))
}

// 获取商户站点配置
func GetPartnerSiteConf(merchantId int) *merchant.SiteConf {
	var v merchant.SiteConf
	var sto gof.Storage = GetKVS()
	var key string = GetPartnerSiteConfCK(merchantId)
	if sto.Get(key, &v) != nil {
		v2 := dps.PartnerService.GetSiteConf(merchantId)
		v2.Logo = format.GetResUrl(v2.Logo)
		if v2 != nil {
			sto.SetExpire(key, *v2, DefaultMaxSeconds)
		}
		return v2
	}
	return &v
}

// 根据主机头识别会员编号
func GetMerchantIdByHost(host string) int {
	merchantId := 0
	key := "cache:host-for:" + host
	sto := GetKVS()
	var err error
	if merchantId, err = sto.GetInt(key); err != nil || merchantId <= 0 {
		merchantId = dps.PartnerService.GetMerchantIdByHost(host)
		if merchantId > 0 {
			sto.SetExpire(key, merchantId, DefaultMaxSeconds)
		}
	}
	return merchantId
}

// 根据API ID获取商户ID
func GetMerchantIdByApiId(apiId string) int {
	var merchantId int
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:id-%s", apiId)
	kvs.Get(key, &merchantId)
	if merchantId == 0 {
		merchantId = dps.PartnerService.GetMerchantIdByApiId(apiId)
		if merchantId != 0 {
			kvs.Set(key, merchantId)
		}
	}
	return merchantId
}

// 获取API 信息
func GetPartnerApiInfo(merchantId int) *merchant.ApiInfo {
	var d *merchant.ApiInfo = new(merchant.ApiInfo)
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:info-%d", merchantId)
	err := kvs.Get(key, &d)
	if err != nil {
		if d = dps.PartnerService.GetApiInfo(merchantId); d != nil {
			kvs.Set(key, d)
		}
	}
	return d
}
