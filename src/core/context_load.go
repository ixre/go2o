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
	connector := db.NewSimpleConnector(driver, connStr, l, 5000, debug)

	//table mapping
	orm := connector.GetOrm()

	/** new **/
	orm.TableMapping(member.ValueMember{}, "mm_member")
	orm.TableMapping(member.IntegralLog{}, "mm_integral_log")
	orm.TableMapping(member.AccountValue{}, "mm_account")
	orm.TableMapping(member.DeliverAddress{}, "mm_deliver_addr")
	orm.TableMapping(member.MemberRelation{}, "mm_relation")
	orm.TableMapping(member.BalanceInfoValue{}, "mm_balance_info")

	orm.TableMapping(member.BankInfo{}, "mm_bank")
	orm.TableMapping(shopping.ValueOrder{}, "pt_order")
	orm.TableMapping(shopping.OrderItem{}, "pt_order_item")
	orm.TableMapping(shopping.OrderCoupon{}, "pt_order_coupon")
	orm.TableMapping(shopping.OrderPromotionBind{}, "pt_order_pb")
	orm.TableMapping(shopping.OrderLog{}, "pt_order_log")
	orm.TableMapping(shopping.ValueCart{}, "sale_cart")
	orm.TableMapping(shopping.ValueCartItem{}, "sale_cart_item")

	/** 销售 **/
	orm.TableMapping(sale.ValueItem{}, "gs_item")
	orm.TableMapping(sale.ValueGoods{}, "gs_goods")
	orm.TableMapping(sale.ValueCategory{}, "gs_category")
	orm.TableMapping(sale.GoodsSnapshot{}, "gs_snapshot")
	orm.TableMapping(sale.ValueSaleTag{}, "gs_sale_tag")
	orm.TableMapping(sale.MemberPrice{}, "gs_member_price")

	/** 商户 **/
	orm.TableMapping(partner.ValuePartner{}, "pt_partner")
	orm.TableMapping(partner.ApiInfo{}, "pt_api")
	orm.TableMapping(partner.SiteConf{}, "pt_siteconf")
	orm.TableMapping(partner.ValueShop{}, "pt_shop")
	orm.TableMapping(partner.SaleConf{}, "pt_saleconf")
	orm.TableMapping(valueobject.MemberLevel{}, "pt_member_level")
	orm.TableMapping(content.ValuePage{}, "pt_page")
	orm.TableMapping(ad.ValueAdvertisement{}, "pt_ad")
	orm.TableMapping(ad.ValueImage{}, "pt_ad_image")
	orm.TableMapping(mss.MailTemplate{}, "pt_mail_template")
	orm.TableMapping(mss.MailTask{}, "pt_mail_queue")

	/** 促销 **/
	orm.TableMapping(promotion.ValueCoupon{}, "pm_coupon")
	orm.TableMapping(promotion.ValueCouponBind{}, "pm_coupon_bind")
	orm.TableMapping(promotion.ValueCouponTake{}, "pm_coupon_take")
	orm.TableMapping(promotion.ValuePromotion{}, "pm_info")
	orm.TableMapping(promotion.ValueCashBack{}, "pm_cash_back")

	/** 配送 **/
	orm.TableMapping(delivery.AreaValue{}, "dlv_area")
	orm.TableMapping(delivery.CoverageValue{}, "dlv_coverage")
	orm.TableMapping(delivery.PartnerDeliverBind{}, "dlv_partner_bind")

	/** 用户 **/
	orm.TableMapping(user.RoleValue{}, "usr_role")
	orm.TableMapping(user.PersonValue{}, "usr_person")
	orm.TableMapping(user.CredentialValue{}, "usr_credential")

	orm.TableMapping(personfinance.RiseInfoValue{}, "pf_riseinfo")
	orm.TableMapping(personfinance.RiseDayInfo{}, "pf_dayrise")
	orm.TableMapping(personfinance.RiseLog{}, "pf_riselog")

	orm.TableMapping(valueobject.Goods{}, "")

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
