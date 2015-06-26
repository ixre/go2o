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
	"database/sql"
	"fmt"
	"github.com/atnet/gof/db"
	"github.com/atnet/gof/log"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/valueobject"
	memberImpl "go2o/src/core/domain/member"
)

var _ member.IMemberRep = new(memberRep)

type memberRep struct {
	db.Connector
}

func NewMemberRep(c db.Connector) member.IMemberRep {
	return &memberRep{
		Connector: c,
	}
}

// 根据用户名获取会员
func (this *memberRep) GetMemberValueByUsr(usr string) *member.ValueMember {
	e := &member.ValueMember{}
	err := this.Connector.GetOrm().GetBy(e, "usr='"+usr+"'")
	if err != nil {
		return nil
	}
	return e
}

func (this *memberRep) GetMember(memberId int) (member.IMember, error) {
	e := &member.ValueMember{}
	err := this.Connector.GetOrm().Get(memberId, e)
	if err != nil {
		return nil, err
	}
	return memberImpl.NewMember(e, this), err
}

// 创建会员
func (this *memberRep) CreateMember(v *member.ValueMember) member.IMember {
	return memberImpl.NewMember(v, this)
}

// 根据邀请码获取会员编号
func (this *memberRep) GetMemberIdByInvitationCode(code string) int {
	var memberId int
	this.ExecScalar("SELECT id FROM mm_member WHERE invitation_code=?", &memberId, code)
	return memberId
}

func (this *memberRep) GetLevel(partnerId, levelValue int) *valueobject.MemberLevel {
	var m valueobject.MemberLevel
	err := this.Connector.GetOrm().GetBy(&m, "partner_id=? AND value = ?", partnerId, levelValue)
	if err != nil {
		return nil
	}
	return &m
}

// 获取下一个等级
func (this *memberRep) GetNextLevel(partnerId, levelVal int) *valueobject.MemberLevel {
	var m valueobject.MemberLevel
	err := this.Connector.GetOrm().GetBy(&m, "partner_id=? AND value>? LIMIT 0,1", partnerId, levelVal)
	if err != nil {
		return nil
	}
	return &m
}

// 获取会员等级
func (this *memberRep) GetMemberLevels(partnerId int) []*valueobject.MemberLevel {
	list := []*valueobject.MemberLevel{}
	this.Connector.GetOrm().Select(&list,
		"partner_id=?", partnerId)
	return list
}

// 删除会员等级
func (this *memberRep) DeleteMemberLevel(partnerId, id int) error {
	_, err := this.Connector.GetOrm().Delete(&valueobject.MemberLevel{},
		"id=? AND partner_id=?", id, partnerId)
	return err
}

// 保存等级
func (this *memberRep) SaveMemberLevel(partnerId int, v *valueobject.MemberLevel) (int, error) {
	orm := this.Connector.GetOrm()
	var err error
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar(`SELECT MAX(id) FROM pt_member_level`, &v.Id)
	}
	return v.Id, err
}

// 获取账户
func (this *memberRep) GetAccount(memberId int) *member.Account {
	e := new(member.Account)
	this.Connector.GetOrm().Get(memberId, e)
	return e
}

// 保存账户，传入会员编号
func (this *memberRep) SaveAccount(a *member.Account) error {
	_, _, err := this.Connector.GetOrm().Save(a.MemberId, a)
	return err
}

// 获取银行信息
func (this *memberRep) GetBankInfo(memberId int) *member.BankInfo {
	e := new(member.BankInfo)
	this.Connector.GetOrm().Get(memberId, e)
	return e
}

// 保存银行信息
func (this *memberRep) SaveBankInfo(v *member.BankInfo) error {
	var err error
	_, _, err = this.Connector.GetOrm().Save(v.MemberId, v)
	return err
}

