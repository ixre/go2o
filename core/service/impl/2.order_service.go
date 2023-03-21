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
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/parser"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

var _ proto.OrderServiceServer = new(orderServiceImpl)

type orderServiceImpl struct {
	repo        order.IOrderRepo
	prodRepo    product.IProductRepo
	itemRepo    item.IItemRepo
	cartRepo    cart.ICartRepo
	mchRepo     merchant.IMerchantRepo
	shopRepo    shop.IShopRepo
	manager     order.IOrderManager
	memberRepo  member.IMemberRepo
	payRepo     payment.IPaymentRepo
	shipRepo    shipment.IShipmentRepo
	expressRepo express.IExpressRepo
	orderQuery  *query.OrderQuery
	serviceUtil
	proto.UnimplementedOrderServiceServer
}

func NewShoppingService(r order.IOrderRepo,
	cartRepo cart.ICartRepo, memberRepo member.IMemberRepo,
	prodRepo product.IProductRepo, goodsRepo item.IItemRepo,
	mchRepo merchant.IMerchantRepo, shopRepo shop.IShopRepo,
	payRepo payment.IPaymentRepo, shipRepo shipment.IShipmentRepo,
	expressRepo express.IExpressRepo,
	orderQuery *query.OrderQuery) *orderServiceImpl {
	return &orderServiceImpl{
		repo:       r,
		prodRepo:   prodRepo,
		cartRepo:   cartRepo,
		memberRepo: memberRepo,
		itemRepo:   goodsRepo,
		mchRepo:    mchRepo,
		shopRepo:   shopRepo,
		payRepo:    payRepo,
		shipRepo:   shipRepo,
		manager:    r.Manager(),
		orderQuery: orderQuery,
	}
}

// 获取购物车
func (s *orderServiceImpl) getShoppingCart(buyerId int64, cartCode string) cart.ICart {
	// 本地的购物车
	var ic cart.ICart
	if len(cartCode) > 0 {
		ic = s.cartRepo.GetShoppingCartByKey(cartCode)
	}
	// 获取用户购物车
	mc := s.cartRepo.GetMyCart(buyerId, cart.KNormal)
	if ic != nil {
		// 绑定临时购物车为会员购物车
		if mc == nil {
			ic.Bind(int(buyerId))
			return ic
		}
		// 会员购物车合并临时购物车
		if mc.GetAggregateRootId() != ic.GetAggregateRootId() {
			mc.(cart.INormalCart).Combine(ic)
			_, _ = mc.Save()
		}
	}
	return mc
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
		Type:          order.OrderType(r.OrderType),
		BuyerId:       r.BuyerId,
		AddressId:     r.AddressId,
		Subject:       r.Subject,
		CouponCode:    r.CouponCode,
		BalanceDeduct: r.BalanceDeduct,
		WalletDeduct:  r.WalletDeduct,
		AffiliateCode: r.AffiliateCode,
		PostedData:    iData,
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
	ic := s.getShoppingCart(r.BuyerId, r.CartCode)
	if ic == nil {
		ic = s.cartRepo.NewTempNormalCart(int(r.BuyerId), r.CartCode)
	}
	if r.Item != nil {
		err := ic.Put(r.Item.ItemId, r.Item.SkuId, r.Item.Quantity, true, true)
		if err == nil {
			_, err = ic.Save()
		}
		if err != nil {
			return &proto.PrepareOrderResponse{
				ErrCode: 1,
				ErrMsg:  err.Error(),
			}, nil
		}
	}
	o, err := s.manager.PrepareNormalOrder(ic)
	if err == nil {
		// 设置收货地址
		if r.AddressId > 0 {
			err = o.SetShipmentAddress(r.AddressId)
		} else {
			arr := s.memberRepo.GetDeliverAddress(r.BuyerId)
			if len(arr) > 0 {
				var addressId int64 = 0
				for _, v := range arr {
					if v.IsDefault == 1 {
						err = o.SetShipmentAddress(v.Id)
						addressId = v.Id
						break
					}
				}
				if addressId == 0 {
					err = o.SetShipmentAddress(arr[0].Id)
				}
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
	// 绑定账户信息
	acc := s.memberRepo.GetMember(r.BuyerId).GetAccount()
	balance := acc.GetValue().Balance
	walletBalance := acc.GetValue().WalletBalance
	var deductAmount int64
	// 使用余额
	if fb, fw := domain.TestFlag(int(r.PaymentFlag), payment.MBalance),
		domain.TestFlag(int(r.PaymentFlag), payment.MWallet); fb || fw {
		// 更新抵扣余额之后的金额
		if fb {
			if balance >= ov.FinalAmount {
				deductAmount += ov.FinalAmount
				ov.FinalAmount = 0
			} else {
				deductAmount += balance
				ov.FinalAmount -= balance
			}
		}
		// 更新抵扣钱包余额之后的金额
		if fw {
			if walletBalance >= ov.FinalAmount {
				deductAmount += ov.FinalAmount
				ov.FinalAmount = 0
			} else {
				deductAmount += walletBalance
				ov.FinalAmount -= walletBalance
			}
		}
	}
	rsp := parser.PrepareOrderDto(ov)
	rsp.DeductAmount = deductAmount
	rsp.BuyerBalance = balance
	rsp.BuyerWallet = walletBalance
	rsp.BuyerIntegral = int64(acc.GetValue().Integral)
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
func (s *orderServiceImpl) GetOrder(_ context.Context, r *proto.OrderRequest) (*proto.SSingleOrder, error) {
	if len(r.OrderNo) == 0 {
		return nil, order.ErrNoSuchOrder
	}
	c := s.manager.Unified(r.OrderNo, true).Complex()
	if c != nil {
		ret := parser.OrderDto(c)
		if r.WithDetail {
			// 获取支付单信息
			po := s.payRepo.GetPaymentOrder(r.OrderNo)
			if po != nil {
				pv := po.Get()
				ret.DeductAmount = int32(pv.DeductAmount)
				ret.FinalAmount = int32(pv.FinalAmount)
				ret.ExpiresTime = pv.ExpiresTime
				ret.PayTime = pv.PaidTime
				for _, t := range po.TradeMethods() {
					pm := s.parseTradeMethodDataDto(t)
					pm.ChanName = po.ChanName(t.Method)
					if len(pm.ChanName) == 0 {
						pm.ChanName = pv.OutTradeSp
					}
					ret.TradeData = append(ret.TradeData, pm)
				}
			}
			// 获取发货单信息
			if c.Status >= order.StatShipped && c.Status <= order.StatCompleted {
				list := s.shipRepo.GetShipOrders(c.OrderId, true)
				for _, v := range list {
					// 绑定快递名称
					ex := s.expressRepo.GetExpressProvider(int32(v.Value().SpId))
					if ex != nil {
						ret.ShipExpressName = ex.Name
					}
					ret.ShipLogisticCode = v.Value().SpOrder
				}
			}
		}
		return ret, nil
	}
	return nil, order.ErrNoSuchOrder
}

func (s *orderServiceImpl) parseTradeMethodDataDto(src *payment.TradeMethodData) *proto.SOrderPayChanData {
	return &proto.SOrderPayChanData{
		ChanId:     int32(src.Method),
		Amount:     src.Amount,
		ChanCode:   src.Code,
		OutTradeNo: src.OutTradeNo,
	}
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
	err := c.Cancel(r.IsBuyerCancel, r.Reason)
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
