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
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/dto"
)

var _ member.IInvitationManager = new(invitationManager)

type invitationManager struct {
	member      *memberImpl
	_memberRepo member.IMemberRepo
}

func (i *invitationManager) getBlockInfo(memberId int) *member.BlockList {
	return i._memberRepo.BlockRepo().FindBy("member_id = ? AND block_member_id = ?",
		i.member.GetAggregateRootId(), memberId)
}

// Block implements member.IInvitationManager.
func (i *invitationManager) Block(memberId int) error {
	v := i.getBlockInfo(memberId)
	if v != nil {
		v.BlockFlag |= member.BlockFlagBlack
	} else {
		v = &member.BlockList{
			Id:            memberId,
			MemberId:      int(i.member.GetAggregateRootId()),
			BlockMemberId: memberId,
			BlockFlag:     member.BlockFlagBlack,
			CreateTime:    int(time.Now().Unix()),
		}
	}
	_, err := i._memberRepo.BlockRepo().Save(v)
	return err
}

// IsBlockOrShield implements member.IInvitationManager.
func (i *invitationManager) IsBlockOrShield(memberId int) (bool, int) {
	v := i.getBlockInfo(memberId)
	if v == nil {
		return false, 0
	}
	return true, v.BlockFlag
}

// Shield implements member.IInvitationManager.
func (i *invitationManager) Shield(memberId int) error {
	v := i.getBlockInfo(memberId)
	if v != nil {
		v.BlockFlag |= member.BlockFlagShield
	} else {
		v = &member.BlockList{
			Id:            memberId,
			MemberId:      int(i.member.GetAggregateRootId()),
			BlockMemberId: memberId,
			BlockFlag:     member.BlockFlagShield,
			CreateTime:    int(time.Now().Unix()),
		}
	}
	_, err := i._memberRepo.BlockRepo().Save(v)
	return err
}

// UnBlock implements member.IInvitationManager.
func (i *invitationManager) UnBlock(memberId int) (err error) {
	v := i.getBlockInfo(memberId)
	if v != nil {
		v.BlockFlag ^= member.BlockFlagBlack
		if v.BlockFlag == 0 {
			err = i._memberRepo.BlockRepo().Delete(v)
		} else {
			_, err = i._memberRepo.BlockRepo().Save(v)
		}
	}
	return err
}

// UnShield implements member.IInvitationManager.
func (i *invitationManager) UnShield(memberId int) (err error) {
	v := i.getBlockInfo(memberId)
	if v != nil {
		v.BlockFlag ^= member.BlockFlagBlack
		if v.BlockFlag == 0 {
			err = i._memberRepo.BlockRepo().Delete(v)
		} else {
			_, err = i._memberRepo.BlockRepo().Save(v)
		}
	}
	return err
}

// 更换邀请人
func (i *invitationManager) UpdateInviter(inviterId int, sync bool) error {
	id := i.member.GetAggregateRootId()
	var rl *member.InviteRelation
	if inviterId > 0 {
		rl = i.member.repo.GetRelation(int64(inviterId))
	}
	// 判断邀请人是否为下级的被邀请会员
	if i.checkInvitation(int64(inviterId), int64(id)) {
		return member.ErrInvalidInviteLevel
	}
	if !sync {
		return i.walkUpdateInvitation(int64(id), rl)
	}
	// 异步更新
	go i.walkUpdateInvitation(int64(id), rl)
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
		arr := m.Invitation().InviterArray(int64(r.InviterId), 2)
		r.InviterD2 = int(arr[0])
		r.InviterD3 = int(arr[1])
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
		arr[di] = int64(rl.InviterId)
		inviterId = arr[di]
		di++
	}
	return arr
}

// 判断是否直接推荐了某个会员
func (i *invitationManager) InvitationBy(memberId int64) bool {
	rl := i.member.GetRelation()
	if rl != nil {
		return int(rl.InviterId) == int(memberId)
	}
	return false
}

// 获取我邀请的会员
func (i *invitationManager) GetInvitationMembers(begin, end int) (
	int, []*dto.InvitationMember) {
	return i.member.repo.GetMyInvitationMembers(
		int64(i.member.GetAggregateRootId()), begin, end)
}

// 获取邀请会员下级邀请数量
func (i *invitationManager) GetSubInvitationNum(memberIdArr []int32) map[int32]int {
	if memberIdArr == nil || len(memberIdArr) == 0 {
		return map[int32]int{}
	}
	return i.member.repo.GetSubInvitationNum(int64(i.member.GetAggregateRootId()),
		memberIdArr)
}

// 获取邀请要的会员
func (i *invitationManager) GetInvitationMeMember() *member.Member {
	return i.member.repo.GetInvitationMeMember(int64(i.member.GetAggregateRootId()))
}

// 是否存在邀请关系
func (i *invitationManager) checkInvitation(inviterId int64, id int64) bool {
	currId := id
	for {
		arr := i.InviterArray(currId, 1)
		if currId == arr[0] {
			return true
		}
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
