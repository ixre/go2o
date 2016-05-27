/**
 * Copyright 2015 @ z3q.net.
 * name : value_rep
 * author : jarryliu
 * date : 2016-05-27 15:32
 * description :
 * history :
 */
package repository

import (
    "github.com/jsix/gof/db"
    "go2o/core/domain/interface/valueobject"
    "errors"
)

var _ valueobject.IValueRep = new(valueRep)

type valueRep struct {
    db.Connector
    _wxConfig *valueobject.WxApiConfig
}

func NewValueRep(conn db.Connector)valueobject.IValueRep{
    return &valueRep{
        Connector:conn,
    }
}

// 获取微信接口配置
func (this *valueRep) GetWxApiConfig() *valueobject.WxApiConfig {
    if this._wxConfig == nil {
        this._wxConfig = &valueobject.WxApiConfig{}
        unMarshalFromFile("conf/core/wx_api", this._wxConfig)
    }
    return this._wxConfig
}

// 保存微信接口配置
func (this *valueRep) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
    if v != nil{
        //todo: 检查证书文件是否存在
        this._wxConfig = v
        marshalToFile("conf/core/wx_api",this._wxConfig)
    }
    return errors.New("nil value")
}





