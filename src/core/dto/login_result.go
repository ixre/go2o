/**
 * Copyright 2015 @ z3q.net.
 * name : login_result.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

import (
	"go2o/src/core/domain/interface/member"
)

type MemberLoginResult struct {
	Result  bool
	Message string
	Member  *member.ValueMember
}
