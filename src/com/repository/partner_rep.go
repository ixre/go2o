/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-12 17:16
 * description :
 * history :
 */

package repository

import (
	"com/domain/interface/partner"
	pt "com/domain/partner"
	"ops/cf/db"
	"fmt"
	"com/share/variable"
	"com/infrastructure"
)

type PartnerRep struct {
	db.Connector
}
func (this *PartnerRep) CreatePartner(v *partner.ValuePartner) partner.IPartner {
	return pt.NewPartner(v, this)
}

func (this *PartnerRep) GetPartner(id int) partner.IPartner {
	e := new(partner.ValuePartner)
	if this.Connector.GetOrm().Get(e, id) == nil {
		return pt.NewPartner(e, this)
	}
	return nil
}

// 获取销售配置
func (this *PartnerRep) GetSaleConf(partnerId int) *partner.SaleConf {
	//10%分成
	//0.2,         #上级
	//0.1,         #上上级
	//0.8          #消费者自己
	var saleConf *partner.SaleConf = new(partner.SaleConf)
	if this.Connector.GetOrm().Get(saleConf, partnerId) == nil {
		return saleConf
	}
	return nil
}


func (this *PartnerRep) SaveSaleConf(v *partner.SaleConf)error{
	_,_,err := this.Connector.GetOrm().Save(v.PtId,v)
	return err
}

// 获取站点配置
func (this *PartnerRep) GetSiteConf(partnerId int) *partner.SiteConf{
	var siteConf partner.SiteConf
	if err := this.Connector.GetOrm().Get(&siteConf, partnerId); err == nil {
		if len(siteConf.Host) == 0 {
			var usr string
			this.Connector.ExecScalar(
				`SELECT usr FROM pt_partner WHERE id=?`,
				&usr, partnerId)
			siteConf.Host = fmt.Sprintf("%s.%s", usr,
				infrastructure.GetContext().Config().
				Get(variable.ServerDomain))
		}
		return &siteConf
	}
	return nil
}

func (this *PartnerRep) SaveSiteConf(v *partner.SiteConf)error{
	_,_,err := this.Connector.GetOrm().Save(v.PtId,v)
	return err
}
