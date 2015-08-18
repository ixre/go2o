/**
 * Copyright 2015 @ z3q.net.
 * name : gallery
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import "go2o/src/core/domain/interface/ad"

var _ ad.IAdvertisement = new(GalleryAd)
var _ ad.IGalleryAd = new(GalleryAd)

type GalleryAd struct {
	*Advertisement
	_adValue ad.ValueGallery
}

// 获取广告数据
func (this *GalleryAd) GetAdValue() ad.ValueGallery {
	if this._adValue == nil {
		if this.GetDomainId() > 0 {
			this._adValue = this.Rep.GetValueGallery(this.GetDomainId())
		} else {
			this._adValue = []*ad.ValueImage{}
		}
	}
	return this._adValue
}

// 获取可用的广告数据
func (this *GalleryAd) GetEnabledAdValue() ad.ValueGallery {
	val := this.GetAdValue()
	newVal := ad.ValueGallery{}
	for _, v := range val {
		if v.Enabled == 1 {
			newVal = append(newVal, v)
		}
	}
	return newVal
}

// 保存广告图片
func (this *GalleryAd) SaveImage(v *ad.ValueImage) (int, error) {
	v.AdvertisementId = this.GetDomainId()
	return this.Rep.SaveAdImageValue(v)
}

// 获取图片项
func (this *GalleryAd) GetImage(id int) *ad.ValueImage {
	return this.Rep.GetValueAdImage(this.GetDomainId(), id)
}

// 删除图片项
func (this *GalleryAd) DelImage(id int) error {
	return this.Rep.DelAdImage(this.GetDomainId(), id)
}
