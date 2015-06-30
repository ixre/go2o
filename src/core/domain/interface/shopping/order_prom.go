/**
 * Copyright 2015 @ S1N1 Team.
 * name : order_prom
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package shopping

type OrderPromotionBind struct {
	// 编号
	Id int `db:"id" pk:"yes" auto:"yes"`

	// 促销编号
	PromotionId int `db:"promotion_id"`

	// 订单号
	OrderNo string  `db:"order_no"`

	// 标题
	Title string `db:"title"`

	// 节省金额
	SaveFee float32 `db:"save_fee"`

	// 赠送积分
	PresentIntegral int `db:"present_integral"`

	// 是否应用
	IsApply int `db:"is_apply"`

	// 是否确认
	IsConfirm int `db:"is_confirm"`
}