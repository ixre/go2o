/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: sub_merchant_manager.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-09-11 17:05:19
 * description: 子商户入网管理
 * history:
 */

package payment

import (
	"errors"
	"fmt"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/domain/eventbus"
)

type subMerchantManagerImpl struct {
	rep      payment.IPaymentRepo
	_mchRepo merchant.IMerchantRepo
	_mmRepo  member.IMemberRepo
}

// InitialMerchant implements payment.ISubMerchantManager.
func (s *subMerchantManagerImpl) InitialMerchant(userType int, userId int) (*payment.PayMerchant, error) {
	if userType != 1 && userType != 2 {
		return nil, errors.New("用户类型错误")
	}
	if userId <= 0 {
		return nil, errors.New("用户编号错误")
	}
	v := s.FindMerchant(userId, userType)
	if v != nil {
		return v, errors.New("商户入网信息已存在")
	}
	code := crypto.Md5([]byte(fmt.Sprintf("%d%d-%d", userType, userId, time.Now().UnixNano())))
	v = &payment.PayMerchant{
		Id:                    userId,
		Code:                  code[8:24],
		UserId:                userId,
		UserType:              userType,
		MchType:               1,
		MchRole:               4, //默认分账接受方
		LicencePic:            "",
		SignName:              "",
		SignType:              1,
		LicenceNo:             "",
		ShortName:             "",
		AccountLicencePic:     "",
		LegalName:             "",
		LegalLicenceType:      0,
		LegalLicenceNo:        "",
		LegalFrontPic:         "",
		LegalBackPic:          "",
		ContactName:           "",
		ContactPhone:          "",
		ContactEmail:          "",
		ContactLicenceNo:      "",
		AccountEmail:          "",
		AccountPhone:          "",
		PrimaryIndustryCode:   "",
		SecondaryIndustryCode: "",
		ProvinceCode:          0,
		CityCode:              0,
		DistrictCode:          0,
		Address:               "",
		SettleDirection:       0,
		SettleBankCode:        "",
		SettleAccountType:     0,
		SettleBankAccount:     "",
		IssueMchNo:            "",
		IssueStatus:           0,
		IssueMessage:          "",
		CreateTime:            0,
		UpdateTime:            0,
	}
	if userType == 1 {
		// 会员
		v.MchType = 2  // 小微商户
		v.SignType = 1 // 个体户
		mm := s._mmRepo.GetMember(int64(userId))
		if mm == nil {
			return nil, errors.New("会员不存在")
		}
		if !mm.ContainFlag(member.FlagTrusted) {
			return nil, errors.New("会员未实名认证")
		}
		auth := mm.Profile().GetCertificationInfo()
		if auth == nil {
			return nil, errors.New("会员未找到实名信息")
		}

		mv := mm.GetValue()

		v.LicenceNo = ""
		v.LicencePic = ""
		v.SignName = auth.RealName
		v.SignType = 0
		v.ShortName = auth.RealName
		v.AccountLicencePic = ""
		v.LegalName = auth.RealName
		v.LegalLicenceType = 1
		v.LegalLicenceNo = auth.CardId
		v.LegalFrontPic = auth.CertFrontPic
		v.LegalBackPic = auth.CertBackPic
		v.ContactName = auth.RealName
		v.ContactPhone = mv.Phone
		v.ContactEmail = mv.Email
		v.ContactLicenceNo = auth.CardId
		v.AccountEmail = mv.Email
		v.AccountPhone = mv.Phone
		v.PrimaryIndustryCode = ""
		v.SecondaryIndustryCode = ""
		v.ProvinceCode = 0
		v.CityCode = 0
		v.DistrictCode = 0
		v.Address = ""
		v.SettleDirection = 2
		v.SettleBankCode = ""
		v.SettleAccountType = 1
		v.SettleBankAccount = ""
	}
	if userType == 2 {
		// 商户
		mch := s._mchRepo.GetMerchant(userId)
		if mch == nil {
			return nil, errors.New("商户不存在")
		}
		auth := mch.ProfileManager().GetAuthenticate()
		if auth == nil {
			return nil, errors.New("商户尚未通过认证")
		}
		mv := mch.GetValue()
		v.MchType = 1  // 企业商户
		v.SignType = 2 // 企业
		v.LicenceNo = auth.LicenceNo
		v.LicencePic = auth.LicencePic
		v.SignName = auth.OrgName
		v.SignType = 0
		v.ShortName = auth.MchName
		v.AccountLicencePic = auth.BankAccountPic
		v.LegalName = auth.PersonName
		v.LegalLicenceType = 1
		v.LegalLicenceNo = auth.PersonId
		v.LegalFrontPic = auth.PersonFrontPic
		v.LegalBackPic = auth.PersonBackPic
		v.ContactName = auth.PersonName
		v.ContactPhone = auth.PersonPhone
		v.ContactEmail = mv.MailAddr
		v.ContactLicenceNo = auth.PersonId
		v.AccountEmail = mv.MailAddr
		v.AccountPhone = auth.PersonPhone
		v.PrimaryIndustryCode = ""
		v.SecondaryIndustryCode = ""
		v.ProvinceCode = auth.Province
		v.CityCode = auth.City
		v.DistrictCode = auth.District
		v.Address = auth.OrgAddress
		v.SettleDirection = 2
		v.SettleBankCode = auth.BankName
		v.SettleAccountType = 1
		v.SettleBankAccount = auth.BankNo
	}
	v.IssueStatus = 0
	v.CreateTime = int(time.Now().Unix())
	v.UpdateTime = v.CreateTime
	return s.rep.MerchantRepo().Save(v)
}

