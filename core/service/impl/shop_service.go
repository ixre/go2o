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
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/query"
	"go2o/core/service/proto"
)

var _ proto.ShopServiceServer = new(shopServiceImpl)

type shopServiceImpl struct {
	repo     shop.IShopRepo
	mchRepo  merchant.IMerchantRepo
	shopRepo shop.IShopRepo
	query    *query.ShopQuery
	serviceUtil
}

func NewShopService(rep shop.IShopRepo, mchRepo merchant.IMerchantRepo,
	shopRepo shop.IShopRepo, query *query.ShopQuery) *shopServiceImpl {
	return &shopServiceImpl{
		repo:     rep,
		mchRepo:  mchRepo,
		shopRepo: shopRepo,
		query:    query,
	}
}

// 保存门店
func (si *shopServiceImpl) SaveOfflineShop(_ context.Context, r *proto.SStore) (*proto.Result, error) {
	mch := si.mchRepo.GetMerchant(int(r.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		mgr := mch.ShopManager()
		store, v := si.parseOfflineShop(r)
		var sp shop.IShop
		if store.Id > 0 {
			// 保存商店
			sp = mgr.GetStore(int(store.Id))
		} else {
			//创建商店
			sp = mgr.CreateShop(store)
		}
		err = sp.SetValue(store)
		if err == nil {
			ofs := sp.(shop.IOfflineShop)
			err = ofs.SetShopValue(v)
			if err == nil {
				err = sp.Save()
			}
		}
	}
	return si.error(err), nil
}

func (si *shopServiceImpl) DeleteStore(_ context.Context, id *proto.StoreId) (*proto.Result, error) {
	panic("implement me")
}

func (si *shopServiceImpl) GetShop(_ context.Context, shopId *proto.ShopId) (*proto.SShop, error) {
	sp := si.shopRepo.GetShop(shopId.Value)
	if sp != nil {
		iop := sp.(shop.IOnlineShop)
		iv := iop.GetShopValue()
		ret := si.parseShopDto(iv)
		ret.ShopTitle = ret.ShopName
		ret.Host = iv.Host
		ret.Logo = iv.Logo
		ret.Alias = iv.Alias
		ret.Telephone = iv.Tel
		return ret, nil
	}
	return nil, nil
}

// 检查商户是否开通店铺
func (si *shopServiceImpl) CheckMerchantShopState(_ context.Context, id *proto.MerchantId) (*proto.CheckShopResponse, error) {
	sp := si.shopRepo.GetOnlineShopOfMerchant(int(id.Value))
	ret := &proto.CheckShopResponse{}
	if sp != nil {
		ret.Status = 1
		ret.Remark = "已开通"
		ret.ShopId = int64(sp.GetDomainId())
	} else {
		//todo: 返回审核中状态
	}
	return ret, nil
}

// 根据主机头获取店铺编号
func (si *shopServiceImpl) QueryShopByHost(_ context.Context, host *proto.String) (*proto.Int64, error) {
	_, shopId := si.query.QueryShopIdByHost(host.Value)
	return &proto.Int64{Value: shopId}, nil
}

// 获取门店
func (si *shopServiceImpl) GetStore(_ context.Context, storeId *proto.StoreId) (*proto.SStore, error) {
	sp := si.shopRepo.GetStore(storeId.Value)
	if sp != nil {
		v := sp.GetValue()
		ifv := sp.(shop.IOfflineShop)
		iv := ifv.GetShopValue()
		ret := &proto.SStore{
			Id:            storeId.Value,
			MerchantId:    v.VendorId,
			Name:          v.Name,
			State:         v.State,
			OpeningState:  v.OpeningState,
			StorePhone:    iv.Tel,
			StoreNotice:   "",
			Province:      iv.Province,
			City:          iv.City,
			District:      iv.District,
			Address:       "",
			DetailAddress: iv.Address,
			Lat:           float64(iv.Lat),
			Lng:           float64(iv.Lng),
			CoverRadius:   int32(iv.CoverRadius),
			SortNum:       v.SortNum,
		}
		return ret, nil
	}
	return nil, nil
}

// 打开或关闭商店
func (si *shopServiceImpl) TurnShop(_ context.Context, r *proto.TurnShopRequest) (*proto.Result, error) {
	var err error
	sp := si.repo.GetShop(r.ShopId)
	if sp == nil {
		err = shop.ErrNoSuchShop
	} else {
		if r.On {
			err = sp.TurnOn()
		} else {
			err = sp.TurnOff(r.Reason)
		}
	}
	return si.result(err), nil
}

// 设置商店是否营业
func (si *shopServiceImpl) OpenShop(_ context.Context, shopId int32, on bool, reason string) (*proto.Result, error) {
	var err error
	sp := si.repo.GetShop(int64(shopId))
	if sp == nil {
		err = shop.ErrNoSuchShop
	} else {
		if on {
			err = sp.Opening()
		} else {
			err = sp.Pause()
		}
	}
	return si.result(err), nil
}

func (si *shopServiceImpl) getMerchantId(shopId int64) int64 {
	return si.query.GetMerchantId(shopId)
}

//todo : remove
func (si *shopServiceImpl) GetMerchantId(shopId int64) int64 {
	return si.query.GetMerchantId(shopId)
}

// 获取商店的数据
func (si *shopServiceImpl) GetShopData(mchId, shopId int64) *shop.ComplexShop {
	mch := si.mchRepo.GetMerchant(int(mchId))
	sp := mch.ShopManager().GetStore(int(shopId))
	if sp != nil {
		return sp.Data()
	}
	return nil
}

func (si *shopServiceImpl) GetShopValueById(mchId, shopId int64) *shop.Shop {
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		v := mch.ShopManager().GetStore(int(shopId)).GetValue()
		return &v
	}
	return nil
}

