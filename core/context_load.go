/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package core

import (
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/crypto"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/log"
	"go2o/core/domain/interface/ad"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/content"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/personfinance"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/variable"
	"strconv"
	"time"
)

func getDb(c *gof.Config, debug bool, l log.ILogger) db.Connector {
	//数据库连接字符串
	//root@tcp(127.0.0.1:3306)/db_name?charset=utf8
	var connStr string
	driver := c.GetString(variable.DbDriver)
	dbCharset := c.GetString(variable.DbCharset)
	if dbCharset == "" {
		dbCharset = "utf8"
	}
	connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local",
		c.GetString(variable.DbUsr),
		c.GetString(variable.DbPwd),
		c.GetString(variable.DbServer),
		c.GetString(variable.DbPort),
		c.GetString(variable.DbName),
		dbCharset,
	)
	connector := db.NewSimpleConnector(driver, connStr, l, 5000, debug)
	OrmMapping(connector)
	return connector
}

func OrmMapping(conn db.Connector) {
	//table mapping
	orm := conn.GetOrm()
	orm.Mapping(valueobject.Area{}, "china_area")
	/* ad */
	orm.Mapping(ad.Ad{}, "ad_list")
	orm.Mapping(ad.Image{}, "ad_image")
	orm.Mapping(ad.HyperLink{}, "ad_hyperlink")
	orm.Mapping(ad.AdGroup{}, "ad_group")
	orm.Mapping(ad.AdPosition{}, "ad_position")
	orm.Mapping(ad.AdUserSet{}, "ad_userset")

	/* MSS */
	orm.Mapping(mss.Message{}, "msg_list")
	orm.Mapping(mss.To{}, "msg_to")
	orm.Mapping(mss.Content{}, "msg_content")
	orm.Mapping(mss.Replay{}, "msg_replay")

	/* 内容 */
	orm.Mapping(content.Page{}, "con_page")
	orm.Mapping(content.Article{}, "con_article")
	orm.Mapping(content.ArticleCategory{}, "con_article_category")

	/** new **/
	orm.Mapping(member.Level{}, "mm_level")
	orm.Mapping(member.Member{}, "mm_member")
	orm.Mapping(member.Profile{}, "mm_profile")
	orm.Mapping(member.IntegralLog{}, "mm_integral_log")
	orm.Mapping(member.BalanceLog{}, "mm_balance_log")
	orm.Mapping(member.PresentLog{}, "mm_present_log")
	orm.Mapping(member.Account{}, "mm_account")
	orm.Mapping(member.DeliverAddress{}, "mm_deliver_addr")
	orm.Mapping(member.Relation{}, "mm_relation")
	orm.Mapping(member.BalanceInfo{}, "mm_balance_info")
	orm.Mapping(member.TrustedInfo{}, "mm_trusted_info")
	orm.Mapping(member.Favorite{}, "mm_favorite")
	orm.Mapping(member.BankInfo{}, "mm_bank")

	//** ORDER **//

	orm.Mapping(order.Order{}, "sale_order")
	orm.Mapping(order.SubOrder{}, "sale_sub_order")

	//orm.Mapping(order.ValueOrder1{}, "pt_order")
	orm.Mapping(order.OrderItem{}, "sale_order_item")
	orm.Mapping(order.OrderCoupon{}, "pt_order_coupon")
	orm.Mapping(order.OrderPromotionBind{}, "pt_order_pb")
	orm.Mapping(order.OrderLog{}, "sale_order_log")
	orm.Mapping(cart.ValueCart{}, "sale_cart")
	orm.Mapping(cart.CartItem{}, "sale_cart_item")

	//** After Sales **/
	orm.Mapping(afterSales.AfterSalesOrder{}, "sale_after_order")
	orm.Mapping(afterSales.ReturnOrder{}, "sale_return")
	orm.Mapping(afterSales.ExchangeOrder{}, "sale_exchange")
	orm.Mapping(afterSales.RefundOrder{}, "sale_refund")

	//** Express **//
	orm.Mapping(express.ExpressProvider{}, "express_provider")
	orm.Mapping(express.ExpressTemplate{}, "express_template")
	orm.Mapping(express.ExpressAreaTemplate{}, "express_area_set")

	//** Shipment **/
	orm.Mapping(shipment.ShipmentOrder{}, "ship_order")
	orm.Mapping(shipment.Item{}, "ship_item")

	/** 销售 **/
	orm.Mapping(item.Item{}, "gs_item")
	orm.Mapping(goods.ValueGoods{}, "gs_goods")
	orm.Mapping(sale.Category{}, "gs_category")
	orm.Mapping(goods.Snapshot{}, "gs_snapshot")
	orm.Mapping(goods.SalesSnapshot{}, "gs_sales_snapshot")
	orm.Mapping(sale.Label{}, "gs_sale_label")
	orm.Mapping(goods.MemberPrice{}, "gs_member_price")

	/** 商户 **/
	orm.Mapping(merchant.Merchant{}, "mch_merchant")
	orm.Mapping(merchant.EnterpriseInfo{}, "mch_enterprise_info")
	orm.Mapping(merchant.ApiInfo{}, "mch_api_info")
	orm.Mapping(shop.Shop{}, "mch_shop")
	orm.Mapping(shop.OnlineShop{}, "mch_online_shop")
	orm.Mapping(shop.OfflineShop{}, "mch_offline_shop")
	orm.Mapping(merchant.SaleConf{}, "mch_sale_conf")
	orm.Mapping(merchant.MemberLevel{}, "pt_member_level")
	orm.Mapping(merchant.Account{}, "mch_account")
	orm.Mapping(merchant.BalanceLog{}, "mch_balance_log")
	orm.Mapping(merchant.MchDayChart{}, "mch_day_chart")
	orm.Mapping(merchant.MchSignUp{}, "mch_sign_up")
	orm.Mapping(mss.MailTemplate{}, "pt_mail_template")
	orm.Mapping(mss.MailTask{}, "pt_mail_queue")

	orm.Mapping(payment.PaymentOrder{}, "pay_order")

	/** 促销 **/
	orm.Mapping(promotion.ValueCoupon{}, "pm_coupon")
	orm.Mapping(promotion.ValueCouponBind{}, "pm_coupon_bind")
	orm.Mapping(promotion.ValueCouponTake{}, "pm_coupon_take")
	orm.Mapping(promotion.PromotionInfo{}, "pm_info")
	orm.Mapping(promotion.ValueCashBack{}, "pm_cash_back")

	/** 配送 **/
	orm.Mapping(delivery.AreaValue{}, "dlv_area")
	orm.Mapping(delivery.CoverageValue{}, "dlv_coverage")
	orm.Mapping(delivery.MerchantDeliverBind{}, "dlv_merchant_bind")

	/** 用户 **/
	orm.Mapping(user.RoleValue{}, "usr_role")
	orm.Mapping(user.PersonValue{}, "usr_person")
	orm.Mapping(user.CredentialValue{}, "usr_credential")

	orm.Mapping(personfinance.RiseInfoValue{}, "pf_riseinfo")
	orm.Mapping(personfinance.RiseDayInfo{}, "pf_riseday")
	orm.Mapping(personfinance.RiseLog{}, "pf_riselog")

	orm.Mapping(valueobject.Goods{}, "")
}

func initTemplate(c *gof.Config) *gof.Template {
	spam := crypto.Md5([]byte(strconv.Itoa(int(time.Now().Unix()))))[8:14]
	return &gof.Template{
		Init: func(m *gof.TemplateDataMap) {
			v := *m
			v["static_serve"] = c.GetString(variable.StaticServer)
			v["img_serve"] = c.GetString(variable.ImageServer)
			v["domain"] = c.GetString(variable.ServerDomain)
			v["version"] = c.GetString(variable.Version)
			v["spam"] = spam
		},
	}
}
