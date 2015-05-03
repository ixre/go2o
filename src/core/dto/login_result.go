/**
 * Copyright 2015 @ S1N1 Team.
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
    Result bool
    Message  string
    Token string
    Member *member.ValueMember
}