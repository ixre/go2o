/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-02 21:34
 * description :
 * history :
 */

package promotion

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/atnet/gof/math"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/infrastructure/format"
	"strconv"
	"time"
)

// 优惠券,一张优惠券有数量，对应唯一的促销码。
// 优惠内容包含：送金额，送积分,订单折扣。仅在消费时有效。
// 使用需要达到最低金额和最低等级。
// 优惠券包含了开始时间和结束时间，超出时间则过期。
// 优惠券可以设置启用和停用
// 是否允许绑定，如果不绑定。则可以任意使用.只要有绑定和使用后，就不允许修改此属性。
type Coupon struct {
	value        *promotion.ValueCoupon
	promRep      promotion.IPromotionRep
	memberRep    member.IMemberRep
	takes        []promotion.ValueCouponTake
	binds        []promotion.ValueCouponBind
	takes_loaded bool
	binds_loaded bool
}

func newCoupon(v *promotion.ValueCoupon, promRep promotion.IPromotionRep,
	memberRep member.IMemberRep) promotion.ICoupon {

	cp := &Coupon{value: v,
		promRep:   promRep,
		memberRep: memberRep,
	}

	cp.releaseCoupon()
	return cp
}

func (this *Coupon) GetDomainId() int {
	return this.value.Id
}

// 释放优惠券
func (this *Coupon) releaseCoupon() {
	// 仅在会员通用情况下才存在占用
	if this.value.NeedBind == 0 &&
		this.value.TotalAmount != this.value.Amount {
		now := time.Now().Unix()
		oriAmount := this.value.Amount
		for _, take := range this.GetTakes() {
			// 未应用到订单，且释放时间小于当前时间，则释放
			if take.IsApply == 0 && now > take.ExtraTime {
				this.value.Amount = this.value.Amount + 1
			}
		}

		//保存新的可用数量
		if oriAmount != this.value.Amount {
			this.Save()
		}
	}
}

func (this *Coupon) GetValue() promotion.ValueCoupon {
	return *this.value
}

// 设置值
func (this *Coupon) SetValue(v *promotion.ValueCoupon) error {

	val := this.value
	if v.TotalAmount < val.TotalAmount-val.Amount {
		return errors.New("优惠券总数必须大于已使用张数")
	}

	val.OverTime = v.OverTime
	val.AllowEnable = v.AllowEnable
	val.BeginTime = v.BeginTime
	val.Code = v.Code
	val.Description = v.Description
	val.Discount = v.Discount
	val.Fee = v.Fee
	val.Integral = v.Integral
	val.MinFee = v.MinFee
	val.MinLevel = v.MinLevel
	val.NeedBind = v.NeedBind
	val.TotalAmount = v.TotalAmount
	val.UpdateTime = time.Now().Unix()
	return nil
}

func (this *Coupon) GetBinds() []promotion.ValueCouponBind {
	if !this.binds_loaded {
		this.binds = this.promRep.GetCouponBinds(this.value.Id)
	}
	return this.binds
}

func (this *Coupon) GetTakes() []promotion.ValueCouponTake {
	if !this.takes_loaded {
		this.takes = this.promRep.GetCouponTakes(this.value.Id)
	}
	return this.takes
}

func (this *Coupon) Save() (id int, err error) {
	if this.value.Id > 0 {
		if this.value.TotalAmount != this.value.Amount {
			errors.New("优惠券已被绑定或使用，不允许修改。")
		}
	} else {
		this.value.Amount = this.value.TotalAmount
	}
	return this.promRep.SaveCoupon(*this.value)
}

// 获取优惠券描述
func (this *Coupon) GetDescribe() string {
	buf := bytes.NewBufferString("")
	v := this.value

	if v.MinLevel != 0 {
		level := this.memberRep.GetLevel(v.MinLevel)
		buf.WriteString("[*" + level.Name + "]")
	}

	if v.MinFee == 0 {
		buf.WriteString("任意订单")
	} else {
		buf.WriteString(fmt.Sprintf("订单满%f", v.MinFee))
	}

	if v.Discount != 0 {
		dis := format.ToDiscountStr(v.Discount)
		buf.WriteString(fmt.Sprintf(",%s折优惠", dis))
		if v.Fee != 0 {
			buf.WriteString(fmt.Sprintf(",另减%d元", v.Fee))
		}
	} else if v.Fee != 0 {
		buf.WriteString(fmt.Sprintf(",减%d元", v.Fee))
	}

	if v.Integral != 0 {
		buf.WriteString(fmt.Sprintf(",赠送积分%d点", v.Integral))
	}

	return buf.String()
}