// FindMerchant implements payment.ISubMerchantManager.
func (s *subMerchantManagerImpl) FindMerchant(userType int, userId int) *payment.PayMerchant {
	return s.rep.MerchantRepo().FindBy("user_type = ? and user_id = ?", userType, userId)
}

// GetMerchant implements payment.ISubMerchantManager.
func (s *subMerchantManagerImpl) GetMerchant(code string) *payment.PayMerchant {
	return s.rep.MerchantRepo().FindBy("code = ?", code)
}

// SaveMerchant implements payment.ISubMerchantManager.
func (s *subMerchantManagerImpl) StageMerchant(mch *payment.PayMerchant) error {
	if mch.Id <= 0 {
		return errors.New("无法新增商户入网资料")
	}
	raw := s.GetMerchant(mch.Code)
	if raw == nil {
		return errors.New("商户入网信息不存在")
	}
	if raw.UserType != mch.UserType {
		return errors.New("用户类型不匹配")
	}
	if raw.UserId != mch.UserId {
		return errors.New("用户编号不匹配")
	}
	_, err := s.rep.MerchantRepo().Save(mch)
	return err
}

// Submit implements payment.ISubMerchantManager.
func (s *subMerchantManagerImpl) Submit(code string) error {
	mch := s.GetMerchant(code)
	if mch == nil {
		return errors.New("商户入网信息不存在")
	}
	if mch.IssueStatus == 1 {
		return errors.New("商户入网信息已提交")
	}
	if !(mch.IssueStatus == 0 || mch.IssueStatus == 2) {
		return errors.New("商户入网信息无法提交或已通过审核")
	}
	if mch.MchType == 1 {
		// 企业商户
		if mch.LicencePic == "" {
			return errors.New("商户营业执照图片不能为空")
		}
		if mch.SignName == "" {
			return errors.New("商户名称不能为空")
		}
		if mch.ShortName == "" {
			return errors.New("商户简称不能为空")
		}
		if mch.AccountLicencePic == "" {
			return errors.New("账户信息表未上传")
		}
		if mch.LegalName == "" {
			return errors.New("法人姓名不能为空")
		}
		if mch.LegalLicenceNo == "" {
			return errors.New("法人身份证号不能为空")
		}
		if mch.LegalFrontPic == "" {
			return errors.New("法人身份证正面图片不能为空")
		}
		if mch.LegalBackPic == "" {
			return errors.New("法人身份证反面图片不能为空")
		}
		if mch.ContactName == "" {
			return errors.New("联系人姓名不能为空")
		}
		if mch.ContactPhone == "" {
			return errors.New("联系人电话不能为空")
		}
		if mch.ContactEmail == "" {
			return errors.New("联系人邮箱不能为空")
		}
		if mch.ContactLicenceNo == "" {
			return errors.New("联系人身份证号不能为空")
		}
		if mch.AccountEmail == "" {
			return errors.New("管理员邮箱不能为空")
		}
		if mch.AccountPhone == "" {
			return errors.New("管理员电话不能为空")
		}
		if mch.PrimaryIndustryCode == "" || mch.SecondaryIndustryCode == "" {
			return errors.New("主营行业不能为空")
		}
		if mch.ProvinceCode == 0 || mch.CityCode == 0 || mch.DistrictCode == 0 {
			return errors.New("经营区域不能为空")
		}
		if mch.Address == "" {
			return errors.New("经营地址不能为空")
		}
		if mch.SettleBankCode == "" {
			return errors.New("结算银行未指定")
		}
		if mch.SettleAccountType == 0 {
			return errors.New("结算账户类型不能为空")
		}
		if mch.SettleBankAccount == "" {
			return errors.New("结算账户不能为空")
		}
	}
	mch.IssueStatus = 1
	mch.UpdateTime = int(time.Now().Unix())
	_, err := s.rep.MerchantRepo().Save(mch)
	if err == nil {
		// 发布商户入网事件
		go eventbus.Dispatch(&payment.PaymentMerchantRegistrationEvent{
			Merchant: mch,
		})
	}
	return err
}

