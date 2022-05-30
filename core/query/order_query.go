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
	"sort"
	"strconv"

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
	if where != "" {
		where = " order_type = 1 AND " + where
	} else {
		where = " order_type = 1"
	}
	if memberId > 0 {
		where += fmt.Sprintf(" AND order_list.buyer_id = %d", memberId)
	}
	if orderBy != "" {
		orderBy = " ORDER BY " + orderBy
	} else {
		orderBy = " ORDER BY order_list.create_time desc "
	}

	if pagination {
		err := d.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM order_list WHERE %s`,
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
	_ = d.Query(fmt.Sprintf(`select order_list.id,order_list.order_no,
			order_list.buyer_id,buyer_user,order_list.item_count,order_list.item_amount,order_list.discount_amount,
			order_list.express_fee,order_list.package_fee,order_list.final_amount,
			order_list.state,order_list.create_time FROM order_list 
         	WHERE %s %s LIMIT $2 OFFSET $1`,
		where, orderBy),
		func(rs *sql.Rows) {
			i := 0
			for rs.Next() {
				e := &dto.MemberPagingOrderDto{}
				//{
				//	Items: []*dto.OrderItem{},
				//}
				err := rs.Scan(&e.OrderId, &e.OrderNo, &e.BuyerId, &e.BuyerUser, &e.ItemCount, &e.ItemAmount,
					&e.DiscountAmount, &e.ExpressFee, &e.PackageFee,
					&e.FinalAmount, &e.State, &e.CreateTime)
				if err != nil {
					log.Println(" normal order list scan error:", err.Error())
				}
				e.SubOrders = make([]*dto.MemberPagingSubOrderDto, 0)
				e.StateText = order.OrderState(e.State).String()
				orderList = append(orderList, e)
				//orderMap[e.Id] = i
				//idBuf.WriteString(strconv.Itoa(int(e.Id)))
				i++
			}
			_ = rs.Close()
		}, begin, size)

	// 获取子订单
	if len(orderList) > 0 {

		orderIdList := make([]int, 0)
		orderMap := make(map[int64]*dto.MemberPagingOrderDto, 0)
		for _, ord := range orderList {
			orderMap[ord.OrderId] = ord
			orderIdList = append(orderIdList, int(ord.OrderId))
		}
		sort.Ints(orderIdList)
		begin := orderIdList[0]
		end := orderIdList[len(orderIdList)-1]
		subOrders := o.queryNormalSubOrd(begin, end)
		if len(subOrders) > 0 {
			orderIdList = make([]int, 0)
			subOrderMap := make(map[int64]*dto.MemberPagingSubOrderDto, 0)
			// 将子订单绑定到父订单
			for _, ord := range subOrders {
				if _, ok := orderMap[ord.ParentOrderId]; ok {
					subOrderMap[ord.OrderId] = ord
					orderIdList = append(orderIdList, int(ord.OrderId))
					orderMap[ord.ParentOrderId].SubOrders = append(orderMap[ord.ParentOrderId].SubOrders, ord)
				}
			}
			// 获取商品
			sort.Ints(orderIdList)
			begin = orderIdList[0]
			end = orderIdList[len(orderIdList)-1]
			items := o.queryNormalOrderItems(begin, end)
			for _, v := range items {
				if _, ok := subOrderMap[v.OrderId]; ok {
					subOrderMap[v.OrderId].Items = append(subOrderMap[v.OrderId].Items, v)
				}
			}
		}
	}
	return num, orderList
}

func (o *OrderQuery) queryNormalSubOrd(begin int, end int) []*dto.MemberPagingSubOrderDto {
	subOrderList := make([]*dto.MemberPagingSubOrderDto, 0)
	_ = o.Connector.Query(`select id,order_id,order_no,shop_id,shop_name,state
			FROM sale_sub_order where order_id between $1 and $2`,
		func(rs *sql.Rows) {
			for rs.Next() {
				e := &dto.MemberPagingSubOrderDto{}
				err := rs.Scan(&e.OrderId, &e.ParentOrderId, &e.OrderNo, &e.ShopId, &e.ShopName, &e.State)
				if err != nil {
					log.Println(" normal sub order list scan error:", err.Error())
				}
				e.StateText = order.OrderState(e.State).String()
				e.Items = make([]*dto.OrderItem, 0)
				subOrderList = append(subOrderList, e)
			}
			_ = rs.Close()
		}, begin, end)

	return subOrderList
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
		orderBy = " ORDER BY po.create_time desc "
	}

	if pagination {
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM sale_sub_order o
		  INNER JOIN order_list po ON po.id=o.order_id WHERE o.vendor_id= $1 %s`,
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
         WHERE o.vendor_id= $1 %s %s LIMIT $3 OFFSET $2`,
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
			rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.ItemTitle,
				&e.Image, &e.Price, &e.Quantity, &e.Amount, &e.FinalAmount)
			e.Image = format.GetResUrl(e.Image)
			e.FinalPrice = int64(float64(e.FinalAmount) / float64(e.Quantity))
			orderList[orderMap[e.OrderId]].Items = append(
				orderList[orderMap[e.OrderId]].Items, e)
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
        wo.package_fee,wo.final_amount,wo.is_paid,wo.state,wo.create_time
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
				//	&e.PackageFee, &e.FinalFee, &e.IsPaid, &e.State,
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
			wo.package_fee,wo.final_amount,wo.is_paid,wo.state,wo.create_time
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
				&e.ItemTitle, &e.Image, &e.Quantity,
				&e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
			e.FinalPrice = int64(float64(e.FinalAmount) / float64(e.Quantity))
			e.Image = format.GetResUrl(e.Image)
			orderList[orderMap[e.OrderId]].Items = append(
				orderList[orderMap[e.OrderId]].Items, e)
		}
	})

	return num, orderList
}

// 查询分页订单
func (o *OrderQuery) PagedTradeOrderOfBuyer(memberId, begin, size int64, pagination bool, where, orderBy string) (int, []*proto.SSingleOrder) {
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
	err := d.Query(fmt.Sprintf(`SELECT o.id,o.order_no,vendor_id,ot.subject,
        ot.order_amount,ot.discount_amount,
        ot.final_amount,ot.cash_pay,ot.ticket_image, o.state,o.create_time
        FROM order_list o INNER JOIN order_trade_order ot ON ot.order_id = o.id
         WHERE o.buyer_id= $1 %s %s LIMIT $3 OFFSET $2`,
		where, orderBy),
		func(rs *sql.Rows) {
			var cashPay int
			var ticket string
			for rs.Next() {
				e := &proto.SSingleOrder{}
				rs.Scan(&e.OrderId, &e.OrderNo, &e.SellerId, &e.Subject,
					&e.ItemAmount, &e.DiscountAmount, &e.FinalAmount,
					&cashPay, &ticket, &e.State, &e.SubmitTime)
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
        ot.final_amount,ot.cash_pay,ot.ticket_image, o.state,o.create_time,
        m.user FROM order_list o INNER JOIN order_trade_order ot ON ot.order_id = o.id
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
					&cashPay, &ticket, &e.State, &e.CreateTime, &user)
				e.Data = map[string]string{
					"StateText":   order.OrderState(e.State).String(),
					"CashPay":     strconv.Itoa(cashPay),
					"TicketImage": ticket,
					"User":        user,
					"CreateTime":  format.UnixTimeStr(e.CreateTime),
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

func (o *OrderQuery) queryNormalOrderItems(begin int, end int) []*dto.OrderItem {
	list := make([]*dto.OrderItem, 0)
	// 查询分页订单的Item
	_ = o.Query(`SELECT si.id,si.order_id,si.snap_id,sn.item_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.return_quantity,si.amount,si.final_amount,
            si.is_shipped FROM sale_order_item si INNER JOIN item_trade_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id BETWEEN $1 AND $2
            ORDER BY si.id ASC`, func(rs *sql.Rows) {
		for rs.Next() {
			e := &dto.OrderItem{}
			_ = rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.ItemId, &e.SkuId, &e.ItemTitle,
				&e.Image, &e.Price, &e.Quantity, &e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
			e.FinalPrice = int64(float64(e.FinalAmount) / float64(e.Quantity))
			e.Image = format.GetGoodsImageUrl(e.Image)
			list = append(list, e)
		}
	}, begin, end)
	return list
}
