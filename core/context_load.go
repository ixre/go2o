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
	"go2o/core/domain/interface/content"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/user"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/personfinance"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/shopping"
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

func OrmMapping(conn db.Connector){
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

	/** new **/
	orm.Mapping(member.Level{}, "mm_level")
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
	orm.Mapping(sale.Item{}, "gs_item")
	orm.Mapping(sale.ValueGoods{}, "gs_goods")
	orm.Mapping(sale.Category{}, "gs_category")
	orm.Mapping(sale.GoodsSnapshot{}, "gs_snapshot")
	orm.Mapping(sale.Label{}, "gs_sale_label")
	orm.Mapping(sale.MemberPrice{}, "gs_member_price")

	/** 商户 **/
	orm.Mapping(merchant.Merchant{}, "mch_merchant")
	orm.Mapping(merchant.EnterpriseInfo{}, "mch_enterprise_info")
	orm.Mapping(merchant.ApiInfo{}, "mch_api_info")
	orm.Mapping(shop.Shop{}, "mch_shop")
	orm.Mapping(shop.OnlineShop{}, "mch_online_shop")
	orm.Mapping(shop.OfflineShop{}, "mch_offline_shop")
	orm.Mapping(merchant.SaleConf{}, "mch_sale_conf")
	orm.Mapping(merchant.MemberLevel{}, "pt_member_level")
	orm.Mapping(content.ValuePage{}, "mch_page")
	orm.Mapping(mss.MailTemplate{}, "pt_mail_template")
	orm.Mapping(mss.MailTask{}, "pt_mail_queue")

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
