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
	"errors"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"go2o/core"
	"go2o/core/domain/interface/delivery"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/shopping"
	"go2o/core/domain/interface/valueobject"
	shoppingImpl "go2o/core/domain/shopping"
	"go2o/core/infrastructure/domain"
	"go2o/core/variable"
	"time"
)

var _ shopping.IShoppingRep = new(shoppingRep)

type shoppingRep struct {
	db.Connector
	_saleRep    sale.ISaleRep
	_goodsRep   goods.IGoodsRep
	_promRep    promotion.IPromotionRep
	_memberRep  member.IMemberRep
	_partnerRep merchant.IMerchantRep
	_deliverRep delivery.IDeliveryRep
	_valRep     valueobject.IValueRep
	_cache      map[int]shopping.IShopping
}

func NewShoppingRep(c db.Connector, ptRep merchant.IMerchantRep,
	saleRep sale.ISaleRep, goodsRep goods.IGoodsRep, promRep promotion.IPromotionRep,
	memRep member.IMemberRep, deliverRep delivery.IDeliveryRep,
	valRep valueobject.IValueRep) shopping.IShoppingRep {
	return (&shoppingRep{
		Connector:   c,
		_saleRep:    saleRep,
		_goodsRep:   goodsRep,
		_promRep:    promRep,
		_memberRep:  memRep,
		_partnerRep: ptRep,
		_deliverRep: deliverRep,
		_valRep:     valRep,
	}).init()
}

func (this *shoppingRep) init() shopping.IShoppingRep {
	this._cache = make(map[int]shopping.IShopping)
	return this
}

func (this *shoppingRep) GetShopping(memberId int) shopping.IShopping {
	if this._saleRep == nil {
		panic("saleRep uninitialize!")
	}
	//v, ok := this._cache[merchantId]
	//if !ok {
	v := shoppingImpl.NewShopping(memberId, this._partnerRep,
		this, this._saleRep, this._goodsRep, this._promRep,
		this._memberRep, this._deliverRep, this._valRep)
	//this._cache[merchantId] = v
	//}
	return v
}

// 获取可用的订单号
func (this *shoppingRep) GetFreeOrderNo(merchantId int) string {
	//todo:实际应用需要预先生成订单号
	d := this.Connector
	var order_no string
	for {
		order_no = domain.NewOrderNo(merchantId)
		var rec int
		if d.ExecScalar(`SELECT COUNT(0) FROM pt_order where order_no=?`,
			&rec, order_no); rec == 0 {
			break
		}
	}
	return order_no
}

// 获取订单项
func (this *shoppingRep) GetOrderItems(orderId int) []*shopping.OrderItem {
	var items = []*shopping.OrderItem{}
	this.Connector.GetOrm().Select(&items, "order_id=?", orderId)
	return items
}

func (this *shoppingRep) SaveOrder(merchantId int, v *shopping.ValueOrder) (int, error) {
	var err error
	var statusIsChanged bool //业务状态是否改变
	d := this.Connector
	v.MerchantId = merchantId

	if v.Id > 0 {
		var oriStatus int
		d.ExecScalar("SELECT status FROM pt_order WHERE id=?", &oriStatus, v.Id)
		statusIsChanged = oriStatus != v.Status // 业务状态是否改变
		_, _, err = d.GetOrm().Save(v.Id, v)
		if v.Status == enum.ORDER_COMPLETED {
			//todo:将去掉下行
			gof.CurrentApp.Storage().Set(variable.KvHaveNewCompletedOrder, enum.TRUE)
		}
	} else {
		//验证Merchant和Member是否有绑定关系
		var num int
		if d.ExecScalar(`SELECT COUNT(0) FROM mm_relation WHERE member_id=? AND reg_merchant_id=?`,
			&num, v.MemberId, v.MerchantId); num != 1 {
			return v.Id, errors.New("error partner and member.")
		}
		_, _, err = d.GetOrm().Save(nil, v)
		if err == nil {
			err = d.ExecScalar(`SELECT MAX(id) FROM pt_order WHERE merchant_id=? AND member_id=?`, &v.Id,
				merchantId, v.MemberId)
		}
		statusIsChanged = true
	}

	if statusIsChanged { //如果业务状态已经发生改变,则提交到队列
		rc := core.GetRedisConn()
		defer rc.Close()
		rc.Do("RPUSH", variable.KvOrderBusinessQueue, v.Id) // push to queue
		//log.Println("-- PUSH - ",v.Id,err)
	}

	// 保存订单项
	unix := time.Now().Unix()
	if err == nil && v.Items != nil {
		orm := d.GetOrm()
		for _, v1 := range v.Items {
			v1.OrderId = v.Id
			v1.UpdateTime = unix
			if v1.Id > 0 {
				orm.Save(v1.Id, v1)
			} else {
				orm.Save(nil, v1)
			}
		}
	}
	return v.Id, err
}

