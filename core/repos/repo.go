package repos

import (
	"log"
	"sync"

	"github.com/ixre/go2o/core/dao/model"
	"github.com/ixre/go2o/core/domain/interface/ad"
	afterSales "github.com/ixre/go2o/core/domain/interface/aftersales"
	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/domain/interface/delivery"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/merchant/user"
	"github.com/ixre/go2o/core/domain/interface/merchant/wholesaler"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/personfinance"
	promodel "github.com/ixre/go2o/core/domain/interface/pro_model"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : orm_mapping
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-21 15:12
 * description :
 * history :
 */

var (
	mux                 sync.Mutex
	DefaultCacheSeconds int64 = 3600
)

// 处理错误
func handleError(err error) error {
	return domain.HandleError(err, "rep")
	//if err != nil && gof.CurrentApp.Debug() {
	//	gof.CurrentApp.Log().Println("[ GO2O][ Repo][ Error] -", err.Error())
	//}
	//return err
}

// 删除指定前缀的缓存
func PrefixDel(sto storage.Interface, prefix string) error {
	_, err := sto.DeleteWith(prefix)
	if err != nil {
		log.Println("[ Cache][ Clean]: clean by prefix ", prefix, " error:", err)
	}
	return err
}

func OrmMapping(orm orm.Orm) {
	//table mapping
	orm.Mapping(valueobject.Area{}, "china_area")
	orm.Mapping(registry.Registry{}, "registry")
	// ad
	orm.Mapping(ad.Ad{}, "ad_list")
	orm.Mapping(ad.Image{}, "ad_image")
	orm.Mapping(ad.HyperLink{}, "ad_hyperlink")
	orm.Mapping(ad.Position{}, "ad_position")
	orm.Mapping(ad.AdUserSet{}, "ad_userset")

	// MSS
	orm.Mapping(mss.Message{}, "msg_list")
	orm.Mapping(mss.To{}, "msg_to")
	orm.Mapping(mss.Content{}, "msg_content")
	orm.Mapping(mss.Replay{}, "msg_replay")

	// 内容
	orm.Mapping(content.Page{}, "arc_page")
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
	orm.Mapping(member.ConsigneeAddress{}, "mm_deliver_addr")
	orm.Mapping(member.InviteRelation{}, "mm_relation")
	orm.Mapping(member.CerticationInfo{}, "mm_cert_info")
	orm.Mapping(member.Favorite{}, "mm_favorite")
	orm.Mapping(member.BankCard{}, "mm_bank_card")
	orm.Mapping(member.ReceiptsCode{}, "mm_receipts_code")
	orm.Mapping(member.LevelUpLog{}, "mm_levelup")
	orm.Mapping(member.BuyerGroup{}, "mm_buyer_group")
	orm.Mapping(member.MmLockInfo{}, "mm_lock_info")
	orm.Mapping(member.MmLockHistory{}, "mm_lock_history")

	// ORDER
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
	orm.Mapping(express.Provider{}, "express_provider")
	orm.Mapping(express.ExpressTemplate{}, "mch_express_template")
	orm.Mapping(express.RegionExpressTemplate{}, "express_area_set")

	// Shipment
	orm.Mapping(shipment.ShipmentOrder{}, "ship_order")
	orm.Mapping(shipment.ShipmentItem{}, "ship_item")

	// 产品
	orm.Mapping(product.Product{}, "product")
	orm.Mapping(item.GoodsItem{}, "item_info")
	orm.Mapping(item.Sku{}, "item_sku")
	orm.Mapping(product.Category{}, "product_category")
	orm.Mapping(promodel.ProductModel{}, "product_model")
	orm.Mapping(promodel.ProductBrand{}, "product_brand")
	orm.Mapping(promodel.ProModelBrand{}, "product_model_brand")
	orm.Mapping(promodel.Attr{}, "product_model_attr")
	orm.Mapping(promodel.AttrItem{}, "product_model_attr_item")
	orm.Mapping(promodel.Spec{}, "product_model_spec")
	orm.Mapping(promodel.SpecItem{}, "product_model_spec_item")
	orm.Mapping(product.AttrValue{}, "product_attr_info")
	orm.Mapping(item.Snapshot{}, "item_snapshot")
	orm.Mapping(item.TradeSnapshot{}, "item_trade_snapshot")
	orm.Mapping(item.Label{}, "gs_sale_label")
	orm.Mapping(item.MemberPrice{}, "gs_member_price")

	// 商户
	orm.Mapping(merchant.Merchant{}, "mch_merchant")
	orm.Mapping(merchant.ApiInfo{}, "mch_api_info")
	orm.Mapping(shop.OnlineShop{}, "mch_online_shop")
	orm.Mapping(shop.OfflineShop{}, "mch_offline_shop")
	orm.Mapping(merchant.SaleConf{}, "mch_sale_conf")
	orm.Mapping(merchant.TradeConf{}, "mch_trade_conf")
	orm.Mapping(merchant.MemberLevel{}, "pt_member_level")
	orm.Mapping(merchant.Account{}, "mch_account")
	orm.Mapping(merchant.BalanceLog{}, "mch_balance_log")
	orm.Mapping(merchant.MchDayChart{}, "mch_day_chart")
	orm.Mapping(merchant.MchSignUp{}, "mch_sign_up")
	orm.Mapping(merchant.MchBuyerGroupSetting{}, "mch_buyer_group")
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
	orm.Mapping(model.QrTemplate{}, "comm_qr_template")
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
