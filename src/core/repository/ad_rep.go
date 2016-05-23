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

var _ ad.IAdRep = new(advertisementRep)

type advertisementRep struct {
	db.Connector
}

// 广告仓储
func NewAdvertisementRep(c db.Connector) ad.IAdRep {
	return &advertisementRep{
		Connector: c,
	}
}

// 获取广告管理器
func (this *advertisementRep) GetAdManager() ad.IAdManager {
	return adImpl.NewAdManager(this)
}

// 获取广告分组
func (this *advertisementRep) GetAdGroups() []*ad.AdGroup {
	panic("")
}

// 删除广告组
func (this *advertisementRep) DelAdGroup(id int) error {
	panic("")
}

// 获取广告位
func (this *advertisementRep) GetAdPositionsByGroupId(adGroupId int) []*ad.AdPosition {
	panic("")
}

// 删除广告位
func (this *advertisementRep) DelAdPosition(id int) error {
	panic("")
}

// 保存广告位
func (this *advertisementRep) SaveAdPosition(a *ad.AdPosition) (int, error) {
	panic("")
}

// 保存
func (this *advertisementRep) SaveAdGroup(value *ad.AdGroup) (int, error) {
	panic("")
}

// 设置用户的广告
func (this *advertisementRep) SetUserAd(adUserId, posId, adId int) error {
	panic("")
}

// 根据名称获取广告编号
func (this *advertisementRep) GetIdByName(merchantId int, name string) int {
	var id int
	this.Connector.ExecScalar("SELECT id FROM ad_list WHERE merchant_id=? AND name=?", &id, merchantId, name)
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
		this.Connector.ExecScalar("SELECT MAX(id) FROM ad_list WHERE merchant_id=?", &v.Id, v.AdUserId)
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
		this.Connector.ExecScalar("SELECT MAX(id) FROM ad_image_ad WHERE ad_id=?", &v.Id, v.AdvertisementId)
	}
	return v.Id, err
}

// 获取广告
func (this *advertisementRep) GetValueAdvertisement(id int) *ad.ValueAdvertisement {
	var e ad.ValueAdvertisement
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 根据名称获取广告
func (this *advertisementRep) GetValueAdvertisementByName(merchantId int, name string) *ad.ValueAdvertisement {
	var e ad.ValueAdvertisement
	if err := this.Connector.GetOrm().GetBy(&e, "merchant_id=? and name=?", merchantId, name); err == nil {
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
func (this *advertisementRep) DelAdvertisement(merchantId, advertisementId int) error {
	_, err := this.Connector.GetOrm().Delete(ad.ValueAdvertisement{}, "merchant_id=? AND id=?", merchantId, advertisementId)
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
