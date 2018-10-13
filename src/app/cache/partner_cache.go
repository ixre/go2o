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
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/format"
	"go2o/src/core/service/dps"
)

// 获取商户信息缓存
func GetValuePartnerCache(partnerId int) *partner.ValuePartner {
	var v partner.ValuePartner
	var sto storage.Interface = GetKVS()
	var key string = GetValuePartnerCacheCK(partnerId)
	if sto.Get(key, &v) != nil {
		v2, err := dps.PartnerService.GetPartner(partnerId)
		if v2 != nil && err == nil {
			sto.SetExpire(key, *v2, DefaultMaxSeconds)
			return v2
		}
	}
	return &v

}

// 设置商户信息缓存
func GetValuePartnerCacheCK(partnerId int) string {
	return fmt.Sprintf("cache:partner:value:%d", partnerId)
}

// 设置商户站点配置
func GetPartnerSiteConfCK(partnerId int) string {
	return fmt.Sprintf("cache:partner:siteconf:%d", partnerId)
}

func DelPartnerCache(partnerId int) {
	kvs := GetKVS()
	kvs.Del(GetValuePartnerCacheCK(partnerId))
	kvs.Del(GetPartnerSiteConfCK(partnerId))
}

// 获取商户站点配置
func GetPartnerSiteConf(partnerId int) *partner.SiteConf {
	var v partner.SiteConf
	var sto storage.Interface = GetKVS()
	var key string = GetPartnerSiteConfCK(partnerId)
	if sto.Get(key, &v) != nil {
		v2 := dps.PartnerService.GetSiteConf(partnerId)
		v2.Logo = format.GetResUrl(v2.Logo)
		if v2 != nil {
			sto.SetExpire(key, *v2, DefaultMaxSeconds)
		}
		return v2
	}
	return &v
}

// 根据主机头识别会员编号
func GetPartnerIdByHost(host string) int {
	partnerId := 0
	key := "cache:host-for:" + host
	sto := GetKVS()
	var err error
	if partnerId, err = sto.GetInt(key); err != nil || partnerId <= 0 {
		partnerId = dps.PartnerService.GetPartnerIdByHost(host)
		if partnerId > 0 {
			sto.SetExpire(key, partnerId, DefaultMaxSeconds)
		}
	}
	return partnerId
}

// 根据API ID获取商户ID
func GetPartnerIdByApiId(apiId string) int {
	var partnerId int
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:id-%s", apiId)
	kvs.Get(key, &partnerId)
	if partnerId == 0 {
		partnerId = dps.PartnerService.GetPartnerIdByApiId(apiId)
		if partnerId != 0 {
			kvs.Set(key, partnerId)
		}
	}
	return partnerId
}

// 获取API 信息
func GetPartnerApiInfo(partnerId int) *partner.ApiInfo {
	var d *partner.ApiInfo = new(partner.ApiInfo)
	kvs := GetKVS()
	key := fmt.Sprintf("cache:partner:api:info-%d", partnerId)
	err := kvs.Get(key, &d)
	if err != nil {
		if d = dps.PartnerService.GetApiInfo(partnerId); d != nil {
			kvs.Set(key, d)
		}
	}
	return d
}
