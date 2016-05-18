/**
 * Copyright 2015 @ z3q.net.
 * name : partner_ad
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import (
	"go2o/src/core/domain/interface/ad"
	"time"
)

var _ ad.IMerchantAdvertisement = new(MerchantAdvertisement)

type MerchantAdvertisement struct {
	_rep        ad.IAdvertisementRep
	_merchantId int
}

func NewMerchantAdvertisement(merchantId int, rep ad.IAdvertisementRep) ad.IMerchantAdvertisement {
	return &MerchantAdvertisement{
		_rep:        rep,
		_merchantId: merchantId,
	}
}

// 初始化默认的广告位
func (this *MerchantAdvertisement) InitInternalAdvertisements() {
	merchantId := this.GetAggregateRootId()
	unix := time.Now().Unix()

	arr := []*ad.ValueAdvertisement{
		&ad.ValueAdvertisement{
			MerchantId: merchantId,
			Name:       "online-shop-slide",
			Type:       ad.TypeGallery,
			Enabled:    1,
		},
		&ad.ValueAdvertisement{
			MerchantId: merchantId,
			Name:       "app-entry-slide",
			Type:       ad.TypeGallery,
			Enabled:    1,
		},
		&ad.ValueAdvertisement{
			MerchantId: merchantId,
			Name:       "online-mobi-shop-slide",
			Type:       ad.TypeGallery,
			Enabled:    1,
		},
	}

	for _, v := range arr {
		v.IsInternal = 1
		v.UpdateTime = unix
		this.CreateAdvertisement(v).Save()
	}
}

// 获取聚合根标识
func (this *MerchantAdvertisement) GetAggregateRootId() int {
	return this._merchantId
}

// 根据编号获取广告
func (this *MerchantAdvertisement) GetById(id int) ad.IAdvertisement {
	v := this._rep.GetValueAdvertisement(this._merchantId, id)
	if v != nil {
		return this.CreateAdvertisement(v)
	}
	return nil
}

// 删除广告
func (this *MerchantAdvertisement) DeleteAdvertisement(advertisementId int) error {
	adv := this.GetById(advertisementId)
	if adv != nil {

		if adv.System() {
			return ad.ErrInternalDisallow
		}

		err := this._rep.DelAdvertisement(this._merchantId, advertisementId)

		this._rep.DelImageDataForAdvertisement(advertisementId)
		this._rep.DelTextDataForAdvertisement(advertisementId)
		return err

	}
	return nil
}

// 根据名称获取广告
func (this *MerchantAdvertisement) GetByName(name string) ad.IAdvertisement {
	v := this._rep.GetValueAdvertisementByName(this._merchantId, name)
	if v != nil {
		return this.CreateAdvertisement(v)
	}
	return nil
}

// 创建广告对象
func (this *MerchantAdvertisement) CreateAdvertisement(v *ad.ValueAdvertisement) ad.IAdvertisement {
	adv := &Advertisement{
		Rep:   this._rep,
		Value: v,
	}

	// 轮播广告
	if v.Type == ad.TypeGallery {
		return &GalleryAd{
			Advertisement: adv,
		}
	}

	//todo: other ad type
	return adv
}
