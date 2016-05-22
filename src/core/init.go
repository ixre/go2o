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
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/merchant"
)

func init() {
	registerTypes()
}

// 注册序列类型
func registerTypes() {
	gob.Register(&member.ValueMember{})
	gob.Register(&merchant.MerchantValue{})
	gob.Register(&merchant.ApiInfo{})
}
