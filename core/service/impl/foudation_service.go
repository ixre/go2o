package impl

/**
 * Copyright 2015 @ 56x.net.
 * name : platform_service
 * author : jarryliu
 * date : 2020-09-05 15:30
 * description :
 * history :
 */

import (
	"context"
	de "github.com/ixre/go2o/core/domain/interface/domain"
	"github.com/ixre/go2o/core/domain/interface/mss/notify"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/infrastructure/tool/sensitive"
	"github.com/ixre/go2o/core/module"
	"github.com/ixre/go2o/core/module/bank"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/storage"
)

var _ proto.FoundationServiceServer = new(foundationService)

// 基础服务
type foundationService struct {
	_rep         valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
	notifyRepo   notify.INotifyRepo
	_s           storage.Interface
	serviceUtil
}

func NewFoundationService(rep valueobject.IValueRepo, registryRepo registry.IRegistryRepo, s storage.Interface, notifyRepo notify.INotifyRepo) *foundationService {
	return &foundationService{
		_rep:         rep,
		_s:           s,
		notifyRepo:   notifyRepo,
		registryRepo: registryRepo,
	}
}

// 检测是否包含敏感词
func (s *foundationService) CheckSensitive(_ context.Context, r *proto.String) (*proto.Bool, error) {
	_, b := sensitive.Singleton().CheckSensitive(r.Value)
	return &proto.Bool{Value: b}, nil
}

// 替换敏感词
func (s *foundationService) ReplaceSensitive(_ context.Context, r *proto.ReplaceSensitiveRequest) (*proto.String, error) {
	v := sensitive.Singleton().ReplaceAll(r.Text, r.Replacement)
	return &proto.String{Value: v}, nil
}

func (s *foundationService) CleanCache(_ context.Context, request *proto.CleanCacheRequest) (*proto.CleanCacheResponse, error) {
	var count = 0
	if len(request.Prefix) > 0 {
		count, _ = s._s.DeleteWith(request.Prefix)
	} else if len(request.Key) > 0 {
		s._s.Delete(request.Key)
		count = 1
	}
	return &proto.CleanCacheResponse{
		Count: int32(count),
	}, nil
}

