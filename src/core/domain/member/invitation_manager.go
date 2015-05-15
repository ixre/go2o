/**
 * Copyright 2015 @ S1N1 Team.
 * name : invitation_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package member

import (
	"go2o/src/core/domain/interface/member"
	"strconv"
	"strings"
)

var _ member.IInvitationManager = new(invitationManager)

type invitationManager struct {
	_member       *Member
	_myInvMembers []*member.ValueMember
}

// 获取我邀请的会员
func (this *invitationManager) GetMyInvitationMembers() []*member.ValueMember {
	this._myInvMembers = this._member._rep.GetMyInvitationMembers(this._member.GetAggregateRootId())
	return this._myInvMembers
}

// 获取我的邀请码
func (this *invitationManager) GetMyInvitationCode() string {
	return this._member.GetValue().InvitationCode
}

// 获取邀请会员下级邀请数量
func (this *invitationManager) GetSubInvitationNum() map[int]int {
	if this._myInvMembers == nil {
		this._myInvMembers = this.GetMyInvitationMembers()
	}

	var ids []string = make([]string, len(this._myInvMembers))
	for i, v := range this._myInvMembers {
		ids[i] = strconv.Itoa(v.Id)
	}

	return this._member._rep.GetSubInvitationNum(strings.Join(ids, ","))
}

// 获取邀请要的会员
func (this *invitationManager) GetInvitationMeMember() *member.ValueMember {
	return this._member._rep.GetInvitationMeMember(this._member.GetAggregateRootId())
}
