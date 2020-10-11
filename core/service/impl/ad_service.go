/**
 * Copyright 2015 @ to2.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types"
	"go2o/core/domain/interface/ad"
	"go2o/core/repos"
	"go2o/core/service/proto"
)

var _ proto.AdvertisementServiceServer = new(adService)

type adService struct {
	_rep    ad.IAdRepo
	storage storage.Interface
	//_query     *shopQuery.ContentQuery
	serviceUtil
}

func (a *adService) GetAdGroupById(_ context.Context, id *proto.Int64) (*proto.SAdGroup, error) {
	ig := a._rep.GetAdManager().GetAdGroup(id.Value)
	if ig != nil {
		return a.parseAdGroupDto(ig.GetValue()), nil
	}
	return nil, nil
}

func (a *adService) DelAdGroup(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	err := a._rep.GetAdManager().DelAdGroup(id.Value)
	return a.error(err), nil
}

func (a *adService) SaveAdGroup(_ context.Context, group *proto.SAdGroup) (*proto.Result, error) {
	m := a._rep.GetAdManager()
	var e ad.IAdGroup
	v := a.parseAdGroup(group)
	if v.ID > 0 {
		e = m.GetAdGroup(v.ID)
		if e == nil {
			return a.error(ad.ErrNoSuchAdGroup), nil
		}
	} else {
		e = m.CreateAdGroup(v.Name)
	}
	err := e.SetValue(v)
	if err == nil {
		_, err = e.Save()
	}
	return a.error(err), nil
}

func (a *adService) GetGroups(_ context.Context, empty *proto.Empty) (*proto.AdGroupListResponse, error) {
	m := a._rep.GetAdManager()
	list := m.GetAdGroups()
	arr := make([]*proto.SAdGroup, 0)
	for _, v := range list {
		arr = append(arr, a.parseAdGroupDto(v.GetValue()))
	}
	return &proto.AdGroupListResponse{Value: arr}, nil
}

func (a *adService) GetPosition(_ context.Context, id *proto.AdPositionId) (*proto.SAdPosition, error) {
	ip := a._rep.GetAdPositionById(id.Value)
	if ip != nil {
		return a.parseAdPositionDto(ip), nil
	}
	return nil, nil
}

func (a *adService) SaveAdPosition(_ context.Context, r *proto.SAdPosition) (*proto.Result, error) {
	var err error
	group := a._rep.GetAdManager().GetAdGroup(r.GroupId)
	if group == nil {
		err = ad.ErrNoSuchAdGroup
	} else {
		v := a.parseAdPosition(r)
		_, err = group.SavePosition(v)
	}
	return a.error(err), nil
}

func (a *adService) DelAdPosition(_ context.Context, id *proto.AdPositionId) (*proto.Result, error) {
	err := a._rep.DelAdPosition(id.Value)
	return a.error(err), nil
}

// 设置广告位的默认广告
func (a *adService) SetDefaultAd(_ context.Context, r *proto.SetDefaultAdRequest) (*proto.Result, error) {
	ig := a._rep.GetAdManager().GetAdGroup(r.GroupId)
	err := ig.SetDefault(r.PosId, r.AdId)
	return a.error(err), nil
}

// 用户投放广告
func (a *adService) SetUserAd(_ context.Context, r *proto.SetUserAdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	ua := a._rep.GetAdManager().GetUserAd(r.AdUserId)
	err := ua.SetAd(r.PosId, r.AdId)
	return a.error(err), nil
}

// 获取广告
func (a *adService) GetAdvertisement(_ context.Context, r *proto.AdIdRequest) (*proto.SAd, error) {
	ig := a.getUserAd(r.AdUserId).GetById(r.AdId)
	if ig != nil {
		return a.parseAdDto(ig.GetValue()), nil
	}
	return nil, nil
}

// 获取广告及广告数据, 用于展示关高
func (a *adService) GetAdAndDataByKey(_ context.Context, r *proto.AdKeyRequest) (*proto.SAdDto, error) {
	//if adv := a.getUserAd(adUserId).GetByPositionKey(key); adv != nil {
	//	switch adv.Type() {
	//	case ad.TypeGallery:
	//		dto := adv.Dto()
	//		gallary := dto.Data.(ad.ValueGallery)
	//		for _, v := range gallary {
	//			v.ImageUrl = format.GetResUrl(v.ImageUrl)
	//		}
	//		dto.Data = gallary
	//		return dto
	//	case ad.TypeImage:
	//		dto := adv.Dto()
	//		img := dto.Data.(*ad.Image)
	//		img.ImageUrl = format.GetResUrl(img.ImageUrl)
	//		return dto
	//	}
	//	return adv.Dto()
	//}
	//return nil
	panic("not implement")
}

// 获取广告数据传输对象
func (a *adService) GetAdDto_(_ context.Context, r *proto.AdIdRequest) (*proto.SAdDto, error) {
	ua := a.getUserAd(r.AdUserId).GetById(r.AdId)
	if ua != nil {
		//return ua.Dto(),nil
	}
	return nil, nil
}

// 保存广告,更新时不允许修改类型
func (a *adService) SaveAd(_ context.Context, r *proto.SaveAdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	pa := a.getUserAd(r.AdUserId)
	var adv ad.IAd
	v := a.parseAd(r.Value)
	if v.Id > 0 {
		adv = pa.GetById(v.Id)
	} else {
		adv = pa.CreateAd(v)
	}
	err := adv.SetValue(v)
	if err == nil {
		_, err = adv.Save()
	}
	return a.error(err), nil
}

func (a *adService) DeleteAd(_ context.Context, r *proto.AdIdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	err := a.getUserAd(r.AdUserId).DeleteAd(r.AdId)
	return a.error(err), nil
}

// 保存图片广告
func (a *adService) SaveHyperLinkAd(_ context.Context, r *proto.SaveLinkAdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	pa := a.getUserAd(r.AdUserId)
	var adv = pa.GetById(r.AdId)
	var err error
	if adv == nil {
		err = ad.ErrNoSuchAd
	} else {
		if adv.Type() == ad.TypeHyperLink {
			g := adv.(ad.IHyperLinkAd)
			v := a.parseHyperLinkAd(r.Value)
			err = g.SetData(v)
			if err == nil {
				_, err = adv.Save()
			}
		}
	}
	return a.error(err), nil
}

// 保存图片广告
func (a *adService) SaveImagOfAd(_ context.Context, r *proto.SaveImageAdRequest) (*proto.Result, error) {
	ia := a.getUserAd(r.AdUserId).GetById(r.AdId)
	var err error
	if ia == nil {
		err = errors.New("no such ad image")
	} else {
		if ia.Type() == ad.TypeImage {
			g := ia.(ad.IImageAd)
			err = g.SetData(a.parseAdImage(r.Value))
			if err == nil {
				_, err = ia.Save()
			}
		}
	}
	return a.error(err), nil
}

// 保存广告图片
func (a *adService) SaveImage(adUserId int64, adId int64, v *ad.Image) (int64, error) {
	defer a.cleanCache(adUserId)
	pa := a.getUserAd(adUserId)
	var adv = pa.GetById(adId)
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
func (a *adService) GetValueAdImage(_ context.Context, r *proto.ImageAdIdRequest) (*proto.SImage, error) {
	ia := a.getUserAd(r.AdUserId).GetById(r.AdId)
	if ia != nil {
		if ia.Type() == ad.TypeGallery {
			gad := ia.(ad.IGalleryAd)
			return a.parseAdImageDto(gad.GetImage(r.ImageId)), nil
		}
	}
	return nil, nil
}

// 删除广告图片
func (a *adService) DelAdImage(_ context.Context, r *proto.ImageAdIdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	pa := a.getUserAd(r.AdUserId)
	var adv = pa.GetById(r.AdId)
	var err error
	if adv == nil {
		err = errors.New("no such ad image")
	} else {
		if adv.Type() == ad.TypeGallery {
			gad := adv.(ad.IGalleryAd)
			err = gad.DelImage(r.ImageId)
		}
	}
	return a.error(err), nil
}

func NewAdvertisementService(rep ad.IAdRepo, storage storage.Interface) *adService {
	return &adService{
		_rep:    rep,
		storage: storage,
	}
}

func (a *adService) getUserAd(adUserId int64) ad.IUserAd {
	return a._rep.GetAdManager().GetUserAd(adUserId)
}

func (a *adService) GetPositionById(adPosId int64) *ad.AdPosition {
	return a._rep.GetAdPositionById(adPosId)
}

func (a *adService) cleanCache(adUserId int64) error {
	return repos.PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", adUserId))
}

func (a *adService) parseAdGroup(v *proto.SAdGroup) *ad.AdGroup {
	return &ad.AdGroup{
		ID:      v.Id,
		Name:    v.Name,
		Opened:  types.IntCond(v.Opened, 1, 0),
		Enabled: types.IntCond(v.Enabled, 1, 0),
	}
}

func (a *adService) parseAdGroupDto(v ad.AdGroup) *proto.SAdGroup {
	return &proto.SAdGroup{
		Id:      v.ID,
		Name:    v.Name,
		Opened:  v.Opened == 1,
		Enabled: v.Enabled == 1,
	}
}

func (a *adService) parseAdPositionDto(v *ad.AdPosition) *proto.SAdPosition {
	return &proto.SAdPosition{
		Id:        v.ID,
		GroupId:   v.GroupId,
		Key:       v.Key,
		Name:      v.Name,
		TypeLimit: int32(v.TypeLimit),
		Opened:    v.Opened == 1,
		Enabled:   v.Enabled == 1,
		DefaultId: v.DefaultId,
	}
}

func (a *adService) parseAdPosition(v *proto.SAdPosition) *ad.AdPosition {
	return &ad.AdPosition{
		ID:        v.Id,
		GroupId:   v.GroupId,
		Key:       v.Key,
		Name:      v.Name,
		TypeLimit: int(v.TypeLimit),
		Opened:    types.IntCond(v.Opened, 1, 0),
		Enabled:   types.IntCond(v.Enabled, 1, 0),
		DefaultId: v.DefaultId,
	}
}

func (a *adService) parseAdDto(v *ad.Ad) *proto.SAd {
	return &proto.SAd{
		Id:         v.Id,
		UserId:     v.UserId,
		Name:       v.Name,
		Type:       int32(v.Type),
		ShowTimes:  int32(v.ShowTimes),
		ClickTimes: int32(v.ClickTimes),
		ShowDays:   int32(v.ShowDays),
		UpdateTime: v.UpdateTime,
	}
}

func (a *adService) parseAd(v *proto.SAd) *ad.Ad {
	return &ad.Ad{
		Id:         v.Id,
		UserId:     v.UserId,
		Name:       v.Name,
		Type:       int(v.Type),
		ShowTimes:  int(v.ShowTimes),
		ClickTimes: int(v.ClickTimes),
		ShowDays:   int(v.ShowDays),
		UpdateTime: v.UpdateTime,
	}
}

func (a *adService) parseHyperLinkAd(v *proto.SHyperLink) *ad.HyperLink {
	return &ad.HyperLink{
		Id:      v.Id,
		AdId:    v.AdId,
		Title:   v.Title,
		LinkUrl: v.LinkUrl,
	}
}

func (a *adService) parseAdImageDto(v *ad.Image) *proto.SImage {
	return &proto.SImage{
		Id:       v.Id,
		AdId:     v.AdId,
		Title:    v.Title,
		LinkUrl:  v.LinkUrl,
		ImageUrl: v.ImageUrl,
		SortNum:  int32(v.SortNum),
		Enabled:  v.Enabled == 1,
	}
}

func (a *adService) parseAdImage(v *proto.SImage) *ad.Image {
	return &ad.Image{
		Id:       v.Id,
		AdId:     v.AdId,
		Title:    v.Title,
		LinkUrl:  v.LinkUrl,
		ImageUrl: v.ImageUrl,
		SortNum:  int(v.SortNum),
		Enabled:  types.IntCond(v.Enabled, 1, 0),
	}
}
