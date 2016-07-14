/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package dps

import (
	"bytes"
	"errors"
	"github.com/jsix/gof/log"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
)

type shoppingService struct {
	_rep      order.IOrderRep
	_itemRep  item.IItemRep
	_goodsRep goods.IGoodsRep
	_saleRep  sale.ISaleRep
	_cartRep  cart.ICartRep
	_mchRep   merchant.IMerchantRep
	_manager  order.IOrderManager
}

func NewShoppingService(r order.IOrderRep,
	saleRep sale.ISaleRep, cartRep cart.ICartRep,
	itemRep item.IItemRep, goodsRep goods.IGoodsRep,
	mchRep merchant.IMerchantRep) *shoppingService {
	return &shoppingService{
		_rep:      r,
		_itemRep:  itemRep,
		_cartRep:  cartRep,
		_goodsRep: goodsRep,
		_saleRep:  saleRep,
		_mchRep:   mchRep,
		_manager:  r.Manager(),
	}
}

/*================ 购物车  ================*/

//  获取购物车
func (this *shoppingService) getShoppingCart(buyerId int,
	cartKey string) cart.ICart {
	var c cart.ICart
	if len(cartKey) > 0 {
		c = this._cartRep.GetShoppingCartByKey(cartKey)
	} else if buyerId > 0 {
		c = this._cartRep.GetMemberCurrentCart(buyerId)
	}
	if c == nil {
		c = this._cartRep.NewCart()
		_, err := c.Save()
		domain.HandleError(err, "service")
	}
	if c.GetValue().BuyerId <= 0 {
		err := c.SetBuyer(buyerId)
		domain.HandleError(err, "service")
	}
	return c
}

// 获取购物车,当购物车编号不存在时,将返回一个新的购物车
func (this *shoppingService) GetShoppingCart(memberId int,
	cartKey string) *dto.ShoppingCart {
	c := this.getShoppingCart(memberId, cartKey)
	return this.parseCart(c)
}

// 创建一个新的购物车
func (this *shoppingService) CreateShoppingCart(memberId int) *dto.ShoppingCart {
	c := this._cartRep.NewCart()
	c.SetBuyer(memberId)
	return cart.ParseToDtoCart(c)
}

func (this *shoppingService) parseCart(c cart.ICart) *dto.ShoppingCart {
	dto := cart.ParseToDtoCart(c)
	for _, v := range dto.Vendors {
		mch, _ := this._mchRep.GetMerchant(v.VendorId)
		v.VendorName = mch.GetValue().Name
		if v.ShopId > 0 {
			v.ShopName = mch.ShopManager().GetShop(v.ShopId).GetValue().Name
		}
	}
	return dto
}

