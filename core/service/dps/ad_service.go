/**
 * Copyright 2015 @ z3q.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dps

import (
	"errors"
	"go2o/core/domain/interface/ad"
	"go2o/core/infrastructure/format"
)

type adService struct {
	_rep ad.IAdRep
	//_query     *query.ContentQuery
}

func NewAdvertisementService(rep ad.IAdRep) *adService {
	return &adService{
		_rep: rep,
	}
}

func (a *adService) getUserAd(adUserId int) ad.IUserAd {
	return a._rep.GetAdManager().GetUserAd(adUserId)
}

func (a *adService) GetAdGroupById(id int) ad.AdGroup {
	return a._rep.GetAdManager().GetAdGroup(id).GetValue()
}

func (a *adService) DelAdGroup(id int) error {
	return a._rep.GetAdManager().DelAdGroup(id)
}

func (a *adService) SaveAdGroup(d *ad.AdGroup) (int, error) {
	m := a._rep.GetAdManager()
	var e ad.IAdGroup
	if d.Id > 0 {
		e = m.GetAdGroup(d.Id)
	} else {
		e = m.CreateAdGroup(d.Name)
	}
	err := e.SetValue(d)
	if err == nil {
		return e.Save()
	}
	return -1, err
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

func (a *adService) GetPosition(groupId, adPosId int) *ad.AdPosition {
	return a._rep.GetAdManager().GetAdGroup(groupId).GetPosition(adPosId)
}

func (a *adService) SaveAdPosition(e *ad.AdPosition) (int, error) {
	group := a._rep.GetAdManager().GetAdGroup(e.GroupId)
	if group == nil {
		return -1, errors.New("no such ad group")
	}
	return group.SavePosition(e)
}

func (a *adService) DelAdPosition(groupId, id int) error {
	return a._rep.GetAdManager().GetAdGroup(groupId).DelPosition(id)
}

// 设置广告位的默认广告
func (a *adService) SetDefaultAd(groupId, posId, adId int) error {
	return a._rep.GetAdManager().GetAdGroup(groupId).SetDefault(posId, adId)
}

// 用户投放广告
func (a *adService) SetUserAd(adUserId int, posId int, adId int) error {
	ua := a._rep.GetAdManager().GetUserAd(adUserId)
	return ua.SetAd(posId, adId)
}

// 获取广告
func (a *adService) GetAdvertisement(adUserId, id int) *ad.Ad {
	if adv := a.getUserAd(adUserId).GetById(id); adv != nil {
		return adv.GetValue()
	}
	return nil
}

// 获取广告及广告数据, 用于展示关高
func (a *adService) GetAdAndDataByKey(adUserId int, key string) *ad.AdDto {
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
func (a *adService) GetAdDto(userId int, id int) *ad.AdDto {
	ua := a.getUserAd(userId)
	if adv := ua.GetById(id); adv != nil {
		return adv.Dto()
	}
	return nil
}

// 保存广告,更新时不允许修改类型
func (a *adService) SaveAd(adUserId int, v *ad.Ad) (int, error) {
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

func (a *adService) DeleteAd(adUserId, adId int) error {
	return a.getUserAd(adUserId).DeleteAd(adId)
}

// 保存图片广告
func (a *adService) SaveHyperLinkAd(adUserId int, v *ad.HyperLink) (int, error) {
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
func (a *adService) SaveImageAd(adUserId int, v *ad.Image) (int, error) {
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
func (a *adService) SaveImage(adUserId int, advertisementId int, v *ad.Image) (int, error) {
	pa := a.getUserAd(adUserId)
	var adv ad.IAd = pa.GetById(advertisementId)
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
func (a *adService) GetValueAdImage(adUserId, advertisementId, imgId int) *ad.Image {
	pa := a.getUserAd(adUserId)
	var adv ad.IAd = pa.GetById(advertisementId)
	if adv != nil {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			return gad.GetImage(imgId)
		}
	}
	return nil
}

// 删除广告图片
func (a *adService) DelAdImage(adUserId, advertisementId, imgId int) error {
	pa := a.getUserAd(adUserId)
	var adv ad.IAd = pa.GetById(advertisementId)
	if adv != nil {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			return gad.DelImage(imgId)
		}
	}
	return nil
}
