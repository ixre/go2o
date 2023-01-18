/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-05 17:53
 * description :
 * history :
 */

package impl

import (
	"bytes"
	"context"
	"errors"

	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/parser"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

var _ proto.OrderServiceServer = new(orderServiceImpl)

type orderServiceImpl struct {
	repo       order.IOrderRepo
	prodRepo   product.IProductRepo
	itemRepo   item.IItemRepo
	cartRepo   cart.ICartRepo
	mchRepo    merchant.IMerchantRepo
	shopRepo   shop.IShopRepo
	manager    order.IOrderManager
	memberRepo member.IMemberRepo
	orderQuery *query.OrderQuery
	serviceUtil
	proto.UnimplementedOrderServiceServer
}

func NewShoppingService(r order.IOrderRepo,
	cartRepo cart.ICartRepo, memberRepo member.IMemberRepo,
	prodRepo product.IProductRepo, goodsRepo item.IItemRepo,
	mchRepo merchant.IMerchantRepo, shopRepo shop.IShopRepo,
	orderQuery *query.OrderQuery) *orderServiceImpl {
	return &orderServiceImpl{
		repo:       r,
		prodRepo:   prodRepo,
		cartRepo:   cartRepo,
		memberRepo: memberRepo,
		itemRepo:   goodsRepo,
		mchRepo:    mchRepo,
		shopRepo:   shopRepo,
		manager:    r.Manager(),
		orderQuery: orderQuery,
	}
}

// 获取购物车
func (s *orderServiceImpl) getShoppingCart(buyerId int64, code string) cart.ICart {
	var c cart.ICart
	var cc cart.ICart
	if len(code) > 0 {
		cc = s.cartRepo.GetShoppingCartByKey(code)
	}
	// 如果传入会员编号，则合并购物车
	if buyerId > 0 {
		c = s.cartRepo.GetMyCart(buyerId, cart.KNormal)
		if cc != nil {
			rc := c.(cart.INormalCart)
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
	c = s.cartRepo.NewTempNormalCart(int(buyerId), code)
	//_, err := c.Save()
	//domain.HandleError(err, "service")
	return c
}

// SubmitOrderV1 提交订单
func (s *orderServiceImpl) SubmitOrder(_ context.Context, r *proto.SubmitOrderRequest) (*proto.OrderSubmitResponse, error) {
	iData := parser.NewPostedData(r.Data, r)
	/* 批发订单
	c := s.cartRepo.GetMyCart(r.BuyerId, cart.KWholesale)
	rd, err := s.repo.Manager().SubmitWholesaleOrder(c, iData)
	if err != nil {
		return &proto.StringMap{Value: map[string]string{
			"error": err.Error(),
		}}, nil
	}
	return &proto.StringMap{Value: rd}, nil
	*/
	/*　交易类订单
	if r.Order.ShopId <= 0 {
		mch := s.mchRepo.GetMerchant(int(r.Order.SellerId))
		if mch != nil {
			sp := mch.ShopManager().GetOnlineShop()
			if sp != nil {
				r.Order.ShopId = int64(sp.GetDomainId())
			} else {
				r.Order.ShopId = 1
			}
		}
	}
	io, err := s.manager.SubmitTradeOrder(parser.ParseTradeOrder(r.Order), r.Rate)
	rs := s.result(err)
	rs.Data = map[string]string{
		"OrderId": strconv.Itoa(int(io.GetAggregateRootId())),
	}
	if err == nil {
		// 返回支付单号
		rs.Data["OrderNo"] = io.OrderNo()
		rs.Data["PaymentOrderNo"] = io.GetPaymentOrder().TradeNo()
	}
	return rs, nil
	*/
	_, rd, err := s.manager.SubmitOrder(order.SubmitOrderData{
		Type:            order.OrderType(r.OrderType),
		BuyerId:         r.BuyerId,
		AddressId:       r.AddressId,
		Subject:         r.Subject,
		CouponCode:      r.CouponCode,
		BalanceDiscount: r.BalanceDiscount,
		AffiliateCode:   r.AffiliateCode,
		PostedData:      iData,
	})
	ret := &proto.OrderSubmitResponse{}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	} else {
		ret.OrderNo = rd.OrderNo
		ret.MergePay = rd.MergePay
		ret.TradeNo = rd.TradeNo
		ret.TradeAmount = rd.TradeAmount
		ret.PaymentOrderNo = rd.PaymentOrderNo
	}
	return ret, nil
}

// PrepareOrder 预生成订单
func (s *orderServiceImpl) PrepareOrder(_ context.Context, r *proto.PrepareOrderRequest) (*proto.PrepareOrderResponse, error) {
	ic := s.getShoppingCart(r.BuyerId, r.CouponCode)
	o, err := s.manager.PrepareNormalOrder(ic)
	if err == nil {
		// 设置收货地址
		if r.AddressId > 0 {
			err = o.SetShipmentAddress(r.AddressId)
		} else {
			arr := s.memberRepo.GetDeliverAddress(r.BuyerId)
			if len(arr) > 0 {
				err = o.SetShipmentAddress(arr[0].Id)
			}
		}
		// 使用优惠券
		// todo:
		//io.ApplyCoupon()

	}
	if err != nil {
		return &proto.PrepareOrderResponse{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}, nil
	}
	ov := o.Complex()

	// 使用余额
	if r.PaymentFlag&payment.MBalance == payment.MBalance {
		acc := s.memberRepo.GetMember(r.BuyerId).GetAccount()
		balance := acc.GetValue().Balance
		if balance >= ov.FinalAmount {
			ov.DiscountAmount = ov.FinalAmount
			ov.FinalAmount = 0
		} else {
			ov.DiscountAmount = balance
			ov.FinalAmount -= balance
		}
	}
	rsp := parser.PrepareOrderDto(ov)
	rsp.Sellers = s.parsePrepareItemsFromCart(ic)
	return rsp, err
}

// 转换购物车数据
func (s *orderServiceImpl) parsePrepareItemsFromCart(c cart.ICart) []*proto.SPrepareOrderGroup {
	arr := parser.ParsePrepareOrderGroups(c)
	for _, v := range arr {
		is := s.shopRepo.GetShop(v.ShopId)
		if is != nil {
			io := is.(shop.IOnlineShop)
			v.ShopName = io.GetShopValue().ShopName
		} else {
			for _, it := range v.Items {
				c.Remove(it.ItemId, it.SkuId, it.Quantity)
			}
		}
	}
	return arr
}

// PrepareOrderWithCoupon_ 预生成订单，使用优惠券
func (s *orderServiceImpl) PrepareOrderWithCoupon_(_ context.Context, r *proto.PrepareOrderRequest) (*proto.StringMap, error) {
	cart := s.getShoppingCart(r.BuyerId, r.CartCode)
	o, err := s.manager.PrepareNormalOrder(cart)
	if err != nil {
		return nil, err
	}
	o.SetShipmentAddress(r.AddressId)
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
	data := make(map[string]string)

	// 取消优惠券
	data["totalFee"] = typeconv.Stringify(v.ItemAmount)
	data["fee"] = typeconv.Stringify(v.ItemAmount)
	data["payFee"] = typeconv.Stringify(v.FinalAmount)
	data["discountFee"] = typeconv.Stringify(discountFee)
	data["expressFee"] = typeconv.Stringify(v.ExpressFee)

	// 设置优惠券的信息
	if r.CartCode != "" {
		// 优惠券没有减金额
		if v.DiscountAmount == 0 {
			data["result"] = typeconv.Stringify(v.DiscountAmount != 0)
			data["message"] = "优惠券无效"
		} else {
			// 成功应用优惠券
			data["couponFee"] = typeconv.Stringify(v.DiscountAmount)
			data["couponDescribe"] = buf.String()
		}
	}

	return &proto.StringMap{Value: data}, err
}

// GetParentOrder 根据编号获取订单
func (s *orderServiceImpl) GetParentOrder(c context.Context, req *proto.OrderNoV2) (*proto.SParentOrder, error) {
	io := s.manager.GetOrderByNo(req.Value)
	if io == nil {
		return nil, errors.New("no such order")
	}
	ord := io.Complex()
	return parser.ParentOrderDto(ord), nil
}

// GetOrder 获取订单和商品项信息
func (s *orderServiceImpl) GetOrder(_ context.Context, orderNo *proto.OrderNoV2) (*proto.SSingleOrder, error) {
	if len(orderNo.Value) == 0 {
		return nil, order.ErrNoSuchOrder
	}
	c := s.manager.Unified(orderNo.Value, true).Complex()
	if c != nil {
		return parser.OrderDto(c), nil
	}
	return nil, order.ErrNoSuchOrder
}

// TradeOrderCashPay 交易单现金支付
func (s *orderServiceImpl) TradeOrderCashPay(_ context.Context, orderId *proto.Int64) (ro *proto.Result, err error) {
	o := s.manager.GetOrderById(orderId.Value)
	if o == nil || o.Type() != order.TTrade {
		err = order.ErrNoSuchOrder
	} else {
		io := o.(order.ITradeOrder)
		err = io.CashPay()
	}
	return s.result(err), nil
}

// TradeOrderUpdateTicket 上传交易单发票
func (s *orderServiceImpl) TradeOrderUpdateTicket(_ context.Context, r *proto.TradeOrderTicketRequest) (rs *proto.Result, err error) {
	o := s.manager.GetOrderById(r.OrderId)
	if o == nil || o.Type() != order.TTrade {
		err = order.ErrNoSuchOrder
	} else {
		io := o.(order.ITradeOrder)
		err = io.UpdateTicket(r.Image)
	}
	return s.result(err), nil
}

// CancelOrder 取消订单
func (s *orderServiceImpl) CancelOrder(_ context.Context, r *proto.CancelOrderRequest) (*proto.Result, error) {
	c := s.manager.Unified(r.OrderNo, r.Sub)
	err := c.Cancel(r.Reason)
	return s.error(err), nil
}

// ConfirmOrder 确定订单
func (s *orderServiceImpl) ConfirmOrder(_ context.Context, r *proto.OrderNo) (*proto.Result, error) {
	c := s.manager.Unified(r.OrderNo, r.Sub)
	err := c.Confirm()
	return s.error(err), nil
}

// PickUp 备货完成
func (s *orderServiceImpl) PickUp(_ context.Context, r *proto.OrderNo) (*proto.Result, error) {
	c := s.manager.Unified(r.OrderNo, r.Sub)
	err := c.PickUp()
	return s.error(err), nil
}

// Ship 订单发货,并记录配送服务商编号及单号
func (s *orderServiceImpl) Ship(_ context.Context, r *proto.OrderShipmentRequest) (*proto.Result, error) {
	c := s.manager.Unified(r.OrderNo, r.Sub)
	err := c.Ship(int32(r.ProviderId), r.ShipOrderNo)
	return s.error(err), nil
}

// BuyerReceived 买家收货
func (s *orderServiceImpl) BuyerReceived(_ context.Context, r *proto.OrderNo) (*proto.Result, error) {
	c := s.manager.Unified(r.OrderNo, r.Sub)
	err := c.BuyerReceived()
	return s.error(err), nil
}

// Forbid implements 删除订单
func (s *orderServiceImpl) Forbid(_ context.Context, r *proto.OrderNo) (*proto.Result, error) {
	c := s.manager.Unified(r.OrderNo, r.Sub)
	err := c.Forbid()
	return s.error(err), nil
}

// ChangeConsignee 更改订单收货人信息
func (s *orderServiceImpl) ChangeShipmentAddress(_ context.Context, r *proto.ChangeOrderAddressRequest) (*proto.Result, error) {
	c := s.manager.Unified(r.OrderNo, r.Sub)
	err := c.ChangeShipmentAddress(r.AddressId)
	return s.error(err), nil
}

// LogBytes 获取订单日志
func (s *orderServiceImpl) LogBytes(_ context.Context, r *proto.OrderNo) (*proto.String, error) {
	c := s.manager.Unified(r.OrderNo, r.Sub)
	return &proto.String{
		Value: string(c.LogBytes()),
	}, nil
}

//
//// 根据商品快照获取订单项
//func (_s *orderServiceImpl) GetOrderItemBySnapshotId(orderId int64, snapshotId int32) *order.SubOrderItem {
//	return _s.repo.GetOrderItemBySnapshotId(orderId, snapshotId)
//}

//// 根据商品快照获取订单项数据传输对象
//func (_s *orderServiceImpl) GetOrderItemDtoBySnapshotId(orderId int64, snapshotId int32) *dto.OrderItem {
//	return _s.repo.GetOrderItemDtoBySnapshotId(orderId, snapshotId)
//}