func NewSubMerchantManager(rep payment.IPaymentRepo, mchRepo merchant.IMerchantRepo, mmRepo member.IMemberRepo) payment.ISubMerchantManager {
	return &subMerchantManagerImpl{
		rep:      rep,
		_mchRepo: mchRepo,
		_mmRepo:  mmRepo,
	}
}

func (s *subMerchantManagerImpl) Update(code string, data *payment.SubMerchantUpdateParams) error {
	mch := s.GetMerchant(code)
	if mch == nil {
		return errors.New("商户入网信息不存在")
	}
	if mch.IssueStatus == 5 {
		return errors.New("商户入网信息已审核通过")
	}
	if !collections.InArray([]int{2, 3, 4, 5}, data.Status) {
		return errors.New("商户入网信息状态错误")
	}
	mch.IssueStatus = data.Status
	mch.IssueMessage = data.Remark
	mch.IssueMchNo = data.MerchantCode
	mch.AgreementSignUrl = data.AgreementSignUrl
	mch.UpdateTime = int(time.Now().Unix())
	_, err := s.rep.MerchantRepo().Save(mch)
	if err == nil {
		if mch.IssueStatus == 5 {
			// 已通过，则更新商户或个人的入网信息
			if mch.UserType == 2 {
				err = s.updateMerchantSubMerchantNo(mch.UserId, mch.IssueMchNo)
			}
			if mch.UserType == 1 {
				err = s.updateMemberSubMerchantNo(mch.UserId, mch.IssueMchNo)
			}
		}
	}
	return err
}

func (s *subMerchantManagerImpl) updateMerchantSubMerchantNo(userId int, subMerchantNo string) error {
	mch := s._mchRepo.GetMerchant(userId)
	if mch == nil {
		return errors.New("商户不存在")
	}
	mgr := mch.ConfManager()
	conf := mgr.GetSettleConf()
	conf.SubMchNo = subMerchantNo
	return mgr.SaveSettleConf(conf)
}

func (s *subMerchantManagerImpl) updateMemberSubMerchantNo(userId int, subMerchantNo string) error {
	//todo: 暂未实现
	return nil
}
