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
	*merchantImpl
	//企业信息列表
	list []*merchant.EnterpriseInfo
}

func newProfileManager(m *merchantImpl) merchant.IProfileManager {
	return &profileManagerImpl{
		merchantImpl: m,
	}
}

func (p *profileManagerImpl) getAll() []*merchant.EnterpriseInfo {
	if p.list == nil {
		p.list = []*merchant.EnterpriseInfo{}
		tmp.Db().GetOrm().Select(&p.list, "mch_id=?", p.GetAggregateRootId())
	}
	return p.list
}

// 获取企业信息
func (p *profileManagerImpl) GetReviewingEnterpriseInfo() *merchant.EnterpriseInfo {
	for _, v := range p.getAll() {
		if v.Reviewed == 0 {
			return v
		}
	}
	return nil
}

// 获取审核过的企业信息
func (p *profileManagerImpl) GetReviewedEnterpriseInfo() *merchant.EnterpriseInfo {
	for _, v := range p.getAll() {
		if v.Reviewed == 1 {
			return v
		}
	}
	return nil
}

func (p *profileManagerImpl) copy(src *merchant.EnterpriseInfo,
	dst *merchant.EnterpriseInfo) {
	// 商户编号
	dst.MerchantId = p.GetAggregateRootId()

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
	dst.PersonImage = src.PersonImage
	// 营业执照图片
	dst.CompanyImage = src.CompanyImage
	// 授权书
	dst.AuthDoc = src.AuthDoc

	//是否已审核
	//dst.Reviewed = src.Reviewed
	// 审核时间
	//dst.ReviewTime = src.ReviewTime
	// 审核备注
	//dst.Remark = src.Remark
}

// 保存企业信息
func (p *profileManagerImpl) SaveEnterpriseInfo(v *merchant.EnterpriseInfo) (int32, error) {
	e := p.GetReviewingEnterpriseInfo()
	if e == nil {
		e = &merchant.EnterpriseInfo{}
	}
	p.copy(v, e)
	dt := time.Now().Unix()
	e.Reviewed = 0
	e.IsHandled = 0
	e.ReviewTime = dt
	e.UpdateTime = dt
	p.list = nil //clean cache
	return orm.I32(orm.Save(tmp.Db().GetOrm(), e, int(e.Id)))
}

// 标记企业为审核通过
func (p *profileManagerImpl) ReviewEnterpriseInfo(pass bool, message string) error {
	var err error
	e := p.GetReviewingEnterpriseInfo()
	if e != nil {
		e.ReviewTime = time.Now().Unix()
		e.IsHandled = 1

		// 通过审核,将审批的记录删除,同时更新到审核数据
		if pass {
			v := p.GetReviewedEnterpriseInfo()
			if v != nil {
				p.copy(e, v)
				tmp.Db().GetOrm().DeleteByPk(e, e.Id)
				err = p.save(v)
			} else {
				e.Reviewed = 1
				e.Remark = ""
				err = p.save(e)
			}
			if err == nil {
				// 保存省、市、区到Merchant
				v := p.merchantImpl.GetValue()
				v.Province = e.Province
				v.City = e.City
				v.District = e.District
				err = p.SetValue(&v)
				if err == nil {
					_, err = p.merchantImpl.Save()
				}
			}
		} else {
			e.Reviewed = 0
			e.Remark = message
			err = p.save(e)
		}
	}
	return err
}

// 修改密码
func (p *profileManagerImpl) ModifyPassword(newPwd, oldPwd string) error {
	if b, err := dm.ChkPwdRight(newPwd); !b {
		return err
	}
	if len(oldPwd) != 0 {
		if newPwd == oldPwd {
			return domain.ErrPwdCannotSame
		}
		if oldPwd != p._value.Pwd {
			return domain.ErrPwdOldPwdNotRight
		}
	}
	p._value.Pwd = newPwd
	_, err := p.Save()
	return err
}

func (p *profileManagerImpl) save(e *merchant.EnterpriseInfo) error {
	_, err := orm.Save(tmp.Db().GetOrm(), e, int(e.Id))
	return err
}
