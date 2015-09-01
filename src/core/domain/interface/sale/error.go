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
	"go2o/src/core/infrastructure/domain"
)

var (
	ErrNoSuchGoods *domain.DomainError = domain.NewDomainError(
		"no_such_goods", "商品不存在")

	ErrLatestSnapshot *domain.DomainError = domain.NewDomainError(
		"latest_snapshot", "已经是最新的快照")

	ErrNoSuchSnapshot *domain.DomainError = domain.NewDomainError(
		"no_such_snapshot", "商品快照不存在")

	ErrNotOnShelves *domain.DomainError = domain.NewDomainError(
		"not_on_shelves", "商品未上架")

	ErrOutOfSalePrice *domain.DomainError = domain.NewDomainError(
		"out_of_sale_price", "超出商品售价")

	ErrInternalDisallow *domain.DomainError = domain.NewDomainError(
		"err_sale_tag_internal_disallow", "不允许删除内置销售标签！")

	ErrCanNotDeleteItem *domain.DomainError = domain.NewDomainError(
		"err_goods_can_not_delete_item", "已售出货品只允许下架。")
)
