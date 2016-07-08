/**
 * Copyright 2015 @ z3q.net.
 * name : order_query
 * author : jarryliu
 * date : 2016-07-08 15:32
 * description :
 * history :
 */
package query

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/jsix/gof/db"
	"go2o/core/dto"
	"strconv"
)

type OrderQuery struct {
	db.Connector
}

func NewOrderQuery(conn db.Connector) *OrderQuery {
	return &OrderQuery{Connector: conn}
}

// 查询分页订单
func (this *OrderQuery) QueryPagerOrder(memberId, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedMemberSubOrder) {
	d := this.Connector
	orderList := []*dto.PagedMemberSubOrder{}
	num := 0
	if size == 0 || begin < 0 {
		return 0, orderList
	}
	if where != "" {
		where = "AND " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	} else {
		orderBy = " ORDER BY po.create_time desc "
	}

	if pagination {
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM sale_sub_order
		  INNER JOIN sale_order ON sale_sub_order.parent_order = sale_order.id
		   WHERE buyer_id=? %s`,
			where), &num, memberId)
		if num == 0 {
			return num, orderList
		}
	}

	orderMap := make(map[int]int) //存储订单编号和对象的索引
	idBuf := bytes.NewBufferString("")

	// 查询分页的订单
	d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,po.order_no as parent_no,
        vendor_id,o.shop_id,s.name as shop_name,
        o.goods_fee,o.discount_fee,o.express_fee,
        o.package_fee,o.final_fee,o.status
         FROM flm.sale_sub_order o INNER JOIN sale_order po ON po.id=o.parent_order
            INNER JOIN mch_shop s ON o.shop_id=s.id
         WHERE buyer_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(rs *sql.Rows) {
			i := 0
			for rs.Next() {
				e := &dto.PagedMemberSubOrder{
					Items: []*dto.OrderItem{},
				}
				rs.Scan(&e.Id, &e.OrderNo, &e.ParentNo, &e.VendorId, &e.ShopId,
					&e.ShopName, &e.GoodsFee, &e.DiscountFee, &e.ExpressFee,
					&e.PackageFee, &e.FinalFee, &e.Status)
				orderList = append(orderList, e)
				orderMap[e.Id] = i
				if i != 0 {
					idBuf.WriteString(",")
				}
				idBuf.WriteString(strconv.Itoa(e.Id))
				i++
			}
			rs.Close()
		}, memberId, begin, size)

	// 查询分页订单的Item
	d.Query(fmt.Sprintf(`SELECT si.id,si.order_id,si.snap_id,sn.sku_id,
            sn.goods_title,sn.img,si.quantity,si.fee,si.final_fee
            FROM sale_order_item si INNER JOIN gs_sales_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id IN(%s)
            ORDER BY si.id ASC`, idBuf.String()), func(rs *sql.Rows) {
		for rs.Next() {
			e := &dto.OrderItem{}
			rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.SkuId, &e.GoodsTitle,
				&e.Image, &e.Quantity, &e.Fee, &e.FinalFee)
			orderList[orderMap[e.OrderId]].Items = append(
				orderList[orderMap[e.OrderId]].Items, e)
		}
	})

	return num, orderList
}
