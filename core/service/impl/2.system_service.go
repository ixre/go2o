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
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	de "github.com/ixre/go2o/core/domain/interface/domain"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/go2o/core/infrastructure/regex"
	"github.com/ixre/go2o/core/infrastructure/util/sensitive"
	"github.com/ixre/go2o/core/module"
	"github.com/ixre/go2o/core/module/bank"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/storage"
)

var _ proto.SystemServiceServer = new(systemServiceImpl)

// 基础服务
type systemServiceImpl struct {
	_repo        valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
	notifyRepo   mss.IMessageRepo
	sysRepo      sys.ISystemRepo
	_s           storage.Interface
	stationRepo  sys.IStationRepo
	serviceUtil
	proto.UnimplementedSystemServiceServer
}

func NewSystemService(rep valueobject.IValueRepo,
	registryRepo registry.IRegistryRepo,
	sysRepo sys.ISystemRepo,
	s storage.Interface,
	notifyRepo mss.IMessageRepo,
	stationRepo sys.IStationRepo,
) proto.SystemServiceServer {
	return &systemServiceImpl{
		_repo:        rep,
		_s:           s,
		notifyRepo:   notifyRepo,
		sysRepo:      sysRepo,
		registryRepo: registryRepo,
		stationRepo:  stationRepo,
	}
}

// GetSystemInfo implements proto.SystemServiceServer.
func (s *systemServiceImpl) GetSystemInfo(context.Context, *proto.Empty) (*proto.SSystemInfo, error) {
	isa := s.sysRepo.GetSystemAggregateRoot()
	return &proto.SSystemInfo{
		LastUpdateTime: uint64(isa.LastUpdateTime()),
	}, nil
}

// 检测是否包含敏感词
func (s *systemServiceImpl) CheckSensitive(_ context.Context, r *proto.String) (*proto.Bool, error) {
	_, b := sensitive.Singleton().CheckSensitive(r.Value)
	return &proto.Bool{Value: b}, nil
}

// 替换敏感词
func (s *systemServiceImpl) ReplaceSensitive(_ context.Context, r *proto.ReplaceSensitiveRequest) (*proto.String, error) {
	v := sensitive.Singleton().ReplaceAll(r.Text, r.Replacement)
	if r.Extra {
		v = regex.ContainPhoneRegex.ReplaceAllString(v, "***")
		v = regex.ContainEmailRegex.ReplaceAllString(v, "***")
		v = regex.ContainUrlRegex.ReplaceAllString(v, "***")
		v = regex.ContainChinaIDCardRegex.ReplaceAllString(v, "***")
		v = regex.ContainBankCardRegex.ReplaceAllString(v, "***")
		v = regex.ContainTelphoneRegex.ReplaceAllString(v, "***")
	}
	return &proto.String{Value: v}, nil
}

