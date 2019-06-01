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

// 获取推荐数组
func (i *invitationManager) InviterArray(memberId int64, depth int) []int64 {
	arr := make([]int64, depth)
	var di int
	inviterId := memberId
	for di <= depth-1 {
		rl := i.member.repo.GetRelation(inviterId)
		if rl == nil || rl.InviterId <= 0 {
			break
		}
		arr[di] = rl.InviterId
		inviterId = arr[di]
		di++
	}
	return arr
}

// 判断是否推荐了某个会员
func (i *invitationManager) InvitationBy(memberId int64) bool {
	rl := i.member.GetRelation()
	if rl != nil {
		return rl.InviterId == memberId
	}
	return false
}

// 获取我邀请的会员
func (i *invitationManager) GetInvitationMembers(begin, end int) (
	int, []*dto.InvitationMember) {
	return i.member.repo.GetMyInvitationMembers(
		i.member.GetAggregateRootId(), begin, end)
}

// 获取我的邀请码
func (i *invitationManager) MyCode() string {
	return i.member.GetValue().InvitationCode
}

// 获取邀请会员下级邀请数量
func (i *invitationManager) GetSubInvitationNum(memberIdArr []int32) map[int32]int {
	if memberIdArr == nil || len(memberIdArr) == 0 {
		return map[int32]int{}
	}
	return i.member.repo.GetSubInvitationNum(i.member.GetAggregateRootId(),
		memberIdArr)
}

// 获取邀请要的会员
func (i *invitationManager) GetInvitationMeMember() *member.Member {
	return i.member.repo.GetInvitationMeMember(i.member.GetAggregateRootId())
}
