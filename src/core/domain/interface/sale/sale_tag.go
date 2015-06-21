/**
 * Copyright 2015 @ S1N1 Team.
 * name : sale_tag
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import "go2o/src/core/domain/interface/valueobject"

type ISaleTag interface {
	GetDomainId() int

	GetValue() *ValueSaleTag

	SetValue(v *ValueSaleTag) error

	Save() (int, error)

	// 获取标签下的商品
	GetValueGoods(begin, end int) []*valueobject.Goods
}
