/**
 * Copyright 2015 @ S1N1 Team.
 * name : invitaton_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package member

type IInvitationManager interface {
	// 获取我邀请的会员
	GetMyInvitationMembers() []*ValueMember

	// 获取我的邀请码
	GetMyInvitationCode() string

	// 获取邀请会员下级邀请数量
	GetSubInvitationNum() map[int]int

	// 获取邀请要的会员
	GetInvitationMeMember() *ValueMember
}
