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
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
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
			g.adValue = []*ad.Data{}
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
func (g *GalleryAd) SaveImage(items []*ad.Data) error {
	nowMap := collections.ToMap(items, func(a *ad.Data) (int64, *ad.Data) {
		return a.Id, a
	})
	// 从旧的数据中筛选出要删除的项
	delList := collections.FilterArray(g.GetAdValue(), func(v *ad.Data) bool {
		return nowMap[v.Id] == nil
	})
	// 删除项
	for _, v := range delList {
		g._rep.DeleteSwiperAdImage(g.GetDomainId(), v.Id)
	}
	// 保存项
	for _, v := range items {
		if v.AdId == 0 {
			v.AdId = g.GetDomainId()
		}
		_, err := g._rep.SaveImageAdData(v)
		if err != nil {
			return err
		}
	}
	g.adValue = nil
	return nil
}

// 获取图片项
func (g *GalleryAd) GetImage(id int64) *ad.Data {
	return g._rep.GetSwiperAdImage(g.GetDomainId(), id)
}

// 转换为数据传输对象
func (g *GalleryAd) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:     g.adImpl.GetDomainId(),
		AdType: g.adImpl.Type(),
		Data:   g.GetEnabledAdValue(),
	}
}
