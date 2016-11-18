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
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/storage"
	"go2o/core/dto"
	"go2o/core/infrastructure"
	"go2o/core/variable"
	"regexp"
)

type ShopQuery struct {
	db.Connector
	Storage        storage.Interface
	commHostRegexp *regexp.Regexp
}

func NewShopQuery(c gof.App) *ShopQuery {
	return &ShopQuery{
		Connector: c.Db(),
		Storage:   c.Storage(),
	}
}

func (s *ShopQuery) getHostRegexp() *regexp.Regexp {
	if s.commHostRegexp == nil {
		s.commHostRegexp = regexp.MustCompile("([^\\.]+)." +
			infrastructure.GetApp().Config().GetString(variable.ServerDomain))
	}
	return s.commHostRegexp
}

// 根据主机查询商店编号
func (s *ShopQuery) QueryShopIdByHost(host string) (mchId int32, shopId int) {
	//  $ 获取合作商ID
	// $ hostname : 域名
	// *.wdian.net  二级域名
	// www.dc1.com  顶级域名

	var err error

	reg := s.getHostRegexp()
	if reg.MatchString(host) {
		matches := reg.FindAllStringSubmatch(host, 1)
		usr := matches[0][1]
		err = s.Connector.QueryRow(`SELECT s.mch_id,o.shop_id FROM mch_online_shop o
		    INNER JOIN mch_shop s ON s.id=o.shop_id WHERE o.alias=?`, func(row *sql.Row) {
			row.Scan(&mchId, &shopId)
		}, usr)
	} else {
		err = s.Connector.ExecScalar(
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
func (s *ShopQuery) GetMerchantId(shopId int) int {
	var mchId int
	s.Connector.ExecScalar(`SELECT mch_id FROM mch_shop WHERE id=?`, &mchId, shopId)
	return mchId
}

// 获取营业中的店铺列表
func (s *ShopQuery) PagedOnBusinessOnlineShops(begin, end int, where string,
	order string) (int, []*dto.ListOnlineShop) {
	var sql string
	total := 0
	if len(where) != 0 {
		where = " AND " + where
	}
	if len(order) != 0 {
		order = "  ORDER BY " + order
	}
	s.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM mch_shop sp INNER JOIN mch_online_shop ol
    ON ol.shop_id=sp.id INNER JOIN mch_merchant mch ON mch.id=sp.mch_id
    WHERE sp.state=2 AND mch.enabled = 1 %s`, where), &total)

	e := []*dto.ListOnlineShop{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT sp.id,sp.name,alias,host,ol.logo,sp.create_time
        FROM mch_shop sp INNER JOIN mch_online_shop ol
        ON ol.shop_id=sp.id INNER JOIN mch_merchant mch ON mch.id=sp.mch_id
        WHERE sp.state=2 AND mch.enabled = 1 %s %s LIMIT ?,?`,
			where, order)
		s.GetOrm().SelectByQuery(&e, sql, begin, (end - begin))
	}
	return total, e
}
