/**
 * Copyright 2015 @ S1N1 Team.
 * name : partner_ad
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad
import "go2o/src/core/domain/interface/ad"

var _ ad.IPartnerAdvertisement = new(PartnerAdvertisement)

type PartnerAdvertisement struct {
	_rep ad.IAdvertisementRep
	_partnerId int
}

func NewPartnerAdvertisement(partnerId int, rep ad.IAdvertisementRep) ad.IPartnerAdvertisement {
	return &PartnerAdvertisement{
		_rep: rep,
		_partnerId:  partnerId,
	}
}

// 获取聚合根标识
func (this *PartnerAdvertisement) GetAggregateRootId() int{
	return this._partnerId
}

// 根据编号获取广告
func (this *PartnerAdvertisement) GetById(id int)ad.IAdvertisement{
	v := this._rep.GetValueAdvertisement(this._partnerId,id)
	if v != nil{
		return this.CreateAdvertisement(v)
	}
	return nil
}

// 根据名称获取广告
func (this *PartnerAdvertisement) GetByName(name string)ad.IAdvertisement{
	v := this._rep.GetValueAdvertisementByName(this._partnerId,name)
	if v != nil{
		return this.CreateAdvertisement(v)
	}
	return nil
}

// 创建广告对象
func (t *PartnerAdvertisement) CreateAdvertisement(v *ad.ValueAdvertisement)ad.IAdvertisement{
   if v.Type == ad.TypeGallery{

   }

	//todo: other ad type
	return nil
}
