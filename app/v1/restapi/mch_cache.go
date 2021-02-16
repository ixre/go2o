package restapi

import (
	"context"
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/storage"
	"go2o/core/service"
	"go2o/core/service/proto"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : mch_cache.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-30 11:46
 * description :
 * history :
 */
// 设置商户信息缓存
func getValueMerchantCacheCK(mchId int) string {
	return fmt.Sprintf("cache:partner:value:%d", mchId)
}

// 设置商户站点配置
func getMerchantSiteConfCK(mchId int) string {
	return fmt.Sprintf("cache:partner:siteconf:%d", mchId)
}

func DelMerchantCache(mchId int) {
	kvs := getKVS()
	kvs.Delete(getValueMerchantCacheCK(mchId))
	kvs.Delete(getMerchantSiteConfCK(mchId))
}

// 根据主机头识别会员编号
func GetMerchantIdByHost(host string) int {
	key := "cache:host-for:" + host
	sto := getKVS()
	id, err := sto.GetInt(key)
	mchId := id
	if err != nil || mchId <= 0 {
		trans, cli, _ := service.MerchantServiceClient()
		defer trans.Close()
		mchId, _ := cli.GetMerchantIdByHost(context.TODO(),
			&proto.String{Value: host})
		if mchId.Value > 0 {
			sto.SetExpire(key, mchId, 3600)
		}
	}
	return mchId
}

// 根据API ID获取商户ID
func GetMerchantIdByApiId(apiId string) int64 {
	var mchId int64
	kvs := getKVS()
	key := fmt.Sprintf("cache:partner:api:id-%s", apiId)
	kvs.Get(key, &mchId)
	if mchId == 0 {
		trans, cli, _ := service.MerchantServiceClient()
		defer trans.Close()
		mchId, _ := cli.GetMerchantIdByApiId(context.TODO(),
			&proto.String{Value: apiId})
		if mchId.Value != 0 {
			kvs.Set(key, mchId)
		}
	}
	return mchId
}

// 获取API 信息
func GetMerchantApiInfo(mchId int64) *proto.SMerchantApiInfo {
	var d = new(proto.SMerchantApiInfo)
	kvs := getKVS()
	key := fmt.Sprintf("cache:partner:api:info-%d", mchId)
	err := kvs.Get(key, &d)
	if err != nil {
		trans, cli, _ := service.MerchantServiceClient()
		defer trans.Close()
		ret, _ := cli.GetApiInfo(context.TODO(), &proto.MerchantId{Value: mchId})
		if ret != nil {
			kvs.Set(key, d)
		}
	}
	return d
}

func getKVS() storage.Interface {
	return gof.CurrentApp.Storage()
}
