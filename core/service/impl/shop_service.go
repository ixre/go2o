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

func (si *shopServiceImpl) SaveShop(_ context.Context, sShop *proto.SShop) (*proto.Result, error) {
	panic("implement me")
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
			sp = mgr.GetShop(int(store.Id))
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

func (si *shopServiceImpl) DeleteStore(_ context.Context, id *proto.ShopId) (*proto.Result, error) {
	panic("implement me")
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

func (si *shopServiceImpl) GetShop(_ context.Context, shopId *proto.Int64) (*proto.SShop, error) {
	sp := si.shopRepo.GetOnlineShop(int(shopId.Value))
	if sp != nil {
		return si.parseShop(sp), nil
	}
	return nil, nil
}

func (si *shopServiceImpl) GetVendorShop(_ context.Context, vendorId *proto.Int64) (*proto.SShop, error) {
	sp := si.shopRepo.GetOnlineShopOfMerchant(int(vendorId.Value))
	if sp != nil {
		return si.parseShop(sp), nil
	}
	return nil, nil
}

// 根据主机头获取店铺编号
func (si *shopServiceImpl) QueryShopByHost(_ context.Context, host *proto.String) (*proto.Int64, error) {
	_, shopId := si.query.QueryShopIdByHost(host.Value)
	return &proto.Int64{Value: shopId}, nil
}

// 获取门店
func (si *shopServiceImpl) GetStore(_ context.Context, storeId *proto.Int64) (*proto.SStore, error) {
	panic("返回门店")
	//mch := si.mchRepo.GetMerchant(int(storeId))
	//if mch != nil {
	//	shop := mch.ShopManager().GetOnlineShop()
	//	if shop != nil {
	//		return parser.ParseOnlineShop(shop), nil
	//	}
	//}
	//return nil, nil
}

func (si *shopServiceImpl) GetStoreById(ctx context.Context, shopId *proto.Int64) (*proto.SStore, error) {
	vendorId := si.query.GetMerchantId(shopId.Value)
	return si.GetStore(ctx, &proto.Int64{Value: vendorId})
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
	sp := mch.ShopManager().GetShop(int(shopId))
	if sp != nil {
		return sp.Data()
	}
	return nil
}

func (si *shopServiceImpl) GetShopValueById(mchId, shopId int64) *shop.Shop {
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		v := mch.ShopManager().GetShop(int(shopId)).GetValue()
		return &v
	}
	return nil
}

// 保存线上商店
func (si *shopServiceImpl) SaveStore(s *proto.SStore) error {
	mch := si.mchRepo.GetMerchant(int(s.MerchantId))
	if mch != nil {
		v, v1 := si.parse2OnlineShop(s)
		mgr := mch.ShopManager()
		sp := mgr.GetOnlineShop()
		// 创建商店
		if sp == nil {
			sp = mgr.CreateShop(v)
		}
		err := sp.SetValue(v)
		if err == nil {
			ofs := sp.(shop.IOnlineShop)
			err = ofs.SetShopValue(v1)
			if err == nil {
				err = sp.Save()
			}
		}
		return err
	}
	return merchant.ErrNoSuchMerchant
}

func (si *shopServiceImpl) SaveShop_(mchId int64, v *shop.Shop) (int64, error) {
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		var shop shop.IShop
		if v.Id > 0 {
			shop = mch.ShopManager().GetShop(int(v.Id))
			if shop == nil {
				return 0, errors.New("门店不存在")
			}
		} else {
			shop = mch.ShopManager().CreateShop(v)
		}
		err := shop.SetValue(v)
		if err != nil {
			return v.Id, err
		}
		err = shop.Save()
		return int64(shop.GetDomainId()), err
	}
	return 0, merchant.ErrNoSuchMerchant
}

func (si *shopServiceImpl) DeleteShop(mchId, shopId int32) error {
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ShopManager().DeleteShop(shopId)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取线上商城配置
func (si *shopServiceImpl) GetOnlineShopConf(shopId int64) *shop.OnlineShop {
	mchId := si.getMerchantId(shopId)
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		s := mch.ShopManager().GetShop(int(shopId))
		if s == nil {
			v := s.(shop.IOnlineShop).GetShopValue()
			return &v
		}
	}
	return nil
}

// 获取商城
func (si *shopServiceImpl) GetOnlineShops(vendorId int64) []*shop.Shop {
	mch := si.mchRepo.GetMerchant(int(vendorId))
	shops := mch.ShopManager().GetShops()
	sv := make([]*shop.Shop, 0)
	for _, v := range shops {
		if v.Type() == shop.TypeOnlineShop {
			vv := v.GetValue()
			sv = append(sv, &vv)
		}
	}
	return sv
}

func (si *shopServiceImpl) parseShop(sp *shop.OnlineShop) *proto.SShop {
	return &proto.SShop{
		Id:         sp.Id,
		MerchantId: sp.VendorId,
		ShopName:   sp.ShopName,
		Config: &proto.SShopConfig{
			Logo:  sp.Logo,
			Host:  sp.Host,
			Alias: sp.Alias,
			Tel:   sp.Tel,
		},
		ShopTitle:  sp.ShopTitle,
		ShopNotice: sp.ShopNotice,
	}
}

func (si *shopServiceImpl) parse2OnlineShop(s *proto.SStore) (*shop.Shop, *shop.OnlineShop) {
	sv := &shop.Shop{
		Id:           s.Id,
		Name:         s.Name,
		VendorId:     s.MerchantId,
		ShopType:     shop.TypeOnlineShop,
		State:        s.State,
		OpeningState: s.OpeningState,
	}
	ov := &shop.OnlineShop{}
	ov.Id = s.Id
	ov.Addr = "" //todo:???
	ov.Tel = s.StorePhone
	ov.Logo = s.Logo
	ov.ShopNotice = s.StoreNotice
	ov.ShopTitle = s.StoreTitle
	return sv, ov
}

func (si *shopServiceImpl) parseOfflineShop(r *proto.SStore) (*shop.Shop, *shop.OfflineShop) {
	return &shop.Shop{
			Id:           r.Id,
			VendorId:     r.MerchantId,
			ShopType:     shop.TypeOfflineShop,
			Name:         r.Name,
			State:        r.State,
			OpeningState: r.OpeningState,
			SortNum:      r.SortNumber,
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
