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
	"go2o/core/domain/interface/order"
	"go2o/core/dto"
	"go2o/core/infrastructure/format"
	"strconv"
)

type OrderQuery struct {
	db.Connector
}

func NewOrderQuery(conn db.Connector) *OrderQuery {
	return &OrderQuery{Connector: conn}
}

func (o *OrderQuery) queryOrderItems(idArr string) []*dto.OrderItem {
	list := []*dto.OrderItem{}
	if idArr != "" {
		// 查询分页订单的Item
		o.Query(fmt.Sprintf(`SELECT si.id,si.order_id,si.snap_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.return_quantity,si.amount,si.final_amount,
            si.is_shipped FROM sale_order_item si INNER JOIN gs_sales_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id IN(%s)
            ORDER BY si.id ASC`, idArr), func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.OrderItem{}
				rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.GoodsTitle,
					&e.Image, &e.Price, &e.Quantity, &e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
				e.FinalPrice = e.FinalAmount / float32(e.Quantity)
				list = append(list, e)
			}
		})
	}
	return list
}

// 获取订单的商品项
func (o *OrderQuery) QueryOrderItems(subOrderId int32) []*dto.OrderItem {
	return o.queryOrderItems(strconv.Itoa(int(subOrderId)))
}

// 查询分页订单
func (o *OrderQuery) QueryPagerOrder(memberId int32, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedMemberSubOrder) {
	d := o.Connector
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
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM sale_sub_order o
		  INNER JOIN sale_order po ON o.parent_order = po.id WHERE o.buyer_id=? %s`,
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
        o.goods_amount,o.discount_amount,o.express_fee,
        o.package_fee,o.final_amount,o.is_paid,o.state,po.create_time
         FROM sale_sub_order o INNER JOIN sale_order po ON po.id=o.parent_order
            INNER JOIN mch_shop s ON o.shop_id=s.id
         WHERE o.buyer_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(rs *sql.Rows) {
			i := 0
			for rs.Next() {
				e := &dto.PagedMemberSubOrder{
					Items: []*dto.OrderItem{},
				}
				rs.Scan(&e.Id, &e.OrderNo, &e.ParentNo, &e.VendorId, &e.ShopId,
					&e.ShopName, &e.GoodsAmount, &e.DiscountAmount, &e.ExpressFee,
					&e.PackageFee, &e.FinalAmount, &e.IsPaid, &e.State, &e.CreateTime)
				e.StateText = order.OrderState(e.State).String()
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
	idArr := idBuf.String()
	if idArr != "" {
		d.Query(fmt.Sprintf(`SELECT si.id,si.order_id,si.snap_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.return_quantity,
            si.amount,si.final_amount,
            si.is_shipped FROM sale_order_item si INNER JOIN gs_sales_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id IN(%s)
            ORDER BY si.id ASC`, idArr), func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.OrderItem{}
				rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.GoodsTitle,
					&e.Image, &e.Price, &e.Quantity, &e.ReturnQuantity,
					&e.Amount, &e.FinalAmount, &e.IsShipped)
				e.FinalPrice = e.FinalAmount / float32(e.Quantity)
				e.Image = format.GetResUrl(e.Image)
				orderList[orderMap[e.OrderId]].Items = append(
					orderList[orderMap[e.OrderId]].Items, e)
			}
		})
	}
	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagedOrdersOfVendor(vendorId int32, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	d := o.Connector
	orderList := []*dto.PagedVendorOrder{}
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
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM sale_sub_order o
		  INNER JOIN sale_order po ON o.parent_order = po.id WHERE o.vendor_id=? %s`,
			where), &num, vendorId)
		if num == 0 {
			return num, orderList
		}
	}

	orderMap := make(map[int]int) //存储订单编号和对象的索引
	idBuf := bytes.NewBufferString("")

	// 查询分页的订单
	d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,po.order_no as parent_no,
		o.buyer_id,mp.name as buyer_name,o.goods_amount,o.discount_amount,o.express_fee,
        o.package_fee,o.final_amount,o.is_paid,o.state,po.create_time
         FROM sale_sub_order o INNER JOIN sale_order po ON po.id=o.parent_order
         INNER JOIN mm_profile mp ON mp.member_id=o.buyer_id
         WHERE o.vendor_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(rs *sql.Rows) {
			i := 0
			for rs.Next() {
				e := &dto.PagedVendorOrder{
					Items: []*dto.OrderItem{},
				}
				rs.Scan(&e.Id, &e.OrderNo, &e.ParentNo, &e.BuyerId, &e.BuyerName,
					&e.GoodsAmount, &e.DiscountAmount, &e.ExpressFee,
					&e.PackageFee, &e.FinalAmount, &e.IsPaid, &e.State, &e.CreateTime)
				e.StateText = order.OrderState(e.State).String()
				orderList = append(orderList, e)
				orderMap[e.Id] = i
				if i != 0 {
					idBuf.WriteString(",")
				}
				idBuf.WriteString(strconv.Itoa(e.Id))
				i++
			}
			rs.Close()
		}, vendorId, begin, size)

	// 查询分页订单的Item
	d.Query(fmt.Sprintf(`SELECT si.id,si.order_id,si.snap_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.amount,si.final_amount
            FROM sale_order_item si INNER JOIN gs_sales_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id IN(%s)
            ORDER BY si.id ASC`, idBuf.String()), func(rs *sql.Rows) {
		for rs.Next() {
			e := &dto.OrderItem{}
			rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.GoodsTitle,
				&e.Image, &e.Price, &e.Quantity, &e.Amount, &e.FinalAmount)
			e.FinalPrice = e.FinalAmount / float32(e.Quantity)
			orderList[orderMap[e.OrderId]].Items = append(
				orderList[orderMap[e.OrderId]].Items, e)
		}
	})

	return num, orderList
}
