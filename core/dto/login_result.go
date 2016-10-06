/**
 * Copyright 2015 @ z3q.net.
 * name : login_result.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

type MemberLoginResult struct {
	Result  bool
	Message string
	Member  *MemberSummary
}
