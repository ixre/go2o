/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-09 17:49
 * description :
 * history :
 */
package dto

type SettleMeta struct {
	PaymentOpt int                `json:"pay_opt"`
	DeliverOpt int                `json:"deliver_opt"`
	Shop       *SettleShopMeta    `json:"shop"`
	Deliver    *SettleDeliverMeta `json:"deliver"`
}
