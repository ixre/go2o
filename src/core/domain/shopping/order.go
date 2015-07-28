 /**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-09 15:03
 * description :
 * history :
 */

package shopping

import (
	"errors"
	"fmt"
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/shopping"
	"go2o/src/core/infrastructure"
	"go2o/src/core/infrastructure/log"
	"go2o/src/core/variable"
	"strings"
	"time"
)

var (
	EXP_BIT float32
)

type Order struct {
	_shopping        shopping.IShopping
	_value           *shopping.ValueOrder
	_cart            shopping.ICart
	_coupons         []promotion.ICouponPromotion
	_availPromotions []promotion.IPromotion
	_orderPbs        []*shopping.OrderPromotionBind
	_memberRep       member.IMemberRep
	_shoppingRep     shopping.IShoppingRep
	_partnerRep      partner.IPartnerRep
	_saleRep         sale.ISaleRep
	_promRep         promotion.IPromotionRep
	_internalSuspend bool // 是否为内部挂起
	_balanceDiscount bool // 余额支付
}

func newOrder(shopping shopping.IShopping, value *shopping.ValueOrder, cart shopping.ICart,
	partnerRep partner.IPartnerRep, shoppingRep shopping.IShoppingRep, saleRep sale.ISaleRep,
	promRep promotion.IPromotionRep, memberRep member.IMemberRep) shopping.IOrder {
	return &Order{
		_shopping:    shopping,
		_value:       value,
		_cart:        cart,
		_memberRep:   memberRep,
		_promRep:     promRep,
		_shoppingRep: shoppingRep,
		_partnerRep:  partnerRep,
		_saleRep:     saleRep,
	}
}

func (this *Order) GetDomainId() int {
	return this._value.Id
}

func (this *Order) GetValue() shopping.ValueOrder {
	return *this._value
}

// 应用优惠券
func (this *Order) ApplyCoupon(coupon promotion.ICouponPromotion) error {
	if this._coupons == nil {
		this._coupons = []promotion.ICouponPromotion{}
	}
	this._coupons = append(this._coupons, coupon)

	//val := coupon.GetValue()
	v := this._value
	//v.CouponCode = val.Code
	//v.CouponDescribe = coupon.GetDescribe()
	v.CouponFee = coupon.GetCouponFee(v.Fee)
	v.PayFee = this.GetPaymentFee()
	v.DiscountFee = v.DiscountFee + v.CouponFee
	return nil
}

// 获取支付金额
func (this *Order) GetPaymentFee() float32 {
	return this._value.TotalFee - this._value.CouponFee
}

// 获取应用的优惠券
func (this *Order) GetCoupons() []promotion.ICouponPromotion {
	if this._coupons == nil {
		return make([]promotion.ICouponPromotion, 0)
	}
	return this._coupons
}

// 获取可用的促销,不包含优惠券
func (this *Order) GetAvailableOrderPromotions() []promotion.IPromotion {
	if this._availPromotions == nil {
		partnerId := this._value.PartnerId
		var vp []*promotion.ValuePromotion = this._promRep.GetPromotionOfPartnerOrder(partnerId)
		var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
		for i, v := range vp {
			proms[i] = this._promRep.CreatePromotion(v)
		}
		return proms
	}
	return this._availPromotions
}

// 获取促销绑定
func (this *Order) GetPromotionBinds() []*shopping.OrderPromotionBind {
	if this._orderPbs == nil {
		this._orderPbs = this._shoppingRep.GetOrderPromotionBinds(this._value.OrderNo)
	}
	return this._orderPbs
}

// 获取最省的促销
func (this *Order) GetBestSavePromotion() (p promotion.IPromotion, saveFee float32, integral int) {
	//todo: not implement
	return nil, 0, 0
}

// 添加备注
func (this *Order) AddRemark(remark string) {
	this._value.Note = remark
}

