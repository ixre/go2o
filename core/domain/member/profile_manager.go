/**
 * Copyright 2015 @ 56x.net.
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
	"regexp"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	dm "github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/domain/util"
	"github.com/ixre/go2o/core/infrastructure/domain/validate"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
	"github.com/ixre/go2o/core/initial/provide"
)

var _ member.IProfileManager = new(profileManagerImpl)
var (
	exampleTrustImageUrl = "res/tru-example.jpg"
	// qqRegex = regexp.MustCompile("^\\d{5,12}$")
	zhNameRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]{2,6}$")
)

type profileManagerImpl struct {
	member        *memberImpl
	memberId      int64
	repo          member.IMemberRepo
	valueRepo     valueobject.IValueRepo
	registryRepo  registry.IRegistryRepo
	bankCards     []member.BankCard
	trustedInfo   *member.CerticationInfo
	profile       *member.Profile
	receiptsCodes []member.ReceiptsCode
}

func (p *profileManagerImpl) ReceiptsCodes() []member.ReceiptsCode {
	if p.receiptsCodes == nil {
		p.receiptsCodes = p.repo.ReceiptsCodes(p.memberId)
	}
	return p.receiptsCodes
}

func (p *profileManagerImpl) SaveReceiptsCode(c *member.ReceiptsCode) error {
	if c.MemberId > 0 && c.MemberId != p.memberId {
		return errors.New("receipts code owner not match")
	}

	c.MemberId = p.memberId
	if len(c.Identity) == 0 {
		return member.ErrReceiptsNoIdentity
	}
	if l := len([]rune(c.Name)); l == 0 {
		return member.ErrReceiptsNoName
	} else if l > 10 {
		return member.ErrReceiptsNameLen
	}
	// 如果未传入ID,对ID赋值
	if c.Id <= 0 {
		for _, v := range p.ReceiptsCodes() {
			if v.Identity == c.Identity {
				c.Id = v.Id
				break
			}
		}
		// 如果新增,默认启用
		if c.Id <= 0 {
			c.State = 1
		}
	}
	_, err := p.repo.SaveReceiptsCode(c, p.memberId)
	if err == nil {
		p.receiptsCodes = nil
	}
	return err
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
		repo:         rep,
		registryRepo: registryRepo,
		valueRepo:    valueRepo,
	}
}

// 手机号码是否占用
func (p *profileManagerImpl) phoneIsExist(phone string) bool {
	return p.repo.CheckPhoneBind(phone, p.memberId)
}

// 验证数据,用v.updateTime > 0 判断是否为新创建用户
func (p *profileManagerImpl) validateProfile(v *member.Profile) error {
	v.Name = strings.TrimSpace(v.Name)
	v.Email = strings.ToLower(strings.TrimSpace(v.Email))
	v.Phone = strings.TrimSpace(v.Phone)
	// 验证昵称
	if len([]rune(v.Name)) < 1 && v.UpdateTime > 0 {
		return member.ErrEmptyNickname
	}
	// 检查区域
	if (v.Province == 0 || v.City == 0 || v.District == 0 ||
		len(v.Address) == 0) && v.UpdateTime > 0 {
		return member.ErrAddress
	}
	// 检查邮箱
	if len(v.Email) != 0 && !validate.IsEmail(v.Email) {
		return member.ErrInvalidEmail
	}
	// 检查手机
	checkPhone := p.registryRepo.Get(registry.MemberCheckPhoneFormat).BoolValue()
	if len(v.Phone) != 0 && checkPhone {
		if !validate.IsPhone(v.Phone) {
			return member.ErrInvalidPhone
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
	dst.Gender = v.Gender
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
		len(v.BirthDay) != 0 && len(v.Address) != 0 && v.Gender != 0 &&
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
		p.profile = p.repo.GetProfile(p.memberId)
	}
	return *p.profile
}

// 保存资料
func (p *profileManagerImpl) SaveProfile(v *member.Profile) error {
	ptr := p.GetProfile()
	err := p.copyProfile(v, &ptr)
	if err == nil {
		ptr.MemberId = p.memberId
		err = p.repo.SaveProfile(&ptr)
		if err == nil {
			// 完善资料通知
			if p.ProfileCompleted() {
				// 标记会员已完善资料
				if !p.member.ContainFlag(member.FlagProfileCompleted) {
					p.member.value.UserFlag |= member.FlagProfileCompleted
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
		return member.ErrInvalidPhone
	}
	used := p.repo.CheckPhoneBind(phone, p.memberId)
	if !used {
		v := p.GetProfile()
		v.Phone = phone
		err := p.repo.SaveProfile(&v)
		if err == nil {
			p.repo.ResetMemberIdCache("phone", p.member.value.Phone)
			//todo: phone as user
			p.member.value.Phone = phone
			_, err = p.member.Save()
		}
		return err
	}
	return member.ErrPhoneHasBind
}

// ChangeNickname 修改昵称
func (p *profileManagerImpl) ChangeNickname(nickname string, limitTime bool) error {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return member.ErrEmptyNickname
	}
	used := p.repo.CheckNicknameIsUse(nickname, p.memberId)
	if !used {
		v := p.GetProfile()
		v.Name = nickname
		err := p.repo.SaveProfile(&v)
		if err == nil {
			p.member.value.Nickname = nickname
			_, err = p.member.Save()
		}
		return err
	}
	return member.ErrPhoneHasBind
}

// 设置头像
func (p *profileManagerImpl) ChangeHeadPortrait(portrait string) error {
	if portrait == "" {
		return member.ErrInvalidHeadPortrait
	}
	v := p.GetProfile()
	if p.profile != nil {
		p.profile.Avatar = portrait
	}
	v.Avatar = portrait
	err := p.repo.SaveProfile(&v)
	if err == nil {
		//todo: phone as user
		p.member.value.Phone = portrait
		_, err = p.member.Save()
	}
	return err
}

// todo: ?? 重构
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
				Type: mss.TypeEmailMessage,
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
			val := &mss.MailMessage{
				Subject: mailTpl.Subject,
				Body:    mailTpl.Body,
			}
			msg := p.member.mssRepo.MessageManager().CreateMessage(v, val)
			//todo:?? data
			var data = map[string]string{
				"Name":       p.profile.Name,
				"InviteCode": p.member.GetValue().UserCode,
			}
			return msg.Send(data)
		}
	}
	return errors.New("no such email template")
}

// ChangePassword 修改密码,旧密码可为空
func (p *profileManagerImpl) ChangePassword(newPassword, oldPwd string) error {
	if b, err := dm.ChkPwdRight(newPassword); !b {
		return err
	}
	if len(oldPwd) != 0 {
		if newPassword == oldPwd {
			return domain.ErrPwdCannotSame
		}
		if oldPwd != p.member.value.Password {
			return domain.ErrPwdOldPwdNotRight
		}
	}
	p.member.value.Password = newPassword
	_, err := p.member.Save()
	return err
}

// ChangeTradePassword 修改交易密码，旧密码可为空
func (p *profileManagerImpl) ChangeTradePassword(newPassword, oldPwd string) error {
	if newPassword == oldPwd {
		return domain.ErrPwdCannotSame
	}
	if b, err := dm.ChkPwdRight(newPassword); !b {
		return err
	}
	// 已经设置过旧密码
	if len(oldPwd) != 0 && p.member.value.TradePassword != oldPwd {
		return domain.ErrPwdOldPwdNotRight
	}
	p.member.value.TradePassword = newPassword
	if p.member.ContainFlag(member.FlagNoTradePasswd) {
		p.member.value.UserFlag ^= member.FlagNoTradePasswd
	}
	_, err := p.member.Save()
	return err
}

// GetBankCards 获取提现银行信息
func (p *profileManagerImpl) GetBankCards() []member.BankCard {
	if p.bankCards == nil {
		p.bankCards = p.repo.BankCards(p.memberId)
		//if p.bank == nil {
		//	p.bank = &member.BankInfo{
		//		MemberId:   p.memberId,
		//		IsLocked:   member.BankNoLock,
		//		State:      0,
		//		UpdateTime: time.Now().Unix(),
		//	}
		//	orm.Save(tmp.Orm, p.bank, 0)
		//}
	}
	return p.bankCards
}

// 获取绑定的银行卡
func (p *profileManagerImpl) GetBankCard(cardNo string) *member.BankCard {
	for _, v := range p.GetBankCards() {
		if v.BankAccount == cardNo {
			return &v
		}
	}
	return nil
}

// 绑定银行信息
func (p *profileManagerImpl) AddBankCard(v *member.BankCard) error {
	if v.MemberId > 0 && v.MemberId != p.memberId {
		return member.ErrNoSuchMember
	}
	if p.bankCardIsExists(v.BankAccount) {
		return member.ErrBankCardIsExists
	}
	trustInfo := p.GetCertificationInfo()
	if trustInfo.ReviewStatus == 0 {
		return member.ErrNotTrusted
	}
	if v.BankAccount == "" || v.BankName == "" {
		return member.ErrBankInfo
	}
	v.AccountName = trustInfo.RealName
	err := p.checkBank(v)
	if err == nil {
		v.CreateTime = time.Now().Unix()
		v.MemberId = p.memberId
		if err = p.repo.SaveBankCard(v); err == nil {
			p.bankCards = nil
		}
	}
	return err
}

// 检查银行信息
func (p *profileManagerImpl) checkBank(v *member.BankCard) error {
	v.BankAccount = strings.TrimSpace(v.BankAccount)
	v.AccountName = strings.TrimSpace(v.AccountName)
	v.Network = strings.TrimSpace(v.Network)
	v.BankName = strings.TrimSpace(v.BankName)
	v.BankCode = strings.TrimSpace(v.BankCode)
	if v.BankName == "" {
		return member.ErrBankName
	}
	if v.AccountName == "" {
		return member.ErrBankAccountName
	}
	if l := len(v.BankAccount); l < 16 || l > 19 {
		return member.ErrBankAccount
	}
	//if v.Network == "" {
	//return member.ErrBankNetwork
	//}
	return nil
}

// 移除银行卡
func (p *profileManagerImpl) RemoveBankCard(cardNo string) error {
	if p.bankCardIsExists(cardNo) {
		err := p.repo.RemoveBankCard(p.memberId, cardNo)
		if err == nil {
			p.bankCards = nil
		}
		return err
	}
	return member.ErrBankNoSuchCard
}

// 创建配送地址
func (p *profileManagerImpl) CreateDeliver(v *member.ConsigneeAddress) member.IDeliverAddress {
	v.MemberId = p.memberId
	return newDeliverAddress(v, p.repo, p.valueRepo)
}

// 获取配送地址
func (p *profileManagerImpl) GetDeliverAddress() []member.IDeliverAddress {
	list := p.repo.GetDeliverAddress(p.memberId)
	var arr = make([]member.IDeliverAddress, len(list))
	for i, v := range list {
		arr[i] = p.CreateDeliver(v)
	}
	return arr
}

// 设置默认地址
func (p *profileManagerImpl) SetDefaultAddress(addressId int64) error {
	for _, v := range p.GetDeliverAddress() {
		vv := v.GetValue()
		o := vv.IsDefault
		if v.GetDomainId() == addressId {
			vv.IsDefault = 1
		} else {
			vv.IsDefault = 0
		}
		if o != vv.IsDefault {
			_, err := p.repo.SaveDeliverAddress(&vv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 获取默认收货地址
func (p *profileManagerImpl) GetDefaultAddress() member.IDeliverAddress {
	list := p.repo.GetDeliverAddress(p.memberId)
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
	v := p.repo.GetSingleDeliverAddress(p.memberId, addressId)
	if v != nil {
		return p.CreateDeliver(v)
	}
	return nil
}

// 删除配送地址
func (p *profileManagerImpl) DeleteAddress(addressId int64) error {
	//todo: 至少保留一个配送地址
	return p.repo.DeleteAddress(p.memberId, addressId)
}

// 拷贝认证信息
func (p *profileManagerImpl) copyCertificationInfo(src member.CerticationInfo, dst *member.CerticationInfo) error {
	if dst == nil {
		dst = &member.CerticationInfo{
			MemberId:     p.memberId,
			ReviewStatus: int(enum.ReviewAwaiting),
		}
	}
	dst.RealName = src.RealName
	dst.CountryCode = src.CountryCode
	dst.CardId = src.CardId
	dst.CardType = src.CardType
	dst.CertImage = src.CertImage
	dst.CertReverseImage = src.CertReverseImage
	dst.ExtraCertFile = src.ExtraCertFile
	dst.ExtraCertExt1 = src.ExtraCertExt1
	dst.ExtraCertExt2 = src.ExtraCertExt2
	dst.TrustImage = src.TrustImage
	dst.ManualReview = src.ManualReview
	return nil
}

// 实名认证信息
func (p *profileManagerImpl) GetCertificationInfo() *member.CerticationInfo {
	if p.trustedInfo == nil {
		p.trustedInfo = p.repo.GetCertificationInfo(int(p.memberId))
		if p.trustedInfo == nil {
			p.trustedInfo = &member.CerticationInfo{
				MemberId:     p.memberId,
				ReviewStatus: int(enum.ReviewAwaiting),
			}
		}
	}
	return p.trustedInfo
}

func (p *profileManagerImpl) checkCardId(cardId string, memberId int64) bool {
	mId := 0
	_db := provide.GetDb()
	_db.ExecScalar(`SELECT COUNT(1) FROM mm_cert_info WHERE 
			review_status= $1 AND card_id= $2 AND member_id <> $3 LIMIT 1`,
		&mId, enum.ReviewPass, cardId, memberId)
	return mId == 0
}

// 保存实名认证信息
func (p *profileManagerImpl) SaveCertificationInfo(v *member.CerticationInfo) error {
	// 验证数据是否完整
	v.CardId = strings.TrimSpace(v.CardId)
	v.RealName = strings.TrimSpace(v.RealName)
	v.TrustImage = strings.TrimSpace(v.TrustImage)
	v.CertImage = strings.TrimSpace(v.CertImage)
	v.CertReverseImage = strings.TrimSpace(v.CertReverseImage)
	v.ExtraCertFile = strings.TrimSpace(v.ExtraCertFile)
	v.ExtraCertExt1 = strings.TrimSpace(v.ExtraCertExt1)
	v.ExtraCertExt2 = strings.TrimSpace(v.ExtraCertExt2)
	if len(v.RealName) == 0 || len(v.CardId) == 0 {
		return member.ErrMissingCertificationInfo
	}
	// 验证姓名
	if !zhNameRegexp.MatchString(v.RealName) {
		return member.ErrRealName
	}
	// 校验身份证号是否正确
	v.CardId = strings.ToUpper(v.CardId)
	err := util.CheckChineseCardID(v.CardId)
	if err != nil {
		return member.ErrCertCardId
	}
	// 检查身份证是否已被占用
	if !p.checkCardId(v.CardId, p.memberId) {
		return member.ErrCarIdExists
	}
	// 检测上传认证图片
	requirePeopleImg := p.registryRepo.Get(registry.MemberTrustRequirePeopleImage).BoolValue()
	if v.TrustImage != "" {
		if len(v.TrustImage) < 10 || v.TrustImage == exampleTrustImageUrl {
			return member.ErrTrustMissingImage
		}
	} else if requirePeopleImg {
		return member.ErrTrustMissingImage
	}
	// 检测证件照片
	requireCardImg := p.registryRepo.Get(registry.MemberTrustRequireCardImage).BoolValue()
	if v.CertImage != "" {
		if len(v.CertImage) < 10 {
			return member.ErrTrustMissingCardImage
		}
	} else if requireCardImg {
		return member.ErrTrustMissingCardImage
	}
	// 是否不需要审核, 自动通过
	autoReviewPass := v.ManualReview == 0
	if v.ManualReview < 0 {
		// 如果外部未指定(默认值),则根据配置决定是否需要审核
		b, _ := p.registryRepo.GetValue(registry.MemberCertificationReviewOff)
		autoReviewPass = b == "1"
		v.ManualReview = types.Ternary(autoReviewPass, 0, 1)
	}
	// 保存
	current := p.GetCertificationInfo()
	err = p.copyCertificationInfo(*v, current)
	if err == nil {
		current.Remark = ""
		current.ReviewStatus = int(enum.ReviewAwaiting) //标记为待处理
		current.UpdateTime = time.Now().Unix()
		p.trustedInfo = current
		_, err = p.repo.SaveCertificationInfo(p.trustedInfo)
		if err == nil {
			if autoReviewPass {
				// 自动审核
				err = p.ReviewCertification(true, "自动审核通过")
			}
		}
		return err
	}
	return err
}

// 审核实名认证,若重复审核将返回错误
func (p *profileManagerImpl) ReviewCertification(pass bool, remark string) error {
	p.GetCertificationInfo()
	if pass {
		p.trustedInfo.ReviewStatus = int(enum.ReviewPass)
		p.member.value.UserFlag |= member.FlagTrusted
		p.member.value.RealName = p.trustedInfo.RealName
	} else {
		remark = strings.TrimSpace(remark)
		if remark == "" {
			return member.ErrEmptyReviewRemark
		}
		p.trustedInfo.ReviewStatus = int(enum.ReviewReject)
		if p.member.ContainFlag(member.FlagTrusted) {
			p.member.value.UserFlag ^= member.FlagTrusted
		}
	}
	unix := time.Now().Unix()
	p.trustedInfo.Remark = remark
	p.trustedInfo.ReviewTime = unix
	_, err := p.repo.SaveCertificationInfo(p.trustedInfo)
	if err == nil {
		_, err = p.member.Save()
	}
	return err
}

// 银行卡是否绑定
func (p *profileManagerImpl) bankCardIsExists(cardNo string) bool {
	for _, v := range p.GetBankCards() {
		if v.BankAccount == cardNo {
			return true
		}
	}
	return false
}

// BindOAuthApp implements member.IProfileManager
func (p *profileManagerImpl) BindOAuthApp(app string, openId string, authToken string) error {
	b := p.GetOAuthBindInfo(app)
	if b != nil {
		return fmt.Errorf("app %s aready has binding info", app)
	}
	b = &member.OAuthAccount{
		MemberId:   p.memberId,
		AppCode:    app,
		OpenId:     openId,
		AuthToken:  authToken,
		HeadImgUrl: "",
		UpdateTime: time.Now().Unix(),
	}
	_, err := p.repo.SaveOAuthAccount(b)
	return err
}

// GetOAuthBindInfo implements member.IProfileManager
func (p *profileManagerImpl) GetOAuthBindInfo(app string) *member.OAuthAccount {
	return p.repo.GetOAuthAccount(int(p.memberId), app)
}

// UnbindOAuthApp implements member.IProfileManager
func (p *profileManagerImpl) UnbindOAuthApp(app string) error {
	b := p.GetOAuthBindInfo(app)
	if b != nil {
		return p.repo.DeleteOAuthAccount(b.Id)
	}
	return nil
}
