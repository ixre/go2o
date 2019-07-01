/**
 * Copyright 2015 @ to2.net.
 * name : platform_service
 * author : jarryliu
 * date : 2016-05-27 15:30
 * description :
 * history :
 */
package rsi

import (
	"context"
	"errors"
	"github.com/ixre/gof"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/module"
	"go2o/core/module/bank"
	"go2o/core/service/auto_gen/rpc/foundation_service"
	"go2o/core/service/auto_gen/rpc/ttype"
)

var _ foundation_service.FoundationService = new(foundationService)

// 基础服务
type foundationService struct {
	_rep         valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
	serviceUtil
}

func NewFoundationService(rep valueobject.IValueRepo, registryRepo registry.IRegistryRepo) *foundationService {
	return &foundationService{
		_rep:         rep,
		registryRepo: registryRepo,
	}
}

// 根据键获取值
func (s *foundationService) GetValue(ctx context.Context, key string) (r string, err error) {
	ir := s.registryRepo.Get(key)
	if ir != nil {
		return ir.StringValue(), nil
	}
	return "", nil
}

// 获取键值存储数据
func (s *foundationService) GetRegistries(ctx context.Context, keys []string) (map[string]string, error) {
	mp := make(map[string]string)
	for _, k := range keys {
		if ir := s.registryRepo.Get(k); ir != nil {
			mp[k] = ir.StringValue()
		} else {
			mp[k] = ""
		}
	}
	return mp, nil
}

// 搜索注册表
func (s *foundationService) SearchRegistry(ctx context.Context, key string) (r []*foundation_service.SRegistry, err error) {
	arr := s.registryRepo.SearchRegistry(key)
	list := make([]*foundation_service.SRegistry, len(arr))
	for i, a := range arr {
		list[i] = &foundation_service.SRegistry{
			Key:         a.Key,
			Value:       a.Value,
			Default:     a.DefaultValue,
			Options:     a.Options,
			UserDefine:  a.UserDefine == 1,
			Description: a.Description,
		}
	}
	return list, nil
}

// 获取数据存储
func (s *foundationService) GetRegistry(ctx context.Context, key string) (string, error) {
	ir := s.registryRepo.Get(key)
	if ir != nil {
		return ir.StringValue(), nil
	}
	return "", nil
}

