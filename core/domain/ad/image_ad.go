/**
 * Copyright 2015 @ z3q.net.
 * name : image_ad
 * author : jarryliu
 * date : 2016-05-25 11:29
 * description :
 * history :
 */
package ad

import (
	"go2o/core/domain/interface/ad"
)

var _ ad.IImageAd = new(ImageAdImpl)

type ImageAdImpl struct {
	_extValue *ad.Image
	*AdImpl
}

// 获取链接广告值
func (this *ImageAdImpl) getData() *ad.Image {
	if this._extValue == nil {
		gallery := this._rep.GetValueGallery(this.GetDomainId())
		if gallery.Len() > 0 {
			this._extValue = gallery[0]
		}

		//如果不存在,则创建一个新的对象
		if this._extValue == nil {
			this._extValue = &ad.Image{
				AdId: this.GetDomainId(),
			}
		}
	}
	return this._extValue
}

func (this *ImageAdImpl) SetData(d *ad.Image) error {
	v := this.getData()
	v.LinkUrl = d.LinkUrl
	v.Title = d.Title
	v.ImageUrl = d.ImageUrl
	v.Enabled = 1 // 图片广告的图片无启用和停用功能
	return nil
}

// 保存广告
func (this *ImageAdImpl) Save() (int, error) {
	id, err := this.AdImpl.Save()
	if err == nil {
		v := this.getData()
		v.AdId = id
		_, err = this._rep.SaveAdImageValue(v)
	}
	return id, err
}

// 转换为数据传输对象
func (this *ImageAdImpl) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:   this.AdImpl.GetDomainId(),
		Type: this.AdImpl.Type(),
		Data: this.getData(),
	}
}
