/**
 * Copyright 2015 @ S1N1 Team.
 * name : value_promotion
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

type ValuePromotion struct {
	// 促销编号
	Id int `db:"id" pk:"yes" auto:"yes"`

	// 商户编号
	PartnerId int `db:"partner_id"`

	// 促销简称
	ShortName string `db:"short_name"`

	// 促销描述
	Description string `db:"description"`

	// 类型位值
	TypeFlag int `db:"type_flag"`

	// 商品编号(为0则应用订单)
	GoodsId int `db:"goods_id"`

	// 是否启用
	Enabled int `db:"enabled"`

	// 修改时间
	UpdateTime int64 `db:"update_time"`
}
