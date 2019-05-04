/**
 * Copyright 2015 @ z3q.net.
 * name : ad_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repos

import (
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	adImpl "go2o/core/domain/ad"
	"go2o/core/domain/interface/ad"
	"sync"
)

var _ ad.IAdRepo = new(advertisementRepo)

type advertisementRepo struct {
	db.Connector
	sync.Mutex
	storage storage.Interface
}

// 广告仓储
func NewAdvertisementRepo(c db.Connector, storage storage.Interface) ad.IAdRepo {
	return &advertisementRepo{
		Connector: c,
		storage:   storage,
	}
}

// 获取广告管理器
func (a *advertisementRepo) GetAdManager() ad.IAdManager {
	return adImpl.NewAdManager(a)
}

// 获取广告分组
func (a *advertisementRepo) GetAdGroups() []*ad.AdGroup {
	var list []*ad.AdGroup
	if err := a.Connector.GetOrm().Select(&list, ""); err != nil {
		handleError(err)
	}
	return list
}

// 删除广告组
func (a *advertisementRepo) DelAdGroup(id int32) error {
	return a.Connector.GetOrm().DeleteByPk(&ad.AdGroup{}, id)
}

// 根据KEY获取广告位
func (a *advertisementRepo) GetAdPositionByKey(key string) *ad.AdPosition {
	e := ad.AdPosition{}
	if err := a.GetOrm().GetBy(&e, "ad_position.key=$1", key); err != nil {
		handleError(err)
		return nil
	}
	return &e
}

// 根据ID获取广告位
func (a *advertisementRepo) GetAdPositionById(adPosId int32) *ad.AdPosition {
	e := ad.AdPosition{}
	if err := a.GetOrm().Get(adPosId, &e); err != nil {
		handleError(err)
		return nil
	}
	return &e
}

// 获取广告位
func (a *advertisementRepo) GetAdPositionsByGroupId(adGroupId int32) []*ad.AdPosition {
	var list []*ad.AdPosition
	if err := a.Connector.GetOrm().Select(&list, "group_id=$1", adGroupId); err != nil {
		handleError(err)
	}
	return list
}

// 删除广告位
func (a *advertisementRepo) DelAdPosition(id int32) error {
	err := a.Connector.GetOrm().DeleteByPk(&ad.AdPosition{}, id)
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", 0))
	}
	return err
}

// 保存广告位
func (a *advertisementRepo) SaveAdPosition(v *ad.AdPosition) (int32, error) {
	id, err := orm.I32(orm.Save(a.GetOrm(), v, int(v.ID)))
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", 0))
	}
	return id, err
}

// 保存
func (a *advertisementRepo) SaveAdGroup(v *ad.AdGroup) (int32, error) {
	return orm.I32(orm.Save(a.GetOrm(), v, int(v.ID)))
}

// 设置用户的广告
func (a *advertisementRepo) SetUserAd(adUserId, posId, adId int32) error {
	v := &ad.AdUserSet{
		AdUserId: adUserId,
		PosId:    posId,
		AdId:     adId,
	}
	a.ExecScalar("SELECT id FROM ad_userset WHERE user_id=$1 AND ad_id=$2", &v.Id, adUserId, adId)
	v.PosId = posId
	_, err := orm.Save(a.GetOrm(), v, int(v.Id))
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", adUserId))
	}
	return err
}

// 根据名称获取广告编号
func (a *advertisementRepo) GetIdByName(userId int32, name string) int {
	var id int
	a.Connector.ExecScalar("SELECT id FROM ad_list WHERE user_id=$1 AND name=$1",
		&id, userId, name)
	return id
}

// 保存广告值
func (a *advertisementRepo) SaveAdValue(v *ad.Ad) (int32, error) {
	id, err := orm.I32(orm.Save(a.GetOrm(), v, int(v.Id)))
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", v.UserId))
	}
	return id, err
}

// 获取超链接广告数据
func (a *advertisementRepo) GetHyperLinkData(adId int32) *ad.HyperLink {
	e := ad.HyperLink{}
	if err := a.GetOrm().GetBy(&e, "ad_id=$1", adId); err != nil {
		handleError(err)
		return nil
	}
	return &e
}

// 保存超链接广告数据
func (a *advertisementRepo) SaveHyperLinkData(v *ad.HyperLink) (int32, error) {
	return orm.I32(orm.Save(a.GetOrm(), v, int(v.Id)))
}

// 保存广告图片
func (a *advertisementRepo) SaveAdImageValue(v *ad.Image) (int32, error) {
	return orm.I32(orm.Save(a.GetOrm(), v, int(v.Id)))
}

// 获取广告
func (a *advertisementRepo) GetValueAd(id int32) *ad.Ad {
	var e ad.Ad
	if err := a.Connector.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 根据名称获取广告
func (a *advertisementRepo) GetAdByKey(userId int32, key string) *ad.Ad {
	e := ad.Ad{}
	const sql string = `select * FROM ad_list
        INNER JOIN ad_userset ON ad_userset.user_id = ad_list.user_id
        INNER JOIN ad_position ON ad_userset.pos_id=ad_position.id
        WHERE ad_list.user_id = $1 AND ad_position.key=$2`
	if err := a.Connector.GetOrm().GetByQuery(&e, sql, userId, key); err == nil {
		return &e
	}
	return nil
}

// 获取轮播广告
func (a *advertisementRepo) GetValueGallery(adId int32) ad.ValueGallery {
	var list = []*ad.Image{}
	if err := a.Connector.GetOrm().Select(&list, "ad_id=$1 ORDER BY sort_num ASC LIMIT 20", adId); err == nil {
		return list
	}
	return nil
}

// 获取图片项
func (a *advertisementRepo) GetValueAdImage(adId, id int32) *ad.Image {
	var e ad.Image
	if err := a.Connector.GetOrm().GetBy(&e, "ad_id=$1 and id=$2", adId, id); err == nil {
		return &e
	}
	return nil
}

// 删除图片项
func (a *advertisementRepo) DelAdImage(adId, imgId int32) error {
	_, err := a.Connector.GetOrm().Delete(ad.Image{}, "ad_id=$1 and id=$2", adId, imgId)
	return err
}

// 删除广告
func (a *advertisementRepo) DelAd(userId, adId int32) error {
	_, err := a.Connector.GetOrm().Delete(ad.Ad{}, "user_id=$1 AND id=$1", userId, adId)
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", userId))
	}
	return err
}

// 删除广告的图片数据
func (a *advertisementRepo) DelImageDataForAdvertisement(adId int32) error {
	_, err := a.Connector.GetOrm().Delete(ad.Image{}, "ad_id=$1", adId)
	return err
}

// 删除广告的文字数据
func (a *advertisementRepo) DelTextDataForAdvertisement(adId int32) error {
	_, err := a.Connector.GetOrm().Delete(ad.HyperLink{}, "ad_id=$1", adId)
	return err
}
