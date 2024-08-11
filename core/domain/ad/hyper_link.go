/**
 * Copyright 2015 @ 56x.net.
 * name : hyper_link.go
 * author : jarryliu
 * date : 2016-05-25 10:50
 * description :
 * history :
 */
package ad

import (
	"github.com/ixre/go2o/core/domain/interface/ad"
)

var _ ad.IHyperLinkAd = new(HyperLinkAdImpl)

type HyperLinkAdImpl struct {
	extValue *ad.Data
	*adImpl
}

// 获取链接广告值
func (h *HyperLinkAdImpl) getData() *ad.Data {
	if h.extValue == nil {
		h.extValue = h._rep.GetTextAdData(h.GetDomainId())
	}
	return h.extValue
}

func (h *HyperLinkAdImpl) SetData(d *ad.Data) error {
	v := h.getData()
	v.AdId = h.adImpl.GetDomainId()
	v.LinkUrl = d.LinkUrl
	v.Title = d.Title
	return nil
}

// 保存广告
func (h *HyperLinkAdImpl) Save() (int64, error) {
	id, err := h.adImpl.Save()
	if err == nil {
		v := h.getData()
		if v == nil {
			v = &ad.Data{}
		}
		v.AdId = id
		_, err = h._rep.SaveTextAdData(v)
	}
	return id, err
}

// 转换为数据传输对象
func (h *HyperLinkAdImpl) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:     h.adImpl.GetDomainId(),
		AdType: h.adImpl.Type(),
		Data:   h.getData(),
	}
}
