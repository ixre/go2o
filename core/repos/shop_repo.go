/**
 * Copyright 2014 @ to2.net.
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
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	shopImpl "go2o/core/domain/merchant/shop"
	"log"
)

var _ shop.IShopRepo = new(shopRepo)

type shopRepo struct {
	db.Connector
	valueRepo    valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
	storage      storage.Interface
}

// 创建电子商城
func (s *shopRepo) CreateShop(v *shop.OnlineShop) shop.IShop {
	return shopImpl.NewShop(v, s, s.valueRepo, s.registryRepo)
}

// 检查商户商城是否存在(创建)
func (s *shopRepo) CheckShopExists(vendorId int) bool {
	num := 0
	s.Connector.ExecScalar(`SELECT COUNT(0) FROM mch_online_shop WHERE
	    vendor_id = $1 LIMIT 1`, &num, vendorId)
	return num > 0
}

func (s *shopRepo) ShopCount(vendorId int32, shopType int32) int {
	num := 0
	s.Connector.ExecScalar(`SELECT COUNT(0) FROM mch_shop WHERE
	    vendor_id= $1 AND shop_type = $2`, &num, vendorId, shopType)
	return num
}

func NewShopRepo(c db.Connector, storage storage.Interface,
	valueRepo valueobject.IValueRepo, registryRepo registry.IRegistryRepo) shop.IShopRepo {
	return &shopRepo{
		Connector:    c,
		valueRepo:    valueRepo,
		storage:      storage,
		registryRepo: registryRepo,
	}
}

// 获取商店
func (s *shopRepo) GetShop(shopId int) shop.IShop {
	v := s.GetValueShop(shopId)
	return shopImpl.NewShop2(v, s, s.valueRepo, s.registryRepo)
}

// 商店别名是否存在
func (s *shopRepo) ShopAliasExists(alias string, shopId int) bool {
	id := 0
	s.Connector.ExecScalar(`SELECT id FROM mch_online_shop WHERE
		alias= $1 AND shop_id<> $2 LIMIT 1`, &id, alias, shopId)
	return id > 0
}

// 获取线上商店
func (s *shopRepo) GetOnlineShop(shopId int) *shop.OnlineShop {
	e := shop.OnlineShop{}
	if s.GetOrm().Get(shopId, &e) != nil {
		return nil
	}
	return &e
}

// 保存线上商店
func (s *shopRepo) SaveOnlineShop(v *shop.OnlineShop) (int, error) {
	id, err := orm.Save(s.GetOrm(), v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchOnlineShop")
	}
	return id, err
}

// 获取线下商店
func (s *shopRepo) GetOfflineShop(shopId int) *shop.OfflineShop {
	e := shop.OfflineShop{}
	if s.GetOrm().Get(shopId, &e) != nil {
		return nil
	}
	return &e
}

// 保存线下商店
func (s *shopRepo) SaveOfflineShop(v *shop.OfflineShop, create bool) error {
	var err error
	if create {
		_, _, err = s.GetOrm().Save(nil, v)
	} else {
		_, _, err = s.GetOrm().Save(v.ShopId, v)
	}
	return err
}

// 保存API信息
func (s *shopRepo) SaveApiInfo(v *merchant.ApiInfo) error {
	_, err := orm.Save(s.GetOrm(), v, int(v.MerchantId))
	return err
}

// 获取API信息
func (s *shopRepo) GetApiInfo(mchId int) *merchant.ApiInfo {
	var d *merchant.ApiInfo = new(merchant.ApiInfo)
	if err := s.GetOrm().Get(mchId, d); err == nil {
		return d
	}
	return nil
}

func (s *shopRepo) SaveShop(v *shop.Shop) (int32, error) {
	id, err := orm.I32(orm.Save(s.GetOrm(), v, int(v.Id)))
	if err == nil {
		s.delCache(int(v.VendorId))
	}
	return id, err
}

func (s *shopRepo) GetValueShop(shopId int) *shop.Shop {
	v := &shop.Shop{}
	err := s.Connector.GetOrm().Get(shopId, v)
	if err == nil {
		return v
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:MchShop")
	}
	return nil
}

func (s *shopRepo) delCache(mchId int) {
	PrefixDel(s.storage, fmt.Sprintf("go2o:repo:shop:%d:*", mchId))
}

func (s *shopRepo) getShopCacheKey(mchId int32) string {
	return fmt.Sprintf("go2o:repo:shop:%d:shops", mchId)
}

func (s *shopRepo) GetOnlineShopOfMerchant(vendorId int) *shop.OnlineShop {
	v := shop.OnlineShop{}
	err := s.Connector.GetOrm().GetBy(&v, "vendor_id= $1 LIMIT 1", vendorId)
	if err == nil {
		return &v
	}
	return nil
}

/**/
func (s *shopRepo) GetShopsOfMerchant(mchId int32) []shop.Shop {
	shops := make([]shop.Shop, 0)
	key := s.getShopCacheKey(mchId)
	jsonStr, err := s.storage.GetString(key)
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &shops)
	}
	if err != nil {
		err = s.Connector.GetOrm().SelectByQuery(&shops,
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

func (s *shopRepo) deleteShop(mchId, shopId int) error {
	_, err := s.Connector.GetOrm().Delete(shop.Shop{},
		"vendor_id= $1 AND id= $2", mchId, shopId)
	s.delCache(mchId)
	return err
}

// 删除线上商店
func (s *shopRepo) DeleteOnlineShop(mchId, shopId int) error {
	err := s.deleteShop(mchId, shopId)
	if err == nil {
		err = s.Connector.GetOrm().DeleteByPk(shop.OnlineShop{}, shopId)
		s.delCache(mchId)
	}
	return err
}

// 删除线下门店
func (s *shopRepo) DeleteOfflineShop(mchId, shopId int) error {
	err := s.deleteShop(mchId, shopId)
	if err == nil {
		err = s.Connector.GetOrm().DeleteByPk(shop.OfflineShop{}, shopId)
		s.delCache(mchId)
	}
	return err
}
