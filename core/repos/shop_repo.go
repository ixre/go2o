/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 17:16
 * description :
 * history :
 */

package repos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	shopImpl "github.com/ixre/go2o/core/domain/merchant/shop"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"log"
)

var _ shop.IShopRepo = new(shopRepo)

type shopRepo struct {
	db.Connector
	valueRepo    valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
	storage      storage.Interface
	o            orm.Orm
}

func (s *shopRepo) GetShopIdByAlias(alias string) int64 {
	e := shop.OnlineShop{}
	if s.o.GetBy(&e, "alias = $1 LIMIT 1", alias) != nil {
		return 0
	}
	return e.Id
}

// CreateShop 创建店铺
func (s *shopRepo) CreateShop(v *shop.OnlineShop) shop.IShop {
	return shopImpl.NewShop(v, s, s.valueRepo, s.registryRepo)
}

// CheckShopExists 检查商户商城是否存在(创建)
func (s *shopRepo) CheckShopExists(vendorId int64) bool {
	num := 0
	s.Connector.ExecScalar(`SELECT COUNT(1) FROM mch_online_shop WHERE
	    vendor_id = $1 LIMIT 1`, &num, vendorId)
	return num > 0
}

func (s *shopRepo) ShopCount(vendorId int64, shopType int32) int {
	num := 0
	s.Connector.ExecScalar(`SELECT COUNT(1) FROM mch_shop WHERE
	    vendor_id= $1 AND shop_type = $2`, &num, vendorId, shopType)
	return num
}

func NewShopRepo(o orm.Orm, storage storage.Interface,
	valueRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) shop.IShopRepo {
	return &shopRepo{
		Connector:    o.Connector(),
		o:            o,
		valueRepo:    valueRepo,
		storage:      storage,
		registryRepo: registryRepo,
	}
}

// 获取商店
func (s *shopRepo) GetShop(shopId int64) shop.IShop {
	e := shop.OnlineShop{}
	if s.o.Get(shopId, &e) != nil {
		return nil
	}
	return s.CreateShop(&e)
}

// 获取门店
func (s *shopRepo) GetStore(storeId int64) shop.IShop {
	v := s.GetValueShop(storeId)
	return shopImpl.NewStore(v, s, s.valueRepo, s.registryRepo)
}

// 商店别名是否存在
func (s *shopRepo) ShopAliasExists(alias string, shopId int) bool {
	id := 0
	s.Connector.ExecScalar(`SELECT id FROM mch_online_shop WHERE
		alias= $1 AND shop_id<> $2 LIMIT 1`, &id, alias, shopId)
	return id > 0
}

// 保存线上商店
func (s *shopRepo) SaveOnlineShop(v *shop.OnlineShop) (int64, error) {
	id, err := orm.Save(s.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchOnlineShop")
	}
	return int64(id), err
}

// 获取线下商店
func (s *shopRepo) GetOfflineShop(shopId int) *shop.OfflineShop {
	e := shop.OfflineShop{}
	if s.o.Get(shopId, &e) != nil {
		return nil
	}
	return &e
}

// 保存线下商店
func (s *shopRepo) SaveOfflineShop(v *shop.OfflineShop, create bool) error {
	var err error
	if create {
		_, _, err = s.o.Save(nil, v)
	} else {
		_, _, err = s.o.Save(v.ShopId, v)
	}
	return err
}

// 保存API信息
func (s *shopRepo) SaveApiInfo(v *merchant.ApiInfo) error {
	_, err := orm.Save(s.o, v, int(v.MerchantId))
	return err
}

// 获取API信息
func (s *shopRepo) GetApiInfo(mchId int) *merchant.ApiInfo {
	var d = new(merchant.ApiInfo)
	if err := s.o.Get(mchId, d); err == nil {
		return d
	}
	return nil
}

func (s *shopRepo) SaveShop(v *shop.Shop) (int64, error) {
	id, err := orm.I32(orm.Save(s.o, v, int(v.Id)))
	if err == nil {
		s.delCache(v.VendorId)
	}
	return int64(id), err
}

func (s *shopRepo) GetValueShop(shopId int64) *shop.Shop {
	v := &shop.Shop{}
	err := s.o.Get(shopId, v)
	if err == nil {
		return v
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchShop")
	}
	return nil
}

func (s *shopRepo) delCache(mchId int64) {
	PrefixDel(s.storage, fmt.Sprintf("go2o:repo:shop:%d:*", mchId))
}

func (s *shopRepo) getShopCacheKey(mchId int64) string {
	return fmt.Sprintf("go2o:repo:shop:%d:shops", mchId)
}

func (s *shopRepo) GetOnlineShopOfMerchant(vendorId int) *shop.OnlineShop {
	v := shop.OnlineShop{}
	err := s.o.GetBy(&v, "vendor_id= $1 LIMIT 1", vendorId)
	if err == nil {
		return &v
	}
	return nil
}

/**/
func (s *shopRepo) GetShopId(mchId int64) []shop.Shop {
	shops := make([]shop.Shop, 0)
	key := s.getShopCacheKey(mchId)
	jsonStr, err := s.storage.GetString(key)
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &shops)
	}
	if err != nil {
		err = s.o.SelectByQuery(&shops,
			"SELECT * FROM mch_shop WHERE vendor_id= $1", mchId)
		if err != nil {
			handleError(err)
			return nil
		} else {
			b, _ := json.Marshal(shops)
			s.storage.Set(key, string(b))
		}
	}
	return shops
}

func (s *shopRepo) deleteShop(mchId, shopId int64) error {
	_, err := s.o.Delete(shop.Shop{},
		"vendor_id= $1 AND id= $2", mchId, shopId)
	s.delCache(mchId)
	return err
}

// 删除线上商店
func (s *shopRepo) DeleteOnlineShop(mchId, shopId int64) error {
	err := s.deleteShop(mchId, shopId)
	if err == nil {
		err = s.o.DeleteByPk(shop.OnlineShop{}, shopId)
		s.delCache(mchId)
	}
	return err
}

// 删除线下门店
func (s *shopRepo) DeleteOfflineShop(mchId, shopId int64) error {
	err := s.deleteShop(mchId, shopId)
	if err == nil {
		err = s.o.DeleteByPk(shop.OfflineShop{}, shopId)
		s.delCache(mchId)
	}
	return err
}
