/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-04 20:39
 * description :
 * history :
 */
package sale

import (
	"go2o/core/infrastructure/domain"
)

var (
	ErrNoSuchGoods *domain.DomainError = domain.NewDomainError(
		"no_such_goods", "商品不存在")

	ErrLatestSnapshot *domain.DomainError = domain.NewDomainError(
		"latest_snapshot", "已经是最新的快照")

	ErrNotOnShelves *domain.DomainError = domain.NewDomainError(
		"not_on_shelves", "商品未上架")
)
