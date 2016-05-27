/**
 * Copyright 2015 @ z3q.net.
 * name : platform_service
 * author : jarryliu
 * date : 2016-05-27 15:30
 * description :
 * history :
 */
package dps

import (
    "go2o/core/domain/interface/valueobject"
)

// 平台服务
type platformService struct {
    _rep            valueobject.IValueRep
}

func NewPlatformService(rep valueobject.IValueRep) *platformService {
    return &platformService{
        _rep:            rep,
    }
}


// 获取微信接口配置
func (this *platformService) GetWxApiConfig()*valueobject.WxApiConfig{
    return this._rep.GetWxApiConfig()
}

// 保存微信接口配置
func (this *platformService) SaveWxApiConfig(v *valueobject.WxApiConfig)error{
    return this._rep.SaveWxApiConfig(v)
}
