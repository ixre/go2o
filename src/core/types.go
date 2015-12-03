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
	"go2o/src/core/domain/interface/partner"
)

// 注册序列号类型
func RegisterTypes() {
	gob.Register(&member.ValueMember{})
	gob.Register(&partner.ValuePartner{})
	gob.Register(&partner.ApiInfo{})
}
