/**
 * Copyright 2015 @ z3q.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package rsi

import (
	"errors"
	"fmt"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/ad"
	"go2o/core/infrastructure/format"
	"go2o/core/repos"
)

type adService struct {
	_rep    ad.IAdRepo
	storage storage.Interface
	//_query     *query.ContentQuery
}

func NewAdvertisementService(rep ad.IAdRepo, storage storage.Interface) *adService {
	return &adService{
		_rep:    rep,
		storage: storage,
	}
}

func (a *adService) getUserAd(adUserId int32) ad.IUserAd {
	return a._rep.GetAdManager().GetUserAd(adUserId)
}

func (a *adService) GetAdGroupById(id int32) ad.AdGroup {
	return a._rep.GetAdManager().GetAdGroup(id).GetValue()
}

func (a *adService) DelAdGroup(id int32) error {
	return a._rep.GetAdManager().DelAdGroup(id)
}

func (a *adService) SaveAdGroup(d *ad.AdGroup) (int32, error) {
	m := a._rep.GetAdManager()
	var e ad.IAdGroup
	if d.ID > 0 {
		e = m.GetAdGroup(d.ID)
	} else {
		e = m.CreateAdGroup(d.Name)
	}
	err := e.SetValue(d)
	if err == nil {
		return e.Save()
	}
	return 0, err
}

func (a *adService) GetGroups() []ad.AdGroup {
	m := a._rep.GetAdManager()
	list := m.GetAdGroups()
	list2 := []ad.AdGroup{}
	for _, v := range list {
		list2 = append(list2, v.GetValue())
	}
	return list2
}

func (a *adService) GetPosition(groupId, adPosId int32) *ad.AdPosition {
	return a._rep.GetAdManager().GetAdGroup(groupId).GetPosition(adPosId)
}

func (a *adService) GetPositionById(adPosId int32) *ad.AdPosition {
	return a._rep.GetAdPositionById(adPosId)
}

func (a *adService) SaveAdPosition(e *ad.AdPosition) (int32, error) {
	group := a._rep.GetAdManager().GetAdGroup(e.GroupId)
	if group == nil {
		return -1, errors.New("no such ad group")
	}
	return group.SavePosition(e)
}

func (a *adService) DelAdPosition(groupId, id int32) error {
	return a._rep.GetAdManager().GetAdGroup(groupId).DelPosition(id)
}

// 设置广告位的默认广告
func (a *adService) SetDefaultAd(groupId, posId, adId int32) error {
	return a._rep.GetAdManager().GetAdGroup(groupId).SetDefault(posId, adId)
}

// 用户投放广告
func (a *adService) SetUserAd(adUserId int32, posId int32, adId int32) error {
	defer a.cleanCache(adUserId)
	ua := a._rep.GetAdManager().GetUserAd(adUserId)
	return ua.SetAd(posId, adId)
}

// 获取广告
func (a *adService) GetAdvertisement(adUserId, id int32) *ad.Ad {
	if adv := a.getUserAd(adUserId).GetById(id); adv != nil {
		return adv.GetValue()
	}
	return nil
}

// 获取广告及广告数据, 用于展示关高
func (a *adService) GetAdAndDataByKey(adUserId int32, key string) *ad.AdDto {
	if adv := a.getUserAd(adUserId).GetByPositionKey(key); adv != nil {
		switch adv.Type() {
		case ad.TypeGallery:
			dto := adv.Dto()
			gallary := dto.Data.(ad.ValueGallery)
			for _, v := range gallary {
				v.ImageUrl = format.GetResUrl(v.ImageUrl)
			}
			dto.Data = gallary
			return dto
		case ad.TypeImage:
			dto := adv.Dto()
			img := dto.Data.(*ad.Image)
			img.ImageUrl = format.GetResUrl(img.ImageUrl)
			return dto
		}
		return adv.Dto()
	}
	return nil
}

// 获取广告数据传输对象
func (a *adService) GetAdDto(userId int32, id int32) *ad.AdDto {
	ua := a.getUserAd(userId)
	if adv := ua.GetById(id); adv != nil {
		return adv.Dto()
	}
	return nil
}

// 保存广告,更新时不允许修改类型
func (a *adService) SaveAd(adUserId int32, v *ad.Ad) (int32, error) {
	defer a.cleanCache(adUserId)
	pa := a.getUserAd(adUserId)
	var adv ad.IAd
	if v.Id > 0 {
		adv = pa.GetById(v.Id)
	} else {
		adv = pa.CreateAd(v)
	}
	err := adv.SetValue(v)
	if err != nil {
		return -1, err
	}
	return adv.Save()
}

func (a *adService) DeleteAd(adUserId, adId int32) error {
	defer a.cleanCache(adUserId)
	return a.getUserAd(adUserId).DeleteAd(adId)
}

// 保存图片广告
func (a *adService) SaveHyperLinkAd(adUserId int32, v *ad.HyperLink) (int32, error) {
	defer a.cleanCache(adUserId)
	pa := a.getUserAd(adUserId)
	var adv ad.IAd = pa.GetById(v.AdId)
	if adv.Type() == ad.TypeHyperLink {
		g := adv.(ad.IHyperLinkAd)
		g.SetData(v)
		return adv.Save()
	}
	return -1, nil
}

// 保存图片广告
func (a *adService) SaveImageAd(adUserId int32, v *ad.Image) (int32, error) {
	pa := a.getUserAd(adUserId)
	var adv ad.IAd = pa.GetById(v.AdId)
	if adv.Type() == ad.TypeImage {
		g := adv.(ad.IImageAd)
		g.SetData(v)
		return adv.Save()
	}
	return -1, nil
}

// 保存广告图片
func (a *adService) SaveImage(adUserId int32, adId int32, v *ad.Image) (int32, error) {
	defer a.cleanCache(adUserId)
	pa := a.getUserAd(adUserId)
	var adv ad.IAd = pa.GetById(adId)
	if adv != nil {
		switch adv.Type() {
		case ad.TypeGallery:
			gad := adv.(ad.IGalleryAd)
			return gad.SaveImage(v)
		case ad.TypeImage:
			gad := adv.(ad.IImageAd)
			gad.SetData(v)
			return adv.Save()
		}
	}
	return -1, ad.ErrNoSuchAd
}

// 获取广告图片
func (a *adService) GetValueAdImage(adUserId, adId, imgId int32) *ad.Image {
	pa := a.getUserAd(adUserId)
	var adv ad.IAd = pa.GetById(adId)
	if adv != nil {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			return gad.GetImage(imgId)
		}
	}
	return nil
}

// 删除广告图片
func (a *adService) DelAdImage(adUserId, adId, imgId int32) error {
	defer a.cleanCache(adUserId)
	pa := a.getUserAd(adUserId)
	var adv ad.IAd = pa.GetById(adId)
	if adv != nil {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			return gad.DelImage(imgId)
		}
	}
	return nil
}

func (a *adService) cleanCache(adUserId int32) error {
	return repos.PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", adUserId))
}
