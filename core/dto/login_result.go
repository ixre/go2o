/**
 * Copyright 2015 @ 56x.net.
 * name : login_result.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

// LoginMember 登录的会员信息
type LoginMember struct {
	ID         int64
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
