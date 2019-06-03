/**
 * Copyright 2015 @ to2.net.
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
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/tmp"
	dm "go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/domain/util"
	"go2o/core/msq"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var _ member.IProfileManager = new(profileManagerImpl)
var (
	exampleTrustImageUrl = "res/tru-example.jpg"
	// qqRegex = regexp.MustCompile("^\\d{5,12}$")
	zhNameRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]{2,6}$")
)

type profileManagerImpl struct {
	member       *memberImpl
	memberId     int64
	rep          member.IMemberRepo
	valueRepo    valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
	bank         *member.BankInfo
	trustedInfo  *member.TrustedInfo
	profile      *member.Profile
}

func newProfileManagerImpl(m *memberImpl, memberId int64,
	rep member.IMemberRepo, registryRepo registry.IRegistryRepo,
	valueRepo valueobject.IValueRepo) member.IProfileManager {
	if memberId == 0 {
		//如果会员不存在,则不应创建服务
		panic(errors.New("member not exists"))
	}
	return &profileManagerImpl{
		member:       m,
		memberId:     memberId,
		rep:          rep,
		registryRepo: registryRepo,
		valueRepo:    valueRepo,
	}
}

// 手机号码是否占用
func (p *profileManagerImpl) phoneIsExist(phone string) bool {
	return p.rep.CheckPhoneBind(phone, p.memberId)
}

// 验证数据,用v.updateTime > 0 判断是否为新创建用户
func (p *profileManagerImpl) validateProfile(v *member.Profile) error {
	v.Name = strings.TrimSpace(v.Name)
	v.Email = strings.ToLower(strings.TrimSpace(v.Email))
	v.Phone = strings.TrimSpace(v.Phone)
	// 验证昵称
	if len([]rune(v.Name)) < 1 && v.UpdateTime > 0 {
		return member.ErrNilNickName
	}
	// 检查区域
	if (v.Province == 0 || v.City == 0 || v.District == 0 ||
		len(v.Address) == 0) && v.UpdateTime > 0 {
		return member.ErrAddress
	}
	// 检查邮箱
	if len(v.Email) != 0 && !emailRegex.MatchString(v.Email) {
		return member.ErrEmailValidErr
	}
	// 检查手机
	checkPhone := p.registryRepo.Get(registry.MemberCheckPhoneFormat).BoolValue()
	if len(v.Phone) != 0 && checkPhone {
		if !phoneRegex.MatchString(v.Phone) {
			return member.ErrPhoneValidErr
		}
	}
	if len(v.Phone) > 0 && p.phoneIsExist(v.Phone) {
		return member.ErrPhoneHasBind
	}
	// 检查IM
	if v.UpdateTime > 0 {
		imRequire := p.registryRepo.Get(registry.MemberImRequired).BoolValue()
		if imRequire && len(v.Im) == 0 {
			return member.ErrMissingIM
		}
	}
	return nil
}

// 拷贝资料
func (p *profileManagerImpl) copyProfile(v, dst *member.Profile) error {
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
	if err := p.validateProfile(v); err != nil {
		return err
	}
	dst.Province = v.Province
	dst.City = v.City
	dst.District = v.District
	dst.Address = v.Address
	dst.BirthDay = v.BirthDay
	dst.Im = v.Im
	dst.Email = v.Email

	//todo: 如果手机不需要验证，则可以随意设置
	if dst.Phone == "" {
		dst.Phone = v.Phone
	}
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

func (p *profileManagerImpl) ProfileCompleted() bool {
	v := p.GetProfile()
	r := len(v.Name) != 0 &&
		len(v.BirthDay) != 0 && len(v.Address) != 0 && v.Sex != 0 &&
		v.Province != 0 && v.City != 0 && v.District != 0
	if r {
		imRequire := p.registryRepo.Get(registry.MemberImRequired).BoolValue()
		if imRequire && len(v.Im) == 0 {
			return false
		}
	}
	return r
}

func (p *profileManagerImpl) CheckProfileComplete() error {
	v := p.GetProfile()
	if v.Phone == "" {
		return errors.New("phone")
	}
	if v.BirthDay == "" {
		return errors.New("birthday")
	}
	if v.Province <= 0 || v.City <= 0 || v.District <= 0 || v.Address == "" {
		return errors.New("address")
	}
	imRequire := p.registryRepo.Get(registry.MemberImRequired).BoolValue()
	if imRequire && len(v.Im) == 0 {
		return errors.New("im")
	}
	return nil
}

//todo: 上传头像方法

// 获取资料
func (p *profileManagerImpl) GetProfile() member.Profile {
	if p.profile == nil {
		p.profile = p.rep.GetProfile(p.memberId)
	}
	return *p.profile
}

// 保存资料
func (p *profileManagerImpl) SaveProfile(v *member.Profile) error {
	ptr := p.GetProfile()
	err := p.copyProfile(v, &ptr)
	if err == nil {
		ptr.MemberId = p.memberId
		err = p.rep.SaveProfile(&ptr)
		if err == nil {
			// 推送资料更新消息
			go msq.PushDelay(msq.MemberProfileUpdated, strconv.Itoa(int(p.memberId)), "", 500)
			// 完善资料通知
			if p.ProfileCompleted() {
				// 标记会员已完善资料
				if !p.member.ContainFlag(member.FlagProfileCompleted) {
					p.member.value.Flag |= member.FlagProfileCompleted
					if err == nil {
						p.member.Save()
					}
				}
				p.notifyOnProfileComplete()
			}
		}

	}
	return err
}

// 更改手机号码
func (p *profileManagerImpl) ChangePhone(phone string) error {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return member.ErrPhoneValidErr
	}
	used := p.rep.CheckPhoneBind(phone, p.memberId)
	if !used {
		v := p.GetProfile()
		v.Phone = phone
		return p.rep.SaveProfile(&v)
	}
	return member.ErrPhoneHasBind
}

// 设置头像
func (p *profileManagerImpl) ChangeAvatar(avatar string) error {
	if avatar == "" {
		return member.ErrNullAvatar
	}
	v := p.GetProfile()
	if p.profile != nil {
		p.profile.Avatar = avatar
	}
	v.Avatar = avatar
	return p.rep.SaveProfile(&v)
}

//todo: ?? 重构
func (p *profileManagerImpl) notifyOnProfileComplete() {
	//rl := p._member.GetRelation()
	//pt, err := p._member._merchantRepo.GetMerchant(rl.RegisterMchId)
	//if err == nil {
	//	key := fmt.Sprintf("profile:complete:id_%d", p._memberId)
	//	if pt.MemberKvManager().GetInt(key) == 0 {
	//		if err := p.sendNotifyMail(pt); err == nil {
	//			pt.MemberKvManager().Set(key, "1")
	//		} else {
	//			fmt.Println(err.Error())
	//		}
	//	}
	//}
}

func (p *profileManagerImpl) sendNotifyMail(pt merchant.IMerchant) error {
	tplId := pt.KvManager().GetInt(merchant.KeyMssTplIdOfProfileComplete)
	if tplId > 0 {
		mailTpl := p.member.mssRepo.GetProvider().GetMailTemplate(int32(tplId))
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
					{
						Role: mss.RoleMember,
						Id:   int32(p.memberId),
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
			msg := p.member.mssRepo.MessageManager().CreateMessage(v, val)
			//todo:?? data
			var data = map[string]string{
				"Name":           p.profile.Name,
				"InvitationCode": p.member.GetValue().InvitationCode,
			}
			return msg.Send(data)
		}
	}
	return errors.New("no such email template")
}

//todo: 密码应独立为credential

// 修改密码,旧密码可为空
func (p *profileManagerImpl) ModifyPassword(newPwd, oldPwd string) error {
	if b, err := dm.ChkPwdRight(newPwd); !b {
		return err
	}
	//log.Println("----",p.member.value.Pwd)
	//log.Println("----",oldPwd)
	//log.Println("----",newPwd)
	//log.Println("---- 123000 / ",dm.MemberSha1Pwd("123000"))
	//log.Println("---- 123456 / ",dm.MemberSha1Pwd("123456"))
	if len(oldPwd) != 0 {
		if newPwd == oldPwd {
			return domain.ErrPwdCannotSame
		}
		if oldPwd != p.member.value.Pwd {
			return domain.ErrPwdOldPwdNotRight
		}
	}
	p.member.value.Pwd = newPwd
	_, err := p.member.Save()
	return err
}

// 修改交易密码，旧密码可为空
func (p *profileManagerImpl) ModifyTradePassword(newPwd, oldPwd string) error {
	if newPwd == oldPwd {
		return domain.ErrPwdCannotSame
	}
	if b, err := dm.ChkPwdRight(newPwd); !b {
		return err
	}
	// 已经设置过旧密码
	if len(oldPwd) != 0 && p.member.value.TradePwd != oldPwd {
		return domain.ErrPwdOldPwdNotRight
	}
	p.member.value.TradePwd = newPwd
	_, err := p.member.Save()
	return err
}

// 获取提现银行信息
func (p *profileManagerImpl) GetBank() member.BankInfo {
	if p.bank == nil {
		p.bank = p.rep.GetBankInfo(p.memberId)
		if p.bank == nil {
			p.bank = &member.BankInfo{
				MemberId:   p.memberId,
				IsLocked:   member.BankNoLock,
				State:      0,
				UpdateTime: time.Now().Unix(),
			}
			orm.Save(tmp.Db().GetOrm(), p.bank, 0)
		}
	}
	return *p.bank
}

// 保存提现银行信息
func (p *profileManagerImpl) SaveBank(v *member.BankInfo) error {
	v.Account = strings.TrimSpace(v.Account)
	v.AccountName = strings.TrimSpace(v.AccountName)
	v.Network = strings.TrimSpace(v.Network)
	v.BankName = strings.TrimSpace(v.BankName)
	if v.Account == "" || v.BankName == "" {
		return member.ErrBankInfo
	}
	trustInfo := p.GetTrustedInfo()
	if trustInfo.ReviewState == 0 {
		return member.ErrNotTrusted
	}
	p.GetBank()
	if p.bank.IsLocked == member.BankLocked {
		return member.ErrBankInfoLocked
	}
	v.AccountName = trustInfo.RealName
	err := p.checkBank(v)
	if err == nil {
		p.bank.Account = v.Account
		p.bank.AccountName = v.AccountName
		p.bank.Network = v.Network
		p.bank.BankName = v.BankName
		p.bank.State = member.StateOk       //todo:???
		p.bank.IsLocked = member.BankLocked //锁定
		p.bank.UpdateTime = time.Now().Unix()
		p.bank.MemberId = p.memberId
		err = p.rep.SaveBankInfo(p.bank)
	}
	return err
}

// 检查银行信息
func (p *profileManagerImpl) checkBank(v *member.BankInfo) error {
	v.Account = strings.TrimSpace(v.Account)
	v.AccountName = strings.TrimSpace(v.AccountName)
	v.Network = strings.TrimSpace(v.Network)
	v.BankName = strings.TrimSpace(v.BankName)

	if v.BankName == "" {
		return member.ErrBankName
	}
	if v.AccountName == "" {
		return member.ErrBankAccountName
	}
	if v.Account == "" || len(v.Account) < 16 {
		return member.ErrBankAccount
	}
	if v.Network == "" {
		//return member.ErrBankNetwork
	}
	return nil
}

// 解锁提现银行卡信息
func (p *profileManagerImpl) UnlockBank() error {
	p.GetBank()
	if p.bank == nil {
		return member.ErrBankInfoNoYetSet
	}
	p.bank.IsLocked = member.BankNoLock
	return p.rep.SaveBankInfo(p.bank)
}

// 创建配送地址
func (p *profileManagerImpl) CreateDeliver(v *member.Address) member.IDeliverAddress {
	return newDeliver(v, p.rep, p.valueRepo)
}

// 获取配送地址
func (p *profileManagerImpl) GetDeliverAddress() []member.IDeliverAddress {
	list := p.rep.GetDeliverAddress(p.memberId)
	var arr []member.IDeliverAddress = make([]member.IDeliverAddress, len(list))
	for i, v := range list {
		arr[i] = p.CreateDeliver(v)
	}
	return arr
}

// 设置默认地址
func (p *profileManagerImpl) SetDefaultAddress(addressId int64) error {
	for _, v := range p.GetDeliverAddress() {
		vv := v.GetValue()
		if v.GetDomainId() == addressId {
			vv.IsDefault = 1
		} else {
			vv.IsDefault = 0
		}
		p.rep.SaveDeliver(&vv)
	}
	return nil
}

// 获取默认收货地址
func (p *profileManagerImpl) GetDefaultAddress() member.IDeliverAddress {
	list := p.rep.GetDeliverAddress(p.memberId)
	// 查找是否有默认地址
	for _, v := range list {
		if v.IsDefault == 1 {
			return p.CreateDeliver(v)
		}
	}
	// 使用第一个地址
	if len(list) > 0 {
		return p.CreateDeliver(list[0])
	}
	return nil
}

// 获取配送地址
func (p *profileManagerImpl) GetAddress(addressId int64) member.IDeliverAddress {
	v := p.rep.GetSingleDeliverAddress(p.memberId, addressId)
	if v != nil {
		return p.CreateDeliver(v)
	}
	return nil
}

// 删除配送地址
func (p *profileManagerImpl) DeleteAddress(addressId int64) error {
	//todo: 至少保留一个配送地址
	return p.rep.DeleteAddress(p.memberId, addressId)
}

// 拷贝认证信息
func (p *profileManagerImpl) copyTrustedInfo(src, dst *member.TrustedInfo) error {
	//if dst.RealName == src.RealName && dst.CardId == src.CardId &&
	//	dst.TrustImage == src.TrustImage {
	//	return member.ErrNoChangedTrustInfo
	//}
	dst.RealName = src.RealName
	dst.CountryCode = src.CountryCode
	dst.CardId = src.CardId
	dst.CardType = src.CardType
	dst.CardImage = src.CardImage
	dst.TrustImage = src.TrustImage
	return nil
}

// 实名认证信息
func (p *profileManagerImpl) GetTrustedInfo() member.TrustedInfo {
	if p.trustedInfo == nil {
		p.trustedInfo = &member.TrustedInfo{
			MemberId:    p.memberId,
			ReviewState: int(enum.ReviewNotSet),
		}
		//如果还没有实名信息,则新建
		orm := tmp.Db().GetOrm()
		if err := orm.Get(p.memberId, p.trustedInfo); err != nil {
			orm.Save(nil, p.trustedInfo)
		}
	}
	// 显示示例图片
	if p.trustedInfo.TrustImage == "" {
		p.trustedInfo.TrustImage = exampleTrustImageUrl
	}
	return *p.trustedInfo
}

func (p *profileManagerImpl) checkCardId(cardId string, memberId int64) bool {
	mId := 0
	tmp.Db().ExecScalar("SELECT id FROM mm_trusted_info WHERE review_state= $1 AND card_id= $2 AND member_id <> $3 LIMIT 1",
		&mId, enum.ReviewPass, cardId, memberId)
	return mId == 0
}

// 保存实名认证信息
func (p *profileManagerImpl) SaveTrustedInfo(v *member.TrustedInfo) error {
	// 验证数据是否完整
	v.TrustImage = strings.TrimSpace(v.TrustImage)
	v.CardId = strings.TrimSpace(v.CardId)
	v.RealName = strings.TrimSpace(v.RealName)

	if len(v.TrustImage) == 0 || len(v.RealName) == 0 ||
		len(v.CardId) == 0 {
		return member.ErrMissingTrustedInfo
	}
	// 验证姓名
	if !zhNameRegexp.MatchString(v.RealName) {
		return member.ErrRealName
	}
	// 校验身份证号是否正确
	v.CardId = strings.ToUpper(v.CardId)
	err := util.CheckChineseCardID(v.CardId)
	if err != nil {
		return member.ErrTrustCardId
	}
	// 检查身份证是否已被占用
	if !p.checkCardId(v.CardId, p.memberId) {
		return member.ErrCarIdExists
	}
	// 检测上传认证图片
	if v.TrustImage != "" {
		if len(v.TrustImage) < 10 || v.TrustImage == exampleTrustImageUrl {
			return member.ErrTrustMissingImage
		}
	}
	// 保存
	p.GetTrustedInfo()
	err = p.copyTrustedInfo(v, p.trustedInfo)
	if err == nil {
		p.trustedInfo.Remark = ""
		p.trustedInfo.ReviewState = int(enum.ReviewAwaiting) //标记为待处理
		p.trustedInfo.UpdateTime = time.Now().Unix()
		_, err = orm.Save(tmp.Db().GetOrm(), p.trustedInfo,
			int(p.trustedInfo.MemberId))
	}
	return err
}

// 审核实名认证,若重复审核将返回错误
func (p *profileManagerImpl) ReviewTrustedInfo(pass bool, remark string) error {
	p.GetTrustedInfo()
	if pass {
		p.trustedInfo.ReviewState = int(enum.ReviewPass)
		p.member.value.Flag |= member.FlagTrusted
	} else {
		remark = strings.TrimSpace(remark)
		if remark == "" {
			return member.ErrEmptyReviewRemark
		}
		p.trustedInfo.ReviewState = int(enum.ReviewReject)
		if p.member.ContainFlag(member.FlagTrusted) {
			p.member.value.Flag ^= member.FlagTrusted
		}
	}
	p.trustedInfo.Remark = remark
	p.trustedInfo.ReviewTime = time.Now().Unix()
	_, err := orm.Save(tmp.Db().GetOrm(), p.trustedInfo,
		int(p.trustedInfo.MemberId))
	if err == nil {
		if _, err = p.member.Save(); err == nil && pass {
			// 通知实名通过
			msq.Push(msq.MemberTrustInfoPassed, strconv.Itoa(int(p.memberId)),
				fmt.Sprintf("%d|%s|%s",
					p.trustedInfo.CardType,
					p.trustedInfo.CardId,
					p.trustedInfo.RealName))
		}
	}
	return err
}

var _ member.IDeliverAddress = new(addressImpl)

type addressImpl struct {
	_value      *member.Address
	_memberRepo member.IMemberRepo
	_valRepo    valueobject.IValueRepo
}

func newDeliver(v *member.Address, memberRepo member.IMemberRepo,
	valRepo valueobject.IValueRepo) member.IDeliverAddress {
	d := &addressImpl{
		_value:      v,
		_memberRepo: memberRepo,
		_valRepo:    valRepo,
	}
	return d
}

func (p *addressImpl) GetDomainId() int64 {
	return p._value.ID
}

func (p *addressImpl) GetValue() member.Address {
	return *p._value
}

func (p *addressImpl) SetValue(v *member.Address) error {
	if p._value.MemberId == v.MemberId {
		if err := p.checkValue(v); err != nil {
			return err
		}
		p._value = v
	}
	return nil
}

// 设置地区中文名
func (p *addressImpl) renewAreaName(v *member.Address) string {
	//names := p._valRepo.GetAreaNames([]int{
	//	v.Province,
	//	v.City,
	//	v.District,
	//})
	//if names[1] == "市辖区" || names[1] == "市辖县" || names[1] == "县" {
	//	return strings.Join([]string{names[0], names[2]}, " ")
	//}
	//return strings.Join(names, " ")

	return p._valRepo.GetAreaString(v.Province, v.City, v.District)
}

func (p *addressImpl) checkValue(v *member.Address) error {
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

func (p *addressImpl) Save() (int64, error) {
	if err := p.checkValue(p._value); err != nil {
		return p.GetDomainId(), err
	}
	p._value.Area = p.renewAreaName(p._value)
	return p._memberRepo.SaveDeliver(p._value)
}
