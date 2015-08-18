/**
 * Copyright 2015 @ z3q.net.
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
	"strings"
)

func HandleCashBackDataTag(m member.IMember, order *shopping.ValueOrder,
	c promotion.ICashBackPromotion, memberRep member.IMemberRep) {
	data := c.GetDataTag()
	var level int = 0
	for k, _ := range data {
		if strings.HasPrefix(k, "G") {
			if l, err := strconv.Atoi(k[1:]); err == nil && l > level {
				level = l
			}
		}
	}
	//log.Println("[ Back][ Level] - ",level)
	cashBack3R(level, m, order, c, memberRep)
}

func cashBack3R(level int,m member.IMember, order *shopping.ValueOrder, c promotion.ICashBackPromotion, memberRep member.IMemberRep) {

	dt := c.GetDataTag()

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

		if fee, err := strconv.Atoi(dt[fmt.Sprintf("G%d",i)]);err == nil {
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

func backCashForMember(m member.IMember, order *shopping.ValueOrder, fee int, refName string) error {
	//更新账户
	acc := m.GetAccount()
	acv := acc.GetValue()
	bFee := float32(fee)
	acv.PresentBalance += bFee // 更新赠送余额
	acv.TotalPresentFee += bFee
	acv.UpdateTime = time.Now().Unix()
	_, err := acc.Save()

	if err == nil {
		//给自己返现
		icLog := &member.IncomeLog{
			MemberId:   m.GetAggregateRootId(),
			OrderId:    order.Id,
			Type:       "backcash",
			Fee:        float32(fee),
			Log:        fmt.Sprintf("推广返现￥%s元,订单号:%s,来源：%s", format.FormatFloat(bFee), order.OrderNo, refName),
			State:      1,
			RecordTime: acv.UpdateTime,
		}
		err = m.SaveIncomeLog(icLog)
	}
	return err
}
