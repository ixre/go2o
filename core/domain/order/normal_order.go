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
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/cart"
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
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/types/typeconv"
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
	_subOrders      []order.ISubOrder
	_payOrder       payment.IPaymentOrder
	// 返利推荐人
	_AffiliateMember member.IMember
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

// Complex 复合的订单信息
func (o *normalOrderImpl) Complex() *order.ComplexOrder {
	co := o.baseOrderImpl.Complex()
	if o.GetAggregateRootId() > 0 {
		subOrders := o.GetSubOrders()
		for _, v := range subOrders {
			if v.GetValue().BreakStatus == order.BreakDefault {
				continue
			}
			co.Details = append(co.Details, parseDetailValue(v))
		}
	}
	return co
}

// ApplyTraderCode 使用返利人代码
func (o *normalOrderImpl) ApplyTraderCode(code string) error {
	memberId := o.memberRepo.GetMemberIdByCode(code)
	if memberId <= 0 {
		return member.ErrNoSuchMember
	}
	im := o.memberRepo.GetMember(memberId)
	// 用户没有返利标志，则不作任何处理
	if im == nil || im.ContainFlag(member.FlagAffiliateDisabled) {
		return nil
	}
	o._AffiliateMember = im
	return nil
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
	// 设置订单标题
	if len(o.baseValue.Subject) == 0 {
		// note: 如果商品标题有空格,在pgsql下引发错误：
		// invalid byte sequence for encoding "UTF8": 0x
		title := strings.ReplaceAll(snap.GoodsTitle, " ", "")
		o.baseValue.Subject = title
		if len(title) > 16 {
			o.baseValue.Subject = title[:15] + "..."
		}
	}

	fee := c.Sku.Price * int64(c.Quantity)
	return &order.SubOrderItem{
		ID:         0,
		VendorId:   c.VendorId,
		ShopId:     c.ShopId,
		SkuId:      c.SkuId,
		ItemId:     c.ItemId,
		SnapshotId: snap.Id,
		// 退回数量
		Quantity:       c.Quantity,
		ReturnQuantity: 0,
		Amount:         fee,
		FinalAmount:    fee,
		//是否配送
		IsShipped:    0,
		ExpressTplId: c.Sku.ExpressTid,
		Weight:       c.Sku.Weight * c.Quantity, //计算重量
		Bulk:         c.Sku.Bulk * c.Quantity,   //计算体积
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
	itemAmount := v.ItemAmount + v.ExpressFee + v.PackageFee
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
		TotalAmount:    finalAmount,
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
			_ = it.Cancel(false, "下单错误自动取消")
		}
	}
	return err
}

