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
	"context"
	"encoding/json"
	"errors"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/cart"
	proItem "go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/product"
	orderImpl "go2o/core/domain/order"
	"go2o/core/dto"
	"go2o/core/query"
	"go2o/core/service/auto_gen/rpc/order_service"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/thrift/parser"
	"strconv"
	"strings"
)

var _ order_service.OrderService = new(orderServiceImpl)

type orderServiceImpl struct {
	repo       order.IOrderRepo
	prodRepo   product.IProductRepo
	itemRepo   proItem.IGoodsItemRepo
	cartRepo   cart.ICartRepo
	mchRepo    merchant.IMerchantRepo
	manager    order.IOrderManager
	memberRepo member.IMemberRepo
	orderQuery *query.OrderQuery
	serviceUtil
}

func NewShoppingService(r order.IOrderRepo,
	cartRepo cart.ICartRepo, memberRepo member.IMemberRepo,
	prodRepo product.IProductRepo, goodsRepo proItem.IGoodsItemRepo,
	mchRepo merchant.IMerchantRepo,
	orderQuery *query.OrderQuery) *orderServiceImpl {
	return &orderServiceImpl{
		repo:       r,
		prodRepo:   prodRepo,
		cartRepo:   cartRepo,
		memberRepo: memberRepo,
		itemRepo:   goodsRepo,
		mchRepo:    mchRepo,
		manager:    r.Manager(),
		orderQuery: orderQuery,
	}
}

/*------------------ 购物车  ------------------*/

/*---------------- 批发购物车 ----------------*/

// 批发购物车接口
func (s *orderServiceImpl) WholesaleCartV1(ctx context.Context, memberId int64, action string, data map[string]string) (*ttype.Result_, error) {
	//todo: check member
	c := s.cartRepo.GetMyCart(memberId, cart.KWholesale)
	if data == nil {
		data = map[string]string{}
	}
	switch action {
	case "GET":
		return s.wsGetCart(c, data)
	case "MINI":
		return s.wsGetSimpleCart(c, data)
	case "PUT":
		return s.wsPutItem(c, data)
	case "UPDATE":
		return s.wsUpdateItem(c, data)
	case "CHECK":
		return s.wsCheckCart(c, data)
	}
	return s.result(errors.New("unknown action")), nil
}

// 转换勾选字典,数据如：{"1":["10","11"],"2":["20","21"]}
func (s *orderServiceImpl) parseCheckedMap(data string) (m map[int64][]int64) {
	if data != "" && data != "{}" {
		src := map[string][]string{}
		err := json.Unmarshal([]byte(data), &src)
		if err == nil {
			m = map[int64][]int64{}
			for k, v := range src {
				itemId, _ := strconv.Atoi(k)
				var skuList []int64
				for _, v2 := range v {
					skuId, _ := strconv.Atoi(v2)
					skuList = append(skuList, int64(skuId))
				}
				m[int64(itemId)] = skuList
			}
			return m
		}
	}
	return nil
}

// 获取可结算的购物车
func (s *orderServiceImpl) wsGetCart(c cart.ICart, data map[string]string) (*ttype.Result_, error) {
	//统计checked
	checked := s.parseCheckedMap(data["checked"])
	checkout := data["checkout"] == "true"
	v := c.(cart.IWholesaleCart).JdoData(checkout, checked)
	if v != nil {
		for _, v2 := range v.Seller {
			mch := s.mchRepo.GetMerchant(v2.SellerId)
			if mch != nil {
				v2.Data["SellerName"] = mch.GetValue().CompanyName
			}
		}
	}
	return s.success(nil), nil
}

// 获取简易的购物车
func (s *orderServiceImpl) wsGetSimpleCart(c cart.ICart, data map[string]string) (*ttype.Result_, error) {
	size, err := strconv.Atoi(data["size"])
	if err != nil {
		size = 5
	}
	qd := c.(cart.IWholesaleCart).QuickJdoData(size)
	return s.success(map[string]string{"JsonData": qd}), nil
}

