/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-09 10:12
 * description :
 * history :
 */

package member

import (
	"com/domain/interface/member"
	"errors"
	"time"
)

type Member struct {
	value    *member.ValueMember
	account  *member.Account
	bank     *member.BankInfo
	rep      member.IMemberRep
	relation *member.MemberRelation
}

func NewMember(val *member.ValueMember, rep member.IMemberRep) member.IMember {
	return &Member{
		value: val,
		rep:   rep,
	}
}

func (this *Member) GetValue() member.ValueMember {
	return *this.value
}

func (this *Member) SetValue(v *member.ValueMember) error {
	this.value.Avatar = v.Avatar
	this.value.Address = v.Address
	this.value.Birthday = v.Birthday
	this.value.Email = v.Email
	this.value.LastLoginTime = v.LastLoginTime
	this.value.Phone = v.Email
	this.value.Pwd = v.Pwd
	this.value.Name = v.Name
	this.value.Sex = v.Sex
	return nil
}

func (this *Member) GetAccount() *member.Account {
	if this.account == nil {
		this.account = this.rep.GetAccount(this.value.Id)
	}
	return this.account
}
func (this *Member) SaveAccount() error {
	a := this.GetAccount()
	a.MemberId = this.value.Id
	return this.rep.SaveAccount(a)
}

// 获取提现银行信息
func (this *Member) GetBank() member.BankInfo {
	if this.bank == nil {
		this.bank = this.rep.GetBankInfo(this.value.Id)
	}
	return *this.bank
}

// 保存提现银行信息
func (this *Member) SaveBank(v *member.BankInfo) error {
	this.bank = v
	return this.rep.SaveBankInfo(this.bank)
}

func (this *Member) SaveIncomeLog(l *member.IncomeLog) error {
	l.MemberId = this.value.Id
	return this.rep.SaveIncomeLog(l)
}
func (this *Member) SaveIntegralLog(l *member.IntegralLog) error {
	l.MemberId = this.value.Id
	return this.rep.SaveIntegralLog(l)
}

// 增加经验值
func (this *Member) AddExp(exp int) error {
	this.value.Exp += exp
	_, err := this.Save()
	return err
}

//　增加积分
// todo:partnerId 不需要
func (this *Member) AddIntegral(partnerId int, backType int,
	integral int, log string) error {

	inteLog := &member.IntegralLog{
		PtId:       partnerId,
		MemberId:   this.value.Id,
		Type:       backType,
		Integral:   integral,
		Log:        log,
		RecordTime: time.Now(),
	}

	err := this.rep.SaveIntegralLog(inteLog)
	if err == nil {
		acc := this.GetAccount()
		acc.Integral = acc.Integral + integral
		err = this.SaveAccount()
	}

	//判断是否升级
	this.checkLevel()

	return err
}

func (this *Member) checkLevel() {
	levelId := this.rep.GetLevelByExp(this.value.Exp)
	if levelId != 0 && this.value.Level < levelId {
		this.value.Level = levelId
		this.Save()
	}
}

// 获取会员关联
func (this *Member) GetRelation() *member.MemberRelation {
	if this.relation == nil {
		this.relation = this.rep.GetRelation(this.value.Id)
	}
	return this.relation
}

// 保存
func (this *Member) Save() (int, error) {

	if this.value.Id > 0 {
		return this.rep.SaveMember(this.value)
	}

	return this.create(this.value)
}

// 创建会员
func (this *Member) create(m *member.ValueMember) (int, error) {
	if this.UsrIsExist() {
		return -1, errors.New("用户名已经被使用")
	}
	t := time.Now()
	m.State = 1
	m.RegTime = t
	m.LastLoginTime = t
	m.Level = 1
	m.Avatar = "share/noavatar.gif"
	m.Birthday = "1970-01-01"
	m.LoginToken = m.Pwd
	m.Exp = 0

	id, err := this.rep.SaveMember(m)
	if id != 0 {
		this.value.Id = id
	}
	return id, err
}

// 用户是否已经存在
func (this *Member) UsrIsExist() bool {
	return this.rep.CheckUsrExist(this.value.Usr)
}

// 创建并初始化
func (this *Member) SaveRelation(r *member.MemberRelation) error {
	this.relation = r
	this.relation.MemberId = this.value.Id
	return this.rep.SaveRelation(this.relation)
}
