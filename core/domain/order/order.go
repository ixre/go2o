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
	"errors"
	"fmt"
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
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure"
	"go2o/core/infrastructure/domain"
	"go2o/core/variable"
	"strconv"
	"strings"
	"time"
)

var (
	EXP_BIT float32
)
var _ order.IOrder = new(orderImpl)

//todo: 促销

type orderImpl struct {
	_manager         order.IOrderManager
	_value           *order.Order
	_cart            cart.ICart //购物车,仅在订单生成时设置
	_coupons         []promotion.ICouponPromotion
	_availPromotions []promotion.IPromotion
	_orderPbs        []*order.OrderPromotionBind
	_memberRep       member.IMemberRep
	_buyer           member.IMember
	_orderRep        order.IOrderRep
	_partnerRep      merchant.IMerchantRep //todo: can delete ?
	_expressRep      express.IExpressRep
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
}

func newOrder(shopping order.IOrderManager, value *order.Order,
	mchRep merchant.IMerchantRep, shoppingRep order.IOrderRep,
	goodsRep goods.IGoodsRep, saleRep sale.ISaleRep,
	promRep promotion.IPromotionRep, memberRep member.IMemberRep,
	expressRep express.IExpressRep, valRep valueobject.IValueRep) order.IOrder {
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
	}
}

func (this *orderImpl) GetAggregateRootId() int {
	return this._value.Id
}

func (this *orderImpl) GetValue() *order.Order {
	return this._value
}

// 设置订单值
func (this *orderImpl) SetValue(v *order.Order) error {
	v.Id = this.GetAggregateRootId()
	this._value = v
	return nil
}

// 应用优惠券
func (this *orderImpl) ApplyCoupon(coupon promotion.ICouponPromotion) error {
	//if this._coupons == nil {
	//	this._coupons = []promotion.ICouponPromotion{}
	//}
	//this._coupons = append(this._coupons, coupon)

	// 添加到促销信息中
	if this._orderPbs == nil {
		this._orderPbs = []*order.OrderPromotionBind{}
	}
	for _, v := range this._orderPbs {
		if v.PromotionId == coupon.GetDomainId() {
			return order.ErrPromotionApplied
		}
	}

	this._orderPbs = append(this._orderPbs, &order.OrderPromotionBind{
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
		SaveFee: coupon.GetCouponFee(this._value.GoodsFee),
		// 赠送积分
		PresentIntegral: 0, //todo;/////
		// 是否应用
		IsApply: 0,
		// 是否确认
		IsConfirm: 0,
	})

	//v := this._value
	//v.CouponCode = val.Code
	//v.CouponDescribe = coupon.GetDescribe()
	//v.CouponFee = coupon.GetCouponFee(v.TotalFee)
	//v.PayFee = this.GetPaymentFee()
	//v.DiscountFee = v.DiscountFee + v.CouponFee
	return nil
}

// 获取支付金额
//func (this *orderImpl) GetPaymentFee() float32 {
//	return this._value.PayFee - this._value.CouponFee
//}

// 获取应用的优惠券
func (this *orderImpl) GetCoupons() []promotion.ICouponPromotion {
	if this._coupons == nil {
		return make([]promotion.ICouponPromotion, 0)
	}
	return this._coupons
}

// 获取可用的促销,不包含优惠券
func (this *orderImpl) GetAvailableOrderPromotions() []promotion.IPromotion {
	if this._availPromotions == nil {
		//merchantId := this._cart.VendorId

		//todo: 将购物车中的vendor均获取出来
		merchantId := -1
		var vp []*promotion.PromotionInfo = this._promRep.GetPromotionOfMerchantOrder(merchantId)
		var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
		for i, v := range vp {
			proms[i] = this._promRep.CreatePromotion(v)
		}
		return proms
	}
	return this._availPromotions
}

// 获取促销绑定
func (this *orderImpl) GetPromotionBinds() []*order.OrderPromotionBind {
	if this._orderPbs == nil {
		this._orderPbs = this._orderRep.GetOrderPromotionBinds(this._value.OrderNo)
	}
	return this._orderPbs
}

// 获取最省的促销
func (this *orderImpl) GetBestSavePromotion() (p promotion.IPromotion, saveFee float32, integral int) {
	//todo: not implement
	return nil, 0, 0
}

