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

	"github.com/ixre/go2o/core/domain/interface/ad"
	"github.com/ixre/go2o/core/repos"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/storage"
)

var _ proto.AdvertisementServiceServer = new(advertisementService)

type advertisementService struct {
	_rep    ad.IAdRepo
	storage storage.Interface
	//_query     *shopQuery.ContentQuery
	serviceUtil
	proto.UnimplementedAdvertisementServiceServer
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
	return nil, ad.ErrNoSuchAdPosition
}

func (a *advertisementService) SaveAdPosition(_ context.Context, r *proto.SAdPosition) (*proto.TxResult, error) {
	v := a.parseAdPosition(r)
	var ap ad.IAdPosition
	if r.Id > 0 {
		ap = a._rep.GetPosition(r.Id)
		if ap == nil {
			return a.errorV2(ad.ErrNoSuchAdPosition), nil
		}
	} else {
		ap = a._rep.CreateAdPosition(v)
	}
	err := ap.SetValue(v)
	if err == nil {
		err = ap.Save()
	}
	return a.errorV2(err), nil
}

func (a *advertisementService) DeleteAdPosition(_ context.Context, id *proto.AdPositionId) (*proto.TxResult, error) {
	err := a._rep.DeleteAdPosition(id.PositionId)
	return a.errorV2(err), nil
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
			AdId: int64(v.Id),
			Name: v.Name,
		}
	}
	return rsp, nil
}

// PutDefaultAd 设置广告位的默认广告
func (a *advertisementService) PutDefaultAd(_ context.Context, r *proto.SetDefaultAdRequest) (*proto.TxResult, error) {
	ig := a._rep.GetAdManager().GetPosition(r.PositionId)
	if ig == nil {
		return a.errorV2(ad.ErrNoSuchAdPosition), nil
	}
	err := ig.PutAd(r.AdId)
	return a.errorV2(err), nil
}

// 用户投放广告
func (a *advertisementService) SetUserAd(_ context.Context, r *proto.SetUserAdRequest) (*proto.TxResult, error) {
	defer a.cleanCache(r.AdUserId)
	ua := a._rep.GetAdManager().GetUserAd(int(r.AdUserId))
	err := ua.SetAd(int(r.PosId), int(r.AdId))
	return a.errorV2(err), nil
}

// 获取广告
func (a *advertisementService) GetAdvertisement(_ context.Context, r *proto.AdIdRequest) (*proto.SAd, error) {
	ia := a.getUserAd(r.AdUserId).GetById(int(r.AdId))
	if ia != nil {
		ret := a.parseAdDto(ia.GetValue())
		if r.ReturnData {
			//ret.Data = a.getAdvertisementDto(ia)
		}
		return ret, nil
	}
	return nil, ad.ErrNoSuchAd
}

// 获取广告数据传输对象
func (a *advertisementService) getAdvertisementPackage(ia ad.IAdAggregateRoot) *proto.SAdvertisementPackage {
	dto := ia.Dto()
	ret := &proto.SAdvertisementPackage{
		Id:    int64(dto.Id),
		Type:  int32(dto.AdType),
		Media: make([]*proto.SAdvertisementMedia, 0),
	}
	switch dto.AdType {
	case ad.TypeText:
		v := dto.Data.(*ad.Data)
		if v != nil {
			ret.Media = append(ret.Media, &proto.SAdvertisementMedia{
				Title:    v.Title,
				Cmd:      "LINK",
				LinkUrl:  v.LinkUrl,
				MediaUrl: "",
			})
		}
	case ad.TypeImage:
		v := dto.Data.(*ad.Data)
		if v != nil {
			ret.Media = append(ret.Media, &proto.SAdvertisementMedia{
				Title:    v.Title,
				Cmd:      "LINK",
				LinkUrl:  v.LinkUrl,
				MediaUrl: v.ImageUrl,
			})
		}
	case ad.TypeSwiper:
		images := dto.Data.(ad.SwiperAd)
		for _, v := range images {
			ret.Media = append(ret.Media, &proto.SAdvertisementMedia{
				Title:    v.Title,
				Cmd:      "LINK",
				LinkUrl:  v.LinkUrl,
				MediaUrl: v.ImageUrl,
			})
		}
	default:
		panic("not support ad type")
	}
	return ret
}

// // 获取广告数据传输对象
// func (a *advertisementService) getAdvertisementDto(ia ad.IAdAggregateRoot) *proto.SAdvertisementDto {
// 	dto := ia.Dto()
// 	ret := &proto.SAdvertisementDto{Id: dto.Id, Type: int32(dto.AdType)}
// 	switch dto.AdType {
// 	case ad.TypeText:
// 		ret.Text = a.parseTextDto(dto)
// 	case ad.TypeImage:
// 		ret.Image = a.parseImageDto(dto)
// 	case ad.TypeSwiper:
// 		ret.Swiper = a.parseSwiperDto(dto)
// 	default:
// 		panic("not support ad type")
// 	}
// 	return ret
// }

func (a *advertisementService) QueryAdvertisementData(_ context.Context, r *proto.QueryAdvertisementDataRequest) (*proto.QueryAdvertisementDataResponse, error) {
	iu := a.getUserAd(r.AdUserId)
	var list = iu.QueryAdvertisement(r.Keys)
	arr := make(map[string]*proto.SAdvertisementPackage, len(r.Keys))
	for _, k := range r.Keys {
		v := list[k]
		if v == nil {
			arr[k] = nil
		} else {
			arr[k] = a.getAdvertisementPackage(v)
		}
	}
	return &proto.QueryAdvertisementDataResponse{
		UserId: r.AdUserId,
		Value:  arr,
	}, nil
}