//　保存订单优惠券绑定
func (this *shoppingRep) SaveOrderCouponBind(val *shopping.OrderCoupon) error {
	_, _, err := this.Connector.GetOrm().Save(nil, val)
	return err
}

// 根据编号获取订单
func (this *shoppingRep) GetOrderById(id int) *shopping.ValueOrder {
	var v = new(shopping.ValueOrder)
	if err := this.Connector.GetOrm().GetBy(v, "id=?", id); err == nil {
		return v
	}
	return nil
}

func (this *shoppingRep) GetOrderByNo(merchantId int, orderNo string) (
	*shopping.ValueOrder, error) {
	var v = new(shopping.ValueOrder)
	err := this.Connector.GetOrm().GetBy(v, "merchant_id=? AND order_no=?", merchantId, orderNo)
	if err != nil {
		return nil, err
	}
	return v, err
}

// 根据订单号获取订单
func (this *shoppingRep) GetValueOrderByNo(orderNo string) *shopping.ValueOrder {
	var v = new(shopping.ValueOrder)
	err := this.Connector.GetOrm().GetBy(v, "order_no=?", orderNo)
	if err == nil {
		return v
	}
	return nil
}

// 获取等待处理的订单
func (this *shoppingRep) GetWaitingSetupOrders(merchantId int) ([]*shopping.ValueOrder, error) {
	dst := []*shopping.ValueOrder{}
	err := this.Connector.GetOrm().Select(&dst,
		"merchant_id=? AND is_suspend=0 AND status IN("+enum.ORDER_SETUP_STATE+")",
		merchantId)
	if err != nil {
		return nil, err
	}
	return dst, err
}

// 保存订单日志
func (this *shoppingRep) SaveOrderLog(v *shopping.OrderLog) error {
	_, _, err := this.Connector.GetOrm().Save(nil, v)
	return err
}

// 获取购物车
func (this *shoppingRep) GetShoppingCart(key string) (*shopping.ValueCart, error) {
	var v = new(shopping.ValueCart)
	err := this.Connector.GetOrm().GetBy(v, "cart_key=?", key)
	if v == nil || err != nil {
		return nil, err
	}

	var items = []*shopping.CartItem{}
	err = this.Connector.GetOrm().Select(&items, "cart_id=?", v.Id)
	if err == nil {
		v.Items = items
	}

	return v, err
}

// 获取最新的购物车
func (this *shoppingRep) GetLatestCart(buyerId int) (*shopping.ValueCart, error) {
	var v = new(shopping.ValueCart)
	err := this.Connector.GetOrm().GetBy(v, "buyer_id=? ORDER BY id DESC", buyerId)
	if v == nil || err != nil {
		return nil, err
	}

	var items = []*shopping.CartItem{}
	err = this.Connector.GetOrm().Select(&items, "cart_id=?", v.Id)
	if err == nil {
		v.Items = items
	}
	return v, err
}

// 保存购物车
func (this *shoppingRep) SaveShoppingCart(v *shopping.ValueCart) (int, error) {
	var err error
	_orm := this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = _orm.Save(v.Id, v)
	} else {
		_, _, err = _orm.Save(nil, v)
		this.Connector.ExecScalar(`SELECT MAX(id) FROM sale_cart`, &v.Id)
	}
	return v.Id, err
}

// 移出购物车项
func (this *shoppingRep) RemoveCartItem(id int) error {
	return this.Connector.GetOrm().DeleteByPk(shopping.CartItem{}, id)
}

// 保存购物车项
func (this *shoppingRep) SaveCartItem(v *shopping.CartItem) (int, error) {
	_orm := this.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = _orm.Save(v.Id, v)
	} else {
		_, _, err = _orm.Save(nil, v)
		this.Connector.ExecScalar(`SELECT MAX(id) FROM sale_cart_item where cart_id=?`, &v.Id, v.CartId)
	}

	return v.Id, err
}

// 清空购物车项
func (this *shoppingRep) EmptyCartItems(id int) error {
	_, err := this.Connector.GetOrm().Delete(shopping.CartItem{}, "cart_id=?", id)
	return err
}

// 删除购物车
func (this *shoppingRep) DeleteCart(id int) error {
	return this.Connector.GetOrm().DeleteByPk(shopping.ValueCart{}, id)
}

// 获取订单的促销绑定
func (this *shoppingRep) GetOrderPromotionBinds(orderNo string) []*shopping.OrderPromotionBind {
	var arr []*shopping.OrderPromotionBind = []*shopping.OrderPromotionBind{}
	err := this.Connector.GetOrm().Select(&arr, "order_no=?", orderNo)
	if err == nil {
		return arr
	}
	return make([]*shopping.OrderPromotionBind, 0)
}

// 保存订单的促销绑定
func (this *shoppingRep) SavePromotionBindForOrder(v *shopping.OrderPromotionBind) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM pt_order_pb WHERE order_no=? AND promotion_id=?",
			&v.Id, v.OrderNo, v.PromotionId)
	}
	return v.Id, err
}
