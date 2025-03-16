/**
 * Copyright 2015 @ 56x.net.
 * name : profile_manager
 * author : jarryliu
 * date : 2016-05-27 10:38
 * description :
 * history :
 */
package merchant

import (
	"time"

	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	dm "github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/i18n"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/go2o/core/infrastructure/regex"
)

var _ merchant.IProfileManager = new(profileManagerImpl)

type profileManagerImpl struct {
	*merchantImpl
	valRepo      valueobject.IValueRepo
	_invoiceRepo invoice.IInvoiceRepo
}

// GetAuthenticate implements merchant.IProfileManager.
func (p *profileManagerImpl) GetAuthenticate() *merchant.Authenticate {
	return p._repo.GetMerchantAuthenticate(p.GetAggregateRootId(), 1)
}

func newProfileManager(m *merchantImpl, valRepo valueobject.IValueRepo,
	invoiceRepo invoice.IInvoiceRepo,
) merchant.IProfileManager {
	return &profileManagerImpl{
		merchantImpl: m,
		valRepo:      valRepo,
		_invoiceRepo: invoiceRepo,
	}
}

func (p *profileManagerImpl) getStagingAuthenticate() *merchant.Authenticate {
	return p._repo.GetMerchantAuthenticate(p.GetAggregateRootId(), 0)
}

// SaveAuthenticate implements merchant.IProfileManager.
func (p *profileManagerImpl) SaveAuthenticate(v *merchant.Authenticate) (int, error) {
	if v.Id > 0 && v.Version > 0 {
		return 0, i18n.Errorf("已经通过认证,不能再次进行提交认证")
	}
	err := p.checkAuthenticate(v)
	if err != nil {
		return 0, err
	}
	v.MchId = int(p.GetAggregateRootId())

	v.ReviewStatus = int(enum.ReviewPending)
	v.ReviewRemark = ""
	v.ReviewTime = 0
	// aName := p.valRepo.GetDistrictNames([]int32{e.Province, e.City, e.District})
	// e.Location = strings.Join(aName, "")

	v.SubmitTime = int(time.Now().Unix())
	v.UpdateTime = v.SubmitTime
	auth := p.getStagingAuthenticate()
	if auth != nil {
		v.Id = auth.Id
	}
	id, err := p._repo.SaveAuthenticate(v)
	if err == nil {
		err = p.applyMerchantWaitAuthStatus()
	}
	return id, err
}

// 设置商户为待认证状态
func (p *profileManagerImpl) applyMerchantWaitAuthStatus() error {
	// 添加待审批标记
	err := p.merchantImpl.GrantFlag(merchant.FlagWaitAuthenticate)
	if err == nil {
		_, err = p.merchantImpl.Save()
	}
	return err
}

// 检查企业认证信息
func (p *profileManagerImpl) checkAuthenticate(v *merchant.Authenticate) error {
	if v == nil {
		return i18n.Errorf("商户认证信息不能为空")
	}
	if len(v.MchName) < 2 {
		return i18n.Errorf("商户名称不能为空")
	}
	if p._repo.IsExistsMerchantName(v.MchName, p.GetAggregateRootId()) {
		return i18n.Errorf("商户名称已被使用")
	}
	if len(v.OrgName) < 2 {
		return i18n.Errorf("企业名称不能为空")
	}
	if p._repo.IsExistsOrganizationName(v.OrgName, p.GetAggregateRootId()) {
		return i18n.Errorf("企业名称已被使用")
	}
	if v.Province == 0 || v.City == 0 || v.District == 0 {
		return i18n.Errorf("请选择所在地区")
	}
	if len(v.LicenceNo) == 0 {
		return i18n.Errorf("企业营业执照号不能为空")
	}
	if len(v.LicencePic) == 0 {
		return i18n.Errorf("企业营业执照图片不能为空")
	}
	if len(v.PersonName) < 2 {
		return i18n.Errorf("法人名称不能为空")
	}
	if len(v.PersonId) != 18 {
		return i18n.Errorf("法人身份证号不正确")
	}
	if len(v.PersonFrontPic) == 0 {
		return i18n.Errorf("法人身份证正面照片不能为空")
	}
	if len(v.PersonBackPic) == 0 {
		return i18n.Errorf("负责人身份证背面照片不能为空")
	}
	if !regex.IsPhone(v.PersonPhone) {
		return i18n.Errorf("负责人联系电话不正确")
	}
	if len(v.ContactName) == 0 {
		return i18n.Errorf("联系人姓名不能为空")
	}
	if !regex.IsPhone(v.ContactPhone) {
		return i18n.Errorf("联系人电话不正确")
	}
	// if len(v.AuthorityPic) == 0 {
	// 	return i18n.Errorf("未上传授权书")
	// }
	return nil
}

