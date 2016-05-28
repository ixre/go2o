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

// 获取全局系统销售设置
func (this *platformService) GetGlobSaleConf() *valueobject.GlobSaleConf {
	return this._rep.GetGlobSaleConf()
}

// 保存全局系统销售设置
func (this *platformService) SaveGlobSaleConf(v *valueobject.GlobSaleConf) error {
	return this._rep.SaveGlobSaleConf(v)
}

// 获取全局商户销售设置
func (this *platformService) GetGlobMerchantSaleConf() *valueobject.GlobMerchantSaleConf {
	return this._rep.GetGlobMerchantSaleConf()
}

// 保存全局商户销售设置
func (this *platformService) SaveGlobMerchantSaleConf(v *valueobject.GlobMerchantSaleConf) error {
	return this._rep.SaveGlobMerchantSaleConf(v)
}
