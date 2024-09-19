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
	_orderRepo       fw.Repository[payment.Order]
	_divideRepo      fw.Repository[payment.PayDivide]
	_subMerchantRepo fw.Repository[payment.PayMerchant]
}

func NewPaymentQuery(o fw.ORM) *PaymentQuery {
	return &PaymentQuery{ORM: o,
		_orderRepo:       &fw.BaseRepository[payment.Order]{ORM: o},
		_divideRepo:      &fw.BaseRepository[payment.PayDivide]{ORM: o},
		_subMerchantRepo: &fw.BaseRepository[payment.PayMerchant]{ORM: o},
	}
}

// 查询可分账的支付订单
func (p *PaymentQuery) QueryDivideOrders(memberId int, orderType int) []*DivideOrderInfo {
	arr := make([]*DivideOrderInfo, 0)
	orders := p._orderRepo.FindList(nil, "buyer_id=? AND order_type=? AND status = ? AND divide_status < ? ORDER BY id ASC",
		memberId,
		orderType,
		payment.StateFinished,
		payment.DivideCompleted,
	)
	payIds := make([]int, 0)
	mp := make(map[int]*DivideOrderInfo)
	for _, v := range orders {
		// 支付金额减去退款金额为实际可分账金额
		amount := v.FinalAmount - v.RefundAmount
		if amount > 0 {
			payIds = append(payIds, v.Id)
			dst := &DivideOrderInfo{
				TradeNo:       v.TradeNo,
				Amount:        amount,
				DividedAmount: 0,
				CreateTime:    v.SubmitTime,
				DivideStatus:  v.DivideStatus,
			}
			arr = append(arr, dst)
			mp[v.Id] = dst
		}
	}
	if len(payIds) > 0 {
		// 查询已分账的记录
		divides := p._divideRepo.FindList(nil, "pay_id IN (?)", payIds)
		for _, v := range divides {
			dst := mp[v.PayId]
			dst.DividedAmount += v.DivideAmount
			mp[v.PayId] = dst
		}
	}
	return arr
}

// 查询待提交的分账记录
func (p *PaymentQuery) QueryAwaitSubmitDivides(unix int64, size int) ([]*payment.PayDivide, error) {
	rows := p._divideRepo.FindList(nil, "submit_status=? AND submit_time < ? AND user_id <> 0", payment.DivideItemStatusPending, unix)
	return rows, nil
}

// GetPaymentOrder 获取支付单
func (p *PaymentQuery) GetPaymentOrder(payId int) *payment.Order {
	return p._orderRepo.Get(payId)
}

// GetPayDivide 获取分账记录

func (p *PaymentQuery) GetPayDivide(divideId int) *payment.PayDivide {
	return p._divideRepo.Get(divideId)
}

// GetSubMerchant 获取子商户
func (p *PaymentQuery) GetSubMerchant(subType int, mchId int) *payment.PayMerchant {
	return p._subMerchantRepo.FindBy("user_type=? AND user_id=?", subType, mchId)
}

// QueryPagingSubMerchant 分页查询入网商户
func (p *PaymentQuery) QueryPagingSubMerchant(params *fw.PagingParams) (*fw.PagingResult, error) {
	return p._subMerchantRepo.QueryPaging(params)
}
