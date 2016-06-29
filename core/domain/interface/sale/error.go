/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-04 20:39
 * description :
 * history :
 */
package sale

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrGoodsNum *domain.DomainError = domain.NewDomainError(
		"err_goods_num", "商品数量错误")

	ErrOutOfSalePrice *domain.DomainError = domain.NewDomainError(
		"out_of_sale_price", "超出商品售价")

	ErrOutOfStock *domain.DomainError = domain.NewDomainError(
		"err_out_of_stock", "库存不足")

	ErrFullOfStock *domain.DomainError = domain.NewDomainError(
		"err_full_of_stock", "商品已经售完")

	ErrInternalDisallow *domain.DomainError = domain.NewDomainError(
		"err_sale_tag_internal_disallow", "不允许删除内置销售标签！")

	ErrCanNotDeleteItem *domain.DomainError = domain.NewDomainError(
		"err_goods_can_not_delete_item", "已售出货品只允许下架。")
)
