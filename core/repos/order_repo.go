/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package repos

import (
	"database/sql"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"go2o/core"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	orderImpl "go2o/core/domain/order"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/variable"
	"log"
)

var _ order.IOrderRepo = new(OrderRepImpl)

type OrderRepImpl struct {
	Storage storage.Interface
	db.Connector
	_orm         orm.Orm
	_productRepo product.IProductRepo
	_goodsRepo   item.IGoodsItemRepo
	_promRepo    promotion.IPromotionRepo
	_memberRepo  member.IMemberRepo
	_mchRepo     merchant.IMerchantRepo
	_deliverRepo delivery.IDeliveryRepo
	_cartRepo    cart.ICartRepo
	_valRepo     valueobject.IValueRepo
	_cache       map[int]order.IOrderManager
	_payRepo     payment.IPaymentRepo
	_manager     order.IOrderManager
	_expressRepo express.IExpressRepo
	_shipRepo    shipment.IShipmentRepo
}

func NewOrderRepo(sto storage.Interface, c db.Connector,
	mchRepo merchant.IMerchantRepo, payRepo payment.IPaymentRepo,
	proRepo product.IProductRepo, cartRepo cart.ICartRepo, goodsRepo item.IGoodsItemRepo,
	promRepo promotion.IPromotionRepo, memRepo member.IMemberRepo,
	deliverRepo delivery.IDeliveryRepo, expressRepo express.IExpressRepo,
	shipRepo shipment.IShipmentRepo,
	valRepo valueobject.IValueRepo) order.IOrderRepo {
	return &OrderRepImpl{
		Storage:      sto,
		Connector:    c,
		_orm:         c.GetOrm(),
		_productRepo: proRepo,
		_goodsRepo:   goodsRepo,
		_promRepo:    promRepo,
		_payRepo:     payRepo,
		_memberRepo:  memRepo,
		_mchRepo:     mchRepo,
		_cartRepo:    cartRepo,
		_deliverRepo: deliverRepo,
		_valRepo:     valRepo,
		_expressRepo: expressRepo,
		_shipRepo:    shipRepo,
	}
}

func (o *OrderRepImpl) SetPaymentRepo(payRepo payment.IPaymentRepo) {
	o._payRepo = payRepo
}

func (o *OrderRepImpl) Manager() order.IOrderManager {
	if o._productRepo == nil {
		panic("saleRepo uninitialize!")
	}
	if o._manager == nil {
		o._manager = orderImpl.NewOrderManager(o._cartRepo, o._mchRepo,
			o, o._payRepo, o._productRepo, o._goodsRepo, o._promRepo,
			o._memberRepo, o._deliverRepo, o._expressRepo, o._shipRepo,
			o._valRepo)
	}
	return o._manager
}

// 生成订单
func (o *OrderRepImpl) CreateOrder(val *order.Order) order.IOrder {
	return orderImpl.FactoryOrder(val, o.Manager(), o, o._mchRepo, o._goodsRepo,
		o._productRepo, o._promRepo, o._memberRepo, o._expressRepo,
		o._shipRepo, o._payRepo, o._cartRepo, o._valRepo)
}

// 生成空白订单,并保存返回对象
func (o *OrderRepImpl) CreateNormalSubOrder(v *order.NormalSubOrder) order.ISubOrder {
	return orderImpl.NewSubNormalOrder(v, o.Manager(), o, o._memberRepo,
		o._goodsRepo, o._shipRepo, o._productRepo, o._payRepo,
		o._valRepo, o._mchRepo)
}

// 获取可用的订单号
func (o *OrderRepImpl) GetFreeOrderNo(vendorId int32) string {
	//todo:实际应用需要预先生成订单号
	d := o.Connector
	var orderNo string
	for {
		orderNo = domain.NewTradeNo(1, int(vendorId))
		var rec int
		if d.ExecScalar(`SELECT id FROM order_list WHERE order_no= $1 LIMIT 1`,
			&rec, orderNo); rec == 0 {
			break
		}
	}
	return orderNo
}

//　保存订单优惠券绑定
func (o *OrderRepImpl) SaveOrderCouponBind(val *order.OrderCoupon) error {
	_, _, err := o.Connector.GetOrm().Save(nil, val)
	return err
}

// 根据编号获取订单
func (o *OrderRepImpl) GetNormalOrderById(orderId int64) *order.NormalOrder {
	e := &order.NormalOrder{}
	k := o.getOrderCk(orderId, false)
	if o.Storage.Get(k, e) != nil {
		if o.Connector.GetOrm().GetBy(e, "order_id= $1", orderId) != nil {
			return nil
		}
		o.Storage.SetExpire(k, *e, DefaultCacheSeconds*10)
	}
	return e
}

// 保存订单日志
func (o *OrderRepImpl) SaveNormalSubOrderLog(v *order.OrderLog) error {
	_, _, err := o.Connector.GetOrm().Save(nil, v)
	return err
}

