/**
 * Copyright 2015 @ 56x.net.
 * name : gallery
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import (
	"sort"

	"github.com/ixre/go2o/core/domain/interface/ad"
)

var _ ad.IAdAggregateRoot = new(GalleryAd)
var _ ad.IGalleryAd = new(GalleryAd)

type GalleryAd struct {
	*adImpl
	adValue ad.SwiperAd
}

// 获取广告数据
func (g *GalleryAd) GetAdValue() ad.SwiperAd {
	if g.adValue == nil {
		if g.GetDomainId() > 0 {
			g.adValue = g._rep.GetSwiperAd(g.GetDomainId())
			sort.Sort(g.adValue)
		} else {
			g.adValue = []*ad.Image{}
		}
	}
	return g.adValue
}

// 获取可用的广告数据
func (g *GalleryAd) GetEnabledAdValue() ad.SwiperAd {
	val := g.GetAdValue()
	newVal := ad.SwiperAd{}
	for _, v := range val {
		if v.Enabled == 1 {
			newVal = append(newVal, v)
		}
	}
	return newVal
}

// 保存广告图片
func (g *GalleryAd) SaveImage(v *ad.Image) (int64, error) {
	v.AdId = g.GetDomainId()
	return g._rep.SaveImageAdData(v)
}

// 获取图片项
func (g *GalleryAd) GetImage(id int64) *ad.Image {
	return g._rep.GetSwiperAdImage(g.GetDomainId(), id)
}

// 删除图片项
func (g *GalleryAd) DeleteItem(id int64) error {
	return g._rep.DeleteSwiperAdImage(g.GetDomainId(), id)
}

// 转换为数据传输对象
func (g *GalleryAd) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:     g.adImpl.GetDomainId(),
		AdType: g.adImpl.Type(),
		Data:   g.GetEnabledAdValue(),
	}
}
