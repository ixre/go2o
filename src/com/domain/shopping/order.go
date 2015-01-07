/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-09 15:03
 * description :
 * history :
 */

package shopping

import (
	"com/domain/interface/enum"
	"com/domain/interface/member"
	"com/domain/interface/partner"
	"com/domain/interface/promotion"
	"com/domain/interface/shopping"
	"com/infrastructure"
	"com/infrastructure/log"
	"com/share/variable"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	EXP_BIT float32
)

type Order struct {
	_shopping   shopping.IShopping
	value       *shopping.ValueOrder
	cart        shopping.ICart
	coupons     []promotion.ICoupon
	memberRep   member.IMemberRep
	shoppingRep shopping.IShoppingRep
	partnerRep  partner.IPartnerRep
}

func newOrder(shopping shopping.IShopping, value *shopping.ValueOrder, cart shopping.ICart,
	partnerRep partner.IPartnerRep, shoppingRep shopping.IShoppingRep,
	memberRep member.IMemberRep) shopping.IOrder {
	return &Order{
		_shopping:   shopping,
		value:       value,
		cart:        cart,
		memberRep:   memberRep,
		shoppingRep: shoppingRep,
		partnerRep:  partnerRep,
	}
}

func (this *Order) GetDomainId() int {
	return this.value.Id
}

func (this *Order) GetValue() shopping.ValueOrder {
	return *this.value
}

func (this *Order) ApplyCoupon(coupon promotion.ICoupon) error {
	if this.coupons == nil {
		this.coupons = []promotion.ICoupon{}
	}
	this.coupons = append(this.coupons, coupon)

	//val := coupon.GetValue()
	v := this.value
	//v.CouponCode = val.Code
	//v.CouponDescribe = coupon.GetDescribe()
	v.CouponFee = coupon.GetCouponFee(v.Fee)
	v.PayFee = v.Fee - v.CouponFee
	v.DiscountFee = v.DiscountFee + v.CouponFee
	return nil
}

// 获取应用的优惠券
func (this *Order) GetCoupons() []promotion.ICoupon {
	if this.coupons == nil {
		return make([]promotion.ICoupon, 0)
	}
	return this.coupons
}

// 添加备注
func (this *Order) AddRemark(remark string) {
	this.value.Note = remark
}

// 设置Shop
func (this *Order) SetShop(shopId int) error {
	//todo:验证Shop
	this.value.ShopId = shopId
	return nil
}

// 设置支付方式
func (this *Order) SetPayment(payment int) {
	this.value.PayMethod = payment
}

// 设置配送地址
func (this *Order) SetDeliver(deliverAddrId int) error {
	d := this.memberRep.GetDeliverAddr(this.value.MemberId, deliverAddrId)
	if d != nil {
		v := this.value
		v.DeliverAddress = d.Address
		v.DeliverName = d.RealName
		v.DeliverPhone = d.Phone
		v.DeliverTime = time.Now().Add(-time.Hour)
		return nil
	}
	return errors.New("Deliver address not exist!")
}

// 提交订单，返回订单号。如有错误则返回
func (this *Order) Submit() (string, error) {
	if this.GetDomainId() != 0 {
		return "", errors.New("订单不允许重复提交！")
	}

	if this.cart == nil || len(this.cart.GetValue().Items) == 0 {
		return "", errors.New("购物车为空！")
	}

	v := this.value
	v.CreateTime = time.Now()
	v.UpdateTime = v.CreateTime
	v.ItemsInfo = this.cart.GetSummary()
	v.OrderNo = this._shopping.GetFreeOrderNo()

	// 应用优惠券
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
			log.PrintErr(err)
			err = errors.New("Code 105:优惠券使用失败")
			return "", err
		}
	}

	// 保存订单
	id, err := this.shoppingRep.SaveOrder(this._shopping.GetAggregateRootId(), v)
	v.Id = id
	if err == nil {
		var oc *shopping.OrderCoupon = new(shopping.OrderCoupon)
		for _, c := range this.GetCoupons() {
			oc.Clone(c, v.Id, v.Fee)
			// 绑定订单与优惠券
			this.shoppingRep.SaveOrderCouponBind(oc)
		}
	}

	return v.OrderNo, err
}

// 保存订单
func (this *Order) Save() error {
	_, err := this.shoppingRep.SaveOrder(
		this._shopping.GetAggregateRootId(), this.value)
	return err
}

// 订单是否已完成
func (this *Order) IsOver() bool {
	s := this.value.Status
	return s == enum.ORDER_CANCEL || s == enum.ORDER_COMPLETED
}

// 处理订单
func (this *Order) Process() error {
	dt := time.Now()
	this.value.Status += 1
	this.value.UpdateTime = dt
	return this.Save()
}

// 确认订单
func (this *Order) Confirm() error {
	this.value.Status = enum.ORDER_CONFIRMED
	this.value.UpdateTime = time.Now()
	return this.Save()
}

// 配送订单
func (this *Order) Deliver() error {
	dt := time.Now()
	this.value.Status += 1
	this.value.DeliverTime = dt
	this.value.UpdateTime = dt
	return this.Save()
}

