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
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/dto"
	"go2o/core/infrastructure"
	"go2o/core/variable"
	"log"
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
func (s *ShopQuery) QueryShopIdByHost(host string) (vendorId int32, shopId int32) {
	//  $ 获取合作商ID
	// $ hostname : 域名
	// *.wdian.net  二级域名
	// www.dc1.com  顶级域名

	var err error
	reg := s.getHostRegexp()
	if reg.MatchString(host) {
		matches := reg.FindAllStringSubmatch(host, 1)
		user := matches[0][1]
		err = s.Connector.QueryRow(`SELECT s.vendor_id,o.shop_id FROM mch_online_shop o
		    INNER JOIN mch_shop s ON s.id=o.shop_id WHERE o.alias= $1`, func(row *sql.Row) error {
			return row.Scan(&vendorId, &shopId)
		}, user)
	} else {
		err = s.Connector.ExecScalar(`SELECT shop_id FROM mch_online_shop WHERE host= $1`,
			&shopId, host)
	}
	if err != nil {
		gof.CurrentApp.Log().Error(err)
	}
	return vendorId, shopId
}

// 获取商户编号
func (s *ShopQuery) GetMerchantId(shopId int32) int32 {
	var vendorId int32
	s.Connector.ExecScalar(`SELECT vendor_id FROM mch_shop WHERE id= $1`, &vendorId, shopId)
	return vendorId
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

	err := s.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM mch_shop sp INNER JOIN mch_online_shop ol
    ON ol.shop_id=sp.id INNER JOIN mch_merchant mch ON mch.id=sp.vendor_id
    WHERE sp.state=%d AND mch.enabled = 1 %s`, shop.StateNormal, where), &total)

	var e []*dto.ListOnlineShop
	if total > 0 && err == nil {
		sql = fmt.Sprintf(`SELECT sp.id,sp.name,alias,host,ol.logo,sp.create_time
        FROM mch_shop sp INNER JOIN mch_online_shop ol
        ON ol.shop_id=sp.id INNER JOIN mch_merchant mch ON mch.id=sp.vendor_id
        WHERE sp.state=%d AND mch.enabled = 1 %s %s LIMIT $2 OFFSET $1`,
			shop.StateNormal, where, order)
		err = s.GetOrm().SelectByQuery(&e, sql, begin, (end - begin))
	}
	if err != nil {
		log.Println("[ Go2o][ Query][ Error]:", err.Error())
	}
	return total, e
}
