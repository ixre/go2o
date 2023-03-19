/**
 * Copyright 2015 @ 56x.net.
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
	"log"
	"strconv"
	"strings"

	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

type OrderQuery struct {
	db.Connector
}

func NewOrderQuery(o orm.Orm) *OrderQuery {
	return &OrderQuery{Connector: o.Connector()}
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
				rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.ItemTitle,
					&e.Image, &e.Price, &e.Quantity, &e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
				e.FinalPrice = int64(float64(e.FinalAmount) / float64(e.Quantity))
				e.Image = format.GetGoodsImageUrl(e.Image)
				list = append(list, e)
			}
		})
	}
	return list
}

// QueryOrderItems 获取订单的商品项
func (o *OrderQuery) QueryOrderItems(subOrderId int64) []*dto.OrderItem {
	return o.queryOrderItems(strconv.Itoa(int(subOrderId)))
}

// QueryPagingNormalOrder 查询分页普通订单
func (o *OrderQuery) QueryPagingNormalOrder(memberId, begin, size int64, pagination bool, where, orderBy string) (int, []*dto.MemberPagingOrderDto) {
	d := o.Connector
	orderList := make([]*dto.MemberPagingOrderDto, 0)
	num := 0
	if size == 0 || begin < 0 {
		return 0, orderList
	}
	if len(where) > 0 {
		where += " AND "
	}
	where += " break_status <> 0"
	if memberId > 0 {
		where += fmt.Sprintf(" AND buyer_id = %d", memberId)
	}
	if orderBy != "" {
		orderBy = " ORDER BY " + orderBy
	} else {
		orderBy = " ORDER BY create_time desc "
	}

	if pagination {
		err := d.ExecScalar(fmt.Sprintf(`
			SELECT COUNT(1) FROM sale_sub_order WHERE is_forbidden = 0 AND %s`,
			where), &num)
		if err != nil {
			log.Println("query order error", err.Error())
			return 0, nil
		}
		if num == 0 {
			return num, orderList
		}
	}

	//orderMap := make(map[int64]int) //存储订单编号和对象的索引
	// 查询分页的订单
	cmd := fmt.Sprintf(`SELECT id,order_no,buyer_id,shop_id,shop_name,express_fee,
	item_count,final_amount,status,create_time
	FROM sale_sub_order  
	 WHERE is_forbidden = 0 AND %s %s LIMIT $2 OFFSET $1`,
		where, orderBy)
	err := d.Query(cmd,
		func(rs *sql.Rows) {
			i := 0
			for rs.Next() {
				e := &dto.MemberPagingOrderDto{}
				err := rs.Scan(&e.OrderId, &e.OrderNo, &e.BuyerId,
					&e.ShopId, &e.ShopName, &e.ExpressFee, &e.ItemCount,
					&e.FinalAmount, &e.Status, &e.CreateTime)
				if err != nil {
					log.Println(" normal order list scan error:", err.Error())
				}
				e.StatusText = order.OrderStatus(e.Status).String()
				orderList = append(orderList, e)
				//orderMap[e.Id] = i
				//idBuf.WriteString(strconv.Itoa(int(e.Id)))
				i++
			}
			_ = rs.Close()
		}, begin, size)
	//log.Println(cmd)
	if err != nil {
		log.Println("[ GO2O][ ERROR]:query order error", err.Error())
	}
	// 获取子订单
	if l := len(orderList); l > 0 {
		idList := make([]string, l)
		orderMap := make(map[int64]*dto.MemberPagingOrderDto, 0)
		for i, ord := range orderList {
			orderMap[ord.OrderId] = ord
			idList[i] = strconv.Itoa(int(ord.OrderId))
		}

		items := o.queryNormalOrderItems(idList)
		for _, v := range items {
			if _, ok := orderMap[v.OrderId]; ok {
				orderMap[v.OrderId].Items = append(orderMap[v.OrderId].Items, v)
			}
		}
	}
	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagedNormalOrderOfVendor(vendorId int64, begin, size int, pagination bool,
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
		orderBy = " ORDER BY o.create_time desc "
	}

	if pagination {
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM sale_sub_order o
		  WHERE o.vendor_id= $1 
		  AND break_status > 0 %s`,
			where), &num, vendorId)
		if num == 0 {
			return num, orderList
		}
	}
	// 查询分页的订单
	err := d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,
		o.buyer_id,o.item_amount,o.discount_amount,o.express_fee,
        o.package_fee,o.final_amount,o.status,o.create_time
         FROM sale_sub_order o 
         WHERE o.vendor_id= $1 AND break_status > 0 %s %s LIMIT $3 OFFSET $2`,
		where, orderBy),
		func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.PagedVendorOrder{Items: []*dto.OrderItem{}}
				rs.Scan(&e.Id, &e.OrderNo, &e.BuyerId,
					&e.ItemAmount, &e.DiscountAmount, &e.ExpressFee,
					&e.PackageFee, &e.FinalAmount, &e.Status, &e.CreateTime)
				e.StatusText = order.OrderStatus(e.Status).String()
				orderList = append(orderList, e)
			}
			rs.Close()
		}, vendorId, begin, size)
	if err != nil {
		log.Println("order query", err.Error())
	}
	orderMap := make(map[int64]*dto.PagedVendorOrder) //存储订单编号和对象的索引
	buf := bytes.NewBufferString("")
	for i, v := range orderList {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(int(v.Id)))
		orderMap[v.Id] = v
	}
	// 查询分页订单的Item
	d.Query(fmt.Sprintf(`SELECT si.id,si.seller_order_id,si.snap_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.amount,si.final_amount
            FROM sale_order_item si INNER JOIN item_trade_snapshot sn
            ON sn.id=si.snap_id WHERE si.seller_order_id IN (%s)
            ORDER BY si.id ASC`, buf.String()), func(rs *sql.Rows) {
		for rs.Next() {
			e := &dto.OrderItem{}
			rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.ItemTitle,
				&e.Image, &e.Price, &e.Quantity, &e.Amount, &e.FinalAmount)
			e.Image = format.GetFileFullUrl(e.Image)
			e.FinalPrice = int64(float64(e.FinalAmount) / float64(e.Quantity))
			orderMap[e.OrderId].Items = append(orderMap[e.OrderId].Items, e)
		}
	})
	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagedWholesaleOrderOfBuyer(memberId, begin, size int64, pagination bool, where, orderBy string) (int, []*dto.MemberPagingOrderDto) {
	d := o.Connector
	var orderList []*dto.MemberPagingOrderDto
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
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM order_list o
		 INNER JOIN order_wholesale_order wo ON wo.order_id = o.id
		 WHERE o.buyer_id= $1 %s`,
			where), &num, memberId)
		if num == 0 {
			return num, orderList
		}
	}

	//orderMap := make(map[int64]int) //存储订单编号和对象的索引
	//idBuf := bytes.NewBufferString("")

	// 查询分页的订单
	d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,
		wo.vendor_id,wo.shop_id,mch.company_name as seller_name,
        wo.item_amount,wo.discount_amount,wo.express_fee,
        wo.package_fee,wo.final_amount,wo.is_paid,wo.status,wo.create_time
         FROM order_list o INNER JOIN order_wholesale_order wo ON wo.order_id = o.id
		INNER JOIN mch_merchant mch ON mch.id= wo.vendor_id
         WHERE o.buyer_id= $1 %s %s LIMIT $3 OFFSET $2`,
		where, orderBy),
		func(rs *sql.Rows) {
			//i := 0
			for rs.Next() {
				//e := &dto.MemberPagingOrderDto{
				//	//Items: []*dto.OrderItem{},
				//}
				//rs.Scan(&e.Id, &e.OrderNo, &e.VendorId, &e.ShopId,
				//	&e.ShopName, &e.ItemAmount, &e.DiscountAmount, &e.ExpressFee,
				//	&e.PackageFee, &e.FinalAmount, &e.IsPaid, &e.State,
				//	&e.CreateTime)
				//e.StateText = order.OrderState(e.State).String()
				//orderList = append(orderList, e)
				//orderMap[e.Id] = i
				//if i != 0 {
				//	idBuf.WriteString(",")
				//}
				//idBuf.WriteString(strconv.Itoa(int(e.Id)))
				//i++
			}
			rs.Close()
		}, memberId, begin, size)

	// 查询分页订单的Item
	//idArr := idBuf.String()
	//if idArr != "" {
	//	d.Query(fmt.Sprintf(` SELECT oi.id,oi.order_id,oi.snapshot_id,sn.item_id,sn.sku_id,
	//        sn.goods_title,sn.img,oi.quantity,oi.return_quantity,
	//        oi.amount,oi.final_amount,oi.is_shipped FROM order_wholesale_item oi
	//        INNER JOIN item_trade_snapshot sn ON sn.id=oi.snapshot_id
	//        WHERE oi.order_id IN(%s) ORDER BY oi.id ASC`, idArr), func(rs *sql.Rows) {
	//		for rs.Next() {
	//			e := &dto.OrderItem{}
	//			rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.ItemTitle,
	//				&e.Image, &e.Quantity, &e.ReturnQuantity,
	//				&e.Amount, &e.FinalAmount, &e.IsShipped)
	//			e.FinalPrice = int64(float64(e.FinalAmount) / float64(e.Quantity))
	//			e.Image = format.GetResUrl(e.Image)
	//			orderList[orderMap[e.OrderId]].Items = append(
	//				orderList[orderMap[e.OrderId]].Items, e)
	//		}
	//	})
	//}
	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagedWholesaleOrderOfVendor(vendorId int64, begin, size int, pagination bool,
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
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM order_list o
		 INNER JOIN order_wholesale_order wo ON wo.order_id = o.id
		 WHERE wo.vendor_id= $1 %s`, where), &num, vendorId)
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
			wo.package_fee,wo.final_amount,wo.is_paid,wo.status,wo.create_time
	    FROM order_list o INNER JOIN order_wholesale_order wo ON wo.order_id = o.id
        INNER JOIN mm_profile mp ON mp.member_id=o.buyer_id
	    WHERE wo.vendor_id= $1 %s %s LIMIT $3 OFFSET $2`,
		where, orderBy),
		func(rs *sql.Rows) {
			i := 0
			for rs.Next() {
				e := &dto.PagedVendorOrder{
					Items: []*dto.OrderItem{},
				}
				rs.Scan(&e.Id, &e.OrderNo, &e.BuyerId, &e.BuyerName,
					&e.ItemAmount, &e.DiscountAmount, &e.ExpressFee,
					&e.PackageFee, &e.FinalAmount, &e.IsPaid, &e.Status, &e.CreateTime)
				e.StatusText = order.OrderStatus(e.Status).String()
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
				&e.ItemTitle, &e.Image, &e.Quantity,
				&e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
			e.FinalPrice = int64(float64(e.FinalAmount) / float64(e.Quantity))
			e.Image = format.GetFileFullUrl(e.Image)
			orderList[orderMap[e.OrderId]].Items = append(
				orderList[orderMap[e.OrderId]].Items, e)
		}
	})

	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagingTradeOrderOfBuyer(memberId, begin, size int64, pagination bool, where, orderBy string) (int, []*proto.SSingleOrder) {
	d := o.Connector
	var orderList []*proto.SSingleOrder
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
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM order_list o
		  INNER JOIN order_trade_order ot ON ot.order_id = o.id WHERE o.buyer_id= $1 %s`,
			where), &num, memberId)
		if num == 0 {
			return num, orderList
		}
	}

	// 查询分页的订单
	err := d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,vendor_id,
        ot.order_amount,ot.discount_amount,
        ot.final_amount,ot.cash_pay,ot.ticket_image, o.status,o.create_time
        FROM order_list o INNER JOIN order_trade_order ot ON ot.order_id = o.id
         WHERE o.buyer_id= $1 %s %s LIMIT $3 OFFSET $2`,
		where, orderBy),
		func(rs *sql.Rows) {
			var cashPay int
			var ticket string
			for rs.Next() {
				e := &proto.SSingleOrder{}
				rs.Scan(&e.OrderId, &e.OrderNo, &e.SellerId, 
					&e.ItemAmount, &e.DiscountAmount, &e.FinalAmount,
					&cashPay, &ticket, &e.Status, &e.SubmitTime)
				e.Data = map[string]string{
					"StatusText":  order.OrderStatus(e.Status).String(),
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
func (o *OrderQuery) PagedTradeOrderOfVendor(vendorId int64, begin, size int, pagination bool,
	where, orderBy string) (int32, []*dto.PagedVendorOrder) {
	d := o.Connector
	var orderList []*dto.PagedVendorOrder
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
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM order_list o
		  INNER JOIN order_trade_order ot ON ot.order_id = o.id WHERE ot.vendor_id= $1 %s`,
			where), &num, vendorId)
		if num == 0 {
			return num, orderList
		}
	}

	// 查询分页的订单
	err := d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,buyer_id,ot.subject,
        ot.order_amount,ot.discount_amount,
        ot.final_amount,ot.cash_pay,ot.ticket_image, o.status,o.create_time,
        m.username FROM order_list o INNER JOIN order_trade_order ot ON ot.order_id = o.id
        LEFT JOIN mm_member m ON m.id = o.buyer_id
         WHERE ot.vendor_id= $1 %s %s LIMIT $3 OFFSET $2`,
		where, orderBy),
		func(rs *sql.Rows) {
			var cashPay int
			var ticket string
			var user string
			for rs.Next() {
				e := &dto.PagedVendorOrder{}
				_ = rs.Scan(&e.Id, &e.OrderNo, &e.BuyerId, &e.Details,
					&e.ItemAmount, &e.DiscountAmount, &e.FinalAmount,
					&cashPay, &ticket, &e.Status, &e.CreateTime, &user)
				e.Data = map[string]string{
					"statusText":  order.OrderStatus(e.Status).String(),
					"cashPay":     strconv.Itoa(cashPay),
					"ticketImage": ticket,
					"user":        user,
					"createTime":  format.UnixTimeStr(e.CreateTime),
				}
				orderList = append(orderList, e)
			}
			_ = rs.Close()
		}, vendorId, begin, size)

	if err != nil && err != sql.ErrNoRows {
		log.Println("QueryPagerTradeOrder: ", err)
	}
	return num, orderList
}

func (o *OrderQuery) queryNormalOrderItems(idArr []string) []*dto.OrderItem {
	list := make([]*dto.OrderItem, 0)
	// 查询分页订单的Item
	_ = o.Query(fmt.Sprintf(`SELECT si.id,si.order_id,si.snap_id,sn.item_id,sn.sku,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.return_quantity,si.amount,si.final_amount,
            si.is_shipped FROM sale_order_item si INNER JOIN item_trade_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id IN (%s)
            ORDER BY si.id ASC`,
		strings.Join(idArr, ",")),
		func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.OrderItem{}
				_ = rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId,
					&e.SpecWord, &e.SkuId, &e.ItemTitle,
					&e.Image, &e.Price, &e.Quantity, &e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
				e.FinalPrice = int64(float64(e.FinalAmount) / float64(e.Quantity))
				e.Image = format.GetGoodsImageUrl(e.Image)
				list = append(list, e)
			}
		})
	return list
}
