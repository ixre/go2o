/**
 * Copyright 2014 @ z3q.net.
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
	"github.com/jsix/gof/log"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"math"
	"strconv"
	"strings"
	"time"
)

var _ order.IOrder = new(orderImpl)

//todo: 促销

type orderImpl struct {
	_manager order.IOrderManager
	_value   *order.Order
	_cart    cart.ICart //购物车,仅在订单生成时设置
	//_payOrder        *payment.IPaymentOrder //支付单
	_paymentOrder    payment.IPaymentOrder
	_coupons         []promotion.ICouponPromotion
	_availPromotions []promotion.IPromotion
	_orderPbs        []*order.OrderPromotionBind
	_memberRep       member.IMemberRep
	_buyer           member.IMember
	_orderRep        order.IOrderRep
	_partnerRep      merchant.IMerchantRep //todo: can delete ?
	_expressRep      express.IExpressRep
	_payRep          payment.IPaymentRep
	_goodsRep        goods.IGoodsRep
	_saleRep         sale.ISaleRep
	_promRep         promotion.IPromotionRep
	_valRep          valueobject.IValueRep
	// 运营商商品映射,用于整理购物车
	_vendorItemsMap map[int][]*order.OrderItem
	// 运营商与邮费的MAP
	_vendorExpressMap map[int]float32
	// 是否为内部挂起
	_internalSuspend bool
	_subList         []order.ISubOrder
}

func newOrder(shopping order.IOrderManager, value *order.Order,
	mchRep merchant.IMerchantRep, shoppingRep order.IOrderRep,
	goodsRep goods.IGoodsRep, saleRep sale.ISaleRep,
	promRep promotion.IPromotionRep, memberRep member.IMemberRep,
	expressRep express.IExpressRep, payRep payment.IPaymentRep,
	valRep valueobject.IValueRep) order.IOrder {
	return &orderImpl{
		_manager:    shopping,
		_value:      value,
		_memberRep:  memberRep,
		_promRep:    promRep,
		_orderRep:   shoppingRep,
		_partnerRep: mchRep,
		_goodsRep:   goodsRep,
		_saleRep:    saleRep,
		_valRep:     valRep,
		_expressRep: expressRep,
		_payRep:     payRep,
	}
}

func (o *orderImpl) GetAggregateRootId() int {
	return o._value.Id
}

func (o *orderImpl) GetValue() *order.Order {
	return o._value
}

// 设置订单值
func (o *orderImpl) SetValue(v *order.Order) error {
	v.Id = o.GetAggregateRootId()
	o._value = v
	return nil
}

// 应用优惠券
func (o *orderImpl) ApplyCoupon(coupon promotion.ICouponPromotion) error {
	//if o._coupons == nil {
	//	o._coupons = []promotion.ICouponPromotion{}
	//}
	//o._coupons = append(o._coupons, coupon)

	// 添加到促销信息中
	if o._orderPbs == nil {
		o._orderPbs = []*order.OrderPromotionBind{}
	}
	for _, v := range o._orderPbs {
		if v.PromotionId == coupon.GetDomainId() {
			return order.ErrPromotionApplied
		}
	}

	o._orderPbs = append(o._orderPbs, &order.OrderPromotionBind{
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
		SaveFee: coupon.GetCouponFee(o._value.GoodsAmount),
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
	//v.CouponFee = coupon.GetCouponFee(v.TotalFee)
	//v.PayFee = o.GetPaymentFee()
	//v.DiscountFee = v.DiscountFee + v.CouponFee
	return nil
}

// 获取支付金额
//func (o *orderImpl) GetPaymentFee() float32 {
//	return o._value.PayFee - o._value.CouponFee
//}

// 获取应用的优惠券
func (o *orderImpl) GetCoupons() []promotion.ICouponPromotion {
	if o._coupons == nil {
		return make([]promotion.ICouponPromotion, 0)
	}
	return o._coupons
}

// 获取可用的促销,不包含优惠券
func (o *orderImpl) GetAvailableOrderPromotions() []promotion.IPromotion {
	if o._availPromotions == nil {
		//merchantId := o._cart.VendorId

		//todo: 将购物车中的vendor均获取出来
		merchantId := -1
		var vp []*promotion.PromotionInfo = o._promRep.GetPromotionOfMerchantOrder(merchantId)
		var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
		for i, v := range vp {
			proms[i] = o._promRep.CreatePromotion(v)
		}
		return proms
	}
	return o._availPromotions
}

// 获取促销绑定
func (o *orderImpl) GetPromotionBinds() []*order.OrderPromotionBind {
	if o._orderPbs == nil {
		o._orderPbs = o._orderRep.GetOrderPromotionBinds(o._value.OrderNo)
	}
	return o._orderPbs
}

// 获取最省的促销
func (o *orderImpl) GetBestSavePromotion() (p promotion.IPromotion, saveFee float32, integral int) {
	//todo: not implement
	return nil, 0, 0
}

// 设置配送地址
func (o *orderImpl) SetDeliver(deliverAddressId int) error {
	return o.setAddress(deliverAddressId)
}

// 设置配送地址
func (o *orderImpl) setAddress(addressId int) error {
	if addressId <= 0 {
		return order.ErrNoAddress
	}
	buyer := o.GetBuyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	addr := buyer.Profile().GetDeliver(addressId)
	if addr == nil {
		return order.ErrNoAddress
	}
	d := addr.GetValue()
	v := o._value
	v.ShippingAddress = strings.Replace(d.Area, " ", "", -1) + d.Address
	v.ConsigneePerson = d.RealName
	v.ConsigneePhone = d.Phone
	v.ShippingTime = time.Now().Add(-time.Hour).Unix()
	return nil
}

// 获取购买的会员
func (o *orderImpl) GetBuyer() member.IMember {
	if o._buyer == nil {
		o._buyer = o._memberRep.GetMember(o._value.BuyerId)
	}
	return o._buyer
}

// 获取支付单
func (o *orderImpl) GetPaymentOrder() payment.IPaymentOrder {
	if o._paymentOrder == nil {
		if o.GetAggregateRootId() <= 0 {
			panic(" Get payment order error ; because of order no yet created!")
		}
		o._paymentOrder = o._payRep.GetPaymentBySalesOrderId(o.GetAggregateRootId())
	}
	return o._paymentOrder
}

//************* 订单提交 ***************//