// 转换提交到购物车的数据(PUT和UPDATE), 数据如：91:1;92:1
func (s *orderServiceImpl) wsParseCartPostedData(skuData string) (arr []*cart.ItemPair) {
	arr = []*cart.ItemPair{}
	splitArr := strings.Split(skuData, ";")
	for _, str := range splitArr {
		i := strings.Index(str, ":")
		if i == -1 {
			continue
		}
		skuId, err := util.I64Err(strconv.Atoi(str[:i]))
		quantity, err1 := util.I32Err(strconv.Atoi(str[i+1:]))
		if err == nil && err1 == nil {
			arr = append(arr, &cart.ItemPair{
				SkuId:    skuId,
				Quantity: quantity,
			})
		}
	}
	return arr
}

// 生成结算提交的数据(立即购买),skuData参考数据：skuId:num;skuId2;num2
func (s *orderServiceImpl) createCheckedData(itemId int64, arr []*cart.ItemPair) string {
	buf := bytes.NewBufferString("{\"")
	buf.WriteString(strconv.Itoa(int(itemId)))
	buf.WriteString("\":[")
	for i, v := range arr {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\"")
		buf.WriteString(strconv.Itoa(int(v.SkuId)))
		buf.WriteString("\"")
	}
	buf.WriteString("]}")
	return buf.String()
}

// 放入商品，data["Data"]
func (s *orderServiceImpl) wsPutItem(c cart.ICart, data map[string]string) (*ttype.Result_, error) {
	aId := c.GetAggregateRootId()
	itemId, err := util.I64Err(strconv.Atoi(data["ItemId"]))
	arr := s.wsParseCartPostedData(data["Data"])
	for _, v := range arr {
		err = c.Put(itemId, v.SkuId, v.Quantity)
		if err != nil {
			break
		}
	}
	if err == nil {
		_, err = c.Save()
		if err == nil {
			mp := make(map[string]string)
			mp["cartId"] = strconv.Itoa(int(aId))
			mp["checked"] = s.createCheckedData(itemId, arr)
			return s.success(mp), nil
		}
	}
	return s.result(err), nil
}

func (s *orderServiceImpl) wsUpdateItem(c cart.ICart, data map[string]string) (*ttype.Result_, error) {
	itemId, err := util.I64Err(strconv.Atoi(data["ItemId"]))
	arr := s.wsParseCartPostedData(data["Data"])
	for _, v := range arr {
		err = c.Update(itemId, v.SkuId, v.Quantity)
		if err != nil {
			break
		}
	}
	if err == nil {
		_, err = c.Save()
	}
	return s.result(err), nil
}

// 勾选购物车，格式如：1:2;1:5
func (s *orderServiceImpl) wsCheckCart(c cart.ICart, data map[string]string) (*ttype.Result_, error) {
	checked := data["Checked"]
	var arr []*cart.ItemPair
	splitArr := strings.Split(checked, ";")
	for _, str := range splitArr {
		i := strings.Index(str, ":")
		if i == -1 {
			continue
		}
		itemId, err := util.I64Err(strconv.Atoi(str[:i]))
		skuId, err1 := util.I64Err(strconv.Atoi(str[i+1:]))
		if err == nil && err1 == nil {
			arr = append(arr, &cart.ItemPair{
				ItemId: itemId,
				SkuId:  skuId,
			})
		}
	}
	err := c.SignItemChecked(arr)
	return s.result(err), nil
}

/*---------------- 零售购物车 ----------------*/

// 零售购物车接口
func (s *orderServiceImpl) RetailCartV1(ctx context.Context, memberId int64, action string, data map[string]string) (*ttype.Result_, error) {
	//todo: check member
	c := s.cartRepo.GetMyCart(memberId, cart.KWholesale)
	if data == nil {
		data = map[string]string{}
	}
	switch action {
	case "GET":
		return s.wsGetCart(c, data)
	case "MINI":
		return s.wsGetSimpleCart(c, data)
	case "PUT":
		return s.wsPutItem(c, data)
	case "UPDATE":
		return s.wsUpdateItem(c, data)
	case "CHECK":
		return s.wsCheckCart(c, data)
	}
	return s.result(errors.New("unknown action")), nil
}

