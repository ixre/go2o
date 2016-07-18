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
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	orderImpl "go2o/core/domain/order"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/variable"
)

var _ order.IOrderRep = new(orderRepImpl)

type orderRepImpl struct {
	db.Connector
	_saleRep    sale.ISaleRep
	_goodsRep   goods.IGoodsRep
	_promRep    promotion.IPromotionRep
	_memberRep  member.IMemberRep
	_mchRep     merchant.IMerchantRep
	_deliverRep delivery.IDeliveryRep
	_cartRep    cart.ICartRep
	_valRep     valueobject.IValueRep
	_cache      map[int]order.IOrderManager
	_payRep     payment.IPaymentRep
	_manager    order.IOrderManager
	_expressRep express.IExpressRep
	_shipRep    shipment.IShipmentRep
}

func NewOrderRep(c db.Connector, ptRep merchant.IMerchantRep, payRep payment.IPaymentRep,
	saleRep sale.ISaleRep, cartRep cart.ICartRep, goodsRep goods.IGoodsRep,
	promRep promotion.IPromotionRep, memRep member.IMemberRep,
	deliverRep delivery.IDeliveryRep, expressRep express.IExpressRep,
	shipRep shipment.IShipmentRep,
	valRep valueobject.IValueRep) *orderRepImpl {
	return &orderRepImpl{
		Connector:   c,
		_saleRep:    saleRep,
		_goodsRep:   goodsRep,
		_promRep:    promRep,
		_payRep:     payRep,
		_memberRep:  memRep,
		_mchRep:     ptRep,
		_cartRep:    cartRep,
		_deliverRep: deliverRep,
		_valRep:     valRep,
		_expressRep: expressRep,
		_shipRep:    shipRep,
	}
}

func (this *orderRepImpl) SetPaymentRep(payRep payment.IPaymentRep) {
	this._payRep = payRep
}

func (this *orderRepImpl) Manager() order.IOrderManager {
	if this._saleRep == nil {
		panic("saleRep uninitialize!")
	}
	if this._manager == nil {
		this._manager = orderImpl.NewOrderManager(this._cartRep, this._mchRep,
			this, this._payRep, this._saleRep, this._goodsRep, this._promRep,
			this._memberRep, this._deliverRep, this._expressRep, this._shipRep,
			this._valRep)
	}
	return this._manager
}

// 获取可用的订单号
func (this *orderRepImpl) GetFreeOrderNo(vendorId int) string {
	//todo:实际应用需要预先生成订单号
	d := this.Connector
	var order_no string
	for {
		order_no = domain.NewOrderNo(vendorId, "")
		var rec int
		if d.ExecScalar(`SELECT COUNT(0) FROM pt_order where order_no=?`,
			&rec, order_no); rec == 0 {
			break
		}
	}
	return order_no
}

//　保存订单优惠券绑定
func (this *orderRepImpl) SaveOrderCouponBind(val *order.OrderCoupon) error {
	_, _, err := this.Connector.GetOrm().Save(nil, val)
	return err
}

// 根据编号获取订单
func (this *orderRepImpl) GetOrderById(id int) *order.Order {
	var v = new(order.Order)
	if err := this.Connector.GetOrm().Get(id, v); err == nil {
		return v
	}
	return nil
}

// 根据订单号获取订单
func (this *orderRepImpl) GetValueOrderByNo(orderNo string) *order.Order {
	e := new(order.Order)
	if this.Connector.GetOrm().GetBy(e, "order_no=?", orderNo) == nil {
		return e
	}
	return nil
}

// 获取等待处理的订单
func (this *orderRepImpl) GetWaitingSetupOrders(vendorId int) ([]*order.Order, error) {
	dst := []*order.Order{}
	err := this.Connector.GetOrm().Select(&dst,
		"(vendor_id <= 0 OR vendor_id=?) AND is_suspend=0 AND status IN("+
			enum.ORDER_SETUP_STATE+")", vendorId)
	if err != nil {
		return nil, err
	}
	return dst, err
}

// 保存订单日志
func (this *orderRepImpl) SaveSubOrderLog(v *order.OrderLog) error {
	_, _, err := this.Connector.GetOrm().Save(nil, v)
	return err
}

