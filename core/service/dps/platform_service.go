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
	_rep valueobject.IValueRep
}

func NewPlatformService(rep valueobject.IValueRep) *platformService {
	return &platformService{
		_rep: rep,
	}
}

// 获取微信接口配置
func (this *platformService) GetWxApiConfig() *valueobject.WxApiConfig {
	return this._rep.GetWxApiConfig()
}

// 保存微信接口配置
func (this *platformService) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	return this._rep.SaveWxApiConfig(v)
}

// 获取注册配置
func (this *platformService) GetRegisterPerm() *valueobject.RegisterPerm {
	return this._rep.GetRegisterPerm()
}

// 保存注册配置
func (this *platformService) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	return this._rep.SaveRegisterPerm(v)
}

// 获取全局系统数值设置
func (this *platformService) GetGlobNumberConf() *valueobject.GlobNumberConf {
	return this._rep.GetGlobNumberConf()
}

// 保存全局系统数值设置
func (this *platformService) SaveGlobNumberConf(v *valueobject.GlobNumberConf) error {
	return this._rep.SaveGlobNumberConf(v)
}

// 获取全局商户设置
func (this *platformService) GetGlobMchConf() *valueobject.GlobMchConf {
	return this._rep.GetGlobMchConf()
}

// 保存全局商户设置
func (this *platformService) SaveGlobMchConf(v *valueobject.GlobMchConf) error {
	return this._rep.SaveGlobMchConf(v)
}

// 获取全局商户销售设置
func (this *platformService) GetGlobMchSaleConf() *valueobject.GlobMchSaleConf {
	return this._rep.GetGlobMchSaleConf()
}

// 保存全局商户销售设置
func (this *platformService) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	return this._rep.SaveGlobMchSaleConf(v)
}
