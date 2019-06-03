/**
 * Copyright 2015 @ to2.net.
 * name : invitaton_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package member

import "go2o/core/dto"

type IInvitationManager interface {
	// 获取邀请人数组
	InviterArray(memberId int64, depth int) []int64

	// 判断是否由会员邀请
	InvitationBy(memberId int64) bool

	// 获取我邀请的会员
	GetInvitationMembers(begin, end int) (total int, rows []*dto.InvitationMember)

	// 获取我的邀请码
	MyCode() string

	// 获取邀请会员下级邀请数量
	GetSubInvitationNum(memberIdArr []int32) map[int32]int

	// 获取邀请我的会员
	GetInvitationMeMember() *Member
}
