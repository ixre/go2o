/**
 * Copyright 2015 @ S1N1 Team.
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

type GalleryAd struct{
	*Advertisement
	_adValue ad.ValueGallery
}



// 获取广告值
func (this *GalleryAd) GetAdValue()ad.ValueGallery {
	if this._adValue == nil {
		if this.GetDomainId() > 0 {
			this._adValue =  this.Rep.GetValueGallery(this.GetDomainId())
		}else{
			this._adValue = []*ad.ValueImage{}
		}
	}
	return this._adValue
}

// 保存广告图片
func (this *GalleryAd) SaveImage(v *ad.ValueImage)(int,error){
	v.AdvertisementId = this.GetDomainId()
	return this.Rep.SaveAdImageValue(v)
}