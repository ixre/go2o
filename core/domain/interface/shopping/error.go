/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 07:58
 * description :
 * history :
 */
package shopping

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrBalanceNotEnough *domain.DomainError = domain.NewDomainError(
		"rtt_balance_not_enough ", "余额不足")

	ErrOrderPayed *domain.DomainError = domain.NewDomainError(
		"err_order_payed ", "订单已支付")

	ErrOrderNotPayed *domain.DomainError = domain.NewDomainError(
		"err_order_not_payed ", "订单未支付")
)
