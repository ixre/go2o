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
	"go2o/core/domain"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/tmp"
	dm "go2o/core/infrastructure/domain"
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
	_valRep      valueobject.IValueRep
	_bank        *member.BankInfo
	_trustedInfo *member.TrustedInfo
	_profile     *member.Profile
}

func newProfileManagerImpl(m *memberImpl, memberId int,
	rep member.IMemberRep, valRep valueobject.IValueRep) member.IProfileManager {
	if memberId == 0 {
		//如果会员不存在,则不应创建服务
		panic(errors.New("member not exists"))
	}
	return &profileManagerImpl{
		_member:   m,
		_memberId: memberId,
		_rep:      rep,
		_valRep:   valRep,
	}
}

// 手机号码是否占用
func (pm *profileManagerImpl) phoneIsExist(phone string) bool {
	return pm._rep.CheckPhoneBind(phone, pm._memberId)
}

// 验证数据
func (pm *profileManagerImpl) validateProfile(v *member.Profile) error {
	v.Name = strings.TrimSpace(v.Name)
	v.Email = strings.ToLower(strings.TrimSpace(v.Email))
	v.Phone = strings.TrimSpace(v.Phone)

	if len([]rune(v.Name)) < 1 {
		return member.ErrNilNickName
	}

	if len(v.Email) != 0 && !emailRegex.MatchString(v.Email) {
		return member.ErrEmailValidErr
	}
	if len(v.Phone) != 0 && !phoneRegex.MatchString(v.Phone) {
		return member.ErrPhoneValidErr
	}

	if len(v.Phone) > 0 && pm.phoneIsExist(v.Phone) {
		return member.ErrPhoneHasBind
	}

	//if len(v.Qq) != 0 && !qqRegex.MatchString(v.Qq) {
	//	return member.ErrQqValidErr
	//}
	return nil
}

