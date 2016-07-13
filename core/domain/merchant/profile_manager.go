/**
 * Copyright 2015 @ z3q.net.
 * name : profile_manager
 * author : jarryliu
 * date : 2016-05-27 10:38
 * description :
 * history :
 */
package merchant

import (
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/tmp"
	dm "go2o/core/infrastructure/domain"
	"time"
)

var _ merchant.IProfileManager = new(profileManagerImpl)

type profileManagerImpl struct {
	*MerchantImpl
	//企业信息列表
	_list []*merchant.EnterpriseInfo
}

func newProfileManager(m *MerchantImpl) merchant.IProfileManager {
	return &profileManagerImpl{
		MerchantImpl: m,
	}
}

func (this *profileManagerImpl) getAll() []*merchant.EnterpriseInfo {
	if this._list == nil {
		this._list = []*merchant.EnterpriseInfo{}
		tmp.Db().GetOrm().Select(&this._list, "mch_id=?", this.GetAggregateRootId())
	}
	return this._list
}

// 获取企业信息
func (this *profileManagerImpl) GetReviewingEnterpriseInfo() *merchant.EnterpriseInfo {
	for _, v := range this.getAll() {
		if v.Reviewed == 0 {
			return v
		}
	}
	return nil
}

// 获取审核过的企业信息
func (this *profileManagerImpl) GetReviewedEnterpriseInfo() *merchant.EnterpriseInfo {
	for _, v := range this.getAll() {
		if v.Reviewed == 1 {
			return v
		}
	}
	return nil
}

func (this *profileManagerImpl) copy(src *merchant.EnterpriseInfo,
	dst *merchant.EnterpriseInfo) {
	// 商户编号
	dst.MerchantId = this.GetAggregateRootId()

	// 公司名称
	dst.Name = src.Name
	// 公司营业执照编号
	dst.CompanyNo = src.CompanyNo
	// 法人
	dst.PersonName = src.PersonName
	// 公司电话
	dst.Tel = src.Tel
	// 公司地址
	dst.Address = src.Address

	dst.Province = src.Province

	dst.City = src.City

	dst.District = src.District

	dst.Location = src.Location
	// 法人身份证
	dst.PersonIdNo = src.PersonIdNo
	// 身份证验证图片(人捧身份证照相)
	dst.PersonImageUrl = src.PersonImageUrl
	// 营业执照图片
	dst.CompanyImageUrl = src.CompanyImageUrl
	//是否已审核
	//dst.Reviewed = src.Reviewed
	// 审核时间
	//dst.ReviewTime = src.ReviewTime
	// 审核备注
	//dst.Remark = src.Remark
}

// 保存企业信息
func (this *profileManagerImpl) SaveEnterpriseInfo(v *merchant.EnterpriseInfo) (int, error) {
	e := this.GetReviewingEnterpriseInfo()
	if e == nil {
		e = &merchant.EnterpriseInfo{}
	}
	this.copy(v, e)
	dt := time.Now().Unix()
	e.Reviewed = 0
	e.IsHandled = 0
	e.ReviewTime = dt
	e.UpdateTime = dt
	this._list = nil //clean cache
	return orm.Save(tmp.Db().GetOrm(), e, e.Id)
}

// 标记企业为审核通过
func (this *profileManagerImpl) ReviewEnterpriseInfo(pass bool, message string) error {
	var err error
	e := this.GetReviewingEnterpriseInfo()
	if e != nil {
		e.ReviewTime = time.Now().Unix()
		e.IsHandled = 1

		// 通过审核,将审批的记录删除,同时更新到审核数据
		if pass {
			v := this.GetReviewedEnterpriseInfo()
			if v != nil {
				this.copy(e, v)
				tmp.Db().GetOrm().DeleteByPk(e, e.Id)
				err = this.save(v)
			} else {
				e.Reviewed = 1
				e.Remark = ""
				err = this.save(e)
			}
			if err == nil {
				// 保存省、市、区到Merchant
				v := this.MerchantImpl.GetValue()
				v.Province = e.Province
				v.City = e.City
				v.District = e.District
				err = this.SetValue(&v)
				if err == nil {
					_, err = this.MerchantImpl.Save()
				}
			}
		} else {
			e.Reviewed = 0
			e.Remark = message
			err = this.save(e)
		}
	}
	return err
}

// 修改密码
func (this *profileManagerImpl) ModifyPassword(newPwd, oldPwd string) error {
	var err error
	if newPwd == oldPwd {
		return domain.ErrPwdCannotSame
	}
	if b, err := dm.ChkPwdRight(newPwd); !b {
		return err
	}
	if len(oldPwd) != 0 && oldPwd != this._value.Pwd {
		return domain.ErrPwdOldPwdNotRight
	}
	this._value.Pwd = dm.MerchantSha1Pwd(this._value.Usr, newPwd)
	_, err = this.Save()
	return err
}

func (this *profileManagerImpl) save(e *merchant.EnterpriseInfo) error {
	_, err := orm.Save(tmp.Db().GetOrm(), e, e.Id)
	return err
}
