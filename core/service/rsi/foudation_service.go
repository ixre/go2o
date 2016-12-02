/**
 * Copyright 2015 @ z3q.net.
 * name : platform_service
 * author : jarryliu
 * date : 2016-05-27 15:30
 * description :
 * history :
 */
package rsi

import (
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/module"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/service/thrift/parser"
	"go2o/core/variable"
)

// 基础服务
type foundationService struct {
	_rep valueobject.IValueRep
}

func NewFoundationService(rep valueobject.IValueRep) *foundationService {
	return &foundationService{
		_rep: rep,
	}
}

// 验证超级用户账号和密码
func (s *foundationService) ValidateSuper(user string, pwd string) (r bool, err error) {
	superPwd := gof.CurrentApp.Config().Get("super_login_md5")
	encPwd := domain.Md5Pwd(pwd, user)
	return superPwd == encPwd, nil
}

// 保存超级用户账号和密码
func (s *foundationService) FlushSuperPwd(user string, pwd string) (err error) {
	conf := gof.CurrentApp.Config()
	encPwd := domain.Md5Pwd(pwd, user)
	conf.Set("super_login_md5", encPwd)
	//conf.Flush()
	return errors.New("暂不支持保存")
}

// 注册单点登录应用,返回值：
//   -  1. 成功，并返回token
//   - -1. 接口地址不正确
//   - -2. 已经注册
func (s *foundationService) RegisterSsoApp(app *define.SsoApp) (r string, err error) {
	sso := module.Get(module.M_SSO).(*module.SSOModule)
	token, err := sso.Register(app)
	if err == nil {
		return "1:" + token, nil
	}
	return err.Error(), nil
}

// 获取单点登录应用
func (s *foundationService) GetAllSsoApp() (r []string, err error) {
	sso := module.Get(module.M_SSO).(*module.SSOModule)
	return sso.Array(), nil
}

// 创建同步登录的地址
func (s *foundationService) GetSyncLoginUrl(returnUrl string) (r string, err error) {
	return fmt.Sprintf("%s://%s%s/auth?return_url=%s",
		variable.DOMAIN_PASSPORT_PROTO, variable.DOMAIN_PREFIX_PASSPORT,
		variable.Domain, returnUrl), nil
}

// 获取数据存储
func (s *foundationService) GetRegistry() valueobject.Registry {
	return s._rep.GetRegistry()
}

// 保存数据存储
func (s *foundationService) SaveRegistry(v *valueobject.Registry) error {
	return s._rep.SaveRegistry(v)
}

// 获取模板配置
func (s *foundationService) GetTemplateConf() valueobject.TemplateConf {
	return s._rep.GetTemplateConf()
}

// 保存模板配置
func (s *foundationService) SaveTemplateConf(v *valueobject.TemplateConf) error {
	return s._rep.SaveTemplateConf(v)
}

// 获取移动应用设置
func (p *foundationService) GetMoAppConf() valueobject.MoAppConf {
	return p._rep.GetMoAppConf()
}

// 保存移动应用设置
func (p *foundationService) SaveMoAppConf(r *valueobject.MoAppConf) error {
	return p._rep.SaveMoAppConf(r)
}

// 获取微信接口配置
func (p *foundationService) GetWxApiConfig() valueobject.WxApiConfig {
	return p._rep.GetWxApiConfig()
}

// 保存微信接口配置
func (p *foundationService) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	return p._rep.SaveWxApiConfig(v)
}

// 获取注册配置
func (p *foundationService) GetRegisterPerm() valueobject.RegisterPerm {
	return p._rep.GetRegisterPerm()
}

// 保存注册配置
func (p *foundationService) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	return p._rep.SaveRegisterPerm(v)
}

// 获取全局系统数值设置
func (p *foundationService) GetGlobNumberConf() valueobject.GlobNumberConf {
	return p._rep.GetGlobNumberConf()
}

// 保存全局系统数值设置
func (p *foundationService) SaveGlobNumberConf(v *valueobject.GlobNumberConf) error {
	return p._rep.SaveGlobNumberConf(v)
}

// 获取资源地址
func (p *foundationService) ResourceUrl(url string) (r string, err error) {
	return format.GetResUrl(url), nil
}

// 获取平台设置
func (p *foundationService) GetPlatformConf() (r *define.PlatformConf, err error) {
	v := p._rep.GetPlatformConf()
	return parser.PlatformConfDto(&v), nil
}

// 保存平台设置
func (p *foundationService) SavePlatformConf(v *valueobject.PlatformConf) error {
	return p._rep.SavePlatformConf(v)
}

// 获取全局商户销售设置
func (p *foundationService) GetGlobMchSaleConf() valueobject.GlobMchSaleConf {
	return p._rep.GetGlobMchSaleConf()
}

// 保存全局商户销售设置
func (p *foundationService) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	return p._rep.SaveGlobMchSaleConf(v)
}

// 获取短信设置
func (p *foundationService) GetSmsApiSet() valueobject.SmsApiSet {
	return p._rep.GetSmsApiSet()
}

// 保存短信API
func (p *foundationService) SaveSmsApiPerm(provider int, s *valueobject.SmsApiPerm) error {
	return p._rep.SaveSmsApiPerm(provider, s)
}

// 获取默认的短信API
func (p *foundationService) GetDefaultSmsApiPerm() (int, *valueobject.SmsApiPerm) {
	return p._rep.GetDefaultSmsApiPerm()
}

// 获取下级区域
func (p *foundationService) GetChildAreas(id int32) []*valueobject.Area {
	return p._rep.GetChildAreas(id)
}

// 获取地区名称
func (p *foundationService) GetAreaNames(id []int32) []string {
	return p._rep.GetAreaNames(id)
}

// 获取省市区字符串
func (p *foundationService) GetAreaString(province, city, district int32) string {
	if province == 0 || city == 0 || district == 0 {
		return ""
	}
	return p._rep.GetAreaString(province, city, district)
}
