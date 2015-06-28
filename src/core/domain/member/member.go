/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:12
 * description :
 * history :
 */

package member

import (
	"errors"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/valueobject"
	partnerImpl "go2o/src/core/domain/partner"
	"go2o/src/core/infrastructure/domain"
	"time"
)

var _ member.IMember = new(Member)

type Member struct {
	_value        *member.ValueMember
	_account      *member.Account
	_bank         *member.BankInfo
	_level        *valueobject.MemberLevel
	_rep          member.IMemberRep
	_relation     *member.MemberRelation
	_invitation   member.IInvitationManager
	_levelManager partner.ILevelManager
}

func NewMember(val *member.ValueMember, rep member.IMemberRep) member.IMember {
	return &Member{
		_value: val,
		_rep:   rep,
	}
}

// 获取聚合根编号
func (this *Member) GetAggregateRootId() int {
	return this._value.Id
}

// 获取值
func (this *Member) GetValue() member.ValueMember {
	return *this._value
}

// 设置值
func (this *Member) SetValue(v *member.ValueMember) error {
	this._value.Avatar = v.Avatar
	this._value.Address = v.Address
	this._value.Birthday = v.Birthday
	this._value.Email = v.Email
	this._value.LastLoginTime = v.LastLoginTime
	this._value.Phone = v.Email
	this._value.Pwd = v.Pwd
	this._value.Name = v.Name
	this._value.Sex = v.Sex
	this._value.RegFrom = v.RegFrom
	if len(this._value.InvitationCode) == 0 {
		this._value.InvitationCode = v.InvitationCode
	}
	return nil
}

// 邀请管理
func (this *Member) Invitation() member.IInvitationManager {
	if this._invitation == nil {
		this._invitation = &invitationManager{
			_member: this,
		}
	}
	return this._invitation
}

// 获取账户
func (this *Member) GetAccount() *member.Account {
	if this._account == nil {
		this._account = this._rep.GetAccount(this._value.Id)
	}
	return this._account
}

// 保护账户
func (this *Member) SaveAccount() error {
	a := this.GetAccount()
	a.MemberId = this._value.Id
	return this._rep.SaveAccount(a)
}

// 获取提现银行信息
func (this *Member) GetBank() member.BankInfo {
	if this._bank == nil {
		this._bank = this._rep.GetBankInfo(this._value.Id)
	}
	return *this._bank
}

// 保存提现银行信息
func (this *Member) SaveBank(v *member.BankInfo) error {
	this.GetBank()

	if this._bank == nil {
		this._bank = v
	} else {
		this._bank.Account = v.Account
		this._bank.AccountName = v.AccountName
		this._bank.Network = v.Network
		this._bank.State = v.State
	}
	this._bank.UpdateTime = time.Now().Unix()
	//this.bank.MemberId = this.value.Id
	return this._rep.SaveBankInfo(this._bank)
}

// 保存返现记录
func (this *Member) SaveIncomeLog(l *member.IncomeLog) error {
	l.MemberId = this._value.Id
	return this._rep.SaveIncomeLog(l)
}

// 保存积分记录
func (this *Member) SaveIntegralLog(l *member.IntegralLog) error {
	l.MemberId = this._value.Id
	return this._rep.SaveIntegralLog(l)
}

// 增加经验值
func (this *Member) AddExp(exp int) error {
	this._value.Exp += exp
	_, err := this.Save()

	//判断是否升级
	this.checkUpLevel()

	return err
}

// 获取等级管理
func (this *Member) getLevelManager() partner.ILevelManager {
	if this._levelManager == nil {
		rl := this.GetRelation()
		parterId := rl.RegisterPartnerId
		this._levelManager = partnerImpl.NewLevelManager(parterId, this._rep)
	}
	return this._levelManager

}

// 获取等级
func (this *Member) GetLevel() *valueobject.MemberLevel {
	if this._level == nil {
		this._level = this.getLevelManager().GetLevelByValue(this._value.Level)
	}
	return this._level
}

//　增加积分
// todo:partnerId 不需要
func (this *Member) AddIntegral(partnerId int, backType int,
	integral int, log string) error {

	inLog := &member.IntegralLog{
		PartnerId:  partnerId,
		MemberId:   this._value.Id,
		Type:       backType,
		Integral:   integral,
		Log:        log,
		RecordTime: time.Now().Unix(),
	}

	err := this._rep.SaveIntegralLog(inLog)
	if err == nil {
		acc := this.GetAccount()
		acc.Integral = acc.Integral + integral
		err = this.SaveAccount()
	}

	return err
}

