/**
 * Copyright 2015 @ z3q.net.
 * name : gallery
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import (
	"go2o/core/domain/interface/ad"
	"sort"
)

var _ ad.IAd = new(GalleryAd)
var _ ad.IGalleryAd = new(GalleryAd)

type GalleryAd struct {
	*adImpl
	_adValue ad.ValueGallery
}

// 获取广告数据
func (this *GalleryAd) GetAdValue() ad.ValueGallery {
	if this._adValue == nil {
		if this.GetDomainId() > 0 {
			this._adValue = this._rep.GetValueGallery(this.GetDomainId())
			sort.Sort(this._adValue)
		} else {
			this._adValue = []*ad.Image{}
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
func (this *GalleryAd) SaveImage(v *ad.Image) (int, error) {
	v.AdId = this.GetDomainId()
	return this._rep.SaveAdImageValue(v)
}

// 获取图片项
func (this *GalleryAd) GetImage(id int) *ad.Image {
	return this._rep.GetValueAdImage(this.GetDomainId(), id)
}

// 删除图片项
func (this *GalleryAd) DelImage(id int) error {
	return this._rep.DelAdImage(this.GetDomainId(), id)
}

// 转换为数据传输对象
func (this *GalleryAd) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:   this.adImpl.GetDomainId(),
		Type: this.adImpl.Type(),
		Data: this.GetEnabledAdValue(),
	}
}
