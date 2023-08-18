/**
 * Copyright 2015 @ 56x.net.
 * name : cash_back
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package order

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/format"
)

// 获取推荐数组
func (o *subOrderImpl) getReferArr(memberId int64, level int32) []int64 {
	arr := make([]int64, level)
	var i int32
	referId := memberId
	for i <= level-1 {
		rl := o.memberRepo.GetRelation(referId)
		if rl == nil || rl.InviterId <= 0 {
			break
		}
		arr[i] = rl.InviterId
		referId = arr[i]
		i++
	}
	return arr
}

func (o *subOrderImpl) handleCashBack() error {
	gobConf := o.valRepo.GetGlobMchSaleConf()
	if !gobConf.FxSalesEnabled {
		return nil
	}
	var err error
	v := o.value
	mch := o.mchRepo.GetMerchant(int(v.VendorId))
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		buyer := o.getBuyer()
		now := time.Now().Unix()

		//******* 返现到账户  ************
		var back_fee int64
		saleConf := mch.ConfManager().GetSaleConf()

		if saleConf.CashBackPercent > 0 {
			back_fee = int64(float32(v.FinalAmount) * saleConf.CashBackPercent)
			//将此次消费记入会员账户
			err = o.updateShoppingMemberBackFee(mch.GetValue().Name, buyer,
				int64(float32(back_fee)*saleConf.CashBackMemberPercent), now)
			domain.HandleError(err, "domain")

		}

		// 处理返现促销
		//todo: ????
		//o.handleCashBackPromotions(mch, m)
		// 三级返现
		if back_fee > 0 {
			err = o.backFor3R(mch, buyer, back_fee, now)
		}
	}
	return err
}

func (o *subOrderImpl) updateMemberAccount(m member.IMemberAggregateRoot,
	ptName, mName string, fee int64, unixTime int64) error {
	if fee > 0 {
		//更新账户
		acc := m.GetAccount()
		acv := acc.GetValue()
		acv.WalletBalance += fee
		acv.TotalWalletAmount += fee
		acv.UpdateTime = unixTime
		_, err := acc.Save()
		if err == nil {
			//给自己返现
			tit := fmt.Sprintf("订单:%s(商户:%s,会员:%s)收入￥%.2f元",
				o.value.OrderNo, ptName, mName, fee)
			_, err = acc.CarryTo(member.AccountWallet, member.AccountOperateData{
				Title:   tit,
				Amount:  int(fee * 100),
				OuterNo: o.value.OrderNo,
				Remark:  "sys",
			}, false, 0)
		}
		return err
	}
	return nil
}

// 三级返现
func (o *subOrderImpl) backFor3R(mch merchant.IMerchant, m member.IMemberAggregateRoot,
	back_fee int64, unixTime int64) (err error) {
	if back_fee > 0 {
		i := 0
		mName := m.Profile().GetProfile().Name
		saleConf := mch.ConfManager().GetSaleConf()
		percent := saleConf.CashBackTg2Percent
		for i < 2 {
			rl := m.GetRelation()
			if rl == nil || rl.InviterId == 0 {
				break
			}

			m = o.memberRepo.GetMember(rl.InviterId)
			if m == nil {
				break
			}

			if i == 1 {
				percent = saleConf.CashBackTg1Percent
			}

			err = o.updateMemberAccount(m, mch.GetValue().Name, mName,
				int64(float32(back_fee)*percent), unixTime)
			if err != nil {
				domain.HandleError(err, "domain")
				break
			}
			i++
		}
	}
	return err
}

func HandleCashBackDataTag(m member.IMemberAggregateRoot, o order.IOrder,
	c promotion.ICashBackPromotion, memberRepo member.IMemberRepo) {
	data := c.GetDataTag()
	level := 0
	for k := range data {
		if strings.HasPrefix(k, "G") {
			if l, err := strconv.Atoi(k[1:]); err == nil && l > level {
				level = l
			}
		}
	}
	//log.Println("[ Back][ Level] - ",level)
	cashBack3R(level, m, o, c, memberRepo)
}

func cashBack3R(level int, m member.IMemberAggregateRoot, o order.IOrder,
	c promotion.ICashBackPromotion, memberRepo member.IMemberRepo) {

	dt := c.GetDataTag()

	cm := m
	pm := m

	// fmt.Println("------ START BACK ------")

	var backFunc = func(m member.IMemberAggregateRoot, parentM member.IMemberAggregateRoot, fee int) {
		// fmt.Println("---------[ back ]",parentM.Value().Name,fee)
		backCashForMember(m, o, fee, parentM.Profile().GetProfile().Name)
	}
	var i = 0
	for true {
		rl := cm.GetRelation()
		// fmt.Println("-------- BACK - Id - ",rl.InvitationMemberId)
		if rl == nil || rl.InviterId == 0 {
			break
		}

		cm = memberRepo.GetMember(rl.InviterId)

		// fmt.Println("-------- BACK ",cm == nil)
		if m == nil {
			break
		}

		if fee, err := strconv.Atoi(dt[fmt.Sprintf("G%d", i)]); err == nil {
			//log.Println("[ Back][ Cash] - ",i," back ",fee)
			backFunc(cm, pm, fee)
		}

		pm = cm

		i++
		if i >= level {
			break
		}
	}
}

func backCashForMember(m member.IMemberAggregateRoot, o order.IOrder,
	fee int, refName string) error {
	//更新账户
	acc := m.GetAccount()
	acv := acc.GetValue()
	bFee := int64(fee)
	acv.WalletBalance += bFee // 更新赠送余额
	acv.TotalWalletAmount += bFee
	acv.UpdateTime = time.Now().Unix()
	_, err := acc.Save()

	if err == nil {
		orderNo := o.OrderNo()
		tit := fmt.Sprintf("推广返现￥%s元,订单号:%s,来源：%s",
			format.FormatIntMoney(bFee), orderNo, refName)
		_, err = acc.CarryTo(member.AccountWallet, member.AccountOperateData{
			Title:   tit,
			Amount:  fee * 100,
			OuterNo: orderNo,
			Remark:  "sys",
		}, false, 0)
	}
	return err
}
