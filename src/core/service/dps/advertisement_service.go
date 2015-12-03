/**
 * Copyright 2015 @ z3q.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dps

import (
	"go2o/src/core/domain/interface/ad"
)

type advertisementService struct {
	_rep ad.IAdvertisementRep
	//_query     *query.ContentQuery
}

func NewAdvertisementService(rep ad.IAdvertisementRep) *advertisementService {
	return &advertisementService{
		_rep: rep,
	}
}

// 获取广告
func (this *advertisementService) GetAdvertisement(partnerId, id int) *ad.ValueAdvertisement {
	pa := this._rep.GetPartnerAdvertisement(partnerId)
	if adv := pa.GetById(id); adv != nil {
		return adv.GetValue()
	}
	return nil
}

// 获取广告及广告数据
func (this *advertisementService) GetAdvertisementAndDataByName(partnerId int, name string) (
	*ad.ValueAdvertisement, interface{}) {
	pa := this._rep.GetPartnerAdvertisement(partnerId)
	if adv := pa.GetByName(name); adv != nil {
		v := adv.GetValue()
		switch adv.Type() {
		case ad.TypeGallery:
			return v, adv.(ad.IGalleryAd).GetEnabledAdValue()
			//todo: 其他的广告类型
		}
		return v, nil
	}
	return nil, nil
}

// 保存广告
func (this *advertisementService) SaveAdvertisement(partnerId int, v *ad.ValueAdvertisement) (int, error) {
	pa := this._rep.GetPartnerAdvertisement(partnerId)
	var adv ad.IAdvertisement
	if v.Id > 0 {
		adv = pa.GetById(v.Id)
		adv.SetValue(v)
	} else {
		adv = pa.CreateAdvertisement(v)
	}
	return adv.Save()
}

func (this *advertisementService) DelAdvertisement(partnerId, advertisementId int) error {
	return this._rep.GetPartnerAdvertisement(partnerId).DeleteAdvertisement(advertisementId)
}

// 保存广告图片
func (this *advertisementService) SaveImage(partnerId int, advertisementId int, v *ad.ValueImage) (int, error) {
	pa := this._rep.GetPartnerAdvertisement(partnerId)
	var adv ad.IAdvertisement = pa.GetById(advertisementId)
	if adv != nil {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			return gad.SaveImage(v)
		}
	}
	return -1, nil
}

// 获取广告图片
func (this *advertisementService) GetValueAdImage(partnerId, advertisementId, imgId int) *ad.ValueImage {
	pa := this._rep.GetPartnerAdvertisement(partnerId)
	var adv ad.IAdvertisement = pa.GetById(advertisementId)
	if adv != nil {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			return gad.GetImage(imgId)
		}
	}
	return nil
}

// 删除广告图片
func (this *advertisementService) DelAdImage(partnerId, advertisementId, imgId int) error {
	pa := this._rep.GetPartnerAdvertisement(partnerId)
	var adv ad.IAdvertisement = pa.GetById(advertisementId)
	if adv != nil {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			return gad.DelImage(imgId)
		}
	}
	return nil
}
