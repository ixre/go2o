/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 15:03
 * description :
 * history :
 */

package order

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/shipment"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/util"
)

var _ order.IOrder = new(normalOrderImpl)
var _ order.INormalOrder = new(normalOrderImpl)

//todo: 促销

type normalOrderImpl struct {
	*baseOrderImpl
	manager         order.IOrderManager
	cart            cart.ICart //购物车,仅在订单生成时设置
	coupons         []promotion.ICouponPromotion
	availPromotions []promotion.IPromotion
	orderPbs        []*order.OrderPromotionBind
	orderRepo       order.IOrderRepo
	expressRepo     express.IExpressRepo
	payRepo         payment.IPaymentRepo
	goodsRepo       item.IItemRepo
	productRepo     product.IProductRepo
	promRepo        promotion.IPromotionRepo
	registryRepo    registry.IRegistryRepo
	valRepo         valueobject.IValueRepo
	cartRepo        cart.ICartRepo
	shopRepo        shop.IShopRepo
	// 运营商商品映射,用于整理购物车
	vendorItemsMap map[int][]*order.SubOrderItem
	// 运营商与邮费的MAP
	vendorExpressMap map[int]int64
	// 是否为内部挂起
	internalSuspend bool
	_list           []order.ISubOrder
	_payOrder  payment.IPaymentOrder
}

func newNormalOrder(shopping order.IOrderManager, base *baseOrderImpl,
	shoppingRepo order.IOrderRepo, goodsRepo item.IItemRepo, productRepo product.IProductRepo,
	promRepo promotion.IPromotionRepo, expressRepo express.IExpressRepo, payRepo payment.IPaymentRepo,
	cartRepo cart.ICartRepo, shopRepo shop.IShopRepo, registryRepo registry.IRegistryRepo,
	valRepo valueobject.IValueRepo) order.IOrder {
	return &normalOrderImpl{
		baseOrderImpl: base,
		manager:       shopping,
		promRepo:      promRepo,
		orderRepo:     shoppingRepo,
		goodsRepo:     goodsRepo,
		productRepo:   productRepo,
		valRepo:       valRepo,
		expressRepo:   expressRepo,
		payRepo:       payRepo,
		cartRepo:      cartRepo,
		shopRepo:      shopRepo,
		registryRepo:  registryRepo,
	}
}

func (o *normalOrderImpl) getBaseOrder() *baseOrderImpl {
	return o.baseOrderImpl
}

// Complex 复合的订单信息
func (o *normalOrderImpl) Complex() *order.ComplexOrder {
	v := o.baseValue
	co := o.baseOrderImpl.Complex()
	co.VendorId = 0
	co.ShopId = 0
	co.ParentOrderId = 0
	co.Consignee = &order.ComplexConsignee{
		ConsigneeName:   v.ConsigneeName,
		ConsigneePhone:  v.ConsigneePhone,
		ShippingAddress: v.ShippingAddress,
	}
	co.DiscountAmount = v.DiscountAmount
	co.ItemCount = v.ItemCount
	co.ItemAmount = v.ItemAmount
	co.ExpressFee = v.ExpressFee
	co.PackageFee = v.PackageFee
	co.FinalAmount = v.FinalAmount
	co.IsBreak = int32(v.IsBreak)
	co.UpdateTime = v.UpdateTime
	return co
}

// ApplyCoupon 应用优惠券
func (o *normalOrderImpl) ApplyCoupon(coupon promotion.ICouponPromotion) error {
	//if o._coupons == nil {
	//	o._coupons = []promotion.ICouponPromotion{}
	//}
	//o._coupons = append(o._coupons, coupon)

	// 添加到促销信息中
	if o.orderPbs == nil {
		o.orderPbs = []*order.OrderPromotionBind{}
	}
	for _, v := range o.orderPbs {
		if v.PromotionId == coupon.GetDomainId() {
			return order.ErrPromotionApplied
		}
	}

	o.orderPbs = append(o.orderPbs, &order.OrderPromotionBind{
		Id: 0,
		// 订单号
		OrderId: 0,
		// 促销编号
		PromotionId: coupon.GetDomainId(),
		// 促销类型
		PromotionType: coupon.(promotion.IPromotion).Type(),
		// 标题
		Title: coupon.GetDescribe(),
		// 节省金额
		SaveFee: int64(coupon.GetCouponFee(int(o.baseValue.ItemAmount))),
		// 赠送积分
		PresentIntegral: 0, //todo;/////
		// 是否应用
		IsApply: 0,
		// 是否确认
		IsConfirm: 0,
	})

	//v := o._value
	//v.CouponCode = val.Code
	//v.CouponDescribe = coupon.GetDescribe()
	//v.CouponFee = coupon.GetCouponFee(v.TotalAmount)
	//v.PayFee = o.GetPaymentFee()
	//v.DiscountFee = v.DiscountFee + v.CouponFee
	return nil
}

// 获取支付金额
//func (o *orderImpl) GetPaymentFee() float32 {
//	return o._value.PayFee - o._value.CouponFee
//}

// GetCoupons 获取应用的优惠券
func (o *normalOrderImpl) GetCoupons() []promotion.ICouponPromotion {
	if o.coupons == nil {
		return make([]promotion.ICouponPromotion, 0)
	}
	return o.coupons
}

// GetAvailableOrderPromotions 获取可用的促销,不包含优惠券
func (o *normalOrderImpl) GetAvailableOrderPromotions() []promotion.IPromotion {
	if o.availPromotions == nil {
		//mchId := o._cart.VendorId

		//todo: 将购物车中的vendor均获取出来
		var mchId int64 = -1
		var vp = o.promRepo.GetPromotionOfMerchantOrder(mchId)
		var proms = make([]promotion.IPromotion, len(vp))
		for i, v := range vp {
			proms[i] = o.promRepo.CreatePromotion(v)
		}
		return proms
	}
	return o.availPromotions
}

// GetPromotionBinds 获取促销绑定
func (o *normalOrderImpl) GetPromotionBinds() []*order.OrderPromotionBind {
	if o.orderPbs == nil {
		orderNo := o.OrderNo()
		o.orderPbs = o.orderRepo.GetOrderPromotionBinds(orderNo)
	}
	return o.orderPbs
}

// GetBestSavePromotion 获取最省的促销
func (o *normalOrderImpl) GetBestSavePromotion() (p promotion.IPromotion, saveFee float32, integral int) {
	//todo: not implement
	return nil, 0, 0
}

//************* 订单提交 ***************//

