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
	"github.com/jsix/gof/storage"
	"go2o/core/domain/interface/merchant"
	"go2o/core/service/dps"
)

// 获取商户信息缓存
func GetValueMerchantCache(merchantId int) *merchant.Merchant {
	var v merchant.Merchant
	var sto storage.Interface = GetKVS()
	var key string = GetValueMerchantCacheCK(merchantId)
	if sto.Get(key, &v) != nil {
		v2, err := dps.MerchantService.GetMerchant(merchantId)
		if v2 != nil && err == nil {
			sto.SetExpire(key, *v2, DefaultMaxSeconds)
			return v2
		}
	}
	return &v

}

// 设置商户信息缓存
func GetValueMerchantCacheCK(merchantId int) string {
	return fmt.Sprintf("cache:partner:value:%d", merchantId)
}

// 设置商户站点配置
func GetMerchantSiteConfCK(merchantId int) string {
	return fmt.Sprintf("cache:partner:siteconf:%d", merchantId)
}

func DelMerchantCache(merchantId int) {
	kvs := GetKVS()
	kvs.Del(GetValueMerchantCacheCK(merchantId))
	kvs.Del(GetMerchantSiteConfCK(merchantId))
}

// 根据主机头识别会员编号
func GetMerchantIdByHost(host string) int {
	merchantId := 0
	key := "cache:host-for:" + host
	sto := GetKVS()
	var err error
	if merchantId, err = sto.GetInt(key); err != nil || merchantId <= 0 {
		merchantId = dps.MerchantService.GetMerchantIdByHost(host)
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
		merchantId = dps.MerchantService.GetMerchantIdByApiId(apiId)
		if merchantId != 0 {
			kvs.Set(key, merchantId)
		}
	}
	return merchantId
}

// 获取API 信息
func GetMerchantApiInfo(merchantId int) *merchant.ApiInfo {
	var d *merchant.ApiInfo = new(merchant.ApiInfo)
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:info-%d", merchantId)
	err := kvs.Get(key, &d)
	if err != nil {
		if d = dps.MerchantService.GetApiInfo(merchantId); d != nil {
			kvs.Set(key, d)
		}
	}
	return d
}
