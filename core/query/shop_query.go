/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
	"database/sql"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"go2o/core/infrastructure"
	"go2o/core/variable"
	"regexp"
)

type ShopQuery struct {
	db.Connector
	gof.Storage
	commHostRegexp *regexp.Regexp
}

func NewShopQuery(c gof.App) *ShopQuery {
	return &ShopQuery{
		Connector: c.Db(),
		Storage:   c.Storage(),
	}
}

var ()

func (this *ShopQuery) getHostRegexp() *regexp.Regexp {
	if this.commHostRegexp == nil {
		this.commHostRegexp = regexp.MustCompile("([^\\.]+)." +
			infrastructure.GetApp().Config().GetString(variable.ServerDomain))
	}
	return this.commHostRegexp
}

// 根据主机查询商店编号
func (this *ShopQuery) QueryShopIdByHost(host string) (mchId int, shopId int) {
	//  $ 获取合作商ID
	// $ hostname : 域名
	// *.wdian.net  二级域名
	// www.dc1.com  顶级域名

	var err error

	reg := this.getHostRegexp()
	if reg.MatchString(host) {
		matches := reg.FindAllStringSubmatch(host, 1)
		usr := matches[0][1]
		err = this.Connector.QueryRow(`SELECT s.mch_id,o.shop_id FROM mch_online_shop o
		    INNER JOIN mch_shop s ON s.id=o.shop_id WHERE o.alias=?`, func(row *sql.Row) {
			row.Scan(&mchId, &shopId)
		}, usr)
	} else {
		err = this.Connector.ExecScalar(
			`SELECT id FROM mch_merchant INNER JOIN pt_siteconf
					 ON pt_siteconf.merchant_id = mch_merchant.id
					 WHERE host=?`, &shopId, host)
	}
	if err != nil {
		gof.CurrentApp.Log().Error(err)
	}
	return mchId, shopId
}

// 获取商户编号
func (this *ShopQuery) GetMerchantId(shopId int) int {
	var mchId int
	this.Connector.ExecScalar(`SELECT mch_id FROM mch_shop WHERE id=?`, &mchId, shopId)
	return mchId
}
