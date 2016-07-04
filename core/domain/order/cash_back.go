/**
 * Copyright 2015 @ z3q.net.
 * name : cash_back
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package order

import (
	"fmt"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/promotion"
	"go2o/core/infrastructure/format"
	"strconv"
	"strings"
	"time"
)

func HandleCashBackDataTag(m member.IMember, order *order.Order,
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

func cashBack3R(level int, m member.IMember, order *order.Order, c promotion.ICashBackPromotion, memberRep member.IMemberRep) {

	dt := c.GetDataTag()

	var cm member.IMember = m
	var pm member.IMember = m

	// fmt.Println("------ START BACK ------")

	var backFunc = func(m member.IMember, parentM member.IMember, fee int) {
		// fmt.Println("---------[ back ]",parentM.GetValue().Name,fee)
		backCashForMember(m, order, fee, parentM.Profile().GetProfile().Name)
	}
	var i int = 0
	for true {
		rl := cm.GetRelation()
		// fmt.Println("-------- BACK - ID - ",rl.InvitationMemberId)
		if rl == nil || rl.RefereesId == 0 {
			break
		}

		cm = memberRep.GetMember(rl.RefereesId)

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

func backCashForMember(m member.IMember, order *order.Order,
	fee int, refName string) error {
	//更新账户
	acc := m.GetAccount()
	acv := acc.GetValue()
	bFee := float32(fee)
	acv.PresentBalance += bFee // 更新赠送余额
	acv.TotalPresentFee += bFee
	acv.UpdateTime = time.Now().Unix()
	_, err := acc.Save()

	if err == nil {
		tit := fmt.Sprintf("推广返现￥%s元,订单号:%s,来源：%s",
			format.FormatFloat(bFee), order.OrderNo, refName)
		err = acc.PresentBalance(tit, order.OrderNo, float32(fee))
	}
	return err
}
