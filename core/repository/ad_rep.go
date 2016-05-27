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
	adImpl "go2o/core/domain/ad"
	"go2o/core/domain/interface/ad"
	"sync"
)

var _ ad.IAdRep = new(advertisementRep)

type advertisementRep struct {
	db.Connector
	sync.Mutex
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
	var list = []*ad.AdGroup{}
	if err := this.Connector.GetOrm().Select(&list, ""); err != nil {
		handleError(err)
	}
	return list
}

// 删除广告组
func (this *advertisementRep) DelAdGroup(id int) error {
	return this.Connector.GetOrm().DeleteByPk(&ad.AdGroup{}, id)
}

// 获取广告位
func (this *advertisementRep) GetAdPositionsByGroupId(adGroupId int) []*ad.AdPosition {
	var list = []*ad.AdPosition{}
	if err := this.Connector.GetOrm().Select(&list, "group_id=?", adGroupId); err != nil {
		handleError(err)
	}
	return list
}

// 删除广告位
func (this *advertisementRep) DelAdPosition(id int) error {
	return this.Connector.GetOrm().DeleteByPk(&ad.AdPosition{}, id)
}

// 保存广告位
func (this *advertisementRep) SaveAdPosition(v *ad.AdPosition) (int, error) {
	var err error
	this.Mutex.Lock()
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM ad_position WHERE group_id=?", &v.Id, v.GroupId)
	}
	this.Mutex.Unlock()
	return v.Id, err
}

// 保存
func (this *advertisementRep) SaveAdGroup(v *ad.AdGroup) (int, error) {
	var err error
	this.Mutex.Lock()
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM ad_group", &v.Id)
	}
	this.Mutex.Unlock()
	return v.Id, err
}

// 设置用户的广告
func (this *advertisementRep) SetUserAd(adUserId, posId, adId int) error {
	_, err := this.Connector.ExecNonQuery("UPDATE ad_userset set ad_id=? WHERE user_id=? AND pos_id=?",
		adId, adUserId, posId)
	return err
}

// 根据名称获取广告编号
func (this *advertisementRep) GetIdByName(merchantId int, name string) int {
	var id int
	this.Connector.ExecScalar("SELECT id FROM ad_list WHERE merchant_id=? AND name=?",
		&id, merchantId, name)
	return id
}

// 保存广告值
func (this *advertisementRep) SaveAdValue(v *ad.Ad) (int, error) {
	var err error
	this.Mutex.Lock()
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM ad_list WHERE merchant_id=?",
			&v.Id, v.UserId)
	}
	this.Mutex.Unlock()
	return v.Id, err
}

// 获取超链接广告数据
func (this *advertisementRep) GetHyperLinkData(adId int) *ad.HyperLink {
	e := ad.HyperLink{}
	if err := this.GetOrm().GetBy(&e, "ad_id=?", adId); err != nil {
		handleError(err)
		return nil
	}
	return &e
}

// 保存超链接广告数据
func (this *advertisementRep) SaveHyperLinkData(v *ad.HyperLink) (int, error) {
	var err error
	this.Mutex.Lock()
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM ad_hyperlink WHERE ad_id=?", &v.Id, v.AdId)
	}
	this.Mutex.Unlock()
	return v.Id, err
}

// 保存广告图片
func (this *advertisementRep) SaveAdImageValue(v *ad.Image) (int, error) {
	var err error
	var orm = this.Connector.GetOrm()
	if v.Id > 0 {
		_, _, err = orm.Save(v.Id, v)
	} else {
		_, _, err = orm.Save(nil, v)
		this.Connector.ExecScalar("SELECT MAX(id) FROM ad_image WHERE ad_id=?", &v.Id, v.AdId)
	}
	return v.Id, err
}

// 获取广告
func (this *advertisementRep) GetValueAd(id int) *ad.Ad {
	var e ad.Ad
	if err := this.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 根据名称获取广告
func (this *advertisementRep) GetAdByName(merchantId int, name string) *ad.Ad {
	var e ad.Ad
	if err := this.Connector.GetOrm().GetBy(&e, "merchant_id=? and name=?", merchantId, name); err == nil {
		return &e
	}
	return nil
}

// 获取轮播广告
func (this *advertisementRep) GetValueGallery(advertisementId int) ad.ValueGallery {
	var list = []*ad.Image{}
	if err := this.Connector.GetOrm().Select(&list, "ad_id=? ORDER BY sort_number ASC", advertisementId); err == nil {
		return list
	}
	return nil
}

// 获取图片项
func (this *advertisementRep) GetValueAdImage(advertisementId, id int) *ad.Image {
	var e ad.Image
	if err := this.Connector.GetOrm().GetBy(&e, "ad_id=? and id=?", advertisementId, id); err == nil {
		return &e
	}
	return nil
}

// 删除图片项
func (this *advertisementRep) DelAdImage(advertisementId, id int) error {
	_, err := this.Connector.GetOrm().Delete(ad.Image{}, "ad_id=? and id=?", advertisementId, id)
	return err
}

// 删除广告
func (this *advertisementRep) DelAd(merchantId, advertisementId int) error {
	_, err := this.Connector.GetOrm().Delete(ad.Ad{}, "user_id=? AND id=?", merchantId, advertisementId)
	return err
}

// 删除广告的图片数据
func (this *advertisementRep) DelImageDataForAdvertisement(advertisementId int) error {
	_, err := this.Connector.GetOrm().Delete(ad.Image{}, "ad_id=?", advertisementId)
	return err
}

// 删除广告的文字数据
func (this *advertisementRep) DelTextDataForAdvertisement(advertisementId int) error {
	_, err := this.Connector.GetOrm().Delete(ad.HyperLink{}, "ad_id=?", advertisementId)
	return err
}