// 获取订单的促销绑定
func (o *OrderRepImpl) GetOrderPromotionBinds(orderNo string) []*order.OrderPromotionBind {
	var arr []*order.OrderPromotionBind
	err := o.Connector.GetOrm().Select(&arr, "order_no= $1", orderNo)
	if err == nil {
		return arr
	}
	return make([]*order.OrderPromotionBind, 0)
}

// 保存订单的促销绑定
func (o *OrderRepImpl) SavePromotionBindForOrder(v *order.OrderPromotionBind) (int32, error) {
	return orm.I32(orm.Save(o.GetOrm(), v, int(v.Id)))
}

// 获取订单项
func (o *OrderRepImpl) GetSubOrderItems(orderId int64) []*order.SubOrderItem {
	var items = []*order.SubOrderItem{}
	o.Connector.GetOrm().Select(&items, "order_id= $1", orderId)
	return items
}

func (o *OrderRepImpl) SaveNormalOrder(v *order.NormalOrder) (int, error) {
	id, err := orm.Save(o.GetOrm(), v, int(v.ID))
	if err == nil {
		v.ID = int64(id)
		// 缓存
		o.Storage.SetExpire(o.getOrderCk(v.OrderId, false), *v, DefaultCacheSeconds*10)
		//o.Storage.Set(o.getOrderCkByNo(v.OrderNo, false), v.ID)
	}
	return id, err
}

func (o *OrderRepImpl) GetSubOrderByOrderNo(orderNo string) order.ISubOrder {
	var e = order.NormalSubOrder{}
	err := o.Connector.GetOrm().GetBy(&e, "order_no= $1", orderNo)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:order_sub_order")
		return nil
	}
	return o.CreateNormalSubOrder(&e)
}

// 获取订单的所有子订单
func (o *OrderRepImpl) GetNormalSubOrders(orderId int64) []*order.NormalSubOrder {
	list := make([]*order.NormalSubOrder, 0)
	o.GetOrm().Select(&list, "order_id= $1", orderId)
	return list
}

// 获取缓存订单的Key
func (o *OrderRepImpl) getOrderCk(id int64, subOrder bool) string {
	if subOrder {
		return fmt.Sprintf("go2o:repo:order:s_%d", id)
	}
	return fmt.Sprintf("go2o:repo:order:%d", id)
}

// 获取缓存订单编号的Key
func (o *OrderRepImpl) getOrderCkByNo(orderNo string, subOrder bool) string {
	if subOrder {
		return fmt.Sprintf("go2o:repo:order:s_%s", orderNo)
	}
	return fmt.Sprintf("go2o:repo:order:%s", orderNo)
}

// 获取订单编号
func (o *OrderRepImpl) GetOrderId(orderNo string, subOrder bool) int64 {
	k := o.getOrderCkByNo(orderNo, subOrder)
	id, err := o.Storage.GetInt64(k)
	if err != nil {
		if subOrder {
			o.Connector.ExecScalar("SELECT id FROM sale_sub_order where order_no= $1", &id, orderNo)
		} else {
			o.Connector.ExecScalar("SELECT id FROM order_list where order_no= $1", &id, orderNo)
		}
		if id > 0 {
			o.Storage.Set(k, id)
		}
	}
	return id
}

// 获取子订单
func (o *OrderRepImpl) GetSubOrder(id int64) *order.NormalSubOrder {
	e := &order.NormalSubOrder{}
	k := o.getOrderCk(id, true)
	if o.Storage.Get(k, e) != nil {
		if o.Connector.GetOrm().Get(id, e) != nil {
			return nil
		}
		o.Storage.SetExpire(k, *e, DefaultCacheSeconds*10)
	}
	return e
}

// 保存子订单的商品项,并返回编号和错误
func (o *OrderRepImpl) SaveOrderItem(subOrderId int64, v *order.SubOrderItem) (int32, error) {
	v.OrderId = subOrderId
	return orm.I32(orm.Save(o.GetOrm(), v, int(v.ID)))
}

// 获取订单的操作记录
func (o *OrderRepImpl) GetSubOrderLogs(orderId int64) []*order.OrderLog {
	list := []*order.OrderLog{}
	o.GetOrm().Select(&list, "order_id= $1", orderId)
	return list
}

// 根据商品快照获取订单项
func (o *OrderRepImpl) GetOrderItemBySnapshotId(orderId int64, snapshotId int32) *order.SubOrderItem {
	e := &order.SubOrderItem{}
	if o.GetOrm().GetBy(e, "order_id= $1 AND snap_id= $2", orderId, snapshotId) == nil {
		return e
	}
	return nil
}

