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

	ErrNotOnShelves *domain.DomainError = domain.NewDomainError(
		"not_on_shelves", "商品未上架")

ErrOutOfSalePrice *domain.DomainError = domain.NewDomainError(
	"out_of_sale_price", "超出商品售价")
)