// ReviewAuthenticate 审核商户企业认证信息
func (p *profileManagerImpl) ReviewAuthenticate(pass bool, message string) error {
	e := p._repo.GetMerchantAuthenticate(p.GetAggregateRootId(), 0)
	if e == nil {
		if !pass {
			// 驳回已审核
			return p.rejectReviewedAuthenticate(message)
		}
		// 只对待审核的资料进行审核
		return i18n.Errorf("无法进行审核通过操作")
	}
	if e.ReviewStatus != int(enum.ReviewPending) {
		return i18n.Errorf("商户认证信息已审核")
	}
	var err error
	e.ReviewTime = int(time.Now().Unix())
	// 通过审核,将审批的记录删除,同时更新到审核数据
	if pass {
		e.ReviewStatus = int(enum.ReviewApproved)
		e.ReviewRemark = ""
		// 更新企业认证信息
		err = p.saveMerchantApprovedAuthenticate(e)
	} else {
		e.ReviewStatus = int(enum.ReviewRejected)
		e.ReviewRemark = message
		_, err = p._repo.SaveAuthenticate(e)
		if err == nil {
			// 添加待认证标志
			err = p.merchantImpl.GrantFlag(merchant.FlagWaitAuthenticate)
			if err == nil {
				_, err = p.merchantImpl.Save()
			}
		}
	}
	return err
}

// rejectReviewedAuthenticate 驳回已审核商户认证信息
func (p *profileManagerImpl) rejectReviewedAuthenticate(message string) error {
	e := p._repo.GetMerchantAuthenticate(p.GetAggregateRootId(), 1)
	if e == nil {
		return i18n.Errorf("商户未提交认证信息")
	}
	e.ReviewStatus = int(enum.ReviewRejected)
	e.ReviewRemark = message
	e.Version = 0
	_, err := p._repo.SaveAuthenticate(e)
	if err == nil {
		// 添加待认证标志
		err = p.merchantImpl.GrantFlag(merchant.FlagWaitAuthenticate)
		if err == nil {
			_, err = p.merchantImpl.Save()
		}
	}
	return err
}

func (p *profileManagerImpl) initInvoiceTitle(e *merchant.Authenticate) error {
	var err error
	tenant := p._invoiceRepo.CreateTenant(&invoice.InvoiceTenant{
		TenantType: int(invoice.TenantMerchant),
		TenantUid:  p.merchantImpl.GetAggregateRootId(),
	})
	if tenant == nil {
		err = i18n.Errorf("创建开票租户失败")
	} else {
		err = tenant.CreateInvoiceTitle(&invoice.InvoiceTitle{
			InvoiceType: invoice.InvoiceTypeNormal,
			IssueType:   invoice.IssueTypeEnterprise,
			TitleName:   e.OrgName,
			TaxCode:     e.LicenceNo,
			SignAddress: e.OrgAddress,
			SignTel:     e.PersonPhone,
			BankName:    e.BankName,
			BankAccount: e.BankAccount,
			Remarks:     "",
			IsDefault:   1,
			CreateTime:  0,
		})
	}
	if err != nil {
		logger.Error("创建开票租户失败: mchId: %d, 错误:%s", p.merchantImpl.GetAggregateRootId(), err.Error())
	}
	return err
}

// 保存企业认证信息
func (p *profileManagerImpl) saveMerchantApprovedAuthenticate(v *merchant.Authenticate) error {
	// 删除之前已认证过的信息
	p._repo.DeleteOthersAuthenticate(p.GetAggregateRootId(), v.Id)
	// 将当前信息作为审核通过的信息
	v.Version = 1
	_, err := p._repo.SaveAuthenticate(v)
	if err == nil {
		// 更新商户信息
		mch := p.merchantImpl.GetValue()
		mch.MchName = v.MchName
		mch.FullName = v.OrgName
		mch.Address = v.OrgAddress
		mch.Province = v.Province
		mch.City = v.City
		mch.District = v.District
		if len(mch.Tel) == 0 {
			mch.Tel = v.PersonPhone
		}
		if mch.Status == 0 {
			// 如果商户状态为待认证,则设置商户已开通
			mch.Status = 1
		}
		err = p.SetValue(&mch)
		if err != nil {
			return err
		}
		// 去除待认证标记
		err = p.merchantImpl.GrantFlag(-merchant.FlagWaitAuthenticate)
		if err == nil {
			_, err = p.merchantImpl.Save()
		}
		// 为商户添加开票信息
		p.initInvoiceTitle(v)
	}
	return err
}

// ChangePassword 修改密码
func (p *profileManagerImpl) ChangePassword(newPassword, oldPassword string) error {
	if len(newPassword) != 32 {
		return i18n.Errorf("密码必须32位Md5")
	}
	if len(oldPassword) != 0 {
		if newPassword == oldPassword {
			return domain.ErrPwdCannotSame
		}
		oldPassword = dm.MerchantSha265Pwd(oldPassword, p.merchantImpl.GetValue().Salt)
		if oldPassword != p._value.Password {
			return domain.ErrPwdOldPwdNotRight
		}
	}
	p._value.Password = dm.MerchantSha265Pwd(newPassword, p.merchantImpl.GetValue().Salt)
	_, err := p.Save()
	return err
}
