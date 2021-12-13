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
	"github.com/ixre/gof/types"
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

func (a *advertisementService) GetAdGroupById(_ context.Context, id *proto.Int64) (*proto.SAdGroup, error) {
	ig := a._rep.GetAdManager().GetAdGroup(id.Value)
	if ig != nil {
		return a.parseAdGroupDto(ig.GetValue()), nil
	}
	return nil, nil
}

func (a *advertisementService) DelAdGroup(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	err := a._rep.GetAdManager().DelAdGroup(id.Value)
	return a.error(err), nil
}

func (a *advertisementService) SaveAdGroup(_ context.Context, group *proto.SAdGroup) (*proto.Result, error) {
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

func (a *advertisementService) GetGroups(_ context.Context, empty *proto.Empty) (*proto.AdGroupListResponse, error) {
	m := a._rep.GetAdManager()
	list := m.GetAdGroups()
	arr := make([]*proto.SAdGroup, 0)
	for _, v := range list {
		arr = append(arr, a.parseAdGroupDto(v.GetValue()))
	}
	return &proto.AdGroupListResponse{Value: arr}, nil
}

func (a *advertisementService) GetGroupsV2(_ context.Context, empty *proto.Empty) (*proto.AdGroupResponse, error) {
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
		Value: make([]*proto.SAd, len(ret)),
	}
	for i, v := range ret {
		rsp.Value[i] = &proto.SAd{
			Id:   v.Id,
			Name: v.Name,
		}
	}
	return rsp, nil
}

