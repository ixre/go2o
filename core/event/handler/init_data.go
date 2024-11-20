package handler

import (
	"log"
	"strings"
	"time"

	"github.com/ixre/go2o/core/dao"
	"github.com/ixre/go2o/core/dao/model"
	"github.com/ixre/go2o/core/domain/interface/content"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/sys"

	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/util"
)

// 子订单推送
func (h EventHandler) HandleAppInitialEvent(data interface{}) {
	v := data.(*events.AppInitialEvent)
	if v == nil {
		return
	}
	initSystemConfig(h.registryRepo)
	initJWTSecret(h.registryRepo)
	initSuperLoginToken(h.registryRepo)
	initPortalNav(h.portalDao)
	initPortalNavGroup(h.portalDao)
	initPages(h.pageRepo)
	// 初始化站点
	go h._sysRepo.GetSystemAggregateRoot().Stations().SyncStations()
	// 初始化通知模板
	h.initNotifyTemplate()
	// 初始化分发应用
	h.initDistributeApp()
}

func initSystemConfig(repo registry.IRegistryRepo) {
	// 初始化系统配置
	re := repo.Get(registry.Domain)
	if re != nil && re.StringValue() == re.Value().DefaultValue {
		logger.Error("主域未配置，会导致部分功能无法使用, 请在系统设置中配置")
	}
	re = repo.Get(registry.FileServerUrl)
	if re != nil && re.StringValue() == re.Value().DefaultValue {
		logger.Error("文件服务器地址未配置，会导致图片无法上传, 请在系统设置中配置")
	}
	re = repo.Get(registry.MchServerUrl)
	if re != nil && re.StringValue() == re.Value().DefaultValue {
		logger.Error("B端商户服务器地址未配置，会导致商户部分功能不能正常使用, 请在系统设置中配置")
	}
}

func initSuperLoginToken(repo registry.IRegistryRepo) {
	value, _ := repo.GetValue(registry.SysSuperLoginToken)
	if strings.TrimSpace(value) == "" {
		pwd := util.RandString(8)
		log.Printf(`[ GO2O][ INFO]: the initial super pwd is '%s', it only show first time. plese save it.\n`, pwd)
		token := domain.Sha1("master" + crypto.Md5([]byte(pwd)))
		_ = repo.UpdateValue(registry.SysSuperLoginToken, token)
	}

}

// 初始化jwt密钥
func initJWTSecret(repo registry.IRegistryRepo) {
	value, _ := repo.GetValue(registry.SysJWTSecret)
	if strings.TrimSpace(value) == "" {
		_, privateKey, _ := crypto.GenRsaKeys(2048)
		_ = repo.UpdateValue(registry.SysJWTSecret, privateKey)
	}
}

// 初始化导航数据
func initPortalNav(repo dao.IPortalDao) {
	nav := repo.SelectNav("")
	if len(nav) == 0 {
		arr := []*model.PortalNav{
			{
				Text:    "超市",
				Url:     "#",
				Image:   "http://m.360buyimg.com/mobilecms/s120x120_jfs/t1/175540/24/19329/6842/60ec0b0aEf35f7384/ec560dbf9b82b90b.png!q70.jpg.dpg",
				NavType: 1,
			},
			{

				Text:    "数码",
				Url:     "#",
				Image:   "http://m.360buyimg.com/mobilecms/s120x120_jfs/t1/178015/31/13828/6862/60ec0c04Ee2fd63ac/ccf74d805a059a44.png!q70.jpg.dpg",
				NavType: 1,
			},
			{

				Text:    "服饰",
				Url:     "#",
				Image:   "http://m.360buyimg.com/mobilecms/s120x120_jfs/t1/41867/2/15966/7116/60ec0e0dE9f50d596/758babcb4f911bf4.png!q70.jpg.dpg",
				NavType: 1,
			},
			{

				Text:    "食品",
				Url:     "#",
				Image:   "http://m.360buyimg.com/mobilecms/s120x120_jfs/t1/177902/16/13776/5658/60ec0e71E801087f2/a0d5a68bf1461e6d.png!q70.jpg.dpg",
				NavType: 1,
			},
			{

				Text:    "家具",
				Url:     "#",
				Image:   "http://m.360buyimg.com/mobilecms/s120x120_jfs/t1/199143/10/8979/4223/614599f5E45cd5464/d15aa650a0ebe596.png!q70.jpg.dpg",
				NavType: 1,
			},
			{

				Text:    "VIP专区",
				Url:     "#",
				Image:   "http://m.360buyimg.com/mobilecms/s120x120_jfs/t1/37709/6/15279/6118/60ec1046E4b5592c6/a7d6b66354efb141.png!q70.jpg.dpg",
				NavType: 1,
			},
			{

				Text:    "优惠券",
				Url:     "#",
				Image:   "http://m.360buyimg.com/mobilecms/s120x120_jfs/t1/186080/16/13681/8175/60ec0fcdE032af6cf/c5acd2f8454c40e1.png!q70.jpg.dpg",
				NavType: 1,
			}, {

				Text:    "点卡",
				Url:     "#",
				Image:   "http://m.360buyimg.com/mobilecms/s120x120_jfs/t1/185733/21/13527/6648/60ec0f31E0fea3e0a/d86d463521140bb6.png!q70.jpg.dpg",
				NavType: 1,
			},
		}
		for _, v := range arr {
			repo.SaveNav(v)
			v.NavType = 2
			v.Id = 0
			repo.SaveNav(v)
		}
	}
}

