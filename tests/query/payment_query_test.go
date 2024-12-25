/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: payment_query_test.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-09-09 08:58:27
 * description: 支付查询测试
 * history:
 */
package query

import (
	"testing"

	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/gof/typeconv"
)

// 测试查询可分账的支付订单
func TestPaymentQuery_QueryDivideOrders(t *testing.T) {
	qs := inject.GetPaymentQueryService()
	orders := qs.QueryDivideOrders(848, payment.TypeRecharge)
	t.Errorf("orders: %s", typeconv.MustJson(orders))
}

func TestPaymentQuery_QueryRefundableOrders(t *testing.T) {
	memberId := 848
	qs := inject.GetPaymentQueryService()
	orders := qs.QueryRefundableRechargeOrders(memberId)
	t.Errorf("orders: %s", typeconv.MustJson(orders))
}