// 设置支付方式
//func (this *orderImpl) SetPayment(payment int) {
//	this._value.PaymentOpt = payment
//}

// 在线交易支付
func (this *orderImpl) PaymentForOnlineTrade(serverProvider string, tradeNo string) error {
	if this._value.IsPaid == 1 {
		return order.ErrOrderPayed
	}
	unix := time.Now().Unix()
	this._value.IsPaid = 1
	this._value.UpdateTime = unix
	this._value.PaidTime = unix
	if this._value.Status == enum.ORDER_WAIT_PAYMENT {
		this._value.Status = enum.ORDER_WAIT_CONFIRM // 设置为待确认状态
	}
	this._manager.SmartConfirmOrder(this) // 确认订单
	_, err := this.Save()
	return err
}

// 设置配送地址
func (this *orderImpl) SetDeliver(deliverAddressId int) error {
	d := this._memberRep.GetSingleDeliverAddress(this._value.BuyerId, deliverAddressId)
	if d != nil {
		v := this._value
		v.ShippingAddress = d.Address
		v.ConsigneePerson = d.RealName
		v.ConsigneePhone = d.Phone
		v.ShippingTime = time.Now().Add(-time.Hour).Unix()
		return nil
	}
	return member.ErrNoSuchDeliverAddress
}

// 获取购买的会员
func (this *orderImpl) GetBuyer() member.IMember {
	if this._buyer == nil {
		//if this._value.BuyerId <= 0 {
		//    panic(errors.New("订单BuyerId非会员或未设置"))
		//}
		this._buyer = this._memberRep.GetMember(this._value.BuyerId)
	}
	return this._buyer
}

//************* 订单提交 ***************//

