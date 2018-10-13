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
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/domain/interface/content"
	"go2o/src/core/domain/interface/delivery"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/partner/mss"
	"go2o/src/core/domain/interface/partner/user"
	"go2o/src/core/domain/interface/personfinance"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/variable"
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
	connector := db.NewConnector(driver, connStr, l, debug)

	//table mapping
	orm := connector.GetOrm()

	/** new **/
	orm.Mapping(member.ValueMember{}, "mm_member")
	orm.Mapping(member.IntegralLog{}, "mm_integral_log")
	orm.Mapping(member.AccountValue{}, "mm_account")
	orm.Mapping(member.DeliverAddress{}, "mm_deliver_addr")
	orm.Mapping(member.MemberRelation{}, "mm_relation")
	orm.Mapping(member.BalanceInfoValue{}, "mm_balance_info")

	orm.Mapping(member.BankInfo{}, "mm_bank")
	orm.Mapping(shopping.ValueOrder{}, "pt_order")
	orm.Mapping(shopping.OrderItem{}, "pt_order_item")
	orm.Mapping(shopping.OrderCoupon{}, "pt_order_coupon")
	orm.Mapping(shopping.OrderPromotionBind{}, "pt_order_pb")
	orm.Mapping(shopping.OrderLog{}, "pt_order_log")
	orm.Mapping(shopping.ValueCart{}, "sale_cart")
	orm.Mapping(shopping.ValueCartItem{}, "sale_cart_item")

	/** 销售 **/
	orm.Mapping(sale.ValueItem{}, "gs_item")
	orm.Mapping(sale.ValueGoods{}, "gs_goods")
	orm.Mapping(sale.ValueCategory{}, "gs_category")
	orm.Mapping(sale.GoodsSnapshot{}, "gs_snapshot")
	orm.Mapping(sale.ValueSaleTag{}, "gs_sale_tag")
	orm.Mapping(sale.MemberPrice{}, "gs_member_price")

	/** 商户 **/
	orm.Mapping(partner.ValuePartner{}, "pt_partner")
	orm.Mapping(partner.ApiInfo{}, "pt_api")
	orm.Mapping(partner.SiteConf{}, "pt_siteconf")
	orm.Mapping(partner.ValueShop{}, "pt_shop")
	orm.Mapping(partner.SaleConf{}, "pt_saleconf")
	orm.Mapping(valueobject.MemberLevel{}, "pt_member_level")
	orm.Mapping(content.ValuePage{}, "pt_page")
	orm.Mapping(ad.ValueAdvertisement{}, "pt_ad")
	orm.Mapping(ad.ValueImage{}, "pt_ad_image")
	orm.Mapping(mss.MailTemplate{}, "pt_mail_template")
	orm.Mapping(mss.MailTask{}, "pt_mail_queue")

	/** 促销 **/
	orm.Mapping(promotion.ValueCoupon{}, "pm_coupon")
	orm.Mapping(promotion.ValueCouponBind{}, "pm_coupon_bind")
	orm.Mapping(promotion.ValueCouponTake{}, "pm_coupon_take")
	orm.Mapping(promotion.ValuePromotion{}, "pm_info")
	orm.Mapping(promotion.ValueCashBack{}, "pm_cash_back")

	/** 配送 **/
	orm.Mapping(delivery.AreaValue{}, "dlv_area")
	orm.Mapping(delivery.CoverageValue{}, "dlv_coverage")
	orm.Mapping(delivery.PartnerDeliverBind{}, "dlv_partner_bind")

	/** 用户 **/
	orm.Mapping(user.RoleValue{}, "usr_role")
	orm.Mapping(user.PersonValue{}, "usr_person")
	orm.Mapping(user.CredentialValue{}, "usr_credential")

	orm.Mapping(personfinance.RiseInfoValue{}, "pf_riseinfo")
	orm.Mapping(personfinance.RiseDayInfo{}, "pf_riseday")
	orm.Mapping(personfinance.RiseLog{}, "pf_riselog")

	orm.Mapping(valueobject.Goods{}, "")

	return connector
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
