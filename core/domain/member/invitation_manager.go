/**
 * Copyright 2015 @ 56x.net.
 * name : invitation_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package member

import (
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/msq"
	"strconv"
	"time"
)

var _ member.IInvitationManager = new(invitationManager)

type invitationManager struct {
	member       *memberImpl
	myInvMembers []*member.Member
}

// 更换邀请人
func (i *invitationManager) UpdateInviter(inviterId int64, sync bool) error {
	id := i.member.GetAggregateRootId()
	var rl *member.InviteRelation
	if inviterId > 0 {
		rl = i.member.repo.GetRelation(inviterId)
	}
	// 判断邀请人是否为下别的被邀请会员
	if i.checkInvitation(id, inviterId) {
		return member.ErrInvalidInviter
	}
	if !sync {
		return i.walkUpdateInvitation(id, rl)
	}
	// 异步更新
	go i.walkUpdateInvitation(id, rl)
	return nil
}

// 递归修改邀请人
func (i *invitationManager) walkUpdateInvitation(id int64, p *member.InviteRelation) error {
	r := i.member.repo.GetRelation(id)
	if p == nil {
		r.InviterId = 0
		r.InviterD2 = 0
		r.InviterD3 = 0
	} else {
		r.InviterId = p.MemberId
		r.InviterD2 = p.InviterId
		r.InviterD3 = p.InviterD2
	}
	err := i.member.repo.SaveRelation(r)
	if err == nil {
		// 推送关系更新消息
		go msq.PushDelay(msq.MemberRelationUpdated, strconv.Itoa(int(r.MemberId)), 500)
		// 更新被邀请会员的邀请关系
		var idList = i.member.repo.GetInviteChildren(id)
		for idx, cid := range idList {
			i.walkUpdateInvitation(cid, r)
			if idx%5 == 0 {
				time.Sleep(time.Second / 10)
			}
		}
	}
	return err
}

// 更新邀请关系
func (m *memberImpl) updateDepthInvite(r *member.InviteRelation) error {
	if r.InviterId > 0 {
		arr := m.Invitation().InviterArray(r.InviterId, 2)
		r.InviterD2 = arr[0]
		r.InviterD3 = arr[1]
	} else {
		r.InviterD2 = 0
		r.InviterD3 = 0
	}
	err := m.repo.SaveRelation(r)
	if err == nil {

	}
	return err
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

// 判断是否直接推荐了某个会员
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
	return i.member.GetValue().InviteCode
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

//  是否存在邀请关系
func (i *invitationManager) checkInvitation(inviterId int64, id int64) bool {
	currId := id
	for {
		arr := i.InviterArray(currId, 1)
		currId = arr[0]
		if currId == inviterId {
			return true
		}
		if currId == 0 {
			break
		}
	}
	return false
}
