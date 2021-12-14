/**
 * Copyright 2015 @ 56x.net.
 * name : ad_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repos

import (
	"database/sql"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	adImpl "go2o/core/domain/ad"
	"go2o/core/domain/interface/ad"
	"log"
	"sync"
)

var _ ad.IAdRepo = new(advertisementRepo)

type advertisementRepo struct {
	db.Connector
	sync.Mutex
	storage storage.Interface
	o       orm.Orm
}

// 广告仓储
func NewAdvertisementRepo(o orm.Orm, storage storage.Interface) ad.IAdRepo {
	return &advertisementRepo{
		Connector: o.Connector(),
		storage:   storage,
		o:         o,
	}
}

// 获取广告管理器
func (a *advertisementRepo) GetAdManager() ad.IAdvertisementManager {
	return adImpl.NewAdManager(a)
}

// GetGroups 获取广告分组
func (a *advertisementRepo) GetGroups() []string {
	var arr []string
	a.o.Connector().Query("select distinct(group_name) from ad_position", func(rows *sql.Rows) {
		var s = ""
		for rows.Next() {
			rows.Scan(&s)
			if len(s) > 0 {
				arr = append(arr, s)
			}
		}
	})
	return arr
}

func (a *advertisementRepo) GetPosition(id int64) ad.IAdPosition {
	e := ad.Position{}
	if err := a.o.Get(id, &e); err != nil {
		handleError(err)
		return nil
	}
	return a.CreateAdPosition(&e)
}

func (a *advertisementRepo) CreateAdPosition(v *ad.Position) ad.IAdPosition {
	return adImpl.NewAdPosition(a, v)
}

// 根据KEY获取广告位
func (a *advertisementRepo) GetAdPositionByKey(key string) *ad.Position {
	e := ad.Position{}
	if err := a.o.GetBy(&e, "ad_position.key=$1", key); err != nil {
		handleError(err)
		return nil
	}
	return &e
}

// 根据ID获取广告位
func (a *advertisementRepo) GetAdPositionById(adPosId int64) *ad.Position {
	e := ad.Position{}
	if err := a.o.Get(adPosId, &e); err != nil {
		handleError(err)
		return nil
	}
	return &e
}

// 获取广告位
func (a *advertisementRepo) GetAdPositionsByGroupId(adGroupId int64) []*ad.Position {
	var list []*ad.Position
	if err := a.o.Select(&list, "group_id=$1", adGroupId); err != nil {
		handleError(err)
	}
	return list
}

// 删除广告位
func (a *advertisementRepo) DeleteAdPosition(id int64) error {
	err := a.o.DeleteByPk(&ad.Position{}, id)
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", 0))
	}
	return err
}

// 保存广告位
func (a *advertisementRepo) SaveAdPosition(v *ad.Position) (int64, error) {
	id, err := orm.I64(orm.Save(a.o, v, int(v.Id)))
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", 0))
	}
	return id, err
}

// 设置用户的广告
func (a *advertisementRepo) SetUserAd(adUserId, posId, adId int64) error {
	v := &ad.AdUserSet{
		AdUserId: adUserId,
		PosId:    posId,
		AdId:     adId,
	}
	a.ExecScalar("SELECT id FROM ad_userset WHERE user_id=$1 AND ad_id=$2", &v.Id, adUserId, adId)
	v.PosId = posId
	_, err := orm.Save(a.o, v, int(v.Id))
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", adUserId))
	}
	return err
}

func (a *advertisementRepo) QueryAdList(keyword string, size int) []*ad.Ad {
	var arr = make([]*ad.Ad, 0)
	err := a.o.Select(&arr, " name LIKE '%"+keyword+"%' LIMIT $1", size)
	if err != nil {
		log.Println("QueryAdList error:", err.Error())
	}
	return arr
}

// 根据名称获取广告编号
func (a *advertisementRepo) GetIdByName(userId int64, name string) int {
	var id int
	a.Connector.ExecScalar("SELECT id FROM ad_list WHERE user_id=$1 AND name=$1",
		&id, userId, name)
	return id
}

// 保存广告值
func (a *advertisementRepo) SaveAdValue(v *ad.Ad) (int64, error) {
	id, err := orm.I64(orm.Save(a.o, v, int(v.Id)))
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", v.UserId))
	}
	return id, err
}

// 获取超链接广告数据
func (a *advertisementRepo) GetTextAdData(adId int64) *ad.HyperLink {
	e := ad.HyperLink{}
	if err := a.o.GetBy(&e, "ad_id=$1", adId); err != nil {
		handleError(err)
		return nil
	}
	return &e
}

// 保存超链接广告数据
func (a *advertisementRepo) SaveTextAdData(v *ad.HyperLink) (int64, error) {
	return orm.I64(orm.Save(a.o, v, int(v.Id)))
}

// 保存广告图片
func (a *advertisementRepo) SaveImageAdData(v *ad.Image) (int64, error) {
	return orm.I64(orm.Save(a.o, v, int(v.Id)))
}

// 获取广告
func (a *advertisementRepo) GetAd(id int64) *ad.Ad {
	var e ad.Ad
	if err := a.o.Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 根据名称获取广告
func (a *advertisementRepo) GetAdByKey(userId int64, key string) *ad.Ad {
	e := ad.Ad{}
	const sql string = `select * FROM ad_list
        INNER JOIN ad_userset ON ad_userset.user_id = ad_list.user_id
        INNER JOIN ad_position ON ad_userset.pos_id=ad_position.id
        WHERE ad_list.user_id = $1 AND ad_position.key=$2`
	if err := a.o.GetByQuery(&e, sql, userId, key); err == nil {
		return &e
	}
	return nil
}

// 获取轮播广告
func (a *advertisementRepo) GetSwiperAd(adId int64) ad.SwiperAd {
	var list = []*ad.Image{}
	if err := a.o.Select(&list, "ad_id=$1 ORDER BY sort_num ASC LIMIT 20", adId); err == nil {
		return list
	}
	return nil
}

// 获取图片项
func (a *advertisementRepo) GetSwiperAdImage(adId, id int64) *ad.Image {
	var e ad.Image
	if err := a.o.GetBy(&e, "ad_id=$1 and id=$2", adId, id); err == nil {
		return &e
	}
	return nil
}

// 删除图片项
func (a *advertisementRepo) DeleteSwiperAdImage(adId, imgId int64) error {
	_, err := a.o.Delete(ad.Image{}, "ad_id=$1 and id=$2", adId, imgId)
	return err
}

// 删除广告
func (a *advertisementRepo) DeleteAd(userId, adId int64) error {
	_, err := a.o.Delete(ad.Ad{}, "user_id=$1 AND id=$1", userId, adId)
	if err == nil {
		//更新用户的广告缓存
		PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", userId))
	}
	return err
}

// 删除广告的图片数据
func (a *advertisementRepo) DeleteImageAdData(adId int64) error {
	_, err := a.o.Delete(ad.Image{}, "ad_id=$1", adId)
	return err
}

// 删除广告的文字数据
func (a *advertisementRepo) DeleteTextAdData(adId int64) error {
	_, err := a.o.Delete(ad.HyperLink{}, "ad_id=$1", adId)
	return err
}

func (a *advertisementRepo) GetPositions() []*ad.Position {
	list := make([]*ad.Position, 0)
	err := a.o.Select(&list, "")
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:Position")
	}
	return list
}
