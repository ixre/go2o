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
	"go2o/core/domain/interface/member"
	"go2o/core/domain/tmp"
	"go2o/core/infrastructure/domain"
	"strings"
	"time"
)

var _ member.IProfileManager = new(profileManagerImpl)

type profileManagerImpl struct {
	_member      *memberImpl
	_memberId    int
	_rep         member.IMemberRep
	_bank        *member.BankInfo
	_trustedInfo *member.TrustedInfo
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
