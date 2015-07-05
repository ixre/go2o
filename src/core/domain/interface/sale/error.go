/**
 * Copyright 2014 @ S1N1 Team.
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
		"err_internal_disallow", "不允许删除内置销售标签！")
)