// 创建用户自定义注册项
func (s *foundationService) CreateUserRegistry(ctx context.Context, key string, defaultValue string, description string) (r *ttype.Result_, err error) {
	if s.registryRepo.Get(key) != nil {
		return s.resultWithCode(-1, "registry is exist"), nil
	}
	rv := &registry.Registry{
		Key:          key,
		Value:        defaultValue,
		DefaultValue: defaultValue,
		Options:      "",
		UserDefine:   1,
		Description:  description,
	}
	ir := s.registryRepo.Create(rv)
	err = ir.Save()
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 更新注册表数据
func (s *foundationService) UpdateRegistry(ctx context.Context, registries map[string]string) (r *ttype.Result_, err error) {
	for k, v := range registries {
		if ir := s.registryRepo.Get(k); ir != nil {
			if err = ir.Update(v); err != nil {
				return s.error(err), nil
			}
		}
	}
	r = s.success(nil)
	return r, nil
}

// 设置键值
func (s *foundationService) SetValue(ctx context.Context, key string, value string) (r *ttype.Result_, err error) {
	err = s._rep.SetValue(key, value)
	return s.result(err), nil
}

// 删除值
func (s *foundationService) DeleteValue(ctx context.Context, key string) (r *ttype.Result_, err error) {
	err = s._rep.DeleteValue(key)
	return s.result(err), nil
}

// 根据前缀获取值
func (s *foundationService) GetValuesByPrefix(ctx context.Context, prefix string) (r map[string]string, err error) {
	return s._rep.GetValues(prefix), nil
}

// 获取键值存储数据
func (s *foundationService) GetRegistryV1(ctx context.Context, keys []string) ([]string, error) {
	return s._rep.GetsRegistry(keys), nil
}

// 验证超级用户账号和密码
func (s *foundationService) SuperValidate(ctx context.Context, user string, pwd string) (r bool, err error) {
	superPwd := gof.CurrentApp.Config().Get("super_login_md5")
	encPwd := domain.Sha1Pwd(pwd + user)
	return superPwd == encPwd, nil
}

// 保存超级用户账号和密码
func (s *foundationService) FlushSuperPwd(ctx context.Context, user string, pwd string) (err error) {
	conf := gof.CurrentApp.Config()
	encPwd := domain.Sha1Pwd(pwd + user)
	conf.Set("super_login_md5", encPwd)
	//conf.Flush()
	return errors.New("暂不支持保存")
}

// 注册单点登录应用,返回值：
//   -  1. 成功，并返回token
//   - -1. 接口地址不正确
//   - -2. 已经注册
func (s *foundationService) RegisterApp(ctx context.Context, app *foundation_service.SSsoApp) (r string, err error) {
	sso := module.Get(module.SSO).(*module.SSOModule)
	token, err := sso.Register(app)
	if err == nil {
		return "1:" + token, nil
	}
	return err.Error(), nil
}

// 获取应用信息
func (s *foundationService) GetApp(ctx context.Context, name string) (r *foundation_service.SSsoApp, err error) {
	sso := module.Get(module.SSO).(*module.SSOModule)
	return sso.Get(name), nil
}

// 获取单点登录应用
func (s *foundationService) GetAllSsoApp(ctx context.Context) (r []string, err error) {
	sso := module.Get(module.SSO).(*module.SSOModule)
	return sso.Array(), nil
}

// 创建同步登录的地址
func (s *foundationService) GetSyncLoginUrl(ctx context.Context, returnUrl string) (r string, err error) {
	panic("not implement")
	//return fmt.Sprintf("%s://%s%s/auth?return_url=%s",
	//	consts.DOMAIN_PASSPORT_PROTO, consts.DOMAIN_PREFIX_PASSPORT,
	//	variable.Domain, returnUrl), nil
}

// 获取移动应用设置
func (s *foundationService) GetMoAppConf() valueobject.MoAppConf {
	return s._rep.GetMoAppConf()
}

// 保存移动应用设置
func (s *foundationService) SaveMoAppConf(r *valueobject.MoAppConf) error {
	return s._rep.SaveMoAppConf(r)
}

// 获取微信接口配置
func (s *foundationService) GetWxApiConfig() valueobject.WxApiConfig {
	return s._rep.GetWxApiConfig()
}

// 保存微信接口配置
func (s *foundationService) SaveWxApiConfig(v *valueobject.WxApiConfig) error {
	return s._rep.SaveWxApiConfig(v)
}

// 获取注册配置
func (s *foundationService) GetRegisterPerm() valueobject.RegisterPerm {
	return s._rep.GetRegisterPerm()
}

// 保存注册配置
func (s *foundationService) SaveRegisterPerm(v *valueobject.RegisterPerm) error {
	return s._rep.SaveRegisterPerm(v)
}

// 获取资源地址
func (s *foundationService) ResourceUrl(ctx context.Context, url string) (r string, err error) {
	return format.GetResUrl(url), nil
}

// 获取全局商户销售设置
func (s *foundationService) GetGlobMchSaleConf() valueobject.GlobMchSaleConf {
	return s._rep.GetGlobMchSaleConf()
}

// 保存全局商户销售设置
func (s *foundationService) SaveGlobMchSaleConf(v *valueobject.GlobMchSaleConf) error {
	return s._rep.SaveGlobMchSaleConf(v)
}

// 获取短信设置
func (s *foundationService) GetSmsApiSet() valueobject.SmsApiSet {
	return s._rep.GetSmsApiSet()
}

// 保存短信API
func (s *foundationService) SaveSmsApiPerm(provider int, perm *valueobject.SmsApiPerm) error {
	return s._rep.SaveSmsApiPerm(provider, perm)
}

// 获取默认的短信API
func (s *foundationService) GetDefaultSmsApiPerm() (int, *valueobject.SmsApiPerm) {
	return s._rep.GetDefaultSmsApiPerm()
}

// 获取下级区域
func (s *foundationService) GetChildAreas(ctx context.Context, code int32) ([]*foundation_service.SArea, error) {
	var arr []*foundation_service.SArea
	for _, v := range s._rep.GetChildAreas(code) {
		arr = append(arr, &foundation_service.SArea{
			Code:   int32(v.Code),
			Parent: int32(v.Parent),
			Name:   v.Name,
		})
	}
	return arr, nil
}

// 获取地区名称
func (s *foundationService) GetAreaNames(ctx context.Context, codes []int32) ([]string, error) {
	return s._rep.GetAreaNames(codes), nil
}

// 获取省市区字符串
func (s *foundationService) GetAreaString(province, city, district int32) string {
	if province == 0 || city == 0 || district == 0 {
		return ""
	}
	return s._rep.GetAreaString(province, city, district)
}

// 获取支付平台
func (s *foundationService) GetPayPlatform() []*bank.PaymentPlatform {
	m := module.Get(module.PAY).(*module.PaymentModule)
	return m.GetPayPlatform()
}
