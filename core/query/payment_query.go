/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: payment_query.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-09-09 08:29:05
 * description: 支付查询
 * history:
 */
package query

import (
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type DivideOrderInfo struct {
	TradeNo       string
	Amount        int
	DividedAmount int
	CreateTime    int
	DivideStatus  int
}

type PaymentQuery struct {
	fw.ORM
	_orderRepo  fw.Repository[payment.Order]
	_divideRepo fw.Repository[payment.PayDivide]
}

func NewPaymentQuery(o fw.ORM) *PaymentQuery {
	return &PaymentQuery{ORM: o,
		_orderRepo:  &fw.BaseRepository[payment.Order]{ORM: o},
		_divideRepo: &fw.BaseRepository[payment.PayDivide]{ORM: o},
	}
}

// 查询可分账的支付订单
func (p *PaymentQuery) QueryDivideOrders(memberId int, orderType int) []*DivideOrderInfo {
	arr := make([]*DivideOrderInfo, 0)
	orders := p._orderRepo.FindList(nil, "buyer_id=? AND order_type=? AND status = ? AND divide_status <> ? ORDER BY id ASC",
		memberId,
		orderType,
		payment.StateFinished,
		payment.DivideFinished,
	)
	payIds := make([]int, 0)
	mp := make(map[int]*DivideOrderInfo)
	for _, v := range orders {
		payIds = append(payIds, v.Id)
		dst := &DivideOrderInfo{
			TradeNo:       v.TradeNo,
			Amount:        v.FinalAmount,
			DividedAmount: 0,
			CreateTime:    v.SubmitTime,
			DivideStatus:  v.DivideStatus,
		}
		arr = append(arr, dst)
		mp[v.Id] = dst
	}
	divides := p._divideRepo.FindList(nil, "pay_id IN (?)", payIds)
	for _, v := range divides {
		dst := mp[v.PayId]
		dst.DividedAmount += v.DivideAmount
		mp[v.PayId] = dst
	}
	return arr
}