// 设置Shop
func (this *Order) SetShop(shopId int) error {
	//todo:验证Shop
	this._value.ShopId = shopId
	return nil
}

// 设置支付方式
func (this *Order) SetPayment(payment int) {
	this._value.PaymentOpt = payment
}

// 使用余额支付
func (this *Order) PaymentWithBalance() error {
	if this._value.IsPaid == 1 {
		return shopping.ErrOrderPayed
	}
	acc := this._memberRep.GetMember(this._value.MemberId).GetAccount()
	if fee := this.getBalanceDiscountFee(acc); fee == 0 {
		return shopping.ErrBalanceNotEnough
	} else {
		this._value.BalanceDiscount = fee
		this._value.PayFee -= fee
		err := acc.OrderDiscount(this._value.OrderNo, fee)
		if err != nil {
			return err
		}
	}
	unix := time.Now().Unix()
	if this._value.PayFee == 0 {
		this._value.IsPaid = 1
	}
	this._value.UpdateTime = unix
	this._value.PaidTime = unix

	_, err := this.Save()
	return err
}

// 在线交易支付
func (this *Order) PaymentOnlineTrade(serverProvider string, tradeNo string) error {
	if this._value.IsPaid == 1 {
		return shopping.ErrOrderPayed
	}
	unix := time.Now().Unix()
	this._value.IsPaid = 1
	this._value.UpdateTime = unix
	this._value.PaidTime = unix

	_, err := this.Save()
	return err
}

// 设置配送地址
func (this *Order) SetDeliver(deliverAddressId int) error {
	d := this._memberRep.GetSingleDeliverAddress(this._value.MemberId, deliverAddressId)
	if d != nil {
		v := this._value
		v.DeliverAddress = d.Address
		v.DeliverName = d.RealName
		v.DeliverPhone = d.Phone
		v.DeliverTime = time.Now().Add(-time.Hour).Unix()
		return nil
	}
	return errors.New("Deliver address not exist!")
}

// 使用余额支付
func (this *Order) UseBalanceDiscount() {
	this._balanceDiscount = true
}

// 提交订单，返回订单号。如有错误则返回
func (this *Order) Submit() (string, error) {
	if this.GetDomainId() != 0 {
		return "", errors.New("订单不允许重复提交！")
	}

	if this._cart == nil || len(this._cart.GetValue().Items) == 0 {
		return "", errors.New("购物车为空！")
	}

	mem := this._memberRep.GetMember(this._value.MemberId)
	if mem == nil {
		return "", member.ErrNoSuchMember
	}
	acc := mem.GetAccount()

	v := this._value
	v.CreateTime = time.Now().Unix()
	v.UpdateTime = v.CreateTime
	v.ItemsInfo = this._cart.GetSummary()
	v.OrderNo = this._shopping.GetFreeOrderNo()

	// 应用优惠券
	if err := this.applyCouponOnSubmit(v); err != nil {
		return "", err
	}

	// 购物车商品
	proms, fee := this.applyCartPromotionObSubmit(v, this._cart)
	if len(proms) != 0 {
		v.DiscountFee += float32(fee)
		v.PayFee -= float32(fee)
		if v.PayFee < 0 {
			v.PayFee = 0
		}
	}

	//todo:
	//prom,fee,integral := this.GetBestSavePromotion()

	// 余额支付
	if this._balanceDiscount {
		if fee := this.getBalanceDiscountFee(acc); fee > 0 {
			v.PayFee -= fee
			v.BalanceDiscount = fee
		}
	}

	// 校验是否支付
	if v.PayFee == 0 {
		v.IsPaid = 1
	}

	// 保存订单
	id, err := this.saveOrderOnSubmit()
	v.Id = id
	if err == nil {
		// 绑定优惠券促销
		this.bindCouponOnSubmit(v.OrderNo)
		// 销毁购物车
		this._cart.Destroy()
		// 绑定购物车商品的促销
		for _, p := range proms {
			this.bindPromotionOnSubmit(v.OrderNo, p)
		}

		// 记录余额支付记录
		if v.BalanceDiscount > 0 {
			err = acc.OrderDiscount(v.OrderNo, v.BalanceDiscount)
		}
	}
	return v.OrderNo, err
}

