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
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/query"
	"strings"
)

type shoppingService struct {
	_rep        order.IOrderRep
	_itemRep    item.IItemRep
	_goodsRep   goods.IGoodsRep
	_saleRep    sale.ISaleRep
	_cartRep    cart.ICartRep
	_mchRep     merchant.IMerchantRep
	_manager    order.IOrderManager
	_orderQuery *query.OrderQuery
}

func NewShoppingService(r order.IOrderRep,
	saleRep sale.ISaleRep, cartRep cart.ICartRep,
	itemRep item.IItemRep, goodsRep goods.IGoodsRep,
	mchRep merchant.IMerchantRep, orderQuery *query.OrderQuery) *shoppingService {
	return &shoppingService{
		_rep:        r,
		_itemRep:    itemRep,
		_cartRep:    cartRep,
		_goodsRep:   goodsRep,
		_saleRep:    saleRep,
		_mchRep:     mchRep,
		_manager:    r.Manager(),
		_orderQuery: orderQuery,
	}
}

/*================ 购物车  ================*/

//  获取购物车
func (s *shoppingService) getShoppingCart(buyerId int,
	cartKey string) cart.ICart {
	var c cart.ICart
	if len(cartKey) > 0 {
		c = s._cartRep.GetShoppingCartByKey(cartKey)
	} else if buyerId > 0 {
		c = s._cartRep.GetMemberCurrentCart(buyerId)
	}
	if c == nil {
		c = s._cartRep.NewCart()
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
func (s *shoppingService) GetShoppingCart(memberId int64,
	cartKey string) *dto.ShoppingCart {
	c := s.getShoppingCart(memberId, cartKey)
	return s.parseCart(c)
}

// 创建一个新的购物车
func (s *shoppingService) CreateShoppingCart(memberId int64) *dto.ShoppingCart {
	c := s._cartRep.NewCart()
	c.SetBuyer(memberId)
	return cart.ParseToDtoCart(c)
}

func (s *shoppingService) parseCart(c cart.ICart) *dto.ShoppingCart {
	dto := cart.ParseToDtoCart(c)
	for _, v := range dto.Vendors {
		mch := s._mchRep.GetMerchant(v.VendorId)
		v.VendorName = mch.GetValue().Name
		if v.ShopId > 0 {
			v.ShopName = mch.ShopManager().GetShop(v.ShopId).GetValue().Name
		}
	}
	return dto
}

//todo: 这里响应较慢,性能?
func (s *shoppingService) AddCartItem(memberId int64, cartKey string,
	skuId, num int, checked bool) (*dto.CartItem, error) {
	c := s.getShoppingCart(memberId, cartKey)
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
		snap := s._goodsRep.GetLatestSnapshot(skuId)
		if snap == nil {
			return nil, goods.ErrNoSuchGoods
		}
		tm := s._itemRep.GetValueItem(snap.ItemId)
		// 检测是否开通商城
		mch := s._mchRep.GetMerchant(tm.VendorId)
		if mch == nil {
			return nil, merchant.ErrNoSuchMerchant
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
func (s *shoppingService) SubCartItem(memberId int64,
	cartKey string, goodsId, num int) error {
	cart := s.getShoppingCart(memberId, cartKey)
	err := cart.RemoveItem(goodsId, num)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

// 勾选商品结算
func (s *shoppingService) CartCheckSign(memberId int64,
	cartKey string, arr []int) error {
	cart := s.getShoppingCart(memberId, cartKey)
	return cart.SignItemChecked(arr)
}

// 更新购物车结算
func (s *shoppingService) PrepareSettlePersist(memberId, shopId,
	paymentOpt, deliverOpt, deliverId int) error {
	var cart = s.getShoppingCart(memberId, "")
	err := cart.SettlePersist(shopId, paymentOpt, deliverOpt, deliverId)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

func (s *shoppingService) GetCartSettle(memberId int64,
	cartKey string) *dto.SettleMeta {
	cart := s.getShoppingCart(memberId, cartKey)
	sp, deliver, payOpt, dlvOpt := cart.GetSettleData()
	st := new(dto.SettleMeta)
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
			Address:    strings.Replace(v.Area, " ", "", -1) + v.Address,
		}
	}

	return st
}

func (s *shoppingService) SetBuyerAddress(buyerId int, cartKey string, addressId int) error {
	cart := s.getShoppingCart(buyerId, cartKey)
	return cart.SetBuyerAddress(addressId)
}

/*================ 订单  ================*/

func (s *shoppingService) PrepareOrder(buyerId int, cartKey string) *order.Order {
	cart := s.getShoppingCart(buyerId, cartKey)
	order, _, err := s._manager.PrepareOrder(cart, "", "")
	if err != nil {
		return nil
	}
	return order.GetValue()
}

func (s *shoppingService) PrepareOrder2(buyerId int, cartKey string,
	subject string, couponCode string) (map[string]interface{}, error) {
	cart := s.getShoppingCart(buyerId, cartKey)
	order, py, err := s._manager.PrepareOrder(cart, subject, couponCode)
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

	discountFee := v.GoodsAmount - po.TotalFee + po.SubAmount
	data := make(map[string]interface{})

	//　取消优惠券
	data["totalFee"] = v.GoodsAmount
	data["fee"] = po.TotalFee
	data["payFee"] = po.FinalAmount
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

func (s *shoppingService) SubmitOrder(buyerId int, cartKey string,
	subject string, couponCode string, balanceDiscount bool) (
	orderNo string, paymentTradeNo string, err error) {
	c := s.getShoppingCart(buyerId, cartKey)
	od, py, err := s._manager.SubmitOrder(c, subject, couponCode, balanceDiscount)
	if err != nil {
		return "", "", err
	}
	return od.GetOrderNo(), py.GetTradeNo(), err
}

func (s *shoppingService) SetDeliverShop(orderNo string,
	shopId int) (err error) {
	o := s._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}

	panic("not implement")
	//if err = o.SetShop(shopId); err == nil {
	//	_, err = o.Save()
	//}
	return err
}

// 根据编号获取订单
func (s *shoppingService) GetOrderById(id int) *order.Order {
	return s._rep.GetOrderById(id)
}

func (s *shoppingService) GetOrderByNo(orderNo string) *order.Order {
	order := s._manager.GetOrderByNo(orderNo)
	if order != nil {
		return order.GetValue()
	}
	return nil
}

// 人工付款
func (s *shoppingService) PayForOrderByManager(orderNo string) error {
	o := s._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.CmPaymentWithBalance()
}

// 根据订单号获取订单
func (s *shoppingService) GetValueOrderByNo(orderNo string) *order.Order {
	return s._rep.GetValueOrderByNo(orderNo)
}

// 获取子订单
func (s *shoppingService) GetSubOrder(id int) *order.SubOrder {
	return s._rep.GetSubOrder(id)
}

// 获取子订单
func (s *shoppingService) GetSubOrderByNo(orderNo string) *order.SubOrder {
	return s._rep.GetSubOrderByNo(orderNo)
}

// 获取订单商品项
func (s *shoppingService) GetSubOrderItems(subOrderId int) []*order.OrderItem {
	return s._rep.GetSubOrderItems(subOrderId)
}

// 获取子订单及商品项
func (s *shoppingService) GetSubOrderAndItems(id int) (*order.SubOrder, []*dto.OrderItem) {
	o := s.GetSubOrder(id)
	if o == nil {
		return o, []*dto.OrderItem{}
	}
	return o, s._orderQuery.QueryOrderItems(id)
}

// 获取子订单及商品项
func (s *shoppingService) GetSubOrderAndItemsByNo(orderNo string) (*order.SubOrder, []*dto.OrderItem) {
	o := s.GetSubOrderByNo(orderNo)
	if o == nil {
		return o, []*dto.OrderItem{}
	}
	return o, s._orderQuery.QueryOrderItems(o.Id)
}

// 取消订单
func (s *shoppingService) CancelOrder(subOrderId int, reason string) error {
	o := s._manager.GetSubOrder(subOrderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.Cancel(reason)
}

// 使用余额为订单付款
func (s *shoppingService) PayForOrderWithBalance(orderNo string) error {
	o := s._manager.GetOrderByNo(orderNo)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.PaymentWithBalance()
}

// 确定订单
func (s *shoppingService) ConfirmOrder(id int) error {
	o := s._manager.GetSubOrder(id)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.Confirm()
}

// 获取订单日志
func (s *shoppingService) GetOrderLogString(id int) []byte {
	o := s._manager.GetSubOrder(id)
	if o == nil {
		return []byte("")
	}
	return o.LogBytes()
}

// 根据父订单编号获取购买的商品项
func (s *shoppingService) GetItemsByParentOrderId(orderId int) []*order.OrderItem {
	return s._manager.GetItemsByParentOrderId(orderId)
}

// 备货完成
func (s *shoppingService) PickUp(subOrderId int) error {
	o := s._manager.GetSubOrder(subOrderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.PickUp()
}

// 订单发货,并记录配送服务商编号及单号
func (s *shoppingService) Ship(subOrderId int, spId int, spOrder string) error {
	o := s._manager.GetSubOrder(subOrderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.Ship(spId, spOrder)
}

// 消费者收货
func (s *shoppingService) BuyerReceived(subOrderId int) error {
	o := s._manager.GetSubOrder(subOrderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.BuyerReceived()
}

// 根据商品快照获取订单项
func (s *shoppingService) GetOrderItemBySnapshotId(orderId int, snapshotId int) *order.OrderItem {
	return s._rep.GetOrderItemBySnapshotId(orderId, snapshotId)
}

// 根据商品快照获取订单项数据传输对象
func (s *shoppingService) GetOrderItemDtoBySnapshotId(orderId int, snapshotId int) *dto.OrderItem {
	return s._rep.GetOrderItemDtoBySnapshotId(orderId, snapshotId)
}

func (s *shoppingService) OrderAutoSetup(mchId int64, f func(error)) {
	s._manager.OrderAutoSetup(f)
}
