/**
 * Copyright 2015 @ 56x.net.
 * name : partner_ad
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import (
	"sync"
	"time"

	"github.com/ixre/go2o/core/domain/interface/ad"
	"github.com/ixre/go2o/core/initial/provide"
)

var _ ad.IUserAd = new(userAdImpl)
var _ ad.IAdvertisementManager = new(adManagerImpl)

type adManagerImpl struct {
	rep       ad.IAdRepo
	defaultAd ad.IUserAd
	mux       sync.Mutex
	cache     map[string]ad.IAdAggregateRoot
}

// GetGroups 获取广告分组
func (a *adManagerImpl) GetGroups() []string {
	return a.rep.GetGroups()
}

func NewAdManager(rep ad.IAdRepo) ad.IAdvertisementManager {
	a := &adManagerImpl{
		rep: rep,
	}
	a.defaultAd = newUserAd(a, rep, 0)
	return a
}

func (a *adManagerImpl) QueryAd(keyword string, size int) []*ad.Ad {
	return a.rep.QueryAdList(keyword, size)
}

// 根据编号获取广告位
func (a *adManagerImpl) GetPosition(id int64) ad.IAdPosition {
	return a.rep.GetPosition(id)
}

// 根据KEY获取广告位
func (a *adManagerImpl) GetPositionByKey(key string) *ad.Position {
	return a.rep.GetAdPositionByKey(key)
}

// 根据广告位KEY获取默认广告
func (a *adManagerImpl) GetAdByPositionKey(key string) ad.IAdAggregateRoot {
	a.mux.Lock()
	defer a.mux.Unlock()
	ok := false
	var iv ad.IAdAggregateRoot
	if a.cache == nil {
		a.cache = make(map[string]ad.IAdAggregateRoot)
	}
	//从缓存中获取
	if iv, ok = a.cache[key]; ok {
		return iv
	}

	pos := a.GetPositionByKey(key)
	if pos != nil && pos.PutAid > 0 {
		iv = a.defaultAd.GetById(pos.PutAid)
	}
	if iv != nil {
		a.cache[key] = iv
	}
	return iv
}

// 获取用户的广告管理
func (a *adManagerImpl) GetUserAd(adUserId int) ad.IUserAd {
	return newUserAd(a, a.rep, adUserId)
}

type userAdImpl struct {
	_rep      ad.IAdRepo
	_manager  ad.IAdvertisementManager
	_adUserId int
	_cache    map[string]ad.IAdAggregateRoot
	_mux      sync.Mutex
}

func newUserAd(m ad.IAdvertisementManager, rep ad.IAdRepo, adUserId int) ad.IUserAd {
	return &userAdImpl{
		_rep:      rep,
		_manager:  m,
		_adUserId: adUserId,
	}
}

// 获取聚合根标识
func (a *userAdImpl) GetAggregateRootId() int {
	return a._adUserId
}

// 根据编号获取广告
func (a *userAdImpl) GetById(id int) ad.IAdAggregateRoot {
	v := a._rep.GetAd(int64(id))
	if v != nil {
		return a.CreateAd(v)
	}
	return nil
}

// 获取广告关联的广告位
func (a *userAdImpl) GetAdPositionsByAdId(adId int64) []*ad.Position {
	var list []*ad.Position
	//todo:
	//for _, v := range a._manager.GetAdGroups() {
	//	for _, p := range v.GetPositions() {
	//		if p.PutAdId == adId {
	//			list = append(list, p)
	//		}
	//	}
	//}
	return list
}

func (a *userAdImpl) QueryAdvertisement(keys []string) map[string]ad.IAdAggregateRoot {
	arr := a._rep.GetPositions()
	keyMap := make(map[string]int, len(keys))
	for _, v := range keys {
		keyMap[v] = 0
	}
	mp := make(map[string]ad.IAdAggregateRoot, 0)
	for _, v := range arr {
		if _, ok := keyMap[v.Key]; ok {
			if v.PutAid <= 0 {
				continue
			}
			if ia := a.GetById(v.PutAid); ia != nil {
				mp[v.Key] = ia
			}
		}
	}
	return mp
}

// 删除广告
func (a *userAdImpl) DeleteAd(adId int64) error {
	adv := a.GetById(int(adId))
	if adv != nil {
		if len(a.GetAdPositionsByAdId(adId)) > 0 {
			return ad.ErrAdUsed
		}
		err := a._rep.DeleteAd(int64(a._adUserId), adId)
		a._rep.DeleteImageAdData(adId)
		a._rep.DeleteTextAdData(adId)
		if err == nil {
			a._cache = nil
		}
		return err
	}
	return nil
}

// 根据KEY获取广告
func (a *userAdImpl) GetByPositionKey(key string) ad.IAdAggregateRoot {
	a._mux.Lock()
	defer a._mux.Unlock()
	ok := false
	var iv ad.IAdAggregateRoot
	if a._cache == nil {
		a._cache = make(map[string]ad.IAdAggregateRoot)
	}
	//从缓存中获取
	if iv, ok = a._cache[key]; ok {
		return iv
	}
	//获取用户的设定,如果没有,则获取平台的设定
	v := a._rep.GetAdByKey(int64(a._adUserId), key)
	if v == nil {
		iv = a._manager.GetAdByPositionKey(key)
	} else {
		iv = a.CreateAd(v)
	}
	//加入到缓存
	if iv != nil {
		a._cache[key] = iv
	}
	return iv
}

// 创建广告对象
func (a *userAdImpl) CreateAd(v *ad.Ad) ad.IAdAggregateRoot {
	adv := &adImpl{
		_rep:   a._rep,
		_value: v,
	}
	switch v.TypeId {
	case ad.TypeSwiper:
		// 轮播广告
		return &GalleryAd{
			adImpl: adv,
		}
	case ad.TypeText:
		// 文本广告
		return &HyperLinkAdImpl{
			adImpl: adv,
		}
	case ad.TypeImage:
		// 图片广告
		return &ImageAdImpl{
			adImpl: adv,
		}
	}
	return adv
}

// 检测广告位是否可以被绑定
func (a *userAdImpl) checkPositionBind(posId int64, adId int64) bool {
	total := 0
	_db := provide.GetDb()
	_db.ExecScalar("SELECT COUNT(1) FROM ad_userset WHERE user_id= $1 AND pos_id= $2 AND ad_id <> $3",
		&total, a._adUserId, posId, adId)
	return total == 0
}

// 设置广告
func (a *userAdImpl) SetAd(posId, adId int) error {
	ap := a._manager.GetPosition(int64(posId))
	if ap == nil {
		return ad.ErrNoSuchAdPosition
	}
	if ap.GetValue().Opened == 0 {
		return ad.ErrNotOpened
	}
	if !a.checkPositionBind(int64(posId), int64(adId)) {
		return ad.ErrUserPositionIsBind
	}
	if a._rep.GetAd(int64(adId)) == nil {
		return ad.ErrNoSuchAd
	}
	err := a._rep.SetUserAd(a.GetAggregateRootId(), posId, adId)
	if err == nil {
		a._cache = nil
	}
	return err
}

var _ ad.IAdAggregateRoot = new(adImpl)

type adImpl struct {
	_rep   ad.IAdRepo
	_value *ad.Ad
}

// 获取领域对象编号
func (a *adImpl) GetDomainId() int {
	if a._value != nil {
		return a._value.Id
	}
	return 0
}

// 是否为系统内置的广告
func (a *adImpl) System() bool {
	return a._value.UserId == 0
}

// 广告类型
func (a *adImpl) Type() int {
	return a._value.TypeId
}

// 广告名称
func (a *adImpl) Name() string {
	return a._value.Name
}

// 设置值
func (a *adImpl) SetValue(v *ad.Ad) error {
	if v.TypeId == 0 {
		return ad.ErrAdType
	}
	if v.TypeId != a.Type() {
		return ad.ErrDisallowModifyAdType
	}
	a._value.Name = v.Name
	return nil
}

// 获取值
func (a *adImpl) GetValue() *ad.Ad {
	return a._value
}

// 保存广告
func (a *adImpl) Save() (int64, error) {
	//id := a.Repo.GetIdByName(a.Value.UserId, a.Value.Name)
	//if id > 0 && id != a.GetDomainId() {
	//	return a.GetDomainId(), ad.ErrNameExists
	//}
	a._value.UpdateTime = int(time.Now().Unix())
	return a._rep.SaveAdValue(a._value)
}

// 增加展现次数
func (a *adImpl) AddShowTimes(times int) error {
	a._value.ShowTimes += times
	return nil
}

// 增加展现次数
func (a *adImpl) AddClickTimes(times int) error {
	a._value.ClickTimes += times
	return nil
}

// 增加展现次数
func (a *adImpl) AddShowDays(days int) error {
	a._value.ShowDays += days
	return nil
}

// 转换为数据传输对象
func (a *adImpl) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:     a.GetDomainId(),
		AdType: a.Type(),
	}
}