func (s *systemServiceImpl) CleanCache(_ context.Context, request *proto.CleanCacheRequest) (*proto.CleanCacheResponse, error) {
	if request.ClearGlobal {
		isa := s.sysRepo.GetSystemAggregateRoot()
		isa.FlushUpdateStatus()
	}
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

// SubmitSystemLog 提交系统日志
func (s *systemServiceImpl) SubmitSystemLog(_ context.Context, r *proto.SubmitSystemLogRequest) (*proto.TxResult, error) {
	ia := s.sysRepo.GetSystemAggregateRoot()

	extraInfo := ""
	if len(r.ExtraInfo) > 0 {
		bytes, _ := json.Marshal(r.ExtraInfo)
		extraInfo = string(bytes)
	}
	err := ia.Application().AddLog(&sys.SysLog{
		UserId:          int(r.UserId),
		Username:        r.Username,
		LogLevel:        int(r.Level),
		Message:         r.Message,
		Arguments:       r.Arguments,
		TerminalName:    r.TerminalName,
		TerminalEntry:   r.TerminalEntry,
		TerminalModel:   r.TerminalModel,
		TerminalVersion: r.TerminalVersion,
		ExtraInfo:       extraInfo,
	})
	return s.errorV2(err), nil
}

func (s *systemServiceImpl) DeleteSystemLog(_ context.Context, r *proto.DeleteSystemLogRequest) (*proto.TxResult, error) {
	ia := s.sysRepo.GetSystemAggregateRoot().Application()
	err := ia.DeleteLog(collections.MapList(r.Id, func(i int64) int {
		return int(i)
	}))
	return s.errorV2(err), nil
}

// 保存短信API凭据
func (s *systemServiceImpl) SaveSmsSetting(_ context.Context, r *proto.SSmsProviderSetting) (*proto.Result, error) {
	manager := s.notifyRepo.NotifyManager()
	if r.HttpExtra == nil {
		r.HttpExtra = &proto.SSmsExtraSetting{}
	}
	perm := &mss.SmsApiPerm{
		Provider:  int(r.Provider),
		Key:       r.Key,
		Secret:    r.Secret,
		Signature: r.Signature,
		Extra: &mss.SmsExtraSetting{
			ApiUrl:       r.HttpExtra.ApiUrl,
			Params:       r.HttpExtra.Params,
			Method:       r.HttpExtra.Method,
			Charset:      r.HttpExtra.Charset,
			SuccessChars: r.HttpExtra.SuccessChars,
		},
	}
	if err := manager.SaveSmsApiPerm(perm); err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 获取短信API凭据, @provider 短信服务商, 默认:http
func (s *systemServiceImpl) GetSmsSetting(_ context.Context, req *proto.GetSmsSettingRequest) (*proto.SSmsProviderSetting, error) {
	manager := s.notifyRepo.NotifyManager()
	perm := manager.GetSmsApiPerm(int(req.Provider))
	if perm == nil {
		return nil, mss.ErrNoSuchSmsProvider
	}
	return &proto.SSmsProviderSetting{
		Provider:  proto.ESmsProvider(perm.Provider),
		Key:       perm.Key,
		Secret:    perm.Secret,
		Signature: perm.Signature,
		HttpExtra: &proto.SSmsExtraSetting{
			ApiUrl:       perm.Extra.ApiUrl,
			Params:       perm.Extra.Params,
			Method:       perm.Extra.Method,
			Charset:      perm.Extra.Charset,
			SuccessChars: perm.Extra.SuccessChars,
		},
	}, nil
}

// 保存面板HOOK数据,这通常是在第三方应用中初始化或调用,参见文档：BoardHooks
func (s *systemServiceImpl) SaveBoardHook(_ context.Context, request *proto.BoardHookSaveRequest) (*proto.Result, error) {

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

// 保存超级用户账号和密码
func (s *systemServiceImpl) UpdateSuperCredential(_ context.Context, user *proto.SuperPassswordRequest) (*proto.TxResult, error) {
	if len(user.NewPassword) != 32 || len(user.OldPassword) != 32 {
		return s.errorV2(de.ErrNotMD5Format), nil
	}
	username := "master"
	oldPassword := domain.SuperPassword(username, user.OldPassword)
	pwd, _ := s.registryRepo.GetValue(registry.SysSuperLoginToken)
	if pwd != oldPassword {
		return s.errorV2(de.ErrPasswordNotMatch), nil
	}
	newPwd := domain.SuperPassword(username, user.NewPassword)
	err := s.registryRepo.UpdateValue(registry.SysSuperLoginToken, newPwd)
	return s.errorV2(err), nil
}

// 注册单点登录应用,返回值：
//   - 1. 成功，并返回token
//   - -1. 接口地址不正确
//   - -2. 已经注册
func (s *systemServiceImpl) RegisterApp(_ context.Context, app *proto.SSsoApp) (*proto.String, error) {
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
func (s *systemServiceImpl) GetApp(_ context.Context, s2 *proto.String) (*proto.SSsoApp, error) {
	sso := module.Get(module.SSO).(*module.SSOModule)
	return sso.Get(s2.Value), nil
}

// 获取单点登录应用
func (s *systemServiceImpl) GetAllSsoApp(_ context.Context, _ *proto.Empty) (*proto.StringListResponse, error) {
	sso := module.Get(module.SSO).(*module.SSOModule)
	return &proto.StringListResponse{
		Value: sso.Array(),
	}, nil
}

// 创建同步登录的地址
func (s *systemServiceImpl) GetSyncLoginUrl(_ context.Context, s2 *proto.String) (*proto.String, error) {
	panic("not implement")
	//return fmt.Sprintf("%_s://%_s%_s/auth?return_url=%_s",
	//	consts.DOMAIN_PASSPORT_PROTO, consts.DOMAIN_PREFIX_PASSPORT,
	//	variable.Domain, returnUrl), nil
}

// 获取移动应用设置
func (s *systemServiceImpl) GetMoAppConf(_ context.Context, _ *proto.Empty) (*proto.SMobileAppConfig, error) {
	c := s._repo.GetMoAppConf()
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
func (s *systemServiceImpl) SaveMoAppConf(_ context.Context, config *proto.SMobileAppConfig) (*proto.Result, error) {
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
	err := s._repo.SaveMoAppConf(dst)
	return s.error(err), nil
}

// 获取微信接口配置
func (s *systemServiceImpl) GetWxApiConfig(_ context.Context, empty *proto.Empty) (*proto.SWxApiConfig, error) {
	c := s._repo.GetWxApiConfig()
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
func (s *systemServiceImpl) SaveWxApiConfig(_ context.Context, cfg *proto.SWxApiConfig) (*proto.Result, error) {
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
	err := s._repo.SaveWxApiConfig(dst)
	return s.error(err), nil
}

// ResourceUrl 获取资源地址
func (s *systemServiceImpl) ResourceUrl(_ context.Context, s2 *proto.String) (*proto.String, error) {
	return &proto.String{Value: s2.Value}, nil
}

// GetSmsApiSet 获取短信设置
func (s *systemServiceImpl) GetSmsApiSet() mss.SmsApiSet {
	//return _s._rep.GetSmsApiSet()
	//todo: will remove
	return mss.SmsApiSet{}
}

// GetChildDistrict 获取下级区域
func (s *systemServiceImpl) GetChildDistrict(_ context.Context, code *proto.DistrictChildrenRequest) (*proto.AreaListResponse, error) {
	isa := s.sysRepo.GetSystemAggregateRoot()
	addr := isa.Address()
	arr := addr.GetChildrenDistricts(int(code.ParentId))
	ret := collections.MapList(arr, func(d *sys.District) *proto.SDistrict {
		return &proto.SDistrict{
			Id:       int64(d.Code),
			ParentId: int64(d.Parent),
			Name:     d.Name,
			IsLeaf:   len(addr.GetChildrenDistricts(d.Code)) == 0,
		}
	})
	return &proto.AreaListResponse{
		Value: ret,
	}, nil
}

// GetDistrictNames 获取地区名称
func (s *systemServiceImpl) GetDistrictNames(_ context.Context, request *proto.GetNamesRequest) (*proto.IntStringMapResponse, error) {
	isa := s.sysRepo.GetSystemAggregateRoot().Address()
	codes := collections.MapList(request.Value, func(i int32) int {
		return int(i)
	})
	mp := isa.GetDistrictNames(codes...)
	retMap := collections.Map(mp, func(k int, v string) (uint64, string) {
		return uint64(k), v
	})
	return &proto.IntStringMapResponse{
		Value: retMap,
	}, nil
}

func (s *systemServiceImpl) SaveGenericOption(_ context.Context, option *proto.SGenericOption) (*proto.TxResult, error) {
	isa := s.sysRepo.GetSystemAggregateRoot().Options()
	err := isa.SaveOption(&sys.GeneralOption{
		Id:      int(option.Id),
		Type:    option.Type,
		Pid:     int(option.Pid),
		Label:   option.Label,
		Value:   option.Value,
		SortNum: int(option.SortNum),
		Enabled: int(option.Enabled),
	})
	return s.errorV2(err), nil
}

// GetOptionNames implements proto.SystemServiceServer.
func (s *systemServiceImpl) GetOptionNames(_ context.Context, req *proto.GetNamesRequest) (*proto.IntStringMapResponse, error) {
	isa := s.sysRepo.GetSystemAggregateRoot().Options()
	codes := collections.MapList(req.Value, func(i int32) int {
		return int(i)
	})
	mp := isa.GetOptionNames(codes...)
	retMap := collections.Map(mp, func(k int, v string) (uint64, string) {
		return uint64(k), v
	})
	return &proto.IntStringMapResponse{
		Value: retMap,
	}, nil
}

// GetChildOptions implements proto.SystemServiceServer.
func (s *systemServiceImpl) GetChildOptions(_ context.Context, req *proto.OptionsRequest) (*proto.OptionsResponse, error) {
	isa := s.sysRepo.GetSystemAggregateRoot().Options()
	options := isa.GetChildOptions(int(req.ParentId), req.TypeName)
	sort.Sort(sys.GeneralOptions(options))
	ret := collections.MapList(options, func(o *sys.GeneralOption) *proto.SOption {
		return &proto.SOption{
			Id:     int64(o.Id),
			Name:   o.Label,
			Value:  o.Value,
			IsLeaf: isa.IsLeaf(o),
		}
	})
	return &proto.OptionsResponse{
		Value: ret,
	}, nil
}

// GetAreaString 获取省市区字符串
func (s *systemServiceImpl) GetAreaString(_ context.Context, r *proto.AreaStringRequest) (*proto.String, error) {
	if r.Province == 0 || r.City == 0 || r.District == 0 {
		return &proto.String{Value: ""}, nil
	}
	str := s._repo.GetAreaString(r.Province, r.City, r.District)
	return &proto.String{Value: str}, nil
}

// FindCity implements proto.SystemServiceServer.
func (s *systemServiceImpl) FindCity(_ context.Context, req *proto.FindAreaRequest) (*proto.SArea, error) {
	is := s.sysRepo.GetSystemAggregateRoot()
	isa := is.Address()
	city := isa.FindCity(req.Name)
	if city == nil {
		return &proto.SArea{}, nil
	}
	ret := &proto.SArea{
		Code:      int64(city.Code),
		Name:      city.Name,
		StationId: 0,
	}
	isn := is.Stations().FindStationByCity(city.Code)
	if isn != nil {
		ret.StationId = int64(isn.GetDomainId())
	}
	return ret, nil
}

// GetStation 获取站点
func (s *systemServiceImpl) GetStation(_ context.Context, req *proto.GetStationRequest) (*proto.SStation, error) {
	isa := s.stationRepo.GetStation(int(req.StationId))
	if isa == nil {
		return &proto.SStation{}, nil
	}

	v := isa.GetValue()
	ia := s.sysRepo.GetSystemAggregateRoot().Address()
	d := ia.GetDistrict(int(v.CityCode))
	ret := &proto.SStation{
		Id:         int64(isa.GetDomainId()),
		CityCode:   int64(v.CityCode),
		Status:     int32(v.Status),
		Letter:     v.Letter,
		IsHot:      v.IsHot == 1,
		CreateTime: int64(v.CreateTime),
	}
	if d != nil {
		ret.CityName = d.Name
	}
	return ret, nil
}

// GetPayPlatform 获取支付平台
func (s *systemServiceImpl) GetPayPlatform(_ context.Context, r *proto.Empty) (*proto.PaymentPlatformResponse, error) {
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
func (s *systemServiceImpl) GetGlobMchSaleConf_(_ context.Context, r *proto.Empty) (*proto.SGlobMchSaleConf, error) {
	c := s._repo.GetGlobMchSaleConf()
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
func (s *systemServiceImpl) SaveGlobMchSaleConf_(_ context.Context, conf *proto.SGlobMchSaleConf) (*proto.Result, error) {
	dst := &valueobject.GlobMchSaleConf{
		FxSalesEnabled:          conf.FxSalesEnabled,
		CashBackPercent:         float64(conf.CashBackPercent),
		CashBackTg1Percent:      float64(conf.CashBackTg1Percent),
		CashBackTg2Percent:      float64(conf.CashBackTg2Percent),
		CashBackMemberPercent:   float64(conf.CashBackMemberPercent),
		AutoSetupOrder:          int(conf.AutoSetupOrder),
		OrderTimeOutMinute:      int(conf.OrderTimeOutMinute),
		OrderConfirmAfterMinute: int(conf.OrderConfirmAfterMinute),
		OrderTimeOutReceiveHour: int(conf.OrderTimeOutReceiveHour),
	}
	err := s._repo.SaveGlobMchSaleConf(dst)
	return s.error(err), nil
}

func (s *systemServiceImpl) parsePayPlatform(v *bank.PaymentPlatform) *proto.PaymentPlatform {
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

func (s *systemServiceImpl) parseBankItem(v *bank.BankItem) *proto.BankItem {
	return &proto.BankItem{
		Id:   v.ID,
		Name: v.Name,
		Sign: v.Sign,
	}
}

// QueryBanks 查询银行列表
func (s *systemServiceImpl) QueryBanks(ctx context.Context, req *proto.QueryBanksRequest) (*proto.OptionsResponse, error) {
	isa := s.sysRepo.GetSystemAggregateRoot()
	arr := isa.GetBanks()

	name := regexp.MustCompile("(银行)(.+)").ReplaceAllString(req.Name, "$1")

	if len(name) == 0 {
		return &proto.OptionsResponse{}, nil
	}
	filters := collections.FilterArray(arr, func(o *sys.GeneralOption) bool {
		return strings.Contains(o.Label, name) || o.Value == name
	})
	ret := collections.MapList(filters, func(o *sys.GeneralOption) *proto.SOption {
		return &proto.SOption{
			Id:     int64(o.Id),
			Name:   o.Label,
			Value:  o.Value,
			IsLeaf: true,
		}
	})
	return &proto.OptionsResponse{
		Value: ret,
	}, nil
}
