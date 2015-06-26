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
func (this *GalleryAd) GetAdValue()ad.ValueGallery{
	return this._adValue
}

// 设置广告值
func (this *GalleryAd) SetAdValue(v ad.ValueGallery)error{
	this._adValue = v
	return nil
}

// 保存广告
func (this *GalleryAd) Save()(int,error){
	id,err := this.Advertisement.Save()
	if this._adValue != nil {
		for _, v := range this._adValue {
			v.AdvertisementId = this.GetDomainId()
			this.Rep.SaveAdImageValue(v)
		}
	}
	return id,err
}
