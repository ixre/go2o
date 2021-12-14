/**
 * Copyright 2015 @ 56x.net.
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
	"go2o/core/domain/interface/ad"
	"go2o/core/repos"
	"go2o/core/service/proto"
)

var _ proto.AdvertisementServiceServer = new(advertisementService)

type advertisementService struct {
	_rep    ad.IAdRepo
	storage storage.Interface
	//_query     *shopQuery.ContentQuery
	serviceUtil
}

func (a *advertisementService) GetGroups(_ context.Context, empty *proto.Empty) (*proto.AdGroupResponse, error) {
	repo := a._rep.GetAdManager()
	var arr = repo.GetGroups()
	return &proto.AdGroupResponse{Value: arr}, nil
}

func (a *advertisementService) GetPosition(_ context.Context, id *proto.AdPositionId) (*proto.SAdPosition, error) {
	ip := a._rep.GetAdPositionById(id.PositionId)
	if ip != nil {
		return a.parseAdPositionDto(ip), nil
	}
	return nil, nil
}

func (a *advertisementService) SaveAdPosition(_ context.Context, r *proto.SAdPosition) (*proto.Result, error) {
	v := a.parseAdPosition(r)
	var ap ad.IAdPosition
	if r.Id > 0 {
		ap = a._rep.GetPosition(r.Id)
		if ap == nil {
			return a.error(ad.ErrNoSuchAdPosition), nil
		}
	} else {
		ap = a._rep.CreateAdPosition(v)
	}
	err := ap.SetValue(v)
	if err == nil {
		err = ap.Save()
	}
	return a.error(err), nil
}

func (a *advertisementService) DeleteAdPosition(_ context.Context, id *proto.AdPositionId) (*proto.Result, error) {
	err := a._rep.DeleteAdPosition(id.PositionId)
	return a.error(err), nil
}

func (a *advertisementService) QueryAd(_ context.Context, request *proto.QueryAdRequest) (*proto.QueryAdResponse, error) {
	if request.Size <= 0 || request.Size > 10 {
		request.Size = 10
	}
	ret := a._rep.GetAdManager().QueryAd(request.Keyword, int(request.Size))
	rsp := &proto.QueryAdResponse{
		Value: make([]*proto.SAdDto, len(ret)),
	}
	for i, v := range ret {
		rsp.Value[i] = &proto.SAdDto{
			Id:   v.Id,
			Name: v.Name,
		}
	}
	return rsp, nil
}

// PutDefaultAd 设置广告位的默认广告
func (a *advertisementService) PutDefaultAd(_ context.Context, r *proto.SetDefaultAdRequest) (*proto.Result, error) {
	ig := a._rep.GetAdManager().GetPosition(r.PositionId)
	if ig == nil {
		return a.error(ad.ErrNoSuchAdPosition), nil
	}
	err := ig.PutAd(r.AdId)
	return a.error(err), nil
}

// 用户投放广告
func (a *advertisementService) SetUserAd(_ context.Context, r *proto.SetUserAdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	ua := a._rep.GetAdManager().GetUserAd(r.AdUserId)
	err := ua.SetAd(r.PosId, r.AdId)
	return a.error(err), nil
}

// 获取广告
func (a *advertisementService) GetAdvertisement(_ context.Context, r *proto.AdIdRequest) (*proto.SAdDto, error) {
	ia := a.getUserAd(r.AdUserId).GetById(r.AdId)
	if ia != nil {
		ret := a.parseAdDto(ia.GetValue())
		if r.ReturnData {
			ret.Data = a.getAdvertisementDto(ia)
		}
		return ret, nil
	}
	return nil, nil
}

//  获取广告数据传输对象
func (a *advertisementService) getAdvertisementDto(ia ad.IAd) *proto.SAdvertisementDto {
	dto := ia.Dto()
	ret := &proto.SAdvertisementDto{Id: dto.Id, AdType: int32(dto.AdType)}
	switch dto.AdType {
	case ad.TypeText:
		ret.Text = a.parseTextDto(dto)
	case ad.TypeImage:
		ret.Image = a.parseImageDto(dto)
	case ad.TypeSwiper:
		ret.Swiper = a.parseSwiperDto(dto)
	default:
		panic("not support ad type")
	}
	return ret
}

func (a *advertisementService) QueryAdvertisementData(_ context.Context, r *proto.QueryAdvertisementDataRequest) (*proto.QueryAdvertisementDataResponse, error) {
	iu := a.getUserAd(r.AdUserId)
	var list = iu.QueryAdvertisement(r.Keys)
	arr := make(map[string]*proto.SAdvertisementDto, len(list))
	for k, v := range list {
		arr[k] = a.getAdvertisementDto(v)
	}
	return &proto.QueryAdvertisementDataResponse{
		AdUserId: r.AdUserId,
		Value:    arr,
	}, nil
}

// 保存广告,更新时不允许修改类型
func (a *advertisementService) SaveAd(_ context.Context, r *proto.SaveAdRequest) (*proto.Result, error) {
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
		if _, err = adv.Save(); err == nil {
			err = a.updateAdData(adv, r.Data)
		}

	}
	return a.error(err), nil
}

func (a *advertisementService) DeleteAd(_ context.Context, r *proto.AdIdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	err := a.getUserAd(r.AdUserId).DeleteAd(r.AdId)
	return a.error(err), nil
}

// 保存图片广告
func (a *advertisementService) SaveSwiperAdImage(_ context.Context, r *proto.SaveSwiperImageRequest) (*proto.Result, error) {
	ia := a.getUserAd(r.AdUserId).GetById(r.AdId)
	var err error
	if ia == nil {
		err = errors.New("no such ad image")
	} else {
		if ia.Type() == ad.TypeSwiper {
			gad := ia.(ad.IGalleryAd)
			v := a.parseAdImage(r.Value)
			_, err = gad.SaveImage(v)
		}
	}
	return a.error(err), nil
}

// 获取广告图片
func (a *advertisementService) GetSwiperAdImage(_ context.Context, r *proto.ImageIdRequest) (*proto.SImageAdData, error) {
	ia := a.getUserAd(r.AdUserId).GetById(r.AdId)
	if ia != nil {
		if ia.Type() == ad.TypeSwiper {
			gad := ia.(ad.IGalleryAd)
			return a.parseAdImageDto(gad.GetImage(r.ImageId)), nil
		}
	}
	return nil, nil
}

// DeleteSwiperAdImage 删除广告图片
func (a *advertisementService) DeleteSwiperAdImage(_ context.Context, r *proto.ImageIdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	pa := a.getUserAd(r.AdUserId)
	var adv = pa.GetById(r.AdId)
	var err error
	if adv == nil {
		err = errors.New("no such ad image")
	} else {
		if adv.Type() == ad.TypeSwiper {
			gad := adv.(ad.IGalleryAd)
			err = gad.DeleteItem(r.ImageId)
		}
	}
	return a.error(err), nil
}

func NewAdvertisementService(rep ad.IAdRepo, storage storage.Interface) *advertisementService {
	return &advertisementService{
		_rep:    rep,
		storage: storage,
	}
}

func (a *advertisementService) getUserAd(adUserId int64) ad.IUserAd {
	return a._rep.GetAdManager().GetUserAd(adUserId)
}

func (a *advertisementService) GetPositionById(adPosId int64) *ad.Position {
	return a._rep.GetAdPositionById(adPosId)
}

func (a *advertisementService) cleanCache(adUserId int64) error {
	return repos.PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", adUserId))
}

func (a *advertisementService) parseAdPositionDto(v *ad.Position) *proto.SAdPosition {
	return &proto.SAdPosition{
		Id:   v.Id,
		Key:  v.Key,
		Name: v.Name,
		Flag: int32(v.Flag),
		//TypeLimit: int32(v.TypeLimit),
		Opened:    int32(v.Opened),
		Enabled:   int32(v.Enabled),
		PutAid:    v.PutAdId,
		GroupName: v.GroupName,
	}
}

func (a *advertisementService) parseAdPosition(v *proto.SAdPosition) *ad.Position {
	return &ad.Position{
		Id:   v.Id,
		Key:  v.Key,
		Name: v.Name,
		Flag: int(v.Flag),
		//TypeLimit: int(v.TypeLimit),
		Opened:    int(v.Opened),
		Enabled:   int(v.Enabled),
		PutAdId:   v.PutAid,
		GroupName: v.GroupName,
	}
}

func (a *advertisementService) parseAdDto(v *ad.Ad) *proto.SAdDto {
	return &proto.SAdDto{
		Id:         v.Id,
		UserId:     v.UserId,
		Name:       v.Name,
		AdType:     int32(v.AdType),
		ShowTimes:  int32(v.ShowTimes),
		ClickTimes: int32(v.ClickTimes),
		ShowDays:   int32(v.ShowDays),
		UpdateTime: v.UpdateTime,
	}
}

func (a *advertisementService) parseAd(v *proto.SAdDto) *ad.Ad {
	return &ad.Ad{
		Id:         v.Id,
		UserId:     v.UserId,
		Name:       v.Name,
		AdType:     int(v.AdType),
		ShowTimes:  int(v.ShowTimes),
		ClickTimes: int(v.ClickTimes),
		ShowDays:   int(v.ShowDays),
		UpdateTime: v.UpdateTime,
	}
}

func (a *advertisementService) parseHyperLinkAd(v *proto.STextAdData) *ad.HyperLink {
	return &ad.HyperLink{
		Id:      v.Id,
		Title:   v.Title,
		LinkUrl: v.LinkURL,
	}
}

func (a *advertisementService) parseAdImageDto(v *ad.Image) *proto.SImageAdData {
	return &proto.SImageAdData{
		Id:       v.Id,
		Title:    v.Title,
		LinkURL:  v.LinkUrl,
		ImageURL: v.ImageUrl,
		SortNum:  int32(v.SortNum),
		Enabled:  int32(v.Enabled),
	}
}

func (a *advertisementService) parseAdImage(v *proto.SImageAdData) *ad.Image {
	return &ad.Image{
		Id:       v.Id,
		Title:    v.Title,
		LinkUrl:  v.LinkURL,
		ImageUrl: v.ImageURL,
		SortNum:  int(v.SortNum),
		Enabled:  int(v.Enabled),
	}
}

func (a *advertisementService) parseTextDto(dto *ad.AdDto) *proto.STextAdData {
	v := dto.Data.(*ad.HyperLink)
	return &proto.STextAdData{
		Id:      v.Id,
		Title:   v.Title,
		LinkURL: v.LinkUrl,
	}
}

func (a *advertisementService) parseImageDto(dto *ad.AdDto) *proto.SImageAdData {
	v := dto.Data.(*ad.Image)
	return a.parseSingleImageDto(v)
}

func (a *advertisementService) parseSingleImageDto(v *ad.Image) *proto.SImageAdData {
	return &proto.SImageAdData{
		Id:       v.Id,
		Title:    v.Title,
		LinkURL:  v.LinkUrl,
		ImageURL: v.ImageUrl,
		Enabled:  int32(v.Enabled),
		SortNum:  int32(v.SortNum),
	}
}

func (a *advertisementService) parseSwiperDto(dto *ad.AdDto) *proto.SSwiperAdData {
	images := dto.Data.(ad.SwiperAd)
	arr := make([]*proto.SImageAdData, len(images))
	for i, v := range images {
		arr[i] = a.parseSingleImageDto(v)
	}
	return &proto.SSwiperAdData{Images: arr}
}

// 更新广告数据
func (a *advertisementService) updateAdData(ia ad.IAd, data *proto.SAdvertisementDto) error {
	if data == nil {
		return nil
	}
	if ia.Type() == ad.TypeText {
		g := ia.(ad.IHyperLinkAd)
		v := a.parseHyperLinkAd(data.Text)
		err := g.SetData(v)
		if err == nil {
			_, err = ia.Save()
		}
		return err
	}
	if ia.Type() == ad.TypeImage {
		g := ia.(ad.IImageAd)
		err := g.SetData(a.parseAdImage(data.Image))
		if err == nil {
			_, err = ia.Save()
		}
		return err
	}
	if ia.Type() == ad.TypeSwiper {
		return errors.New("please use SaveSwiperImage to update")
	}
	return nil
}
