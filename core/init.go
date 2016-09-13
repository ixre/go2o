/**
 * Copyright 2015 @ z3q.net.
 * name : types.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package core

import (
	"encoding/gob"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/payment"
	"go2o/core/dto"
)

func init() {
	registerTypes()
}

// 注册序列类型
func registerTypes() {
	gob.Register(&member.Member{})
	gob.Register(&merchant.Merchant{})
	gob.Register(&merchant.ApiInfo{})
	gob.Register(&shop.OnlineShop{})
	gob.Register(&shop.OfflineShop{})
	gob.Register(&shop.ShopDto{})
	gob.Register(&member.Account{})
	gob.Register(&payment.PaymentOrderBean{})
	gob.Register(&member.Relation{})
	gob.Register(&dto.ListOnlineShop{})
}
