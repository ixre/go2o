/**
 * Copyright 2015 @ 56x.net.
 * name : member
 * author : jarryliu
 * date : 2015-10-29 15:06
 * description :
 * history :
 */
package dto

type (
	SimpleMember struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		User   string `json:"user"`
		Phone  string `json:"phone"`
		Avatar string `json:"avatar"`
	}

	// 邀请会员数据
	InvitationMember struct {
		// 会员编号
		MemberId int32
		// 用户名
		User string
		// 等级
		Level int32
		// 头像
		Avatar string
		// 昵称
		Nickname string
		// 电话
		Phone string
		// 注册时间
		RegTime int64
		// 邀请人数
		InvitationNum int
	}
)
