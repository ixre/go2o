/**
 * Copyright 2015 @ to2.net.
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
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/domain/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	itemImpl "go2o/core/domain/item"
	"go2o/core/infrastructure/format"
	"log"
)

var _ item.IGoodsItemRepo = new(goodsRepo)

type goodsRepo struct {
	db.Connector
	_orm         orm.Orm
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

// 商品仓储
func NewGoodsItemRepo(c db.Connector, catRepo product.ICategoryRepo,
	proRepo product.IProductRepo, proMRepo promodel.IProductModelRepo,
	itemWsRepo item.IItemWholesaleRepo, expressRepo express.IExpressRepo,
	registryRepo registry.IRegistryRepo, shopRepo shop.IShopRepo) *goodsRepo {
	return &goodsRepo{
		Connector:    c,
		_orm:         c.GetOrm(),
		catRepo:      catRepo,
		proRepo:      proRepo,
		proMRepo:     proMRepo,
		itemWsRepo:   itemWsRepo,
		expressRepo:  expressRepo,
		registryRepo: registryRepo,
		shopRepo:     shopRepo,
	}
}

// 获取SKU服务
func (i *goodsRepo) SkuService() item.ISkuService {
	if i._skuService == nil {
		i._skuService = itemImpl.NewSkuServiceImpl(i, i.proMRepo)
	}
	return i._skuService
}

// 获取快照服务
func (i *goodsRepo) SnapshotService() item.ISnapshotService {
	if i._snapService == nil {
		i._snapService = itemImpl.NewSnapshotServiceImpl(i)
	}
	return i._snapService
}

// 创建商品
func (i *goodsRepo) CreateItem(v *item.GoodsItem) item.IGoodsItem {
	return itemImpl.NewItem(i.proRepo, i.catRepo, nil, v, i.registryRepo, i,
		i.proMRepo, i.itemWsRepo, i.expressRepo, i.shopRepo, nil)
}

// 获取商品
func (i *goodsRepo) GetItem(itemId int64) item.IGoodsItem {
	v := i.GetValueGoodsById(itemId)
	if v != nil {
		return i.CreateItem(v)
	}
	return nil
}

// 根据SKU-ID获取商品,SKU-ID为商品ID
func (i *goodsRepo) GetItemBySkuId(skuId int64) interface{} {
	snap := i.GetLatestSnapshot(skuId)
	if snap != nil {
		return i.GetItem(skuId)
	}
	return nil
}

// 获取商品
func (i *goodsRepo) GetValueGoods(itemId, skuId int64) *item.GoodsItem {
	var e = new(item.GoodsItem)
	if i.Connector.GetOrm().GetBy(e, "product_id= $1 AND sku_id= $2", itemId, skuId) == nil {
		return e
	}
	return nil
}

// 获取商品
func (i *goodsRepo) GetValueGoodsById(itemId int64) *item.GoodsItem {
	var e = new(item.GoodsItem)
	if i.Connector.GetOrm().Get(itemId, e) == nil {
		return e
	}
	return nil
}

// 根据SKU获取商品
func (i *goodsRepo) GetValueGoodsBySku(productId, skuId int64) *item.GoodsItem {
	var e = new(item.GoodsItem)
	if i.Connector.GetOrm().GetBy(e, "product_id= $1 AND sku_id= $2", productId, skuId) == nil {
		return e
	}
	return nil
}

// 根据编号获取商品
func (i *goodsRepo) GetGoodsByIds(ids ...int64) ([]*valueobject.Goods, error) {
	var items []*valueobject.Goods
	err := i.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM item_info INNER JOIN product ON item_info.product_id=product.id
     WHERE item_info.id IN (`+format.I64ArrStrJoin(ids)+`)`)

	return items, err
}

// 获取会员价
func (i *goodsRepo) GetGoodSMemberLevelPrice(goodsId int64) []*item.MemberPrice {
	var items []*item.MemberPrice
	if i.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_member_price WHERE goods_id = $1`, goodsId) == nil {
		return items
	}
	return nil
}

// 保存会员价
func (i *goodsRepo) SaveGoodSMemberLevelPrice(v *item.MemberPrice) (int32, error) {
	return orm.I32(orm.Save(i.GetOrm(), v, v.Id))
}

// 移除会员价
func (i *goodsRepo) RemoveGoodSMemberLevelPrice(id int) error {
	return i.Connector.GetOrm().DeleteByPk(item.MemberPrice{}, id)
}

// 保存商品
func (i *goodsRepo) SaveValueGoods(v *item.GoodsItem) (int64, error) {
	return orm.I64(orm.Save(i.GetOrm(), v, int(v.Id)))
}

// 获取已上架的商品
func (i *goodsRepo) GetPagedOnShelvesGoods(shopId int64, catIds []int,
	start, end int, where, orderBy string) (int, []*valueobject.Goods) {
	var s string
	total := 0
	catIdStr := ""
	if catIds != nil && len(catIds) > 0 {
		catIdStr = fmt.Sprintf(" AND cat.id IN (%s)",
			format.IntArrStrJoin(catIds))
	}

	if len(where) != 0 {
		where = " AND " + where
	}
	if len(orderBy) != 0 {
		orderBy += ","
	}

	var list = make([]*valueobject.Goods, 0)
	err := i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM item_info it
	  INNER JOIN product_category cat ON it.cat_id=cat.id
		 WHERE ($1 <=0 OR it.shop_id = $2) AND it.review_state= $3
		  AND it.shelve_state= $4  %s %s`,
		catIdStr, where), &total, shopId, shopId, enum.ReviewPass, item.ShelvesOn)

	if total > 0 {
		s = fmt.Sprintf(`SELECT it.* FROM item_info it INNER JOIN product_category cat ON it.cat_id=cat.id
		 WHERE ($1 <=0 OR it.shop_id = $2) %s AND it.review_state= $3 AND it.shelve_state= $4
		  %s ORDER BY %s it.sort_num DESC,it.update_time DESC LIMIT $6 OFFSET $5`, catIdStr, where, orderBy)
		err = i.Connector.GetOrm().SelectByQuery(&list, s, shopId, shopId,
			enum.ReviewPass, item.ShelvesOn, start, end-start)
	}
	if err != nil {
		log.Println("[ Go2o][ Repo][ Error]:", err.Error())
	}
	return total, list
}

// 获取指定数量已上架的商品
func (i *goodsRepo) GetOnShelvesGoods(mchId int64, start, end int, sortBy string) []*valueobject.Goods {
	var e []*valueobject.Goods
	s := fmt.Sprintf(`SELECT * FROM item_info INNER JOIN product ON product.id = item_info.product_id
		 INNER JOIN product_category ON product.cat_id=product_category.id
		 WHERE supplier_id= $1 AND product.review_state= $2 AND product.shelve_state= $3
		 ORDER BY %s,update_time DESC LIMIT $5 OFFSET $4`,
		sortBy)

	_ = i.Connector.GetOrm().SelectByQuery(&e, s, mchId, enum.ReviewPass,
		item.ShelvesOn, start, end-start)
	return e
}

// 保存快照
func (i *goodsRepo) SaveSnapshot(v *item.Snapshot) (int64, error) {
	r, _, err := i.Connector.GetOrm().Save(v.ItemId, v)
	if r == 0 {
		_, _, err = i.Connector.GetOrm().Save(nil, v)
	}
	return v.ItemId, err
}

// 获取最新的商品快照
func (i *goodsRepo) GetLatestSnapshot(itemId int64) *item.Snapshot {
	e := &item.Snapshot{}
	if i.Connector.GetOrm().Get(itemId, e) == nil {
		return e
	}
	return nil
}

// 根据指定商品快照
func (i *goodsRepo) GetSnapshots(skuIdArr []int64) []item.Snapshot {
	var list []item.Snapshot
	_ = i.Connector.GetOrm().Select(&list, `item_id IN (`+
		format.I64ArrStrJoin(skuIdArr)+`)`)
	return list
}

// 获取最新的商品销售快照
func (i *goodsRepo) GetLatestSalesSnapshot(skuId int64) *item.TradeSnapshot {
	e := new(item.TradeSnapshot)
	if i.Connector.GetOrm().GetBy(e, "sku_id= $1 ORDER BY id DESC", skuId) == nil {
		return e
	}
	return nil
}

// 获取指定的商品销售快照
func (i *goodsRepo) GetSalesSnapshot(id int64) *item.TradeSnapshot {
	e := new(item.TradeSnapshot)
	if i.Connector.GetOrm().Get(id, e) == nil {
		return e
	}
	return nil
}

// 根据Key获取商品销售快照
func (i *goodsRepo) GetSaleSnapshotByKey(key string) *item.TradeSnapshot {
	var e = new(item.TradeSnapshot)
	if i.Connector.GetOrm().GetBy(e, "key= $1", key) == nil {
		return e
	}
	return nil
}

// 保存商品销售快照
func (i *goodsRepo) SaveSalesSnapshot(v *item.TradeSnapshot) (int64, error) {
	return orm.I64(orm.Save(i.Connector.GetOrm(), v, int(v.Id)))
}

// Get ItemSku
func (i *goodsRepo) GetItemSku(primary interface{}) *item.Sku {
	e := item.Sku{}
	err := i._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return nil
}

// Select ItemSku
func (i *goodsRepo) SelectItemSku(where string, v ...interface{}) []*item.Sku {
	var list []*item.Sku
	err := i._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return list
}

// Save ItemSku
func (i *goodsRepo) SaveItemSku(v *item.Sku) (int, error) {
	id, err := orm.Save(i._orm, v, int(v.ID))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return id, err
}

// Delete ItemSku
func (i *goodsRepo) DeleteItemSku(primary interface{}) error {
	err := i._orm.DeleteByPk(item.Sku{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return err
}

// Batch Delete ItemSku
func (i *goodsRepo) BatchDeleteItemSku(where string, v ...interface{}) (int64, error) {
	r, err := i._orm.Delete(item.Sku{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:ItemSku")
	}
	return r, err
}
