/**
 * Copyright 2015 @ to2.net.
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
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/log"
	"go2o/core/dao/model"
	"go2o/core/domain/interface/ad"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/content"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
	"go2o/core/domain/interface/merchant/wholesaler"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/personfinance"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/interface/wallet"
	"go2o/core/dto"
	"go2o/core/module/express/kdniao"
	"go2o/core/msq"
	"go2o/core/service"
	"go2o/core/service/proto"
	"go2o/core/variable"
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
	gob.Register(ad.ValueGallery{})
	gob.Register(ad.Ad{})
	gob.Register([]*valueobject.Goods{})
	gob.Register(valueobject.Goods{})
	gob.Register(ad.HyperLink{})
	gob.Register(ad.Image{})
}

func Init(a *AppImpl, debug, trace bool) bool {
	a._debugMode = debug
	if trace {
		a.Db().GetOrm().SetTrace(a._debugMode)
	}
	OrmMapping(a.Db())
	// 初始化变量
	variable.Domain = a._config.GetString(variable.ServerDomain)
	a.Loaded = true
	for _, f := range startJobs {
		f()
	}
	return true
}

func OrmMapping(conn db.Connector) {
	//table mapping
	orm := conn.GetOrm()
	orm.Mapping(valueobject.Area{}, "china_area")
	orm.Mapping(registry.Registry{}, "registry")
	// ad
	orm.Mapping(ad.Ad{}, "ad_list")
	orm.Mapping(ad.Image{}, "ad_image")
	orm.Mapping(ad.HyperLink{}, "ad_hyperlink")
	orm.Mapping(ad.AdGroup{}, "ad_group")
	orm.Mapping(ad.AdPosition{}, "ad_position")
	orm.Mapping(ad.AdUserSet{}, "ad_userset")

	// MSS
	orm.Mapping(mss.Message{}, "msg_list")
	orm.Mapping(mss.To{}, "msg_to")
	orm.Mapping(mss.Content{}, "msg_content")
	orm.Mapping(mss.Replay{}, "msg_replay")

	// 内容
	orm.Mapping(content.Page{}, "ex_page")
	orm.Mapping(content.Article{}, "article_list")
	orm.Mapping(content.ArticleCategory{}, "article_category")

	// new
	orm.Mapping(member.Level{}, "mm_level")
	orm.Mapping(member.Member{}, "mm_member")
	orm.Mapping(member.Profile{}, "mm_profile")
	orm.Mapping(member.IntegralLog{}, "mm_integral_log")
	orm.Mapping(member.BalanceLog{}, "mm_balance_log")
	orm.Mapping(member.WalletAccountLog{}, "mm_wallet_log")
	orm.Mapping(member.FlowAccountLog{}, "mm_flow_log")
	orm.Mapping(member.Account{}, "mm_account")
	orm.Mapping(member.Address{}, "mm_deliver_addr")
	orm.Mapping(member.InviteRelation{}, "mm_relation")
	orm.Mapping(member.TrustedInfo{}, "mm_trusted_info")
	orm.Mapping(member.Favorite{}, "mm_favorite")
	orm.Mapping(member.BankInfo{}, "mm_bank")
	orm.Mapping(member.ReceiptsCode{}, "mm_receipts_code")
	orm.Mapping(member.LevelUpLog{}, "mm_levelup")
	orm.Mapping(member.BuyerGroup{}, "mm_buyer_group")
	orm.Mapping(member.MmLockInfo{}, "mm_lock_info")
	orm.Mapping(member.MmLockHistory{}, "mm_lock_history")

	// ORDER
	orm.Mapping(order.NormalOrder{}, "sale_order")
	orm.Mapping(order.NormalSubOrder{}, "sale_sub_order")
	orm.Mapping(order.SubOrderItem{}, "sale_order_item")
	orm.Mapping(order.OrderCoupon{}, "pt_order_coupon")
	orm.Mapping(order.OrderPromotionBind{}, "pt_order_pb")
	orm.Mapping(order.OrderLog{}, "sale_order_log")
	orm.Mapping(cart.NormalCart{}, "sale_cart")
	orm.Mapping(cart.NormalCartItem{}, "sale_cart_item")
	orm.Mapping(order.Order{}, "order_list")
	orm.Mapping(order.WholesaleOrder{}, "order_wholesale_order")
	orm.Mapping(order.WholesaleItem{}, "order_wholesale_item")
	orm.Mapping(order.TradeOrder{}, "order_trade_order")

	// After Sales
	orm.Mapping(afterSales.AfterSalesOrder{}, "sale_after_order")
	orm.Mapping(afterSales.ReturnOrder{}, "sale_return")
	orm.Mapping(afterSales.ExchangeOrder{}, "sale_exchange")
	orm.Mapping(afterSales.RefundOrder{}, "sale_refund")

	// Express
	orm.Mapping(express.ExpressProvider{}, "express_provider")
	orm.Mapping(express.ExpressTemplate{}, "mch_express_template")
	orm.Mapping(express.ExpressAreaTemplate{}, "express_area_set")

	// Shipment
	orm.Mapping(shipment.ShipmentOrder{}, "ship_order")
	orm.Mapping(shipment.Item{}, "ship_item")

	// 产品
	orm.Mapping(product.Product{}, "product")
	orm.Mapping(item.GoodsItem{}, "item_info")
	orm.Mapping(item.Sku{}, "item_sku")
	orm.Mapping(product.Category{}, "product_category")
	orm.Mapping(promodel.ProModel{}, "product_model")
	orm.Mapping(promodel.ProBrand{}, "product_brand")
	orm.Mapping(promodel.ProModelBrand{}, "product_model_brand")
	orm.Mapping(promodel.Attr{}, "product_model_attr")
	orm.Mapping(promodel.AttrItem{}, "product_model_attr_item")
	orm.Mapping(promodel.Spec{}, "product_model_spec")
	orm.Mapping(promodel.SpecItem{}, "product_model_spec_item")
	orm.Mapping(product.Attr{}, "product_attr_info")
	orm.Mapping(item.Snapshot{}, "item_snapshot")
	orm.Mapping(item.TradeSnapshot{}, "item_trade_snapshot")
	orm.Mapping(item.Label{}, "gs_sale_label")
	orm.Mapping(item.MemberPrice{}, "gs_member_price")

	// 商户
	orm.Mapping(merchant.Merchant{}, "mch_merchant")
	orm.Mapping(merchant.EnterpriseInfo{}, "mch_enterprise_info")
	orm.Mapping(merchant.ApiInfo{}, "mch_api_info")
	orm.Mapping(shop.Shop{}, "mch_shop")
	orm.Mapping(shop.OnlineShop{}, "mch_online_shop")
	orm.Mapping(shop.OfflineShop{}, "mch_offline_shop")
	orm.Mapping(merchant.SaleConf{}, "mch_sale_conf")
	orm.Mapping(merchant.TradeConf{}, "mch_trade_conf")
	orm.Mapping(merchant.MemberLevel{}, "pt_member_level")
	orm.Mapping(merchant.Account{}, "mch_account")
	orm.Mapping(merchant.BalanceLog{}, "mch_balance_log")
	orm.Mapping(merchant.MchDayChart{}, "mch_day_chart")
	orm.Mapping(merchant.MchSignUp{}, "mch_sign_up")
	orm.Mapping(merchant.MchBuyerGroup{}, "mch_buyer_group")
	orm.Mapping(mss.MailTemplate{}, "pt_mail_template")
	orm.Mapping(mss.MailTask{}, "pt_mail_queue")

	// 批发
	orm.Mapping(wholesaler.WsWholesaler{}, "ws_wholesaler")
	orm.Mapping(wholesaler.WsRebateRate{}, "ws_rebate_rate")
	orm.Mapping(item.WsItem{}, "ws_item")
	orm.Mapping(item.WsItemDiscount{}, "ws_item_discount")
	orm.Mapping(item.WsSkuPrice{}, "ws_sku_price")
	orm.Mapping(cart.WsCart{}, "ws_cart")
	orm.Mapping(cart.WsCartItem{}, "ws_cart_item")

	// 支付
	orm.Mapping(payment.Order{}, "pay_order")
	orm.Mapping(payment.TradeMethodData{}, "pay_trade_data")
	orm.Mapping(payment.MergeOrder{}, "pay_merge_order")
	orm.Mapping(payment.PaySpTrade{}, "pay_sp_trade")

	// 促销
	orm.Mapping(promotion.ValueCoupon{}, "pm_coupon")
	orm.Mapping(promotion.ValueCouponBind{}, "pm_coupon_bind")
	orm.Mapping(promotion.ValueCouponTake{}, "pm_coupon_take")
	orm.Mapping(promotion.PromotionInfo{}, "pm_info")
	orm.Mapping(promotion.ValueCashBack{}, "pm_cash_back")

	// 配送
	orm.Mapping(delivery.AreaValue{}, "dlv_area")
	orm.Mapping(delivery.CoverageValue{}, "dlv_coverage")
	orm.Mapping(delivery.MerchantDeliverBind{}, "dlv_merchant_bind")

	// 用户
	orm.Mapping(user.RoleValue{}, "user_role")
	orm.Mapping(user.PersonValue{}, "user_person")
	orm.Mapping(user.CredentialValue{}, "user_credential")

	orm.Mapping(personfinance.RiseInfoValue{}, "pf_riseinfo")
	orm.Mapping(personfinance.RiseDayInfo{}, "pf_riseday")
	orm.Mapping(personfinance.RiseLog{}, "pf_riselog")

	// 通用模块
	orm.Mapping(model.CommQrTemplate{}, "comm_qr_template")
	orm.Mapping(model.PortalNav{}, "portal_nav")
	orm.Mapping(model.PortalNavType{}, "portal_nav_type")
	orm.Mapping(model.PortalFloorAd{}, "portal_floor_ad")
	orm.Mapping(model.PortalFloorLink{}, "portal_floor_link")
	orm.Mapping(valueobject.Goods{}, "")

	// 钱包
	orm.Mapping(wallet.Wallet{}, "wal_wallet")
	orm.Mapping(wallet.WalletLog{}, "wal_wallet_log")
	// KV
	orm.Mapping(valueobject.SysKeyValue{}, "sys_kv")
}

func AppDispose() {
	//GetRedisPool().Close()
	msq.Close()
}

func InitialModules() {
	initExpressAPI()
	initBankB4eAPI()
	initSSOModule()
}

func initSSOModule() {
	//domain := variable.Domain
	trans, _, err := service.RegistryServeClient()
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
			registries, _ := cli.GetRegistries(context.TODO(),&proto.StringArray{Value:  keys})
			_, _ = s.Register(&proto.SSsoApp{
				ID:   1,
				Name: "RetailPortal",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[0]], domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				ID:   2,
				Name: "WholesalePortal",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[1]], domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				ID:   3,
				Name: "HApi",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[2]], domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				ID:   4,
				Name: "Member",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[3]], domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				ID:   5,
				Name: "MemberMobile",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[4]],
					domain),
			})
			_, _ = s.Register(&proto.SSsoApp{
				ID:   6,
				Name: "RetailPortalMobile",
				ApiUrl: fmt.Sprintf("//%s%s/user/sync_m.p",
					registries.Value[keys[5]], domain),
			})

		*/
	}
}

func initBankB4eAPI() {
	trans, cli, err := service.RegistryServeClient()
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

		//data, _ := cli.GetRegistries(ctx, &proto.StringArray{Value: keys})
		//b.open, _ = strconv.ParseBool(data.Value[keys[0]])
		//b.appKey = data.Value[keys[1]]
	}
}

func initExpressAPI() {
	trans, cli, err := service.RegistryServeClient()
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
		data, _ := cli.GetRegistries(context.TODO(), &proto.StringArray{Value: keys})
		kdniao.EBusinessID = data.Value[keys[0]]
		kdniao.AppKey = data.Value[keys[1]]
	} else {
		log.Println("intialize express module error:", err.Error())
		os.Exit(1)
	}
}
