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
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/db"
	"github.com/jrsix/gof/log"
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/partner/mss"
	"go2o/src/core/domain/interface/partner/user"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/variable"
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

	//table mapping
	orm := connector.GetOrm()

	/** new **/
	orm.CreateTableMap(member.ValueMember{}, "mm_member")
	orm.CreateTableMap(member.IncomeLog{}, "mm_income_log")
	orm.CreateTableMap(member.IntegralLog{}, "mm_integral_log")
	orm.CreateTableMap(member.AccountValue{}, "mm_account")
	orm.CreateTableMap(member.DeliverAddress{}, "mm_deliver_addr")
	orm.CreateTableMap(member.MemberRelation{}, "mm_relation")
	orm.CreateTableMap(member.BalanceInfoValue{}, "mm_balance_info")

	orm.CreateTableMap(member.BankInfo{}, "mm_bank")
	orm.CreateTableMap(shopping.ValueOrder{}, "pt_order")
	orm.CreateTableMap(shopping.OrderItem{}, "pt_order_item")
	orm.CreateTableMap(shopping.OrderCoupon{}, "pt_order_coupon")
	orm.CreateTableMap(shopping.OrderPromotionBind{}, "pt_order_pb")
	orm.CreateTableMap(shopping.OrderLog{}, "pt_order_log")
	orm.CreateTableMap(shopping.ValueCart{}, "sale_cart")
	orm.CreateTableMap(shopping.ValueCartItem{}, "sale_cart_item")

	/** 销售 **/
	orm.CreateTableMap(sale.ValueItem{}, "gs_item")
	orm.CreateTableMap(sale.ValueGoods{}, "gs_goods")
	orm.CreateTableMap(sale.ValueCategory{}, "gs_category")
	orm.CreateTableMap(sale.GoodsSnapshot{}, "gs_snapshot")
	orm.CreateTableMap(sale.ValueSaleTag{}, "gs_sale_tag")
	orm.CreateTableMap(sale.MemberPrice{}, "gs_member_price")

	/** 商户 **/
	orm.CreateTableMap(partner.ValuePartner{}, "pt_partner")
	orm.CreateTableMap(partner.ApiInfo{}, "pt_api")
	orm.CreateTableMap(partner.SiteConf{}, "pt_siteconf")
	orm.CreateTableMap(partner.ValueShop{}, "pt_shop")
	orm.CreateTableMap(partner.SaleConf{}, "pt_saleconf")
	orm.CreateTableMap(valueobject.MemberLevel{}, "pt_member_level")
	orm.CreateTableMap(content.ValuePage{}, "pt_page")
	orm.CreateTableMap(ad.ValueAdvertisement{}, "pt_ad")
	orm.CreateTableMap(ad.ValueImage{}, "pt_ad_image")
	orm.CreateTableMap(mss.MailTemplate{}, "pt_mail_template")
	orm.CreateTableMap(mss.MailTask{}, "pt_mail_queue")

	/** 促销 **/
	orm.CreateTableMap(promotion.ValueCoupon{}, "pm_coupon")
	orm.CreateTableMap(promotion.ValueCouponBind{}, "pm_coupon_bind")
	orm.CreateTableMap(promotion.ValueCouponTake{}, "pm_coupon_take")
	orm.CreateTableMap(promotion.ValuePromotion{}, "pm_info")
	orm.CreateTableMap(promotion.ValueCashBack{}, "pm_cash_back")

	/** 配送 **/
	orm.CreateTableMap(delivery.AreaValue{}, "dlv_area")
	orm.CreateTableMap(delivery.CoverageValue{}, "dlv_coverage")
	orm.CreateTableMap(delivery.PartnerDeliverBind{}, "dlv_partner_bind")

	/** 用户 **/
	orm.CreateTableMap(user.RoleValue{}, "usr_role")
	orm.CreateTableMap(user.PersonValue{}, "usr_person")
	orm.CreateTableMap(user.CredentialValue{}, "usr_credential")

	orm.CreateTableMap(valueobject.Goods{}, "")

	return connector
}

func initTemplate(c *gof.Config) *gof.Template {
	return &gof.Template{
		Init: func(m *gof.TemplateDataMap) {
			v := *m
			v["static_serve"] = c.GetString(variable.StaticServer)
			v["img_serve"] = c.GetString(variable.ImageServer)
			v["domain"] = c.GetString(variable.ServerDomain)
			v["version"] = c.GetString(variable.Version)
		},
	}
}
