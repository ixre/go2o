/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package rsi

import (
	"bytes"
	"go2o/core/domain/interface/cart"
	proItem "go2o/core/domain/interface/item"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/product"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/query"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/service/thrift/parser"
	"strings"
)

var _ define.SaleService = new(shoppingService)

type shoppingService struct {
	_rep        order.IOrderRepo
	_itemRepo   product.IProductRepo
	_goodsRepo  proItem.IGoodsItemRepo
	_cartRepo   cart.ICartRepo
	_mchRepo    merchant.IMerchantRepo
	_manager    order.IOrderManager
	_orderQuery *query.OrderQuery
}

func NewShoppingService(r order.IOrderRepo,
	cartRepo cart.ICartRepo,
	itemRepo product.IProductRepo, goodsRepo proItem.IGoodsItemRepo,
	mchRepo merchant.IMerchantRepo, orderQuery *query.OrderQuery) *shoppingService {
	return &shoppingService{
		_rep:        r,
		_itemRepo:   itemRepo,
		_cartRepo:   cartRepo,
		_goodsRepo:  goodsRepo,
		_mchRepo:    mchRepo,
		_manager:    r.Manager(),
		_orderQuery: orderQuery,
	}
}

/*================ 购物车  ================*/

//  获取购物车
func (s *shoppingService) getShoppingCart(buyerId int32,
	cartCode string) cart.ICart {
	var c cart.ICart
	if len(cartCode) > 0 {
		c = s._cartRepo.GetShoppingCartByKey(cartCode)
	} else if buyerId > 0 {
		c = s._cartRepo.GetMemberCurrentCart(buyerId)
	}

	if c == nil {
		c = s._cartRepo.NewCart()
		_, err := c.Save()
		domain.HandleError(err, "service")
	}
	if c.GetValue().BuyerId <= 0 && buyerId > 0 {
		err := c.SetBuyer(buyerId)
		domain.HandleError(err, "service")
	}
	return c
}

// 获取购物车,当购物车编号不存在时,将返回一个新的购物车
func (s *shoppingService) GetShoppingCart(memberId int32,
	cartCode string) *define.ShoppingCart {
	c := s.getShoppingCart(memberId, cartCode)
	return s.parseCart(c)
}

// 创建一个新的购物车
func (s *shoppingService) CreateShoppingCart(memberId int32) *define.ShoppingCart {
	c := s._cartRepo.NewCart()
	c.SetBuyer(memberId)
	return cart.ParseToDtoCart(c)
}

// 转换购物车数据
func (s *shoppingService) parseCart(c cart.ICart) *define.ShoppingCart {
	dto := cart.ParseToDtoCart(c)
	for _, v := range dto.Shops {

		//todo: 改为不依赖vendor

		mch := s._mchRepo.GetMerchant(v.VendorId)
		if v.ShopId > 0 {
			v.ShopName = mch.ShopManager().
				GetShop(v.ShopId).GetValue().Name
		}
	}
	return dto
}

// 放入购物车
func (s *shoppingService) PutInCart(cartId, itemId, skuId,
	quantity int32) (*define.ShoppingCartItem, error) {
	c := s._cartRepo.GetCart(cartId)
	if c == nil {
		return nil, cart.ErrNoSuchCart
	}
	item, err := c.Put(itemId, skuId, quantity)
	if err == nil {
		if _, err = c.Save(); err == nil {
			return cart.ParseCartItem(item), err
		}
	}
	return nil, err
}
func (s *shoppingService) SubCartItem(cartId, itemId, skuId,
	quantity int32) error {
	c := s._cartRepo.GetCart(cartId)
	if c == nil {
		return cart.ErrNoSuchCart
	}
	err := c.Remove(itemId, skuId, quantity)
	if err == nil {
		_, err = c.Save()
	}
	return err
}

// 勾选商品结算
func (s *shoppingService) CartCheckSign(memberId int32,
	cartCode string, arr []*define.ShoppingCartItem) error {
	c := s.getShoppingCart(memberId, cartCode)
	list := make([]*cart.CartItem, len(arr))
	for i, v := range arr {
		list[i] = parser.ShoppingCartItem(v)
	}
	return c.SignItemChecked(list)
}

