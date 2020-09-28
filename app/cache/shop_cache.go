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
	"context"
	"fmt"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/util"
	"go2o/core/service"
	"go2o/core/service/impl"
	"go2o/core/service/proto"
)

// 设置商户站点配置
func GetShopDataKey(shopId int64) string {
	return fmt.Sprintf("go2o:cache:online-shop:siteconf:%d", shopId)
}

// 清除在线商店缓存
func CleanShopData(shopId int64) {
	if shopId > 0 {
		key := GetShopDataKey(shopId)
		GetKVS().Del(key)
	}
}

// 删除店铺缓存
func DelShopCache(mchId int) {
	kvs := GetKVS()
	kvs.Del(GetValueMerchantCacheCK(mchId))
	kvs.Del(GetMerchantSiteConfCK(mchId))
}

// 根据主机头识别商店编号
func GetShopIdByHost(host string) int64 {
	//去除"m."
	//host = strings.Replace(host, variable.DOMAIN_PREFIX_MOBILE, "", -1)
	key := "go2o:cache:shop-host:" + host
	sto := GetKVS()
	shopId, err := util.I64Err(sto.GetInt(key))
	if err != nil || shopId <= 0 {
		trans, cli, err := service.ShopServiceClient()
		if err == nil {
			defer trans.Close()
			v, _ := cli.QueryShopByHost(context.TODO(), &proto.String{Value: host})
			shopId = v.Value
			if shopId > 0 {
				_ = sto.SetExpire(key, shopId, DefaultMaxSeconds)
			}
		}
	}
	return shopId
}

// 根据商城编号获取商户编号
// todo: ?? int 和 int32
func GetMchIdByShopId(shopId int64) int64 {
	key := fmt.Sprintf("go2o:cache:mch-by-shop:%d", shopId)
	sto := GetKVS()
	tmpId, err := sto.GetInt(key)
	mchId := int64(tmpId)
	trans,cli,_ := service.ShopServiceClient()
	defer trans.Close()
	if err != nil || mchId <= 0 {
		mchId = cli.GetMerchantId(shopId)
		if mchId > 0 {
			sto.SetExpire(key, mchId, DefaultMaxSeconds)
		} else {
			log.Println("[ Shop][ Exception] - shop", shopId, " no merchant bind!")
		}
	}
	return mchId
}
