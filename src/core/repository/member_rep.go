/**
 * Copyright 2014 @ z3q.net.
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
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/log"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/valueobject"
	memberImpl "go2o/src/core/domain/member"
	"go2o/src/core/variable"
)

var _ member.IMemberRep = new(MemberRep)

type MemberRep struct {
	db.Connector
	_partnerRep partner.IPartnerRep
}

func NewMemberRep(c db.Connector) *MemberRep {
	return &MemberRep{
		Connector: c,
	}
}

func (this *MemberRep) SetPartnerRep(partnerRep partner.IPartnerRep) {
	this._partnerRep = partnerRep
}

// 根据用户名获取会员
func (this *MemberRep) GetMemberValueByUsr(usr string) *member.ValueMember {
	e := &member.ValueMember{}
	err := this.Connector.GetOrm().GetBy(e, "usr=?", usr)
	if err != nil {
		return nil
	}
	return e
}

// 根据手机号码获取会员
func (this *MemberRep) GetMemberValueByPhone(phone string) *member.ValueMember {
	e := &member.ValueMember{}
	err := this.Connector.GetOrm().GetBy(e, "phone=?", phone)
	if err != nil {
		return nil
	}
	return e
}

// 获取会员
func (this *MemberRep) GetMember(memberId int) member.IMember {
	e := &member.ValueMember{}
	err := this.Connector.GetOrm().Get(memberId, e)
	if err == nil {
		return this.CreateMember(e)
	}
	return nil
}

func (this *MemberRep) GetMemberIdByUser(user string) int {
	var id int
	this.Connector.ExecScalar("SELECT id FROM mm_member WHERE usr = ?", &id, user)
	return id
}

// 创建会员
func (this *MemberRep) CreateMember(v *member.ValueMember) member.IMember {
	return memberImpl.NewMember(v, this, this._partnerRep)
}

// 根据邀请码获取会员编号
func (this *MemberRep) GetMemberIdByInvitationCode(code string) int {
	var memberId int
	this.ExecScalar("SELECT id FROM mm_member WHERE invitation_code=?", &memberId, code)
	return memberId
}

func (this *MemberRep) GetLevel(partnerId, levelValue int) *valueobject.MemberLevel {
	var m valueobject.MemberLevel
	err := this.Connector.GetOrm().GetBy(&m, "partner_id=? AND value = ?", partnerId, levelValue)
	if err != nil {
		return nil
	}
	return &m
}

// 获取下一个等级
func (this *MemberRep) GetNextLevel(partnerId, levelVal int) *valueobject.MemberLevel {
	var m valueobject.MemberLevel
	err := this.Connector.GetOrm().GetBy(&m, "partner_id=? AND value>? LIMIT 0,1", partnerId, levelVal)
	if err != nil {
		return nil
	}
	return &m
}

// 获取会员等级
func (this *MemberRep) GetMemberLevels(partnerId int) []*valueobject.MemberLevel {
	list := []*valueobject.MemberLevel{}
	this.Connector.GetOrm().Select(&list,
		"partner_id=?", partnerId)
	return list
}

// 删除会员等级
func (this *MemberRep) DeleteMemberLevel(partnerId, id int) error {
	_, err := this.Connector.GetOrm().Delete(&valueobject.MemberLevel{},
		"id=? AND partner_id=?", id, partnerId)
	return err
}

// 保存等级
func (this *MemberRep) SaveMemberLevel(partnerId int, v *valueobject.MemberLevel) (int, error) {
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

// 获取会员最后更新时间
func (this *MemberRep) GetMemberLatestUpdateTime(id int) int64 {
	var updateTime int64
	this.Connector.ExecScalar(`SELECT update_time FROM mm_member where id=?`, &updateTime, id)
	return updateTime
}

// 获取账户
func (this *MemberRep) GetAccount(memberId int) *member.AccountValue {
	e := new(member.AccountValue)
	this.Connector.GetOrm().Get(memberId, e)
	return e
}

// 保存账户，传入会员编号
func (this *MemberRep) SaveAccount(a *member.AccountValue) (int, error) {
	_, _, err := this.Connector.GetOrm().Save(a.MemberId, a)
	return a.MemberId, err
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
func (this *MemberRep) GetLevelValueByExp(partnerId int, exp int) int {
	var levelId int
	this.Connector.ExecScalar(`SELECT lv.value FROM pt_member_level lv
	 	where lv.partner_id=? AND lv.require_exp <= ? AND lv.enabled=1
	 	 ORDER BY lv.require_exp DESC LIMIT 0,1`,
		&levelId, partnerId, exp)
	return levelId

}

// 锁定会员
func (this *MemberRep) LockMember(id int, state int) error {
	_, err := this.Connector.ExecNonQuery("update mm_member set state=? WHERE id=?", state, id)
	return err
}

// 保存会员
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

	// 更新会员数
	var total = 0
	this.Connector.ExecScalar("SELECT COUNT(0) FROM mm_member", &total)
	gof.CurrentApp.Storage().Set(variable.KvTotalMembers, total)

	return id, err
}

func (this *MemberRep) getLatestId() int {
	var id int
	this.Connector.ExecScalar("SELECT MAX(id) FROM mm_member", &id)
	return id
}

func (this *MemberRep) initMember(id int, v *member.ValueMember) {

	orm := this.Connector.GetOrm()
	orm.Save(nil, &member.AccountValue{
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
		MemberId:          id,
		CardId:            "",
		RefereesId:        0,
		RegisterPartnerId: 0,
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
func (this *MemberRep) GetDeliverAddress(memberId int) []*member.DeliverAddress {
	addresses := []*member.DeliverAddress{}
	this.Connector.GetOrm().Select(&addresses, "member_id=?", memberId)
	return addresses
}

// 获取配送地址
func (this *MemberRep) GetSingleDeliverAddress(memberId, deliverId int) *member.DeliverAddress {
	var address member.DeliverAddress
	err := this.Connector.GetOrm().Get(deliverId, &address)

	if err == nil && address.MemberId == memberId {
		return &address
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
	this.Connector.GetOrm().SelectByQuery(&arr,
		"SELECT * FROM mm_member WHERE id IN (SELECT member_id FROM mm_relation WHERE invi_member_id=?)", memberId)
	return arr
}

// 获取下级会员数量
func (this *MemberRep) GetSubInvitationNum(memberIds string) map[int]int {
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
func (this *MemberRep) GetInvitationMeMember(memberId int) *member.ValueMember {
	var d *member.ValueMember = new(member.ValueMember)
	err := this.Connector.GetOrm().GetByQuery(d,
		"SELECT * FROM mm_member WHERE id =(SELECT invi_member_id FROM mm_relation  WHERE id=?)",
		memberId)

	if err != nil {
		return nil
	}
	return d
}

// 根据编号获取余额变动信息
func (this *MemberRep) GetBalanceInfo(id int) *member.BalanceInfoValue {
	var e member.BalanceInfoValue
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 根据号码获取余额变动信息
func (this *MemberRep) GetBalanceInfoByNo(tradeNo string) *member.BalanceInfoValue {
	var e member.BalanceInfoValue
	if err := this.Connector.GetOrm().GetBy(&e, "trade_no=?", tradeNo); err == nil {
		return &e
	}
	return nil
}

// 保存余额变动信息
func (this *MemberRep) SaveBalanceInfo(v *member.BalanceInfoValue) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM mm_balance_info WHERE member_id=?", &v.Id, v.MemberId)
	}
	return v.Id, err
}
