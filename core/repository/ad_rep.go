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
    "github.com/jsix/gof/db/orm"
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
func (a *advertisementRep) GetAdManager() ad.IAdManager {
    return adImpl.NewAdManager(a)
}

// 获取广告分组
func (a *advertisementRep) GetAdGroups() []*ad.AdGroup {
    var list = []*ad.AdGroup{}
    if err := a.Connector.GetOrm().Select(&list, ""); err != nil {
        handleError(err)
    }
    return list
}

// 删除广告组
func (a *advertisementRep) DelAdGroup(id int) error {
    return a.Connector.GetOrm().DeleteByPk(&ad.AdGroup{}, id)
}

// 根据KEY获取广告位
func (a *advertisementRep) GetAdPositionByKey(key string) *ad.AdPosition {
    e := ad.AdPosition{}
    if err := a.GetOrm().GetBy(&e, "ad_position.key=?", key); err != nil {
        handleError(err)
        return nil
    }
    return &e
}

// 获取广告位
func (a *advertisementRep) GetAdPositionsByGroupId(adGroupId int) []*ad.AdPosition {
    var list = []*ad.AdPosition{}
    if err := a.Connector.GetOrm().Select(&list, "group_id=?", adGroupId); err != nil {
        handleError(err)
    }
    return list
}

// 删除广告位
func (a *advertisementRep) DelAdPosition(id int) error {
    return a.Connector.GetOrm().DeleteByPk(&ad.AdPosition{}, id)
}

// 保存广告位
func (a *advertisementRep) SaveAdPosition(v *ad.AdPosition) (int, error) {
    var err error
    a.Mutex.Lock()
    var orm = a.Connector.GetOrm()
    if v.Id > 0 {
        _, _, err = orm.Save(v.Id, v)
    } else {
        _, _, err = orm.Save(nil, v)
        a.Connector.ExecScalar("SELECT MAX(id) FROM ad_position WHERE group_id=?", &v.Id, v.GroupId)
    }
    a.Mutex.Unlock()
    return v.Id, err
}

// 保存
func (a *advertisementRep) SaveAdGroup(v *ad.AdGroup) (int, error) {
    var err error
    a.Mutex.Lock()
    var orm = a.Connector.GetOrm()
    if v.Id > 0 {
        _, _, err = orm.Save(v.Id, v)
    } else {
        _, _, err = orm.Save(nil, v)
        a.Connector.ExecScalar("SELECT MAX(id) FROM ad_group", &v.Id)
    }
    a.Mutex.Unlock()
    return v.Id, err
}

// 设置用户的广告
func (a *advertisementRep) SetUserAd(adUserId, posId, adId int) error {
    e := &ad.AdUserSet{
        AdUserId: adUserId,
        PosId:    posId,
        AdId:     adId,
    }
    a.ExecScalar("SELECT id FROM ad_userset WHERE user_id=? AND ad_id=?", &e.Id, adUserId, adId)
    e.PosId = posId
    _, err := orm.Save(a.GetOrm(), e, e.Id)
    return err
}

// 根据名称获取广告编号
func (a *advertisementRep) GetIdByName(userId int, name string) int {
    var id int
    a.Connector.ExecScalar("SELECT id FROM ad_list WHERE user_id=? AND name=?",
        &id, userId, name)
    return id
}

// 保存广告值
func (a *advertisementRep) SaveAdValue(v *ad.Ad) (int, error) {
    var err error
    a.Mutex.Lock()
    var orm = a.Connector.GetOrm()
    if v.Id > 0 {
        _, _, err = orm.Save(v.Id, v)
    } else {
        _, _, err = orm.Save(nil, v)
        a.Connector.ExecScalar("SELECT MAX(id) FROM ad_list WHERE user_id=?",
            &v.Id, v.UserId)
    }
    a.Mutex.Unlock()
    return v.Id, err
}

// 获取超链接广告数据
func (a *advertisementRep) GetHyperLinkData(adId int) *ad.HyperLink {
    e := ad.HyperLink{}
    if err := a.GetOrm().GetBy(&e, "ad_id=?", adId); err != nil {
        handleError(err)
        return nil
    }
    return &e
}

// 保存超链接广告数据
func (a *advertisementRep) SaveHyperLinkData(v *ad.HyperLink) (int, error) {
    var err error
    a.Mutex.Lock()
    var orm = a.Connector.GetOrm()
    if v.Id > 0 {
        _, _, err = orm.Save(v.Id, v)
    } else {
        _, _, err = orm.Save(nil, v)
        a.Connector.ExecScalar("SELECT MAX(id) FROM ad_hyperlink WHERE ad_id=?", &v.Id, v.AdId)
    }
    a.Mutex.Unlock()
    return v.Id, err
}

// 保存广告图片
func (a *advertisementRep) SaveAdImageValue(v *ad.Image) (int, error) {
    var err error
    var orm = a.Connector.GetOrm()
    if v.Id > 0 {
        _, _, err = orm.Save(v.Id, v)
    } else {
        _, _, err = orm.Save(nil, v)
        a.Connector.ExecScalar("SELECT MAX(id) FROM ad_image WHERE ad_id=?", &v.Id, v.AdId)
    }
    return v.Id, err
}

// 获取广告
func (a *advertisementRep) GetValueAd(id int) *ad.Ad {
    var e ad.Ad
    if err := a.Connector.GetOrm().Get(id, &e); err == nil {
        return &e
    }
    return nil
}

// 根据名称获取广告
func (a *advertisementRep) GetAdByKey(userId int, key string) *ad.Ad {
    e := ad.Ad{}
    const sql string = `select * FROM ad_list
        INNER JOIN ad_userset ON ad_userset.user_id = ad_list.user_id
        INNER JOIN ad_position ON ad_userset.pos_id=ad_position.id
        WHERE ad_list.user_id = ? AND ad_position.key=?`
    if err := a.Connector.GetOrm().GetByQuery(&e, sql, userId, key); err == nil {
        return &e
    }
    return nil
}

// 获取轮播广告
func (a *advertisementRep) GetValueGallery(advertisementId int) ad.ValueGallery {
    var list = []*ad.Image{}
    if err := a.Connector.GetOrm().Select(&list, "ad_id=? ORDER BY sort_number ASC", advertisementId); err == nil {
        return list
    }
    return nil
}

// 获取图片项
func (a *advertisementRep) GetValueAdImage(advertisementId, id int) *ad.Image {
    var e ad.Image
    if err := a.Connector.GetOrm().GetBy(&e, "ad_id=? and id=?", advertisementId, id); err == nil {
        return &e
    }
    return nil
}

// 删除图片项
func (a *advertisementRep) DelAdImage(advertisementId, id int) error {
    _, err := a.Connector.GetOrm().Delete(ad.Image{}, "ad_id=? and id=?", advertisementId, id)
    return err
}

// 删除广告
func (a *advertisementRep) DelAd(userId, advertisementId int) error {
    _, err := a.Connector.GetOrm().Delete(ad.Ad{}, "user_id=? AND id=?", userId, advertisementId)
    return err
}

// 删除广告的图片数据
func (a *advertisementRep) DelImageDataForAdvertisement(advertisementId int) error {
    _, err := a.Connector.GetOrm().Delete(ad.Image{}, "ad_id=?", advertisementId)
    return err
}

// 删除广告的文字数据
func (a *advertisementRep) DelTextDataForAdvertisement(advertisementId int) error {
    _, err := a.Connector.GetOrm().Delete(ad.HyperLink{}, "ad_id=?", advertisementId)
    return err
}
