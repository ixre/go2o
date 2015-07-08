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
	"go2o/src/core/infrastructure/format"
	"strconv"
	"time"
)

func HandleCashBackDataTag(m member.IMember, order *shopping.ValueOrder,
	c promotion.ICashBackPromotion, memberRep member.IMemberRep) {
	data := c.GetDataTag()
	if _, ok := data["G1"]; ok {
		cashBack3R(m, order, c, memberRep)
	}
	// fmt.Println("---------[ xxx ]",len(data),data)
}

func cashBack3R(m member.IMember, order *shopping.ValueOrder, c promotion.ICashBackPromotion, memberRep member.IMemberRep) {

	var fee1 int
	var fee2 int

	dt := c.GetDataTag()
	fee1, _ = strconv.Atoi(dt["G1"])
	fee2, _ = strconv.Atoi(dt["G2"])

	var cm member.IMember = m
	var pm member.IMember = m

	// fmt.Println("------ START BACK ------")

	var backFunc = func(m member.IMember, parentM member.IMember, fee int) {
		// fmt.Println("---------[ back ]",parentM.GetValue().Name,fee)
		backCashForMember(m, order, fee, parentM.GetValue().Name)
	}
	var i int = 0
	for true {
		rl := cm.GetRelation()
		// fmt.Println("-------- BACK - ID - ",rl.InvitationMemberId)
		if rl == nil || rl.InvitationMemberId == 0 {
			break
		}

		cm = memberRep.GetMember(rl.InvitationMemberId)

		// fmt.Println("-------- BACK ",cm == nil)
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
	bfee := float32(fee)
	acc.PresentBalance += bfee // 更新赠送余额
	acc.TotalPresentFee += bfee
	acc.UpdateTime = time.Now().Unix()
	err := m.SaveAccount()

	if err == nil {
		//给自己返现
		icLog := &member.IncomeLog{
			MemberId:   m.GetAggregateRootId(),
			OrderId:    order.Id,
			Type:       "backcash",
			Fee:        float32(fee),
			Log:        fmt.Sprintf("推广返现￥%s元,订单号:%s,来源：%s", format.FormatFloat(bfee), order.OrderNo, refName),
			State:      1,
			RecordTime: acc.UpdateTime,
		}
		err = m.SaveIncomeLog(icLog)
	}
	return err
}
