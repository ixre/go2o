/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-02 21:34
 * description :
 * history :
 */

package promotion

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ixre/gof/math"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/promotion"
	"go2o/core/infrastructure/format"
	"strconv"
	"time"
)

var _ promotion.IPromotion = new(Coupon)
var _ promotion.ICouponPromotion = new(Coupon)

// 优惠券,一张优惠券有数量，对应唯一的促销码。
// 优惠内容包含：送金额，送积分,订单折扣。仅在消费时有效。
// 使用需要达到最低金额和最低等级。
// 优惠券包含了开始时间和结束时间，超出时间则过期。
// 优惠券可以设置启用和停用
// 是否允许绑定，如果不绑定。则可以任意使用.只要有绑定和使用后，就不允许修改此属性。
type Coupon struct {
	*promotionImpl
	detailsValue *promotion.ValueCoupon
	promRepo     promotion.IPromotionRepo
	memberRepo   member.IMemberRepo
	takes        []promotion.ValueCouponTake
	binds        []promotion.ValueCouponBind
	takesLoaded  bool
	bindsLoaded  bool
}

func newCoupon(p *promotionImpl, v *promotion.ValueCoupon, promRepo promotion.IPromotionRepo,
	memberRepo member.IMemberRepo) *Coupon {
	cp := &Coupon{
		detailsValue:  v,
		promotionImpl: p,
		promRepo:      promRepo,
		memberRepo:    memberRepo,
	}
	cp.releaseCoupon()
	return cp
}

func (c *Coupon) GetDomainId() int32 {
	return c.detailsValue.Id
}

// 释放优惠券
func (c *Coupon) releaseCoupon() {
	// 仅在会员通用情况下才存在占用
	if c.detailsValue.NeedBind == 0 &&
		c.detailsValue.TotalAmount != c.detailsValue.Amount {
		now := time.Now().Unix()
		oriAmount := c.detailsValue.Amount
		for _, take := range c.GetTakes() {
			// 未应用到订单，且释放时间小于当前时间，则释放
			if take.IsApply == 0 && now > take.ExtraTime {
				c.detailsValue.Amount = c.detailsValue.Amount + 1
			}
		}

		//保存新的可用数量
		if oriAmount != c.detailsValue.Amount {
			c.Save()
		}
	}
}

// 获取相关的值
func (c *Coupon) GetRelationValue() interface{} {
	return c.detailsValue
}

// 促销类型
func (c *Coupon) TypeName() string {
	return "优惠券"
}

// 获取促销内容
func (c *Coupon) GetDetailsValue() promotion.ValueCoupon {
	return *c.detailsValue
}

// 设置促销内容
func (c *Coupon) SetDetailsValue(v *promotion.ValueCoupon) error {

	val := c.detailsValue
	if v.TotalAmount < val.TotalAmount-val.Amount {
		return errors.New("优惠券总数必须大于已使用张数")
	}
	//	if c._detailsValue.TotalAmount != c._detailsValue.Amount {
	//		return c.GetAggregateRootId(), errors.New("优惠券已被绑定或使用，不允许修改数量。")
	//	}

	val.OverTime = v.OverTime
	val.BeginTime = v.BeginTime
	val.Code = v.Code
	val.Discount = v.Discount
	val.Fee = v.Fee
	val.Integral = v.Integral
	val.MinFee = v.MinFee
	val.MinLevel = v.MinLevel
	val.NeedBind = v.NeedBind
	val.TotalAmount = v.TotalAmount
	return nil
}

func (c *Coupon) GetBinds() []promotion.ValueCouponBind {
	if !c.bindsLoaded {
		c.binds = c.promRepo.GetCouponBinds(c.detailsValue.Id)
	}
	return c.binds
}

func (c *Coupon) GetTakes() []promotion.ValueCouponTake {
	if !c.takesLoaded {
		c.takes = c.promRepo.GetCouponTakes(c.detailsValue.Id)
	}
	return c.takes
}

func (c *Coupon) Save() (int32, error) {

	if c.GetRelationValue() == nil {
		return c.GetAggregateRootId(), promotion.ErrCanNotApplied
	}

	if c.detailsValue.Id <= 0 {
		c.detailsValue.Amount = c.detailsValue.TotalAmount
	}

	var isCreate bool = c.GetAggregateRootId() == 0

	id, err := c.promotionImpl.Save()
	c.value.Id = id

	if err == nil {
		c.detailsValue.Id = c.GetAggregateRootId()
		return c.promRepo.SaveValueCoupon(c.detailsValue, isCreate)
	}
	return id, err
}