// 保存短信API凭据
func (s *foundationService) SaveSmsApi(_ context.Context, r *proto.SmsApiSaveRequest) (*proto.Result, error) {
	manager := s.notifyRepo.Manager()
	a := r.Api
	perm := &notify.SmsApiPerm{
		Key:         a.Key,
		Secret:      a.Secret,
		ApiUrl:      a.ApiUrl,
		Params:      a.Params,
		Method:      a.Method,
		Charset:     a.Charset,
		Signature:   a.Signature,
		SuccessChar: a.SuccessChar,
	}
	if err := manager.SaveSmsApiPerm(r.Provider, perm); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 获取短信API凭据, @provider 短信服务商, 默认:http
func (s *foundationService) GetSmsApi(_ context.Context, provider *proto.String) (*proto.SSmsApi, error) {
	manager := s.notifyRepo.Manager()
	perm := manager.GetSmsApiPerm(provider.Value)
	if perm != nil {
		return &proto.SSmsApi{
			ApiUrl:      perm.ApiUrl,
			Key:         perm.Key,
			Secret:      perm.Secret,
			Params:      perm.Params,
			Method:      perm.Method,
			Charset:     perm.Charset,
			Signature:   perm.Signature,
			SuccessChar: perm.SuccessChar,
		}, nil
	}
	return &proto.SSmsApi{
		ApiUrl:      "",
		Key:         "",
		Secret:      "",
		Params:      "",
		Method:      "",
		Charset:     "",
		Signature:   "",
		SuccessChar: "",
	}, nil
}

// 保存面板HOOK数据,这通常是在第三方应用中初始化或调用,参见文档：BoardHooks
func (s *foundationService) SaveBoardHook(_ context.Context, request *proto.BoardHookSaveRequest) (*proto.Result, error) {

	mp := map[string]string{
		registry.BoardHookURL:   request.HookURL,
		registry.BoardHookToken: request.Token,
	}
	for k, v := range mp {
		if ir := s.registryRepo.Get(k); ir != nil {
			err := ir.Update(v)
			if err == nil {
				err = ir.Save()
			}
			if err != nil {
				return s.error(err), nil
			}
		}
	}
	return s.success(nil), nil
}

//
//// 删除值
//func (_s *foundationService) DeleteValue(_ context.Context, s2 *proto.String) (*proto.Result, error) {
//	err := _s._rep.DeleteValue(s2.Value)
//	return _s.result(err), nil
//}
//
//// 根据前缀获取值
//func (_s *foundationService) GetValuesByPrefix(_ context.Context, s2 *proto.String) (*proto.StringMap, error) {
//	return &proto.StringMap{
//		Value: _s._rep.GetValues(s2.Value),
//	}, nil
//}

// 验证超级用户账号和密码
func (s *foundationService) SuperValidate(_ context.Context, user *proto.UserPwd) (*proto.SuperLoginResponse, error) {
	if len(user.Pwd) != 32 {
		return &proto.SuperLoginResponse{
			ErrMsg:  de.ErrNotMD5Format.Error(),
			ErrCode: 2}, nil
	}
	superPwd, _ := s.registryRepo.GetValue(registry.SysSuperLoginToken)
	encPwd := domain.Sha1Pwd(user.User+user.Pwd, "")
	if superPwd != encPwd {
		return &proto.SuperLoginResponse{
			ErrMsg:  de.ErrCredential.Error(),
			ErrCode: 1}, nil
	}
	return &proto.SuperLoginResponse{
		ErrCode: 0,
		Role:    1,
	}, nil
}

// 保存超级用户账号和密码
func (s *foundationService) FlushSuperPwd(_ context.Context, user *proto.UserPwd) (*proto.Result, error) {
	if len(user.Pwd) != 32 {
		return s.error(de.ErrNotMD5Format), nil
	}
	encPwd := domain.Sha1Pwd(user.User+user.Pwd, "")
	err := s.registryRepo.UpdateValue(registry.SysSuperLoginToken,
		encPwd)
	return s.error(err), nil
}

// 注册单点登录应用,返回值：
//   -  1. 成功，并返回token
//   - -1. 接口地址不正确
//   - -2. 已经注册
func (s *foundationService) RegisterApp(_ context.Context, app *proto.SSsoApp) (*proto.String, error) {
	sso := module.Get(module.SSO).(*module.SSOModule)
	token, err := sso.Register(app)
	if err == nil {
		return &proto.String{
			Value: "1:" + token,
		}, nil
	}
	return &proto.String{
		Value: err.Error(),
	}, nil
}

// 获取应用信息
func (s *foundationService) GetApp(_ context.Context, s2 *proto.String) (*proto.SSsoApp, error) {
	sso := module.Get(module.SSO).(*module.SSOModule)
	return sso.Get(s2.Value), nil
}

// 获取单点登录应用
func (s *foundationService) GetAllSsoApp(_ context.Context, _ *proto.Empty) (*proto.StringListResponse, error) {
	sso := module.Get(module.SSO).(*module.SSOModule)
	return &proto.StringListResponse{
		Value: sso.Array(),
	}, nil
}

// 创建同步登录的地址
func (s *foundationService) GetSyncLoginUrl(_ context.Context, s2 *proto.String) (*proto.String, error) {
	panic("not implement")
	//return fmt.Sprintf("%_s://%_s%_s/auth?return_url=%_s",
	//	consts.DOMAIN_PASSPORT_PROTO, consts.DOMAIN_PREFIX_PASSPORT,
	//	variable.Domain, returnUrl), nil
}

// 获取移动应用设置
func (s *foundationService) GetMoAppConf(_ context.Context, _ *proto.Empty) (*proto.SMobileAppConfig, error) {
	c := s._rep.GetMoAppConf()
	return &proto.SMobileAppConfig{
		AppName:           c.AppName,
		AppIcon:           c.AppIcon,
		Description:       c.Description,
		AndroidVersion:    c.AndroidVersion,
		AndroidReleaseUrl: c.AndroidReleaseUrl,
		IosVersion:        c.IosVersion,
		IosReleaseUrl:     c.IosReleaseUrl,
		ShowTplPath_:      c.ShowTplPath,
	}, nil
}

// 保存移动应用设置
func (s *foundationService) SaveMoAppConf(_ context.Context, config *proto.SMobileAppConfig) (*proto.Result, error) {
	dst := &valueobject.MoAppConf{
		AppName:           config.AppName,
		AppIcon:           config.AppIcon,
		Description:       config.Description,
		ShowTplPath:       config.ShowTplPath_,
		AndroidVersion:    config.AndroidVersion,
		AndroidReleaseUrl: config.AndroidReleaseUrl,
		IosVersion:        config.IosVersion,
		IosReleaseUrl:     config.IosReleaseUrl,
		WpVersion:         "",
		WpReleaseUrl:      "",
	}
	err := s._rep.SaveMoAppConf(dst)
	return s.error(err), nil
}

// 获取微信接口配置
func (s *foundationService) GetWxApiConfig(_ context.Context, empty *proto.Empty) (*proto.SWxApiConfig, error) {
	c := s._rep.GetWxApiConfig()
	return &proto.SWxApiConfig{
		AppId:               c.AppId,
		AppSecret:           c.AppSecret,
		MpToken:             c.MpToken,
		MpAesKey:            c.MpAesKey,
		OriId:               c.OriId,
		MchId:               c.MchId,
		MchApiKey:           c.MchApiKey,
		MchCertPath:         c.MchCertPath,
		MchCertKeyPath:      c.MchCertKeyPath,
		RedPackEnabled:      c.RedPackEnabled,
		RedPackAmountLimit:  float64(c.RedPackAmountLimit),
		RedPackDayTimeLimit: int32(c.RedPackDayTimeLimit),
	}, nil
}

// 保存微信接口配置
func (s *foundationService) SaveWxApiConfig(_ context.Context, cfg *proto.SWxApiConfig) (*proto.Result, error) {
	dst := &valueobject.WxApiConfig{
		AppId:               cfg.AppId,
		AppSecret:           cfg.AppSecret,
		MpToken:             cfg.MpToken,
		MpAesKey:            cfg.MpAesKey,
		OriId:               cfg.OriId,
		MchId:               cfg.MchId,
		MchApiKey:           cfg.MchApiKey,
		MchCertPath:         cfg.MchCertPath,
		MchCertKeyPath:      cfg.MchCertKeyPath,
		RedPackEnabled:      cfg.RedPackEnabled,
		RedPackAmountLimit:  float32(cfg.RedPackAmountLimit),
		RedPackDayTimeLimit: int(cfg.RedPackDayTimeLimit),
	}
	err := s._rep.SaveWxApiConfig(dst)
	return s.error(err), nil
}

// ResourceUrl 获取资源地址
func (s *foundationService) ResourceUrl(_ context.Context, s2 *proto.String) (*proto.String, error) {
	return &proto.String{Value: format.GetResUrl(s2.Value)}, nil
}

// GetSmsApiSet 获取短信设置
func (s *foundationService) GetSmsApiSet() notify.SmsApiSet {
	//return _s._rep.GetSmsApiSet()
	//todo: will remove
	return notify.SmsApiSet{}
}

// GetChildAreas 获取下级区域
func (s *foundationService) GetChildAreas(_ context.Context, code *proto.Int32) (*proto.AreaListResponse, error) {
	var arr []*proto.SArea
	for _, v := range s._rep.GetChildAreas(code.Value) {
		arr = append(arr, &proto.SArea{
			Id:       v.Code,
			ParentId: v.Parent,
			Name:     v.Name,
		})
	}
	return &proto.AreaListResponse{
		Value: arr,
	}, nil
}

// GetAreaNames 获取地区名称
func (s *foundationService) GetAreaNames(_ context.Context, request *proto.GetAreaNamesRequest) (*proto.StringListResponse, error) {
	return &proto.StringListResponse{
		Value: s._rep.GetAreaNames(request.Value),
	}, nil
}

// GetAreaString 获取省市区字符串
func (s *foundationService) GetAreaString(_ context.Context, r *proto.AreaStringRequest) (*proto.String, error) {
	if r.Province == 0 || r.City == 0 || r.District == 0 {
		return &proto.String{Value: ""}, nil
	}
	str := s._rep.GetAreaString(r.Province, r.City, r.District)
	return &proto.String{Value: str}, nil
}

// GetPayPlatform 获取支付平台
func (s *foundationService) GetPayPlatform(_ context.Context, r *proto.Empty) (*proto.PaymentPlatformResponse, error) {
	m := module.Get(module.PAY).(*module.PaymentModule)
	pf := m.GetPayPlatform()
	ret := &proto.PaymentPlatformResponse{
		Value: make([]*proto.PaymentPlatform, len(pf)),
	}
	for i, v := range pf {
		ret.Value[i] = s.parsePayPlatform(v)
	}
	return ret, nil
}

// GetGlobMchSaleConf_ 获取全局商户销售设置
func (s *foundationService) GetGlobMchSaleConf_(_ context.Context, r *proto.Empty) (*proto.SGlobMchSaleConf, error) {
	c := s._rep.GetGlobMchSaleConf()
	return &proto.SGlobMchSaleConf{
		FxSalesEnabled:          c.FxSalesEnabled,
		CashBackPercent:         float64(c.CashBackPercent),
		CashBackTg1Percent:      float64(c.CashBackTg1Percent),
		CashBackTg2Percent:      float64(c.CashBackTg2Percent),
		CashBackMemberPercent:   float64(c.CashBackMemberPercent),
		AutoSetupOrder:          int32(c.AutoSetupOrder),
		OrderTimeOutMinute:      int32(c.OrderTimeOutMinute),
		OrderConfirmAfterMinute: int32(c.OrderConfirmAfterMinute),
		OrderTimeOutReceiveHour: int32(c.OrderTimeOutReceiveHour),
	}, nil
}

// 保存全局商户销售设置
func (s *foundationService) SaveGlobMchSaleConf_(_ context.Context, conf *proto.SGlobMchSaleConf) (*proto.Result, error) {
	dst := &valueobject.GlobMchSaleConf{
		FxSalesEnabled:          conf.FxSalesEnabled,
		CashBackPercent:         float32(conf.CashBackPercent),
		CashBackTg1Percent:      float32(conf.CashBackTg1Percent),
		CashBackTg2Percent:      float32(conf.CashBackTg2Percent),
		CashBackMemberPercent:   float32(conf.CashBackMemberPercent),
		AutoSetupOrder:          int(conf.AutoSetupOrder),
		OrderTimeOutMinute:      int(conf.OrderTimeOutMinute),
		OrderConfirmAfterMinute: int(conf.OrderConfirmAfterMinute),
		OrderTimeOutReceiveHour: int(conf.OrderTimeOutReceiveHour),
	}
	err := s._rep.SaveGlobMchSaleConf(dst)
	return s.error(err), nil
}

func (s *foundationService) parsePayPlatform(v *bank.PaymentPlatform) *proto.PaymentPlatform {
	dst := &proto.PaymentPlatform{
		Id:    v.ID,
		Name:  v.Name,
		Sign:  v.Sign,
		Items: make([]*proto.BankItem, len(v.Bank)),
	}
	if v.Bank == nil {
		for i, v := range v.Bank {
			dst.Items[i] = s.parseBankItem(v)
		}
	}
	return dst
}

func (s *foundationService) parseBankItem(v *bank.BankItem) *proto.BankItem {
	return &proto.BankItem{
		Id:   v.ID,
		Name: v.Name,
		Sign: v.Sign,
	}
}
