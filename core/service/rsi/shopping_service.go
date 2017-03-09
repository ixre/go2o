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
	"go2o/core/query"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/service/thrift/parser"
	"strings"
)

var _ define.SaleService = new(shoppingService)

type shoppingService struct {
	_repo       order.IOrderRepo
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
		_repo:       r,
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
func (s *shoppingService) getShoppingCart(buyerId int64, code string) cart.ICart {
	var c cart.ICart
	var cc cart.ICart
	if len(code) > 0 {
		cc = s._cartRepo.GetShoppingCartByKey(code)
	}
	// 如果传入会员编号，则合并购物车
	if buyerId > 0 {
		c = s._cartRepo.GetMyCart(buyerId, cart.KRetail)
		if cc != nil {
			rc := c.(cart.IRetailCart)
			rc.Combine(cc)
			c.Save()
		}
		return c
	}
	// 如果只传入code,且购物车存在，直接返回。
	if cc != nil {
		return cc
	}
	// 不存在，则新建购物车
	c = s._cartRepo.NewRetailCart(code)
	//_, err := c.Save()
	//domain.HandleError(err, "service")
	return c
}

// 获取购物车,当购物车编号不存在时,将返回一个新的购物车
func (s *shoppingService) GetShoppingCart(memberId int64,
	cartCode string) *define.ShoppingCart {
	c := s.getShoppingCart(memberId, cartCode)
	return s.parseCart(c)
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
func (s *shoppingService) PutInCart(memberId int64, code string, itemId, skuId,
	quantity int32) (*define.ShoppingCartItem, error) {
	c := s.getShoppingCart(memberId, code)
	if c == nil {
		return nil, cart.ErrNoSuchCart
	}
	err := c.Put(itemId, skuId, quantity)
	if err == nil {
		if _, err = c.Save(); err == nil {
			rc := c.(cart.IRetailCart)
			item := rc.GetItem(itemId, skuId)
			return cart.ParseCartItem(item), err
		}
	}
	return nil, err
}
func (s *shoppingService) SubCartItem(memberId int64, code string, itemId, skuId,
	quantity int32) error {
	c := s.getShoppingCart(memberId, code)
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
func (s *shoppingService) CartCheckSign(memberId int64,
	cartCode string, arr []*define.ShoppingCartItem) error {
	c := s.getShoppingCart(memberId, cartCode)
	items := make([]*cart.ItemPair, len(arr))
	for i, v := range arr {
		items[i] = &cart.ItemPair{
			ItemId:  v.ItemId,
			SkuId:   v.SkuId,
			Checked: 1,
		}
	}
	err := c.SignItemChecked(items)
	if err == nil {
		_, err = c.Save()
	}
	return err
}

// 更新购物车结算
func (s *shoppingService) PrepareSettlePersist(memberId int64, shopId int32,
	paymentOpt, deliverOpt int32, deliverId int64) error {
	var cart = s.getShoppingCart(memberId, "")
	err := cart.SettlePersist(shopId, paymentOpt, deliverOpt, deliverId)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

func (s *shoppingService) GetCartSettle(memberId int64,
	cartCode string) *dto.SettleMeta {
	cart := s.getShoppingCart(memberId, cartCode)
	sp, deliver, payOpt := cart.GetSettleData()
	st := new(dto.SettleMeta)
	st.PaymentOpt = payOpt
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

func (s *shoppingService) SetBuyerAddress(buyerId int64, cartCode string, addressId int64) error {
	cart := s.getShoppingCart(buyerId, cartCode)
	return cart.SetBuyerAddress(addressId)
}

/*================ 订单  ================*/

func (s *shoppingService) PrepareOrder(buyerId int64, addressId int64,
	cartCode string) *order.ComplexOrder {
	cart := s.getShoppingCart(buyerId, cartCode)
	o, err := s._manager.PrepareNormalOrder(cart)
	if err == nil {
		no := o.(order.INormalOrder)
		err = no.SetAddress(addressId)
	}
	return o.Complex()
}

// 预生成订单，使用优惠券
func (s *shoppingService) PrepareOrderWithCoupon(buyerId int64, cartCode string,
	addressId int64, subject string, couponCode string) (map[string]interface{}, error) {
	cart := s.getShoppingCart(buyerId, cartCode)
	o, err := s._manager.PrepareNormalOrder(cart)
	if err != nil {
		return nil, err
	}
	no := o.(order.INormalOrder)
	no.SetAddress(addressId)
	//todo: 应用优惠码
	v := o.Complex()
	buf := bytes.NewBufferString("")

	if o.Type() != order.TRetail {
		panic("not support order type")
	}
	io := o.(order.INormalOrder)
	for _, v := range io.GetCoupons() {
		buf.WriteString(v.GetDescribe())
		buf.WriteString("\n")
	}

	discountFee := v.ItemAmount - v.FinalAmount + v.DiscountAmount
	data := make(map[string]interface{})

	//　取消优惠券
	data["totalFee"] = v.ItemAmount
	data["fee"] = v.ItemAmount
	data["payFee"] = v.FinalAmount
	data["discountFee"] = discountFee
	data["expressFee"] = v.ExpressFee

	// 设置优惠券的信息
	if couponCode != "" {
		// 优惠券没有减金额
		if v.DiscountAmount == 0 {
			data["result"] = v.DiscountAmount != 0
			data["message"] = "优惠券无效"
		} else {
			// 成功应用优惠券
			data["couponFee"] = v.DiscountAmount
			data["couponDescribe"] = buf.String()
		}
	}

	return data, err
}

func (s *shoppingService) SubmitOrder(buyerId int64, cartCode string,
	addressId int64, subject string, couponCode string, balanceDiscount bool) (
	orderNo string, paymentTradeNo string, err error) {
	c := s.getShoppingCart(buyerId, cartCode)
	od, err := s._manager.SubmitOrder(c, addressId, couponCode, balanceDiscount)
	if err != nil {
		return "", "", err
	}
	py := od.(order.INormalOrder).GetPaymentOrder()
	return od.OrderNo(), py.GetTradeNo(), err
}

// 根据编号获取订单
func (s *shoppingService) GetOrder(orderId int64, sub bool) (*define.ComplexOrder, error) {
	c := s._manager.Unified(orderId, sub).Complex()
	if c != nil {
		return parser.OrderDto(c), nil
	}
	return nil, nil
}

// 根据编号获取订单
func (s *shoppingService) GetOrderById(id int64) *order.ComplexOrder {
	o := s._manager.GetOrderById(id)
	if o != nil {
		return o.Complex()
	}
	return nil
}

func (s *shoppingService) GetOrderByNo(orderNo string) *order.ComplexOrder {
	o := s._manager.GetOrderByNo(orderNo)
	if o != nil {
		return o.Complex()
	}
	return nil
}

// 人工付款
func (s *shoppingService) PayForOrderByManager(orderNo string) error {
	//todo: 对支付单进行人工付款
	panic("应使用支付单进行人工付款")
	//o := s._manager.GetOrderByNo(orderNo)
	//if o == nil {
	//	return order.ErrNoSuchOrder
	//}
	//return o.CmPaymentWithBalance()
}

// 根据订单号获取订单
func (s *shoppingService) GetNormalOrderByNo(orderNo string) *order.NormalOrder {
	return s._repo.GetNormalOrderByNo(orderNo)
}

// 获取子订单
func (s *shoppingService) GetSubOrder(id int64) (r *define.ComplexOrder, err error) {
	o := s._repo.GetSubOrder(id)
	if o != nil {
		return parser.SubOrderDto(o), nil
	}
	return nil, nil
}

// 根据订单号获取子订单
func (s *shoppingService) GetSubOrderByNo(orderNo string) (r *define.ComplexOrder, err error) {
	o := s._repo.GetSubOrderByNo(orderNo)
	if o != nil {
		return parser.SubOrderDto(o), nil
	}
	return nil, nil
}

// 获取订单商品项
func (s *shoppingService) GetSubOrderItems(subOrderId int64) ([]*define.ComplexItem, error) {
	list := s._repo.GetSubOrderItems(subOrderId)
	arr := make([]*define.ComplexItem, len(list))
	for i, v := range list {
		arr[i] = parser.SubOrderItemDto(v)
	}
	return arr, nil
}

// 获取子订单及商品项
func (s *shoppingService) GetSubOrderAndItems(id int64) (*order.NormalSubOrder, []*dto.OrderItem) {
	o := s._repo.GetSubOrder(id)
	if o == nil {
		return o, []*dto.OrderItem{}
	}
	return o, s._orderQuery.QueryOrderItems(id)
}

// 获取子订单及商品项
func (s *shoppingService) GetSubOrderAndItemsByNo(orderNo string) (*order.NormalSubOrder, []*dto.OrderItem) {
	o := s._repo.GetSubOrderByNo(orderNo)
	if o == nil {
		return o, []*dto.OrderItem{}
	}
	return o, s._orderQuery.QueryOrderItems(o.ID)
}

// 获取订单日志
func (s *shoppingService) LogBytes(orderId int64, sub bool) []byte {
	c := s._manager.Unified(orderId, sub)
	return c.LogBytes()
}

// 提交订单
func (s *shoppingService) SubmitTradeOrder(o *define.ComplexOrder, rate float64) (r *define.Result64, err error) {
	if o.ShopId <= 0 {
		mch := s._mchRepo.GetMerchant(o.VendorId)
		if mch != nil {
			o.ShopId = mch.ShopManager().GetOnlineShop().GetDomainId()
		}
	}
	io, err := s._manager.SubmitTradeOrder(parser.Order(o), rate)
	r = parser.Result64(io.GetAggregateRootId(), err)
	if err == nil {
		// 返回支付单号
		ro := io.(order.ITradeOrder)
		r.Message = ro.GetPaymentOrder().GetTradeNo()
	}
	return r, nil
}

// 提交订单
func (s *shoppingService) TradeOrderCashPay(orderId int64) (r *define.Result64, err error) {
	o := s._manager.GetOrderById(orderId)
	if o == nil || o.Type() != order.TTrade {
		err = order.ErrNoSuchOrder
	} else {
		io := o.(order.ITradeOrder)
		err = io.CashPay()
	}
	return parser.Result64(o.GetAggregateRootId(), err), nil
}

// 取消订单
func (s *shoppingService) CancelOrder(orderId int64, sub bool, reason string) error {
	c := s._manager.Unified(orderId, sub)
	return c.Cancel(reason)
}

// 确定订单
func (s *shoppingService) ConfirmOrder(orderId int64, sub bool) error {
	c := s._manager.Unified(orderId, sub)
	return c.Confirm()
}

// 备货完成
func (s *shoppingService) PickUp(orderId int64, sub bool) error {
	c := s._manager.Unified(orderId, sub)
	return c.PickUp()
}

// 订单发货,并记录配送服务商编号及单号
func (s *shoppingService) Ship(orderId int64, sub bool, spId int32, spOrder string) error {
	c := s._manager.Unified(orderId, sub)
	return c.Ship(spId, spOrder)
}

// 消费者收货
func (s *shoppingService) BuyerReceived(orderId int64, sub bool) error {
	c := s._manager.Unified(orderId, sub)
	return c.BuyerReceived()
}

// 根据商品快照获取订单项
func (s *shoppingService) GetOrderItemBySnapshotId(orderId int64, snapshotId int32) *order.SubOrderItem {
	return s._repo.GetOrderItemBySnapshotId(orderId, snapshotId)
}

// 根据商品快照获取订单项数据传输对象
func (s *shoppingService) GetOrderItemDtoBySnapshotId(orderId int64, snapshotId int32) *dto.OrderItem {
	return s._repo.GetOrderItemDtoBySnapshotId(orderId, snapshotId)
}
