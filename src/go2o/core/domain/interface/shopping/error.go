/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 07:58
 * description :
 * history :
 */
package shopping

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrEmptyShoppingCart *domain.DomainError = domain.NewDomainError(
		"empty_shopping_cart", "购物车没有商品")

	ErrCartBuyerBinded *domain.DomainError = domain.NewDomainError(
		"cart_buyer_binded ", "购物车已绑定")

	ErrDisallowBindForCart *domain.DomainError = domain.NewDomainError(
		"cart_disallow_bind ", "无法为购物车绑定订单")
)