// 提交订单
func (s *orderServiceImpl) SubmitOrderV1(ctx context.Context, buyerId int64, cartType int32,
	data map[string]string) (map[string]string, error) {
	c := s.cartRepo.GetMyCart(buyerId, cart.KWholesale)
	iData := orderImpl.NewPostedData(data)
	rd, err := s.repo.Manager().SubmitWholesaleOrder(c, iData)
	if err != nil {
		return map[string]string{
			"error": err.Error(),
		}, nil
	}
	return rd, nil
}

//  获取购物车
func (s *orderServiceImpl) getShoppingCart(buyerId int64, code string) cart.ICart {
	var c cart.ICart
	var cc cart.ICart
	if len(code) > 0 {
		cc = s.cartRepo.GetShoppingCartByKey(code)
	}
	// 如果传入会员编号，则合并购物车
	if buyerId > 0 {
		c = s.cartRepo.GetMyCart(buyerId, cart.KRetail)
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
	c = s.cartRepo.NewRetailCart(code)
	//_, err := c.Save()
	//domain.HandleError(err, "service")
	return c
}

// 获取购物车,当购物车编号不存在时,将返回一个新的购物车
func (s *orderServiceImpl) GetShoppingCart(memberId int64,
	cartCode string) *ttype.SShoppingCart {
	c := s.getShoppingCart(memberId, cartCode)
	return s.parseCart(c)
}

// 转换购物车数据
func (s *orderServiceImpl) parseCart(c cart.ICart) *ttype.SShoppingCart {
	dto := cart.ParseToDtoCart(c)
	for _, v := range dto.Shops {

		//todo: 改为不依赖vendor

		mch := s.mchRepo.GetMerchant(v.VendorId)
		if v.ShopId > 0 {
			v.ShopName = mch.ShopManager().
				GetShop(v.ShopId).GetValue().Name
		}
	}
	return dto
}

// 放入购物车
func (s *orderServiceImpl) PutInCart(memberId int64, code string,
	itemId, skuId int64, quantity int32) (*ttype.SShoppingCartItem, error) {
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
func (s *orderServiceImpl) SubCartItem(memberId int64, code string,
	itemId, skuId int64, quantity int32) error {
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
func (s *orderServiceImpl) CartCheckSign(memberId int64,
	cartCode string, arr []*ttype.SShoppingCartItem) error {
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
func (s *orderServiceImpl) PrepareSettlePersist(memberId int64, shopId int32,
	paymentOpt, deliverOpt int32, deliverId int64) error {
	var cart = s.getShoppingCart(memberId, "")
	err := cart.SettlePersist(shopId, paymentOpt, deliverOpt, deliverId)
	if err == nil {
		_, err = cart.Save()
	}
	return err
}

func (s *orderServiceImpl) GetCartSettle(memberId int64,
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
			Tel:  ols.GetShopValue().ServiceTel,
		}
	}

	if deliver != nil {
		v := deliver.GetValue()
		st.Deliver = &dto.SettleDeliverMeta{
			Id:         v.ID,
			PersonName: v.RealName,
			Phone:      v.Phone,
			Address:    strings.Replace(v.Area, " ", "", -1) + v.Address,
		}
	}

	return st
}

func (s *orderServiceImpl) SetBuyerAddress(buyerId int64, cartCode string, addressId int64) error {
	cart := s.getShoppingCart(buyerId, cartCode)
	return cart.SetBuyerAddress(addressId)
}

/*================ 订单  ================*/

func (s *orderServiceImpl) PrepareOrder(buyerId int64, addressId int64,
	cartCode string) (*order.ComplexOrder, error) {
	ic := s.getShoppingCart(buyerId, cartCode)
	o, err := s.manager.PrepareNormalOrder(ic)
	if err == nil {
		no := o.(order.INormalOrder)
		if addressId > 0 {
			err = no.SetAddress(addressId)
		} else {
			arr := s.memberRepo.GetDeliverAddress(buyerId)
			if len(arr) > 0 {
				err = no.SetAddress(arr[0].ID)
			}
		}
	}
	if err == nil {
		//log.Println("-------",o == nil,err)
		return o.Complex(), err
	}
	return nil, err
}

// 预生成订单，使用优惠券
func (s *orderServiceImpl) PrepareOrderWithCoupon(buyerId int64, cartCode string,
	addressId int64, subject string, couponCode string) (map[string]interface{}, error) {
	cart := s.getShoppingCart(buyerId, cartCode)
	o, err := s.manager.PrepareNormalOrder(cart)
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

func (s *orderServiceImpl) SubmitOrder_V1(buyerId int64, cartCode string,
	addressId int64, subject string, couponCode string, balanceDiscount bool) (*order.SubmitReturnData, error) {
	//todo: 临时取消余额支付
	//balanceDiscount = falseorder_manager.go
	c := s.getShoppingCart(buyerId, cartCode)
	_, rd, err := s.manager.SubmitOrder(c, addressId, couponCode, balanceDiscount)
	return rd, err
}

// 根据编号获取订单
func (s *orderServiceImpl) GetOrder(ctx context.Context, orderNo string, sub bool) (*order_service.SComplexOrder, error) {
	c := s.manager.Unified(orderNo, sub).Complex()
	if c != nil {
		return parser.OrderDto(c), nil
	}
	return nil, nil
}

// 获取订单和商品项信息
func (s *orderServiceImpl) GetOrderAndItems(ctx context.Context, orderNo string, sub bool) (*order_service.SComplexOrder, error) {
	c := s.manager.Unified(orderNo, sub).Complex()
	if c != nil {
		return parser.OrderDto(c), nil
	}
	return nil, nil
}

// 根据编号获取订单
func (s *orderServiceImpl) GetOrderById(id int64) *order.ComplexOrder {
	o := s.manager.GetOrderById(id)
	if o != nil {
		return o.Complex()
	}
	return nil
}

func (s *orderServiceImpl) GetOrderByNo(orderNo string) *order.ComplexOrder {
	o := s.manager.GetOrderByNo(orderNo)
	if o != nil {
		return o.Complex()
	}
	return nil
}

// 人工付款
func (s *orderServiceImpl) PayForOrderByManager(orderNo string) error {
	//todo: 对支付单进行人工付款
	panic("应使用支付单进行人工付款")
	//o := s._manager.GetOrderByNo(orderNo)
	//if o == nil {
	//	return order.ErrNoSuchOrder
	//}
	//return o.CmPaymentWithBalance()
}

// 获取子订单
func (s *orderServiceImpl) GetSubOrder(ctx context.Context, id int64) (r *order_service.SComplexOrder, err error) {
	o := s.repo.GetSubOrder(id)
	if o != nil {
		return parser.SubOrderDto(o), nil
	}
	return nil, nil
}

// 根据订单号获取子订单
func (s *orderServiceImpl) GetSubOrderByNo(ctx context.Context, orderNo string) (r *order_service.SComplexOrder, err error) {
	orderId := s.repo.GetOrderId(orderNo, true)
	o := s.repo.GetSubOrder(orderId)
	if o != nil {
		return parser.SubOrderDto(o), nil
	}
	return nil, nil
}

// 获取订单商品项
func (s *orderServiceImpl) GetSubOrderItems(ctx context.Context, subOrderId int64) ([]*order_service.SComplexItem, error) {
	list := s.repo.GetSubOrderItems(subOrderId)
	arr := make([]*order_service.SComplexItem, len(list))
	for i, v := range list {
		arr[i] = parser.SubOrderItemDto(v)
	}
	return arr, nil
}

// 获取子订单及商品项
func (s *orderServiceImpl) GetSubOrderAndItems(id int64) (*order.NormalSubOrder, []*dto.OrderItem) {
	o := s.repo.GetSubOrder(id)
	if o == nil {
		return o, []*dto.OrderItem{}
	}
	return o, s.orderQuery.QueryOrderItems(id)
}

// 获取子订单及商品项
func (s *orderServiceImpl) GetSubOrderAndItemsByNo(orderNo string) (*order.NormalSubOrder, []*dto.OrderItem) {
	orderId := s.repo.GetOrderId(orderNo, true)
	o := s.repo.GetSubOrder(orderId)
	if o == nil {
		return o, []*dto.OrderItem{}
	}
	return o, s.orderQuery.QueryOrderItems(orderId)
}

// 获取订单日志
func (s *orderServiceImpl) LogBytes(orderNo string, sub bool) []byte {
	c := s.manager.Unified(orderNo, sub)
	return c.LogBytes()
}

// 提交订单
func (s *orderServiceImpl) SubmitTradeOrder(ctx context.Context, o *order_service.SComplexOrder, rate float64) (*ttype.Result_, error) {
	if o.ShopId <= 0 {
		mch := s.mchRepo.GetMerchant(o.VendorId)
		if mch != nil {
			sp := mch.ShopManager().GetOnlineShop()
			if sp != nil {
				o.ShopId = sp.GetDomainId()
			} else {
				o.ShopId = 1
			}
		}
	}
	io, err := s.manager.SubmitTradeOrder(parser.Order(o), rate)
	r := s.result(err)
	r.Data = map[string]string{
		"OrderId": strconv.Itoa(int(io.GetAggregateRootId())),
	}
	if err == nil {
		// 返回支付单号
		ro := io.(order.ITradeOrder)
		r.Data["OrderNo"] = io.OrderNo()
		r.Data["PaymentOrderNo"] = ro.GetPaymentOrder().TradeNo()
	}
	return r, nil
}

// 交易单现金支付
func (s *orderServiceImpl) TradeOrderCashPay(ctx context.Context, orderId int64) (r *ttype.Result_, err error) {
	o := s.manager.GetOrderById(orderId)
	if o == nil || o.Type() != order.TTrade {
		err = order.ErrNoSuchOrder
	} else {
		io := o.(order.ITradeOrder)
		err = io.CashPay()
	}
	return s.result(err), nil
}

// 上传交易单发票
func (s *orderServiceImpl) TradeOrderUpdateTicket(ctx context.Context, orderId int64, img string) (r *ttype.Result_, err error) {
	o := s.manager.GetOrderById(orderId)
	if o == nil || o.Type() != order.TTrade {
		err = order.ErrNoSuchOrder
	} else {
		io := o.(order.ITradeOrder)
		err = io.UpdateTicket(img)
	}
	return s.result(err), nil
}

// 取消订单
func (s *orderServiceImpl) CancelOrder(orderNo string, sub bool, reason string) error {
	c := s.manager.Unified(orderNo, sub)
	return c.Cancel(reason)
}

// 确定订单
func (s *orderServiceImpl) ConfirmOrder(orderNo string, sub bool) error {
	c := s.manager.Unified(orderNo, sub)
	return c.Confirm()
}

// 备货完成
func (s *orderServiceImpl) PickUp(orderNo string, sub bool) error {
	c := s.manager.Unified(orderNo, sub)
	return c.PickUp()
}

// 订单发货,并记录配送服务商编号及单号
func (s *orderServiceImpl) Ship(orderNo string, sub bool, spId int32, spOrder string) error {
	c := s.manager.Unified(orderNo, sub)
	return c.Ship(spId, spOrder)
}

// 消费者收货
func (s *orderServiceImpl) BuyerReceived(orderNo string, sub bool) error {
	c := s.manager.Unified(orderNo, sub)
	return c.BuyerReceived()
}

// 根据商品快照获取订单项
func (s *orderServiceImpl) GetOrderItemBySnapshotId(orderId int64, snapshotId int32) *order.SubOrderItem {
	return s.repo.GetOrderItemBySnapshotId(orderId, snapshotId)
}

// 根据商品快照获取订单项数据传输对象
func (s *orderServiceImpl) GetOrderItemDtoBySnapshotId(orderId int64, snapshotId int32) *dto.OrderItem {
	return s.repo.GetOrderItemDtoBySnapshotId(orderId, snapshotId)
}
