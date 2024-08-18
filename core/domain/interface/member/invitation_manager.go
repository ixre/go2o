/**
 * Copyright 2015 @ 56x.net.
 * name : invitaton_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package member

import "github.com/ixre/go2o/core/dto"

type IInvitationManager interface {
	// 获取邀请人数组
	InviterArray(memberId int64, depth int) []int64

	// 判断是否由会员邀请
	InvitationBy(memberId int64) bool

	// 获取我邀请的会员
	GetInvitationMembers(begin, end int) (total int, rows []*dto.InvitationMember)

	// 获取邀请会员下级邀请数量
	GetSubInvitationNum(memberIdArr []int32) map[int32]int

	// 获取邀请我的会员
	GetInvitationMeMember() *Member

	// 更换邀请人,async是否异步更新
	UpdateInviter(inviterId int64, sync bool) error

	// 屏蔽
	Shield(memberId int) error
	// 取消屏蔽
	UnShield(memberId int) error
	// 拉黑
	Block(memberId int) error
	// 取消拉黑
	UnBlock(memberId int) error
	// 是否被屏蔽或拉黑
	IsBlockOrShield(memberId int) (bool, int)
}

const (
	// 屏蔽
	BlockFlagShield = 1 << iota
	// 拉黑
	BlockFlagBlack = 2
)

// MmBlockList 会员拉黑列表
type BlockList struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 会员编号
	MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" bson:"memberId"`
	// 拉黑会员编号
	BlockMemberId int `json:"blockMemberId" db:"block_member_id" gorm:"column:block_member_id" bson:"blockMemberId"`
	// 拉黑标志，1: 屏蔽  2: 拉黑
	BlockFlag int `json:"blockFlag" db:"block_flag" gorm:"column:block_flag" bson:"blockFlag"`
	// 拉黑时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (m BlockList) TableName() string {
	return "mm_block_list"
}