// 获取订单的促销绑定
func (this *orderRepImpl) GetOrderPromotionBinds(orderNo string) []*order.OrderPromotionBind {
	var arr []*order.OrderPromotionBind = []*order.OrderPromotionBind{}
	err := this.Connector.GetOrm().Select(&arr, "order_no=?", orderNo)
	if err == nil {
		return arr
	}
	return make([]*order.OrderPromotionBind, 0)
}

// 保存订单的促销绑定
func (this *orderRepImpl) SavePromotionBindForOrder(v *order.OrderPromotionBind) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		var id64 int64
		_, id64, err = orm.Save(nil, v)
		v.Id = int(id64)
	}
	return v.Id, err
}

// 获取订单项
func (this *orderRepImpl) GetOrderItems(orderId int) []*order.OrderItem {
	var items = []*order.OrderItem{}
	this.Connector.GetOrm().Select(&items, "order_id=?", orderId)
	return items
}

// 根据父订单编号获取购买的商品项
func (this *orderRepImpl) GetItemsByParentOrderId(orderId int) []*order.OrderItem {
	var items = []*order.OrderItem{}
	this.Connector.GetOrm().SelectByQuery(&items,
		"order_id IN (SELECT id FROM sale_sub_order WHERE parent_order=?)", orderId)
	return items
}

func (this *orderRepImpl) SaveOrder(v *order.Order) (int, error) {
	return orm.Save(this.GetOrm(), v, v.Id)
}

// 获取订单的所有子订单
func (this *orderRepImpl) GetSubOrdersByParentId(orderId int) []*order.SubOrder {
	list := []*order.SubOrder{}
	this.GetOrm().Select(&list, "parent_order=?", orderId)
	return list
}

// 获取订单编号
func (this *orderRepImpl) GetOrderId(orderNo string, subOrder bool) int {
	id := 0
	if subOrder {
		this.Connector.ExecScalar("SELECT id FROM sale_sub_order where order_no=?", &id, orderNo)
	} else {
		this.Connector.ExecScalar("SELECT id FROM sale_order where order_no=?", &id, orderNo)
	}
	return id
}

// 获取子订单
func (this *orderRepImpl) GetSubOrder(id int) *order.SubOrder {
	e := &order.SubOrder{}
	if this.GetOrm().Get(id, e) == nil {
		return e
	}
	return nil
}

// 根据订单号获取子订单
func (this *orderRepImpl) GetSubOrderByNo(orderNo string) *order.SubOrder {
	e := &order.SubOrder{}
	if this.GetOrm().GetBy(e, "order_no=?", orderNo) == nil {
		return e
	}
	return nil
}

// 保存子订单
func (o *orderRepImpl) SaveSubOrder(v *order.SubOrder) (int, error) {
	var err error
	var statusIsChanged bool //业务状态是否改变
	d := o.Connector

	if v.Id <= 0 {
		statusIsChanged = true
	} else {
		var oriStatus int
		d.ExecScalar("SELECT state FROM sale_sub_order WHERE id=?", &oriStatus, v.Id)
		statusIsChanged = oriStatus != v.State // 业务状态是否改变
	}
	v.Id, err = orm.Save(o.GetOrm(), v, v.Id)

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
func (this *orderRepImpl) SaveOrderItem(subOrderId int, v *order.OrderItem) (int, error) {
	v.OrderId = subOrderId
	return orm.Save(this.GetOrm(), v, v.Id)
}

// 获取订单的操作记录
func (this *orderRepImpl) GetSubOrderLogs(orderId int) []*order.OrderLog {
	list := []*order.OrderLog{}
	this.GetOrm().Select(&list, "order_id=?", orderId)
	return list
}

// 根据商品快照获取订单项
func (o *orderRepImpl) GetOrderItemBySnapshotId(orderId int, snapshotId int) *order.OrderItem {
	e := &order.OrderItem{}
	if o.GetOrm().GetBy(e, "order_id=? AND snap_id=?", orderId, snapshotId) == nil {
		return e
	}
	return nil
}

// 根据商品快照获取订单项数据传输对象
func (o *orderRepImpl) GetOrderItemDtoBySnapshotId(orderId int, snapshotId int) *dto.OrderItem {
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