// 根据商品快照获取订单项数据传输对象
func (o *OrderRepImpl) GetOrderItemDtoBySnapshotId(orderId int64, snapshotId int32) *dto.OrderItem {
	e := &dto.OrderItem{}
	err := o.QueryRow(`SELECT si.id,si.order_id,si.snap_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.return_quantity,si.amount,si.final_amount,
            si.is_shipped FROM sale_order_item si INNER JOIN item_trade_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id = $1 AND si.snap_id= $2`, func(rs *sql.Row) error {
		err := rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.SkuId, &e.GoodsTitle,
			&e.Image, &e.Price, &e.Quantity, &e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
		e.FinalPrice = e.FinalAmount / float32(e.Quantity)
		return err
	}, orderId, snapshotId)
	if err == nil {
		return e
	}
	return nil
}

// Get OrderList
func (o *OrderRepImpl) GetOrder(where string, arg ...interface{}) *order.Order {
	e := order.Order{}
	err := o._orm.GetBy(&e, where, arg...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:OrderList")
	}
	return nil
}

// 加入到订单通知队列,如果为子订单,则带上sub
func (o *OrderRepImpl) pushOrderQueue(orderNo string, sub bool) {
	rc := core.GetRedisConn()
	if sub {
		content := fmt.Sprintf("sub!%s", orderNo)
		rc.Do("RPUSH", variable.KvOrderBusinessQueue, content)
	} else {
		rc.Do("RPUSH", variable.KvOrderBusinessQueue, orderNo)
	}
	rc.Close()

	//log.Println("----- order notify ! orderNo:", orderNo, " sub:", sub)
}

// Save OrderList
func (o *OrderRepImpl) saveOrder(v *order.Order) (int, error) {
	id, err := orm.Save(o._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:OrderList")
	}
	return id, err
}

// Save OrderList
func (o *OrderRepImpl) SaveOrder(v *order.Order) (int, error) {
	// 零售订单或已拆单的订单不进行通知
	if v.OrderType == int32(order.TRetail) ||
		v.State == int32(order.StatBreak) {
		return o.saveOrder(v)
	}
	// 判断业务状态是否改变
	statusIsChanged := true
	if v.ID <= 0 {
		statusIsChanged = true
	} else {
		origin := o.GetOrder("id= $1", v.ID)
		statusIsChanged = origin.State != v.State
	}
	// log.Println("--- save order:", v.ID, "; state:",
	// v.State, ";", statusIsChanged)
	id, err := o.saveOrder(v)
	if err == nil {
		v.ID = int64(id)
		//如果业务状态已经发生改变,则提交到队列
		if statusIsChanged {
			o.pushOrderQueue(v.OrderNo, false)
		}
	}
	return id, err
}

func (o *OrderRepImpl) saveSubOrder(v *order.NormalSubOrder) (int, error) {
	id, err := orm.Save(o._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:SaleSubOrder")
	}
	if err == nil {
		v.ID = int64(id)
		// 缓存订单号
		o.Storage.Set(o.getOrderCkByNo(v.OrderNo, true), v.ID)
		// 缓存订单
		o.Storage.SetExpire(o.getOrderCk(v.ID, true), *v, DefaultCacheSeconds*10)
	}
	return id, err
}

// 保存子订单
func (o *OrderRepImpl) SaveSubOrder(v *order.NormalSubOrder) (int, error) {
	// 判断业务状态是否改变
	statusIsChanged := true
	if v.ID <= 0 {
		statusIsChanged = true
	} else {
		origin := o.GetSubOrder(v.ID)
		statusIsChanged = origin.State != v.State
	}
	id, err := o.saveSubOrder(v)
	if err == nil {
		v.ID = int64(id)
		//如果业务状态已经发生改变,则提交到队列
		if statusIsChanged {
			o.pushOrderQueue(v.OrderNo, true)
		}
	}
	return id, err
}

// Get WholesaleOrder
func (o *OrderRepImpl) GetWholesaleOrder(where string, v ...interface{}) *order.WholesaleOrder {
	e := order.WholesaleOrder{}
	err := o._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WholesaleOrder")
	}
	return nil
}

// Save WholesaleOrder
func (o *OrderRepImpl) SaveWholesaleOrder(v *order.WholesaleOrder) (int, error) {
	id, err := orm.Save(o._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WholesaleOrder")
	}
	return id, err
}

// Select WholesaleItem
func (o *OrderRepImpl) SelectWholesaleItem(where string, v ...interface{}) []*order.WholesaleItem {
	list := make([]*order.WholesaleItem, 0)
	err := o._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WholesaleItem")
	}
	return list
}

// Save WholesaleItem
func (o *OrderRepImpl) SaveWholesaleItem(v *order.WholesaleItem) (int, error) {
	id, err := orm.Save(o._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:WholesaleItem")
	}
	return id, err
}

// Get OrderTradeOrder
func (o *OrderRepImpl) GetTradeOrder(where string, v ...interface{}) *order.TradeOrder {
	e := order.TradeOrder{}
	err := o._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:OrderTradeOrder")
	}
	return nil
}

// Save OrderTradeOrder
func (o *OrderRepImpl) SaveTradeOrder(v *order.TradeOrder) (int, error) {
	id, err := orm.Save(o._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:OrderTradeOrder")
	}
	return id, err
}