// 保存广告,更新时不允许修改类型
func (a *advertisementService) SaveAd(_ context.Context, req *proto.SAd) (*proto.TxResult, error) {
	defer a.cleanCache(req.UserId)
	pa := a.getUserAd(req.UserId)
	var adv ad.IAdAggregateRoot
	v := &ad.Ad{
		UserId: int(req.UserId),
		Name:   req.Name,
		TypeId: int(req.AdType),
	}
	if req.AdId > 0 {
		adv = pa.GetById(v.Id)
		if adv == nil {
			return a.errorV2(ad.ErrNoSuchAd), nil
		}
	} else {
		adv = pa.CreateAd(v)
	}
	// 保存广告
	err := adv.SetValue(v)
	if err == nil {
		_, err = adv.Save()
	}
	if err != nil {
		return a.errorV2(err), nil
	}
	// 保存广告数据
	err = a.updateAdData(adv, req.Data)
	return a.errorV2(err), nil
}

func (a *advertisementService) DeleteAd(_ context.Context, r *proto.AdIdRequest) (*proto.TxResult, error) {
	defer a.cleanCache(r.AdUserId)
	err := a.getUserAd(r.AdUserId).DeleteAd(r.AdId)
	return a.errorV2(err), nil
}

// 保存图片广告
func (a *advertisementService) SaveSwiperAdImage(_ context.Context, r *proto.SaveSwiperImageRequest) (*proto.Result, error) {
	ia := a.getUserAd(r.AdUserId).GetById(int(r.AdId))
	var err error
	if ia == nil {
		err = errors.New("no such ad image")
	} else {
		if ia.Type() == ad.TypeSwiper {

		}
	}
	return a.error(err), nil
}

func NewAdvertisementService(rep ad.IAdRepo, storage storage.Interface) proto.AdvertisementServiceServer {
	return &advertisementService{
		_rep:    rep,
		storage: storage,
	}
}

func (a *advertisementService) getUserAd(adUserId int64) ad.IUserAd {
	return a._rep.GetAdManager().GetUserAd(int(adUserId))
}

func (a *advertisementService) GetPositionById(adPosId int64) *ad.Position {
	return a._rep.GetAdPositionById(adPosId)
}

func (a *advertisementService) cleanCache(adUserId int64) error {
	return repos.PrefixDel(a.storage, fmt.Sprintf("go2o:repo:ad:%d:*", adUserId))
}

func (a *advertisementService) parseAdPositionDto(v *ad.Position) *proto.SAdPosition {
	return &proto.SAdPosition{
		Id:   int64(v.Id),
		Key:  v.Key,
		Name: v.Name,
		Flag: int32(v.Flag),
		//TypeLimit: int32(v.TypeLimit),
		Opened:    int32(v.Opened),
		Enabled:   int32(v.Enabled),
		PutAid:    int64(v.PutAid),
		GroupName: v.GroupName,
	}
}

func (a *advertisementService) parseAdPosition(v *proto.SAdPosition) *ad.Position {
	return &ad.Position{
		Id:   int(v.Id),
		Key:  v.Key,
		Name: v.Name,
		Flag: int(v.Flag),
		//TypeLimit: int(v.TypeLimit),
		Opened:    int(v.Opened),
		Enabled:   int(v.Enabled),
		PutAid:    int(v.PutAid),
		GroupName: v.GroupName,
	}
}

func (a *advertisementService) parseAdDto(v *ad.Ad) *proto.SAd {
	return &proto.SAd{
		AdId:   int64(v.Id),
		UserId: int64(v.UserId),
		Name:   v.Name,
		AdType: int32(v.TypeId),
	}
}

func (a *advertisementService) parseAdImage(v *proto.SAdData) *ad.Data {
	return &ad.Data{
		Id:       int(v.Id),
		Title:    v.Title,
		LinkUrl:  v.LinkUrl,
		ImageUrl: v.ImageUrl,
		SortNum:  int(v.SortNum),
		Enabled:  int(v.Enabled),
	}
}

func (a *advertisementService) parseImageDto(dto *ad.AdDto) *proto.SAdData {
	v := dto.Data.(*ad.Data)
	if v == nil {
		return &proto.SAdData{}
	}
	return a.parseSingleImageDto(v)
}

func (a *advertisementService) parseSingleImageDto(v *ad.Data) *proto.SAdData {
	return &proto.SAdData{
		Id:       int64(v.Id),
		Title:    v.Title,
		LinkUrl:  v.LinkUrl,
		ImageUrl: v.ImageUrl,
		Enabled:  int32(v.Enabled),
		SortNum:  int32(v.SortNum),
	}
}

// 更新广告数据
func (a *advertisementService) updateAdData(ia ad.IAdAggregateRoot, data []*proto.SAdData) error {
	if data == nil {
		return nil
	}
	if ia.Type() == ad.TypeText {
		g := ia.(ad.IHyperLinkAd)
		err := g.SetData(a.parseAdImage(data[0]))
		if err == nil {
			_, err = ia.Save()
		}
		return err
	}
	if ia.Type() == ad.TypeImage {
		g := ia.(ad.IImageAd)
		err := g.SetData(a.parseAdImage(data[0]))
		if err == nil {
			_, err = ia.Save()
		}
		return err
	}
	if ia.Type() == ad.TypeSwiper {
		gad := ia.(ad.IGalleryAd)
		arr := make([]*ad.Data, len(data))
		for i, v := range data {
			arr[i] = a.parseAdImage(v)
		}
		err := gad.SaveImage(arr)
		if err == nil {
			_, err = ia.Save()
		}
		return err
	}
	return nil
}