// 初始化导航分组
func initPortalNavGroup(repo dao.IPortalDao) {
	group := repo.SelectNavGroup("")
	if len(group) == 0 {
		arr := []string{"头部导航", "友情链接"}
		for _, v := range arr {
			repo.SaveNavGroup(&model.NavGroup{
				Name: v,
			})
		}
	}
}

// 初始化内置页面
func initPages(repo content.IPageRepo) {
	ip := repo.GetPageByCode(0, "privacy")
	if ip == nil {
		pages := []*content.Page{
			{
				Title: "隐私政策",
				Code:  "privacy",
				Content: "请您务必审慎阅读,并充分理解\"服务协议\"和\"隐私政策\"各条款，为了向您提供相关服务，我们需要收集你您的设备信息、操作日志等个人信息。" +
					"您可以在\"设置\"中查看、变更和删除个人信息并管理您的授权。您可阅读《服务协议》和《隐私政策》了解详细信息。" +
					"如果您同意，请点击“同意”开始使用我们的服务。",
			},
			{
				Title:   "用户条款",
				Code:    "agreement",
				Content: "暂无内容",
			},
			{
				Title:   "关于平台",
				Code:    "about",
				Content: "",
			},
			{
				Title:   "联系我们",
				Code:    "contact",
				Content: "",
			},
		}
		for _, v := range pages {
			v.Flag |= content.FCategoryInternal
			v.Enabled = 1
			v.UpdateTime = int(time.Now().Unix())
			repo.SavePage(0, v)
		}
	}
}

// 初始化通知模板
func (h *EventHandler) initNotifyTemplate() {
	arr := h.messageRepo.NotifyRepo().GetAllNotifyTemplate()
	// 初始化短信模板
	tplArr := collections.FilterArray(arr, func(t *mss.NotifyTemplate) bool {
		return t.TplType == 2
	})
	if len(tplArr) == 0 {
		for _, v := range mss.InternalSmsTemplate {
			h.messageRepo.NotifyRepo().SaveNotifyTemplate(v)
		}
	}
	// 初始化邮件模板
	tplArr = collections.FilterArray(arr, func(t *mss.NotifyTemplate) bool {
		return t.TplType == 3
	})
	if len(tplArr) == 0 {
		for _, v := range mss.InternalMailTemplate {
			h.messageRepo.NotifyRepo().SaveNotifyTemplate(v)
		}
	}
}

// 初始化应用分发
func (h *EventHandler) initDistributeApp() {
	ia := h._sysRepo.GetSystemAggregateRoot().Application()
	arr := ia.GetAllAppDistributions()
	if len(arr) == 0 {
		err := ia.SaveAppDistribution(&sys.SysAppDistribution{
			AppName:        "app",
			UpdateMode:     1,
			UrlScheme:      "go2o://net.fze.go2o/open",
			DistributeName: "默认APP",
		})
		if err != nil {
			logger.Error("初始化应用分发失败:%s", err.Error())
		}
	}
}
