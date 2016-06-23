/**
 * Copyright 2015 @ z3q.net.
 * name : member_profile.go
 * author : jarryliu
 * date : 2016-06-23 16:31
 * description :
 * history :
 */
package member

import (
	"errors"
	"fmt"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/tmp"
	"go2o/core/infrastructure/domain"
	"regexp"
	"strings"
	"time"
)

var _ member.IProfileManager = new(profileManagerImpl)
var (
	qqRegex = regexp.MustCompile("^\\d{5,12}$")
)

type profileManagerImpl struct {
	_member      *memberImpl
	_memberId    int
	_rep         member.IMemberRep
	_bank        *member.BankInfo
	_trustedInfo *member.TrustedInfo
	_profile     *member.Profile
}

func newProfileManagerImpl(m *memberImpl, memberId int,
	rep member.IMemberRep) member.IProfileManager {
	if memberId == 0 {
		//如果会员不存在,则不应创建服务
		panic(errors.New("member not exists"))
	}
	return &profileManagerImpl{
		_member:   m,
		_memberId: memberId,
		_rep:      rep,
	}
}

// 手机号码是否占用
func (this *profileManagerImpl) phoneIsExist(phone string) bool {
	return this._rep.CheckPhoneBind(phone, this._memberId)
}

// 验证数据
func (this *profileManagerImpl) validateProfile(v *member.Profile) error {
	v.Name = strings.TrimSpace(v.Name)
	v.Email = strings.ToLower(strings.TrimSpace(v.Email))
	v.Phone = strings.TrimSpace(v.Phone)

	if len([]rune(v.Name)) < 2 {
		return member.ErrPersonName
	}

	if len(v.Email) != 0 && !emailRegex.MatchString(v.Email) {
		return member.ErrEmailValidErr
	}
	if len(v.Phone) != 0 && !phoneRegex.MatchString(v.Phone) {
		return member.ErrPhoneValidErr
	}

	if len(v.Phone) > 0 && this.phoneIsExist(v.Phone) {
		return member.ErrPhoneHasBind
	}

	//if len(v.Qq) != 0 && !qqRegex.MatchString(v.Qq) {
	//	return member.ErrQqValidErr
	//}
	return nil
}

// 拷贝资料
func (this *profileManagerImpl) copyProfile(v, dst *member.Profile) error {
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
	if err := this.validateProfile(v); err != nil {
		return err
	}

	//pro.Avatar = "resource/no_avatar.gif"
	//pro.BirthDay = "1970-01-01"
	//
	//// 如果昵称为空，则跟用户名相同
	//if len(pro.Name) == 0 {
	//    pro.Name = m.Usr
	//}

	dst.Address = v.Address
	dst.BirthDay = v.BirthDay
	dst.Im = v.Im
	dst.Email = v.Email
	dst.Phone = v.Phone
	dst.Name = v.Name
	dst.Sex = v.Sex
	dst.Remark = v.Remark
	dst.Ext1 = v.Ext1
	dst.Ext2 = v.Ext2
	dst.Ext3 = v.Ext3
	dst.Ext4 = v.Ext4
	dst.Ext5 = v.Ext5
	dst.Ext6 = v.Ext6
	if v.Avatar != "" {
		dst.Avatar = v.Avatar
	}
	return nil
}

func (this *profileManagerImpl) ProfileCompleted() bool {
	v := this.GetProfile()
	return len(v.Name) != 0 && len(v.Im) != 0 &&
		len(v.Email) != 0 && len(v.BirthDay) != 0 &&
		len(v.Address) != 0 && len(v.Phone) != 0 &&
		len(v.Avatar) != 0 && v.Sex != 0
}

// 获取资料
func (this *profileManagerImpl) GetProfile() member.Profile {
	if this._profile == nil {
		this._profile = this._rep.GetProfile(this._memberId)
	}
	return *this._profile
}

// 保存资料
func (this *profileManagerImpl) SaveProfile(v *member.Profile) error {
	ptr := this.GetProfile()
	err := this.copyProfile(v, &ptr)
	if err == nil {
		ptr.MemberId = this._memberId
		err = this._rep.SaveProfile(&ptr)
		if this.ProfileCompleted() {
			// 已完善资料
			this.notifyOnProfileComplete()
		}
	}
	return err
}

