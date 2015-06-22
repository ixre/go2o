/**
 * Copyright 2015 @ S1N1 Team.
 * name : sale_goods.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import "go2o/src/core/domain/interface/valueobject"

// 商品
type IGoods interface {
	// 获取领域对象编号
	GetDomainId() int

	// 获取货品
	GetItem() IItem

	// 设置值
	GetValue() *ValueGoods

	// 获取包装过的商品信息
	GetPackedValue() *valueobject.Goods

	// 获取促销价
	GetPromotionPrice(level int) float32

	// 获取会员价销价,返回是否有会原价及价格
	GetLevelPrice(level int) (bool, float32)

	// 获取促销描述
	GetPromotionDescribe() map[string]string

	// 获取会员价
	GetLevelPrices() []*MemberPrice

	// 保存会员价
	SaveLevelPrice(*MemberPrice) (int, error)

	// 设置值
	SetValue(*ValueGoods) error

	// 保存
	Save() (int, error)

	// 生成快照
	GenerateSnapshot() (int, error)

	// 获取最新的快照
	GetLatestSnapshot() *GoodsSnapshot
}
