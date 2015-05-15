/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:13
 * description :
 * history :
 */

package repository

import (
	"fmt"
	"github.com/atnet/gof/db"
	"github.com/atnet/gof/log"
	"go2o/src/core/domain/interface/member"
	memberImpl "go2o/src/core/domain/member"
)

var _ member.IMemberRep = new(MemberRep)

type MemberRep struct {
	db.Connector
}

func NewMemberRep(c db.Connector) member.IMemberRep {
	return &MemberRep{
		Connector: c,
	}
}

// 根据用户名获取会员
func (this *MemberRep) GetMemberValueByUsr(usr string) *member.ValueMember {
	e := &member.ValueMember{}
	err := this.Connector.GetOrm().GetBy(e, "usr='"+usr+"'")
	if err != nil {
		return nil
	}
	return e
}

func (this *MemberRep) GetMember(memberId int) (member.IMember, error) {
	e := &member.ValueMember{}
	err := this.Connector.GetOrm().Get(memberId, e)
	if err != nil {
		return nil, err
	}
	return memberImpl.NewMember(e, this), err
}

// 创建会员
func (this *MemberRep) CreateMember(v *member.ValueMember) member.IMember {
	return memberImpl.NewMember(v, this)
}

// 根据邀请码获取会员编号
func (this *MemberRep) GetMemberIdByInvitationCode(code string) int {
	var memberId int
	this.ExecScalar("SELECT id FROM mm_member WHERE invitation_code=?", &memberId, code)
	return memberId
}

func (this *MemberRep) GetLevel(levelVal int) *member.MemberLevel {
	var m member.MemberLevel
	this.Connector.GetOrm().Get(levelVal, &m)
	return &m
}

// 获取下一个等级
func (this *MemberRep) GetNextLevel(levelVal int) *member.MemberLevel {
	var m member.MemberLevel
	err := this.Connector.GetOrm().GetBy(&m, fmt.Sprintf("value>%d LIMIT 0,1", levelVal))
	if err != nil {
		return nil
	}
	return &m
}

// 获取账户
func (this *MemberRep) GetAccount(memberId int) *member.Account {
	e := new(member.Account)
	this.Connector.GetOrm().Get(memberId, e)
	return e
}

// 保存账户，传入会员编号
func (this *MemberRep) SaveAccount(a *member.Account) error {
	_, _, err := this.Connector.GetOrm().Save(a.MemberId, a)
	return err
}

// 获取银行信息
func (this *MemberRep) GetBankInfo(memberId int) *member.BankInfo {
	e := new(member.BankInfo)
	this.Connector.GetOrm().Get(memberId, e)
	return e
}

// 保存银行信息
func (this *MemberRep) SaveBankInfo(v *member.BankInfo) error {
	var err error
	_, _, err = this.Connector.GetOrm().Save(v.MemberId, v)
	return err
}

// 保存返现记录
func (this *MemberRep) SaveIncomeLog(l *member.IncomeLog) error {
	orm := this.Connector.GetOrm()
	var err error
	if l.Id > 0 {
		_, _, err = orm.Save(l.Id, l)
	} else {
		_, _, err = orm.Save(nil, l)
	}
	return err
}

// 保存积分记录
func (this *MemberRep) SaveIntegralLog(l *member.IntegralLog) error {
	orm := this.Connector.GetOrm()
	var err error
	if l.Id > 0 {
		_, _, err = orm.Save(l.Id, l)
	} else {
		_, _, err = orm.Save(nil, l)
	}
	return err
}

// 获取会员关联
func (this *MemberRep) GetRelation(memberId int) *member.MemberRelation {
	e := new(member.MemberRelation)
	if this.Connector.GetOrm().Get(memberId, e) != nil {
		return nil
	}
	return e
}

// 获取积分对应的等级
func (this *MemberRep) GetLevelByExp(exp int) int {
	var levelId int

	this.Connector.ExecScalar(`SELECT id FROM conf_member_level
	 	where require_exp<=? AND enabled=1 ORDER BY require_exp DESC LIMIT 0,1`,
		&levelId, exp)
	return levelId

}

func (this *MemberRep) SaveMember(v *member.ValueMember) (int, error) {
	if v.Id > 0 {
		_, _, err := this.Connector.GetOrm().Save(v.Id, v)
		return v.Id, err
	}
	return this.createMember(v)
}

func (this *MemberRep) createMember(v *member.ValueMember) (int, error) {

	_, _, err := this.Connector.GetOrm().Save(nil, v)
	if err != nil {
		return -1, err
	}

	id := this.getLatestId()

	this.initMember(id, v)

	return id, err
}

func (this *MemberRep) getLatestId() int {
	var id int
	this.Connector.ExecScalar("SELECT MAX(id) FROM mm_member", &id)
	return id
}

func (this *MemberRep) initMember(id int, v *member.ValueMember) {

	orm := this.Connector.GetOrm()
	orm.Save(nil, &member.Account{
		MemberId:    id,
		Balance:     0,
		TotalFee:    0,
		TotalCharge: 0,
		TotalPay:    0,
		UpdateTime:  v.RegTime,
	})

	orm.Save(nil, &member.BankInfo{
		MemberId: id,
		State:    1,
	})

	orm.Save(nil, &member.MemberRelation{
		MemberId:           id,
		CardId:             "",
		InvitationMemberId: 0,
		RegisterPartnerId:  0,
	})
}

// 用户名是否存在
func (this *MemberRep) CheckUsrExist(usr string) bool {
	var c int
	this.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member WHERE usr=?", &c, usr)
	return c != 0
}

// 保存绑定
func (this *MemberRep) SaveRelation(v *member.MemberRelation) error {
	_, _, err := this.Connector.GetOrm().Save(v.MemberId, v)
	return err
}

// 保存地址
func (this *MemberRep) SaveDeliver(v *member.DeliverAddress) (int, error) {
	orm := this.Connector.GetOrm()
	if v.Id <= 0 {
		_, id, err := orm.Save(nil, v)
		return int(id), err
	} else {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	}
}

// 获取全部配送地址
func (this *MemberRep) GetDeliverAddrs(memberId int) []member.DeliverAddress {
	addresses := []member.DeliverAddress{}
	this.Connector.GetOrm().Select(&addresses, "member_id=?", memberId)
	return addresses
}

// 获取配送地址
func (this *MemberRep) GetDeliverAddr(memberId, deliverId int) *member.DeliverAddress {
	var addr member.DeliverAddress
	err := this.Connector.GetOrm().Get(deliverId, &addr)

	if err == nil && addr.MemberId == memberId {
		return &addr
	}
	return nil
}

// 删除配送地址
func (this *MemberRep) DeleteDeliver(memberId, deliverId int) error {
	_, err := this.Connector.ExecNonQuery(
		"DELETE FROM mm_deliver_addr WHERE member_id=? AND id=?",
		memberId, deliverId)
	return err
}

// 邀请
func (this *MemberRep) GetMyInvitationMembers(memberId int) []*member.ValueMember {
	arr := []*member.ValueMember{}
	err := this.Connector.GetOrm().SelectByQuery(&arr,
		"SELECT * FROM mm_relation WHERE invi_member_id=?", memberId)
	if err != nil {
		log.PrintErr(err)
		return nil
	}

	return arr
}

// 获取下级会员数量
func (this *MemberRep) GetSubInvitationNum(memberIds string) map[int]int {
	return nil
}

// 获取推荐我的人
func (this *MemberRep) GetInvitationMeMember(memberId int) *member.ValueMember {
	return nil
}
