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
	"com/infrastructure"
	"com/infrastructure/log"
	"com/share/variable"
	"fmt"
	"github.com/atnet/gof/db"
)

var _ partner.IPartnerRep = new(partnerRep)

type partnerRep struct {
	db.Connector
	cache map[int]partner.IPartner
}

func NewPartnerRep(c db.Connector) partner.IPartnerRep {
	return &partnerRep{
		Connector: c,
		cache:     make(map[int]partner.IPartner),
	}
}

func (this *partnerRep) CreatePartner(v *partner.ValuePartner) partner.IPartner {
	return pt.NewPartner(v, this)
}

func (this *partnerRep) renew(partnerId int) {
	delete(this.cache, partnerId)
}

func (this *partnerRep) GetPartner(id int) partner.IPartner {
	v, ok := this.cache[id]
	if !ok {
		e := new(partner.ValuePartner)
		if this.Connector.GetOrm().Get(id, e) == nil {
			v = pt.NewPartner(e, this)
			this.cache[id] = v
		}
	}
	return v
}

// 获取销售配置
func (this *partnerRep) GetSaleConf(partnerId int) *partner.SaleConf {
	//10%分成
	//0.2,         #上级
	//0.1,         #上上级
	//0.8          #消费者自己
	var saleConf *partner.SaleConf = new(partner.SaleConf)
	if this.Connector.GetOrm().Get(partnerId, saleConf) == nil {
		return saleConf
	}
	return nil
}

func (this *partnerRep) SaveSaleConf(v *partner.SaleConf) error {
	defer this.renew(v.PartnerId)
	_, _, err := this.Connector.GetOrm().Save(v.PartnerId, v)
	return err
}

// 获取站点配置
func (this *partnerRep) GetSiteConf(partnerId int) *partner.SiteConf {
	var siteConf partner.SiteConf
	if err := this.Connector.GetOrm().Get(partnerId, &siteConf); err == nil {
		if len(siteConf.Host) == 0 {
			var usr string
			this.Connector.ExecScalar(
				`SELECT usr FROM pt_partner WHERE id=?`,
				&usr, partnerId)
			siteConf.Host = fmt.Sprintf("%s.%s", usr,
				infrastructure.GetContext().Config().
					GetString(variable.ServerDomain))
		}
		return &siteConf
	}
	return nil
}

func (this *partnerRep) SaveSiteConf(v *partner.SiteConf) error {
	defer this.renew(v.PartnerId)
	_, _, err := this.Connector.GetOrm().Save(v.PartnerId, v)
	return err
}

func (this *partnerRep) SaveShop(v *partner.ValueShop) (int, error) {
	defer this.renew(v.PartnerId)
	orm := this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	} else {
		_, _, err := orm.Save(nil, v)

		//todo: return id
		return 0, err
	}
}

func (this *partnerRep) GetValueShop(partnerId, shopId int) *partner.ValueShop {
	var v *partner.ValueShop = new(partner.ValueShop)
	err := this.Connector.GetOrm().Get(shopId, v)
	if err == nil &&
		v.PartnerId == partnerId {
		return v
	} else {
		log.PrintErr(err)
	}
	return nil
}

func (this *partnerRep) GetShopsOfPartner(partnerId int) []*partner.ValueShop {
	shops := []*partner.ValueShop{}
	err := this.Connector.GetOrm().SelectByQuery(&shops,
		"SELECT * FROM pt_shop WHERE pt_id=?", partnerId)

	if err != nil {
		log.PrintErr(err)
		return nil
	}

	return shops
}

func (this *partnerRep) DeleteShop(partnerId, shopId int) error {
	defer this.renew(partnerId)
	_, err := this.Connector.GetOrm().Delete(partner.ValueShop{},
		"pt_id=? AND id=?", partnerId, shopId)
	return err
}
