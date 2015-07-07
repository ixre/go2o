/**
 * Copyright 2015 @ S1N1 Team.
 * name : cash_back
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package shopping

import (
	"fmt"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/promotion"
	"go2o/src/core/domain/interface/shopping"
	"strconv"
	"time"
)

func HandleCashBackDataTag(m member.IMember, order *shopping.ValueOrder,
	c promotion.ICashBackPromotion, memberRep member.IMemberRep) {
	data := c.GetDataTag()
	if _, ok := data["G1"]; ok {
		cashBack3R(m, order, c, memberRep)
	}
}

func cashBack3R(m member.IMember, order *shopping.ValueOrder, c promotion.ICashBackPromotion, memberRep member.IMemberRep) {

	var fee1 int
	var fee2 int

	dt := c.GetDataTag()
	fee1, _ = strconv.Atoi(dt["G1"])
	fee2, _ = strconv.Atoi(dt["G2"])

	var cm member.IMember = m
	var pm member.IMember = m

	var backFunc = func(m member.IMember, parentM member.IMember, fee int) {
		backCashForMember(m, order, fee, parentM.GetValue().Name)
	}
	var i int = 0
	for true {
		rl := cm.GetRelation()
		if rl == nil || rl.InvitationMemberId == 0 {
			break
		}

		cm = memberRep.GetMember(rl.InvitationMemberId)
		if m == nil {
			break
		}

		if i == 0 {
			backFunc(cm, pm, fee2)
		} else if i == 1 {
			backFunc(cm, pm, fee1)
		}

		pm = cm

		i++
		if i > 1 {
			break
		}
	}
}

func backCashForMember(m member.IMember, order *shopping.ValueOrder, fee int, refName string) error {
	//更新账户
	acc := m.GetAccount()
	acc.PresentBalance = acc.PresentBalance + float32(fee) // 更新赠送余额
	acc.UpdateTime = time.Now().Unix()
	err := m.SaveAccount()

	if err == nil {
		//给自己返现
		icLog := &member.IncomeLog{
			MemberId:   m.GetAggregateRootId(),
			OrderId:    order.Id,
			Type:       "backcash",
			Fee:        float32(fee),
			Log:        fmt.Sprintf("推广返现￥%.2f元,订单号:%s,来源：%s", order.OrderNo, fee, refName),
			State:      1,
			RecordTime: acc.UpdateTime,
		}
		err = m.SaveIncomeLog(icLog)
	}
	return err
}