// ChangeShipmentAddress implements order.IOrder
func (o *normalOrderImpl) ChangeShipmentAddress(addressId int64) error {
	return o.baseOrderImpl.ChangeShipmentAddress(addressId)
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
	breakStatus := order.BreakNoBreak // 如果只有一个子订单,则不需要拆分
	if newOrderNo {
		orderNo = o.manager.GetFreeOrderNo(int64(vendorId))
		breakStatus = order.BreakAwaitBreak // 待拆分
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
		IsForbidden:    0,
		BuyerComment:   "",
		Remark:         "",
		Status:         order.StatAwaitingPayment,
		BreakStatus:    breakStatus,
		UpdateTime:     o.baseValue.UpdateTime,
		Items:          items,
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
	so := o.repo.CreateNormalSubOrder(v)
	o.createAffiliateRebateOrder(so)
	return so
}

// 创建返利订单
func (o *normalOrderImpl) createAffiliateRebateOrder(so order.ISubOrder) {
	if o._AffiliateMember != nil {
		// 未开启返利
		rv, _ := o.registryRepo.GetValue(registry.OrderEnableAffiliateRebate)
		if v, _ := strconv.ParseBool(rv); !v {
			return
		}
		// 获取返利比例
		rv, err := o.registryRepo.GetValue(registry.OrderGlobalAffiliateRebateRate)
		if err != nil {
			log.Println("[ warning]: affiliate rebate rate error", err.Error())
			return
		}
		rate := typeconv.MustFloat(rv)
		if rate <= 0 {
			return
		}
		ov := so.GetValue()
		unix := time.Now().Unix()
		v := &order.AffiliateDistribution{
			Id:               0,
			PlanId:           0,
			BuyerId:          o._AffiliateMember.GetAggregateRootId(),
			OwnerId:          o._AffiliateMember.GetAggregateRootId(),
			Flag:             1,
			IsRead:           0,
			AffiliateCode:    o._AffiliateMember.GetValue().UserCode,
			OrderNo:          ov.OrderNo,
			OrderSubject:     ov.Subject,
			OrderAmount:      ov.FinalAmount,
			DistributeAmount: int64(float64(ov.FinalAmount) * rate),
			Status:           1,
			CreateTime:       unix,
			UpdateTime:       unix,
		}
		o.orderRepo.SaveOrderRebate(v)
	}
}

// 根据运营商拆单,返回拆单结果,及拆分的订单数组
// 如果拆单,需要生成一个用于支付的子订单,支付完成后删除
func (o *normalOrderImpl) breakUpByVendor() ([]order.ISubOrder, error) {
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
	// 生成一个用于支付的子订单
	orderId := 0
	if l > 1 {
		iso, err := o.createPaymentSubOrder()
		if err != nil {
			log.Println("生成子订单失败:" + err.Error())
			return nil, err
		}
		list = append(list, iso)
		orderId = int(iso.GetDomainId())
	}

	buyerId := o.buyer.GetAggregateRootId()
	for k, v := range o.vendorItemsMap {
		// 绑定商品项的订单编号到支付单
		for _, it := range v {
			it.OrderId = int64(orderId)
		}
		// log.Println("----- vendor ", k, len(v),l)
		list[i] = o.createSubOrderByVendor(parentOrderId, buyerId, k, l > 1, v)
		if _, err := list[i].Submit(); err != nil {
			_ = domain.HandleError(err, "domain")
		}
		i++
	}
	// 设置已拆分的订单
	o._subOrders = list
	if l > 1 {
		// 设置订单为已拆分状态
		o.saveOrderState(order.StatBreak)
	}
	return list, nil
}

// createPaymentSubOrder 生成一个用于合并支付的子订单
func (o *normalOrderImpl) createPaymentSubOrder() (order.ISubOrder, error) {
	orderNo := o.OrderNo()
	breakStatus := order.BreakDefault
	vo := o.baseValue
	v := &order.NormalSubOrder{
		OrderNo:  orderNo,
		BuyerId:  o.baseValue.BuyerId,
		VendorId: 0,
		OrderId:  o.GetAggregateRootId(),
		Subject:  "支付子订单",
		ShopId:   0,
		ShopName: "",
		// 总金额
		ItemAmount: vo.ItemAmount,
		// 减免金额(包含优惠券金额)
		DiscountAmount: vo.DiscountAmount,
		ExpressFee:     vo.ExpressFee,
		PackageFee:     vo.PackageFee,
		FinalAmount:    vo.FinalAmount,
		BuyerComment:   "",
		Remark:         "",
		Status:         order.StatAwaitingPayment,
		BreakStatus:    breakStatus,
		UpdateTime:     o.baseValue.UpdateTime,
	}
	isp := o.repo.CreateNormalSubOrder(v)
	_, err := isp.Submit()
	return isp, err
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
	if o._subOrders == nil {
		list := o.orderRepo.GetNormalSubOrders(orderId)
		for _, v := range list {
			sub := o.repo.CreateNormalSubOrder(v)
			o._subOrders = append(o._subOrders, sub)
		}
	}
	return o._subOrders
}

// 在线支付交易完成
func (o *normalOrderImpl) OnlinePaymentTradeFinish() (err error) {
	// 排除支付子订单
	for i, so := range o.GetSubOrders() {
		// 销毁支付子订单
		ov := so.GetValue()
		if ov.BreakStatus == order.BreakDefault {
			if err := so.Destory(); err != nil {
				log.Println("销毁支付子订单失败:" + err.Error())
			}
			o._subOrders = append(o._subOrders[:i], o._subOrders[i+1:]...)
		}
	}
	for _, so := range o.GetSubOrders() {
		if err = so.PaymentFinishByOnlineTrade(); err != nil {
			return err
		}
	}
	o.baseValue.IsPaid = 1
	o.baseOrderImpl.saveOrder()
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