// 读取购物车数据,用于预生成订单
func (o *orderImpl) RequireCart(c cart.ICart) error {
	if o.GetAggregateRootId() > 0 || o._cart != nil {
		return order.ErrRequireCart
	}
	items := c.GetValue().Items
	if len(items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	// 绑定结算购物车
	o._cart = c
	// 将购物车的商品分类整理
	o._vendorItemsMap = o.buildVendorItemMap(items)
	// 更新订单的金额
	o._vendorExpressMap = o.updateOrderFee(o._vendorItemsMap)
	// 状态设为待支付
	o._value.State = 1

	return nil
}

// 加入运费计算器
func (o *orderImpl) addItemToExpressCalculator(ue express.IUserExpress,
	item *order.OrderItem, cul express.IExpressCalculator) {
	tpl := ue.GetTemplate(item.ExpressTplId)
	if tpl != nil {
		var err error
		v := tpl.Value()
		switch v.Basis {
		case express.BasisByNumber:
			err = cul.Add(item.ExpressTplId, float32(item.Quantity))
		case express.BasisByWeight:
			err = cul.Add(item.ExpressTplId, item.Weight)
		case express.BasisByVolume:
			err = cul.Add(item.ExpressTplId, item.Weight)
		}
		if err != nil {
			log.Println("[ Order][ Express][ Error]:", err)
		}
	}
}

// 更新订单金额,并返回运费
func (o *orderImpl) updateOrderFee(mp map[int][]*order.OrderItem) map[int]float32 {
	o._value.GoodsAmount = 0
	expCul := make(map[int]express.IExpressCalculator)
	expressMap := make(map[int]float32)
	for k, v := range mp {
		userExpress := o._expressRep.GetUserExpress(k)
		expCul[k] = userExpress.CreateCalculator()
		for _, item := range v {
			//计算商品总金额
			o._value.GoodsAmount += item.Amount
			//计算商品优惠金额
			o._value.DiscountAmount += item.Amount - item.FinalAmount
			//加入运费计算器
			o.addItemToExpressCalculator(userExpress, item, expCul[k])
		}
		//计算商户的运费
		expCul[k].Calculate("") //todo: 传入城市地区编号
		expressMap[k] = expCul[k].Total()
		//叠加运费
		o._value.ExpressFee += expressMap[k]
	}
	o._value.PackageFee = 0
	//计算最终金额
	o._value.FinalAmount = o._value.GoodsAmount - o._value.DiscountAmount +
		o._value.ExpressFee + o._value.PackageFee
	return expressMap
}

// 根据运营商获取商品和运费信息,限未生成的订单
func (o *orderImpl) GetByVendor() (items map[int][]*order.OrderItem,
	expressFeeMap map[int]float32) {
	if o._vendorItemsMap == nil {
		panic("订单尚未读取购物车!")
	}
	if o._vendorExpressMap == nil {
		panic("订单尚未计算金额")
	}
	items = o._vendorItemsMap
	expressFeeMap = o._vendorExpressMap
	return items, expressFeeMap
}

// 检查购物车
func (o *orderImpl) checkCart() error {
	if o._cart == nil || len(o._cart.GetValue().Items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	return o._cart.Check()
}

// 生成运营商与订单商品的映射
func (o *orderImpl) buildVendorItemMap(items []*cart.CartItem) map[int][]*order.OrderItem {
	mp := make(map[int][]*order.OrderItem)
	for _, v := range items {
		//必须勾选为结算
		if v.Checked == 1 {
			item := o.parseCartToOrderItem(v)
			if item == nil {
				domain.HandleError(errors.New("转换购物车商品到订单商品时出错: 商品SKU"+
					strconv.Itoa(v.SkuId)), "domain")
				continue
			}
			list, ok := mp[v.VendorId]
			if !ok {
				list = []*order.OrderItem{}
			}
			mp[v.VendorId] = append(list, item)
			//log.Println("--- vendor map len:", len(mp[v.VendorId]))
		}
	}
	return mp
}

// 转换购物车的商品项为订单项目
func (o *orderImpl) parseCartToOrderItem(c *cart.CartItem) *order.OrderItem {
	gs := o._saleRep.GetSale(c.VendorId).GoodsManager().CreateGoods(
		&goods.ValueGoods{Id: c.SkuId, SkuId: c.SkuId})
	// 获取商品已销售快照
	snap := gs.SnapshotManager().GetLatestSaleSnapshot()
	if snap == nil {
		domain.HandleError(errors.New("商品快照生成失败："+
			strconv.Itoa(c.SkuId)), "domain")
		return nil
	}
	fee := c.SalePrice * float32(c.Quantity)
	return &order.OrderItem{
		Id:          0,
		VendorId:    c.VendorId,
		ShopId:      c.ShopId,
		SkuId:       c.SkuId,
		SnapshotId:  snap.Id,
		Quantity:    c.Quantity,
		Amount:      fee,
		FinalAmount: fee,
		//是否配送
		IsShipped: 0,
		// 退回数量
		ReturnQuantity: 0,
		ExpressTplId:   c.Snapshot.ExpressTplId,
		Weight:         c.Snapshot.Weight * float32(c.Quantity), //计算重量
	}
}

// 提交订单，返回订单号。如有错误则返回
func (o *orderImpl) Submit() (string, error) {
	if o.GetAggregateRootId() != 0 {
		return "", errors.New("订单不允许重复提交")
	}
	err := o.checkCart()
	if err != nil {
		return "", err
	}
	buyer := o.GetBuyer()
	if buyer == nil {
		return "", member.ErrNoSuchMember
	}
	cv := o._cart.GetValue()

	err = o.setAddress(cv.DeliverId)
	if err != nil {
		return "", err
	}

	v := o._value

	//todo: best promotion , 优惠券和返现这里需要重构,直接影响到订单金额
	//prom,fee,integral := o.GetBestSavePromotion()

	// 应用优惠券
	if err := o.applyCouponOnSubmit(v); err != nil {
		return "", err
	}

	// 判断商品的优惠促销,如返现等
	proms, fee := o.applyCartPromotionOnSubmit(v, o._cart)
	if len(proms) != 0 {
		v.DiscountAmount += float32(fee)
		v.FinalAmount = v.GoodsAmount - v.DiscountAmount
		if v.FinalAmount < 0 {
			// 如果出现优惠券多余的金额也一并使用
			v.FinalAmount = 0
		}
	}

	o.avgDiscountToItem()

	// 检查是否已支付完成
	o.checkNewOrderPayment()

	// 保存订单
	orderId, err := o.saveNewOrderOnSubmit()
	v.Id = orderId
	if err == nil {
		// 绑定优惠券促销
		o.bindCouponOnSubmit(v.OrderNo)
		// 绑定购物车商品的促销
		for _, p := range proms {
			o.bindPromotionOnSubmit(v.OrderNo, p)
		}
		// 扣除库存
		o.applyGoodsNum()
		// 拆单
		o.breakUpByVendor()

		// 记录余额支付记录
		//todo: 扣减余额
		//if v.BalanceDiscount > 0 {
		//	err = acc.PaymentDiscount(v.OrderNo, v.BalanceDiscount)
		//}
	}
	return v.OrderNo, err
}

// 平均优惠抵扣金额到商品
func (o *orderImpl) avgDiscountToItem() {
	if o._vendorItemsMap == nil {
		panic(errors.New("仅能在下单时进行商品抵扣均分"))
	}
	if o._value.DiscountAmount > 0 {
		totalFee := o._value.GoodsAmount
		disFee := o._value.DiscountAmount
		for _, items := range o._vendorItemsMap {
			for _, v := range items {
				v.FinalAmount = v.Amount - (v.Amount/totalFee)*disFee
			}
		}
	}
}

// 绑定促销优惠
func (o *orderImpl) bindPromotionOnSubmit(orderNo string,
	prom promotion.IPromotion) (int, error) {
	var title string
	var integral int
	var fee int

	//todo: 需要重构,其他促销
	if prom.Type() == promotion.TypeFlagCashBack {
		fee = prom.GetRelationValue().(*promotion.ValueCashBack).BackFee
		title = prom.TypeName() + ":" + prom.GetValue().ShortName
	}

	v := &order.OrderPromotionBind{
		PromotionId:     prom.GetAggregateRootId(),
		PromotionType:   prom.Type(),
		OrderId:         o.GetAggregateRootId(),
		Title:           title,
		SaveFee:         float32(fee),
		PresentIntegral: integral,
		IsConfirm:       1,
		IsApply:         0,
	}
	return o._orderRep.SavePromotionBindForOrder(v)
}

// 应用购物车内商品的促销
func (o *orderImpl) applyCartPromotionOnSubmit(vo *order.Order,
	cart cart.ICart) ([]promotion.IPromotion, int) {
	var proms []promotion.IPromotion = make([]promotion.IPromotion, 0)
	var prom promotion.IPromotion
	var saveFee int
	var totalSaveFee int
	var intOrderFee = int(vo.FinalAmount)
	var rightBack bool

	for _, v := range cart.GetCartGoods() {
		prom = nil
		saveFee = 0
		rightBack = false

		// 判断商品的最省促销
		for _, v1 := range v.GetPromotions() {

			// 返现
			if v1.Type() == promotion.TypeFlagCashBack {
				vc := v1.GetRelationValue().(*promotion.ValueCashBack)
				if vc.MinFee < intOrderFee {
					if vc.BackFee > saveFee {
						prom = v1
						saveFee = vc.BackFee
						rightBack = vc.BackType == promotion.BackUseForOrder // 是否立即抵扣
					}
				}
			}

			//todo: 其他促销
		}

		if prom != nil {
			proms = append(proms, prom)
			if rightBack {
				totalSaveFee += saveFee
			}
		}
	}

	return proms, totalSaveFee
}

// 绑定订单与优惠券
func (o *orderImpl) bindCouponOnSubmit(orderNo string) {
	var oc *order.OrderCoupon = new(order.OrderCoupon)
	for _, c := range o.GetCoupons() {
		oc.Clone(c, o.GetAggregateRootId(), o._value.FinalAmount)
		o._orderRep.SaveOrderCouponBind(oc)
		// 绑定促销
		o.bindPromotionOnSubmit(orderNo, c.(promotion.IPromotion))
	}
}

// 在提交订单时应用优惠券
func (o *orderImpl) applyCouponOnSubmit(v *order.Order) error {
	var err error
	var t *promotion.ValueCouponTake
	var b *promotion.ValueCouponBind
	for _, c := range o.GetCoupons() {
		if c.CanTake() {
			t, err = c.GetTake(v.BuyerId)
			if err == nil {
				err = c.ApplyTake(t.Id)
			}
		} else {
			b, err = c.GetBind(v.BuyerId)
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
func (o *orderImpl) getBalanceDiscountFee(acc member.IAccount) float32 {
	if o._value.FinalAmount <= 0 || math.IsNaN(float64(o._value.FinalAmount)) {
		return 0
	}
	acv := acc.GetValue()
	if acv.Balance >= o._value.FinalAmount {
		return o._value.FinalAmount
	} else {
		return acv.Balance
	}
	return 0
}

// 检查新订单的支付结果,如果最终付款为0,则设置为已支付
// 有可能为多余的, 应等到支付单支付完成后,再通知订单支付完成。
func (o *orderImpl) checkNewOrderPayment() {
	// 校验是否支付
	//todo:  线下支付应设为等待确认
	//|| v.PaymentOpt == enum.PaymentOfflineCashPay ||
	//v.PaymentOpt == enum.PaymentRemit {
	//v.PaymentSign = 1

	// 设置订单状态
	if o._value.FinalAmount == 0 {
		o._value.IsPaid = 1
		o._value.PaidTime = time.Now().Unix()
		o._value.State = order.StatAwaitingConfirm
	} else if o._value.State == 0 {
		o._value.State = order.StatAwaitingPayment
	}
}

// 保存订单
func (o *orderImpl) saveNewOrderOnSubmit() (int, error) {
	unix := time.Now().Unix()
	o._value.ItemsInfo = string(o._cart.GetJsonItems())
	o._value.OrderNo = o._manager.GetFreeOrderNo(0)
	o._value.CreateTime = unix
	o._value.UpdateTime = unix

	id, err := o._orderRep.SaveOrder(o._value)
	if err == nil {
		o._value.Id = id
		// 释放购物车并销毁
		if o._cart.Release() {
			o._cart.Destroy()
		}
	}
	return id, err
}

// 保存订单
func (o *orderImpl) Save() (int, error) {
	// 有操作后解除挂起状态
	// todo: ???
	//if o._value.IsSuspend == 1 && !o._internalSuspend {
	//    o._value.IsSuspend = 0
	//}

	if o._value.Id > 0 {
		return o._orderRep.SaveOrder(o._value)
	}
	o._internalSuspend = false
	return 0, errors.New("please use Order.Submit() save new order.")
}

// 根据运营商生成子订单
func (o *orderImpl) createSubOrderByVendor(parentOrderId int, buyerId int,
	vendorId int, newOrderNo bool, items []*order.OrderItem) order.ISubOrder {
	orderNo := o.GetOrderNo()
	if newOrderNo {
		orderNo = o._manager.GetFreeOrderNo(vendorId)
	}

	if len(items) == 0 {
		domain.HandleError(errors.New("拆分订单,运营商下未获取到商品,订单:"+
			o.GetOrderNo()), "domain")
		return nil
	}

	v := &order.SubOrder{
		OrderNo:   orderNo,
		BuyerId:   buyerId,
		VendorId:  vendorId,
		ParentId:  parentOrderId,
		Subject:   "子订单",
		ShopId:    items[0].ShopId,
		ItemsInfo: "",
		// 总金额
		GoodsAmount: 0,
		// 减免金额(包含优惠券金额)
		DiscountAmount: 0,
		ExpressFee:     0,
		FinalAmount:    0,
		// 是否挂起，如遇到无法自动进行的时挂起，来提示人工确认。
		IsSuspend:  0,
		Note:       "",
		Remark:     "",
		State:      order.StatAwaitingPayment,
		UpdateTime: o._value.UpdateTime,
		Items:      items,
	}
	// 计算订单金额
	for _, item := range items {
		//计算商品金额
		v.GoodsAmount += item.Amount
		//计算商品优惠金额
		v.DiscountAmount += item.Amount - item.FinalAmount
	}
	// 设置运费
	v.ExpressFee = o._vendorExpressMap[vendorId]
	// 设置包装费
	v.PackageFee = 0
	// 最终金额 = 商品金额 - 商品抵扣金额(促销折扣) + 包装费 + 快递费
	v.FinalAmount = v.GoodsAmount - v.DiscountAmount + v.PackageFee + v.ExpressFee
	// 判断是否已支付
	if o._value.IsPaid == 1 {
		v.State = enum.ORDER_WAIT_CONFIRM
	}
	return o._manager.CreateSubOrder(v)
}

//根据运营商拆单,返回拆单结果,及拆分的订单数组
func (o *orderImpl) breakUpByVendor() []order.ISubOrder {
	parentOrderId := o.GetAggregateRootId()
	if parentOrderId <= 0 ||
		o._vendorItemsMap == nil ||
		len(o._vendorItemsMap) == 0 {
		//todo: 订单要取消掉
		panic(fmt.Sprintf("订单异常: 订单未生成或VendorItemMap为空,"+
			"订单编号:%d,订单号:%s,vendor len:%d",
			parentOrderId, o._value.OrderNo, len(o._vendorItemsMap)))
	}
	l := len(o._vendorItemsMap)
	list := make([]order.ISubOrder, l)
	i := 0
	buyerId := o._buyer.GetAggregateRootId()
	for k, v := range o._vendorItemsMap {
		//log.Println("----- vendor ", k, len(v),l)
		list[i] = o.createSubOrderByVendor(parentOrderId, buyerId, k, l > 1, v)
		if _, err := list[i].Save(); err != nil {
			domain.HandleError(err, "domain")
		}
		i++
	}
	return list
}

// 扣除库存
func (o *orderImpl) applyGoodsNum() {
	for _, v := range o._vendorItemsMap {
		for _, v2 := range v {
			o.takeGoodsStock(v2.VendorId, v2.SkuId, v2.Quantity)
		}
	}
}

//****************  订单提交结束 **************//

// 设置支付方式
//func (o *orderImpl) SetPayment(payment int) {
//	o._value.PaymentOpt = payment
//}

// 获取子订单列表
func (o *orderImpl) GetSubOrders() []order.ISubOrder {
	if o.GetAggregateRootId() <= 0 {
		panic(order.ErrNoYetCreated)
	}
	if o._subList == nil {
		subList := o._orderRep.GetSubOrdersByParentId(o.GetAggregateRootId())
		for _, v := range subList {
			o._subList = append(o._subList,
				o._manager.CreateSubOrder(v))
		}
	}
	return o._subList
}

// 在线支付交易完成
func (o *orderImpl) OnlinePaymentTradeFinish() (err error) {
	for _, o := range o.GetSubOrders() {
		err = o.PaymentFinishByOnlineTrade()
		if err != nil {
			return err
		}
	}
	return nil

	//todo:
	if o._value.IsPaid == 1 {
		return order.ErrOrderPayed
	}
	unix := time.Now().Unix()
	o._value.IsPaid = 1
	o._value.UpdateTime = unix
	o._value.PaidTime = unix
	o._value.State = order.StatAwaitingConfirm

	o._manager.SmartConfirmOrder(o) // 确认订单

	_, err = o.Save()
	return err
}

// 使用余额支付
func (o *orderImpl) paymentWithBalance(buyerType int) error {
	if o._value.IsPaid == 1 {
		return order.ErrOrderPayed
	}
	acc := o._memberRep.GetMember(o._value.BuyerId).GetAccount()
	if fee := o.getBalanceDiscountFee(acc); fee == 0 {
		return member.ErrAccountBalanceNotEnough
	} else {
		o._value.DiscountAmount = fee
		o._value.FinalAmount -= fee
		err := acc.PaymentDiscount(o.GetOrderNo(), fee, "")
		if err != nil {
			return err
		}
	}
	unix := time.Now().Unix()
	if o._value.FinalAmount == 0 {
		o._value.IsPaid = 1
		// o._value.PaymentSign = buyerType
		o._value.State = enum.ORDER_WAIT_CONFIRM
	}
	o._value.UpdateTime = unix
	o._value.PaidTime = unix
	_, err := o.Save()
	return err
}

// 使用余额支付
func (o *orderImpl) PaymentWithBalance() error {
	return o.paymentWithBalance(payment.PaymentByBuyer)
}

// 客服使用余额支付
func (o *orderImpl) CmPaymentWithBalance() error {
	return o.paymentWithBalance(payment.PaymentByCM)
}

// 添加日志
func (o *orderImpl) AppendLog(l *order.OrderLog) error {
	if o.GetAggregateRootId() <= 0 {
		return errors.New("order not created.")
	}
	l.OrderId = o.GetAggregateRootId()
	l.RecordTime = time.Now().Unix()
	return o._orderRep.SaveSubOrderLog(l)
}

// 订单是否已完成
func (o *orderImpl) IsOver() bool {
	s := o._value.State
	return s == enum.ORDER_CANCEL || s == enum.ORDER_COMPLETED
}

// 处理订单
func (o *orderImpl) Process() error {
	dt := time.Now()
	o._value.State += 1
	o._value.UpdateTime = dt.Unix()

	_, err := o.Save()
	return err
}

// 确认订单
func (o *orderImpl) Confirm() error {
	return nil
}

// 扣减商品库存
func (o *orderImpl) takeGoodsStock(vendorId, skuId, quantity int) error {
	gds := o._saleRep.GetSale(vendorId).GoodsManager().GetGoods(skuId)
	if gds == nil {
		return goods.ErrNoSuchGoods
	}
	return gds.TakeStock(quantity)
}

// 配送订单
func (o *orderImpl) Deliver(spId int, spNo string) error {
	//todo: 记录快递配送信息
	dt := time.Now()
	o._value.State += 1
	o._value.ShippingTime = dt.Unix()
	o._value.UpdateTime = dt.Unix()

	_, err := o.Save()
	if err == nil {
		err = o.AppendLog(&order.OrderLog{
			Type:       int(order.LogSetup),
			OrderState: order.StatShipped,
			IsSystem:   0,
			Message:    "",
		})
	}
	return err
}

// 获取订单号
func (o *orderImpl) GetOrderNo() string {
	return o.GetValue().OrderNo
}

// 更新返现到会员账户
func (o *orderImpl) updateShoppingMemberBackFee(pt merchant.IMerchant,
	m member.IMember, fee float32, unixTime int64) error {
	if fee == 0 {
		return nil
	}
	v := o.GetValue()
	pv := pt.GetValue()

	//更新账户
	acc := m.GetAccount()
	acv := acc.GetValue()
	//acc.TotalFee += o._value.Fee
	//acc.TotalPay += o._value.PayFee
	acv.PresentBalance += fee // 更新赠送余额
	acv.TotalPresentFee += fee
	acv.UpdateTime = unixTime
	_, err := acc.Save()
	if err == nil {
		//给自己返现
		tit := fmt.Sprintf("订单:%s(商户:%s)返现￥%.2f元", v.OrderNo, pv.Name, fee)
		err = acc.ChargeForPresent(tit, v.OrderNo, float32(fee),
			member.DefaultRelateUser)
	}
	return err
}

// 处理返现促销
func (o *orderImpl) handleCashBackPromotions(pt merchant.IMerchant,
	m member.IMember) error {
	proms := o.GetPromotionBinds()
	for _, v := range proms {
		if v.PromotionType == promotion.TypeFlagCashBack {
			c := o._promRep.GetPromotion(v.PromotionId)
			return o.handleCashBackPromotion(pt, m, v, c)
		}
	}
	return nil
}

// 处理返现促销
func (o *orderImpl) handleCashBackPromotion(pt merchant.IMerchant,
	m member.IMember,
	v *order.OrderPromotionBind, pm promotion.IPromotion) error {
	cpv := pm.GetRelationValue().(*promotion.ValueCashBack)

	//更新账户
	bFee := float32(cpv.BackFee)
	acc := m.GetAccount()
	acv := acc.GetValue()
	acv.PresentBalance += bFee // 更新赠送余额
	acv.TotalPresentFee += bFee
	// 赠送金额，不应该计入到余额，可采取充值到余额
	//acc.Balance += float32(cpv.BackFee)                            // 更新账户余额

	acv.UpdateTime = time.Now().Unix()
	_, err := acc.Save()

	if err == nil {
		// 优惠绑定生效
		v.IsApply = 1
		o._orderRep.SavePromotionBindForOrder(v)

		// 处理自定义返现
		c := pm.(promotion.ICashBackPromotion)
		HandleCashBackDataTag(m, o._value, c, o._memberRep)

		//给自己返现
		tit := fmt.Sprintf("返现￥%d元,订单号:%s", cpv.BackFee, o._value.OrderNo)
		err = acc.ChargeForPresent(tit, o.GetOrderNo(), float32(cpv.BackFee),
			member.DefaultRelateUser)
	}
	return err
}

//todo: ?? 自动收货功能

var _ order.ISubOrder = new(subOrderImpl)

// 子订单实现
type subOrderImpl struct {
	_value           *order.SubOrder
	_parent          order.IOrder
	_buyer           member.IMember
	_internalSuspend bool //内部挂起
	_rep             order.IOrderRep
	_memberRep       member.IMemberRep
	_goodsRep        goods.IGoodsRep
	_saleRep         sale.ISaleRep
	_manager         order.IOrderManager
	_shipRep         shipment.IShipmentRep
	_valRep          valueobject.IValueRep
	_mchRep          merchant.IMerchantRep
}

func NewSubOrder(v *order.SubOrder,
	manager order.IOrderManager, rep order.IOrderRep,
	mmRep member.IMemberRep, goodsRep goods.IGoodsRep,
	shipRep shipment.IShipmentRep, saleRep sale.ISaleRep,
	valRep valueobject.IValueRep,
	mchRep merchant.IMerchantRep) order.ISubOrder {
	return &subOrderImpl{
		_value:     v,
		_manager:   manager,
		_rep:       rep,
		_memberRep: mmRep,
		_goodsRep:  goodsRep,
		_saleRep:   saleRep,
		_shipRep:   shipRep,
		_valRep:    valRep,
		_mchRep:    mchRep,
	}
}

// 获取领域对象编号
func (o *subOrderImpl) GetDomainId() int {
	return o._value.Id
}

// 获取值对象
func (o *subOrderImpl) GetValue() *order.SubOrder {
	return o._value
}

// 获取商品项
func (o *subOrderImpl) Items() []*order.OrderItem {
	if (o._value.Items == nil || len(o._value.Items) == 0) &&
		o.GetDomainId() > 0 {
		o._value.Items = o._rep.GetSubOrderItems(o.GetDomainId())
	}
	return o._value.Items
}

// 获取父订单
func (o *subOrderImpl) Parent() order.IOrder {
	if o._parent == nil {
		o._parent = o._manager.GetOrderById(o._value.ParentId)
	}
	return o._parent
}

// 获取购买的会员
func (o *subOrderImpl) GetBuyer() member.IMember {
	if o._buyer == nil {
		//if o._value.BuyerId <= 0 {
		//    panic(errors.New("订单BuyerId非会员或未设置"))
		//}
		o._buyer = o._memberRep.GetMember(o._value.BuyerId)
	}
	return o._buyer
}

// 添加备注
func (o *subOrderImpl) AddRemark(remark string) {
	o._value.Remark = remark
}

// 设置Shop
func (o *subOrderImpl) SetShop(shopId int) error {
	//todo:验证Shop
	o._value.ShopId = shopId
	if o._value.State == enum.ORDER_WAIT_CONFIRM {
		panic("not impl")
		// o.Confirm()
	}
	return nil
}

func (o *subOrderImpl) saveOrderItems() error {
	unix := time.Now().Unix()
	id := o.GetDomainId()
	for _, v := range o.Items() {
		v.OrderId = id
		v.UpdateTime = unix
		_, err := o._rep.SaveOrderItem(id, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// 保存订单
func (o *subOrderImpl) Save() (int, error) {
	unix := time.Now().Unix()
	o._value.UpdateTime = unix
	if o.GetDomainId() > 0 {
		return o._rep.SaveSubOrder(o._value)
	}
	if o._value.CreateTime <= 0 {
		o._value.CreateTime = unix
	}
	id, err := o._rep.SaveSubOrder(o._value)
	if err == nil {
		o._value.Id = id
		err = o.saveOrderItems()
		o.AppendLog(order.LogSetup, true, "{created}")
	}
	return id, err
}

// 订单完成支付
func (o *subOrderImpl) orderFinishPaid() error {
	if o._value.IsPaid == 1 {
		return order.ErrOrderPayed
	}
	if o._value.State == order.StatAwaitingPayment {
		o._value.IsPaid = 1
		o._value.State = order.StatAwaitingConfirm
		err := o.AppendLog(order.LogSetup, true, "{finish_pay}")
		if err == nil {
			_, err = o.Save()
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
	o._value.IsSuspend = 1
	o._internalSuspend = true
	o._value.UpdateTime = time.Now().Unix()
	_, err := o.Save()
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
		OrderState: int(o._value.State),
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return o._rep.SaveSubOrderLog(l)
}

// 确认订单
func (o *subOrderImpl) Confirm() (err error) {
	//todo: 线下交易,自动确认
	//if o._value.PaymentOpt == enum.PaymentOnlinePay &&
	//o._value.IsPaid == enum.FALSE {
	//    return order.ErrOrderNotPayed
	//}
	if o._value.State < order.StatAwaitingConfirm {
		return order.ErrOrderNotPayed
	}
	if o._value.State >= order.StatAwaitingPickup {
		return order.ErrOrderHasConfirm
	}
	o._value.State = order.StatAwaitingPickup
	o._value.UpdateTime = time.Now().Unix()
	_, err = o.Save()
	if err == nil {
		go o.addSalesNum()
		err = o.AppendLog(order.LogSetup, false, "{confirm}")
	}
	return err
}

// 增加商品的销售数量
func (o *subOrderImpl) addSalesNum() {
	gm := o._saleRep.GetSale(o._value.VendorId).GoodsManager()
	for _, v := range o.Items() {
		gds := gm.GetGoods(v.SkuId)
		gds.AddSalesNum(v.Quantity)
	}
}

// 捡货(备货)
func (o *subOrderImpl) PickUp() error {
	if o._value.State < order.StatAwaitingPickup {
		return order.ErrOrderNotConfirm
	}
	if o._value.State >= order.StatAwaitingShipment {
		return order.ErrOrderHasPickUp
	}
	o._value.State = order.StatAwaitingShipment
	o._value.UpdateTime = time.Now().Unix()
	_, err := o.Save()
	if err == nil {
		err = o.AppendLog(order.LogSetup, true, "{pickup}")
	}
	return err
}

// 发货
func (o *subOrderImpl) Ship(spId int, spOrder string) error {
	//so := o._shipRep.GetOrders()
	//todo: 可进行发货修改
	if o._value.State < order.StatAwaitingShipment {
		return order.ErrOrderNotPickUp
	}
	if o._value.State >= order.StatShipped {
		return order.ErrOrderShipped
	}

	if list := o._shipRep.GetOrders(o.GetDomainId()); len(list) > 0 {
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
		o._value.State = order.StatShipped
		o._value.UpdateTime = time.Now().Unix()
		if _, err = o.Save(); err != nil {
			return err
		}
		// 保存商品的发货状态
		err = o.saveOrderItems()
		o.AppendLog(order.LogSetup, true, "{shipped}")
	}
	return err
}

func (o *subOrderImpl) createShipmentOrder(items []*order.OrderItem) shipment.IShipmentOrder {
	if items == nil || len(items) == 0 {
		return nil
	}
	unix := time.Now().Unix()
	so := &shipment.ShipmentOrder{
		Id:         0,
		OrderId:    o.GetDomainId(),
		ExpressLog: "",
		ShipTime:   unix,
		Stat:       shipment.StatAwaitingShipment,
		UpdateTime: unix,
		Items:      []*shipment.Item{},
	}
	for _, v := range items {
		if v.IsShipped == 1 {
			continue
		}
		so.Amount += v.Amount
		so.FinalAmount += v.FinalAmount
		so.Items = append(so.Items, &shipment.Item{
			Id:          0,
			GoodsSnapId: v.SnapshotId,
			Quantity:    v.Quantity,
			Amount:      v.Amount,
			FinalAmount: v.FinalAmount,
		})
		v.IsShipped = 1
	}
	return o._shipRep.CreateShipmentOrder(so)
}

// 已收货
func (o *subOrderImpl) BuyerReceived() error {
	var err error
	if o._value.State < order.StatShipped {
		return order.ErrOrderNotShipped
	}
	if o._value.State >= order.StatCompleted {
		return order.ErrIsCompleted
	}
	dt := time.Now()
	o._value.State = order.StatCompleted
	o._value.UpdateTime = dt.Unix()
	o._value.IsSuspend = 0
	if _, err = o.Save(); err != nil {
		return err
	}
	err = o.AppendLog(order.LogSetup, true, "{completed}")
	if err == nil {
		go o.vendorSettle()
		// 执行其他的操作
		if err2 := o.onOrderComplete(); err != nil {
			domain.HandleError(err2, "domain")
		}
	}
	return err
}

func (s *subOrderImpl) getOrderAmount() (amount float32, refund float32) {
	items := s.Items()
	for _, item := range items {
		if item.ReturnQuantity > 0 {
			a := item.Amount / float32(item.Quantity) * float32(item.ReturnQuantity)
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
		amount += s._value.ExpressFee + s._value.PackageFee
	}
	return amount, refund
}

// 获取订单的成本
func (s *subOrderImpl) getOrderCost() float32 {
	var cost float32
	items := s.Items()
	for _, item := range items {
		snap := s._goodsRep.GetSaleSnapshot(item.SnapshotId)
		cost += snap.Cost * float32(item.Quantity-item.ReturnQuantity)
	}
	//如果非全部退货、退款,则加上运费及包装费
	if cost > 0 {
		cost += s._value.ExpressFee + s._value.PackageFee
	}
	return cost
}

// 商户结算
func (s *subOrderImpl) vendorSettle() error {
	vendor := s._mchRep.GetMerchant(s._value.VendorId)
	if vendor != nil {
		conf := s._valRep.GetGlobMchSaleConf()
		switch conf.MchOrderSettleMode {
		case enum.MchModeSettleByCost:
			return s.vendorSettleByCost(vendor)
		case enum.MchModeSettleByRate:
			return s.vendorSettleByRate(vendor, conf.MchOrderSettleRate)
		}

	}
	return nil
}

// 根据供货价进行商户结算
func (s *subOrderImpl) vendorSettleByCost(vendor merchant.IMerchant) error {
	_, refund := s.getOrderAmount()
	sAmount := s.getOrderCost()
	if sAmount > 0 {
		return vendor.Account().SettleOrder(s._value.OrderNo,
			sAmount, 0, refund, "订单结算")
	}
	return nil
}

// 根据比例进行商户结算
func (s *subOrderImpl) vendorSettleByRate(vendor merchant.IMerchant, rate float32) error {
	amount, refund := s.getOrderAmount()
	sAmount := amount * rate
	if sAmount > 0 {
		return vendor.Account().SettleOrder(s._value.OrderNo,
			sAmount, 0, refund, "订单结算")
	}
	return nil
}

// 获取订单的日志
func (o *subOrderImpl) LogBytes() []byte {
	buf := bytes.NewBufferString("")
	list := o._rep.GetSubOrderLogs(o.GetDomainId())
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
	if o._value.State != order.StatCompleted {
		return order.ErrUnusualOrderStat
	}
	var err error
	ov := o._value
	conf := o._valRep.GetGlobNumberConf()
	registry := o._valRep.GetRegistry()
	amount := ov.FinalAmount
	acc := m.GetAccount()

	// 增加经验
	if registry.MemberExperienceEnabled {
		rate := conf.ExperienceRateByOrder
		if exp := int(amount * rate); exp > 0 {
			if err = m.AddExp(exp); err != nil {
				return err
			}
		}
	}

	// 增加积分
	//todo: 增加阶梯的返积分,比如订单满30送100积分
	integral := int(amount*conf.IntegralRateByConsumption) + conf.IntegralBackExtra
	// 赠送积分
	if integral > 0 {
		err = m.GetAccount().AddIntegral(member.TypeIntegralShoppingPresent,
			o._value.OrderNo, integral, "")
		if err != nil {
			return err
		}
	}
	acv := acc.GetValue()
	acv.TotalConsumption += ov.GoodsAmount
	acv.TotalPay += ov.FinalAmount
	acv.UpdateTime = time.Now().Unix()
	_, err = acc.Save()
	return err
}

// 取消商品
func (o *subOrderImpl) cancelGoods() error {
	gm := o._saleRep.GetSale(o._value.VendorId).GoodsManager()
	for _, v := range o.Items() {
		snapshot := o._goodsRep.GetSaleSnapshot(v.SnapshotId)
		if snapshot == nil {
			return goods.ErrNoSuchSnapshot
		}
		var gds sale.IGoods = gm.GetGoods(snapshot.SkuId)
		if gds != nil {
			// 释放库存
			gds.FreeStock(v.Quantity)
			// 如果订单已付款，则取消销售数量
			if o._value.IsPaid == 1 {
				gds.CancelSale(snapshot.SkuId, o._value.OrderNo)
			}
		}
	}
	return nil
}

// 取消支付单
func (o *subOrderImpl) cancelPaymentOrder() error {
	po := o.Parent().GetPaymentOrder()
	if po != nil {
		// 订单金额为0,则取消订单
		if po.GetValue().FinalAmount-o._value.FinalAmount == 0 {
			return po.Cancel()
		}
		return po.Adjust(o._value.FinalAmount)
	}
	return nil
}

// 订单退款
func (o *subOrderImpl) backupPayment() (err error) {
	po := o.Parent().GetPaymentOrder()
	if po == nil {
		panic(errors.New("无法获取支付单,订单号:" + o.Parent().GetOrderNo()))
	}
	if pv := po.GetValue(); pv.BalanceDiscount > 0 {
		//退回账户余额抵扣
		acc := o.GetBuyer().GetAccount()
		err = acc.ChargeForBalance(member.ChargeByRefund, "订单退款",
			o._value.OrderNo, o._value.DiscountAmount, member.DefaultRelateUser)
	}
	if o._value.FinalAmount > 0 {
		//todo: 其他支付方式退还,如网银???
	}
	return nil
}

// 取消订单
func (o *subOrderImpl) Cancel(reason string) error {
	if o._value.State == order.StatCancelled {
		return order.ErrOrderCanNotCancel
	}
	if o._value.State > order.StatAwaitingPayment {
		return order.ErrDisallowCancel
	}

	o._value.State = order.StatCancelled
	o._value.UpdateTime = time.Now().Unix()
	_, err := o.Save()
	if err == nil {
		domain.HandleError(o.AppendLog(order.LogSetup, true, reason), "domain")
		// 取消支付单
		domain.HandleError(o.cancelPaymentOrder(), "domain")
		// 取消商品
		err = o.cancelGoods()
		//如果已付款,则取消订单
		if err == nil {
			err = o.backupPayment()
		}
	}
	return err
}

// 退回商品
func (o *subOrderImpl) Return(snapshotId int, quantity int) error {
	for _, v := range o.Items() {
		if v.SnapshotId == snapshotId {
			if v.Quantity-v.ReturnQuantity < quantity {
				return order.ErrOutOfQuantity
			}
			v.ReturnQuantity += quantity
			_, err := o._rep.SaveOrderItem(o.GetDomainId(), v)
			return err
		}
	}
	return order.ErrNoSuchGoodsOfOrder
}

// 撤销退回商品
func (o *subOrderImpl) RevertReturn(snapshotId int, quantity int) error {
	for _, v := range o.Items() {
		if v.SnapshotId == snapshotId {
			if v.ReturnQuantity < quantity {
				return order.ErrOutOfQuantity
			}
			v.ReturnQuantity -= quantity
			_, err := o._rep.SaveOrderItem(o.GetDomainId(), v)
			return err
		}
	}
	return order.ErrNoSuchGoodsOfOrder
}

// 申请退款
func (o *subOrderImpl) SubmitRefund(reason string) error {
	if o._value.State == order.StatAwaitingPayment {
		return o.Cancel("订单主动申请取消,原因:" + reason)
	}
	if o._value.State >= order.StatShipped ||
		o._value.State >= order.StatCancelled {
		return order.ErrOrderCanNotCancel
	}
	o._value.State = order.StatAwaitingCancel
	o._value.UpdateTime = time.Now().Unix()
	_, err := o.Save()
	return err
}

// 取消退款申请
func (o *subOrderImpl) CancelRefund() error {
	if o._value.State != order.StatAwaitingCancel || o._value.IsPaid == 0 {
		panic(errors.New("订单已经取消,不允许再退款"))
	}
	o._value.State = order.StatAwaitingConfirm
	o._value.UpdateTime = time.Now().Unix()
	_, err := o.Save()
	return err
}

// 谢绝订单
func (o *subOrderImpl) Decline(reason string) error {
	if o._value.State == order.StatAwaitingPayment {
		return o.Cancel("商户取消,原因:" + reason)
	}
	if o._value.State >= order.StatShipped ||
		o._value.State >= order.StatCancelled {
		return order.ErrOrderCanNotCancel
	}
	o._value.State = order.StatDeclined
	o._value.UpdateTime = time.Now().Unix()
	_, err := o.Save()
	return err
}

// 退款
func (o *subOrderImpl) Refund() error {
	// 已退款
	if o._value.State == order.StatRefunded ||
		o._value.State == order.StatCancelled {
		return order.ErrHasRefund
	}
	// 不允许退款
	if o._value.State != order.StatAwaitingCancel &&
		o._value.State != order.StatDeclined {
		return order.ErrDisallowRefund
	}
	o._value.State = order.StatRefunded
	o._value.UpdateTime = time.Now().Unix()
	_, err := o.Save()
	if err == nil {
		err = o.backupPayment()
	}
	return err
}

// 退款申请
//func (o *subOrderImpl) refund(reason string) error {
//    //todo: 商户谢绝订单,现仅处理用户提交的退款
//    ov := o._value
//    unix := time.Now().Unix()
//    rv := &afterSales.RefundOrder{
//        Id: 0,
//        // 订单编号
//        OrderId: o.GetDomainId(),
//        // 金额
//        Amount: ov.FinalAmount,
//        // 退款方式：1.退回余额  2: 原路退回
//        RefundType: 1,
//        // 是否为全部退款
//        AllRefund: 1,
//        // 退款的商品项编号
//        ItemId: 0,
//        // 联系人
//        PersonName: "",
//        // 联系电话
//        PersonPhone: "",
//        // 退款原因
//        Reason: reason,
//        // 退款单备注(系统)
//        Remark: "",
//        // 运营商备注
//        VendorRemark: "",
//        // 退款状态
//        State: afterSales.RefundStatAwaitingVendor,
//        // 提交时间
//        CreateTime: unix,
//        // 更新时间
//        UpdateTime: unix,
//    }
//    ro := o._afterSalesRep.CreateRefundOrder(rv)
//    return ro.Submit()
//}

// 完成订单
func (o *subOrderImpl) onOrderComplete() error {
	// 更新发货单
	soList := o._shipRep.GetOrders(o.GetDomainId())
	for _, v := range soList {
		domain.HandleError(v.Completed(), "domain")
	}

	// 获取消费者消息
	m := o.GetBuyer()
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
	m member.IMember, fee float32, unixTime int64) error {
	if fee <= 0 || math.IsNaN(float64(fee)) {
		return nil
	}
	v := o.GetValue()

	//更新账户
	acc := m.GetAccount()
	acv := acc.GetValue()
	//acc.TotalFee += o._value.Fee
	//acc.TotalPay += o._value.PayFee
	acv.PresentBalance += fee // 更新赠送余额
	acv.TotalPresentFee += fee
	acv.UpdateTime = unixTime
	acc.Save()

	//给自己返现
	tit := fmt.Sprintf("订单:%s(商户:%s)返现￥%.2f元", v.OrderNo, mchName, fee)
	return acc.ChargeForPresent(tit, v.OrderNo, float32(fee),
		member.DefaultRelateUser)
}