// 获取优惠券描述
func (c *Coupon) GetDescribe() string {
	buf := bytes.NewBufferString("")
	v := c.detailsValue

	if v.MinLevel != 0 {
		level := c.memberRepo.GetManager().LevelManager().GetLevelById(v.MinLevel)
		buf.WriteString("[*" + level.Name + "]")
	}

	if v.MinFee == 0 {
		buf.WriteString("任意订单")
	} else {
		buf.WriteString(fmt.Sprintf("订单满%d", v.MinFee))
	}

	if v.Discount != 0 && v.Discount != 100 {
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
func (c *Coupon) GetCouponFee(orderFee float32) float32 {
	if float32(c.detailsValue.MinFee) > orderFee {
		return 0
	}
	var couponFee float32

	if c.detailsValue.Fee != 0 {
		couponFee = couponFee + float32(c.detailsValue.Fee)
	}

	if c.detailsValue.Discount != 100 {
		couponFee = couponFee + orderFee*
			(float32(100-c.detailsValue.Discount)/100)
	}
	return math.Round32(couponFee, 2)
}

// 是否可用
func (c *Coupon) CanUse(m member.IMember, fee float32) (bool, error) {
	mv := m.GetValue()
	cv := c.GetDetailsValue()

	if c.value.Enabled == 0 {
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
func (c *Coupon) CanTake() bool {
	return c.detailsValue.NeedBind == 0
}

//获取占用
func (c *Coupon) GetTake(memberId int64) (*promotion.ValueCouponTake, error) {
	return c.promRepo.GetCouponTakeByMemberId(c.detailsValue.Id, memberId)
}

//占用
func (c *Coupon) Take(memberId int64) error {
	if c.detailsValue.Amount == 0 {
		return errors.New("优惠券不足!")
	}

	dt := time.Now()

	valTake := &promotion.ValueCouponTake{
		MemberId:  memberId,
		CouponId:  c.detailsValue.Id,
		TakeTime:  dt.Unix(),
		ExtraTime: dt.Add(time.Hour * 4).Unix(), //4小时过期
		IsApply:   0,
		ApplyTime: dt.Add(-time.Hour).Unix(),
	}

	err := c.promRepo.SaveCouponTake(valTake)
	if err == nil {
		c.detailsValue.Amount -= 1
		c.Save()
	}
	return err
}

//应用到订单
func (c *Coupon) ApplyTake(couponTakeId int32) error {
	valTake := c.promRepo.GetCouponTake(c.detailsValue.Id, couponTakeId)
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

	return c.promRepo.SaveCouponTake(valTake)
}

/********  绑定  *********/

//绑定
func (c *Coupon) Bind(memberId int64) error {
	if c.detailsValue.Amount == 0 {
		return errors.New("优惠券不足")
	}

	var now time.Time = time.Now()

	valBind := &promotion.ValueCouponBind{
		MemberId: memberId,
		CouponId: c.detailsValue.Id,
		BindTime: now.Unix(),
		IsUsed:   0,
		UseTime:  now.Add(-time.Hour * 24).Unix(),
	}

	err := c.promRepo.SaveCouponBind(valBind)
	if err == nil {
		c.detailsValue.Amount -= 1
		_, err = c.Save()
	}
	return err
}

//获取绑定
func (c *Coupon) GetBind(memberId int64) (*promotion.ValueCouponBind, error) {
	return c.promRepo.GetCouponBindByMemberId(c.detailsValue.Id, memberId)
}

func (c *Coupon) Binds(memberIds []string) error {
	if len(memberIds) > c.detailsValue.Amount {
		return errors.New(fmt.Sprintf("优惠券不足%s张，还剩%d张",
			len(memberIds), c.detailsValue.Amount))
	}

	for _, v := range memberIds {
		memberId, err := util.I64Err(strconv.Atoi(v))
		if err != nil {
			return err
		}

		err = c.Bind(memberId)
		if err != nil {
			return err
		}
	}
	return nil
}

//使用优惠券
func (c *Coupon) UseCoupon(couponBindId int32) error {
	valBind := c.promRepo.GetCouponBind(c.detailsValue.Id, couponBindId)

	if valBind == nil {
		return errors.New("优惠券无效")
	}
	if valBind.IsUsed == 1 {
		return errors.New("优惠券已使用!")
	}

	valBind.UseTime = time.Now().Unix()
	valBind.IsUsed = 1
	return c.promRepo.SaveCouponBind(valBind)
}
