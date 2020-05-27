/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-09 17:49
 * description :
 * history :
 */
package dto

type SettleMeta struct {
	PaymentOpt int32              `json:"pay_opt"`
	DeliverOpt int32              `json:"deliver_opt"`
	Shop       *SettleShopMeta    `json:"shop"`
	Deliver    *SettleDeliverMeta `json:"deliver"`
}
