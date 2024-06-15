package app

import (
	"context"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/module/express/kdniao"
	"github.com/ixre/go2o/core/service/proto"
)

func InitialModules() {
	initExpressAPI()
	initBankB4eAPI()
	initSSOModule()
}

func initSSOModule() {

	//domain := variable.Domain
	//service := inject.GetRegistryService()

	// keys := []string{
	// 	registry.DomainPrefixPortal,
	// 	registry.DomainPrefixWholesalePortal,
	// 	registry.DomainPrefixHApi,
	// 	registry.DomainPrefixMember,
	// 	registry.DomainPrefixMobileMember,
	// 	registry.DomainPrefixMobilePortal,
	// }

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

func initBankB4eAPI() {
	cli := inject.GetRegistryService()

	_, _ = cli.CreateRegistry(context.TODO(), &proto.RegistryCreateRequest{
		Key:          "bank4e_trust_on",
		DefaultValue: "false",
		Description:  "是否开启四要素实名认证",
	})
	_, _ = cli.CreateRegistry(context.TODO(), &proto.RegistryCreateRequest{
		Key:          "bank4e_jd_app_key",
		DefaultValue: "",
		Description:  "京东银行四要素接口KEY",
	})

	//todo: etcd

	//data, _ := cli.GetValues(ctx, &proto.StringArray{Value: keys})
	//b.open, _ = strconv.ParseBool(data.Value[keys[0]])
	//b.appKey = data.Value[keys[1]]
}

func initExpressAPI() {
	cli := inject.GetRegistryService()
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

}
