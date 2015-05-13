/**
 * Copyright 2015 @ S1N1 Team.
 * name : partner_cache
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package cache

import (
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/storage"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
)

// 获取商户信息缓存
func GetValuePartnerCache(partnerId int) *partner.ValuePartner {
	var v *partner.ValuePartner
	var sto gof.Storage = GetKVS()
	var key string = fmt.Sprintf("cache:partner:value:%d", partnerId)

	if sto.Driver() == storage.DriveHashStorage {
		if obj := GetKVS().GetRaw(key); obj != nil {
			v = obj.(*partner.ValuePartner)
		}
	} else if sto.Driver() == storage.DriveRedisStorage {
		sto.Get(key, &v)
	}

	return v
}

// 设置商户信息缓存
func SetValuePartnerCache(partnerId int, v *partner.ValuePartner) {
	GetKVS().Set(fmt.Sprintf("cache:partner:value:%d", partnerId), v)
}

// 获取商户站点配置
func GetPartnerSiteConf(partnerId int) *partner.SiteConf {
	var v *partner.SiteConf
	var sto gof.Storage = GetKVS()
	var key string = fmt.Sprintf("cache:partner:siteconf:%d", partnerId)

	if sto.Driver() == storage.DriveHashStorage {
		if obj := GetKVS().GetRaw(key); obj != nil {
			v = obj.(*partner.SiteConf)
		}
	} else if sto.Driver() == storage.DriveRedisStorage {
		sto.Get(key, &v)
	}
	return v
}

// 设置商户站点配置
func SetPartnerSiteConf(partnerId int, v *partner.SiteConf) {
	GetKVS().Set(fmt.Sprintf("cache:partner:siteconf:%d", partnerId), v)
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
