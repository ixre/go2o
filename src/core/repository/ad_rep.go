/**
 * Copyright 2015 @ S1N1 Team.
 * name : ad_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repository
import (
	"go2o/src/core/domain/interface/ad"
	"github.com/atnet/gof/db"
	adImpl "go2o/src/core/domain/ad"
)

var _ ad.IAdvertisementRep = new(advertisementRep)
type advertisementRep struct{
	db.Connector
}

// 广告仓储
func NewAdvertisementRep(c db.Connector) ad.IAdvertisementRep {
	return &advertisementRep{
		Connector: c,
	}
}

// 获取商户的广告管理
func (this *advertisementRep) GetPartnerAdvertisement(partnerId int)ad.IPartnerAdvertisement{
	return adImpl.NewPartnerAdvertisement(partnerId,this)
}

// 保存广告值
func (this *advertisementRep) SaveAdvertisementValue(v *ad.ValueAdvertisement)(int,error){
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0{
		_,_,err = orm.Save(v.Id,v)
	}else{
		_,_,err = orm.Save(nil,v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM pt_ad WHERE partner_id=?",&v.Id,v.PartnerId)
	}
	return v.Id,err
}

// 保存广告图片
func (this *advertisementRep) SaveAdImageValue(v *ad.ValueImage)(int,error){
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0{
		_,_,err = orm.Save(v.Id,v)
	}else{
		_,_,err = orm.Save(nil,v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM pt_ad_image WHERE ad_id=? AND ",&v.Id,v.AdvertisementId)
	}
	return v.Id,err
}

// 获取广告
func (this *advertisementRep) GetValueAdvertisement(partnerId,id int)*ad.ValueAdvertisement{
	var e ad.ValueAdvertisement
	if err := this.Connector.GetOrm().Get(id, &e);
	err == nil && e.PartnerId == partnerId {
		return &e
	}
	return nil
}

// 根据名称获取广告
func (this *advertisementRep) GetValueAdvertisementByName(partnerId int,name string)*ad.ValueAdvertisement{
	var e ad.ValueAdvertisement
	if err := this.Connector.GetOrm().GetBy(&e,"partner_id=? and name=?",partnerId,name);err == nil {
		return &e
	}
	return nil
}