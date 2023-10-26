/**
 * Copyright 2015 @ 56x.net.
 * name : goods_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repos

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	promodel "github.com/ixre/go2o/core/domain/interface/pro_model"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	itemImpl "github.com/ixre/go2o/core/domain/item"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

var _ item.IItemRepo = new(itemRepoImpl)

type itemRepoImpl struct {
	db.Connector
	o            orm.Orm
	_skuService  item.ISkuService
	_snapService item.ISnapshotService
	catRepo      product.ICategoryRepo
	proRepo      product.IProductRepo
	itemWsRepo   item.IItemWholesaleRepo
	expressRepo  express.IExpressRepo
	proMRepo     promodel.IProductModelRepo
	registryRepo registry.IRegistryRepo
	shopRepo     shop.IShopRepo
}

var ormMapped = false

// NewGoodsItemRepo 商品仓储
func NewGoodsItemRepo(o orm.Orm, catRepo product.ICategoryRepo,
	proRepo product.IProductRepo, proMRepo promodel.IProductModelRepo,
	itemWsRepo item.IItemWholesaleRepo, expressRepo express.IExpressRepo,
	registryRepo registry.IRegistryRepo, shopRepo shop.IShopRepo) *itemRepoImpl {
	if !ormMapped {
		_ = o.Mapping(item.Image{}, "item_image")
		ormMapped = true
	}
	return &itemRepoImpl{
		Connector:    o.Connector(),
		o:            o,
		catRepo:      catRepo,
		proRepo:      proRepo,
		proMRepo:     proMRepo,
		itemWsRepo:   itemWsRepo,
		expressRepo:  expressRepo,
		registryRepo: registryRepo,
		shopRepo:     shopRepo,
	}
}

// SkuService 获取SKU服务
func (i *itemRepoImpl) SkuService() item.ISkuService {
	if i._skuService == nil {
		i._skuService = itemImpl.NewSkuServiceImpl(i, i.proMRepo)
	}
	return i._skuService
}

// SnapshotService 获取快照服务
func (i *itemRepoImpl) SnapshotService() item.ISnapshotService {
	if i._snapService == nil {
		i._snapService = itemImpl.NewSnapshotServiceImpl(i)
	}
	return i._snapService
}

// CreateItem 创建商品
func (i *itemRepoImpl) CreateItem(v *item.GoodsItem) item.IGoodsItemAggregateRoot {
	return itemImpl.NewItem(i.proRepo, i.catRepo, nil, v, i.registryRepo, i,
		i.proMRepo, i.itemWsRepo, i.expressRepo, i.shopRepo, nil)
}

// GetItem 获取商品
func (i *itemRepoImpl) GetItem(itemId int64) item.IGoodsItemAggregateRoot {
	v := i.GetValueGoodsById(itemId)
	if v != nil {
		return i.CreateItem(v)
	}
	return nil
}

// GetItemBySkuId 根据SKU-ID获取商品,SKU-ID为商品ID
func (i *itemRepoImpl) GetItemBySkuId(skuId int64) interface{} {
	snap := i.GetLatestSnapshot(skuId)
	if snap != nil {
		return i.GetItem(skuId)
	}
	return nil
}

// GetValueGoods 获取商品
func (i *itemRepoImpl) GetValueGoods(itemId, skuId int64) *item.GoodsItem {
	var e = new(item.GoodsItem)
	if i.o.GetBy(e, "product_id= $1 AND sku_id= $2", itemId, skuId) == nil {
		return e
	}
	return nil
}

// GetItemImages  获取商品图片
func (i *itemRepoImpl) GetItemImages(itemId int64) []*item.Image {
	list := make([]*item.Image, 0)
	err := i.o.Select(&list, "item_id=$1", itemId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemImage")
	}
	return list
}

// SaveItemImage 保存商品图片
func (i *itemRepoImpl) SaveItemImage(v *item.Image) (int, error) {
	id, err := orm.Save(i.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemImage")
	}
	return id, err
}

// DeleteItemImage 删除商品图片
func (i *itemRepoImpl) DeleteItemImage(id int64) error {
	err := i.o.DeleteByPk(item.Image{}, id)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemImage")
	}
	return err
}

// GetValueGoodsById 获取商品
func (i *itemRepoImpl) GetValueGoodsById(itemId int64) *item.GoodsItem {
	var e = new(item.GoodsItem)
	if i.o.Get(itemId, e) == nil {
		return e
	}
	return nil
}

// GetValueGoodsBySku 根据SKU获取商品
func (i *itemRepoImpl) GetValueGoodsBySku(productId, skuId int64) *item.GoodsItem {
	var e = new(item.GoodsItem)
	if i.o.GetBy(e, "product_id= $1 AND sku_id= $2", productId, skuId) == nil {
		return e
	}
	return nil
}

// GetGoodsByIds 根据编号获取商品
func (i *itemRepoImpl) GetGoodsByIds(ids ...int64) ([]*valueobject.Goods, error) {
	var items []*valueobject.Goods
	err := i.o.SelectByQuery(&items,
		`SELECT * FROM item_info INNER JOIN product ON item_info.product_id=product.id
     WHERE item_info.id IN (`+format.I64ArrStrJoin(ids)+`)`)

	return items, err
}

// GetGoodSMemberLevelPrice 获取会员价
func (i *itemRepoImpl) GetGoodSMemberLevelPrice(goodsId int64) []*item.MemberPrice {
	var items []*item.MemberPrice
	if i.o.SelectByQuery(&items,
		`SELECT * FROM gs_member_price WHERE goods_id = $1`, goodsId) == nil {
		return items
	}
	return nil
}

// SaveGoodSMemberLevelPrice 保存会员价
func (i *itemRepoImpl) SaveGoodSMemberLevelPrice(v *item.MemberPrice) (int32, error) {
	return orm.I32(orm.Save(i.o, v, v.Id))
}

// RemoveGoodSMemberLevelPrice 移除会员价
func (i *itemRepoImpl) RemoveGoodSMemberLevelPrice(id int) error {
	return i.o.DeleteByPk(item.MemberPrice{}, id)
}

// SaveValueGoods 保存商品
func (i *itemRepoImpl) SaveValueGoods(v *item.GoodsItem) (int64, error) {
	return orm.I64(orm.Save(i.o, v, int(v.Id)))
}

// GetOnShelvesGoods 获取指定数量已上架的商品
func (i *itemRepoImpl) GetOnShelvesGoods(mchId int64, start, end int, sortBy string) []*valueobject.Goods {
	var e []*valueobject.Goods
	s := fmt.Sprintf(`SELECT * FROM item_info INNER JOIN product ON product.id = item_info.product_id
		 INNER JOIN product_category ON product.cat_id=product_category.id
		 WHERE supplier_id= $1 AND product.review_state= $2 AND product.shelve_state= $3
		 ORDER BY %s,update_time DESC LIMIT $5 OFFSET $4`,
		sortBy)

	_ = i.o.SelectByQuery(&e, s, mchId, enum.ReviewPass,
		item.ShelvesOn, start, end-start)
	return e
}

// SaveSnapshot 保存快照
func (i *itemRepoImpl) SaveSnapshot(v *item.Snapshot) (int64, error) {
	_, r, err := i.o.Save(v.ItemId, v)
	if r == 0 {
		_, _, err = i.o.Save(nil, v)
	}
	return v.ItemId, err
}

// DeleteSnapshot 删除商品快照
func (i *itemRepoImpl) DeleteSnapshot(itemId int64) error {
	return i.o.DeleteByPk(&item.Snapshot{}, itemId)
}

// GetLatestSnapshot 获取最新的商品快照
func (i *itemRepoImpl) GetLatestSnapshot(itemId int64) *item.Snapshot {
	e := &item.Snapshot{}
	if i.o.Get(itemId, e) == nil {
		return e
	}
	return nil
}

// GetSnapshots 根据指定商品快照
func (i *itemRepoImpl) GetSnapshots(skuIdArr []int64) []item.Snapshot {
	var list []item.Snapshot
	_ = i.o.Select(&list, `item_id IN (`+
		format.I64ArrStrJoin(skuIdArr)+`)`)
	return list
}

// GetLatestSalesSnapshot 获取最新的商品销售快照
func (i *itemRepoImpl) GetLatestSalesSnapshot(itemId int64, skuId int64) *item.TradeSnapshot {
	e := new(item.TradeSnapshot)
	if i.o.GetBy(e, "item_id= $1 AND sku_id= $2 ORDER BY id DESC", itemId, skuId) == nil {
		return e
	}
	return nil
}

// GetSalesSnapshot 获取指定的商品销售快照
func (i *itemRepoImpl) GetSalesSnapshot(id int64) *item.TradeSnapshot {
	e := new(item.TradeSnapshot)
	if i.o.Get(id, e) == nil {
		return e
	}
	return nil
}

// GetSaleSnapshotByKey 根据Key获取商品销售快照
func (i *itemRepoImpl) GetSaleSnapshotByKey(key string) *item.TradeSnapshot {
	var e = new(item.TradeSnapshot)
	if i.o.GetBy(e, "key= $1", key) == nil {
		return e
	}
	return nil
}

// 保存商品销售快照
func (i *itemRepoImpl) SaveSalesSnapshot(v *item.TradeSnapshot) (int64, error) {
	return orm.I64(orm.Save(i.o, v, int(v.Id)))
}

// Get ItemSku
func (i *itemRepoImpl) GetItemSku(primary interface{}) *item.Sku {
	e := item.Sku{}
	err := i.o.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return nil
}

// Select ItemSku
func (i *itemRepoImpl) SelectItemSku(where string, v ...interface{}) []*item.Sku {
	var list []*item.Sku
	err := i.o.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return list
}

// Save ItemSku
func (i *itemRepoImpl) SaveItemSku(v *item.Sku) (int, error) {
	id, err := orm.Save(i.o, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return id, err
}

// DeleteItem 删除商品
func (i *itemRepoImpl) DeleteItem(itemId int) error {
	return i.o.DeleteByPk(&item.GoodsItem{}, itemId)
}

// Delete ItemSku
func (i *itemRepoImpl) DeleteItemSku(primary interface{}) error {
	err := i.o.DeleteByPk(item.Sku{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return err
}

// Batch Delete ItemSku
func (i *itemRepoImpl) BatchDeleteItemSku(where string, v ...interface{}) (int64, error) {
	r, err := i.o.Delete(item.Sku{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return r, err
}