func (this *Order) bindPromotionOnSubmit(orderNo string, prom promotion.IPromotion) (int, error) {
	var title string
	var integral int
	var fee int

	//todo: 需要重构,其他促销
	if prom.Type() == promotion.TypeFlagCashBack {
		fee = prom.GetRelationValue().(*promotion.ValueCashBack).BackFee
		title = prom.TypeName() + ":" + prom.GetValue().ShortName
	}

	v := &shopping.OrderPromotionBind{
		PromotionId:     prom.GetAggregateRootId(),
		PromotionType:   prom.Type(),
		OrderNo:         orderNo,
		Title:           title,
		SaveFee:         float32(fee),
		PresentIntegral: integral,
		IsConfirm:       1,
		IsApply:         0,
	}
	return this._shoppingRep.SavePromotionBindForOrder(v)
}

// 应用购物车内商品的促销
func (this *Order) applyCartPromotionObSubmit(vo *shopping.ValueOrder, cart shopping.ICart) ([]promotion.IPromotion, int) {
	var proms []promotion.IPromotion = make([]promotion.IPromotion, 0)
	var prom promotion.IPromotion
	var saveFee int
	var totalSaveFee int
	var intOrderFee = int(vo.Fee)
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
func (this *Order) bindCouponOnSubmit(orderNo string) {
	var oc *shopping.OrderCoupon = new(shopping.OrderCoupon)
	for _, c := range this.GetCoupons() {
		oc.Clone(c, this.GetDomainId(), this._value.Fee)
		this._shoppingRep.SaveOrderCouponBind(oc)

		// 绑定促销
		this.bindPromotionOnSubmit(orderNo, c.(promotion.IPromotion))
	}
}

// 在提交订单时应用优惠券
func (this *Order) applyCouponOnSubmit(v *shopping.ValueOrder) error {
	var err error
	var t *promotion.ValueCouponTake
	var b *promotion.ValueCouponBind
	for _, c := range this.GetCoupons() {
		if c.CanTake() {
			t, err = c.GetTake(v.MemberId)
			if err == nil {
				err = c.ApplyTake(t.Id)
			}
		} else {
			b, err = c.GetBind(v.MemberId)
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
func (this *Order) getBalanceDiscountFee(acc member.IAccount) float32 {
	if this._value.PayFee <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if acv.Balance >= this._value.PayFee {
		return this._value.PayFee
	} else {
		return acv.Balance
	}
	return 0
}

// 保存订单
func (this *Order) saveOrderOnSubmit() (int, error) {
	cartItems := this._cart.GetValue().Items
	if this._value.Items == nil {
		this._value.Items = make([]*shopping.OrderItem, len(cartItems))
	}
	var sl sale.ISale = this._saleRep.GetSale(this._value.PartnerId)
	var item sale.IItem
	var snap *sale.GoodsSnapshot
	for i, v := range cartItems {
		snap = sl.GetGoodsSnapshot(cartItems[i].SnapshotId)
		if snap == nil {
			return 0, errors.New("商品缺少快照：" + item.GetValue().Name)
		}

		this._value.Items[i] = &shopping.OrderItem{
			Id:         0,
			SnapshotId: snap.Id,
			Quantity:   v.Num,
			Sku:        "",
			Fee:        v.SalePrice * float32(v.Num),
		}
	}

	return this._shoppingRep.SaveOrder(this._shopping.GetAggregateRootId(), this._value)
}

// 保存订单
func (this *Order) Save() (int, error) {
	// 有操作后解除挂起状态
	if this._value.IsSuspend == 1 && !this._internalSuspend {
		this._value.IsSuspend = 0
	}

	if this._value.Id > 0 {
		return this._shoppingRep.SaveOrder(
			this._shopping.GetAggregateRootId(), this._value)
	}
	this._internalSuspend = false
	return 0, errors.New("please use Order.Submit() save new order.")
}

// 添加日志
func (this *Order) AppendLog(t enum.OrderLogType, system bool, message string) error {
	if this.GetDomainId() <= 0 {
		return errors.New("order not created.")
	}

	var systemInt int
	if system {
		systemInt = 1
	} else {
		systemInt = 0
	}

	var ol *shopping.OrderLog = &shopping.OrderLog{
		OrderId:    this.GetDomainId(),
		Type:       int(t),
		IsSystem:   systemInt,
		Message:    message,
		RecordTime: time.Now().Unix(),
	}
	return this._shoppingRep.SaveOrderLog(ol)
}

// 订单是否已完成
func (this *Order) IsOver() bool {
	s := this._value.Status
	return s == enum.ORDER_CANCEL || s == enum.ORDER_COMPLETED
}

// 处理订单
func (this *Order) Process() error {
	dt := time.Now()
	this._value.Status += 1
	this._value.UpdateTime = dt.Unix()

	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, false, "订单处理中")
	}
	return err
}

// 确认订单
func (this *Order) Confirm() error {
	this._value.Status = enum.ORDER_CONFIRMED
	this._value.UpdateTime = time.Now().Unix()

	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, false, "订单已经确认")
	}
	return err
}

