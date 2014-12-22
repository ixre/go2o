/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-12 16:55
 * description :
 * history :
 */

package partner

import (
	"com/domain/interface/partner"
)

type Partner struct {
	value    *partner.ValuePartner
	saleConf *partner.SaleConf
	siteConf *partner.SiteConf
	rep      partner.IPartnerRep
}

func NewPartner(v *partner.ValuePartner, rep partner.IPartnerRep) partner.IPartner {
	return &Partner{
		value: v,
		rep:   rep,
	}
}

func (this *Partner) GetAggregateRootId() int {
	return this.value.Id
}
func (this *Partner) GetValue() partner.ValuePartner {
	return *this.value
}

// 获取销售配置
func (this *Partner) GetSaleConf() partner.SaleConf {
	if this.saleConf == nil {

		//10%分成
		//0.2,         #上级
		//0.1,         #上上级
		//0.8          #消费者自己

		this.saleConf = this.rep.GetSaleConf(
			this.GetAggregateRootId())
	}
	return *this.saleConf
}



// 保存销售配置
func (this *Partner) SaveSaleConf(v *partner.SaleConf)error{
	this.saleConf = v
	this.saleConf.PtId= this.value.Id
	return this.rep.SaveSaleConf(this.saleConf)
}

// 获取站点配置
func (this *Partner) GetSiteConf()partner.SiteConf{
	if this.siteConf == nil {
		this.siteConf = this.rep.GetSiteConf(this.GetAggregateRootId())
	}
	return *this.siteConf
}

// 保存站点配置
func (this *Partner) SaveSiteConf(v *partner.SiteConf)error {
	this.siteConf = v
	this.saleConf.PtId = this.value.Id
	return this.rep.SaveSiteConf(this.siteConf)
}
