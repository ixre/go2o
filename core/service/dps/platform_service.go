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

// 获取数据存储
func (s *platformService) GetRegistry() valueobject.Registry {
	return s._rep.GetRegistry()
}

// 保存数据存储
func (s *platformService) SaveRegistry(v *valueobject.Registry) error {
	return s._rep.SaveRegistry(v)
}

// 获取模板配置
func (s *platformService) GetTemplateConf() valueobject.TemplateConf {
	return s._rep.GetTemplateConf()
}

// 保存模板配置
func (s *platformService) SaveTemplateConf(v *valueobject.TemplateConf) error {
	return s._rep.SaveTemplateConf(v)
}

// 获取移动应用设置
func (p *platformService) GetMoAppConf() valueobject.MoAppConf {
	return p._rep.GetMoAppConf()
}

// 保存移动应用设置
func (p *platformService) SaveMoAppConf(r *valueobject.MoAppConf) error {
	return p._rep.SaveMoAppConf(r)
}

// 获取微信接口配置
func (p *platformService) GetWxApiConfig() valueobject.WxApiConfig {
	return p._rep.GetWxApiConfig()
}

// 保存微信接口配置
func (p *platformService) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	return p._rep.SaveWxApiConfig(v)
}

// 获取注册配置
func (p *platformService) GetRegisterPerm() valueobject.RegisterPerm {
	return p._rep.GetRegisterPerm()
}

// 保存注册配置
func (p *platformService) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	return p._rep.SaveRegisterPerm(v)
}

// 获取全局系统数值设置
func (p *platformService) GetGlobNumberConf() valueobject.GlobNumberConf {
	return p._rep.GetGlobNumberConf()
}

// 保存全局系统数值设置
func (p *platformService) SaveGlobNumberConf(v *valueobject.GlobNumberConf) error {
	return p._rep.SaveGlobNumberConf(v)
}

// 获取平台设置
func (p *platformService) GetPlatformConf() valueobject.PlatformConf {
	return p._rep.GetPlatformConf()
}

// 保存平台设置
func (p *platformService) SavePlatformConf(v *valueobject.PlatformConf) error {
	return p._rep.SavePlatformConf(v)
}

// 获取全局商户销售设置
func (p *platformService) GetGlobMchSaleConf() valueobject.GlobMchSaleConf {
	return p._rep.GetGlobMchSaleConf()
}

// 保存全局商户销售设置
func (p *platformService) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	return p._rep.SaveGlobMchSaleConf(v)
}

// 获取短信设置
func (p *platformService) GetSmsApiSet() valueobject.SmsApiSet {
	return p._rep.GetSmsApiSet()
}

// 保存短信API
func (p *platformService) SaveSmsApiPerm(provider int, s *valueobject.SmsApiPerm) error {
	return p._rep.SaveSmsApiPerm(provider, s)
}

// 获取默认的短信API
func (p *platformService) GetDefaultSmsApiPerm() (int, *valueobject.SmsApiPerm) {
	return p._rep.GetDefaultSmsApiPerm()
}

// 获取下级区域
func (p *platformService) GetChildAreas(id int) []*valueobject.Area {
	return p._rep.GetChildAreas(id)
}

// 获取地区名称
func (p *platformService) GetAreaNames(id []int64) []string {
	return p._rep.GetAreaNames(id)
}

// 获取省市区字符串
func (p *platformService) GetAreaString(province, city, district int64) string {
	if province == 0 || city == 0 || district == 0 {
		return ""
	}
	return p._rep.GetAreaString(province, city, district)
}
