/**
 * Copyright 2015 @ z3q.net.
 * name : types.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package core

import (
	"encoding/gob"
	"github.com/jsix/gof/db"
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
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/interface/wallet"
	"go2o/core/dto"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/variable"
	"strings"
)

func init() {
	registerTypes()
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
	gob.Register(&member.Relation{})
	gob.Register(&dto.ListOnlineShop{})
	gob.Register([]*dto.ListOnlineShop{})
	gob.Register(&member_service.SMember{})
	gob.Register(&member_service.SProfile{})
}

func Init(a *AppImpl, debug, trace bool) bool {
	a._debugMode = debug
	if trace {
		a.Db().GetOrm().SetTrace(a._debugMode)
	}
	OrmMapping(a.Db())
	// 初始化变量
	variable.Domain = a._config.GetString(variable.ServerDomain)
	confPath := strings.TrimSpace(a._config.GetString("conf_path"))
	if len(confPath) > 0 {
		variable.ConfPath = confPath
	}
	a.Loaded = true
	return true
}

func OrmMapping(conn db.Connector) {
	//table mapping
	orm := conn.GetOrm()
	orm.Mapping(valueobject.Area{}, "china_area")

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
	orm.Mapping(member.MWalletLog{}, "mm_wallet_log")
	orm.Mapping(member.Account{}, "mm_account")
	orm.Mapping(member.Address{}, "mm_deliver_addr")
	orm.Mapping(member.Relation{}, "mm_relation")
	orm.Mapping(member.BalanceInfo{}, "mm_balance_info")
	orm.Mapping(member.TrustedInfo{}, "mm_trusted_info")
	orm.Mapping(member.Favorite{}, "mm_favorite")
	orm.Mapping(member.BankInfo{}, "mm_bank")
	orm.Mapping(member.LevelUpLog{}, "mm_levelup")
	orm.Mapping(member.BuyerGroup{}, "mm_buyer_group")

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
	orm.Mapping(express.ExpressTemplate{}, "express_template")
	orm.Mapping(express.ExpressAreaTemplate{}, "express_area_set")

	// Shipment
	orm.Mapping(shipment.ShipmentOrder{}, "ship_order")
	orm.Mapping(shipment.Item{}, "ship_item")

	// 产品
	orm.Mapping(product.Product{}, "pro_product")
	orm.Mapping(item.GoodsItem{}, "item_info")
	orm.Mapping(item.Sku{}, "item_sku")
	orm.Mapping(product.Category{}, "pro_category")
	orm.Mapping(promodel.ProModel{}, "pro_model")
	orm.Mapping(promodel.ProModelBrand{}, "pro_model_brand")
	orm.Mapping(promodel.ProBrand{}, "pro_brand")
	orm.Mapping(promodel.Attr{}, "pro_attr")
	orm.Mapping(promodel.AttrItem{}, "pro_attr_item")
	orm.Mapping(promodel.Spec{}, "pro_spec")
	orm.Mapping(promodel.SpecItem{}, "pro_spec_item")
	orm.Mapping(product.Attr{}, "pro_attr_info")
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
	orm.Mapping(user.RoleValue{}, "usr_role")
	orm.Mapping(user.PersonValue{}, "usr_person")
	orm.Mapping(user.CredentialValue{}, "usr_credential")

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