// 保存返现记录
func (this *memberRep) SaveIncomeLog(l *member.IncomeLog) error {
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
func (this *memberRep) SaveIntegralLog(l *member.IntegralLog) error {
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
func (this *memberRep) GetRelation(memberId int) *member.MemberRelation {
	e := new(member.MemberRelation)
	if this.Connector.GetOrm().Get(memberId, e) != nil {
		return nil
	}
	return e
}

// 获取积分对应的等级
func (this *memberRep) GetLevelValueByExp(partnerId int,exp int) int {
	var levelId int
	this.Connector.ExecScalar(`SELECT lv.value FROM pt_member_level lv
	 	where lv.partner_id=? AND lv.require_exp <= ? AND lv.enabled=1
	 	 ORDER BY lv.require_exp DESC LIMIT 0,1`,
		&levelId,partnerId,exp)
	return levelId

}

func (this *memberRep) SaveMember(v *member.ValueMember) (int, error) {
	if v.Id > 0 {
		_, _, err := this.Connector.GetOrm().Save(v.Id, v)
		return v.Id, err
	}
	return this.createMember(v)
}

func (this *memberRep) createMember(v *member.ValueMember) (int, error) {

	_, _, err := this.Connector.GetOrm().Save(nil, v)
	if err != nil {
		return -1, err
	}

	id := this.getLatestId()

	this.initMember(id, v)

	return id, err
}

func (this *memberRep) getLatestId() int {
	var id int
	this.Connector.ExecScalar("SELECT MAX(id) FROM mm_member", &id)
	return id
}

func (this *memberRep) initMember(id int, v *member.ValueMember) {

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
func (this *memberRep) CheckUsrExist(usr string) bool {
	var c int
	this.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member WHERE usr=?", &c, usr)
	return c != 0
}

// 保存绑定
func (this *memberRep) SaveRelation(v *member.MemberRelation) error {
	_, _, err := this.Connector.GetOrm().Save(v.MemberId, v)
	return err
}

// 保存地址
func (this *memberRep) SaveDeliver(v *member.DeliverAddress) (int, error) {
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
func (this *memberRep) GetDeliverAddress(memberId int) []*member.DeliverAddress {
	addresses := []*member.DeliverAddress{}
	this.Connector.GetOrm().Select(&addresses, "member_id=?", memberId)
	return addresses
}

// 获取配送地址
func (this *memberRep) GetSingleDeliverAddress(memberId, deliverId int) *member.DeliverAddress {
	var addr member.DeliverAddress
	err := this.Connector.GetOrm().Get(deliverId, &addr)

	if err == nil && addr.MemberId == memberId {
		return &addr
	}
	return nil
}

// 删除配送地址
func (this *memberRep) DeleteDeliver(memberId, deliverId int) error {
	_, err := this.Connector.ExecNonQuery(
		"DELETE FROM mm_deliver_addr WHERE member_id=? AND id=?",
		memberId, deliverId)
	return err
}

// 邀请
func (this *memberRep) GetMyInvitationMembers(memberId int) []*member.ValueMember {
	arr := []*member.ValueMember{}
	this.Connector.GetOrm().SelectByQuery(&arr,
		"SELECT * FROM mm_member WHERE id IN (SELECT member_id FROM mm_relation WHERE invi_member_id=?)", memberId)
	return arr
}

// 获取下级会员数量
func (this *memberRep) GetSubInvitationNum(memberIds string) map[int]int {
	var d map[int]int = make(map[int]int)
	err := this.Connector.Query(fmt.Sprintf("SELECT r1.member_id,"+
		"(SELECT COUNT(0) FROM mm_relation r2 WHERE r2.invi_member_id=r1.member_id)"+
		"as num FROM mm_relation r1 WHERE r1.member_id IN(%s)", memberIds),
		func(rows *sql.Rows) {
			var id, num int
			for rows.Next() {
				rows.Scan(&id, &num)
				d[id] = num
			}
			rows.Close()
		})

	if err != nil {
		log.PrintErr(err)
	}
	return d
}

// 获取推荐我的人
func (this *memberRep) GetInvitationMeMember(memberId int) *member.ValueMember {
	var d *member.ValueMember = new(member.ValueMember)
	err := this.Connector.GetOrm().GetByQuery(d,
		"SELECT * FROM mm_member WHERE id =(SELECT invi_member_id FROM mm_relation  WHERE id=?)",
		memberId)

	if err != nil {
		return nil
	}
	return d
}
