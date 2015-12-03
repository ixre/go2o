/**
 * Copyright 2015 @ z3q.net.
 * name : ad_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	adImpl "go2o/src/core/domain/ad"
	"go2o/src/core/domain/interface/ad"
)

var _ ad.IAdvertisementRep = new(advertisementRep)

type advertisementRep struct {
	db.Connector
}

// 广告仓储
func NewAdvertisementRep(c db.Connector) ad.IAdvertisementRep {
	return &advertisementRep{
		Connector: c,
	}
}

// 获取商户的广告管理
func (this *advertisementRep) GetPartnerAdvertisement(partnerId int) ad.IPartnerAdvertisement {
	return adImpl.NewPartnerAdvertisement(partnerId, this)
}

// 根据名称获取广告编号
func (this *advertisementRep) GetIdByName(partnerId int, name string) int {
	var id int
	this.Connector.ExecScalar("SELECT id FROM pt_ad WHERE partner_id=? AND name=?", &id, partnerId, name)
	return id
}

// 保存广告值
func (this *advertisementRep) SaveAdvertisementValue(v *ad.ValueAdvertisement) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM pt_ad WHERE partner_id=?", &v.Id, v.PartnerId)
	}
	return v.Id, err
}

// 保存广告图片
func (this *advertisementRep) SaveAdImageValue(v *ad.ValueImage) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM pt_ad_image WHERE ad_id=?", &v.Id, v.AdvertisementId)
	}
	return v.Id, err
}

// 获取广告
func (this *advertisementRep) GetValueAdvertisement(partnerId, id int) *ad.ValueAdvertisement {
	var e ad.ValueAdvertisement
	if err := this.Connector.GetOrm().Get(id, &e); err == nil && e.PartnerId == partnerId {
		return &e
	}
	return nil
}

// 根据名称获取广告
func (this *advertisementRep) GetValueAdvertisementByName(partnerId int, name string) *ad.ValueAdvertisement {
	var e ad.ValueAdvertisement
	if err := this.Connector.GetOrm().GetBy(&e, "partner_id=? and name=?", partnerId, name); err == nil {
		return &e
	}
	return nil
}

// 获取轮播广告
func (this *advertisementRep) GetValueGallery(advertisementId int) ad.ValueGallery {
	var list = []*ad.ValueImage{}
	if err := this.Connector.GetOrm().Select(&list, "ad_id=? ORDER BY sort_number ASC", advertisementId); err == nil {
		return list
	}
	return nil
}

// 获取图片项
func (this *advertisementRep) GetValueAdImage(advertisementId, id int) *ad.ValueImage {
	var e ad.ValueImage
	if err := this.Connector.GetOrm().GetBy(&e, "ad_id=? and id=?", advertisementId, id); err == nil {
		return &e
	}
	return nil
}

// 删除图片项
func (this *advertisementRep) DelAdImage(advertisementId, id int) error {
	_, err := this.Connector.GetOrm().Delete(ad.ValueImage{}, "ad_id=? and id=?", advertisementId, id)
	return err
}

// 删除广告
func (this *advertisementRep) DelAdvertisement(partnerId, advertisementId int) error {
	_, err := this.Connector.GetOrm().Delete(ad.ValueAdvertisement{}, "partner_id=? AND id=?", partnerId, advertisementId)
	return err
}

// 删除广告的图片数据
func (this *advertisementRep) DelImageDataForAdvertisement(advertisementId int) error {
	_, err := this.Connector.GetOrm().Delete(ad.ValueImage{}, "ad_id=?", advertisementId)
	return err
}

// 删除广告的文字数据
func (this *advertisementRep) DelTextDataForAdvertisement(advertisementId int) error {
	//todo:
	return nil
}
