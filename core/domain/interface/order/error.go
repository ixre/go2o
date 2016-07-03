/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 07:58
 * description :
 * history :
 */
package order

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrNoSuchOrder *domain.DomainError = domain.NewDomainError(
		"err_no_such_order ", "订单不存在")

	ErrOrderPayed *domain.DomainError = domain.NewDomainError(
		"err_order_payed ", "订单已支付")

	ErrOrderNotPayed *domain.DomainError = domain.NewDomainError(
		"err_order_not_payed ", "订单未支付")

	ErrOrderDelved *domain.DomainError = domain.NewDomainError(
		"err_order_delved ", "订单已发货")
)
