/**
 * Copyright 2015 @ to2.net.
 * name : shop_cache
 * author : jarryliu
 * date : 2016-06-03 17:22
 * description :
 * history :
 */
package cache

import (
	"fmt"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/util"
	"go2o/core/service/rsi"
)

// 设置商户站点配置
func GetShopDataKey(shopId int32) string {
	return fmt.Sprintf("go2o:cache:online-shop:siteconf:%d", shopId)
}

// 清除在线商店缓存
func CleanShopData(shopId int32) {
	if shopId > 0 {
		key := GetShopDataKey(shopId)
		GetKVS().Del(key)
	}
}

// 删除商铺缓存
func DelShopCache(mchId int32) {
	kvs := GetKVS()
	kvs.Del(GetValueMerchantCacheCK(mchId))
	kvs.Del(GetMerchantSiteConfCK(mchId))
}

// 根据主机头识别商店编号
func GetShopIdByHost(host string) int32 {
	//去除"m."
	//host = strings.Replace(host, variable.DOMAIN_PREFIX_MOBILE, "", -1)
	key := "go2o:cache:shop-host:" + host
	sto := GetKVS()
	shopId, err := util.I32Err(sto.GetInt(key))
	if err != nil || shopId <= 0 {
		_, shopId = rsi.ShopService.GetShopIdByHost(host)
		if shopId > 0 {
			sto.SetExpire(key, shopId, DefaultMaxSeconds)
		}
	}
	return shopId
}

// 根据商城编号获取商户编号
// todo: ?? int 和 int32
func GetMchIdByShopId(shopId int32) int32 {
	key := fmt.Sprintf("go2o:cache:mch-by-shop:%d", shopId)
	sto := GetKVS()
	tmpId, err := sto.GetInt(key)
	mchId := int32(tmpId)
	if err != nil || mchId <= 0 {
		mchId = rsi.ShopService.GetMerchantId(shopId)
		if mchId > 0 {
			sto.SetExpire(key, mchId, DefaultMaxSeconds)
		} else {
			log.Println("[ Shop][ Exception] - shop", shopId, " no merchant bind!")
		}
	}
	return mchId
}
