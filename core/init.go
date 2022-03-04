/**
 * Copyright 2015 @ 56x.net.
 * name : types.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package core

import (
	"context"
	"encoding/gob"
	"github.com/ixre/go2o/core/domain/interface/ad"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/module/express/kdniao"
	"github.com/ixre/go2o/core/msq"
	"github.com/ixre/go2o/core/service"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof/log"
	"os"
)

func init() {
	registerTypes()
}

var startJobs = make([]func(), 0)

func Startup(job func()) {
	startJobs = append(startJobs, job)
}

// 注册序列类型
func registerTypes() {
	gob.Register(&member.Member{})
	gob.Register(&merchant.Merchant{})
	gob.Register(&merchant.ApiInfo{})
	gob.Register(&shop.OnlineShop{})
	gob.Register(&shop.OfflineShop{})
	gob.Register(&shop.ComplexShop{})
	gob.Register(&member.Account{})
	gob.Register(&payment.Order{})
	gob.Register(&member.InviteRelation{})
	gob.Register(&dto.ListOnlineShop{})
	gob.Register([]*dto.ListOnlineShop{})
	gob.Register(&proto.SMember{})
	gob.Register(&proto.SProfile{})
	init2()
}

func init2() {
	gob.Register(map[string]map[string]interface{}{})
	gob.Register(ad.SwiperAd{})
	gob.Register(ad.Ad{})
	gob.Register([]*valueobject.Goods{})
	gob.Register(valueobject.Goods{})
	gob.Register(ad.HyperLink{})
	gob.Register(ad.Image{})
}

func Init(a *AppImpl, debug, trace bool) bool {
	a._debugMode = debug
	// 初始化clickhouse
	//clickhouse.Initialize(a)
	// 初始化变量
	variable.Domain = a.Config().GetString(variable.ServerDomain)
	a.Loaded = true
	for _, f := range startJobs {
		f()
	}
	return true
}

func AppDispose() {
	//GetRedisPool().Close()
	msq.Close()
	//if clickhouse.ConnInstance != nil{
	//	clickhouse.ConnInstance.Close()
	//}
}

func InitialModules() {
	initExpressAPI()
	initBankB4eAPI()
	initSSOModule()
}

func initSSOModule() {
	//domain := variable.Domain
	trans, _, err := service.RegistryServiceClient()
	if err == nil {
		defer trans.Close()
		keys := []string{
			registry.DomainPrefixPortal,
			registry.DomainPrefixWholesalePortal,
			registry.DomainPrefixHApi,
			registry.DomainPrefixMember,
			registry.DomainPrefixMobileMember,
			registry.DomainPrefixMobilePortal,
		}

		println(len(keys))
		//todo: to etcd
		/*
			registries, _ := cli.GetValues(context.TODO(),&proto.StringArray{Value:  keys})
			_, _ = s.Register(&proto.SSsoApp{
				Id:   1,
				Name: "RetailPortal",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[0]], domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				Id:   2,
				Name: "WholesalePortal",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[1]], domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				Id:   3,
				Name: "HApi",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[2]], domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				Id:   4,
				Name: "Member",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[3]], domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				Id:   5,
				Name: "MemberMobile",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[4]],
					domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				Id:   6,
				Name: "RetailPortalMobile",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[5]], domain),
			})

		*/
	}
}

func initBankB4eAPI() {
	trans, cli, err := service.RegistryServiceClient()
	if err == nil {
		ctx := context.TODO()
		defer trans.Close()
		_, _ = cli.CreateRegistry(ctx, &proto.RegistryCreateRequest{
			Key:          "bank4e_trust_on",
			DefaultValue: "false",
			Description:  "是否开启四要素实名认证",
		})
		_, _ = cli.CreateRegistry(ctx, &proto.RegistryCreateRequest{
			Key:          "bank4e_jd_app_key",
			DefaultValue: "",
			Description:  "京东银行四要素接口KEY",
		})

		//todo: etcd

		//data, _ := cli.GetValues(ctx, &proto.StringArray{Value: keys})
		//b.open, _ = strconv.ParseBool(data.Value[keys[0]])
		//b.appKey = data.Value[keys[1]]
	}
}

func initExpressAPI() {
	trans, cli, err := service.RegistryServiceClient()
	if err == nil {
		defer trans.Close()
		keys := []string{"express_kdn_business_id", "express_kdn_api_key"}
		_, _ = cli.CreateRegistry(context.TODO(),
			&proto.RegistryCreateRequest{
				Key:          keys[0],
				DefaultValue: "1314567",
				Description:  "快递鸟接口业务ID",
			})
		_, _ = cli.CreateRegistry(context.TODO(),
			&proto.RegistryCreateRequest{
				Key:          keys[1],
				DefaultValue: "27d809c3-51b6-479c-9b77-6b98d7f3d41",
				Description:  "快递鸟接口KEY",
			})
		data, _ := cli.GetValues(context.TODO(), &proto.StringArray{Value: keys})
		kdniao.EBusinessID = data.Value[keys[0]]
		kdniao.AppKey = data.Value[keys[1]]
	} else {
		log.Println("intialize express module error:", err.Error())
		os.Exit(1)
	}
}