//todo: 这里响应较慢,性能?
func (this *shoppingService) AddCartItem(memberId int, cartKey string,
	skuId, num int, checked bool) (*dto.CartItem, error) {
	c := this.getShoppingCart(memberId, cartKey)
	var item *cart.CartItem
	var err error
	// 从购物车中添加
	for k, v := range c.Items() {
		if k == skuId {
			item, err = c.AddItem(v.VendorId, v.ShopId, skuId, num, checked)
			break
		}
	}
	// 将新商品加入到购物车
	if item == nil {
		snap := this._goodsRep.GetLatestSnapshot(skuId)
		if snap == nil {
			return nil, goods.ErrNoSuchGoods
		}
		tm := this._itemRep.GetValueItem(snap.ItemId)

		// 检测是否开通商城
		mch, err2 := this._mchRep.GetMerchant(tm.VendorId)
		if err2 != nil {
			return nil, err2
		}
		shops := mch.ShopManager().GetShops()
		shopId := 0
		for _, v := range shops {
			if v.Type() == shop.TypeOnlineShop {
				shopId = v.GetDomainId()
				break
			}
		}
		if shopId == 0 {
			return nil, errors.New("商户还未开通商城")
		}

		// 加入购物车
		item, err = c.AddItem(snap.VendorId, shopId, skuId, num, checked)
	}

	if err == nil {
		if _, err = c.Save(); err == nil {
			return cart.ParseCartItem(item), err
		}
	}
	return nil, err
}
func (this *shoppingService) SubCartItem(memberId int,
	cartKey string, goodsId, num int) error {
	cart := this.getShoppingCart(memberId, cartKey)
	err := cart.RemoveItem(goodsId, num)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

// 勾选商品结算
func (this *shoppingService) CartCheckSign(memberId int,
	cartKey string, arr []int) error {
	cart := this.getShoppingCart(memberId, cartKey)
	return cart.SignItemChecked(arr)
}

// 更新购物车结算
func (this *shoppingService) PrepareSettlePersist(memberId, shopId,
	paymentOpt, deliverOpt, deliverId int) error {
	var cart = this.getShoppingCart(memberId, "")
	err := cart.SettlePersist(shopId, paymentOpt, deliverOpt, deliverId)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

func (this *shoppingService) GetCartSettle(memberId int,
	cartKey string) *dto.SettleMeta {
	var cart = this.getShoppingCart(memberId, cartKey)
	sp, deliver, payOpt, dlvOpt := cart.GetSettleData()
	var st *dto.SettleMeta = new(dto.SettleMeta)
	st.PaymentOpt = payOpt
	st.DeliverOpt = dlvOpt
	if sp != nil {
		v := sp.GetValue()
		ols := sp.(shop.IOnlineShop)
		st.Shop = &dto.SettleShopMeta{
			Id:   v.Id,
			Name: v.Name,
			Tel:  ols.GetShopValue().Tel,
		}
	}

	if deliver != nil {
		v := deliver.GetValue()
		st.Deliver = &dto.SettleDeliverMeta{
			Id:         v.Id,
			PersonName: v.RealName,
			Phone:      v.Phone,
			Address:    v.Address,
		}
	}

	return st
}

/*================ 订单  ================*/

func (this *shoppingService) PrepareOrder(buyerId int, cartKey string) *order.Order {
	cart := this.getShoppingCart(buyerId, cartKey)
	order, _, err := this._manager.PrepareOrder(cart, "", "")
	if err != nil {
		return nil
	}
	return order.GetValue()
}

func (this *shoppingService) PrepareOrder2(buyerId int, cartKey string,
	subject string, couponCode string) (map[string]interface{}, error) {
	cart := this.getShoppingCart(buyerId, cartKey)
	order, py, err := this._manager.PrepareOrder(cart, subject, couponCode)
	if err != nil {
		return nil, err
	}

	v := order.GetValue()
	po := py.GetValue()
	buf := bytes.NewBufferString("")

	for _, v := range order.GetCoupons() {
		buf.WriteString(v.GetDescribe())
		buf.WriteString("\n")
	}

	discountFee := v.GoodsFee - po.TotalFee + po.SubFee
	data := make(map[string]interface{})

	//　取消优惠券
	data["totalFee"] = v.GoodsFee
	data["fee"] = po.TotalFee
	data["payFee"] = po.FinalFee
	data["discountFee"] = discountFee
	data["expressFee"] = v.ExpressFee

	// 设置优惠券的信息
	if couponCode != "" {
		// 优惠券没有减金额
		if po.CouponDiscount == 0 {
			data["result"] = po.CouponDiscount != 0
			data["message"] = "优惠券无效"
		} else {
			// 成功应用优惠券
			data["couponFee"] = po.CouponDiscount
			data["couponDescribe"] = buf.String()
		}
	}

	return data, err
}

func (this *shoppingService) SubmitOrder(buyerId int, cartKey string,
	subject string, couponCode string, balanceDiscount bool) (
	orderNo string, paymentTradeNo string, err error) {
	c := this.getShoppingCart(buyerId, cartKey)
	od, py, err := this._manager.SubmitOrder(c, subject, couponCode, balanceDiscount)
	if err != nil {
		return "", "", err
	}
	return od.GetOrderNo(), py.GetTradeNo(), err
}

func (this *shoppingService) SetDeliverShop(orderNo string,
	shopId int) (err error) {
	o := this._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}

	panic("not implement")
	//if err = o.SetShop(shopId); err == nil {
	//	_, err = o.Save()
	//}
	return err
}

func (this *shoppingService) HandleOrder(orderNo string) (err error) {
	o := this._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	b := o.IsOver()
	if b {
		return errors.New("订单已经完成!")
	}

	panic("not implement")

	status := o.GetValue().State
	switch status + 1 {
	case enum.ORDER_WAIT_CONFIRM:
		err = o.Confirm()
	case enum.ORDER_WAIT_DELIVERY:
		err = o.Process()
	case enum.ORDER_WAIT_RECEIVE:
		err = o.Deliver(0, "")
		//case enum.ORDER_RECEIVED:
		//	err = o.SignReceived()
		//case enum.ORDER_COMPLETED:
		//	err = o.Complete()
	}
	return err
}

// 根据编号获取订单
func (this *shoppingService) GetOrderById(id int) *order.Order {
	return this._rep.GetOrderById(id)
}

func (this *shoppingService) GetOrderByNo(orderNo string) *order.Order {
	order := this._manager.GetOrderByNo(orderNo)
	if order != nil {
		return order.GetValue()
	}
	return nil
}

// 根据订单号获取订单
func (this *shoppingService) GetValueOrderByNo(orderNo string) *order.Order {
	return this._rep.GetValueOrderByNo(orderNo)
}

// 获取子订单
func (this *shoppingService) GetSubOrder(id int) *order.SubOrder {
	return this._manager.GetSubOrder(id).GetValue()
}

// 获取子订单
func (this *shoppingService) GetSubOrderByNo(orderNo string) *order.SubOrder {
	id := this._rep.GetOrderId(orderNo, true)
	if id <= 0 {
		return nil
	}
	return this.GetSubOrder(id)
}

func (this *shoppingService) GetMinifySubOrder(id int) *order.SubOrder {
	return this._rep.GetSubOrder(id)
}

func (this *shoppingService) GetMinifySubOrderByNo(orderNo string) *order.SubOrder {
	return this._rep.GetSubOrderByNo(orderNo)
}

func (this *shoppingService) CancelOrder(orderId int, reason string) error {
	log.Println("---取消订单", orderId)
	o := this._manager.GetOrderById(orderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}

	panic("not implement")
	return nil
	//return o.Cancel(reason)
}

func (this *shoppingService) CancelOrderByNo(orderNo string, reason string) error {
	o := this._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}

	panic("not implement")
	return nil
	//return o.Cancel(reason)
}

