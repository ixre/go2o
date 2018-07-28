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
	"go2o/core/service/auto_gen/rpc/order_service"
	"log"
	"strconv"
)

type OrderQuery struct {
	db.Connector
}

func NewOrderQuery(conn db.Connector) *OrderQuery {
	return &OrderQuery{Connector: conn}
}

func (o *OrderQuery) queryOrderItems(idArr string) []*dto.OrderItem {
	list := make([]*dto.OrderItem, 0)
	if idArr != "" && len(idArr) > 0 {
		// 查询分页订单的Item
		o.Query(fmt.Sprintf(`SELECT si.id,si.order_id,si.snap_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.return_quantity,si.amount,si.final_amount,
            si.is_shipped FROM sale_order_item si INNER JOIN item_trade_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id IN(%s)
            ORDER BY si.id ASC`, idArr), func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.OrderItem{}
				rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.GoodsTitle,
					&e.Image, &e.Price, &e.Quantity, &e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
				e.FinalPrice = e.FinalAmount / float32(e.Quantity)
				e.Image = format.GetGoodsImageUrl(e.Image)
				list = append(list, e)
			}
		})
	}
	return list
}

// 获取订单的商品项
func (o *OrderQuery) QueryOrderItems(subOrderId int64) []*dto.OrderItem {
	return o.queryOrderItems(strconv.Itoa(int(subOrderId)))
}

