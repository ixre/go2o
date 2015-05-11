/**
 * Copyright 2015 @ S1N1 Team.
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
	"go2o/src/core/dto"
)

// 注册序列号类型
func RegisterTypes() {
	gob.Register(&member.ValueMember{})
	gob.Register(&partner.ValuePartner{})
	gob.Register(&dto.PartnerApiInfo{})
}
