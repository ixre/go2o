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
func (this *adManagerImpl) GetAdGroups() []ad.IAdGroup {
	if this._groups == nil {
		list := this._rep.GetAdGroups()
		this._groups = make([]ad.IAdGroup, len(list))
		for i, v := range list {
			this._groups[i] = newAdGroup(this, this._rep, v)
		}
	}
	return this._groups
}

// 获取单个广告分组
func (this *adManagerImpl) GetAdGroup(id int) ad.IAdGroup {
	list := this.GetAdGroups()
	for _, v := range list {
		if v.GetDomainId() == id {
			return v
		}
	}
	return nil
}

// 删除广告组
func (this *adManagerImpl) DelAdGroup(id int) error {
	this._groups = nil
	return this._rep.DelAdGroup(id)
}

// 创建广告组
func (this *adManagerImpl) CreateAdGroup(name string) ad.IAdGroup {
	return newAdGroup(this, this._rep, &ad.AdGroup{
		Id:      0,
		Name:    name,
		Opened:  1,
		Enabled: 1,
	})
}

// 根据KEY获取广告位
func (this *adManagerImpl) GetAdPositionByKey(key string) *ad.AdPosition {
	return this._rep.GetAdPositionByKey(key)
}

// 根据广告位KEY获取默认广告
func (this *adManagerImpl) GetAdByPositionKey(key string) ad.IAd {
	this._mux.Lock()
	defer this._mux.Unlock()
	ok := false
	var iv ad.IAd
	if this._cache == nil {
		this._cache = make(map[string]ad.IAd)
	}
	//从缓存中获取
	if iv, ok = this._cache[key]; ok {
		return iv
	}

	pos := this.GetAdPositionByKey(key)
	if pos != nil && pos.DefaultId > 0 {
		iv = this._defaultAd.GetById(pos.DefaultId)
	}
	if iv != nil {
		this._cache[key] = iv
	}
	return iv
}

