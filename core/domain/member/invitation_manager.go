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
)

var _ member.IInvitationManager = new(invitationManager)

type invitationManager struct {
	_member       *memberImpl
	_myInvMembers []*member.ValueMember
}

// 判断是否推荐了某个会员
func (this *invitationManager) InvitationBy(memberId int) bool {
	rl := this._member.GetRelation()
	if rl != nil {
		return rl.RefereesId == memberId
	}
	return false
}

// 获取我邀请的会员
func (this *invitationManager) GetInvitationMembers(begin, end int) (
	int, []*member.ValueMember) {
	return this._member._rep.GetMyInvitationMembers(
		this._member.GetAggregateRootId(), begin, end)
}

// 获取我的邀请码
func (this *invitationManager) GetMyInvitationCode() string {
	return this._member.GetValue().InvitationCode
}

// 获取邀请会员下级邀请数量
func (this *invitationManager) GetSubInvitationNum(memberIdArr []int) map[int]int {
	return this._member._rep.GetSubInvitationNum(this._member.GetAggregateRootId(),
		memberIdArr)
}

// 获取邀请要的会员
func (this *invitationManager) GetInvitationMeMember() *member.ValueMember {
	return this._member._rep.GetInvitationMeMember(this._member.GetAggregateRootId())
}
