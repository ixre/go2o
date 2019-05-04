/**
 * Copyright 2015 @ z3q.net.
 * name : aftersales_query
 * author : jarryliu
 * date : 2016-07-18 19:27
 * description :
 * history :
 */
package query

import (
	"database/sql"
	"github.com/ixre/gof/db"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/dto"
	"go2o/core/infrastructure/format"
)

type AfterSalesQuery struct {
	db.Connector
}

func NewAfterSalesQuery(db db.Connector) *AfterSalesQuery {
	return &AfterSalesQuery{
		Connector: db,
	}
}

// 获取分页售后单
func (a *AfterSalesQuery) QueryPagerAfterSalesOrderOfMember(memberId int64, begin,
	size int, where string) (int, []*dto.PagedMemberAfterSalesOrder) {
	list := []*dto.PagedMemberAfterSalesOrder{}
	total := 0
	if len(where) > 0 {
		where = " AND " + where
	}
	a.ExecScalar(`SELECT COUNT(0) FROM sale_after_order ao
	INNER JOIN sale_sub_order so ON so.id=ao.order_id
	INNER JOIN mch_merchant mch ON so.vendor_id = mch.id
	INNER JOIN item_trade_snapshot sn ON sn.id = ao.snap_id
	WHERE ao.buyer_id= $1 `+where, &total, memberId)
	if total > 0 {
		a.Query(`SELECT ao.id,ao.type,so.order_no,so.vendor_id,mch.name as vendor_name,
 ao.snap_id,ao.quantity,sn.sku_id,sn.goods_title,sn.img,ao.state,
 ao.create_time,ao.update_time FROM sale_after_order ao
INNER JOIN sale_sub_order so ON so.id=ao.order_id
INNER JOIN mch_merchant mch ON so.vendor_id = mch.id
INNER JOIN item_trade_snapshot sn ON sn.id = ao.snap_id
WHERE ao.buyer_id= $1 ORDER BY ao.create_time DESC LIMIT $3 OFFSET $2`, func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.PagedMemberAfterSalesOrder{}
				rs.Scan(&e.Id, &e.Type, &e.OrderNo, &e.VendorId, &e.VendorName,
					&e.SnapshotId, &e.Quantity, &e.SkuId, &e.GoodsTitle,
					&e.GoodsImage, &e.State, &e.CreateTime, &e.UpdateTime)
				e.StateText = afterSales.Stat(e.State).String()
				e.GoodsImage = format.GetResUrl(e.GoodsImage)
				list = append(list, e)
			}
		}, memberId, begin, size)
	}
	return total, list
}

// 获取分页售后单
func (a *AfterSalesQuery) QueryPagerAfterSalesOrderOfVendor(vendorId int32, begin,
	size int, where string) (int, []*dto.PagedVendorAfterSalesOrder) {
	var list []*dto.PagedVendorAfterSalesOrder
	total := 0
	if len(where) > 0 {
		where = " AND " + where
	}
	a.ExecScalar(`SELECT COUNT(0) FROM sale_after_order ao
	INNER JOIN sale_sub_order so ON so.id=ao.order_id
	INNER JOIN mm_profile mp ON mp.member_id = so.buyer_id
	INNER JOIN item_trade_snapshot sn ON sn.id = ao.snap_id
	WHERE ao.vendor_id= $1 `+where, &total, vendorId)

	if total > 0 {
		a.Query(`SELECT ao.id,ao.type,so.order_no,so.buyer_id,mp.name as buyer_name,
 ao.snap_id,ao.quantity,sn.sku_id,sn.goods_title,sn.img,ao.state,
 ao.create_time,ao.update_time FROM sale_after_order ao
INNER JOIN sale_sub_order so ON so.id=ao.order_id
INNER JOIN mm_profile mp ON mp.member_id = so.buyer_id
INNER JOIN item_trade_snapshot sn ON sn.id = ao.snap_id
WHERE ao.vendor_id= $1 `+where+" ORDER BY id DESC LIMIT $3 OFFSET $2", func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.PagedVendorAfterSalesOrder{}
				rs.Scan(&e.Id, &e.Type, &e.OrderNo, &e.BuyerId, &e.BuyerName,
					&e.SnapshotId, &e.Quantity, &e.SkuId, &e.GoodsTitle,
					&e.GoodsImage, &e.State, &e.CreateTime, &e.UpdateTime)
				e.StateText = afterSales.Stat(e.State).String()
				e.GoodsImage = format.GetResUrl(e.GoodsImage)
				list = append(list, e)
			}
		}, vendorId, begin, size)
	}
	return total, list
}