// 配送订单
func (this *Order) Deliver() error {
	dt := time.Now()
	this._value.Status += 1
	this._value.DeliverTime = dt.Unix()
	this._value.UpdateTime = dt.Unix()

	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, false, "订单开始配送")
	}
	return err
}

// 取消订单
func (this *Order) Cancel(reason string) error {
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

	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, true, "订单已取消,原因："+reason)
	}

	return err
}

// 挂起
func (this *Order) Suspend(reason string) error {
	this._value.IsSuspend = 1
	this._internalSuspend = true
	this._value.UpdateTime = time.Now().Unix()
	_, err := this.Save()
	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, true, "订单已锁定"+reason)
	}
	return err
}

// 标记收货
func (this *Order) SignReceived() error {
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
func updateAccountForOrder(m member.IMember, order shopping.IOrder) {
	acc := m.GetAccount()
	ov := order.GetValue()
	acv := acc.GetValue()
	acv.TotalFee += ov.Fee
	acv.TotalPay += ov.PayFee
	acv.UpdateTime = time.Now().Unix()
	acc.Save()
}

// 完成订单
func (this *Order) Complete() error {
	now := time.Now().Unix()
	v := this._value
	m := this._memberRep.GetMember(v.MemberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	var err error
	var ptl partner.IPartner
	ptl, err = this._partnerRep.GetPartner(v.PartnerId)
	if err != nil {
		log.Println("供应商异常!", v.PartnerId)
		log.PrintErr(err)
		return err
	}

	pv := ptl.GetValue()
	if pv.ExpiresTime < time.Now().Unix() {
		return errors.New("您的账户已经过期!")
	}

	// 增加经验
	if EXP_BIT == 0 {
		fv := infrastructure.GetApp().Config().GetFloat(variable.EXP_BIT)
		if fv <= 0 {
			panic("[WANNING]:Exp_bit not set!")
		}
		EXP_BIT = float32(fv)
	}
	if err = m.AddExp(int(v.Fee * EXP_BIT)); err != nil {
		return err
	}

	// 更新账户
	updateAccountForOrder(m, this)

	//******* 返现到账户  ************
	var back_fee float32
	saleConf := ptl.GetSaleConf()
	if saleConf.CashBackPercent > 0 {
		back_fee = v.Fee * saleConf.CashBackPercent

		//将此次消费记入会员账户
		this.updateShoppingMemberBackFee(ptl, m,
			back_fee*saleConf.CashBackMemberPercent, now)

		//todo: 增加阶梯的返积分,比如订单满30送100积分
		backIntegral := int(v.Fee)*saleConf.IntegralBackNum +
			saleConf.IntegralBackExtra

		// 赠送积分
		if backIntegral != 0 {
			err = m.AddIntegral(v.PartnerId, enum.INTEGRAL_TYPE_ORDER,
				backIntegral, fmt.Sprintf("订单返积分%d个", backIntegral))
			if err != nil {
				return err
			}
		}
	}

	this._value.Status = enum.ORDER_COMPLETED
	this._value.IsSuspend = 0
	this._value.UpdateTime = now

	_, err = this.Save()

	if err == nil {
		err = this.AppendLog(enum.ORDER_LOG_SETUP, false, "订单已完成")
		// 处理返现促销
		this.handleCashBackPromotions(ptl, m)
		// 三级返现
		if back_fee > 0 {
			this.backFor3R(ptl, m, back_fee, now)
		}
	}
	return err
}

// 更新返现到会员账户
func (this *Order) updateShoppingMemberBackFee(pt partner.IPartner,
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
	icLog := &member.IncomeLog{
		MemberId:   this._value.MemberId,
		OrderId:    v.Id,
		Type:       "backcash",
		Fee:        fee,
		Log:        fmt.Sprintf("订单:%s(商家:%s)返现￥%.2f元", v.OrderNo, pv.Name, fee),
		State:      1,
		RecordTime: unixTime,
	}
	m.SaveIncomeLog(icLog)
}

// 处理返现促销
func (this *Order) handleCashBackPromotions(pt partner.IPartner, m member.IMember) error {
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
func (this *Order) handleCashBackPromotion(pt partner.IPartner, m member.IMember,
	v *shopping.OrderPromotionBind, pm promotion.IPromotion) error {
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
		this._shoppingRep.SavePromotionBindForOrder(v)

		// 处理自定义返现
		c := pm.(promotion.ICashBackPromotion)
		HandleCashBackDataTag(m, this._value, c, this._memberRep)

		//给自己返现
		icLog := &member.IncomeLog{
			MemberId:   this._value.MemberId,
			OrderId:    this.GetDomainId(),
			Type:       "backcash",
			Fee:        float32(cpv.BackFee),
			Log:        fmt.Sprintf("返现￥%d元,订单号:%s", cpv.BackFee, this._value.OrderNo),
			State:      1,
			RecordTime: acv.UpdateTime,
		}
		err = m.SaveIncomeLog(icLog)
	}
	return err
}

// 三级返现
func (this *Order) backFor3R(pt partner.IPartner, m member.IMember,
	back_fee float32, unixTime int64) {
	if back_fee == 0 {
		return
	}

	i := 0
	mName := m.GetValue().Name
	saleConf := pt.GetSaleConf()
	percent := saleConf.CashBackTg2Percent
	for i < 2 {
		rl := m.GetRelation()
		if rl == nil || rl.InvitationMemberId == 0 {
			break
		}

		m = this._memberRep.GetMember(rl.InvitationMemberId)
		if m == nil {
			break
		}

		if i == 1 {
			percent = saleConf.CashBackTg1Percent
		}

		this.updateMemberAccount(m, pt.GetValue().Name, mName,
			back_fee*percent, unixTime)
		i++
	}
}

func (this *Order) updateMemberAccount(m member.IMember,
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
	icLog := &member.IncomeLog{
		MemberId: this._value.MemberId,
		Type:     "backcash",
		Fee:      fee,
		Log: fmt.Sprintf("订单:%s(商家:%s,会员:%s)收入￥%.2f元",
			this._value.OrderNo, ptName, mName, fee),
		State:      1,
		RecordTime: unixTime,
	}
	m.SaveIncomeLog(icLog)
}
