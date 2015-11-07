/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 10:12
 * description :
 * history :
 */

package member

//todo: 要注意UpdateTime的更新

import (
	"errors"
	"fmt"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/valueobject"
	partnerImpl "go2o/src/core/domain/partner"
	"go2o/src/core/infrastructure/domain"
	"regexp"
	"strings"
	"time"
)

var _ member.IMember = new(Member)

type Member struct {
	_value        *member.ValueMember
	_account      member.IAccount
	_bank         *member.BankInfo
	_level        *valueobject.MemberLevel
	_rep          member.IMemberRep
	_partnerRep   partner.IPartnerRep
	_relation     *member.MemberRelation
	_invitation   member.IInvitationManager
	_levelManager partner.ILevelManager
}

func NewMember(val *member.ValueMember, rep member.IMemberRep, partnerRep partner.IPartnerRep) member.IMember {
	return &Member{
		_value:      val,
		_rep:        rep,
		_partnerRep: partnerRep,
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

var (
	userRegex  = regexp.MustCompile("^[a-zA-Z0-9_]{6,}$")
	emailRegex = regexp.MustCompile("\\w+([-+.']\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*")
	phoneRegex = regexp.MustCompile("^(13[0-9]|15[0|1|2|3|4|5|6|8|9]|18[0|1|2|3|5|6|7|8|9]|17[0|6|7|8])(\\d{8})$")
	qqRegex    = regexp.MustCompile("^\\d{5,12}$")
)

func (this *Member) validate(v *member.ValueMember) error {
	v.Usr = strings.ToLower(strings.TrimSpace(v.Usr)) // 小写并删除空格
	v.Name = strings.TrimSpace(v.Name)
	v.Email = strings.ToLower(strings.TrimSpace(v.Email))
	v.Phone = strings.TrimSpace(v.Phone)

	if len([]rune(v.Usr)) < 6 {
		return member.ErrUserLength
	}
	if !userRegex.MatchString(v.Usr) {
		return member.ErrUserValidErr
	}

	if len([]rune(v.Name)) < 2 {
		return member.ErrPersonName
	}

	if len(v.Email) != 0 && !emailRegex.MatchString(v.Email) {
		return member.ErrEmailValidErr
	}
	if len(v.Phone) != 0 && !phoneRegex.MatchString(v.Phone) {
		return member.ErrPhoneValidErr
	}
	if len(v.Qq) != 0 && !qqRegex.MatchString(v.Qq) {
		return member.ErrQqValidErr
	}
	return nil
}

// 设置值
func (this *Member) SetValue(v *member.ValueMember) error {
	v.Usr = this._value.Usr
	if err := this.validate(v); err != nil {
		return err
	}
	this._value.Address = v.Address
	this._value.Birthday = v.Birthday
	this._value.Qq = v.Qq
	this._value.Email = v.Email
	this._value.LastLoginTime = v.LastLoginTime
	this._value.Phone = v.Phone
	this._value.Pwd = v.Pwd
	this._value.Name = v.Name
	this._value.Sex = v.Sex
	this._value.RegFrom = v.RegFrom
	if v.Avatar != "" {
		this._value.Avatar = v.Avatar
	}
	if len(this._value.InvitationCode) == 0 {
		this._value.InvitationCode = v.InvitationCode
	}

	if v.Exp != 0 {
		this._value.Exp = v.Exp
	}

	if v.Level > 0 {
		this._value.Level = v.Level
	}

	if len(v.TradePwd) == 0 {
		this._value.TradePwd = v.TradePwd
	}

	if len(this._value.Qq) != 0 && len(this._value.Email) != 0 &&
		len(this._value.Birthday) != 0 && len(this._value.Address) != 0 &&
		len(this._value.Phone) != 0 && len(this._value.Avatar) != 0 &&
		this._value.Sex != 0 {
		this.notifyOnProfileComplete()
	}
	return nil
}

func (this *Member) notifyOnProfileComplete() {
	rl := this.GetRelation()
	pt, err := this._partnerRep.GetPartner(rl.RegisterPartnerId)
	if err == nil {
		key := fmt.Sprintf("profile:complete:id_%d", this.GetAggregateRootId())
		if pt.MemberKvManager().GetInt(key) == 0 {
			if err := this.sendNotifyMail(pt); err == nil {
				pt.MemberKvManager().Set(key, "1")
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func (this *Member) sendNotifyMail(pt partner.IPartner) error {
	tplId := pt.KvManager().GetInt(partner.KeyMssTplIdOfProfileComplete)
	if tplId > 0 {
		mailTpl := pt.MssManager().GetMailTemplate(tplId)
		if mailTpl != nil {
			tpl, err := pt.MssManager().CreateMsgTemplate(mailTpl)
			if err != nil {
				return err
			}

			//todo:?? data
			var data = map[string]string{
				"Name":           this._value.Name,
				"InvitationCode": this._value.InvitationCode,
			}

			return pt.MssManager().Send(tpl, data, []string{this._value.Email})
		}
	}
	return errors.New("no such email template")
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
func (this *Member) GetAccount() member.IAccount {
	if this._account == nil {
		v := this._rep.GetAccount(this._value.Id)
		return NewAccount(v, this._rep)
	}
	return this._account
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
		this._bank.Name = v.Name
	}
	this._bank.UpdateTime = time.Now().Unix()
	//this._bank.MemberId = this.value.Id
	return this._rep.SaveBankInfo(this._bank)
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
		partnerId := rl.RegisterPartnerId
		this._levelManager = partnerImpl.NewLevelManager(partnerId, this._rep)
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
		acv := acc.GetValue()
		acv.Integral += integral
		_, err = acc.Save()
	}
	return err
}

// 检查升级
func (this *Member) checkUpLevel() bool {
	levelValue := this.getLevelManager().GetLevelValueByExp(this._value.Exp)
	if levelValue != 0 && this._value.Level < levelValue {
		this._value.Level = levelValue
		this.Save()
		this._level = nil
		return true
	}
	return false
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
	this._value.UpdateTime = time.Now().Unix() // 更新时间，数据以更新时间触发
	if this._value.Id > 0 {
		return this._rep.SaveMember(this._value)
	}

	if err := this.validate(this._value); err != nil {
		return this.GetAggregateRootId(), err
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
	if newPwd == oldPwd {
		return member.ErrPwdCannotSame
	}

	if b, err := domain.ChkPwdRight(newPwd); !b {
		return err
	}

	if len(oldPwd) != 0 && oldPwd != this._value.Pwd {
		return member.ErrPwdOldPwdNotRight

	}

	this._value.Pwd = newPwd
	_, err = this.Save()

	return err
}

// 修改交易密码，旧密码可为空
func (this *Member) ModifyTradePassword(newPwd, oldPwd string) error {
	var err error
	if newPwd == oldPwd {
		return member.ErrPwdCannotSame
	}

	if b, err := domain.ChkPwdRight(newPwd); !b {
		return err
	}

	// 已经设置过旧密码
	if len(this._value.TradePwd) != 0 && this._value.TradePwd != oldPwd {
		return member.ErrPwdOldPwdNotRight
	}

	this._value.TradePwd = newPwd
	_, err = this.Save()

	return err
}

// 创建会员
func (this *Member) create(m *member.ValueMember) (int, error) {
	if this.UsrIsExist() {
		return -1, errors.New("用户名已经被使用")
	}
	if len(m.Phone) > 0 && this.PhoneIsExist(m.Phone) {
		return -1, member.ErrPhoneHasBind
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
	return this._rep.CheckUsrExist(this._value.Usr, this.GetAggregateRootId())
}

// 手机号码是否占用
func (this *Member) PhoneIsExist(phone string) bool {
	return this._rep.CheckPhoneBind(this._value.Usr, this.GetAggregateRootId())
}

// 创建并初始化
func (this *Member) SaveRelation(r *member.MemberRelation) error {
	this._relation = r
	this._relation.MemberId = this._value.Id
	return this._rep.SaveRelation(this._relation)
}

// 创建配送地址
func (this *Member) CreateDeliver(v *member.DeliverAddress)(member.IDeliver,error) {
	return newDeliver(v, this._rep)
}

// 获取配送地址
func (this *Member) GetDeliverAddress() []member.IDeliver {
	var vls []*member.DeliverAddress
	vls = this._rep.GetDeliverAddress(this.GetAggregateRootId())
	var arr []member.IDeliver = make([]member.IDeliver, len(vls))
	for i, v := range vls {
		arr[i],_  = this.CreateDeliver(v)
	}
	return arr
}

// 获取配送地址
func (this *Member) GetDeliver(deliverId int) member.IDeliver {
	v := this._rep.GetSingleDeliverAddress(this.GetAggregateRootId(), deliverId)
	if v != nil {
		d,_ := this.CreateDeliver(v)
		return d
	}
	return nil
}

// 删除配送地址
func (this *Member) DeleteDeliver(deliverId int) error {
	return this._rep.DeleteDeliver(this.GetAggregateRootId(), deliverId)
}
