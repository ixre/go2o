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

	// 获取值
	GetValue() *ValueSaleTag

	// 设置值
	SetValue(v *ValueSaleTag) error

	// 保存
	Save() (int, error)

	// 是否为系统内置
	System() bool

	// 获取标签下的商品
	GetValueGoods(begin, end int) []*valueobject.Goods

	// 获取标签下的分页商品
	GetPagedValueGoods(begin, end int) (total int, goods []*valueobject.Goods)
}