// 获取优惠的金额(四舍五入)
func (this *Coupon) GetCouponFee(orderFee float32) float32 {
	if float32(this.value.MinFee) > orderFee {
		return 0
	}
	var couponFee float32

	if this.value.Fee != 0 {
		couponFee = couponFee + float32(this.value.Fee)
	}

	if this.value.Discount != 100 {
		couponFee = couponFee + orderFee*
			(float32(100-this.value.Discount)/100)
	}
	return math.Round32(couponFee, 2)
}

// 是否可用
func (this *Coupon) CanUse(m member.IMember, fee float32) (bool, error) {
	mv := m.GetValue()
	cv := this.GetValue()

	if cv.AllowEnable == 0 {
		return false, errors.New("无效的优惠券")
	}

	dtUnix := time.Now().Unix()
	stUnix := cv.BeginTime
	ovUnix := cv.OverTime

	if dtUnix < stUnix {
		return false, errors.New(fmt.Sprintf("优惠券必须在%s~%s使用",
			time.Unix(cv.BeginTime, 0).Format("2006-01-02"),
			time.Unix(cv.OverTime, 0).Format("2006-01-02")),
		)
	} else if dtUnix > ovUnix {
		return false, errors.New("优惠拳已过期")
	}

	if cv.NeedBind == 0 && cv.Amount == 0 {
		return false, errors.New("优惠券不足")
	}

	if mv.Level < cv.MinLevel {
		return false, errors.New("会员等级不满足要求")
	}

	if fee < float32(cv.MinFee) {
		return false, errors.New(fmt.Sprintf("订单金额需达到￥%d", cv.MinFee))
	}

	return true, nil
}

/********  占用  *********/

//是否允许占用
func (this *Coupon) CanTake() bool {
	return this.value.NeedBind == 0
}

//获取占用
func (this *Coupon) GetTake(memberId int) (*promotion.ValueCouponTake, error) {
	return this.promRep.GetCouponTakeByMemberId(this.value.Id, memberId)
}

//占用
func (this *Coupon) Take(memberId int) error {
	if this.value.Amount == 0 {
		return errors.New("优惠券不足!")
	}

	dt := time.Now()

	valTake := &promotion.ValueCouponTake{
		MemberId:  memberId,
		CouponId:  this.value.Id,
		TakeTime:  dt.Unix(),
		ExtraTime: dt.Add(time.Hour * 4).Unix(), //4小时过期
		IsApply:   0,
		ApplyTime: dt.Add(-time.Hour).Unix(),
	}

	err := this.promRep.SaveCouponTake(valTake)
	if err == nil {
		this.value.Amount -= 1
		this.Save()
	}
	return err
}

//应用到订单
func (this *Coupon) ApplyTake(couponTakeId int) error {
	valTake := this.promRep.GetCouponTake(this.value.Id, couponTakeId)
	if valTake == nil {
		return errors.New("优惠券无效")
	}
	if valTake.IsApply == 1 {
		return errors.New("优惠券已使用")
	}

	now := time.Now().Unix()
	if now > valTake.ExtraTime {
		return errors.New("优惠券占用超时")
	}

	valTake.IsApply = 1
	valTake.ApplyTime = now

	return this.promRep.SaveCouponTake(valTake)
}

/********  绑定  *********/

//绑定
func (this *Coupon) Bind(memberId int) error {
	if this.value.Amount == 0 {
		return errors.New("优惠券不足")
	}

	var now time.Time = time.Now()

	valBind := &promotion.ValueCouponBind{
		MemberId: memberId,
		CouponId: this.value.Id,
		BindTime: now.Unix(),
		IsUsed:   0,
		UseTime:  now.Add(-time.Hour * 24).Unix(),
	}

	err := this.promRep.SaveCouponBind(valBind)
	if err == nil {
		this.value.Amount -= 1
		this.Save()
	}
	return err
}

//获取绑定
func (this *Coupon) GetBind(memberId int) (*promotion.ValueCouponBind, error) {
	return this.promRep.GetCouponBindByMemberId(this.value.Id, memberId)
}

func (this *Coupon) Binds(memberIds []string) error {
	if len(memberIds) > this.value.Amount {
		return errors.New(fmt.Sprintf("优惠券不足%s张，还剩%d张",
			len(memberIds), this.value.Amount))
	}

	for _, v := range memberIds {
		memberId, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		err = this.Bind(memberId)
		if err != nil {
			return err
		}
	}
	return nil
}

//使用优惠券
func (this *Coupon) UseCoupon(couponBindId int) error {
	valBind := this.promRep.GetCouponBind(this.value.Id, couponBindId)

	if valBind == nil {
		return errors.New("优惠券无效")
	}
	if valBind.IsUsed == 1 {
		return errors.New("优惠券已使用!")
	}

	valBind.UseTime = time.Now().Unix()
	valBind.IsUsed = 1
	return this.promRep.SaveCouponBind(valBind)
}
