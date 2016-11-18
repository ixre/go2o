/**
 * Copyright 2015 @ z3q.net.
 * name : invitation_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package member

import (
	"go2o/core/domain/interface/member"
	"go2o/core/dto"
)

var _ member.IInvitationManager = new(invitationManager)

type invitationManager struct {
	member       *memberImpl
	myInvMembers []*member.Member
}

// 判断是否推荐了某个会员
func (i *invitationManager) InvitationBy(memberId int32) bool {
	rl := i.member.GetRelation()
	if rl != nil {
		return rl.RefereesId == memberId
	}
	return false
}

// 获取我邀请的会员
func (i *invitationManager) GetInvitationMembers(begin, end int) (
	int, []*dto.InvitationMember) {
	return i.member.rep.GetMyInvitationMembers(
		i.member.GetAggregateRootId(), begin, end)
}

// 获取我的邀请码
func (i *invitationManager) MyCode() string {
	return i.member.GetValue().InvitationCode
}

// 获取邀请会员下级邀请数量
func (i *invitationManager) GetSubInvitationNum(memberIdArr []int64) map[int64]int {
	if memberIdArr == nil || len(memberIdArr) == 0 {
		return map[int64]int{}
	}
	return i.member.rep.GetSubInvitationNum(i.member.GetAggregateRootId(),
		memberIdArr)
}

// 获取邀请要的会员
func (i *invitationManager) GetInvitationMeMember() *member.Member {
	return i.member.rep.GetInvitationMeMember(i.member.GetAggregateRootId())
}
