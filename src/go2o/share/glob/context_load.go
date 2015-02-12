/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package glob

import (
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/db"
	"github.com/atnet/gof/log"
	"github.com/atnet/gof/web"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/partner"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/shopping"
	"go2o/core/infrastructure/alipay"
	"go2o/share/variable"
)

func getDb(c *gof.Config, l log.ILogger) db.Connector {
	//数据库连接字符串
	//root@tcp(127.0.0.1:3306)/foodording?charset=utf8
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
	connector := db.NewCommonConnector(driver, connStr, l, 30)

	//table mapping
	orm := connector.GetOrm()

	/** new **/
	orm.CreateTableMap(member.ValueMember{}, "mm_member")
	orm.CreateTableMap(member.IncomeLog{}, "mm_income_log")
	orm.CreateTableMap(member.IntegralLog{}, "mm_integral_log")
	orm.CreateTableMap(member.Account{}, "mm_account")
	orm.CreateTableMap(member.MemberLevel{}, "conf_member_level")
	orm.CreateTableMap(member.DeliverAddress{}, "mm_deliver_addr")
	orm.CreateTableMap(member.MemberRelation{}, "mm_relation")

	orm.CreateTableMap(member.BankInfo{}, "mm_bank")
	orm.CreateTableMap(shopping.ValueOrder{}, "pt_order")
	orm.CreateTableMap(shopping.OrderCoupon{}, "pt_order_coupon")
	orm.CreateTableMap(shopping.OrderCoupon{}, "pt_order_coupon")
	orm.CreateTableMap(shopping.OrderLog{}, "pt_order_log")
	orm.CreateTableMap(shopping.ValueCart{}, "sale_cart")
	orm.CreateTableMap(shopping.ValueCartItem{}, "sale_cart_item")

	orm.CreateTableMap(sale.ValueGoods{}, "gs_goods")
	orm.CreateTableMap(sale.ValueCategory{}, "gs_category")
	orm.CreateTableMap(sale.GoodsSnapshot{}, "gs_snapshot")

	orm.CreateTableMap(partner.ValuePartner{}, "pt_partner")
	orm.CreateTableMap(partner.SiteConf{}, "pt_siteconf")
	orm.CreateTableMap(partner.ValueShop{}, "pt_shop")
	orm.CreateTableMap(partner.SaleConf{}, "pt_saleconf")

	/** 促销 **/
	orm.CreateTableMap(promotion.ValueCoupon{}, "pm_coupon")
	orm.CreateTableMap(promotion.ValueCouponBind{}, "pm_coupon_bind")
	orm.CreateTableMap(promotion.ValueCouponTake{}, "pm_coupon_take")

	return connector
}

func initTemplate(c *gof.Config) *web.TemplateWrapper {
	return &web.TemplateWrapper{
		Init: func(m *map[string]interface{}) {
			v := *m
			v["static_serv"] = c.GetString(variable.StaticServer)
			v["img_serv"] = c.GetString(variable.ImageServer)
			v["domain"] = c.GetString(variable.ServerDomain)
			v["version"] = c.GetString(variable.Version)
		},
	}
}

func paymentCfg(c *gof.Config) {
	alipay.Configure(c.GetString(variable.Alipay_Partner),
		c.GetString(variable.Alipay_Key),
		c.GetString(variable.Alipay_Seller))
}

//	MasterToken: crypto.EncodeUsrPwd("master", "123456"),
//		deviceDir := filepath.Dir(options.DevicePatchDir)
//		fi,err := os.Stat(options.DevicePatchDir)
//
//		if fi == nil || err == os.ErrNotExist {
//			os.MkdirAll(deviceDir,os.ModePerm)
//		}else if !fi.IsDir() {
//			os.Remove(options.DevicePatchDir)
//			os.MkdirAll(deviceDir,os.ModePerm)
//}
