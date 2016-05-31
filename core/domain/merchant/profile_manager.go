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
	"go2o/core/domain"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/tmp"
	"time"
)

var _ merchant.IProfileManager = new(profileManagerImpl)

type profileManagerImpl struct {
	*MerchantImpl
}

func newProfileManager(m *MerchantImpl) merchant.IProfileManager {
	return &profileManagerImpl{
		MerchantImpl: m,
	}
}

// 获取企业信息
func (this *profileManagerImpl) GetEnterpriseInfo() merchant.EnterpriseInfo {
	// 优先获取未审核的,若没有返回审核通过的
	orm := tmp.Db().GetOrm()
	e := merchant.EnterpriseInfo{
		MerchantId: this.GetAggregateRootId(),
	}
	err := orm.GetBy(&e, "mch_id=? AND reviewed=0", this.GetAggregateRootId())
	if err != nil {
		return this.GetReviewedEnterpriseInfo()
	}
	return e
}

// 获取审核过的企业信息
func (this *profileManagerImpl) GetReviewedEnterpriseInfo() merchant.EnterpriseInfo {
	orm := tmp.Db().GetOrm()
	e := merchant.EnterpriseInfo{
		MerchantId: this.GetAggregateRootId(),
	}
	orm.GetBy(&e, "mch_id=? AND reviewed=1", this.GetAggregateRootId())
	return e
}

func (this *profileManagerImpl) copy(src *merchant.EnterpriseInfo,
	dst *merchant.EnterpriseInfo) {
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

	// 身份证验证图片(人捧身份证照相)
	dst.PersonImageUrl = src.PersonImageUrl
	// 营业执照图片
	dst.CompanyImageUrl = src.CompanyImageUrl
	//是否已审核
	dst.Reviewed = src.Reviewed
	// 审核时间
	dst.ReviewTime = src.ReviewTime
	// 审核备注
	dst.Remark = src.Remark
}

// 保存企业信息
func (this *profileManagerImpl) SaveEnterpriseInfo(v *merchant.EnterpriseInfo) (int, error) {
	e := this.GetEnterpriseInfo()
	this.copy(v, &e)
	// 如已审核,则创建一个待审核
	dt := time.Now().Unix()
	if e.Reviewed == 1 {
		e.Id = 0
		e.Reviewed = 0
		e.ReviewTime = dt
	}
	e.UpdateTime = dt

	// ============================
	orm := tmp.Db().GetOrm()
	var err error
	var id int64
	if e.Id > 0 {
		_, _, err = orm.Save(e.Id, &e)
	} else {
		_, id, err = orm.Save(nil, &e)
		e.Id = int(id)
	}
	return e.Id, err
}

// 标记企业为审核通过
func (this *profileManagerImpl) ReviewEnterpriseInfo(pass bool, message string) error {
	e := this.GetEnterpriseInfo()
	if e.Reviewed == 1 {
		return domain.ErrReviewed //已经审核通过
	}

	e.ReviewTime = time.Now().Unix()
	// 通过审核,将审批的记录删除,同时更新到审核数据
	if pass {
		var err error
		e.Reviewed = 1
		e.Remark = ""
		e2 := this.GetReviewedEnterpriseInfo()
		this.copy(&e, &e2)

		// 删除待审核数据
		err = tmp.Db().GetOrm().DeleteByPk(merchant.EnterpriseInfo{}, e.Id)
		if err != nil {
			return err
		}

		//将e2作为新的保存
		if err = this.save(&e2); err == nil {
			// 保存省、市、区到Merchant
			v := this.GetValue()
			v.Province = e2.Province
			v.City = e2.City
			v.District = e2.District
			err = this.SetValue(&v)
			if err == nil {
				_, err = this.Save()
			}
		}

		return err
	}

	e.Reviewed = 0
	e.Remark = message
	return this.save(&e)
}

func (this *profileManagerImpl) save(e *merchant.EnterpriseInfo) error {
	//==================================
	var err error
	orm := tmp.Db().GetOrm()
	if e.Id > 0 {
		_, _, err = orm.Save(e.Id, e)
	} else {
		_, _, err = orm.Save(nil, e)
	}
	return err
}
