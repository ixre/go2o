/**
 * Copyright 2015 @ to2.net.
 * name : hyper_link.go
 * author : jarryliu
 * date : 2016-05-25 10:50
 * description :
 * history :
 */
package ad

import (
	"go2o/core/domain/interface/ad"
)

var _ ad.IHyperLinkAd = new(HyperLinkAdImpl)

type HyperLinkAdImpl struct {
	extValue *ad.HyperLink
	*adImpl
}

// 获取链接广告值
func (h *HyperLinkAdImpl) getData() *ad.HyperLink {
	if h.extValue == nil {
		h.extValue = h._rep.GetHyperLinkData(h.GetDomainId())

		//如果不存在,则创建一个新的对象
		if h.extValue == nil {
			h.extValue = &ad.HyperLink{
				AdId: h.GetDomainId(),
			}
		}
	}
	return h.extValue
}

func (h *HyperLinkAdImpl) SetData(d *ad.HyperLink) error {
	v := h.getData()
	v.AdId = h.adImpl.GetDomainId()
	v.LinkUrl = d.LinkUrl
	v.Title = d.Title
	return nil
}

// 保存广告
func (h *HyperLinkAdImpl) Save() (int32, error) {
	id, err := h.adImpl.Save()
	if err == nil {
		v := h.getData()
		v.AdId = id
		_, err = h._rep.SaveHyperLinkData(v)
	}
	return id, err
}

// 转换为数据传输对象
func (h *HyperLinkAdImpl) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:   h.adImpl.GetDomainId(),
		Type: h.adImpl.Type(),
		Data: h.getData(),
	}
}
