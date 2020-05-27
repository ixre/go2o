/**
 * Copyright 2015 @ to2.net.
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
	adValue ad.ValueGallery
}

// 获取广告数据
func (g *GalleryAd) GetAdValue() ad.ValueGallery {
	if g.adValue == nil {
		if g.GetDomainId() > 0 {
			g.adValue = g._rep.GetValueGallery(g.GetDomainId())
			sort.Sort(g.adValue)
		} else {
			g.adValue = []*ad.Image{}
		}
	}
	return g.adValue
}

// 获取可用的广告数据
func (g *GalleryAd) GetEnabledAdValue() ad.ValueGallery {
	val := g.GetAdValue()
	newVal := ad.ValueGallery{}
	for _, v := range val {
		if v.Enabled == 1 {
			newVal = append(newVal, v)
		}
	}
	return newVal
}

// 保存广告图片
func (g *GalleryAd) SaveImage(v *ad.Image) (int32, error) {
	v.AdId = g.GetDomainId()
	return g._rep.SaveAdImageValue(v)
}

// 获取图片项
func (g *GalleryAd) GetImage(id int32) *ad.Image {
	return g._rep.GetValueAdImage(g.GetDomainId(), id)
}

// 删除图片项
func (g *GalleryAd) DelImage(id int32) error {
	return g._rep.DelAdImage(g.GetDomainId(), id)
}

// 转换为数据传输对象
func (g *GalleryAd) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:   g.adImpl.GetDomainId(),
		Type: g.adImpl.Type(),
		Data: g.GetEnabledAdValue(),
	}
}
