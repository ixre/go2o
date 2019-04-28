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
	"errors"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/util"
	"go2o/core/domain"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/domain/tmp"
	dm "go2o/core/infrastructure/domain"
	"strings"
	"time"
)

var _ merchant.IProfileManager = new(profileManagerImpl)

type profileManagerImpl struct {
	*merchantImpl
	valRepo valueobject.IValueRepo
	//企业信息列表
	ent *merchant.EnterpriseInfo
}

func newProfileManager(m *merchantImpl, valRepo valueobject.IValueRepo) merchant.IProfileManager {
	return &profileManagerImpl{
		merchantImpl: m,
		valRepo:      valRepo,
	}
}

// 获取企业信息
func (p *profileManagerImpl) GetEnterpriseInfo() *merchant.EnterpriseInfo {
	if p.ent == nil {
		p.ent = p._rep.GetMchEnterpriseInfo(p.GetAggregateRootId())
	}
	return p.ent
}

func (p *profileManagerImpl) copy(src *merchant.EnterpriseInfo,
	dst *merchant.EnterpriseInfo) {
	// 商户编号
	dst.MchId = p.GetAggregateRootId()
	// 公司名称
	dst.CompanyName = src.CompanyName
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
	// 法人身份证
	dst.PersonIdNo = src.PersonIdNo
	// 身份证验证图片(人捧身份证照相)
	dst.PersonImage = src.PersonImage
	// 营业执照图片
	dst.CompanyImage = src.CompanyImage
	// 授权书
	dst.AuthDoc = src.AuthDoc
}

// 保存企业信息
func (p *profileManagerImpl) SaveEnterpriseInfo(v *merchant.EnterpriseInfo) (int32, error) {
	e := p.GetEnterpriseInfo()
	if e == nil {
		e = &merchant.EnterpriseInfo{}
	}
	p.copy(v, e)
	dt := time.Now().Unix()
	e.Reviewed = enum.ReviewAwaiting
	aName := p.valRepo.GetAreaNames([]int32{e.Province, e.City, e.District})
	e.Location = strings.Join(aName, "")
	e.ReviewTime = dt
	e.UpdateTime = dt
	p.ent = nil //clean cache
	return util.I32Err(p._rep.SaveMchEnterpriseInfo(e))
}

// 标记企业为审核通过
func (p *profileManagerImpl) ReviewEnterpriseInfo(pass bool, message string) error {
	var err error
	e := p.GetEnterpriseInfo()
	if e == nil {
		return errors.New("no such enterprise info for reviewed")
	}
	e.ReviewTime = time.Now().Unix()
	// 通过审核,将审批的记录删除,同时更新到审核数据
	if pass {
		e.Reviewed = enum.ReviewPass
		e.ReviewRemark = ""
		_, err = p._rep.SaveMchEnterpriseInfo(e)
		if err == nil {
			// 保存省、市、区到Merchant
			v := p.merchantImpl.GetValue()
			v.CompanyName = e.CompanyName
			v.Province = e.Province
			v.City = e.City
			v.District = e.District
			err = p.SetValue(&v)
			if err == nil {
				_, err = p.merchantImpl.Save()
			}
		}
	} else {
		e.Reviewed = enum.ReviewReject
		e.ReviewRemark = message
		_, err = p._rep.SaveMchEnterpriseInfo(e)
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
	_, err := orm.Save(tmp.Db().GetOrm(), e, int(e.ID))
	return err
}