// 获取用户的广告管理
func (this *adManagerImpl) GetUserAd(adUserId int) ad.IUserAd {
	return newUserAd(this, this._rep, adUserId)
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
func (this *AdGroupImpl) GetDomainId() int {
	return this._value.Id
}

// 获取值
func (this *AdGroupImpl) GetValue() ad.AdGroup {
	return *this._value
}

// 设置值
func (this *AdGroupImpl) SetValue(v *ad.AdGroup) error {
	if v != nil {
		this._value.Name = v.Name
		this._value.Enabled = v.Enabled
		this._value.Opened = v.Opened
	}
	return nil
}

// 获取广告位
func (this *AdGroupImpl) GetPositions() []*ad.AdPosition {
	if this._positions == nil {
		this._positions = this._rep.GetAdPositionsByGroupId(this.GetDomainId())
	}
	return this._positions
}

// 根据Id获取广告位
func (this *AdGroupImpl) GetPosition(id int) *ad.AdPosition {
	for _, v := range this.GetPositions() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 删除广告位
func (this *AdGroupImpl) DelPosition(id int) error {
	//todo: 广告位已投放广告,不允许删除
	//if this.getAdPositionBindNum(id) > 0{
	//	return ad.err
	//}
	this._positions = nil
	this._manager._cache = nil
	return this._rep.DelAdPosition(id)
}

// 保存广告位
func (this *AdGroupImpl) SavePosition(a *ad.AdPosition) (int, error) {
	if pos := this._manager.GetAdPositionByKey(a.Key); pos != nil && pos.Id != a.Id {
		return 0, ad.ErrKeyExists
	}
	a.GroupId = this.GetDomainId()
	this._positions = nil
	this._manager._cache = nil
	return this._rep.SaveAdPosition(a)
}

// 保存,需调用Save()保存
func (this *AdGroupImpl) Save() (int, error) {
	return this._rep.SaveAdGroup(this._value)
}

// 开放,需调用Save()保存
func (this *AdGroupImpl) Open() error {
	this._value.Opened = 1
	return nil
}

// 关闭,需调用Save()保存
func (this *AdGroupImpl) Close() error {
	this._value.Opened = 0
	return nil
}

// 启用,需调用Save()保存
func (this *AdGroupImpl) Enabled() error {
	this._value.Enabled = 1
	return nil
}

// 禁用,需调用Save()保存
func (this *AdGroupImpl) Disabled() error {
	this._value.Enabled = 0
	return nil
}

// 设置默认广告
func (this *AdGroupImpl) SetDefault(adPosId int, adId int) error {
	if v := this.GetPosition(adPosId); v != nil {
		// if this._rep.GetValueAdvertisement()
		//todo: 检测广告是否存在
		v.DefaultId = adId
		_, err := this.SavePosition(v)
		this._manager._cache = nil
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
func (this *userAdImpl) GetAggregateRootId() int {
	return this._adUserId
}

// 根据编号获取广告
func (this *userAdImpl) GetById(id int) ad.IAd {
	v := this._rep.GetValueAd(id)
	if v != nil {
		return this.CreateAd(v)
	}
	return nil
}

// 获取广告关联的广告位
func (this *userAdImpl) GetAdPositionsByAdId(adId int) []*ad.AdPosition {
	list := []*ad.AdPosition{}
	for _, v := range this._manager.GetAdGroups() {
		for _, p := range v.GetPositions() {
			if p.DefaultId == adId {
				list = append(list, p)
			}
		}
	}
	return list
}

// 删除广告
func (this *userAdImpl) DeleteAd(adId int) error {
	adv := this.GetById(adId)
	if adv != nil {
		if len(this.GetAdPositionsByAdId(adId)) > 0 {
			return ad.ErrAdUsed
		}
		err := this._rep.DelAd(this._adUserId, adId)
		this._rep.DelImageDataForAdvertisement(adId)
		this._rep.DelTextDataForAdvertisement(adId)
		if err == nil {
			this._cache = nil
		}
		return err
	}
	return nil
}

// 根据KEY获取广告
func (this *userAdImpl) GetByPositionKey(key string) ad.IAd {
	this._mux.Lock()
	defer this._mux.Unlock()
	ok := false
	var iv ad.IAd
	if this._cache == nil {
		this._cache = make(map[string]ad.IAd)
	}
	//从缓存中获取
	if iv, ok = this._cache[key]; ok {
		return iv
	}
	//获取用户的设定,如果没有,则获取平台的设定
	v := this._rep.GetAdByKey(this._adUserId, key)
	if v == nil {
		iv = this._manager.GetAdByPositionKey(key)
	} else {
		iv = this.CreateAd(v)
	}
	//加入到缓存
	if iv != nil {
		this._cache[key] = iv
	}
	return iv
}

// 创建广告对象
func (this *userAdImpl) CreateAd(v *ad.Ad) ad.IAd {
	adv := &AdImpl{
		_rep:   this._rep,
		_value: v,
	}
	switch v.Type {
	case ad.TypeGallery:
		// 轮播广告
		return &GalleryAd{
			AdImpl: adv,
		}
	case ad.TypeHyperLink:
		// 文本广告
		return &HyperLinkAdImpl{
			AdImpl: adv,
		}
	case ad.TypeImage:
		// 图片广告
		return &ImageAdImpl{
			AdImpl: adv,
		}
	}
	return adv
}

// 设置广告
func (this *userAdImpl) SetAd(posId, adId int) error {
	if this._manager.GetAdGroup(posId) == nil {
		return ad.ErrNoSuchAdPosition
	}
	if this._rep.GetValueAd(adId) == nil {
		return ad.ErrNoSuchAd
	}
	err := this._rep.SetUserAd(this.GetAggregateRootId(), posId, adId)
	if err == nil {
		this._cache = nil
	}
	return err
}

var _ ad.IAd = new(AdImpl)

type AdImpl struct {
	_rep   ad.IAdRep
	_value *ad.Ad
}

// 获取领域对象编号
func (this *AdImpl) GetDomainId() int {
	if this._value != nil {
		return this._value.Id
	}
	return 0
}

// 是否为系统内置的广告
func (this *AdImpl) System() bool {
	return this._value.UserId == 0
}

// 广告类型
func (this *AdImpl) Type() int {
	return this._value.Type
}

// 广告名称
func (this *AdImpl) Name() string {
	return this._value.Name
}

// 设置值
func (this *AdImpl) SetValue(v *ad.Ad) error {
	if v.Type != this.Type() {
		return ad.ErrDisallowModifyAdType
	}
	this._value.Name = v.Name
	return nil
}

// 获取值
func (this *AdImpl) GetValue() *ad.Ad {
	return this._value
}

// 保存广告
func (this *AdImpl) Save() (int, error) {
	//id := this.Rep.GetIdByName(this.Value.UserId, this.Value.Name)
	//if id > 0 && id != this.GetDomainId() {
	//	return this.GetDomainId(), ad.ErrNameExists
	//}
	this._value.UpdateTime = time.Now().Unix()
	return this._rep.SaveAdValue(this._value)
}

// 增加展现次数
func (this *AdImpl) AddShowTimes(times int) error {
	this._value.ShowTimes += times
	return nil
}

// 增加展现次数
func (this *AdImpl) AddClickTimes(times int) error {
	this._value.ClickTimes += times
	return nil
}

// 增加展现次数
func (this *AdImpl) AddShowDays(days int) error {
	this._value.ShowDays += days
	return nil
}

// 转换为数据传输对象
func (this *AdImpl) Dto() *ad.AdDto {
	return &ad.AdDto{
		Id:   this.GetDomainId(),
		Type: this.Type(),
	}
}
