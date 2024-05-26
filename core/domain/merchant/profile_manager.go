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
	"errors"
	"time"

	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	dm "github.com/ixre/go2o/core/infrastructure/domain"
)

var _ merchant.IProfileManager = new(profileManagerImpl)

type profileManagerImpl struct {
	*merchantImpl
	valRepo valueobject.IValueRepo
}

func newProfileManager(m *merchantImpl, valRepo valueobject.IValueRepo) merchant.IProfileManager {
	return &profileManagerImpl{
		merchantImpl: m,
		valRepo:      valRepo,
	}
}

// SaveAuthenticate implements merchant.IProfileManager.
func (p *profileManagerImpl) SaveAuthenticate(v *merchant.Authenticate) (int, error) {
	err := p.checkAuthenticate(v)
	if err != nil {
		return 0, err
	}
	v.MchId = int(p.GetAggregateRootId())
	v.ReviewStatus = int(enum.ReviewAwaiting)
	v.ReviewRemark = ""
	v.ReviewTime = 0
	// aName := p.valRepo.GetAreaNames([]int32{e.Province, e.City, e.District})
	// e.Location = strings.Join(aName, "")
	v.UpdateTime = int(time.Now().Unix())
	e := p._repo.GetMerchantAuthenticate(p.GetAggregateRootId(), 0)
	if e != nil {
		v.Id = e.Id
	}
	id, err := p._repo.SaveAuthenticate(v)
	if err == nil {
		// 添加待审批标记
		err = p.merchantImpl.GrantFlag(merchant.FlagAuthenticate)
		if err == nil {
			p.merchantImpl.Save()
		}
	}
	return id, err
}

// 检查企业认证信息
func (p *profileManagerImpl) checkAuthenticate(v *merchant.Authenticate) error {
	if v == nil || len(v.OrgName) < 2 {
		return errors.New("企业名称不能为空")
	}
	if len(v.OrgNo) == 0 {
		return errors.New("企业营业执照号不能为空")
	}
	if len(v.OrgPic) == 0 {
		return errors.New("企业营业执照图片不能为空")
	}
	if len(v.PersonName) < 2 {
		return errors.New("法人名称不能为空")
	}
	if len(v.PersonId) != 18 {
		return errors.New("法人身份证号不正确")
	}
	if len(v.PersonPic) == 0 {
		return errors.New("法人身份证照片不能为空")
	}
	// if len(v.AuthorityPic) == 0 {
	// 	return errors.New("未上传授权书")
	// }
	return nil
}

// ReviewAuthenticate 审核商户企业认证信息
func (p *profileManagerImpl) ReviewAuthenticate(pass bool, message string) error {
	var err error
	e := p._repo.GetMerchantAuthenticate(p.GetAggregateRootId(), 0)
	if e == nil {
		return errors.New("未找到企业认证信息")
	}
	if e.ReviewStatus != int(enum.ReviewAwaiting) {
		return errors.New("企业认证信息已审核")
	}
	e.ReviewTime = int(time.Now().Unix())
	// 通过审核,将审批的记录删除,同时更新到审核数据
	if pass {
		e.ReviewStatus = int(enum.ReviewPass)
		e.ReviewRemark = ""
		_, err = p._repo.SaveAuthenticate(e)
		if err == nil {
			// 更新企业认证信息
			err = p.saveMerchantAuthenticate(e)
			if err != nil {
				return err
			}
			// 保存商户信息
			v := p.merchantImpl.GetValue()
			v.MchName = e.OrgName
			if v.Status == 0 {
				// 如果商户状态为待认证,则设置商户已开通
				v.Status = 1
			}
			if err = p.SetValue(&v); err != nil {
				return err
			}
			// 去除待认证标记
			err = p.merchantImpl.GrantFlag(-merchant.FlagAuthenticate)
			if err == nil {
				_, err = p.merchantImpl.Save()
			}
		}
	} else {
		e.ReviewStatus = int(enum.ReviewReject)
		e.ReviewRemark = message
		_, err = p._repo.SaveAuthenticate(e)
	}
	return err
}

// 保存企业认证信息
func (p *profileManagerImpl) saveMerchantAuthenticate(v *merchant.Authenticate) error {
	dst := &merchant.Authenticate{
		Id:               0,
		MchId:            p.GetAggregateRootId(),
		OrgName:          v.OrgName,
		OrgNo:            v.OrgNo,
		OrgPic:           v.OrgPic,
		WorkCity:         v.WorkCity,
		QualificationPic: v.QualificationPic,
		PersonId:         v.PersonId,
		PersonName:       v.PersonName,
		PersonPic:        v.PersonPic,
		PersonPhone:      v.PersonPhone,
		AuthorityPic:     v.AuthorityPic,
		BankName:         v.BankName,
		BankAccount:      v.BankAccount,
		BankNo:           v.BankNo,
		ExtraData:        v.ExtraData,
		ReviewStatus:     v.ReviewStatus,
		ReviewRemark:     "",
		ReviewTime:       v.ReviewTime,
		Version:          1,
		UpdateTime:       v.ReviewTime,
	}
	e := p._repo.GetMerchantAuthenticate(p.GetAggregateRootId(), 1)
	if e != nil {
		dst.Id = e.Id
	}
	_, err := p._repo.SaveAuthenticate(dst)
	return err
}

// ChangePassword 修改密码
func (p *profileManagerImpl) ChangePassword(newPwd, oldPwd string) error {
	if len(newPwd) != 32 {
		return errors.New("密码必须32位Md5")
	}
	if len(oldPwd) != 0 {
		if newPwd == oldPwd {
			return domain.ErrPwdCannotSame
		}
		oldPwd = dm.MerchantSha1Pwd(oldPwd, p.merchantImpl.GetValue().Salt)
		if oldPwd != p._value.Password {
			return domain.ErrPwdOldPwdNotRight
		}
	}
	p._value.Password = dm.MerchantSha1Pwd(newPwd, p.merchantImpl.GetValue().Salt)
	_, err := p.Save()
	return err
}