// 检查升级
func (this *Member) checkUpLevel() {
	levelId := this.getLevelManager().GetLevelValueByExp(this._value.Exp)
	if levelId != 0 && this._value.Level < levelId {
		this._value.Level = levelId
		this.Save()
		this._level = nil
	}
}

// 获取会员关联
func (this *Member) GetRelation() *member.MemberRelation {
	if this._relation == nil {
		this._relation = this._rep.GetRelation(this._value.Id)
	}
	return this._relation
}

// 保存
func (this *Member) Save() (int, error) {

	if this._value.Id > 0 {
		return this._rep.SaveMember(this._value)
	}

	return this.create(this._value)
}

// 锁定会员
func (this *Member) Lock() error {
	return this._rep.LockMember(this.GetAggregateRootId(), 0)
}

// 解锁会员
func (this *Member) Unlock() error {
	return this._rep.LockMember(this.GetAggregateRootId(), 1)
}

// 修改密码,旧密码可为空
func (this *Member) ModifyPassword(newPwd, oldPwd string) error {
	var err error

	if b, err := domain.ChkPwdRight(newPwd); !b {
		return err
	}

	if len(oldPwd) != 0 {
		dyp := domain.Md5MemberPwd(this._value.Usr, oldPwd)
		if dyp != this._value.Pwd {
			return errors.New("原密码不正确")
		}
	}

	this._value.Pwd = domain.Md5MemberPwd(this._value.Usr, newPwd)
	_, err = this.Save()

	return err
}

// 创建会员
func (this *Member) create(m *member.ValueMember) (int, error) {
	if this.UsrIsExist() {
		return -1, errors.New("用户名已经被使用")
	}
	t := time.Now().Unix()
	m.State = 1
	m.RegTime = t
	m.LastLoginTime = t
	m.Level = 1
	m.Exp = 1
	m.Avatar = "resource/no_avatar.gif"
	m.Birthday = "1970-01-01"
	m.DynamicToken = m.Pwd
	m.Exp = 0

	if len(m.RegFrom) == 0 {
		m.RegFrom = "API-INTERNAL"
	}

	// 如果昵称为空，则跟用户名相同
	if len(m.Name) == 0 {
		m.Name = m.Usr
	}
	m.InvitationCode = this.generateInvitationCode() // 创建一个邀请码

	id, err := this._rep.SaveMember(m)
	if id != 0 {
		this._value.Id = id
	}
	return id, err
}

// 创建邀请码
func (this *Member) generateInvitationCode() string {
	var code string
	for {
		code = domain.GenerateInvitationCode()
		if memberId := this._rep.GetMemberIdByInvitationCode(code); memberId == 0 {
			break
		}
	}
	return code
}

// 用户是否已经存在
func (this *Member) UsrIsExist() bool {
	return this._rep.CheckUsrExist(this._value.Usr)
}

// 创建并初始化
func (this *Member) SaveRelation(r *member.MemberRelation) error {
	this._relation = r
	this._relation.MemberId = this._value.Id
	return this._rep.SaveRelation(this._relation)
}

// 创建配送地址
func (this *Member) CreateDeliver(v *member.DeliverAddress) member.IDeliver {
	return newDeliver(v, this._rep)
}

// 获取配送地址
func (this *Member) GetDeliverAddress() []member.IDeliver {
	var vls []*member.DeliverAddress
	vls = this._rep.GetDeliverAddress(this.GetAggregateRootId())
	var arr []member.IDeliver = make([]member.IDeliver, len(vls))
	for i, v := range vls {
		arr[i] = this.CreateDeliver(v)
	}
	return arr
}

// 获取配送地址
func (this *Member) GetDeliver(deliverId int) member.IDeliver {
	v := this._rep.GetSingleDeliverAddress(this.GetAggregateRootId(), deliverId)
	if v != nil {
		return this.CreateDeliver(v)
	}
	return nil
}

// 删除配送地址
func (this *Member) DeleteDeliver(deliverId int) error {
	return this._rep.DeleteDeliver(this.GetAggregateRootId(), deliverId)
}
