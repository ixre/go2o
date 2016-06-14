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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/mss"
	"go2o/core/infrastructure/domain"
	"regexp"
	"strings"
	"time"
)

//todo: 依赖商户的 MSS 发送通知消息,应去掉
var _ member.IMember = new(memberImpl)

type memberImpl struct {
	_manager     member.IMemberManager
	_value       *member.ValueMember
	_account     member.IAccount
	_bank        *member.BankInfo
	_level       *member.Level
	_rep         member.IMemberRep
	_merchantRep merchant.IMerchantRep
	_relation    *member.MemberRelation
	_invitation  member.IInvitationManager
	_mssProvider mss.IMessageProvider
}

func NewMember(manager member.IMemberManager, val *member.ValueMember, rep member.IMemberRep,
	mp mss.IMessageProvider, merchantRep merchant.IMerchantRep) member.IMember {
	return &memberImpl{
		_manager:     manager,
		_value:       val,
		_rep:         rep,
		_mssProvider: mp,
		_merchantRep: merchantRep,
	}
}

// 获取聚合根编号
func (this *memberImpl) GetAggregateRootId() int {
	return this._value.Id
}

// 获取值
func (this *memberImpl) GetValue() member.ValueMember {
	return *this._value
}

var (
	userRegex  = regexp.MustCompile("^[a-zA-Z0-9_]{6,}$")
	emailRegex = regexp.MustCompile("\\w+([-+.']\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*")
	phoneRegex = regexp.MustCompile("^(13[0-9]|15[0|1|2|3|4|5|6|8|9]|18[0|1|2|3|5|6|7|8|9]|17[0|6|7|8]|14[7])(\\d{8})$")
	qqRegex    = regexp.MustCompile("^\\d{5,12}$")
)