//todo: ?? 重构
func (this *profileManagerImpl) notifyOnProfileComplete() {
	rl := this._member.GetRelation()
	pt, err := this._member._merchantRep.GetMerchant(rl.RegisterMerchantId)
	if err == nil {
		key := fmt.Sprintf("profile:complete:id_%d", this._memberId)
		if pt.MemberKvManager().GetInt(key) == 0 {
			if err := this.sendNotifyMail(pt); err == nil {
				pt.MemberKvManager().Set(key, "1")
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func (this *profileManagerImpl) sendNotifyMail(pt merchant.IMerchant) error {
	tplId := pt.KvManager().GetInt(merchant.KeyMssTplIdOfProfileComplete)
	if tplId > 0 {
		mailTpl := this._member._mssRep.GetProvider().GetMailTemplate(tplId)
		if mailTpl != nil {
			v := &mss.Message{
				// 消息类型
				Type: mss.TypeEmailMessage,
				// 消息用途
				UseFor: mss.UseForNotify,
				// 发送人角色
				SenderRole: mss.RoleSystem,
				// 发送人编号
				SenderId: 0,
				// 发送的目标
				To: []mss.User{
					mss.User{
						Role: mss.RoleMember,
						Id:   this._memberId,
					},
				},
				// 发送的用户角色
				ToRole: -1,
				// 全系统接收
				AllUser: -1,
				// 是否只能阅读
				Readonly: 1,
			}
			val := &mss.MailMessage{
				Subject: mailTpl.Subject,
				Body:    mailTpl.Body,
			}
			msg := this._member._mssRep.GetManager().CreateMessage(v, val)
			//todo:?? data
			var data = map[string]string{
				"Name":           this._profile.Name,
				"InvitationCode": this._member.GetValue().InvitationCode,
			}
			return msg.Send(data)
		}
	}
	return errors.New("no such email template")
}

//todo: 密码应独立为credential

// 修改密码,旧密码可为空
func (this *profileManagerImpl) ModifyPassword(newPwd, oldPwd string) error {
	var err error
	if newPwd == oldPwd {
		return member.ErrPwdCannotSame
	}

	if b, err := domain.ChkPwdRight(newPwd); !b {
		return err
	}

	if len(oldPwd) != 0 && oldPwd != this._member._value.Pwd {
		return member.ErrPwdOldPwdNotRight
	}

	this._member._value.Pwd = newPwd
	_, err = this._member.Save()

	return err
}

// 修改交易密码，旧密码可为空
func (this *profileManagerImpl) ModifyTradePassword(newPwd, oldPwd string) error {
	var err error
	if newPwd == oldPwd {
		return member.ErrPwdCannotSame
	}
	if b, err := domain.ChkPwdRight(newPwd); !b {
		return err
	}
	// 已经设置过旧密码
	if len(this._member._value.TradePwd) != 0 && this._member._value.TradePwd != oldPwd {
		return member.ErrPwdOldPwdNotRight
	}
	this._member._value.TradePwd = newPwd
	_, err = this._member.Save()
	return err
}

// 获取提现银行信息
func (this *profileManagerImpl) GetBank() member.BankInfo {
	if this._bank == nil {
		this._bank = this._rep.GetBankInfo(this._memberId)
		if this._bank == nil {
			return member.BankInfo{}
		}
	}
	return *this._bank
}

// 保存提现银行信息
func (this *profileManagerImpl) SaveBank(v *member.BankInfo) error {
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
func (this *profileManagerImpl) UnlockBank() error {
	this.GetBank()
	if this._bank == nil {
		return member.ErrBankInfoNoYetSet
	}
	this._bank.IsLocked = member.BankNoLock
	return this._rep.SaveBankInfo(this._bank)
}

// 创建配送地址
func (this *profileManagerImpl) CreateDeliver(v *member.DeliverAddress) (member.IDeliver, error) {
	return newDeliver(v, this._rep)
}

// 获取配送地址
func (this *profileManagerImpl) GetDeliverAddress() []member.IDeliver {
	var vls []*member.DeliverAddress
	vls = this._rep.GetDeliverAddress(this._memberId)
	var arr []member.IDeliver = make([]member.IDeliver, len(vls))
	for i, v := range vls {
		arr[i], _ = this.CreateDeliver(v)
	}
	return arr
}

// 获取配送地址
func (this *profileManagerImpl) GetDeliver(deliverId int) member.IDeliver {
	v := this._rep.GetSingleDeliverAddress(this._memberId, deliverId)
	if v != nil {
		d, _ := this.CreateDeliver(v)
		return d
	}
	return nil
}

// 删除配送地址
func (this *profileManagerImpl) DeleteDeliver(deliverId int) error {
	//todo: 至少保留一个配送地址
	return this._rep.DeleteDeliver(this._memberId, deliverId)
}

// 拷贝认证信息
func (this *profileManagerImpl) copyTrustedInfo(src, dst *member.TrustedInfo) {
	dst.RealName = src.RealName
	dst.BodyNumber = src.BodyNumber
	dst.TrustImage = src.TrustImage
}

// 实名认证信息
func (this *profileManagerImpl) GetTrustedInfo() member.TrustedInfo {
	if this._trustedInfo == nil {
		//如果还没有实名信息,则新建
		orm := tmp.Db().GetOrm()
		if err := orm.Get(this._memberId, &this._trustedInfo); err != nil {
			this._trustedInfo = &member.TrustedInfo{
				MemberId: this._memberId,
			}
			orm.Save(nil, this._trustedInfo)
		}
	}
	return *this._trustedInfo
}

// 保存实名认证信息
func (this *profileManagerImpl) SaveTrustedInfo(v *member.TrustedInfo) error {
	this.GetTrustedInfo()
	v.TrustImage = strings.TrimSpace(v.TrustImage)
	v.BodyNumber = strings.TrimSpace(v.BodyNumber)
	v.RealName = strings.TrimSpace(v.RealName)
	if len(v.TrustImage) == 0 || len(v.RealName) == 0 ||
		len(v.BodyNumber) == 0 {
		return member.ErrMissingTrustedInfo
	}
	this.copyTrustedInfo(v, this._trustedInfo)
	this._trustedInfo.IsHandle = 0 //标记为未处理
	this._trustedInfo.UpdateTime = time.Now().Unix()
	_, _, err := tmp.Db().GetOrm().Save(nil, this._trustedInfo)
	return err
}

// 审核实名认证,若重复审核将返回错误
func (this *profileManagerImpl) ReviewTrustedInfo(pass bool, remark string) error {
	this.GetTrustedInfo()
	if pass {
		this._trustedInfo.Reviewed = 1
	} else {
		this._trustedInfo.Reviewed = 0
	}
	this._trustedInfo.IsHandle = 1 //标记为已处理
	this._trustedInfo.Remark = remark
	this._trustedInfo.ReviewTime = time.Now().Unix()
	_, _, err := tmp.Db().GetOrm().Save(nil, this._trustedInfo)
	return err
}
