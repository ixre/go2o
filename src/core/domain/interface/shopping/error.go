/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 07:58
 * description :
 * history :
 */
package shopping

import (
	"go2o/src/core/infrastructure/domain"
)

var (
	ErrEmptyShoppingCart *domain.DomainError = domain.NewDomainError(
		"empty_shopping_cart", "购物车没有商品")

	ErrCartBuyerBinded *domain.DomainError = domain.NewDomainError(
		"cart_buyer_binded ", "购物车已绑定")

	ErrDisallowBindForCart *domain.DomainError = domain.NewDomainError(
		"cart_disallow_bind ", "无法为购物车绑定订单")

	ErrBalanceNotEnough *domain.DomainError = domain.NewDomainError(
		"rtt_balance_not_enough ", "余额不足")

	ErrOrderPayed *domain.DomainError = domain.NewDomainError(
		"err_order_payed ", "订单已支付")

	ErrOrderNotPayed *domain.DomainError = domain.NewDomainError(
		"err_order_not_payed ", "订单未支付")
)