// 设置广告位的默认广告
func (a *advertisementService) SetDefaultAd(_ context.Context, r *proto.SetDefaultAdRequest) (*proto.Result, error) {
	ig := a._rep.GetAdManager().GetAdGroup(r.GroupId)
	err := ig.SetDefault(r.PosId, r.AdId)
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
func (a *advertisementService) GetAdvertisement(_ context.Context, r *proto.AdIdRequest) (*proto.SAd, error) {
	ig := a.getUserAd(r.AdUserId).GetById(r.AdId)
	if ig != nil {
		return a.parseAdDto(ig.GetValue()), nil
	}
	return nil, nil
}

// 获取广告及广告数据, 用于展示关高
func (a *advertisementService) GetAdAndDataByKey(_ context.Context, r *proto.AdKeyRequest) (*proto.SAdvertisementDto, error) {
	//if adv := a.getUserAd(adUserId).GetByPositionKey(key); adv != nil {
	//	switch adv.AdType() {
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

// GetAdvertisementDto 获取广告数据传输对象
func (a *advertisementService) GetAdvertisementDto(_ context.Context, r *proto.AdIdRequest) (*proto.SAdvertisementDto, error) {
	ua := a.getUserAd(r.AdUserId).GetById(r.AdId)
	if ua == nil {
		return nil, nil
	}
	dto :=  ua.Dto()
	ret := &proto.SAdvertisementDto{Id: dto.Id,AdType: int32(dto.AdType),}
	switch dto.AdType {
	case ad.TypeText:
		ret.Text = a.parseTextDto(dto)
	case ad.TypeImage:
		ret.Image = a.parseImageDto(dto)
	case ad.TypeGallery:
		ret.Swiper = a.parseSwiperDto(dto)
	default:
		panic("not support ad type")
	}
	return ret,nil
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
		_, err = adv.Save()
	}
	return a.error(err), nil
}

func (a *advertisementService) DeleteAd(_ context.Context, r *proto.AdIdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	err := a.getUserAd(r.AdUserId).DeleteAd(r.AdId)
	return a.error(err), nil
}

// 保存图片广告
func (a *advertisementService) SaveHyperLinkAd(_ context.Context, r *proto.SaveLinkAdRequest) (*proto.Result, error) {
	defer a.cleanCache(r.AdUserId)
	pa := a.getUserAd(r.AdUserId)
	var adv = pa.GetById(r.AdId)
	var err error
	if adv == nil {
		err = ad.ErrNoSuchAd
	} else {
		if adv.Type() == ad.TypeText {
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
func (a *advertisementService) SaveImagOfAd(_ context.Context, r *proto.SaveImageAdRequest) (*proto.Result, error) {
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
func (a *advertisementService) SaveImage(adUserId int64, adId int64, v *ad.Image) (int64, error) {
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
func (a *advertisementService) GetValueAdImage(_ context.Context, r *proto.ImageAdIdRequest) (*proto.SImageAdData, error) {
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
func (a *advertisementService) DelAdImage(_ context.Context, r *proto.ImageAdIdRequest) (*proto.Result, error) {
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

func (a *advertisementService) parseAdGroup(v *proto.SAdGroup) *ad.AdGroup {
	return &ad.AdGroup{
		ID:      v.Id,
		Name:    v.Name,
		Flag:    int(v.Flag),
		Opened:  types.ElseInt(v.Opened, 1, 0),
		Enabled: types.ElseInt(v.Enabled, 1, 0),
	}
}

func (a *advertisementService) parseAdGroupDto(v ad.AdGroup) *proto.SAdGroup {
	return &proto.SAdGroup{
		Id:      v.ID,
		Name:    v.Name,
		Flag:    int32(v.Flag),
		Opened:  v.Opened == 1,
		Enabled: v.Enabled == 1,
	}
}

func (a *advertisementService) parseAdPositionDto(v *ad.Position) *proto.SAdPosition {
	return &proto.SAdPosition{
		Id:      v.Id,
		GroupId: v.GroupId,
		Key:     v.Key,
		Name:    v.Name,
		Flag:    int32(v.Flag),
		//TypeLimit: int32(v.TypeLimit),
		Opened:    int32(v.Opened),
		Enabled:   int32(v.Enabled),
		PutAid:    v.PutAdId,
		GroupName: v.GroupName,
	}
}

func (a *advertisementService) parseAdPosition(v *proto.SAdPosition) *ad.Position {
	return &ad.Position{
		Id:      v.Id,
		GroupId: v.GroupId,
		Key:     v.Key,
		Name:    v.Name,
		Flag:    int(v.Flag),
		//TypeLimit: int(v.TypeLimit),
		Opened:    int(v.Opened),
		Enabled:   int(v.Enabled),
		PutAdId:   v.PutAid,
		GroupName: v.GroupName,
	}
}

func (a *advertisementService) parseAdDto(v *ad.Ad) *proto.SAd {
	return &proto.SAd{
		Id:         v.Id,
		UserId:     v.UserId,
		Name:       v.Name,
		AdType:       int32(v.AdType),
		ShowTimes:  int32(v.ShowTimes),
		ClickTimes: int32(v.ClickTimes),
		ShowDays:   int32(v.ShowDays),
		UpdateTime: v.UpdateTime,
	}
}

func (a *advertisementService) parseAd(v *proto.SAd) *ad.Ad {
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
		Enabled:  v.Enabled == 1,
	}
}

func (a *advertisementService) parseAdImage(v *proto.SImageAdData) *ad.Image {
	return &ad.Image{
		Id:       v.Id,
		Title:    v.Title,
		LinkUrl:  v.LinkURL,
		ImageUrl: v.ImageURL,
		SortNum:  int(v.SortNum),
		Enabled:  types.ElseInt(v.Enabled, 1, 0),
	}
}

func (a *advertisementService) parseTextDto(dto *ad.AdDto) *proto.STextAdData {
	v := dto.Data.(*ad.HyperLink)
	return &proto.STextAdData{
		Id:                   v.Id,
		Title:                v.Title,
		LinkURL:              v.LinkUrl,
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
		Enabled:  v.Enabled == 1,
		SortNum:  int32(v.SortNum),
	}
}

func (a *advertisementService) parseSwiperDto(dto *ad.AdDto) *proto.SSwiperAdData {
	images := dto.Data.(ad.ValueGallery)
	arr := make([]*proto.SImageAdData,len(images))
	for i,v := range images{
		arr[i] = a.parseSingleImageDto(v)
	}
	return &proto.SSwiperAdData{Images: arr}
}
