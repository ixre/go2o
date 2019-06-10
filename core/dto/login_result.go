/**
 * Copyright 2015 @ to2.net.
 * name : login_result.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

// 登录的会员信息
type LoginMember struct {
	ID         int
	Code       string
	Token      string
	UpdateTime int64
}

// 会员登录返回结果
type MemberLoginResult struct {
	ErrCode int
	ErrMsg  string
	Member  *LoginMember
}
