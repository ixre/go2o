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
	"time"
)

var _ ad.IUserAd = new(UserAdImpl)
var _ ad.IAdGroup = new(AdGroupImpl)
var _ ad.IAdManager = new(AdManagerImpl)

type AdManagerImpl struct {
	_rep    ad.IAdRep
	_groups []ad.IAdGroup
}

func NewAdManager(rep ad.IAdRep) ad.IAdManager {
	return &AdManagerImpl{
		_rep: rep,
	}
}

// 获取广告分组
func (this *AdManagerImpl) GetAdGroups() []ad.IAdGroup {
	if this._groups == nil {
		list := this._rep.GetAdGroups()
		this._groups = make([]ad.IAdGroup, len(list))
		for i, v := range list {
			this._groups[i] = newAdGroup(this._rep, v)
		}
	}
	return this._groups
}

// 获取单个广告分组
func (this *AdManagerImpl) GetAdGroup(id int) ad.IAdGroup {
	list := this.GetAdGroups()
	for _, v := range list {
		if v.GetDomainId() == id {
			return v
		}
	}
	return nil
}

// 删除广告组
func (this *AdManagerImpl) DelAdGroup(id int) error {
	this._groups = nil
	return this._rep.DelAdGroup(id)
}

// 创建广告组
func (this *AdManagerImpl) CreateAdGroup(name string) ad.IAdGroup {
	return newAdGroup(this._rep, &ad.AdGroup{
		Id:      0,
		Name:    name,
		Opened:  1,
		Enabled: 1,
	})
}

// 获取用户的广告管理
func (this *AdManagerImpl) GetUserAd(adUserId int) ad.IUserAd {
	return newUserAd(this, this._rep, adUserId)
}

type AdGroupImpl struct {
	_rep       ad.IAdRep
	_value     *ad.AdGroup
	_positions []*ad.AdPosition
}

func newAdGroup(rep ad.IAdRep, v *ad.AdGroup) ad.IAdGroup {
	return &AdGroupImpl{
		_rep:   rep,
		_value: v,
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
	return this._rep.DelAdPosition(id)
}

// 保存广告位
func (this *AdGroupImpl) SavePosition(a *ad.AdPosition) (int, error) {
	a.GroupId = this.GetDomainId()
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
		v.Default = adId
		_, err := this.SavePosition(v)
		return err
	}
	return ad.ErrNoSuchAd
}

type UserAdImpl struct {
	_rep      ad.IAdRep
	_manager  ad.IAdManager
	_adUserId int
}

func newUserAd(m ad.IAdManager, rep ad.IAdRep, adUserId int) ad.IUserAd {
	return &UserAdImpl{
		_rep:      rep,
		_manager:  m,
		_adUserId: adUserId,
	}
}

// 获取聚合根标识
func (this *UserAdImpl) GetAggregateRootId() int {
	return this._adUserId
}

// 根据编号获取广告
func (this *UserAdImpl) GetById(id int) ad.IAd {
	v := this._rep.GetValueAdvertisement(id)
	if v != nil {
		return this.CreateAdvertisement(v)
	}
	return nil
}

// 删除广告
func (this *UserAdImpl) DeleteAdvertisement(advertisementId int) error {
	adv := this.GetById(advertisementId)
	if adv != nil {
		if adv.System() {
			return ad.ErrInternalDisallow
		}
		err := this._rep.DelAdvertisement(this._adUserId, advertisementId)
		this._rep.DelImageDataForAdvertisement(advertisementId)
		this._rep.DelTextDataForAdvertisement(advertisementId)
		return err

	}
	return nil
}

// 根据名称获取广告
func (this *UserAdImpl) GetByName(name string) ad.IAd {
	v := this._rep.GetValueAdvertisementByName(this._adUserId, name)
	if v != nil {
		return this.CreateAdvertisement(v)
	}
	return nil
}

// 创建广告对象
func (this *UserAdImpl) CreateAdvertisement(v *ad.Ad) ad.IAd {
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
func (this *UserAdImpl) SetAd(posId, adId int) error {
	if this._manager.GetAdGroup(posId) == nil {
		return ad.ErrNoSuchAdPosition
	}
	if this._rep.GetValueAdvertisement(adId) == nil {
		return ad.ErrNoSuchAd
	}
	return this._rep.SetUserAd(this.GetAggregateRootId(), posId, adId)
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
	// 如果为系统内置广告，不能修改名称
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
