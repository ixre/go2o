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
    "go2o/src/core/domain/interface/partner"
    "fmt"
    "github.com/atnet/gof"
    "github.com/atnet/gof/storage"
)

func GetValuePartnerCache(partnerId int)*partner.ValuePartner{
    var v *partner.ValuePartner;
    var sto gof.Storage = GetKVS()
    var key string = fmt.Sprintf("cache:partner:value:%d", partnerId)

    if sto.Driver() == storage.DriveHashStorage {
        if obj :=  GetKVS().GetRaw(key);obj != nil {
            v = obj.(*partner.ValuePartner)
        }
    }else if(sto.Driver() == storage.DriveRedisStorage){
        sto.Get(key,&v)
    }

    return v
}

func SetValuePartnerCache(partnerId int,v *partner.ValuePartner){
    GetKVS().Set(fmt.Sprintf("cache:partner:value:%d",partnerId),v)
}

func GetPartnerSiteConf(partnerId int)*partner.SiteConf{
    var v *partner.SiteConf;
    var sto gof.Storage = GetKVS()
    var key string = fmt.Sprintf("cache:partner:siteconf:%d", partnerId)

    if sto.Driver() == storage.DriveHashStorage {
        if obj :=  GetKVS().GetRaw(key);obj != nil {
            v = obj.(*partner.SiteConf)
        }
    }else if(sto.Driver() == storage.DriveRedisStorage) {
        sto.Get(key, &v)
    }
    return v
}

func SetPartnerSiteConf(partnerId int,v *partner.SiteConf){
    GetKVS().Set(fmt.Sprintf("cache:partner:siteconf:%d",partnerId),v)
}