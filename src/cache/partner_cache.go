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
)

func GetValuePartnerCache(partnerId int)*partner.ValuePartner{
    var v *partner.ValuePartner
    GetKVS().Get(fmt.Sprintf("cache:partner:value:%d",partnerId),&v)
    return v
}

func SetValuePartnerCache(partnerId int,v *partner.ValuePartner){
    GetKVS().Set(fmt.Sprintf("cache:partner:value:%d",partnerId),v)
}