// 保存线上商店
func (si *shopServiceImpl) SaveShop(_ context.Context, s *proto.SShop) (*proto.Result, error) {
	mch := si.mchRepo.GetMerchant(int(s.MerchantId))
	var err error
	if mch == nil {
		err = merchant.ErrNoSuchMerchant
	} else {
		v1 := si.parse2OnlineShop(s)
		mgr := mch.ShopManager()
		sp := mgr.GetOnlineShop()
		if sp == nil {
			err = merchant.ErrNoSuchShop
		} else {
			ofs := sp.(shop.IOnlineShop)
			err := ofs.SetShopValue(v1)
			if err == nil {
				err = sp.Save()
			}
		}
	}
	return si.error(err), nil
}

func (si *shopServiceImpl) DeleteShop(mchId, shopId int32) error {
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ShopManager().DeleteShop(shopId)
	}
	return merchant.ErrNoSuchMerchant
}

func (si *shopServiceImpl) parse2OnlineShop(s *proto.SShop) *shop.OnlineShop {
	return &shop.OnlineShop{
		Id:         s.Id,
		VendorId:   s.MerchantId,
		ShopName:   s.ShopName,
		Logo:       s.Logo,
		Host:       s.Host,
		Tel:        s.Telephone,
		ShopTitle:  s.ShopTitle,
		ShopNotice: s.ShopNotice,
		State:      int16(s.State),
	}
}

func (si *shopServiceImpl) parseOfflineShop(r *proto.SStore) (*shop.Shop, *shop.OfflineShop) {
	return &shop.Shop{
			Id:           r.Id,
			VendorId:     r.MerchantId,
			ShopType:     shop.TypeOfflineShop,
			Name:         r.Name,
			State:        r.State,
			OpeningState: r.OpeningState,
			SortNum:      r.SortNum,
		}, &shop.OfflineShop{
			ShopId:      int(r.Id),
			Tel:         r.StorePhone,
			Province:    r.Province,
			City:        r.City,
			District:    r.District,
			Address:     r.DetailAddress,
			Lng:         float32(r.Lng),
			Lat:         float32(r.Lat),
			CoverRadius: int(r.CoverRadius),
		}
}

func (si *shopServiceImpl) parseShopDto(v shop.OnlineShop) *proto.SShop {
	return &proto.SShop{
		Id:         v.Id,
		MerchantId: v.VendorId,
		ShopName:   v.ShopName,
		ShopTitle:  v.ShopTitle,
		ShopNotice: v.ShopNotice,
		Flag:       int32(v.Flag),
		State:      int32(v.State),
	}
}
