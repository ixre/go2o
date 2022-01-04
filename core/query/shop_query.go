/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
	"fmt"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"log"
	"regexp"
)

type ShopQuery struct {
	db.Connector
	o              orm.Orm
	Storage        storage.Interface
	commHostRegexp *regexp.Regexp
}

func NewShopQuery(o orm.Orm, s storage.Interface) *ShopQuery {
	return &ShopQuery{
		Connector: o.Connector(),
		o:         o,
		Storage:   s,
	}
}


// QueryShopIdByHost 根据主机查询商店编号
func (s *ShopQuery) QueryShopIdByHost(host string) (shopId int64) {
	err := s.Connector.ExecScalar(`SELECT id FROM mch_online_shop WHERE (host= $1 OR alias = $1)`,
		&shopId, host)
	if err != nil {
		log.Println(err.Error())
	}
	return shopId
}

// GetMerchantId 获取商户编号
func (s *ShopQuery) GetMerchantId(shopId int64) int64 {
	var vendorId int64
	s.Connector.ExecScalar(`SELECT vendor_id FROM mch_shop WHERE id= $1`, &vendorId, shopId)
	return vendorId
}

// PagedOnBusinessOnlineShops 获取营业中的店铺列表
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

	err := s.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM mch_online_shop sp
 	INNER JOIN mch_merchant mch ON mch.id=sp.vendor_id
    WHERE sp.state=%d AND mch.enabled = 1 %s`, shop.StateNormal, where), &total)

	var e = make([]*dto.ListOnlineShop, 0)
	if total > 0 && err == nil {
		sql = fmt.Sprintf(`SELECT sp.id,sp.shop_name,alias,host,sp.logo,sp.create_time
        FROM  mch_online_shop sp INNER JOIN mch_merchant mch ON mch.id=sp.vendor_id
        WHERE sp.state=%d AND mch.enabled = 1 %s %s LIMIT $2 OFFSET $1`,
			shop.StateNormal, where, order)
		err = s.o.SelectByQuery(&e, sql, begin, end-begin)
	}
	if err != nil {
		log.Println("[ Go2o][ Params][ Error]:", err.Error())
	}
	return total, e
}
