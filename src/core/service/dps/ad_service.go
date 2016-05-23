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
	"go2o/src/core/infrastructure/format"
)

type adService struct {
	_rep ad.IAdRep
	//_query     *query.ContentQuery
}

func NewAdvertisementService(rep ad.IAdRep) *adService {
	return &adService{
		_rep: rep,
	}
}

func (this *adService) getUserAd(adUserId int) ad.IUserAd {
	return this._rep.GetAdManager().GetUserAd(adUserId)
}

// 获取广告
func (this *adService) GetAdvertisement(adUserId, id int) *ad.ValueAdvertisement {
	if adv := this.getUserAd(adUserId).GetById(id); adv != nil {
		return adv.GetValue()
	}
	return nil
}

// 获取广告及广告数据
func (this *adService) GetAdvertisementAndDataByName(adUserId int, name string) (
	*ad.ValueAdvertisement, interface{}) {
	if adv := this.getUserAd(adUserId).GetByName(name); adv != nil {
		v := adv.GetValue()
		switch adv.Type() {
		case ad.TypeGallery:
			gallary := adv.(ad.IGalleryAd).GetEnabledAdValue()
			for _, v := range gallary {
				v.ImageUrl = format.GetResUrl(v.ImageUrl)
			}
			return v, gallary
			//todo: 其他的广告类型
		}
		return v, nil
	}
	return nil, nil
}

// 保存广告
func (this *adService) SaveAdvertisement(adUserId int, v *ad.ValueAdvertisement) (int, error) {
	pa := this.getUserAd(adUserId)
	var adv ad.IAdvertisement
	if v.Id > 0 {
		adv = pa.GetById(v.Id)
		adv.SetValue(v)
	} else {
		adv = pa.CreateAdvertisement(v)
	}
	return adv.Save()
}

func (this *adService) DelAdvertisement(adUserId, adId int) error {
	return this.getUserAd(adUserId).DeleteAdvertisement(adId)
}

// 保存广告图片
func (this *adService) SaveImage(adUserId int, advertisementId int, v *ad.ValueImage) (int, error) {
	pa := this.getUserAd(adUserId)
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
func (this *adService) GetValueAdImage(adUserId, advertisementId, imgId int) *ad.ValueImage {
	pa := this.getUserAd(adUserId)
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
func (this *adService) DelAdImage(adUserId, advertisementId, imgId int) error {
	pa := this.getUserAd(adUserId)
	var adv ad.IAdvertisement = pa.GetById(advertisementId)
	if adv != nil {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			return gad.DelImage(imgId)
		}
	}
	return nil
}
