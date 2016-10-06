/**
 * Copyright 2015 @ z3q.net.
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
	_extValue *ad.HyperLink
	*adImpl
}

// 获取链接广告值
func (this *HyperLinkAdImpl) getData() *ad.HyperLink {
	if this._extValue == nil {
		this._extValue = this._rep.GetHyperLinkData(this.GetDomainId())

		//如果不存在,则创建一个新的对象
		if this._extValue == nil {
			this._extValue = &ad.HyperLink{
				AdId: this.GetDomainId(),
			}
		}
	}
	return this._extValue
}

func (this *HyperLinkAdImpl) SetData(d *ad.HyperLink) error {
	v := this.getData()
	v.AdId = this.adImpl.GetDomainId()
	v.LinkUrl = d.LinkUrl
	v.Title = d.Title
	return nil
}

// 保存广告
func (this *HyperLinkAdImpl) Save() (int, error) {
	id, err := this.adImpl.Save()
	if err == nil {
		v := this.getData()
		v.AdId = id
		_, err = this._rep.SaveHyperLinkData(v)
	}
	return id, err
}

// 转换为数据传输对象
func (this *HyperLinkAdImpl) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:   this.adImpl.GetDomainId(),
		Type: this.adImpl.Type(),
		Data: this.getData(),
	}
}
