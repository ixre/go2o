/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package repository

import (
	"database/sql"
	"fmt"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/storage"
	"go2o/core"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/enum"
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
)

var _ order.IOrderRepo = new(orderRepImpl)

type orderRepImpl struct {
	Storage storage.Interface
	db.Connector
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
	valRepo valueobject.IValueRepo) *orderRepImpl {
	return &orderRepImpl{
		Storage:      sto,
		Connector:    c,
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

func (o *orderRepImpl) SetPaymentRepo(payRepo payment.IPaymentRepo) {
	o._payRepo = payRepo
}

func (o *orderRepImpl) Manager() order.IOrderManager {
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

// 获取可用的订单号
func (o *orderRepImpl) GetFreeOrderNo(vendorId int32) string {
	//todo:实际应用需要预先生成订单号
	d := o.Connector
	var order_no string
	for {
		order_no = domain.NewOrderNo(int(vendorId), "")
		var rec int
		if d.ExecScalar(`SELECT COUNT(0) FROM pt_order where order_no=?`,
			&rec, order_no); rec == 0 {
			break
		}
	}
	return order_no
}

//　保存订单优惠券绑定
func (o *orderRepImpl) SaveOrderCouponBind(val *order.OrderCoupon) error {
	_, _, err := o.Connector.GetOrm().Save(nil, val)
	return err
}

// 根据编号获取订单
func (o *orderRepImpl) GetOrderById(id int32) *order.Order {
	e := &order.Order{}
	k := o.getOrderCk(id, false)
	if o.Storage.Get(k, e) != nil {
		if o.Connector.GetOrm().Get(id, e) != nil {
			return nil
		}
		o.Storage.SetExpire(k, *e, DefaultCacheSeconds*10)
	}
	return e
}

// 根据订单号获取订单
func (o *orderRepImpl) GetValueOrderByNo(orderNo string) *order.Order {
	id := o.GetOrderId(orderNo, false)
	if id > 0 {
		return o.GetOrderById(id)
	}
	return nil
}

// 获取等待处理的订单
func (o *orderRepImpl) GetWaitingSetupOrders(vendorId int32) ([]*order.Order, error) {
	dst := []*order.Order{}
	err := o.Connector.GetOrm().Select(&dst,
		"(vendor_id <= 0 OR vendor_id=?) AND is_suspend=0 AND status IN("+
			enum.ORDER_SETUP_STATE+")", vendorId)
	if err != nil {
		return nil, err
	}
	return dst, err
}

// 保存订单日志
func (o *orderRepImpl) SaveSubOrderLog(v *order.OrderLog) error {
	_, _, err := o.Connector.GetOrm().Save(nil, v)
	return err
}

// 获取订单的促销绑定
func (o *orderRepImpl) GetOrderPromotionBinds(orderNo string) []*order.OrderPromotionBind {
	var arr []*order.OrderPromotionBind = []*order.OrderPromotionBind{}
	err := o.Connector.GetOrm().Select(&arr, "order_no=?", orderNo)
	if err == nil {
		return arr
	}
	return make([]*order.OrderPromotionBind, 0)
}

// 保存订单的促销绑定
func (o *orderRepImpl) SavePromotionBindForOrder(v *order.OrderPromotionBind) (int32, error) {
	return orm.I32(orm.Save(o.GetOrm(), v, int(v.Id)))
}

// 获取订单项
func (o *orderRepImpl) GetSubOrderItems(orderId int32) []*order.OrderItem {
	var items = []*order.OrderItem{}
	o.Connector.GetOrm().Select(&items, "order_id=?", orderId)
	return items
}

// 根据父订单编号获取购买的商品项
func (o *orderRepImpl) GetItemsByParentOrderId(orderId int32) []*order.OrderItem {
	var items = []*order.OrderItem{}
	o.Connector.GetOrm().SelectByQuery(&items,
		"order_id IN (SELECT id FROM sale_sub_order WHERE parent_order=?)", orderId)
	return items
}

func (o *orderRepImpl) SaveOrder(v *order.Order) (int32, error) {
	id, err := orm.I32(orm.Save(o.GetOrm(), v, int(v.Id)))
	if err == nil {
		v.Id = id
		// 缓存
		o.Storage.SetExpire(o.getOrderCk(v.Id, false), *v, DefaultCacheSeconds*10)
		o.Storage.Set(o.getOrderCkByNo(v.OrderNo, false), v.Id)
	}
	return id, err
}

// 获取订单的所有子订单
func (o *orderRepImpl) GetSubOrdersByParentId(orderId int32) []*order.SubOrder {
	list := []*order.SubOrder{}
	o.GetOrm().Select(&list, "parent_order=?", orderId)
	return list
}

// 获取缓存订单的Key
func (o *orderRepImpl) getOrderCk(id int32, subOrder bool) string {
	if subOrder {
		return fmt.Sprintf("go2o:rep:order:s_%d", id)
	}
	return fmt.Sprintf("go2o:rep:order:%d", id)
}

// 获取缓存订单编号的Key
func (o *orderRepImpl) getOrderCkByNo(orderNo string, subOrder bool) string {
	if subOrder {
		return fmt.Sprintf("go2o:rep:order:s_%s", orderNo)
	}
	return fmt.Sprintf("go2o:rep:order:%s", orderNo)
}

// 获取订单编号
func (o *orderRepImpl) GetOrderId(orderNo string, subOrder bool) int32 {
	k := o.getOrderCkByNo(orderNo, subOrder)
	id, err := o.Storage.GetInt(k)
	if err != nil {
		if subOrder {
			o.Connector.ExecScalar("SELECT id FROM sale_sub_order where order_no=?", &id, orderNo)
		} else {
			o.Connector.ExecScalar("SELECT id FROM sale_order where order_no=?", &id, orderNo)
		}
		if id > 0 {
			o.Storage.Set(k, id)
		}
	}
	return int32(id)
}

// 获取子订单
func (o *orderRepImpl) GetSubOrder(id int32) *order.SubOrder {
	e := &order.SubOrder{}
	k := o.getOrderCk(id, true)
	if o.Storage.Get(k, e) != nil {
		if o.Connector.GetOrm().Get(id, e) != nil {
			return nil
		}
		o.Storage.SetExpire(k, *e, DefaultCacheSeconds*10)
	}
	return e
}

// 根据订单号获取子订单
func (o *orderRepImpl) GetSubOrderByNo(orderNo string) *order.SubOrder {
	id := o.GetOrderId(orderNo, true)
	if id > 0 {
		return o.GetSubOrder(id)
	}
	return nil
}

// 保存子订单
func (o *orderRepImpl) SaveSubOrder(v *order.SubOrder) (int32, error) {
	var err error
	var statusIsChanged bool //业务状态是否改变
	if v.Id <= 0 {
		statusIsChanged = true
	} else {
		origin := o.GetSubOrder(v.Id)
		statusIsChanged = origin.State != v.State // 业务状态是否改变
	}
	v.Id, err = orm.I32(orm.Save(o.GetOrm(), v, int(v.Id)))

	// 缓存订单号
	o.Storage.Set(o.getOrderCkByNo(v.OrderNo, true), v.Id)
	// 缓存订单
	o.Storage.SetExpire(o.getOrderCk(v.Id, true), *v, DefaultCacheSeconds*10)

	//如果业务状态已经发生改变,则提交到队列
	if statusIsChanged && v.Id > 0 {
		rc := core.GetRedisConn()
		rc.Do("RPUSH", variable.KvOrderBusinessQueue, v.Id)
		rc.Close()
		//log.Println("-----order ",v.Id,v.Status,statusIsChanged,err)
	}
	return v.Id, err
}

// 保存子订单的商品项,并返回编号和错误
func (o *orderRepImpl) SaveOrderItem(subOrderId int32, v *order.OrderItem) (int32, error) {
	v.OrderId = subOrderId
	return orm.I32(orm.Save(o.GetOrm(), v, int(v.Id)))
}

// 获取订单的操作记录
func (o *orderRepImpl) GetSubOrderLogs(orderId int32) []*order.OrderLog {
	list := []*order.OrderLog{}
	o.GetOrm().Select(&list, "order_id=?", orderId)
	return list
}

// 根据商品快照获取订单项
func (o *orderRepImpl) GetOrderItemBySnapshotId(orderId int32, snapshotId int32) *order.OrderItem {
	e := &order.OrderItem{}
	if o.GetOrm().GetBy(e, "order_id=? AND snap_id=?", orderId, snapshotId) == nil {
		return e
	}
	return nil
}

// 根据商品快照获取订单项数据传输对象
func (o *orderRepImpl) GetOrderItemDtoBySnapshotId(orderId int32, snapshotId int32) *dto.OrderItem {
	e := &dto.OrderItem{}
	err := o.QueryRow(`SELECT si.id,si.order_id,si.snap_id,sn.sku_id,
            sn.goods_title,sn.img,sn.price,si.quantity,si.return_quantity,si.amount,si.final_amount,
            si.is_shipped FROM sale_order_item si INNER JOIN gs_sales_snapshot sn
            ON sn.id=si.snap_id WHERE si.order_id = ? AND si.snap_id=?`, func(rs *sql.Row) {
		rs.Scan(&e.Id, &e.OrderId, &e.SnapshotId, &e.SkuId, &e.GoodsTitle,
			&e.Image, &e.Price, &e.Quantity, &e.ReturnQuantity, &e.Amount, &e.FinalAmount, &e.IsShipped)
		e.FinalPrice = e.FinalAmount / float32(e.Quantity)
	}, orderId, snapshotId)
	if err == nil {
		return e
	}
	return nil
}