// RequireCart 读取购物车数据,用于预生成订单
func (o *normalOrderImpl) RequireCart(c cart.ICart) error {
	if o.GetAggregateRootId() > 0 || o.cart != nil {
		return order.ErrRequireCart
	}
	if c.Kind() != cart.KNormal {
		panic("购物车非零售")
	}
	rc := c.(cart.INormalCart)
	items := rc.Value().Items
	if len(items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	// 绑定结算购物车
	o.cart = c
	// 将购物车的商品分类整理
	o.vendorItemsMap = o.buildVendorItemMap(items)
	if len(o.vendorItemsMap) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	// 更新订单的金额
	o.vendorExpressMap = o.updateOrderFee(o.vendorItemsMap)
	return nil
}

// 加入运费计算器
func (o *normalOrderImpl) addItemToExpressCalculator(ue express.IUserExpress,
	item *order.SubOrderItem, cul express.IExpressCalculator) {
	tid := int(item.ExpressTplId)
	tpl := ue.GetTemplate(tid)
	if tpl != nil {
		var err error
		v := tpl.Value()
		switch v.Basis {
		case express.BasisByNumber:
			err = cul.Add(tid, int(item.Quantity))
		case express.BasisByWeight:
			err = cul.Add(tid, int(item.Weight))
		case express.BasisByVolume:
			err = cul.Add(tid, int(item.Weight))
		}
		if err != nil {
			log.Println("[ Order][ Express][ Error]:", err)
		}
	}
}

// 更新订单金额,并返回运费
func (o *normalOrderImpl) updateOrderFee(mp map[int][]*order.SubOrderItem) map[int]int64 {
	o.baseValue.ItemAmount = 0
	expCul := make(map[int]express.IExpressCalculator)
	expressMap := make(map[int]int64)
	for k, v := range mp {
		userExpress := o.expressRepo.GetUserExpress(k)
		expCul[k] = userExpress.CreateCalculator()
		for _, it := range v {
			o.baseValue.ItemCount += int(it.Quantity)
			//计算商品总金额
			o.baseValue.ItemAmount += it.Amount
			//计算商品优惠金额
			o.baseValue.DiscountAmount += it.Amount - it.FinalAmount
			//加入运费计算器
			o.addItemToExpressCalculator(userExpress, it, expCul[k])
		}
		//计算商户的运费
		expCul[k].Calculate("") //todo: 传入城市地区编号
		expressMap[k] = expCul[k].Total()
		//叠加运费
		o.baseValue.ExpressFee += expressMap[k]
	}
	o.baseValue.PackageFee = 0
	//计算最终金额
	o.baseValue.FinalAmount = o.baseValue.ItemAmount - o.baseValue.DiscountAmount +
		o.baseValue.ExpressFee + o.baseValue.PackageFee
	return expressMap
}

// GetByVendor 根据运营商获取商品和运费信息,限未生成的订单
func (o *normalOrderImpl) GetByVendor() (items map[int][]*order.SubOrderItem,
	expressFeeMap map[int]int64) {
	if o.vendorItemsMap == nil {
		panic("订单尚未读取购物车!")
	}
	if o.vendorExpressMap == nil {
		panic("订单尚未计算金额")
	}
	items = o.vendorItemsMap
	expressFeeMap = o.vendorExpressMap
	return items, expressFeeMap
}

// 检查购物车
func (o *normalOrderImpl) checkCart() error {
	if o.cart == nil {
		return cart.ErrEmptyShoppingCart
	}
	switch o.cart.Kind() {
	case cart.KNormal:
		rc := o.cart.(cart.INormalCart)
		if l := len(rc.Items()); l == 0 {
			return cart.ErrEmptyShoppingCart
		}
	default:
		panic("购物车非零售")
	}
	return o.cart.Prepare()
}

// 生成运营商与订单商品的映射
func (o *normalOrderImpl) buildVendorItemMap(items []*cart.NormalCartItem) map[int][]*order.SubOrderItem {
	mp := make(map[int][]*order.SubOrderItem)
	for _, v := range items {
		//必须勾选为结算
		if v.Checked == 1 {
			item := o.parseCartToOrderItem(v)
			if item == nil {
				_ = domain.HandleError(errors.New("转换购物车商品到订单商品时出错: 商品SKU"+
					strconv.Itoa(int(v.SkuId))), "domain")
				continue
			}
			list, ok := mp[int(v.VendorId)]
			if !ok {
				list = []*order.SubOrderItem{}
			}
			mp[int(v.VendorId)] = append(list, item)
			//log.Println("--- vendor map len:", len(mp[v.VendorId]))
		}
	}
	return mp
}

// 转换购物车的商品项为订单项目
func (o *normalOrderImpl) parseCartToOrderItem(c *cart.NormalCartItem) *order.SubOrderItem {
	// 获取商品已销售快照
	snap := o.goodsRepo.SnapshotService().GetLatestSalesSnapshot(c.ItemId, c.SkuId)
	if snap == nil {
		_ = domain.HandleError(errors.New("商品快照生成失败："+
			strconv.Itoa(int(c.SkuId))), "domain")
		return nil
	}

	fee := c.Sku.Price * int64(c.Quantity)
	return &order.SubOrderItem{
		ID:          0,
		VendorId:    c.VendorId,
		ShopId:      c.ShopId,
		SkuId:       c.SkuId,
		ItemId:      c.ItemId,
		SnapshotId:  snap.Id,
		Quantity:    c.Quantity,
		Amount:      fee,
		FinalAmount: fee,
		//是否配送
		IsShipped: 0,
		// 退回数量
		ReturnQuantity: 0,
		ExpressTplId:   c.Sku.ExpressTid,
		Weight:         c.Sku.Weight * c.Quantity, //计算重量
		Bulk:           c.Sku.Bulk * c.Quantity,   //计算体积
	}
}

// 检查买家及收货地址
func (o *normalOrderImpl) checkBuyer() error {
	buyer := o.Buyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	if buyer.GetValue().State == 0 {
		return member.ErrMemberLocked
	}
	if o.baseValue.ShippingAddress == "" ||
		o.baseValue.ConsigneePhone == "" ||
		o.baseValue.ConsigneeName == "" {
		return order.ErrMissingShipAddress
	}
	return nil
}

// Submit 提交订单，返回订单号。如有错误则返回
func (o *normalOrderImpl) Submit() error {
	if o.GetAggregateRootId() > 0 {
		return errors.New("订单不允许重复提交")
	}
	err := o.checkBuyer()
	if err == nil {
		err = o.checkCart()
	}
	if err != nil {
		return err
	}
	v := o.baseOrderImpl.baseValue
	//todo: best promotion , 优惠券和返现这里需要重构,直接影响到订单金额
	//prom,fee,integral := o.GetBestSavePromotion()

	// 应用优惠券
	if err := o.applyCouponOnSubmit(); err != nil {
		return err
	}

	// 判断商品的优惠促销,如返现等
	proms, fee := o.applyCartPromotionOnSubmit(o.cart)
	if len(proms) != 0 {
		v.DiscountAmount += int64(fee)
		v.FinalAmount = v.ItemAmount - v.DiscountAmount
		if v.FinalAmount < 0 {
			// 如果出现优惠券多余的金额也一并使用
			v.FinalAmount = 0
		}
	}
	// 均摊优惠折扣到商品
	o.avgDiscountToItem()
	// 保存订单
	err = o.baseOrderImpl.Submit()
	if err != nil {
		return err
	}
	orderNo := o.OrderNo()
	// 保存订单
	err = o.destroyCart()
	if err == nil {
		// 绑定优惠券促销
		o.bindCouponOnSubmit(orderNo)
		// 绑定购物车商品的促销
		for _, p := range proms {
			_, _ = o.bindPromotionOnSubmit(orderNo, p)
		}
		// 扣除库存
		o.applyItemStock()
		// 拆单
		o.breakUpByVendor()
		// 生成支付单
		err = o.createPaymentForOrder()
	}
	return err
}


// GetPaymentOrder implements order.IOrder
func (o *normalOrderImpl) GetPaymentOrder() payment.IPaymentOrder {
    if o._payOrder == nil {
		if o.GetAggregateRootId() <= -1 {
			panic(" Get payment order error ; because of order no yet created!")
		}
		o._payOrder = o.payRepo.GetPaymentOrderByOrderNo(
			int(order.TRetail), o.OrderNo())
	}
	return o._payOrder	
}


// BuildCart 通过订单创建购物车
func (o *normalOrderImpl) BuildCart() cart.ICart {
	bv := o.baseOrderImpl.baseValue
	//v := o.value
	unix := time.Now().Unix()
	vc := &cart.NormalCart{
		BuyerId:    bv.BuyerId,
		PaymentOpt: 1,
		CreateTime: unix,
		UpdateTime: unix,
		Items:      []*cart.NormalCartItem{},
	}
	for _, s := range o.GetSubOrders() {
		for _, v := range s.Items() {
			vc.Items = append(vc.Items, &cart.NormalCartItem{
				VendorId: s.GetValue().VendorId,
				ShopId:   s.GetValue().ShopId,
				ItemId:   v.ItemId,
				SkuId:    v.SkuId,
				Quantity: v.Quantity,
				Checked:  1,
			})
		}
	}
	return o.cartRepo.CreateNormalCart(vc)
}

// 平均优惠抵扣金额到商品
func (o *normalOrderImpl) avgDiscountToItem() {
	if o.vendorItemsMap == nil {
		panic(errors.New("仅能在下单时进行商品抵扣均分"))
	}
	if o.baseValue.DiscountAmount > 0 {
		totalFee := o.baseValue.ItemAmount
		disFee := o.baseValue.DiscountAmount
		for _, items := range o.vendorItemsMap {
			for _, v := range items {
				v.FinalAmount = v.Amount - (v.Amount/totalFee)*disFee
			}
		}
	}
}

// 为所有子订单生成支付单
func (o *normalOrderImpl) createPaymentForOrder() error {
	v := o.baseOrderImpl.baseValue
	itemAmount := v.ItemAmount
	finalAmount := v.FinalAmount
	disAmount := v.DiscountAmount
	po := &payment.Order{
		SellerId:       0,
		TradeNo:        v.OrderNo,
		SubOrder:       0,
		OrderType:      int(order.TRetail),
		OutOrderNo:     v.OrderNo,
		Subject:        v.Subject,
		BuyerId:        v.BuyerId,
		PayUid:         v.BuyerId,
		ItemAmount:     itemAmount,
		DiscountAmount: disAmount,
		DeductAmount:   0,
		AdjustAmount:   0,
		FinalFee:       finalAmount,
		PayFlag:        payment.PAllFlag,
		TradeChannel:   0,
		ExtraData:      "",
		OutTradeSp:     "",
		OutTradeNo:     "",
		State:          payment.StateAwaitingPayment,
		SubmitTime:     v.CreateTime,
		ExpiresTime:    0,
		PaidTime:       0,
		UpdateTime:     v.CreateTime,
		TradeMethods:   []*payment.TradeMethodData{},
	}
	o._payOrder = o.payRepo.CreatePaymentOrder(po)
	err := o._payOrder.Submit()
	if err != nil {
		orders := o.GetSubOrders()
		for _, it := range orders {
			_ = it.Cancel("")
		}
	}
	return err
}

// 绑定促销优惠
func (o *normalOrderImpl) bindPromotionOnSubmit(orderNo string,
	prom promotion.IPromotion) (int32, error) {
	var title string
	var integral int
	var fee int64

	//todo: 需要重构,其他促销
	if prom.Type() == promotion.TypeFlagCashBack {
		fee = prom.GetRelationValue().(*promotion.ValueCashBack).BackFee
		title = prom.TypeName() + ":" + prom.GetValue().ShortName
	}

	v := &order.OrderPromotionBind{
		PromotionId:     prom.GetAggregateRootId(),
		PromotionType:   prom.Type(),
		OrderId:         int32(o.GetAggregateRootId()),
		Title:           title,
		SaveFee:         fee,
		PresentIntegral: integral,
		IsConfirm:       1,
		IsApply:         0,
	}
	return o.orderRepo.SavePromotionBindForOrder(v)
}

// 应用购物车内商品的促销
func (o *normalOrderImpl) applyCartPromotionOnSubmit(cart cart.ICart) ([]promotion.IPromotion, int) {
	//todo: 促销
	var proms = make([]promotion.IPromotion, 0)
	//var prom promotion.IPromotion
	//var saveFee int
	var totalSaveFee int
	//var intOrderFee = int(vo.FinalFee)
	//var rightBack bool
	//
	//for _, v := range cart.GetCartGoods() {
	//	prom = nil
	//	saveFee = 0
	//	rightBack = false
	//
	//	// 判断商品的最省促销
	//	for _, v1 := range v.GetPromotions() {
	//
	//		// 返现
	//		if v1.DbType() == promotion.TypeFlagCashBack {
	//			vc := v1.GetRelationValue().(*promotion.ValueCashBack)
	//			if vc.MinFee < intOrderFee {
	//				if vc.BackFee > saveFee {
	//					prom = v1
	//					saveFee = vc.BackFee
	//					rightBack = vc.BackType == promotion.BackUseForOrder // 是否立即抵扣
	//				}
	//			}
	//		}
	//
	//		//todo: 其他促销
	//	}
	//
	//	if prom != nil {
	//		proms = append(proms, prom)
	//		if rightBack {
	//			totalSaveFee += saveFee
	//		}
	//	}
	//}

	return proms, totalSaveFee
}

func (o *normalOrderImpl) cloneCoupon(src *order.OrderCoupon, coupon promotion.ICouponPromotion,
	orderId int32, orderFee int) *order.OrderCoupon {
	v := coupon.GetDetailsValue()
	src.CouponCode = v.Code
	src.CouponId = v.Id
	src.OrderId = orderId
	src.Fee = int64(coupon.GetCouponFee(orderFee))
	src.Describe = coupon.GetDescribe()
	src.SendIntegral = v.Integral
	return src
}

// 绑定订单与优惠券
func (o *normalOrderImpl) bindCouponOnSubmit(orderNo string) {
	var oc = new(order.OrderCoupon)
	for _, c := range o.GetCoupons() {
		o.cloneCoupon(oc, c, int32(o.GetAggregateRootId()),
			int(o.baseValue.FinalAmount))
		o.orderRepo.SaveOrderCouponBind(oc)
		// 绑定促销
		o.bindPromotionOnSubmit(orderNo, c.(promotion.IPromotion))
	}
}

// 在提交订单时应用优惠券
func (o *normalOrderImpl) applyCouponOnSubmit() error {
	var err error
	var t *promotion.ValueCouponTake
	var b *promotion.ValueCouponBind
	buyerId := o.buyer.GetAggregateRootId()
	for _, c := range o.GetCoupons() {
		if c.CanTake() {
			t, err = c.GetTake(buyerId)
			if err == nil {
				err = c.ApplyTake(t.Id)
			}
		} else {
			b, err = c.GetBind(buyerId)
			if err == nil {
				err = c.UseCoupon(b.Id)
			}
		}
		if err != nil {
			return errors.New("Code 105:优惠券使用失败," + err.Error())
		}
	}
	return err
}

// 应用余额支付
func (o *normalOrderImpl) getBalanceDiscountFee(acc member.IAccount) int64 {
	if o.baseValue.FinalAmount <= 0 || math.IsNaN(float64(o.baseValue.FinalAmount)) {
		return 0
	}
	acv := acc.GetValue()
	if acv.Balance >= o.baseValue.FinalAmount {
		return o.baseValue.FinalAmount
	}
	return acv.Balance
}

// 获取Json格式的商品数据
func (o *normalOrderImpl) getJsonItems() []byte {
	//todo:??? 订单商品JSON表示
	return []byte("{}")
	//var goods []*order.OrderGoods = make([]*order.OrderGoods, len(c.value.Items))
	//for i, v := range c.cart.Items {
	//	goods[i] = &order.OrderGoods{
	//		GoodsId:    v.SkuId,
	//		GoodsImage: v.Sku.Image,
	//		Quantity:   v.Quantity,
	//		Name:       v.Sku.Title,
	//	}
	//}
	//d, _ := json.Marshal(goods)
	//return d
}

// 释放购物车并销毁
func (o *normalOrderImpl) destroyCart() error {
	if o.cart.Release(nil) {
		return o.cart.Destroy()
	}
	return nil
}

// 根据运营商生成子订单
func (o *normalOrderImpl) createSubOrderByVendor(parentOrderId int64, buyerId int64,
	vendorId int, newOrderNo bool, items []*order.SubOrderItem) order.ISubOrder {
	orderNo := o.OrderNo()
	if newOrderNo {
		orderNo = o.manager.GetFreeOrderNo(int64(vendorId))
	}
	if len(items) == 0 {
		domain.HandleError(errors.New("拆分订单,运营商下未获取到商品,订单:"+
			orderNo), "domain")
		return nil
	}
	isp := o.shopRepo.GetShop(items[0].ShopId).(shop.IOnlineShop)
	shopName := isp.GetShopValue().ShopName
	v := &order.NormalSubOrder{
		OrderNo:  orderNo,
		BuyerId:  buyerId,
		VendorId: int64(vendorId),
		OrderId:  o.GetAggregateRootId(),
		Subject:  "子订单",
		ShopId:   items[0].ShopId,
		ShopName: shopName,
		// 总金额
		ItemAmount: 0,
		// 减免金额(包含优惠券金额)
		DiscountAmount: 0,
		ExpressFee:     0,
		FinalAmount:    0,
		// 是否挂起，如遇到无法自动进行的时挂起，来提示人工确认。
		IsSuspend:    0,
		BuyerComment: "",
		Remark:       "",
		State:        order.StatAwaitingPayment,
		UpdateTime:   o.baseValue.UpdateTime,
		Items:        items,
	}
	// 计算订单金额
	for _, iit := range items {
		//计算商品金额
		v.ItemAmount += iit.Amount
		//计算商品优惠金额
		v.DiscountAmount += iit.Amount - iit.FinalAmount
	}
	// 设置运费
	v.ExpressFee = o.vendorExpressMap[vendorId]
	// 设置包装费
	v.PackageFee = 0
	// 最终金额 = 商品金额 - 商品抵扣金额(促销折扣) + 包装费 + 快递费
	v.FinalAmount = v.ItemAmount - v.DiscountAmount +
		v.PackageFee + v.ExpressFee
	return o.repo.CreateNormalSubOrder(v)
}

//根据运营商拆单,返回拆单结果,及拆分的订单数组
func (o *normalOrderImpl) breakUpByVendor() []order.ISubOrder {
	parentOrderId := o.GetAggregateRootId()
	if parentOrderId <= 0 ||
		o.vendorItemsMap == nil ||
		len(o.vendorItemsMap) == 0 {
		//todo: 订单要取消掉
		panic(fmt.Sprintf("订单异常: 订单未生成或VendorItemMap为空,"+
			"订单编号:%d,订单号:%s,vendor len:%d",
			parentOrderId, o.OrderNo(), len(o.vendorItemsMap)))
	}
	l := len(o.vendorItemsMap)
	list := make([]order.ISubOrder, l)
	i := 0
	buyerId := o.buyer.GetAggregateRootId()
	for k, v := range o.vendorItemsMap {
		//log.Println("----- vendor ", k, len(v),l)
		list[i] = o.createSubOrderByVendor(parentOrderId, buyerId, k, l > 1, v)
		if _, err := list[i].Submit(); err != nil {
			_ = domain.HandleError(err, "domain")
		}
		i++
	}
	// 设置已拆分的订单
	o._list = list
	// 设置订单为已拆分状态
	if l > 1 {
		o.saveOrderState(order.StatBreak)
	}
	return list
}

// 扣除库存
func (o *normalOrderImpl) applyItemStock() {
	for _, v := range o.vendorItemsMap {
		for _, v2 := range v {
			o.takeGoodsStock(v2.ItemId, v2.SkuId, v2.Quantity)
		}
	}
}

//****************  订单提交结束 **************//

// 设置支付方式
//func (o *orderImpl) SetPayment(payment int) {
//	o._value.PaymentOpt = payment
//}

// GetSubOrders 获取子订单列表
func (o *normalOrderImpl) GetSubOrders() []order.ISubOrder {
	orderId := o.GetAggregateRootId()
	if orderId <= 0 {
		panic(order.ErrNoYetCreated)
	}
	if o._list == nil {
		list := o.orderRepo.GetNormalSubOrders(orderId)
		for _, v := range list {
			sub := o.repo.CreateNormalSubOrder(v)
			o._list = append(o._list, sub)
		}
	}
	return o._list
}

// 在线支付交易完成
func (o *normalOrderImpl) OnlinePaymentTradeFinish() (err error) {
	for _, o := range o.GetSubOrders() {
		if err = o.PaymentFinishByOnlineTrade(); err != nil {
			return err
		}
	}
	return nil
}

// 扣减商品库存
func (o *normalOrderImpl) takeGoodsStock(itemId, skuId int64, quantity int32) error {
	gds := o.goodsRepo.GetItem(itemId)
	if gds == nil {
		return item.ErrNoSuchItem
	}
	return gds.TakeStock(skuId, quantity)
}

// 更新返现到会员账户
func (o *normalOrderImpl) updateShoppingMemberBackFee(mch merchant.IMerchant,
	m member.IMember, fee int64, unixTime int64) error {
	if fee == 0 {
		return nil
	}
	pv := mch.GetValue()

	//更新账户
	acc := m.GetAccount()
	acv := acc.GetValue()
	//acc.TotalAmount += o._value.Fee
	//acc.TotalPay += o._value.PayFee
	acv.WalletBalance += fee // 更新赠送余额
	acv.TotalWalletAmount += fee
	acv.UpdateTime = unixTime
	_, err := acc.Save()
	if err == nil {
		orderNo := o.OrderNo()
		//给自己返现
		tit := fmt.Sprintf("订单:%s(商户:%s)返现￥%.2f元", orderNo, pv.Name, fee)
		_, err = acc.CarryTo(member.AccountWallet,
			member.AccountOperateData{
				Title:   tit,
				Amount:  int(fee * 100),
				OuterNo: orderNo,
				Remark:  "sys",
			}, false, 0)
	}
	return err
}

// 处理返现促销
func (o *normalOrderImpl) handleCashBackPromotions(pt merchant.IMerchant,
	m member.IMember) error {
	proms := o.GetPromotionBinds()
	for _, v := range proms {
		if v.PromotionType == promotion.TypeFlagCashBack {
			c := o.promRepo.GetPromotion(v.PromotionId)
			return o.handleCashBackPromotion(pt, m, v, c)
		}
	}
	return nil
}

// 处理返现促销
func (o *normalOrderImpl) handleCashBackPromotion(pt merchant.IMerchant,
	m member.IMember,
	v *order.OrderPromotionBind, pm promotion.IPromotion) error {
	cpv := pm.GetRelationValue().(*promotion.ValueCashBack)

	//更新账户
	bFee := cpv.BackFee
	acc := m.GetAccount()
	acv := acc.GetValue()
	acv.WalletBalance += bFee // 更新赠送余额
	acv.TotalWalletAmount += bFee
	// 赠送金额，不应该计入到余额，可采取充值到余额
	//acc.Balance += float32(cpv.BackFee)                            // 更新账户余额

	acv.UpdateTime = time.Now().Unix()
	_, err := acc.Save()

	if err == nil {
		orderNo := o.OrderNo()
		// 优惠绑定生效
		v.IsApply = 1
		o.orderRepo.SavePromotionBindForOrder(v)

		// 处理自定义返现
		c := pm.(promotion.ICashBackPromotion)
		HandleCashBackDataTag(m, o, c, o.memberRepo)

		//给自己返现
		tit := fmt.Sprintf("返现￥%d元,订单号:%s", cpv.BackFee, orderNo)
		_, err = acc.CarryTo(member.AccountWallet, member.AccountOperateData{
			Title:   tit,
			Amount:  int(cpv.BackFee * 100),
			OuterNo: orderNo,
			Remark:  "sys",
		}, false, 0)
	}
	return err
}

//todo: ?? 自动收货功能

var _ order.ISubOrder = new(subOrderImpl)

// 子订单实现
type subOrderImpl struct {
	value           *order.NormalSubOrder
	parent          order.IOrder
	buyer           member.IMember
	internalSuspend bool //内部挂起
	paymentRepo     payment.IPaymentRepo
	repo            order.IOrderRepo
	memberRepo      member.IMemberRepo
	itemRepo        item.IItemRepo
	productRepo     product.IProductRepo
	manager         order.IOrderManager
	shipRepo        shipment.IShipmentRepo
	valRepo         valueobject.IValueRepo
	mchRepo         merchant.IMerchantRepo
	registryRepo    registry.IRegistryRepo
}

func NewSubNormalOrder(v *order.NormalSubOrder,
	manager order.IOrderManager, rep order.IOrderRepo,
	mmRepo member.IMemberRepo, goodsRepo item.IItemRepo,
	shipRepo shipment.IShipmentRepo, productRepo product.IProductRepo,
	paymentRepo payment.IPaymentRepo, valRepo valueobject.IValueRepo,
	mchRepo merchant.IMerchantRepo, registryRepo registry.IRegistryRepo) order.ISubOrder {
	return &subOrderImpl{
		value:        v,
		manager:      manager,
		repo:         rep,
		memberRepo:   mmRepo,
		itemRepo:     goodsRepo,
		productRepo:  productRepo,
		shipRepo:     shipRepo,
		paymentRepo:  paymentRepo,
		valRepo:      valRepo,
		mchRepo:      mchRepo,
		registryRepo: registryRepo,
	}
}

// GetDomainId 获取领域对象编号
func (o *subOrderImpl) GetDomainId() int64 {
	return o.value.Id
}

// GetValue 获取值对象
func (o *subOrderImpl) GetValue() *order.NormalSubOrder {
	return o.value
}

// Complex 复合的订单信息
func (o *subOrderImpl) Complex() *order.ComplexOrder {
	v := o.GetValue()
	co := o.baseOrder().Complex()
	co.OrderId = o.GetDomainId()
	co.ParentOrderId = v.OrderId
	co.VendorId = v.VendorId
	co.ShopId = v.ShopId
	co.SubOrder = true
	co.OrderNo = o.value.OrderNo
	co.Subject = v.Subject
	co.DiscountAmount = v.DiscountAmount
	co.ItemAmount = v.ItemAmount
	co.ExpressFee = v.ExpressFee
	co.PackageFee = v.PackageFee
	co.FinalAmount = v.FinalAmount
	co.UpdateTime = v.UpdateTime
	co.State = v.State
	co.Items = []*order.ComplexItem{}
	for _, v := range o.Items() {
		co.Items = append(co.Items, o.parseComplexItem(v))
	}
	return co
}

// 转换订单商品
func (o *subOrderImpl) parseComplexItem(i *order.SubOrderItem) *order.ComplexItem {
	snap := o.itemRepo.GetSalesSnapshot(i.SnapshotId)
	it := &order.ComplexItem{
		ID:             i.ID,
		ItemId:         i.ItemId,
		SkuId:          i.SkuId,
		SkuWord:        snap.Sku,
		SnapshotId:     i.SnapshotId,
		ItemTitle:      snap.GoodsTitle,
		MainImage:      snap.Image,
		Price:          snap.Price,
		FinalPrice:     snap.Price,
		Quantity:       i.Quantity,
		ReturnQuantity: i.ReturnQuantity,
		Amount:         i.Amount,
		FinalAmount:    i.FinalAmount,
		IsShipped:      i.IsShipped,
		Data:           make(map[string]string),
	}
	base := o.baseOrder().(*normalOrderImpl)
	base.baseOrderImpl.bindItemInfo(it)
	return it
}

// 获取商品项
func (o *subOrderImpl) Items() []*order.SubOrderItem {
	if (o.value.Items == nil || len(o.value.Items) == 0) &&
		o.GetDomainId() > 0 {
		o.value.Items = o.repo.GetSubOrderItems(o.GetDomainId())
	}
	return o.value.Items
}

// 获取订单
func (o *subOrderImpl) baseOrder() order.IOrder {
	if o.parent == nil {
		o.parent = o.manager.GetOrderById(o.value.OrderId)
	}
	return o.parent
}

// 获取购买的会员
func (o *subOrderImpl) getBuyer() member.IMember {
	return o.baseOrder().Buyer()
}

// 添加备注
func (o *subOrderImpl) AddRemark(remark string) {
	o.value.Remark = remark
}

func (o *subOrderImpl) saveOrderItems() error {
	unix := time.Now().Unix()
	id := o.GetDomainId()
	for _, v := range o.Items() {
		v.OrderId = id
		v.UpdateTime = unix
		_, err := o.repo.SaveOrderItem(id, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// 提交子订单
func (o *subOrderImpl) Submit() (int64, error) {
	if o.GetDomainId() > 0 {
		panic("suborder is created!")
	}
	if o.value.CreateTime <= 0 {
		unix := time.Now().Unix()
		o.value.CreateTime = unix
		o.value.UpdateTime = unix
	}
	id, err := util.I64Err(o.repo.SaveSubOrder(o.value))
	if err == nil {
		o.value.Id = id
		err = o.saveOrderItems()
		o.AppendLog(order.LogSetup, true, "{created}")
	}
	return id, err
}

// 保存订单
func (o *subOrderImpl) saveSubOrder() error {
	unix := time.Now().Unix()
	o.value.UpdateTime = unix
	if o.GetDomainId() <= 0 {
		panic("please use Submit() to create new suborder!")
	}
	_, err := o.repo.SaveSubOrder(o.value)
	if err == nil {
		o.syncOrderState()
	}
	return err
}

// 同步订单状态
func (o *subOrderImpl) syncOrderState() {
	if bo := o.baseOrder(); bo != nil {
		oi := bo.(*normalOrderImpl).baseOrderImpl
		if oi.State() != order.StatBreak {
			oi.saveOrderState(order.OrderState(o.value.State))
		}
	}

}

// 订单完成支付
func (o *subOrderImpl) orderFinishPaid() error {
	if o.value.IsPaid == 1 {
		return order.ErrOrderPayed
	}
	if o.value.State == order.StatAwaitingPayment {
		o.value.IsPaid = 1
		o.value.State = order.StatAwaitingConfirm
		err := o.AppendLog(order.LogSetup, true, "{finish_pay}")
		if err == nil {
			err = o.saveSubOrder()
		}
		return err
	}
	return order.ErrUnusualOrderStat
}

// 在线支付交易完成
func (o *subOrderImpl) PaymentFinishByOnlineTrade() error {
	return o.orderFinishPaid()
}

// 挂起
func (o *subOrderImpl) Suspend(reason string) error {
	o.value.IsSuspend = 1
	o.internalSuspend = true
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveSubOrder()
	if err == nil {
		err = o.AppendLog(order.LogSetup, true, "订单已锁定"+reason)
	}
	return err
}

// 添加日志
func (o *subOrderImpl) AppendLog(logType order.LogType, system bool, message string) error {
	if o.GetDomainId() <= 0 {
		return errors.New("order not created.")
	}
	var systemInt int
	if system {
		systemInt = 1
	} else {
		systemInt = 0
	}
	l := &order.OrderLog{
		OrderId:    o.GetDomainId(),
		Type:       int(logType),
		IsSystem:   systemInt,
		OrderState: int(o.value.State),
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return o.repo.SaveNormalSubOrderLog(l)
}

// 确认订单
func (o *subOrderImpl) Confirm() (err error) {
	//todo: 线下交易,自动确认
	//if o._value.PaymentOpt == enum.PaymentOnlinePay &&
	//o._value.IsPaid == enum.FALSE {
	//    return order.ErrOrderNotPayed
	//}
	if o.value.State < order.StatAwaitingConfirm {
		return order.ErrOrderNotPayed
	}
	if o.value.State >= order.StatAwaitingPickup {
		return order.ErrOrderHasConfirm
	}
	o.value.State = order.StatAwaitingPickup
	o.value.UpdateTime = time.Now().Unix()
	err = o.saveSubOrder()
	if err == nil {
		go o.addItemSalesNum()
		err = o.AppendLog(order.LogSetup, false, "{confirm}")
	}
	return err
}

// 增加商品的销售数量
func (o *subOrderImpl) addItemSalesNum() {
	//log.Println("---订单：",o.value.OrderNo," 商品：",len(o.Items()))
	for _, v := range o.Items() {
		it := o.itemRepo.GetItem(v.ItemId)
		err := it.AddSalesNum(v.SkuId, v.Quantity)
		if err != nil {
			log.Println("---增加销售数量：", v.ItemId,
				" sku:", v.SkuId, " error:", err.Error())
		}
	}
}

// PickUp 捡货(备货)
func (o *subOrderImpl) PickUp() error {
	if o.value.State < order.StatAwaitingPickup {
		return order.ErrOrderNotConfirm
	}
	if o.value.State >= order.StatAwaitingShipment {
		return order.ErrOrderHasPickUp
	}
	o.value.State = order.StatAwaitingShipment
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveSubOrder()
	if err == nil {
		err = o.AppendLog(order.LogSetup, true, "{pickup}")
	}
	return err
}

// Ship 发货
func (o *subOrderImpl) Ship(spId int32, spOrder string) error {
	//so := o._shipRepo.GetOrders()
	//todo: 可进行发货修改
	if o.value.State >= order.StatShipped {
		return order.ErrOrderShipped
	}
	// 如果没有备货完成,则发货前自动完成备货
	if o.value.State < order.StatAwaitingShipment {
		o.value.State = order.StatAwaitingShipment
		o.value.UpdateTime = time.Now().Unix()
		_ = o.AppendLog(order.LogSetup, true, "{pickup}")
		//return order.ErrOrderNotPickUp
	}
	if list := o.shipRepo.GetShipOrders(o.GetDomainId(), true); len(list) > 0 {
		return order.ErrPartialShipment
	}
	if spId <= 0 || spOrder == "" {
		return shipment.ErrMissingSpInfo
	}

	so := o.createShipmentOrder(o.Items())
	if so == nil {
		return order.ErrUnusualOrder
	}
	// 生成发货单并发货
	err := so.Ship(spId, spOrder)
	if err == nil {
		o.value.State = order.StatShipped
		o.value.UpdateTime = time.Now().Unix()
		err = o.saveSubOrder()
		if err == nil {
			// 保存商品的发货状态
			err = o.saveOrderItems()
			_ = o.AppendLog(order.LogSetup, true, "{shipped}")
		}
	}
	return err
}

func (o *subOrderImpl) createShipmentOrder(items []*order.SubOrderItem) shipment.IShipmentOrder {
	if items == nil || len(items) == 0 {
		return nil
	}
	unix := time.Now().Unix()
	orderId := o.baseOrder().GetAggregateRootId()
	subOrderId := o.GetDomainId()
	so := &shipment.ShipmentOrder{
		ID:          0,
		OrderId:     orderId,
		SubOrderId:  subOrderId,
		ShipmentLog: "",
		ShipTime:    unix,
		State:       shipment.StatAwaitingShipment,
		UpdateTime:  unix,
		Items:       []*shipment.ShipmentItem{},
	}
	for _, v := range items {
		if v.IsShipped == 1 {
			continue
		}
		so.Amount += float64(v.Amount)
		so.FinalAmount += float64(v.FinalAmount)
		so.Items = append(so.Items, &shipment.ShipmentItem{
			ID:          0,
			SnapshotId:  v.SnapshotId,
			Quantity:    v.Quantity,
			Amount:      float64(v.Amount),
			FinalAmount: float64(v.FinalAmount),
		})
		v.IsShipped = 1
	}
	return o.shipRepo.CreateShipmentOrder(so)
}

// 已收货
func (o *subOrderImpl) BuyerReceived() error {
	if o.value.State < order.StatShipped {
		return order.ErrOrderNotShipped
	}
	if o.value.State >= order.StatCompleted {
		return order.ErrIsCompleted
	}
	dt := time.Now()
	o.value.State = order.StatCompleted
	o.value.UpdateTime = dt.Unix()
	o.value.IsSuspend = 0
	err := o.saveSubOrder()
	if err == nil {
		err = o.AppendLog(order.LogSetup, true, "{completed}")
		if err == nil {
			go o.vendorSettle()
			// 执行其他的操作
			if err2 := o.onOrderComplete(); err != nil {
				domain.HandleError(err2, "domain")
			}
		}
	}
	return err
}

func (o *subOrderImpl) getOrderAmount() (amount int64, refund int64) {
	items := o.Items()
	for _, item := range items {
		if item.ReturnQuantity > 0 {
			a := item.Amount / int64(item.Quantity) * int64(item.ReturnQuantity)
			if item.ReturnQuantity != item.Quantity {
				amount += item.Amount - a
			}
			refund += a
		} else {
			amount += item.Amount
		}
	}
	//如果非全部退货、退款,则加上运费及包装费
	if amount > 0 {
		amount += o.value.ExpressFee + o.value.PackageFee
	}
	return amount, refund
}

// 获取订单的成本
func (o *subOrderImpl) getOrderCost() int64 {
	var cost int64
	items := o.Items()
	for _, item := range items {
		snap := o.itemRepo.GetSalesSnapshot(item.SnapshotId)
		cost += snap.Cost * int64(item.Quantity-item.ReturnQuantity)
	}
	//如果非全部退货、退款,则加上运费及包装费
	if cost > 0 {
		cost += o.value.ExpressFee + o.value.PackageFee
	}
	return cost
}

// 商户结算
func (o *subOrderImpl) vendorSettle() error {
	vendor := o.mchRepo.GetMerchant(int(o.value.VendorId))
	if vendor != nil {
		settleMode := o.registryRepo.Get(registry.MchOrderSettleMode).IntValue()
		switch enum.MchSettleMode(settleMode) {
		case enum.MchModeSettleByCost:
			return o.vendorSettleByCost(vendor)
		case enum.MchModeSettleByRate:
			return o.vendorSettleByRate(vendor)
		case enum.MchModeSettleByOrderQuantity:
			return o.vendorSettleByOrderQuantity(vendor)
		}
	}
	return nil
}

// 根据供货价进行商户结算
func (o *subOrderImpl) vendorSettleByCost(vendor merchant.IMerchant) error {
	_, refund := o.getOrderAmount()
	sAmount := o.getOrderCost()
	if sAmount > 0 {
		totalAmount := int(float32(sAmount) * float32(enum.RATE_AMOUNT))
		refundAmount := int(float32(refund) * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKNormalOrder, totalAmount)
		return vendor.Account().SettleOrder(o.value.OrderNo,
			totalAmount, tradeFee, refundAmount, "零售订单结算")
	}
	return nil
}

// 根据比例进行商户结算
func (o *subOrderImpl) vendorSettleByRate(vendor merchant.IMerchant) error {
	rate := o.registryRepo.Get(registry.MchOrderSettleRate).FloatValue()
	amount, refund := o.getOrderAmount()
	sAmount := int64(float64(amount) * rate)
	if sAmount > 0 {
		totalAmount := int(float32(sAmount) * float32(enum.RATE_AMOUNT))
		refundAmount := int(float32(refund) * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKNormalOrder, totalAmount)
		return vendor.Account().SettleOrder(o.value.OrderNo,
			totalAmount, tradeFee, refundAmount, "零售订单结算")

	}
	return nil
}
func (o *subOrderImpl) vendorSettleByOrderQuantity(vendor merchant.IMerchant) error {
	fee := o.registryRepo.Get(registry.MchSingleOrderServiceFee).FloatValue()
	amount, refund := o.getOrderAmount()
	if fee > 0 {
		totalAmount := int(math.Min(float64(amount), fee) * float64(enum.RATE_AMOUNT))
		refundAmount := int(float32(refund) * float32(enum.RATE_AMOUNT))
		tradeFee, _ := vendor.SaleManager().MathTradeFee(
			merchant.TKNormalOrder, totalAmount)
		return vendor.Account().SettleOrder(o.value.OrderNo,
			totalAmount, tradeFee, refundAmount, "零售订单结算")

	}
	return nil
}

// 获取订单的日志
func (o *subOrderImpl) LogBytes() []byte {
	buf := bytes.NewBufferString("")
	list := o.repo.GetSubOrderLogs(o.GetDomainId())
	for _, v := range list {
		buf.WriteString(time.Unix(v.RecordTime, 0).Format("2006-01-02 15:04:05"))
		buf.WriteString("  ")
		if v.Message[:1] == "{" {
			if msg := o.getLogStringByStat(v.OrderState); len(msg) > 0 {
				v.Message = msg
			}
		}
		buf.WriteString(v.Message)
		buf.Write([]byte("\n"))
	}
	return buf.Bytes()
}

func (o *subOrderImpl) getLogStringByStat(stat int) string {
	switch stat {
	case order.StatAwaitingPayment:
		return "订单已提交..."
	case order.StatAwaitingConfirm:
		return "订单已支付,等待商户确认。"
	case order.StatAwaitingPickup:
		return "订单已确认,备货中..."
	case order.StatAwaitingShipment:
		return "备货完成,即将发货。"
	case order.StatShipped:
		return "订单已发货,请等待收货。"
	case order.StatCompleted:
		return "已收货,订单完成。"
	}
	return ""
}

// 更新账户
func (o *subOrderImpl) updateAccountForOrder(m member.IMember) error {
	if o.value.State != order.StatCompleted {
		return order.ErrUnusualOrderStat
	}
	var err error
	ov := o.value
	amount := ov.FinalAmount
	acc := m.GetAccount()

	// 增加经验
	expEnabled := o.registryRepo.Get(registry.ExperienceEnabled).BoolValue()
	if expEnabled {
		rate := o.registryRepo.Get(registry.ExperienceRateByOrder).FloatValue()
		if exp := int(float64(amount) * rate); exp > 0 {
			if err = m.AddExp(exp); err != nil {
				return err
			}
		}
	}

	// 增加积分
	//todo: 增加阶梯的返积分,比如订单满30送100积分, 不考虑额外赠送,额外的当做补贴
	rate := o.registryRepo.Get(registry.IntegralRateByWholesaleOrder).FloatValue()
	integral := int(float64(amount) * rate)
	// 赠送积分
	if integral > 0 {
		_, err = m.GetAccount().CarryTo(member.AccountIntegral,
			member.AccountOperateData{
				Title:   "购物消费赠送积分",
				Amount:  integral,
				OuterNo: o.value.OrderNo,
				Remark:  "sys",
			}, false, 0)
		if err != nil {
			return err
		}
	}
	acv := acc.GetValue()
	acv.TotalExpense += ov.ItemAmount
	acv.TotalPay += ov.FinalAmount
	acv.UpdateTime = time.Now().Unix()
	_, err = acc.Save()
	return err
}

// 取消订单
func (o *subOrderImpl) Cancel(reason string) error {
	if o.value.State == order.StatCancelled {
		return order.ErrOrderCancelled
	}
	// 已发货订单无法取消
	if o.value.State >= order.StatShipped {
		return order.ErrOrderShippedCancel
	}

	o.value.State = order.StatCancelled
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveSubOrder()
	if err == nil {
		domain.HandleError(o.AppendLog(order.LogSetup, true, reason), "domain")
		// 取消支付单
		err = o.cancelPaymentOrder()
		if err == nil {
			// 取消商品
			err = o.cancelGoods()
		}
	}
	return err
}

// 取消商品
func (o *subOrderImpl) cancelGoods() error {
	for _, v := range o.Items() {
		snapshot := o.itemRepo.GetSalesSnapshot(v.SnapshotId)
		if snapshot == nil {
			return item.ErrNoSuchSnapshot
		}
		gds := o.itemRepo.GetItem(snapshot.SkuId)
		if gds != nil {
			// 释放库存
			gds.ReleaseStock(v.SkuId, v.Quantity)
			// 如果订单已付款，则取消销售数量
			if o.value.IsPaid == 1 {
				gds.CancelSale(v.SkuId, v.Quantity, o.value.OrderNo)
			}
		}
	}
	return nil
}

// 取消支付单
func (o *subOrderImpl) cancelPaymentOrder() error {
	od := o.baseOrder()
	if od.Type() != order.TRetail {
		panic("not support order type")
	}
	ip := od.GetPaymentOrder()
	if ip != nil {
		return ip.Cancel() //todo: there have a bug, when other order has shipmented. all sub order will be cancelled.
	}
	return nil
}

// 退回商品
func (o *subOrderImpl) Return(snapshotId int64, quantity int32) error {
	for _, v := range o.Items() {
		if v.SnapshotId == snapshotId {
			if v.Quantity-v.ReturnQuantity < quantity {
				return order.ErrOutOfQuantity
			}
			v.ReturnQuantity += quantity
			_, err := o.repo.SaveOrderItem(o.GetDomainId(), v)
			return err
		}
	}
	return order.ErrNoSuchGoodsOfOrder
}

// 撤销退回商品
func (o *subOrderImpl) RevertReturn(snapshotId int64, quantity int32) error {
	for _, v := range o.Items() {
		if v.SnapshotId == snapshotId {
			if v.ReturnQuantity < quantity {
				return order.ErrOutOfQuantity
			}
			v.ReturnQuantity -= quantity
			_, err := o.repo.SaveOrderItem(o.GetDomainId(), v)
			return err
		}
	}
	return order.ErrNoSuchGoodsOfOrder
}

// 申请退款
func (o *subOrderImpl) SubmitRefund(reason string) error {
	if o.value.State == order.StatAwaitingPayment {
		return o.Cancel("订单主动申请取消,原因:" + reason)
	}
	if o.value.State >= order.StatShipped ||
		o.value.State >= order.StatCancelled {
		return order.ErrOrderCancelled
	}
	o.value.State = order.StatAwaitingCancel
	o.value.UpdateTime = time.Now().Unix()
	return o.saveSubOrder()
}

// 谢绝订单
func (o *subOrderImpl) Decline(reason string) error {
	if o.value.State == order.StatAwaitingPayment {
		return o.Cancel("商户取消,原因:" + reason)
	}
	if o.value.State >= order.StatShipped ||
		o.value.State >= order.StatCancelled {
		return order.ErrOrderCancelled
	}
	o.value.State = order.StatDeclined
	o.value.UpdateTime = time.Now().Unix()
	return o.saveSubOrder()
}

// 退款 todo: will delete,代码供取消订单参考
func (o *subOrderImpl) refund() error {
	// 已退款
	if o.value.State == order.StatRefunded ||
		o.value.State == order.StatCancelled {
		return order.ErrHasRefund
	}
	// 不允许退款
	if o.value.State != order.StatAwaitingCancel &&
		o.value.State != order.StatDeclined {
		return order.ErrDisallowRefund
	}
	o.value.State = order.StatRefunded
	o.value.UpdateTime = time.Now().Unix()
	err := o.saveSubOrder()
	if err == nil {
		err = o.cancelPaymentOrder()
	}
	return err
}

// 取消退款申请
func (o *subOrderImpl) CancelRefund() error {
	if o.value.State != order.StatAwaitingCancel || o.value.IsPaid == 0 {
		panic(errors.New("订单已经取消,不允许再退款"))
	}
	o.value.State = order.StatAwaitingConfirm
	o.value.UpdateTime = time.Now().Unix()
	return o.saveSubOrder()
}

// 完成订单
func (o *subOrderImpl) onOrderComplete() error {
	// 更新发货单
	soList := o.shipRepo.GetShipOrders(o.GetDomainId(), true)
	for _, v := range soList {
		domain.HandleError(v.Completed(), "domain")
	}

	// 获取消费者消息
	m := o.getBuyer()
	if m == nil {
		return member.ErrNoSuchMember
	}

	// 更新会员账户
	err := o.updateAccountForOrder(m)
	if err != nil {
		return err
	}

	// 处理返现
	err = o.handleCashBack()

	return err
}

// 更新返现到会员账户
func (o *subOrderImpl) updateShoppingMemberBackFee(mchName string,
	m member.IMember, fee int64, unixTime int64) error {
	if fee <= 0 || math.IsNaN(float64(fee)) {
		return nil
	}
	v := o.GetValue()

	//更新账户
	acc := m.GetAccount()
	//给自己返现
	tit := fmt.Sprintf("订单:%s(商户:%s)返现￥%.2f元", v.OrderNo, mchName, fee)
	_, err := acc.CarryTo(member.AccountWallet,
		member.AccountOperateData{
			Title:   tit,
			Amount:  int(fee * 100),
			OuterNo: o.value.OrderNo,
			Remark:  "sys",
		}, false, 0)
	return err
}