// 取消订单
func (this *Order) Cancel(reason string) error {
	if len(strings.TrimSpace(reason)) == 0 {
		return errors.New("取消原因不能为空")
	}
	status := this.value.Status
	if status == enum.ORDER_COMPLETED {
		return errors.New("订单已经完成!")
	}
	if status == enum.ORDER_CANCEL {
		return errors.New("订单已经被取消!")
	}

	this.value.Status = enum.ORDER_CANCEL
	this.value.Remark += "取消原因:" + reason
	return this.Save()
}

// 完成订单
func (this *Order) Complete() error {
	now := time.Now()
	v := this.value
	m, err := this.memberRep.GetMember(v.MemberId)
	if err == nil {
		var ptl partner.IPartner
		ptl = this.partnerRep.GetPartner(v.PartnerId)
		if ptl == nil {
			log.Println("供应商不存在!", v.PartnerId)
			return errors.New("供应商不存在!")
		}
		pv := ptl.GetValue()
		if pv.Expires.UTC().Unix() < time.Now().Unix() {
			return errors.New("您的账户已经过期!")
		}

		//返现比例
		saleConf := ptl.GetSaleConf()
		back_fee := v.Fee * saleConf.CashBackPercent

		//将此次消费记入会员账户
		this.updateShoppingMemberAccount(ptl, m,
			back_fee*saleConf.CashBackMemberPercent, now)

		//todo: 增加阶梯的返积分,比如订单满30送100积分
		backIntegral := int(v.Fee)*saleConf.IntegralBackNum +
			saleConf.IntegralBackExtra

		//判断是否满足升级条件
		if backIntegral != 0 {
			err = m.AddIntegral(v.PartnerId, enum.INTEGRAL_TYPE_ORDER,
				backIntegral, fmt.Sprintf("订单返积分%d个", backIntegral))
			if err != nil {
				return err
			}
		}

		// 增加经验
		if EXP_BIT == 0 {
			fv := infrastructure.GetContext().
				Config().GetFloat(variable.EXP_BIT)
			EXP_BIT = float32(fv)
		}

		if EXP_BIT == 0 {
			log.Println("[WANNING]:Exp_bit not set!")
		}

		err = m.AddExp(int(v.Fee * EXP_BIT))
		if err != nil {
			return err
		}

		// 三级返现
		this.backFor3R(ptl, m, back_fee, now)

		this.value.Status = enum.ORDER_COMPLETED
		this.value.UpdateTime = now
		err = this.Save()
	}
	return err
}

// 更新会员账户
func (this *Order) updateShoppingMemberAccount(pt partner.IPartner,
	m member.IMember, fee float32, t time.Time) {
	if fee == 0 {
		return
	}
	v := this.GetValue()
	pv := pt.GetValue()
	//更新账户
	acc := m.GetAccount()
	acc.TotalFee = acc.TotalFee + this.value.Fee
	acc.TotalPay = acc.TotalPay + this.value.PayFee
	acc.PresentBalance = acc.PresentBalance + fee //更新赠送余额
	acc.UpdateTime = t
	m.SaveAccount()

	//给自己返现
	icLog := &member.IncomeLog{
		MemberId:   this.value.MemberId,
		OrderId:    v.Id,
		Type:       "backcash",
		Fee:        fee,
		Log:        fmt.Sprintf("订单:%s(商家:%s)返现￥%.2f元", v.OrderNo, pv.Name, fee),
		State:      1,
		RecordTime: t,
	}
	m.SaveIncomeLog(icLog)
}

// 三级返现
func (this *Order) backFor3R(pt partner.IPartner, m member.IMember,
	back_fee float32, now time.Time) {
	if back_fee == 0 {
		return
	}

	i := 0
	mName := m.GetValue().Name
	saleConf := pt.GetSaleConf()
	percent := saleConf.CashBackTg2Percent
	for i < 2 {
		rl := m.GetRelation()
		if rl == nil || rl.TgId == 0 {
			break
		}

		m, _ = this.memberRep.GetMember(rl.TgId)
		if m == nil {
			break
		}

		if i == 1 {
			percent = saleConf.CashBackTg1Percent
		}

		this.updateMemberAccount(m, pt.GetValue().Name, mName,
			back_fee*percent, now)
		i++
	}
}

func (this *Order) updateMemberAccount(m member.IMember,
	ptName, mName string, fee float32, t time.Time) {
	if fee == 0 {
		return
	}

	//更新账户
	acc := m.GetAccount()
	acc.PresentBalance = acc.PresentBalance + fee
	acc.UpdateTime = time.Now()
	m.SaveAccount()

	//给自己返现
	icLog := &member.IncomeLog{
		MemberId: this.value.MemberId,
		Type:     "backcash",
		Fee:      fee,
		Log: fmt.Sprintf("订单:%s(商家:%s,会员:%s)收入￥%.2f元",
			this.value.OrderNo, ptName, mName, fee),
		State:      1,
		RecordTime: t,
	}
	m.SaveIncomeLog(icLog)
}
