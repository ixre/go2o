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
	"errors"
	"github.com/jsix/gof/db"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/valueobject"
)

var _ valueobject.IValueRep = new(valueRep)

type valueRep struct {
	db.Connector
	_wxConf *valueobject.WxApiConfig
	_rpConf *valueobject.RegisterPerm
}

func NewValueRep(conn db.Connector) valueobject.IValueRep {
	return &valueRep{
		Connector: conn,
	}
}

// 获取微信接口配置
func (this *valueRep) GetWxApiConfig() *valueobject.WxApiConfig {
	if this._wxConf == nil {
		this._wxConf = &valueobject.WxApiConfig{}
		unMarshalFromFile("conf/core/wx_api", this._wxConf)
	}
	return this._wxConf
}

// 保存微信接口配置
func (this *valueRep) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	if v != nil {
		//todo: 检查证书文件是否存在
		this._wxConf = v
		return marshalToFile("conf/core/wx_api", this._wxConf)
	}
	return errors.New("nil value")
}

// 获取注册权限
func (this *valueRep) GetRegisterPerm() *valueobject.RegisterPerm {
	if this._rpConf == nil {
		this._rpConf = &valueobject.RegisterPerm{
			RegisterMode:        member.RegisterModeNormal,
			AnonymousRegistered: true,
		} // 默认值
		unMarshalFromFile("conf/core/wx_api", this._rpConf)
	}
	return this._rpConf
}

// 保存注册权限
func (this *valueRep) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	if v != nil {
		this._rpConf = v
		return marshalToFile("conf/core/wx_api", this._rpConf)
	}
	return errors.New("nil value")
}