// 拷贝资料
func (pm *profileManagerImpl) copyProfile(v, dst *member.Profile) error {
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
	if err := pm.validateProfile(v); err != nil {
		return err
	}

	//pro.Avatar = "res/no_avatar.gif"
	//pro.BirthDay = "1970-01-01"
	//
	//// 如果昵称为空，则跟用户名相同
	//if len(pro.Name) == 0 {
	//    pro.Name = m.Usr
	//}
	dst.Province = v.Province
	dst.City = v.City
	dst.District = v.District
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

func (pm *profileManagerImpl) ProfileCompleted() bool {
	v := pm.GetProfile()
	return len(v.Name) != 0 && len(v.Im) != 0 &&
		len(v.Email) != 0 && len(v.BirthDay) != 0 &&
		len(v.Address) != 0 && len(v.Phone) != 0 &&
		len(v.Avatar) != 0 && v.Sex != 0
}

// 获取资料
func (pm *profileManagerImpl) GetProfile() member.Profile {
	if pm._profile == nil {
		pm._profile = pm._rep.GetProfile(pm._memberId)
	}
	return *pm._profile
}

// 保存资料
func (pm *profileManagerImpl) SaveProfile(v *member.Profile) error {
	ptr := pm.GetProfile()
	err := pm.copyProfile(v, &ptr)
	if err == nil {
		ptr.MemberId = pm._memberId
		err = pm._rep.SaveProfile(&ptr)
		if pm.ProfileCompleted() {
			// 已完善资料
			pm.notifyOnProfileComplete()
		}
	}
	return err
}

//todo: ?? 重构
func (pm *profileManagerImpl) notifyOnProfileComplete() {
	rl := pm._member.GetRelation()
	pt, err := pm._member._merchantRep.GetMerchant(rl.RegisterMerchantId)
	if err == nil {
		key := fmt.Sprintf("profile:complete:id_%d", pm._memberId)
		if pt.MemberKvManager().GetInt(key) == 0 {
			if err := pm.sendNotifyMail(pt); err == nil {
				pt.MemberKvManager().Set(key, "1")
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func (pm *profileManagerImpl) sendNotifyMail(pt merchant.IMerchant) error {
	tplId := pt.KvManager().GetInt(merchant.KeyMssTplIdOfProfileComplete)
	if tplId > 0 {
		mailTpl := pm._member._mssRep.GetProvider().GetMailTemplate(tplId)
		if mailTpl != nil {
			v := &mss.Message{
				// 消息类型
				Type: notify.TypeEmailMessage,
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
						Id:   pm._memberId,
					},
				},
				// 发送的用户角色
				ToRole: -1,
				// 全系统接收
				AllUser: -1,
				// 是否只能阅读
				Readonly: 1,
			}
			val := &notify.MailMessage{
				Subject: mailTpl.Subject,
				Body:    mailTpl.Body,
			}
			msg := pm._member._mssRep.MessageManager().CreateMessage(v, val)
			//todo:?? data
			var data = map[string]string{
				"Name":           pm._profile.Name,
				"InvitationCode": pm._member.GetValue().InvitationCode,
			}
			return msg.Send(data)
		}
	}
	return errors.New("no such email template")
}

//todo: 密码应独立为credential

// 修改密码,旧密码可为空
func (pm *profileManagerImpl) ModifyPassword(newPwd, oldPwd string) error {
	var err error
	if newPwd == oldPwd {
		return domain.ErrPwdCannotSame
	}

	if b, err := dm.ChkPwdRight(newPwd); !b {
		return err
	}

	if len(oldPwd) != 0 &&
		dm.MemberSha1Pwd(oldPwd) != pm._member._value.Pwd {
		return domain.ErrPwdOldPwdNotRight
	}

	pm._member._value.Pwd = dm.MemberSha1Pwd(newPwd)
	_, err = pm._member.Save()

	return err
}

// 修改交易密码，旧密码可为空
func (pm *profileManagerImpl) ModifyTradePassword(newPwd, oldPwd string) error {
	var err error
	if newPwd == oldPwd {
		return domain.ErrPwdCannotSame
	}
	if b, err := dm.ChkPwdRight(newPwd); !b {
		return err
	}
	// 已经设置过旧密码
	if len(pm._member._value.TradePwd) != 0 &&
		pm._member._value.TradePwd != dm.MemberSha1Pwd(oldPwd) {
		return domain.ErrPwdOldPwdNotRight
	}
	pm._member._value.TradePwd = dm.MemberSha1Pwd(newPwd)
	_, err = pm._member.Save()
	return err
}

// 获取提现银行信息
func (pm *profileManagerImpl) GetBank() member.BankInfo {
	if pm._bank == nil {
		pm._bank = pm._rep.GetBankInfo(pm._memberId)
		if pm._bank == nil {
			return member.BankInfo{}
		}
	}
	return *pm._bank
}

// 保存提现银行信息
func (pm *profileManagerImpl) SaveBank(v *member.BankInfo) error {
	pm.GetBank()
	if pm._bank == nil {
		pm._bank = v
	} else {
		if pm._bank.IsLocked == member.BankLocked {
			return member.ErrBankInfoLocked
		}
		pm._bank.Account = v.Account
		pm._bank.AccountName = v.AccountName
		pm._bank.Network = v.Network
		pm._bank.State = v.State
		pm._bank.Name = v.Name
	}
	pm._bank.State = member.StateOk       //todo:???
	pm._bank.IsLocked = member.BankLocked //锁定
	pm._bank.UpdateTime = time.Now().Unix()
	//pm._bank.MemberId = pm.value.Id
	return pm._rep.SaveBankInfo(pm._bank)
}

// 解锁提现银行卡信息
func (pm *profileManagerImpl) UnlockBank() error {
	pm.GetBank()
	if pm._bank == nil {
		return member.ErrBankInfoNoYetSet
	}
	pm._bank.IsLocked = member.BankNoLock
	return pm._rep.SaveBankInfo(pm._bank)
}

// 创建配送地址
func (pm *profileManagerImpl) CreateDeliver(v *member.DeliverAddress) member.IDeliverAddress {
	return newDeliver(v, pm._rep, pm._valRep)
}

// 获取配送地址
func (pm *profileManagerImpl) GetDeliverAddress() []member.IDeliverAddress {
	var vls []*member.DeliverAddress
	vls = pm._rep.GetDeliverAddress(pm._memberId)
	var arr []member.IDeliverAddress = make([]member.IDeliverAddress, len(vls))
	for i, v := range vls {
		arr[i] = pm.CreateDeliver(v)
	}
	return arr
}

// 获取配送地址
func (pm *profileManagerImpl) GetDeliver(deliverId int) member.IDeliverAddress {
	v := pm._rep.GetSingleDeliverAddress(pm._memberId, deliverId)
	if v != nil {
		return pm.CreateDeliver(v)
	}
	return nil
}

// 删除配送地址
func (pm *profileManagerImpl) DeleteDeliver(deliverId int) error {
	//todo: 至少保留一个配送地址
	return pm._rep.DeleteDeliver(pm._memberId, deliverId)
}

// 拷贝认证信息
func (pm *profileManagerImpl) copyTrustedInfo(src, dst *member.TrustedInfo) {
	dst.RealName = src.RealName
	dst.BodyNumber = src.BodyNumber
	dst.TrustImage = src.TrustImage
}

// 实名认证信息
func (pm *profileManagerImpl) GetTrustedInfo() member.TrustedInfo {
	if pm._trustedInfo == nil {
		pm._trustedInfo = &member.TrustedInfo{
			MemberId: pm._memberId,
		}
		//如果还没有实名信息,则新建
		orm := tmp.Db().GetOrm()
		if err := orm.Get(pm._memberId, pm._trustedInfo); err != nil {
			orm.Save(nil, pm._trustedInfo)
		}
	}
	return *pm._trustedInfo
}

// 保存实名认证信息
func (pm *profileManagerImpl) SaveTrustedInfo(v *member.TrustedInfo) error {
	pm.GetTrustedInfo()
	v.TrustImage = strings.TrimSpace(v.TrustImage)
	v.BodyNumber = strings.TrimSpace(v.BodyNumber)
	v.RealName = strings.TrimSpace(v.RealName)
	if len(v.TrustImage) == 0 || len(v.RealName) == 0 ||
		len(v.BodyNumber) == 0 {
		return member.ErrMissingTrustedInfo
	}
	pm.copyTrustedInfo(v, pm._trustedInfo)
	pm._trustedInfo.IsHandle = 0 //标记为未处理
	pm._trustedInfo.UpdateTime = time.Now().Unix()
	_, _, err := tmp.Db().GetOrm().Save(nil, pm._trustedInfo)
	return err
}

// 审核实名认证,若重复审核将返回错误
func (pm *profileManagerImpl) ReviewTrustedInfo(pass bool, remark string) error {
	pm.GetTrustedInfo()
	if pass {
		pm._trustedInfo.Reviewed = 1
	} else {
		pm._trustedInfo.Reviewed = 0
	}
	pm._trustedInfo.IsHandle = 1 //标记为已处理
	pm._trustedInfo.Remark = remark
	pm._trustedInfo.ReviewTime = time.Now().Unix()
	_, _, err := tmp.Db().GetOrm().Save(nil, pm._trustedInfo)
	return err
}

var _ member.IDeliverAddress = new(deliverAddressImpl)

type deliverAddressImpl struct {
	_value     *member.DeliverAddress
	_memberRep member.IMemberRep
	_valRep    valueobject.IValueRep
}

func newDeliver(v *member.DeliverAddress, memberRep member.IMemberRep,
	valRep valueobject.IValueRep) member.IDeliverAddress {
	d := &deliverAddressImpl{
		_value:     v,
		_memberRep: memberRep,
		_valRep:    valRep,
	}
	return d
}

func (pm *deliverAddressImpl) GetDomainId() int {
	return pm._value.Id
}

func (pm *deliverAddressImpl) GetValue() member.DeliverAddress {
	return *pm._value
}

func (pm *deliverAddressImpl) SetValue(v *member.DeliverAddress) error {
	if pm._value.MemberId == v.MemberId {
		if err := pm.checkValue(v); err != nil {
			return err
		}
		pm._value = v
	}
	return nil
}

// 设置地区中文名
func (pm *deliverAddressImpl) renewAreaName(v *member.DeliverAddress) string {
	names := pm._valRep.GetAreaNames([]int{
		v.Province,
		v.City,
		v.District,
	})
	return strings.Join(names, " ")
}

func (pm *deliverAddressImpl) checkValue(v *member.DeliverAddress) error {
	v.Address = strings.TrimSpace(v.Address)
	v.RealName = strings.TrimSpace(v.RealName)
	v.Phone = strings.TrimSpace(v.Phone)

	if len([]rune(v.RealName)) < 2 {
		return member.ErrDeliverContactPersonName
	}

	if v.Province <= 0 || v.City <= 0 || v.District <= 0 {
		return member.ErrNotSetArea
	}

	if !phoneRegex.MatchString(v.Phone) {
		return member.ErrDeliverContactPhone
	}

	if len([]rune(v.Address)) < 6 {
		// 判断字符长度
		return member.ErrDeliverAddressLen
	}

	return nil
}

func (pm *deliverAddressImpl) Save() (int, error) {
	if err := pm.checkValue(pm._value); err != nil {
		return pm.GetDomainId(), err
	}
	pm._value.Area = pm.renewAreaName(pm._value)
	return pm._memberRep.SaveDeliver(pm._value)
}