// 查询分页订单
func (o *OrderQuery) QueryPagerOrder(memberId int64, begin, size int, pagination bool,
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
		  INNER JOIN order_list po ON o.order_id = po.id WHERE o.buyer_id=? %s`,
			where), &num, memberId)
		if num == 0 {
			return num, orderList
		}
	}

	orderMap := make(map[int64]int) //存储订单编号和对象的索引
	idBuf := bytes.NewBufferString("")

	// 查询分页的订单
	d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,po.order_no as parent_no,
        o.vendor_id,o.shop_id,s.name as shop_name,
        o.item_amount,o.discount_amount,o.express_fee,
        o.package_fee,o.final_amount,o.is_paid,o.state,po.create_time
         FROM sale_sub_order o INNER JOIN order_list po ON o.order_id = po.id
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
					&e.ShopName, &e.ItemAmount, &e.DiscountAmount, &e.ExpressFee,
					&e.PackageFee, &e.FinalAmount, &e.IsPaid, &e.State, &e.CreateTime)
				e.StateText = order.OrderState(e.State).String()
				orderList = append(orderList, e)
				orderMap[e.Id] = i
				if i != 0 {
					idBuf.WriteString(",")
				}
				idBuf.WriteString(strconv.Itoa(int(e.Id)))
				i++
			}
			rs.Close()
		}, memberId, begin, size)

	// 查询分页订单的Item
	idArr := idBuf.String()
	if len(idArr) > 0 {
		d.Query(fmt.Sprintf(`SELECT si.id,si.order_id,si.snap_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.return_quantity,
            si.amount,si.final_amount,
            si.is_shipped FROM sale_order_item si INNER JOIN item_trade_snapshot sn
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
func (o *OrderQuery) PagedNormalOrderOfVendor(vendorId int32, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	d := o.Connector
	var orderList []*dto.PagedVendorOrder
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
		  INNER JOIN order_list po ON po.id=o.order_id WHERE o.vendor_id=? %s`,
			where), &num, vendorId)
		if num == 0 {
			return num, orderList
		}
	}
	orderMap := make(map[int64]int) //存储订单编号和对象的索引
	idBuf := bytes.NewBufferString("")
	// 查询分页的订单
	d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,po.order_no as parent_no,
		o.buyer_id,mp.name as buyer_name,o.item_amount,o.discount_amount,o.express_fee,
        o.package_fee,o.final_amount,o.is_paid,o.state,po.create_time
         FROM sale_sub_order o INNER JOIN order_list po ON po.id=o.order_id
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
					&e.ItemAmount, &e.DiscountAmount, &e.ExpressFee,
					&e.PackageFee, &e.FinalAmount, &e.IsPaid, &e.State, &e.CreateTime)
				e.StateText = order.OrderState(e.State).String()
				orderList = append(orderList, e)
				orderMap[e.Id] = i
				if i != 0 {
					idBuf.WriteString(",")
				}
				idBuf.WriteString(strconv.Itoa(int(e.Id)))
				i++
			}
			rs.Close()
		}, vendorId, begin, size)

	// 查询分页订单的Item
	d.Query(fmt.Sprintf(`SELECT si.id,si.order_id,si.snap_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.amount,si.final_amount
            FROM sale_order_item si INNER JOIN item_trade_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id IN(%s)
            ORDER BY si.id ASC`, idBuf.String()), func(rs *sql.Rows) {
		for rs.Next() {
			e := &dto.OrderItem{}
			rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.GoodsTitle,
				&e.Image, &e.Price, &e.Quantity, &e.Amount, &e.FinalAmount)
			e.Image = format.GetResUrl(e.Image)
			e.FinalPrice = e.FinalAmount / float32(e.Quantity)
			orderList[orderMap[e.OrderId]].Items = append(
				orderList[orderMap[e.OrderId]].Items, e)
		}
	})
	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagedWholesaleOrderOfBuyer(memberId int64, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedMemberSubOrder) {
	d := o.Connector
	var orderList []*dto.PagedMemberSubOrder
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
		orderBy = " ORDER BY wo.create_time desc "
	}

	if pagination {
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM order_list o
		 INNER JOIN order_wholesale_order wo ON wo.order_id = o.id
		 WHERE o.buyer_id=? %s`,
			where), &num, memberId)
		if num == 0 {
			return num, orderList
		}
	}

	orderMap := make(map[int64]int) //存储订单编号和对象的索引
	idBuf := bytes.NewBufferString("")

	// 查询分页的订单
	d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,
		wo.vendor_id,wo.shop_id,mch.company_name as seller_name,
        wo.item_amount,wo.discount_amount,wo.express_fee,
        wo.package_fee,wo.final_amount,wo.is_paid,wo.state,wo.create_time
         FROM order_list o INNER JOIN order_wholesale_order wo ON wo.order_id = o.id
		INNER JOIN mch_merchant mch ON mch.id= wo.vendor_id
         WHERE o.buyer_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(rs *sql.Rows) {
			i := 0
			for rs.Next() {
				e := &dto.PagedMemberSubOrder{
					Items: []*dto.OrderItem{},
				}
				rs.Scan(&e.Id, &e.OrderNo, &e.VendorId, &e.ShopId,
					&e.ShopName, &e.ItemAmount, &e.DiscountAmount, &e.ExpressFee,
					&e.PackageFee, &e.FinalAmount, &e.IsPaid, &e.State,
					&e.CreateTime)
				e.StateText = order.OrderState(e.State).String()
				orderList = append(orderList, e)
				orderMap[e.Id] = i
				if i != 0 {
					idBuf.WriteString(",")
				}
				idBuf.WriteString(strconv.Itoa(int(e.Id)))
				i++
			}
			rs.Close()
		}, memberId, begin, size)

	// 查询分页订单的Item
	idArr := idBuf.String()
	if idArr != "" {
		d.Query(fmt.Sprintf(` SELECT oi.id,oi.order_id,oi.snapshot_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,oi.quantity,oi.return_quantity,
            oi.amount,oi.final_amount,oi.is_shipped FROM order_wholesale_item oi
            INNER JOIN item_trade_snapshot sn ON sn.id=oi.snapshot_id
            WHERE oi.order_id IN(%s) ORDER BY oi.id ASC`, idArr), func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.OrderItem{}
				rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.GoodsTitle,
					&e.Image, &e.Quantity, &e.ReturnQuantity,
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
func (o *OrderQuery) PagedWholesaleOrderOfVendor(vendorId int32, begin, size int, pagination bool,
	where, orderBy string) (int, []*dto.PagedVendorOrder) {
	d := o.Connector
	var orderList []*dto.PagedVendorOrder
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
		orderBy = " ORDER BY wo.create_time desc "
	}

	if pagination {
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM order_list o
		 INNER JOIN order_wholesale_order wo ON wo.order_id = o.id
		 WHERE wo.vendor_id=? %s`, where), &num, vendorId)
		if num == 0 {
			return num, orderList
		}
	}

	orderMap := make(map[int64]int) //存储订单编号和对象的索引
	idBuf := bytes.NewBufferString("")

	// 查询分页的订单
	d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,
			o.buyer_id,mp.name as buyer_name,
			wo.item_amount,wo.discount_amount,wo.express_fee,
			wo.package_fee,wo.final_amount,wo.is_paid,wo.state,wo.create_time
	    FROM order_list o INNER JOIN order_wholesale_order wo ON wo.order_id = o.id
        INNER JOIN mm_profile mp ON mp.member_id=o.buyer_id
	    WHERE wo.vendor_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(rs *sql.Rows) {
			i := 0
			for rs.Next() {
				e := &dto.PagedVendorOrder{
					Items: []*dto.OrderItem{},
				}
				rs.Scan(&e.Id, &e.OrderNo, &e.BuyerId, &e.BuyerName,
					&e.ItemAmount, &e.DiscountAmount, &e.ExpressFee,
					&e.PackageFee, &e.FinalAmount, &e.IsPaid, &e.State, &e.CreateTime)
				e.StateText = order.OrderState(e.State).String()
				orderList = append(orderList, e)
				orderMap[e.Id] = i
				if i != 0 {
					idBuf.WriteString(",")
				}
				idBuf.WriteString(strconv.Itoa(int(e.Id)))
				i++
			}
			rs.Close()
		}, vendorId, begin, size)

	// 查询分页订单的Item
	d.Query(fmt.Sprintf(` SELECT oi.id,oi.order_id,oi.snapshot_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,oi.quantity,oi.return_quantity,
            oi.amount,oi.final_amount,oi.is_shipped FROM order_wholesale_item oi
            INNER JOIN item_trade_snapshot sn ON sn.id=oi.snapshot_id
            WHERE oi.order_id IN(%s) ORDER BY oi.id ASC`, idBuf.String()), func(rs *sql.Rows) {
		for rs.Next() {
			e := &dto.OrderItem{}
			rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId,
				&e.GoodsTitle, &e.Image, &e.Quantity,
				&e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
			e.FinalPrice = e.FinalAmount / float32(e.Quantity)
			e.Image = format.GetResUrl(e.Image)
			orderList[orderMap[e.OrderId]].Items = append(
				orderList[orderMap[e.OrderId]].Items, e)
		}
	})

	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagedTradeOrderOfBuyer(memberId int64, begin, size int, pagination bool,
	where, orderBy string) (int, []*order_service.SComplexOrder) {
	d := o.Connector
	var orderList []*order_service.SComplexOrder
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
		orderBy = " ORDER BY o.create_time desc "
	}

	if pagination {
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM order_list o
		  INNER JOIN order_trade_order ot ON ot.order_id = o.id WHERE o.buyer_id=? %s`,
			where), &num, memberId)
		if num == 0 {
			return num, orderList
		}
	}

	// 查询分页的订单
	err := d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,vendor_id,ot.subject,
        ot.order_amount,ot.discount_amount,
        ot.final_amount,ot.cash_pay,ot.ticket_image, o.state,o.create_time
        FROM order_list o INNER JOIN order_trade_order ot ON ot.order_id = o.id
         WHERE o.buyer_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(rs *sql.Rows) {
			var cashPay int
			var ticket string
			for rs.Next() {
				e := &order_service.SComplexOrder{}
				rs.Scan(&e.OrderId, &e.OrderNo, &e.VendorId, &e.Subject,
					&e.ItemAmount, &e.DiscountAmount, &e.FinalAmount,
					&cashPay, &ticket, &e.State, &e.CreateTime)
				e.Data = map[string]string{
					"StateText":   order.OrderState(e.State).String(),
					"CashPay":     strconv.Itoa(cashPay),
					"TicketImage": ticket,
				}

				orderList = append(orderList, e)
			}
			rs.Close()
		}, memberId, begin, size)

	if err != nil && err != sql.ErrNoRows {
		log.Println("QueryPagerTradeOrder: ", err)
	}
	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagedTradeOrderOfVendor(vendorId int32, begin, size int, pagination bool,
	where, orderBy string) (int32, []*order_service.SComplexOrder) {
	d := o.Connector
	orderList := []*order_service.SComplexOrder{}
	var num int32
	if size == 0 || begin < 0 {
		return 0, orderList
	}
	if where != "" {
		where = "AND " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	} else {
		orderBy = " ORDER BY o.create_time desc "
	}

	if pagination {
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM order_list o
		  INNER JOIN order_trade_order ot ON ot.order_id = o.id WHERE ot.vendor_id=? %s`,
			where), &num, vendorId)
		if num == 0 {
			return num, orderList
		}
	}

	// 查询分页的订单
	err := d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,vendor_id,ot.subject,
        ot.order_amount,ot.discount_amount,
        ot.final_amount,ot.cash_pay,ot.ticket_image, o.state,o.create_time,
        m.usr FROM order_list o INNER JOIN order_trade_order ot ON ot.order_id = o.id
        LEFT JOIN mm_member m ON m.id = o.buyer_id
         WHERE ot.vendor_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(rs *sql.Rows) {
			var cashPay int
			var ticket string
			var usr string
			for rs.Next() {
				e := &order_service.SComplexOrder{}
				rs.Scan(&e.OrderId, &e.OrderNo, &e.VendorId, &e.Subject,
					&e.ItemAmount, &e.DiscountAmount, &e.FinalAmount,
					&cashPay, &ticket, &e.State, &e.CreateTime, &usr)
				e.Data = map[string]string{
					"StateText":   order.OrderState(e.State).String(),
					"CashPay":     strconv.Itoa(cashPay),
					"TicketImage": ticket,
					"Usr":         usr,
					"CreateTime":  format.UnixTimeStr(e.CreateTime),
				}
				orderList = append(orderList, e)
			}
			rs.Close()
		}, vendorId, begin, size)

	if err != nil && err != sql.ErrNoRows {
		log.Println("QueryPagerTradeOrder: ", err)
	}
	return num, orderList
}
