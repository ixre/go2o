/**
 * Copyright 2015 @ to2.net.
 * name : image_ad
 * author : jarryliu
 * date : 2016-05-25 11:29
 * description :
 * history :
 */
package ad

import (
	"go2o/core/domain/interface/ad"
	"go2o/core/infrastructure/format"
)

var _ ad.IImageAd = new(ImageAdImpl)

type ImageAdImpl struct {
	extValue *ad.Image
	*adImpl
}

// 获取链接广告值
func (i *ImageAdImpl) getData() *ad.Image {
	if i.extValue == nil {
		gallery := i._rep.GetValueGallery(i.GetDomainId())
		if gallery.Len() > 0 {
			i.extValue = gallery[0]
		}

		//如果不存在,则创建一个新的对象
		if i.extValue == nil {
			i.extValue = &ad.Image{
				AdId:     i.GetDomainId(),
				ImageUrl: format.GetNoPicPath(),
				Enabled:  1,
			}
		}
	}
	return i.extValue
}

func (i *ImageAdImpl) SetData(d *ad.Image) error {
	v := i.getData()
	v.LinkUrl = d.LinkUrl
	v.Title = d.Title
	v.ImageUrl = d.ImageUrl
	v.Enabled = 1 // 图片广告的图片无启用和停用功能
	return nil
}

// 保存广告
func (i *ImageAdImpl) Save() (int32, error) {
	id, err := i.adImpl.Save()
	if err == nil {
		v := i.getData()
		v.AdId = id
		_, err = i._rep.SaveAdImageValue(v)
	}
	return id, err
}

// 转换为数据传输对象
func (i *ImageAdImpl) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:   i.adImpl.GetDomainId(),
		Type: i.adImpl.Type(),
		Data: i.getData(),
	}
}