// 更新购物车结算
func (s *shoppingService) PrepareSettlePersist(memberId, shopId int32,
	paymentOpt, deliverOpt, deliverId int32) error {
	var cart = s.getShoppingCart(memberId, "")
	err := cart.SettlePersist(shopId, paymentOpt, deliverOpt, deliverId)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

func (s *shoppingService) GetCartSettle(memberId int32,
	cartCode string) *dto.SettleMeta {
	cart := s.getShoppingCart(memberId, cartCode)
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

func (s *shoppingService) SetBuyerAddress(buyerId int32, cartCode string, addressId int32) error {
	cart := s.getShoppingCart(buyerId, cartCode)
	return cart.SetBuyerAddress(addressId)
}

/*================ 订单  ================*/

func (s *shoppingService) PrepareOrder(buyerId int32, addressId int32, cartCode string) *order.Order {
	cart := s.getShoppingCart(buyerId, cartCode)
	order, _, err := s._manager.PrepareOrder(cart, addressId, "", "")
	if err != nil {
		return nil
	}
	return order.GetValue()
}

func (s *shoppingService) PrepareOrder2(buyerId int32, cartCode string,
	addressId int32, subject string, couponCode string) (map[string]interface{}, error) {
	cart := s.getShoppingCart(buyerId, cartCode)
	order, py, err := s._manager.PrepareOrder(cart, addressId,
		subject, couponCode)
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

func (s *shoppingService) SubmitOrder(buyerId int32, cartCode string,
	addressId int32, subject string, couponCode string, balanceDiscount bool) (
	orderNo string, paymentTradeNo string, err error) {
	c := s.getShoppingCart(buyerId, cartCode)
	od, py, err := s._manager.SubmitOrder(c, addressId, subject, couponCode, balanceDiscount)
	if err != nil {
		return "", "", err
	}
	return od.GetOrderNo(), py.GetTradeNo(), err
}

func (s *shoppingService) SetDeliverShop(orderNo string,
	shopId int32) (err error) {
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
func (s *shoppingService) GetOrderById(id int32) *order.Order {
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
func (s *shoppingService) GetSubOrder(id int32) (r *define.SubOrder, err error) {
	o := s._rep.GetSubOrder(id)
	if o != nil {
		return parser.SubOrderDto(o), nil
	}
	return nil, nil
}

// 根据订单号获取子订单
func (s *shoppingService) GetSubOrderByNo(orderNo string) (r *define.SubOrder, err error) {
	o := s._rep.GetSubOrderByNo(orderNo)
	if o != nil {
		return parser.SubOrderDto(o), nil
	}
	return nil, nil
}

// 获取订单商品项
func (s *shoppingService) GetSubOrderItems(subOrderId int32) ([]*define.OrderItem, error) {
	list := s._rep.GetSubOrderItems(subOrderId)
	arr := make([]*define.OrderItem, len(list))
	for i, v := range list {
		arr[i] = parser.OrderItemDto(v)
	}
	return arr, nil
}

// 获取子订单及商品项
func (s *shoppingService) GetSubOrderAndItems(id int32) (*order.SubOrder, []*dto.OrderItem) {
	o := s._rep.GetSubOrder(id)
	if o == nil {
		return o, []*dto.OrderItem{}
	}
	return o, s._orderQuery.QueryOrderItems(id)
}

// 获取子订单及商品项
func (s *shoppingService) GetSubOrderAndItemsByNo(orderNo string) (*order.SubOrder, []*dto.OrderItem) {
	o := s._rep.GetSubOrderByNo(orderNo)
	if o == nil {
		return o, []*dto.OrderItem{}
	}
	return o, s._orderQuery.QueryOrderItems(o.ID)
}

// 取消订单
func (s *shoppingService) CancelOrder(subOrderId int32, reason string) error {
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
func (s *shoppingService) ConfirmOrder(id int32) error {
	o := s._manager.GetSubOrder(id)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.Confirm()
}

// 获取订单日志
func (s *shoppingService) GetOrderLogString(id int32) []byte {
	o := s._manager.GetSubOrder(id)
	if o == nil {
		return []byte("")
	}
	return o.LogBytes()
}

// 根据父订单编号获取购买的商品项
func (s *shoppingService) GetItemsByParentOrderId(orderId int32) []*order.OrderItem {
	return s._manager.GetItemsByParentOrderId(orderId)
}

// 备货完成
func (s *shoppingService) PickUp(subOrderId int32) error {
	o := s._manager.GetSubOrder(subOrderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.PickUp()
}

// 订单发货,并记录配送服务商编号及单号
func (s *shoppingService) Ship(subOrderId int32, spId int32, spOrder string) error {
	o := s._manager.GetSubOrder(subOrderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.Ship(spId, spOrder)
}

// 消费者收货
func (s *shoppingService) BuyerReceived(subOrderId int32) error {
	o := s._manager.GetSubOrder(subOrderId)
	if o == nil {
		return order.ErrNoSuchOrder
	}
	return o.BuyerReceived()
}

// 根据商品快照获取订单项
func (s *shoppingService) GetOrderItemBySnapshotId(orderId int32, snapshotId int32) *order.OrderItem {
	return s._rep.GetOrderItemBySnapshotId(orderId, snapshotId)
}

// 根据商品快照获取订单项数据传输对象
func (s *shoppingService) GetOrderItemDtoBySnapshotId(orderId int32, snapshotId int32) *dto.OrderItem {
	return s._rep.GetOrderItemDtoBySnapshotId(orderId, snapshotId)
}

func (s *shoppingService) OrderAutoSetup(mchId int32, f func(error)) {
	s._manager.OrderAutoSetup(f)
}
