/**
 * Copyright 2015 @ z3q.net.
 * name : shop_cache
 * author : jarryliu
 * date : 2016-06-03 17:22
 * description :
 * history :
 */
package cache

import (
	"fmt"
	"github.com/jsix/gof/log"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/infrastructure/format"
	"go2o/core/service/rsi"
	"go2o/core/variable"
	"strings"
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
	host = strings.Replace(host, variable.DOMAIN_PREFIX_MOBILE, "", -1)
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

func getRdShopData(shopId int32) *shop.ShopDto {
	mchId := GetMchIdByShopId(shopId)
	v2 := rsi.ShopService.GetShopData(mchId, shopId)
	if v2 != nil {
		v3 := v2.Data.(shop.OnlineShop)
		v3.Logo = format.GetResUrl(v3.Logo)
		v2.Data = &v3
		if v2 != nil {
			//sto.SetExpire(key, *v2, DefaultMaxSeconds)
		}
		return v2
	}
	return v2
}

// 获取商城的数据
func GetOnlineShopData(shopId int32) *shop.ShopDto {
	//return getRdShopData(shopId)
	var v shop.ShopDto
	sto := GetKVS()
	key := GetShopDataKey(shopId)
	if err := sto.Get(key, &v); err != nil {
		mchId := GetMchIdByShopId(shopId)
		if v2 := rsi.ShopService.GetShopData(mchId, shopId); v2 != nil {
			v3 := v2.Data.(shop.OnlineShop)
			v3.Logo = format.GetResUrl(v3.Logo)
			v2.Data = &v3
			if v2 != nil {
				sto.SetExpire(key, *v2, DefaultMaxSeconds)
			}
			return v2
		}
	}
	return &v
}
