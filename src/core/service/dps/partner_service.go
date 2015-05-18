/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-19 22:49
 * description :
 * history :
 */

package dps

import (
	"errors"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/infrastructure/log"
	"go2o/src/core/query"
	"go2o/src/core/domain/interface/valueobject"
)

type partnerService struct {
	_partnerRep partner.IPartnerRep
	_query      *query.PartnerQuery
}

func NewPartnerService(r partner.IPartnerRep, q *query.PartnerQuery) *partnerService {
	return &partnerService{
		_partnerRep: r,
		_query:      q,
	}
}

// 验证用户密码并返回编号
func (this *partnerService) Verify(usr, pwd string) int {
	ep := domain.Md5PartnerPwd(usr, pwd)
	return this._query.Verify(usr, ep)
}

func (this *partnerService) GetPartner(partnerId int) (*partner.ValuePartner, error) {
	pt, err := this._partnerRep.GetPartner(partnerId)

	if pt != nil {
		v := pt.GetValue()
		return &v, err
	}
	return nil, err
}

func (this *partnerService) SavePartner(partnerId int, v *partner.ValuePartner) (int, error) {
	var pt partner.IPartner
	var err error
	v.Id = partnerId

	if partnerId > 0 {
		pt, _ = this._partnerRep.GetPartner(partnerId)
		if pt == nil {
			err = errors.New("no such partner")
		} else {
			err = pt.SetValue(v)
		}
	} else {
		pt, err = this._partnerRep.CreatePartner(v)
	}

	if err != nil {
		return 0, err
	}

	return pt.Save()
}

// 根据主机查询商户编号
func (this *partnerService) GetPartnerIdByHost(host string) int {
	return this._query.QueryPartnerIdByHost(host)
}

// 获取商户的域名
func (this *partnerService) GetPartnerMajorHost(partnerId int) string {
	pt, err := this._partnerRep.GetPartner(partnerId)
	if err != nil {
		log.PrintErr(err)
	}
	return pt.GetMajorHost()
}

func (this *partnerService) SaveSiteConf(partnerId int, v *partner.SiteConf) error {
	pt, _ := this._partnerRep.GetPartner(partnerId)
	return pt.SaveSiteConf(v)
}

func (this *partnerService) SaveSaleConf(partnerId int, v *partner.SaleConf) error {
	pt, _ := this._partnerRep.GetPartner(partnerId)
	return pt.SaveSaleConf(v)
}

func (this *partnerService) GetSaleConf(partnerId int) *partner.SaleConf {
	pt, err := this._partnerRep.GetPartner(partnerId)
	if err != nil {
		log.PrintErr(err)
	}
	conf := pt.GetSaleConf()
	return &conf
}

func (this *partnerService) GetSiteConf(partnerId int) *partner.SiteConf {
	pt, err := this._partnerRep.GetPartner(partnerId)
	if err != nil {
		log.PrintErr(err)
	}
	conf := pt.GetSiteConf()
	return &conf
}

func (this *partnerService) GetShopsOfPartner(partnerId int) []*partner.ValueShop {
	pt, err := this._partnerRep.GetPartner(partnerId)
	if err != nil {
		log.PrintErr(err)
	}
	shops := pt.GetShops()
	sv := make([]*partner.ValueShop, len(shops))
	for i, v := range shops {
		vv := v.GetValue()
		sv[i] = &vv
	}
	return sv
}

func (this *partnerService) GetShopValueById(partnerId, shopId int) *partner.ValueShop {
	pt, err := this._partnerRep.GetPartner(partnerId)
	if err != nil {
		log.PrintErr(err)
	}
	v := pt.GetShop(shopId).GetValue()
	return &v
}

func (this *partnerService) SaveShop(partnerId int, v *partner.ValueShop) (int, error) {
	pt, err := this._partnerRep.GetPartner(partnerId)
	if err != nil {
		log.PrintErr(err)
		return 0, err
	}
	var shop partner.IShop
	if v.Id > 0 {
		shop = pt.GetShop(v.Id)
		if shop == nil {
			return 0, errors.New("门店不存在")
		}
	} else {
		shop = pt.CreateShop(v)
	}
	err = shop.SetValue(v)
	if err != nil {
		return 0, err
	}
	return shop.Save()
}

func (this *partnerService) DeleteShop(partnerId, shopId int) error {
	pt, err := this._partnerRep.GetPartner(partnerId)
	if err != nil {
		log.PrintErr(err)
	}
	return pt.DeleteShop(shopId)
}

func (this *partnerService) GetPartnersId() []int {
	return this._partnerRep.GetPartnersId()
}

// 保存API信息
func (this *partnerService) SaveApiInfo(partnerId int, d *partner.ApiInfo) error {
	pt, _ := this._partnerRep.GetPartner(partnerId)
	return pt.SaveApiInfo(d)
}

// 获取API接口
func (this *partnerService) GetApiInfo(partnerId int) *partner.ApiInfo {
	pt, _ := this._partnerRep.GetPartner(partnerId)
	v := pt.GetApiInfo()
	return &v
}

// 根据API ID获取PartnerId
func (this *partnerService) GetPartnerIdByApiId(apiId string) int {
	return this._partnerRep.GetPartnerIdByApiId(apiId)
}

func (this *partnerService) GetMemberLevels(partnerId int)[]*valueobject.MemberLevel{
	pt,_ := this._partnerRep.GetPartner(partnerId)
	return pt.LevelManager().GetLevelSet()
}

// 根据编号获取会员等级信息
func (this *partnerService) GetMemberLevelById(partnerId,id int)*valueobject.MemberLevel{
	pt,_ := this._partnerRep.GetPartner(partnerId)
	return pt.LevelManager().GetLevelById(id)
}

// 保存会员等级信息
func (this *partnerService)  SaveMemberLevel(partnerId int,v *valueobject.MemberLevel)(int,error){
	pt,_ :=this._partnerRep.GetPartner(partnerId)
	return pt.LevelManager().SaveLevel(v)
}

// 删除会员等级
func (this *partnerService) DelMemberLevel(partnerId,levelId int)error{
	pt,_ :=this._partnerRep.GetPartner(partnerId)
	return pt.LevelManager().DeleteLevel(levelId)
}

// 获取等级
func (this *partnerService) GetLevel(partnerId,levelId int)*valueobject.MemberLevel{
	pt,_ := this._partnerRep.GetPartner(partnerId)
	return pt.LevelManager().GetLevelById(levelId)
}

// 获取下一个等级
func (this *partnerService) GetNextLevel(partnerId,levelValue int) *valueobject.MemberLevel {
	pt,_ := this._partnerRep.GetPartner(partnerId)
	return pt.LevelManager().GetNextLevel(levelValue)
}