// 使用余额为订单付款
func (this *shoppingService) PayForOrderWithBalance(orderNo string) error {
	o := this._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.PaymentWithBalance()
}

// 人工付款
func (this *shoppingService) PayForOrderByManager(orderNo string) error {
	o := this._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.CmPaymentWithBalance()
}

// 确定订单
func (this *shoppingService) ConfirmOrder(id int) error {
	o := this._manager.GetSubOrder(id)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.Confirm()
}

// 根据父订单编号获取购买的商品项
func (this *shoppingService) GetItemsByParentOrderId(orderId int) []*order.OrderItem {
	return this._manager.GetItemsByParentOrderId(orderId)
}

//todo: 非必须的orderNo改为orderId
// 配送订单,并记录配送服务商编号及单号
func (this *shoppingService) DeliveryOrder(orderNo string,
	deliverySpId int, deliverySpNo string) error {
	//todo:配送订单,并记录配送服务商编号及单号
	o := this._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	if o.GetValue().State == enum.ORDER_WAIT_DELIVERY {
		return o.Deliver(deliverySpId, deliverySpNo)
	}
	return order.ErrOrderDelved
}

// 标记订单已经收货
func (this *shoppingService) SignOrderReceived(orderNo string) error {
	o := this._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	if o.GetValue().State == enum.ORDER_WAIT_RECEIVE {

		panic("not implement")
		return nil
		//return o.SignReceived()
	}
	return nil
}

// 标记订单已经完成
func (this *shoppingService) SignOrderCompleted(orderNo string) error {
	o := this._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	if o.GetValue().State == enum.ORDER_RECEIVED {
		panic("not implement")
		//return o.Complete()
	}
	return nil
}

func (this *shoppingService) OrderAutoSetup(merchantId int, f func(error)) {
	this._manager.OrderAutoSetup(f)
}