func (this *memberImpl) validate(v *member.ValueMember) error {
	v.Usr = strings.ToLower(strings.TrimSpace(v.Usr)) // 小写并删除空格
	v.Name = strings.TrimSpace(v.Name)
	v.Email = strings.ToLower(strings.TrimSpace(v.Email))
	v.Phone = strings.TrimSpace(v.Phone)

	if len([]rune(v.Usr)) < 6 {
		return member.ErrUsrLength
	}
	if !userRegex.MatchString(v.Usr) {
		return member.ErrUsrValidErr
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
	//if len(v.Qq) != 0 && !qqRegex.MatchString(v.Qq) {
	//	return member.ErrQqValidErr
	//}
	return nil
}

// 设置值
func (this *memberImpl) SetValue(v *member.ValueMember) error {
	v.Usr = this._value.Usr
	v.Address = strings.TrimSpace(v.Address)
	v.Im = strings.TrimSpace(v.Im)
	v.Email = strings.TrimSpace(v.Email)
	v.Phone = strings.TrimSpace(v.Phone)
	v.Name = strings.TrimSpace(v.Name)
	v.Ext1 = strings.TrimSpace(v.Ext1)
	v.Ext2 = strings.TrimSpace(v.Ext2)
	v.Ext3 = strings.TrimSpace(v.Ext3)
	v.Ext4 = strings.TrimSpace(v.Ext4)
	v.Ext5 = strings.TrimSpace(v.Ext5)
	v.Ext6 = strings.TrimSpace(v.Ext6)
	if err := this.validate(v); err != nil {
		return err
	}
	this._value.Address = v.Address
	this._value.BirthDay = v.BirthDay
	this._value.Im = v.Im
	this._value.Email = v.Email
	this._value.LastLoginTime = v.LastLoginTime
	this._value.Phone = v.Phone
	this._value.Pwd = v.Pwd
	this._value.Name = v.Name
	this._value.Sex = v.Sex
	this._value.RegFrom = v.RegFrom
	this._value.Remark = v.Remark
	this._value.Ext1 = v.Ext1
	this._value.Ext2 = v.Ext2
	this._value.Ext3 = v.Ext3
	this._value.Ext4 = v.Ext4
	this._value.Ext5 = v.Ext5
	this._value.Ext6 = v.Ext6
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

	if this.ProfileCompleted() {
		// 已完善资料
		this.notifyOnProfileComplete()
	}
	return nil
}

func (this *memberImpl) ProfileCompleted() bool {
	return len(this._value.Name) != 0 && len(this._value.Im) != 0 &&
		len(this._value.Email) != 0 && len(this._value.BirthDay) != 0 &&
		len(this._value.Address) != 0 && len(this._value.Phone) != 0 &&
		len(this._value.Avatar) != 0 && this._value.Sex != 0
}

func (this *memberImpl) notifyOnProfileComplete() {
	rl := this.GetRelation()
	pt, err := this._merchantRep.GetMerchant(rl.RegisterMerchantId)
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

func (this *memberImpl) sendNotifyMail(pt merchant.IMerchant) error {
	tplId := pt.KvManager().GetInt(merchant.KeyMssTplIdOfProfileComplete)
	if tplId > 0 {
		mailTpl := this._mssProvider.GetMailTemplate(tplId)
		if mailTpl != nil {
			//todo:
			v := &mss.Message{

			}
			val := &mss.ValueMailMessage{
				Subject:mailTpl.Subject,
				Body:mailTpl.Body,
			}
			msg := this._mssProvider.CreateMessage(v)
			//todo:?? data
			var data = map[string]string{
				"Name":           this._value.Name,
				"InvitationCode": this._value.InvitationCode,
			}
			return msg.Send(val,data)
		}
	}
	return errors.New("no such email template")
}

// 邀请管理
func (this *memberImpl) Invitation() member.IInvitationManager {
	if this._invitation == nil {
		this._invitation = &invitationManager{
			_member: this,
		}
	}
	return this._invitation
}

// 获取账户
func (this *memberImpl) GetAccount() member.IAccount {
	if this._account == nil {
		v := this._rep.GetAccount(this._value.Id)
		return NewAccount(v, this._rep)
	}
	return this._account
}

// 获取提现银行信息
func (this *memberImpl) GetBank() member.BankInfo {
	if this._bank == nil {
		this._bank = this._rep.GetBankInfo(this._value.Id)
		if this._bank == nil {
			return member.BankInfo{}
		}
	}
	return *this._bank
}

// 保存提现银行信息
func (this *memberImpl) SaveBank(v *member.BankInfo) error {
	this.GetBank()
	if this._bank == nil {
		this._bank = v
	} else {
		if this._bank.IsLocked == member.BankLocked {
			return member.ErrBankInfoLocked
		}
		this._bank.Account = v.Account
		this._bank.AccountName = v.AccountName
		this._bank.Network = v.Network
		this._bank.State = v.State
		this._bank.Name = v.Name
	}
	this._bank.State = member.StateOk       //todo:???
	this._bank.IsLocked = member.BankLocked //锁定
	this._bank.UpdateTime = time.Now().Unix()
	//this._bank.MemberId = this.value.Id
	return this._rep.SaveBankInfo(this._bank)
}

// 解锁提现银行卡信息
func (this *memberImpl) UnlockBank() error {
	this.GetBank()
	if this._bank == nil {
		return member.ErrBankInfoNoYetSet
	}
	this._bank.IsLocked = member.BankNoLock
	return this._rep.SaveBankInfo(this._bank)
}

// 保存积分记录
func (this *memberImpl) SaveIntegralLog(l *member.IntegralLog) error {
	l.MemberId = this._value.Id
	return this._rep.SaveIntegralLog(l)
}

// 增加经验值
func (this *memberImpl) AddExp(exp int) error {
	this._value.Exp += exp
	_, err := this.Save()
	//判断是否升级
	this.checkUpLevel()

	return err
}

// 获取等级
func (this *memberImpl) GetLevel() *member.Level {
	if this._level == nil {
		this._level = this._manager.LevelManager().
			GetLevelById(this._value.Level)
	}
	return this._level
}

//　增加积分
// todo:merchantId 不需要
func (this *memberImpl) AddIntegral(merchantId int, backType int,
	integral int, log string) error {
	inLog := &member.IntegralLog{
		MerchantId: merchantId,
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
func (this *memberImpl) checkUpLevel() bool {
	lg := this._manager.LevelManager()
	levelId := lg.GetLevelIdByExp(this._value.Exp)
	if levelId != 0 && this._value.Level < levelId {
		this._value.Level = levelId
		this.Save()
		this._level = nil
		return true
	}
	return false
}

// 获取会员关联
func (this *memberImpl) GetRelation() *member.MemberRelation {
	if this._relation == nil {
		this._relation = this._rep.GetRelation(this._value.Id)
	}
	return this._relation
}

// 更换用户名
func (this *memberImpl) ChangeUsr(usr string) error {
	if usr == this._value.Usr {
		return member.ErrSameUsr
	}
	if len([]rune(usr)) < 6 {
		return member.ErrUsrLength
	}
	if !userRegex.MatchString(usr) {
		return member.ErrUsrValidErr
	}
	if this.usrIsExist(usr) {
		return member.ErrUsrExist
	}
	this._value.Usr = usr
	_, err := this.Save()
	return err
}

// 保存
func (this *memberImpl) Save() (int, error) {
	this._value.UpdateTime = time.Now().Unix() // 更新时间，数据以更新时间触发
	if this._value.Id > 0 {
		return this._rep.SaveMember(this._value)
	}

	if len(this._value.Name) == 0 {
		this._value.Name = this._value.Usr
	}
	if err := this.validate(this._value); err != nil {
		return this.GetAggregateRootId(), err
	}
	return this.create(this._value)
}

// 锁定会员
func (this *memberImpl) Lock() error {
	return this._rep.LockMember(this.GetAggregateRootId(), 0)
}

// 解锁会员
func (this *memberImpl) Unlock() error {
	return this._rep.LockMember(this.GetAggregateRootId(), 1)
}

// 修改密码,旧密码可为空
func (this *memberImpl) ModifyPassword(newPwd, oldPwd string) error {
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
func (this *memberImpl) ModifyTradePassword(newPwd, oldPwd string) error {
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
func (this *memberImpl) create(m *member.ValueMember) (int, error) {

	//todo: 获取推荐人编号
	//todo: 检测是否有注册权限
	//if err := this._manager.RegisterPerm(this._relation.RefereesId);err != nil{
	//	return -1,err
	//}

	if this.usrIsExist(m.Usr) {
		return -1, member.ErrUsrExist
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
	m.BirthDay = "1970-01-01"
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
func (this *memberImpl) generateInvitationCode() string {
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
func (this *memberImpl) usrIsExist(usr string) bool {
	return this._rep.CheckUsrExist(usr, this.GetAggregateRootId())
}

// 手机号码是否占用
func (this *memberImpl) PhoneIsExist(phone string) bool {
	return this._rep.CheckPhoneBind(this._value.Usr, this.GetAggregateRootId())
}

// 创建并初始化
func (this *memberImpl) SaveRelation(r *member.MemberRelation) error {
	this._relation = r
	this._relation.MemberId = this._value.Id
	return this._rep.SaveRelation(this._relation)
}

// 创建配送地址
func (this *memberImpl) CreateDeliver(v *member.DeliverAddress) (member.IDeliver, error) {
	return newDeliver(v, this._rep)
}

// 获取配送地址
func (this *memberImpl) GetDeliverAddress() []member.IDeliver {
	var vls []*member.DeliverAddress
	vls = this._rep.GetDeliverAddress(this.GetAggregateRootId())
	var arr []member.IDeliver = make([]member.IDeliver, len(vls))
	for i, v := range vls {
		arr[i], _ = this.CreateDeliver(v)
	}
	return arr
}

// 获取配送地址
func (this *memberImpl) GetDeliver(deliverId int) member.IDeliver {
	v := this._rep.GetSingleDeliverAddress(this.GetAggregateRootId(), deliverId)
	if v != nil {
		d, _ := this.CreateDeliver(v)
		return d
	}
	return nil
}

// 删除配送地址
func (this *memberImpl) DeleteDeliver(deliverId int) error {
	return this._rep.DeleteDeliver(this.GetAggregateRootId(), deliverId)
}