// 读取购物车数据,用于预生成订单
func (this *orderImpl) RequireCart(c cart.ICart) error {
	if this.GetAggregateRootId() > 0 || this._cart != nil {
		return order.ErrRequireCart
	}
	items := c.GetValue().Items
	if len(items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	// 绑定结算购物车
	this._cart = c
	// 将购物车的商品分类整理
	this._vendorItemsMap = this.buildVendorItemMap(items)
	// 更新订单的金额
	this._vendorExpressMap = this.updateOrderFee(this._vendorItemsMap)
	// 状态设为待支付
	this._value.Status = 1

	return nil
}

// 更新订单金额,并返回运费
func (this *orderImpl) updateOrderFee(mp map[int][]*order.OrderItem) map[int]float32 {
	this._value.GoodsFee = 0
	weightMap := make(map[int]int) //重量
	for k, v := range mp {
		weightMap[k] = 0
		for _, item := range v {
			//计算商品总金额
			this._value.GoodsFee += item.Fee
			//计算商品优惠金额
			this._value.DiscountFee += item.Fee - item.FinalFee
			//计重
			weightMap[k] += item.Weight
		}
	}
	// 计算运费
	expressMap := make(map[int]float32)
	for k, weight := range weightMap {
		//todo: 计算运费需从外部传入参数
		unit := weight / 1000 //转换为kg
		expressMap[k] = this._expressRep.GetUserExpress(k).
			GetExpressFee(-1, "1000", unit)
		//叠加运费
		this._value.ExpressFee += expressMap[k]
	}

	this._value.PackageFee = 0
	//计算最终金额
	this._value.FinalFee = this._value.GoodsFee - this._value.DiscountFee +
		this._value.ExpressFee + this._value.PackageFee
	return expressMap
}

// 根据运营商获取商品和运费信息,限未生成的订单
func (this *orderImpl) GetByVendor() (items map[int][]*order.OrderItem,
	expressFee map[int]float32) {
	if this._vendorItemsMap == nil {
		panic("订单尚未读取购物车!")
	}
	if this._vendorExpressMap == nil {
		panic("订单尚未计算金额")
	}
	return this._vendorItemsMap, this._vendorExpressMap
}

// 检查购物车
func (this *orderImpl) checkCart() error {
	if this._cart == nil || len(this._cart.GetValue().Items) == 0 {
		return cart.ErrEmptyShoppingCart
	}
	return this._cart.Check()
}

// 生成运营商与订单商品的映射
func (this *orderImpl) buildVendorItemMap(items []*cart.CartItem) map[int][]*order.OrderItem {
	mp := make(map[int][]*order.OrderItem)
	for _, v := range items {
		//必须勾选为结算
		if v.Checked == 1 {
			item := this.parseCartToOrderItem(v)
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
func (this *orderImpl) parseCartToOrderItem(c *cart.CartItem) *order.OrderItem {
	gs := this._saleRep.GetSale(c.VendorId).GoodsManager().CreateGoods(
		&goods.ValueGoods{Id: c.SkuId, SkuId: c.SkuId})
	snap := gs.SnapshotManager().GetLatestSaleSnapshot()
	if snap == nil {
		domain.HandleError(errors.New("商品快照生成失败："+
			strconv.Itoa(c.SkuId)), "domain")
		return nil
	}
	fee := c.SalePrice * float32(c.Quantity)
	return &order.OrderItem{
		Id:         0,
		VendorId:   c.VendorId,
		ShopId:     c.ShopId,
		SkuId:      c.SkuId,
		SnapshotId: snap.Id,
		Quantity:   c.Quantity,
		Fee:        fee,
		FinalFee:   fee,
		Weight:     c.Snapshot.Weight * c.Quantity, //计算重量
	}
}

// 提交订单，返回订单号。如有错误则返回
func (this *orderImpl) Submit() (string, error) {
	if this.GetAggregateRootId() != 0 {
		return "", errors.New("订单不允许重复提交")
	}
	if err := this.checkCart(); err != nil {
		return "", err
	}
	buyer := this.GetBuyer()
	if buyer == nil {
		return "", member.ErrNoSuchMember
	}

	v := this._value

	//todo: best promotion , 优惠券和返现这里需要重构,直接影响到订单金额
	//prom,fee,integral := this.GetBestSavePromotion()

	// 应用优惠券
	if err := this.applyCouponOnSubmit(v); err != nil {
		return "", err
	}

	// 判断商品的优惠促销,如返现等
	proms, fee := this.applyCartPromotionOnSubmit(v, this._cart)
	if len(proms) != 0 {
		v.DiscountFee += float32(fee)
		v.FinalFee = v.GoodsFee - v.DiscountFee
		if v.FinalFee < 0 {
			// 如果出现优惠券多余的金额也一并使用
			v.FinalFee = 0
		}
	}

	this.avgDiscountToItem()

	// 检查是否已支付完成
	this.checkNewOrderPayment()

	// 保存订单
	orderId, err := this.saveNewOrderOnSubmit()
	v.Id = orderId
	if err == nil {
		// 绑定优惠券促销
		this.bindCouponOnSubmit(v.OrderNo)
		// 绑定购物车商品的促销
		for _, p := range proms {
			this.bindPromotionOnSubmit(v.OrderNo, p)
		}
		// 扣除库存
		this.applyGoodsNum()
		// 拆单
		this.breakUpByVendor()

		// 记录余额支付记录
		//todo: 扣减余额
		//if v.BalanceDiscount > 0 {
		//	err = acc.PaymentDiscount(v.OrderNo, v.BalanceDiscount)
		//}
	}
	return v.OrderNo, err
}

// 平均优惠抵扣金额到商品
func (this *orderImpl) avgDiscountToItem() {
	if this._vendorItemsMap == nil {
		panic(errors.New("仅能在下单时进行商品抵扣均分"))
	}
	if this._value.DiscountFee > 0 {
		totalFee := this._value.GoodsFee
		disFee := this._value.DiscountFee
		for _, items := range this._vendorItemsMap {
			for _, v := range items {
				v.FinalFee = v.Fee - (v.Fee/totalFee)*disFee
			}
		}
	}
}

// 绑定促销优惠
func (this *orderImpl) bindPromotionOnSubmit(orderNo string,
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
		OrderId:         this.GetAggregateRootId(),
		Title:           title,
		SaveFee:         float32(fee),
		PresentIntegral: integral,
		IsConfirm:       1,
		IsApply:         0,
	}
	return this._orderRep.SavePromotionBindForOrder(v)
}

// 应用购物车内商品的促销
func (this *orderImpl) applyCartPromotionOnSubmit(vo *order.Order,
	cart cart.ICart) ([]promotion.IPromotion, int) {
	var proms []promotion.IPromotion = make([]promotion.IPromotion, 0)
	var prom promotion.IPromotion
	var saveFee int
	var totalSaveFee int
	var intOrderFee = int(vo.FinalFee)
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
func (this *orderImpl) bindCouponOnSubmit(orderNo string) {
	var oc *order.OrderCoupon = new(order.OrderCoupon)
	for _, c := range this.GetCoupons() {
		oc.Clone(c, this.GetAggregateRootId(), this._value.FinalFee)
		this._orderRep.SaveOrderCouponBind(oc)
		// 绑定促销
		this.bindPromotionOnSubmit(orderNo, c.(promotion.IPromotion))
	}
}

// 在提交订单时应用优惠券
func (this *orderImpl) applyCouponOnSubmit(v *order.Order) error {
	var err error
	var t *promotion.ValueCouponTake
	var b *promotion.ValueCouponBind
	for _, c := range this.GetCoupons() {
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
func (this *orderImpl) getBalanceDiscountFee(acc member.IAccount) float32 {
	if this._value.FinalFee <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if acv.Balance >= this._value.FinalFee {
		return this._value.FinalFee
	} else {
		return acv.Balance
	}
	return 0
}

// 检查新订单的支付结果,如果最终付款为0,则设置为已支付
// 有可能为多余的, 应等到支付单支付完成后,再通知订单支付完成。
func (this *orderImpl) checkNewOrderPayment() {
	// 校验是否支付
	if this._value.FinalFee == 0 {
		this._value.IsPaid = 1
		this._value.PaidTime = time.Now().Unix()
	}
	// 设置订单状态
	if this._value.IsPaid == 1 {
		//todo:  线下支付应设为等待确认
		//|| v.PaymentOpt == enum.PaymentOfflineCashPay ||
		//v.PaymentOpt == enum.PaymentRemit {
		//v.PaymentSign = 1
		this._value.Status = enum.ORDER_WAIT_CONFIRM
	} else {
		this._value.Status = enum.ORDER_WAIT_PAYMENT
	}
}

// 保存订单
func (this *orderImpl) saveNewOrderOnSubmit() (int, error) {
	unix := time.Now().Unix()
	this._value.ItemsInfo = string(this._cart.GetJsonItems())
	this._value.OrderNo = this._manager.GetFreeOrderNo(0)
	this._value.CreateTime = unix
	this._value.UpdateTime = unix

	id, err := this._orderRep.SaveOrder(this._value)
	if err == nil {
		this._value.Id = id
		// 释放购物车并销毁
		if this._cart.Release() {
			this._cart.Destroy()
		}
	}
	return id, err
}

// 保存订单
func (this *orderImpl) Save() (int, error) {
	// 有操作后解除挂起状态
	// todo: ???
	//if this._value.IsSuspend == 1 && !this._internalSuspend {
	//    this._value.IsSuspend = 0
	//}

	if this._value.Id > 0 {
		return this._orderRep.SaveOrder(this._value)
	}
	this._internalSuspend = false
	return 0, errors.New("please use Order.Submit() save new order.")
}

// 根据运营商生成子订单
func (this *orderImpl) createSubOrderByVendor(parentOrderId int,
	vendorId int, newOrderNo bool, items []*order.OrderItem) order.ISubOrder {
	orderNo := this.GetOrderNo()
	if newOrderNo {
		orderNo = this._manager.GetFreeOrderNo(vendorId)
	}

	if len(items) == 0 {
		domain.HandleError(errors.New("拆分订单,运营商下未获取到商品,订单:"+
			this.GetOrderNo()), "domain")
		return nil
	}

	v := &order.SubOrder{
		OrderNo:   orderNo,
		VendorId:  vendorId,
		ParentId:  parentOrderId,
		Subject:   "子订单",
		ShopId:    items[0].ShopId,
		ItemsInfo: "",
		// 总金额
		GoodsFee: 0,
		// 减免金额(包含优惠券金额)
		DiscountFee: 0,
		ExpressFee:  0,
		FinalFee:    0,
		// 是否挂起，如遇到无法自动进行的时挂起，来提示人工确认。
		IsSuspend:  0,
		Note:       "",
		Remark:     "",
		Status:     enum.ORDER_WAIT_PAYMENT,
		UpdateTime: this._value.UpdateTime,
		Items:      items,
	}
	// 计算订单金额
	for _, item := range items {
		//计算商品金额
		v.GoodsFee += item.Fee
		//计算商品优惠金额
		v.DiscountFee += item.Fee - item.FinalFee
	}
	// 设置运费
	v.ExpressFee = this._vendorExpressMap[vendorId]
	// 设置包装费
	v.PackageFee = 0
	// 最终金额 = 商品金额 - 商品抵扣金额(促销折扣) + 包装费 + 快递费
	v.FinalFee = v.GoodsFee - v.DiscountFee + v.PackageFee + v.ExpressFee
	// 判断是否已支付
	if this._value.IsPaid == 1 {
		v.Status = enum.ORDER_WAIT_CONFIRM
	}
	return this._manager.CreateSubOrder(v)
}

//根据运营商拆单,返回拆单结果,及拆分的订单数组
func (this *orderImpl) breakUpByVendor() []order.ISubOrder {
	parentOrderId := this.GetAggregateRootId()
	if parentOrderId <= 0 ||
		this._vendorItemsMap == nil ||
		len(this._vendorItemsMap) == 0 {
		//todo: 订单要取消掉
		panic(fmt.Sprintf("订单异常: 订单未生成或VendorItemMap为空,"+
			"订单编号:%d,订单号:%s,vendor len:%d",
			parentOrderId, this._value.OrderNo, len(this._vendorItemsMap)))
	}
	l := len(this._vendorItemsMap)
	list := make([]order.ISubOrder, l)
	i := 0
	for k, v := range this._vendorItemsMap {
		//log.Println("----- vendor ", k, len(v),l)
		list[i] = this.createSubOrderByVendor(parentOrderId, k, l > 1, v)
		if _, err := list[i].Save(); err != nil {
			domain.HandleError(err, "domain")
		}
		i++
	}
	return list
}

// 扣除库存
func (this *orderImpl) applyGoodsNum() {
	for _, v := range this._vendorItemsMap {
		for _, v2 := range v {
			this.addGoodsSaleNum(v2.VendorId, v2.SkuId, v2.Quantity)
		}
	}
}

//****************  订单提交结束 **************//

// 使用余额支付
func (this *orderImpl) paymentWithBalance(buyerType int) error {
	if this._value.IsPaid == 1 {
		return order.ErrOrderPayed
	}
	acc := this._memberRep.GetMember(this._value.BuyerId).GetAccount()
	if fee := this.getBalanceDiscountFee(acc); fee == 0 {
		return member.ErrAccountBalanceNotEnough
	} else {
		this._value.DiscountFee = fee
		this._value.FinalFee -= fee
		err := acc.PaymentDiscount(this.GetOrderNo(), fee)
		if err != nil {
			return err
		}
	}
	unix := time.Now().Unix()
	if this._value.FinalFee == 0 {
		this._value.IsPaid = 1
		// this._value.PaymentSign = buyerType
		if this._value.Status == enum.ORDER_WAIT_PAYMENT {
			this._value.Status = enum.ORDER_WAIT_CONFIRM
		}
	}
	this._value.UpdateTime = unix
	this._value.PaidTime = unix

	_, err := this.Save()
	return err
}

// 使用余额支付
func (this *orderImpl) PaymentWithBalance() error {
	return this.paymentWithBalance(payment.PaymentByBuyer)
}

// 客服使用余额支付
func (this *orderImpl) CmPaymentWithBalance() error {
	return this.paymentWithBalance(payment.PaymentByCM)
}

// 添加日志
func (this *orderImpl) AppendLog(t enum.OrderLogType, system bool, message string) error {
	if this.GetAggregateRootId() <= 0 {
		return errors.New("order not created.")
	}

	var systemInt int
	if system {
		systemInt = 1
	} else {
		systemInt = 0
	}

	var ol *order.OrderLog = &order.OrderLog{
		OrderId:    this.GetAggregateRootId(),
		Type:       int(t),
		IsSystem:   systemInt,
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return this._orderRep.SaveOrderLog(ol)
}

// 订单是否已完成
func (this *orderImpl) IsOver() bool {
	s := this._value.Status
	return s == enum.ORDER_CANCEL || s == enum.ORDER_COMPLETED
}

// 处理订单
func (this *orderImpl) Process() error {
	dt := time.Now()
	this._value.Status += 1
	this._value.UpdateTime = dt.Unix()

	_, err := this.Save()
	return err
}

// 确认订单
func (this *orderImpl) Confirm() error {
	panic("not implement")
	//if this._value.PaymentOpt == enum.PaymentOnlinePay &&
	//this._value.IsPaid == enum.FALSE {
	//    return order.ErrOrderNotPayed
	//}
	//if this._value.Status == enum.ORDER_WAIT_CONFIRM {
	//    this._value.Status = enum.ORDER_WAIT_DELIVERY
	//    this._value.UpdateTime = time.Now().Unix()
	//    _, err := this.Save()
	//    if err == nil {
	//        err = this.AppendLog(enum.ORDER_LOG_SETUP, false, "订单已经确认")
	//    }
	//    return err
	//}
	return nil
}

// 添加商品销售数量
func (this *orderImpl) addGoodsSaleNum(vendorId, skuId, quantity int) error {
	gds := this._saleRep.GetSale(vendorId).
		GoodsManager().GetGoods(skuId)
	if gds == nil {
		return goods.ErrNoSuchGoods
	}
	return gds.AddSaleNum(quantity)
}

// 配送订单
func (this *orderImpl) Deliver(spId int, spNo string) error {
	//todo: 记录快递配送信息
	dt := time.Now()
	this._value.Status += 1
	this._value.ShippingTime = dt.Unix()
	this._value.UpdateTime = dt.Unix()

	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, false, "订单开始配送")
	}
	return err
}

// 获取订单号
func (this *orderImpl) GetOrderNo() string {
	return this.GetValue().OrderNo
}

func (this *orderImpl) backupPayment() error {
	if this._value.DiscountFee > 0 {
		//退回账户余额抵扣
		acc := this._memberRep.GetMember(this._value.BuyerId).GetAccount()
		return acc.ChargeBalance(member.TypeBalanceOrderRefund, "订单退款",
			this.GetOrderNo(),
			this._value.DiscountFee)
	}
	if this._value.FinalFee > 0 {
		//todo: 其他支付方式退还,如网银???
	}
	return nil
}

// 更新返现到会员账户
func (this *orderImpl) updateShoppingMemberBackFee(pt merchant.IMerchant,
	m member.IMember, fee float32, unixTime int64) {
	if fee == 0 {
		return
	}
	v := this.GetValue()
	pv := pt.GetValue()

	//更新账户
	acc := m.GetAccount()
	acv := acc.GetValue()
	//acc.TotalFee += this._value.Fee
	//acc.TotalPay += this._value.PayFee
	acv.PresentBalance += fee // 更新赠送余额
	acv.TotalPresentFee += fee
	acv.UpdateTime = unixTime
	acc.Save()

	//给自己返现
	tit := fmt.Sprintf("订单:%s(商户:%s)返现￥%.2f元", v.OrderNo, pv.Name, fee)
	acc.PresentBalance(tit, v.OrderNo, float32(fee))
}

// 处理返现促销
func (this *orderImpl) handleCashBackPromotions(pt merchant.IMerchant,
	m member.IMember) error {
	proms := this.GetPromotionBinds()
	for _, v := range proms {
		if v.PromotionType == promotion.TypeFlagCashBack {
			c := this._promRep.GetPromotion(v.PromotionId)
			return this.handleCashBackPromotion(pt, m, v, c)
		}
	}
	return nil
}

// 处理返现促销
func (this *orderImpl) handleCashBackPromotion(pt merchant.IMerchant,
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
		this._orderRep.SavePromotionBindForOrder(v)

		// 处理自定义返现
		c := pm.(promotion.ICashBackPromotion)
		HandleCashBackDataTag(m, this._value, c, this._memberRep)

		//给自己返现
		tit := fmt.Sprintf("返现￥%d元,订单号:%s", cpv.BackFee, this._value.OrderNo)
		err = acc.PresentBalance(tit, this.GetOrderNo(), float32(cpv.BackFee))
	}
	return err
}

// 三级返现
func (this *orderImpl) backFor3R(mch merchant.IMerchant, m member.IMember,
	back_fee float32, unixTime int64) {
	if back_fee == 0 {
		return
	}

	i := 0
	mName := m.Profile().GetProfile().Name
	saleConf := mch.ConfManager().GetSaleConf()
	percent := saleConf.CashBackTg2Percent
	for i < 2 {
		rl := m.GetRelation()
		if rl == nil || rl.RefereesId == 0 {
			break
		}

		m = this._memberRep.GetMember(rl.RefereesId)
		if m == nil {
			break
		}

		if i == 1 {
			percent = saleConf.CashBackTg1Percent
		}

		this.updateMemberAccount(m, mch.GetValue().Name, mName,
			back_fee*percent, unixTime)
		i++
	}
}

func (this *orderImpl) updateMemberAccount(m member.IMember,
	ptName, mName string, fee float32, unixTime int64) {
	if fee == 0 {
		return
	}

	//更新账户
	acc := m.GetAccount()
	acv := acc.GetValue()
	acv.PresentBalance += fee
	acv.TotalPresentFee += fee
	acv.UpdateTime = unixTime
	acc.Save()

	//给自己返现
	tit := fmt.Sprintf("订单:%s(商户:%s,会员:%s)收入￥%.2f元",
		this._value.OrderNo, ptName, mName, fee)
	acc.PresentBalance(tit, this._value.OrderNo, fee)
}

var _ order.ISubOrder = new(subOrderImpl)

// 子订单实现
type subOrderImpl struct {
	_value           *order.SubOrder
	_internalSuspend bool //内部挂起
	_rep             order.IOrderRep
	_memberRep       member.IMemberRep
	_goodsRep        goods.IGoodsRep
	_saleRep         sale.ISaleRep
	_manager         order.IOrderManager
}

func NewSubOrder(v *order.SubOrder,
	manager order.IOrderManager, rep order.IOrderRep,
	mmRep member.IMemberRep, goodsRep goods.IGoodsRep,
	saleRep sale.ISaleRep) order.ISubOrder {
	return &subOrderImpl{
		_value:     v,
		_manager:   manager,
		_rep:       rep,
		_memberRep: mmRep,
		_goodsRep:  goodsRep,
		_saleRep:   saleRep,
	}
}

// 获取领域对象编号
func (this *subOrderImpl) GetDomainId() int {
	return this._value.Id
}

// 获取值对象
func (this *subOrderImpl) GetValue() *order.SubOrder {
	return this._value
}

// 添加备注
func (this *subOrderImpl) AddRemark(remark string) {
	this._value.Remark = remark
}

// 设置Shop
func (this *subOrderImpl) SetShop(shopId int) error {
	//todo:验证Shop
	this._value.ShopId = shopId
	if this._value.Status == enum.ORDER_WAIT_CONFIRM {
		panic("not impl")
		// this.Confirm()
	}
	return nil
}

// 保存订单
func (this *subOrderImpl) Save() (int, error) {
	if this.GetDomainId() > 0 {
		return this._rep.SaveSubOrder(this._value)
	}
	id, err := this._rep.SaveSubOrder(this._value)
	if err == nil {
		this._value.Id = id
		unix := time.Now().Unix()
		for _, v := range this._value.Items {
			v.OrderId = id
			v.UpdateTime = unix
			this._rep.SaveOrderItem(id, v)
		}
	}
	return id, err
}

// 挂起
func (this *subOrderImpl) Suspend(reason string) error {
	this._value.IsSuspend = 1
	this._internalSuspend = true
	this._value.UpdateTime = time.Now().Unix()
	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, true, "订单已锁定"+reason)
	}
	return err
}

// 添加日志
func (this *subOrderImpl) AppendLog(t enum.OrderLogType,
	system bool, message string) error {
	if this.GetDomainId() <= 0 {
		return errors.New("order not created.")
	}
	var systemInt int
	if system {
		systemInt = 1
	} else {
		systemInt = 0
	}
	var ol *order.OrderLog = &order.OrderLog{
		OrderId:    this.GetDomainId(),
		Type:       int(t),
		IsSystem:   systemInt,
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return this._rep.SaveOrderLog(ol)
}

// 标记收货
func (this *subOrderImpl) SignReceived() error {
	dt := time.Now()
	this._value.Status = enum.ORDER_RECEIVED
	this._value.UpdateTime = dt.Unix()

	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, false, "已收货")
	}
	return err
}

// 更新账户
func (this *subOrderImpl) updateAccountForOrder(m member.IMember,
	order order.ISubOrder) {
	acc := m.GetAccount()
	ov := order.GetValue()
	acv := acc.GetValue()
	acv.TotalFee += ov.GoodsFee
	acv.TotalPay += ov.FinalFee
	acv.UpdateTime = time.Now().Unix()
	acc.Save()
}

// 完成订单
func (this *subOrderImpl) Complete() error {
	// now := time.Now().Unix()
	// v := this._value

	po := this._rep.Manager().GetOrderById(this._value.ParentId)
	pv := po.GetValue()
	m := this._memberRep.GetMember(pv.BuyerId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	var err error

	//todo: ???

	//var mch merchant.IMerchant
	////mch, err = this._partnerRep.GetMerchant(v.VendorId)
	//pv := mch.GetValue()
	//if pv.ExpiresTime < time.Now().Unix() {
	//    return errors.New("您的账户已经过期!")
	//}
	//
	//if err != nil {
	//    log.Println("供应商异常!", v.VendorId)
	//    return err
	//}

	// 增加经验
	if EXP_BIT == 0 {
		fv := infrastructure.GetApp().Config().GetFloat(variable.EXP_BIT)
		if fv <= 0 {
			panic("[WANNING]:Exp_bit not set!")
		}
		EXP_BIT = float32(fv)
	}
	if err = m.AddExp(int(pv.FinalFee * EXP_BIT)); err != nil {
		return err
	}

	// 更新账户
	this.updateAccountForOrder(m, this)

	panic("not implement")

	//todo: 获取商户的销售比例

	//******* 返现到账户  ************
	//var back_fee float32
	//saleConf := mch.ConfManager().GetSaleConf()
	//globSaleConf := this._valRep.GetGlobNumberConf()
	//if saleConf.CashBackPercent > 0 {
	//    back_fee = v.Fee * saleConf.CashBackPercent
	//
	//    //将此次消费记入会员账户
	//    this.updateShoppingMemberBackFee(mch, m,
	//        back_fee * saleConf.CashBackMemberPercent, now)
	//
	//    //todo: 增加阶梯的返积分,比如订单满30送100积分
	//    backIntegral := int(v.Fee) * globSaleConf.IntegralBackNum +
	//    globSaleConf.IntegralBackExtra
	//
	//    // 赠送积分
	//    if backIntegral != 0 {
	//        err = m.GetAccount().AddIntegral(v.VendorId, enum.INTEGRAL_TYPE_ORDER,
	//            backIntegral, fmt.Sprintf("订单返积分%d个", backIntegral))
	//        if err != nil {
	//            return err
	//        }
	//    }
	//}
	//
	//this._value.Status = enum.ORDER_COMPLETED
	//this._value.IsSuspend = 0
	//this._value.UpdateTime = now
	//
	//_, err = this.Save()
	//
	//if err == nil {
	//    err = this.AppendLog(enum.ORDER_LOG_SETUP, false, "订单已完成")
	//    // 处理返现促销
	//    this.handleCashBackPromotions(mch, m)
	//    // 三级返现
	//    if back_fee > 0 {
	//        this.backFor3R(mch, m, back_fee, now)
	//    }
	//}
	//return err
}

// 取消商品
func (this *subOrderImpl) cancelGoods() error {
	for _, v := range this._value.Items {
		snapshot := this._goodsRep.GetSaleSnapshot(v.SnapshotId)
		if snapshot == nil {
			return goods.ErrNoSuchSnapshot
		}
		var gds sale.IGoods = this._saleRep.GetSale(this._value.VendorId).
			GoodsManager().GetGoods(snapshot.SkuId)
		if gds != nil {
			gds.CancelSale(v.Quantity, this._value.OrderNo)
		}
	}
	return nil
}

// 取消订单
func (this *subOrderImpl) Cancel(reason string) error {
	if len(strings.TrimSpace(reason)) == 0 {
		return errors.New("取消原因不能为空")
	}
	status := this._value.Status
	if status == enum.ORDER_COMPLETED {
		return errors.New("订单已经完成!")
	}
	if status == enum.ORDER_CANCEL {
		return errors.New("订单已经被取消!")
	}

	this._value.Status = enum.ORDER_CANCEL
	this._value.UpdateTime = time.Now().Unix()

	//todo: 应同时取消支付单

	this.cancelGoods()
	//todo: 退款

	//this.backupPayment()

	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, true, "订单已取消,原因："+reason)
	}

	return err
}
