package impl

import (
	"github.com/ixre/go2o/core/dao/impl"
	"github.com/ixre/go2o/core/dao/model"
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/repos"
	"github.com/ixre/gof/db/orm"
	"time"
)

func InitData(o orm.Orm) {
	(&dataInitializer{
		o: o,
	}).init()
}

type dataInitializer struct {
	o orm.Orm
}

func (i dataInitializer) init() {
	i.initPortalNav()
	i.initPortalNavGroup()
	i.initPages()
}

// 初始化导航数据
func (i dataInitializer) initPortalNav() {
	repo := impl.NewPortalDao(i.o)
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
func (i dataInitializer) initPortalNavGroup() {
	repo := impl.NewPortalDao(i.o)
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
func (i dataInitializer) initPages() {
	repo := repos.Repo.GetContentRepo()
	ip := repo.GetPageByCode(0, "privacy")
	if ip == nil{
		pages := []*content.Page{
			{
				Title:       "隐私政策",
				Code:        "privacy",
				Content:     "请您务必审慎阅读,并充分理解\"服务协议\"和\"隐私政策\"各条款，为了向您提供相关服务，我们需要收集你您的设备信息、操作日志等个人信息。" +
					"您可以在\"设置\"中查看、变更和删除个人信息并管理您的授权。您可阅读《服务协议》和《隐私政策》了解详细信息。" +
					"如果您同意，请点击“同意”开始使用我们的服务。",
			},
			{
				Title:       "用户服务协议",
				Code:        "protocol",
				Content:     "",
			},
			{
				Title:       "关于平台",
				Code:        "about",
				Content:     "",
			},
			{
				Title:       "联系我们",
				Code:        "contact",
				Content:     "",
			},
		}
		for _,v := range pages{
			v.Flag |= content.FlagInternal
			v.Enabled = 1
			v.UpdateTime = time.Now().Unix()
			repo.SavePage(0,v)
		}
	}
}
