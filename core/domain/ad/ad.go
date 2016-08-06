/**
 * Copyright 2015 @ z3q.net.
 * name : partner_ad
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import (
	"go2o/core/domain/interface/ad"
	"sync"
	"time"
)

var _ ad.IUserAd = new(userAdImpl)
var _ ad.IAdGroup = new(AdGroupImpl)
var _ ad.IAdManager = new(adManagerImpl)

type adManagerImpl struct {
	_rep       ad.IAdRep
	_defaultAd ad.IUserAd
	_groups    []ad.IAdGroup
	_mux       sync.Mutex
	_cache     map[string]ad.IAd
}

func NewAdManager(rep ad.IAdRep) ad.IAdManager {
	a := &adManagerImpl{
		_rep: rep,
	}
	a._defaultAd = newUserAd(a, rep, 0)
	return a
}

// 获取广告分组
func (a *adManagerImpl) GetAdGroups() []ad.IAdGroup {
	if a._groups == nil {
		list := a._rep.GetAdGroups()
		a._groups = make([]ad.IAdGroup, len(list))
		for i, v := range list {
			a._groups[i] = newAdGroup(a, a._rep, v)
		}
	}
	return a._groups
}

// 获取单个广告分组
func (a *adManagerImpl) GetAdGroup(id int) ad.IAdGroup {
	list := a.GetAdGroups()
	for _, v := range list {
		if v.GetDomainId() == id {
			return v
		}
	}
	return nil
}

// 删除广告组
func (a *adManagerImpl) DelAdGroup(id int) error {
	a._groups = nil
	return a._rep.DelAdGroup(id)
}

// 创建广告组
func (a *adManagerImpl) CreateAdGroup(name string) ad.IAdGroup {
	return newAdGroup(a, a._rep, &ad.AdGroup{
		Id:      0,
		Name:    name,
		Opened:  1,
		Enabled: 1,
	})
}

// 根据KEY获取广告位
func (a *adManagerImpl) GetAdPositionByKey(key string) *ad.AdPosition {
	return a._rep.GetAdPositionByKey(key)
}

// 根据广告位KEY获取默认广告
func (a *adManagerImpl) GetAdByPositionKey(key string) ad.IAd {
	a._mux.Lock()
	defer a._mux.Unlock()
	ok := false
	var iv ad.IAd
	if a._cache == nil {
		a._cache = make(map[string]ad.IAd)
	}
	//从缓存中获取
	if iv, ok = a._cache[key]; ok {
		return iv
	}

	pos := a.GetAdPositionByKey(key)
	if pos != nil && pos.DefaultId > 0 {
		iv = a._defaultAd.GetById(pos.DefaultId)
	}
	if iv != nil {
		a._cache[key] = iv
	}
	return iv
}

// 获取用户的广告管理
func (a *adManagerImpl) GetUserAd(adUserId int) ad.IUserAd {
	return newUserAd(a, a._rep, adUserId)
}

type AdGroupImpl struct {
	_rep       ad.IAdRep
	_manager   *adManagerImpl
	_value     *ad.AdGroup
	_positions []*ad.AdPosition
}

func newAdGroup(m *adManagerImpl, rep ad.IAdRep, v *ad.AdGroup) ad.IAdGroup {
	return &AdGroupImpl{
		_rep:     rep,
		_manager: m,
		_value:   v,
	}
}

// 获取领域编号
func (a *AdGroupImpl) GetDomainId() int {
	return a._value.Id
}

// 获取值
func (a *AdGroupImpl) GetValue() ad.AdGroup {
	return *a._value
}

// 设置值
func (a *AdGroupImpl) SetValue(v *ad.AdGroup) error {
	if v != nil {
		a._value.Name = v.Name
		a._value.Enabled = v.Enabled
		a._value.Opened = v.Opened
	}
	return nil
}

// 获取广告位
func (a *AdGroupImpl) GetPositions() []*ad.AdPosition {
	if a._positions == nil {
		a._positions = a._rep.GetAdPositionsByGroupId(a.GetDomainId())
	}
	return a._positions
}

// 根据Id获取广告位
func (a *AdGroupImpl) GetPosition(id int) *ad.AdPosition {
	for _, v := range a.GetPositions() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 删除广告位
func (a *AdGroupImpl) DelPosition(id int) error {
	//todo: 广告位已投放广告,不允许删除
	//if a.getAdPositionBindNum(id) > 0{
	//	return ad.err
	//}
	a._positions = nil
	a._manager._cache = nil
	return a._rep.DelAdPosition(id)
}

// 保存广告位
func (ag *AdGroupImpl) SavePosition(a *ad.AdPosition) (int, error) {
	if pos := ag._manager.GetAdPositionByKey(a.Key); pos != nil && pos.Id != a.Id {
		return 0, ad.ErrKeyExists
	}
	a.GroupId = ag.GetDomainId()
	ag._positions = nil
	ag._manager._cache = nil
	return ag._rep.SaveAdPosition(a)
}

// 保存,需调用Save()保存
func (a *AdGroupImpl) Save() (int, error) {
	return a._rep.SaveAdGroup(a._value)
}

// 开放,需调用Save()保存
func (a *AdGroupImpl) Open() error {
	a._value.Opened = 1
	return nil
}

// 关闭,需调用Save()保存
func (a *AdGroupImpl) Close() error {
	a._value.Opened = 0
	return nil
}

// 启用,需调用Save()保存
func (a *AdGroupImpl) Enabled() error {
	a._value.Enabled = 1
	return nil
}

// 禁用,需调用Save()保存
func (a *AdGroupImpl) Disabled() error {
	a._value.Enabled = 0
	return nil
}

// 设置默认广告
func (a *AdGroupImpl) SetDefault(adPosId int, adId int) error {
	if v := a.GetPosition(adPosId); v != nil {
		// if a._rep.GetValueAdvertisement()
		//todo: 检测广告是否存在
		v.DefaultId = adId
		_, err := a.SavePosition(v)
		a._manager._cache = nil
		return err
	}
	return ad.ErrNoSuchAd
}

type userAdImpl struct {
	_rep      ad.IAdRep
	_manager  ad.IAdManager
	_adUserId int
	_cache    map[string]ad.IAd
	_mux      sync.Mutex
}

func newUserAd(m ad.IAdManager, rep ad.IAdRep, adUserId int) ad.IUserAd {
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
func (a *userAdImpl) GetById(id int) ad.IAd {
	v := a._rep.GetValueAd(id)
	if v != nil {
		return a.CreateAd(v)
	}
	return nil
}

// 获取广告关联的广告位
func (a *userAdImpl) GetAdPositionsByAdId(adId int) []*ad.AdPosition {
	list := []*ad.AdPosition{}
	for _, v := range a._manager.GetAdGroups() {
		for _, p := range v.GetPositions() {
			if p.DefaultId == adId {
				list = append(list, p)
			}
		}
	}
	return list
}

// 删除广告
func (a *userAdImpl) DeleteAd(adId int) error {
	adv := a.GetById(adId)
	if adv != nil {
		if len(a.GetAdPositionsByAdId(adId)) > 0 {
			return ad.ErrAdUsed
		}
		err := a._rep.DelAd(a._adUserId, adId)
		a._rep.DelImageDataForAdvertisement(adId)
		a._rep.DelTextDataForAdvertisement(adId)
		if err == nil {
			a._cache = nil
		}
		return err
	}
	return nil
}

// 根据KEY获取广告
func (a *userAdImpl) GetByPositionKey(key string) ad.IAd {
	a._mux.Lock()
	defer a._mux.Unlock()
	ok := false
	var iv ad.IAd
	if a._cache == nil {
		a._cache = make(map[string]ad.IAd)
	}
	//从缓存中获取
	if iv, ok = a._cache[key]; ok {
		return iv
	}
	//获取用户的设定,如果没有,则获取平台的设定
	v := a._rep.GetAdByKey(a._adUserId, key)
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
func (a *userAdImpl) CreateAd(v *ad.Ad) ad.IAd {
	adv := &adImpl{
		_rep:   a._rep,
		_value: v,
	}
	switch v.Type {
	case ad.TypeGallery:
		// 轮播广告
		return &GalleryAd{
			adImpl: adv,
		}
	case ad.TypeHyperLink:
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

// 设置广告
func (a *userAdImpl) SetAd(posId, adId int) error {
	if a._manager.GetAdGroup(posId) == nil {
		return ad.ErrNoSuchAdPosition
	}
	if a._rep.GetValueAd(adId) == nil {
		return ad.ErrNoSuchAd
	}
	err := a._rep.SetUserAd(a.GetAggregateRootId(), posId, adId)
	if err == nil {
		a._cache = nil
	}
	return err
}

var _ ad.IAd = new(adImpl)

type adImpl struct {
	_rep   ad.IAdRep
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
	return a._value.Type
}

// 广告名称
func (a *adImpl) Name() string {
	return a._value.Name
}

// 设置值
func (a *adImpl) SetValue(v *ad.Ad) error {
	if v.Type == 0 {
		return ad.ErrAdType
	}
	if v.Type != a.Type() {
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
func (a *adImpl) Save() (int, error) {
	//id := a.Rep.GetIdByName(a.Value.UserId, a.Value.Name)
	//if id > 0 && id != a.GetDomainId() {
	//	return a.GetDomainId(), ad.ErrNameExists
	//}
	a._value.UpdateTime = time.Now().Unix()
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
		Id:   a.GetDomainId(),
		Type: a.Type(),
	}
}
