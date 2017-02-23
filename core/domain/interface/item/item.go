/**
 * Copyright 2015 @ z3q.net.
 * name : goods
 * author : jarryliu
 * date : 2016-06-28 23:54
 * description :
 * history :
 */
package item

import (
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
)

const (
	// 仓库中的商品
	ShelvesInWarehouse int32 = 0
	// 已下架
	ShelvesDown int32 = 1
	// 已上架
	ShelvesOn int32 = 2
	// 已拒绝上架 (不允许上架)
	ShelvesIncorrect int32 = 3
)

var (
	ErrNoSuchGoods *domain.DomainError = domain.NewDomainError(
		"no_such_goods", "商品不存在")

	ErrLatestSnapshot *domain.DomainError = domain.NewDomainError(
		"latest_snapshot", "已经是最新的快照")

	ErrNoSuchSnapshot *domain.DomainError = domain.NewDomainError(
		"no_such_snapshot", "商品快照不存在")

	ErrNotOnShelves *domain.DomainError = domain.NewDomainError(
		"not_on_shelves", "商品已下架")

	ErrGoodsMinProfitRate *domain.DomainError = domain.NewDomainError(
		"err_goods_min_profit_rate", "商品利润率不能低于%s")
)

type (
	// 商品仓储
	IGoodsItemRepo interface {
		// 获取SKU服务
		SkuService() ISkuService
		// 获取快照服务
		SnapshotService() ISnapshotService

		// 创建商品
		CreateItem(v *GoodsItem) IGoodsItem

		// 获取商品
		GetItem(itemId int32) IGoodsItem

		// 获取商品
		GetValueGoods(itemId int32, skuId int32) *GoodsItem

		// 根据SKU-ID获取商品,SKU-ID为商品ID
		//todo: 循环引有,故为interface{}
		GetGoodsBySkuId(skuId int32) interface{}

		// 获取商品
		GetValueGoodsById(goodsId int32) *GoodsItem

		// 根据SKU获取商品
		GetValueGoodsBySku(itemId, sku int32) *GoodsItem

		// 保存商品
		SaveValueGoods(*GoodsItem) (int32, error)

		// 获取在货架上的商品
		GetOnShelvesGoods(mchId int32, start, end int,
			sortBy string) []*valueobject.Goods

		// 获取在货架上的商品
		GetPagedOnShelvesGoods(mchId int32, catIds []int32, start, end int,
			where, orderBy string) (total int, goods []*valueobject.Goods)

		// 根据编号获取商品
		GetGoodsByIds(ids ...int32) ([]*valueobject.Goods, error)

		// 获取会员价
		GetGoodsLevelPrice(goodsId int32) []*MemberPrice

		// 保存会员价
		SaveGoodsLevelPrice(*MemberPrice) (int32, error)

		// 移除会员价
		RemoveGoodsLevelPrice(id int32) error

		// 保存快照
		SaveSnapshot(*Snapshot) (int32, error)

		// 根据指定商品快照
		GetSnapshots(skuIdArr []int32) []Snapshot

		// 获取最新的商品快照
		GetLatestSnapshot(itemId int32) *Snapshot

		// 获取指定的商品快照
		GetSalesSnapshot(id int32) *TradeSnapshot

		// 根据Key获取商品快照
		GetSaleSnapshotByKey(key string) *TradeSnapshot

		// 获取最新的商品销售快照
		GetLatestSalesSnapshot(skuId int32) *TradeSnapshot

		// 保存商品销售快照
		SaveSalesSnapshot(*TradeSnapshot) (int32, error)

		// Get ItemSku
		GetItemSku(primary interface{}) *Sku
		// Select ItemSku
		SelectItemSku(where string, v ...interface{}) []*Sku
		// Save ItemSku
		SaveItemSku(v *Sku) (int, error)
		// Delete ItemSku
		DeleteItemSku(primary interface{}) error
		// Batch Delete ItemSku
		BatchDeleteItemSku(where string, v ...interface{}) (int64, error)
	}

	// 商品,临时改方便辨别
	GoodsItem struct {
		// 商品编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 产品编号
		ProductId int32 `db:"product_id"`
		// 促销标志
		PromFlag int32 `db:"prom_flag"`
		// 分类编号
		CatId int32 `db:"cat_id"`
		// 供货商编号
		VendorId int32 `db:"vendor_id"`
		// 品牌编号(冗余)
		BrandId int32 `db:"brand_id"`
		// 商铺编号
		ShopId int32 `db:"shop_id"`
		// 商铺分类编号
		ShopCatId int32 `db:"shop_cat_id"`
		// 快递模板编号
		ExpressTid int32 `db:"express_tid"`
		// 商品标题
		Title string `db:"title"`
		// 短标题
		ShortTitle string `db:"short_title"`
		// 供货商编码
		Code string `db:"code"`
		// 主图
		Image string `db:"image"`
		// 是否为赠品
		IsPresent int32 `db:"is_present"`
		// 销售价格区间
		PriceRange string `db:"price_range"`
		// 总库存
		StockNum int32 `db:"stock_num"`
		// 销售数量
		SaleNum int32 `db:"sale_num"`
		// SKU数量
		SkuNum int32 `db:"sku_num"`
		// 默认SKU编号
		SkuId int32 `db:"sku_id"`
		// 成本价
		Cost float32 `db:"cost"`
		// 销售价
		Price float32 `db:"price"`
		// 零售价
		RetailPrice float32 `db:"retail_price"`
		// 重量:克(g)
		Weight int32 `db:"weight"`
		// 体积:毫升(ml)
		Bulk int32 `db:"bulk"`
		// 是否上架
		ShelveState int32 `db:"shelve_state"`
		// 审核状态
		ReviewState int32 `db:"review_state"`
		// 审核备注
		ReviewRemark string `db:"review_remark"`
		// 排序序号
		SortNum int32 `db:"sort_num"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
		// 促销价
		PromPrice float32 `db:"-"`

		SkuArray []*Sku `db:"-"`
	}

	// 会员价
	MemberPrice struct {
		Id      int32   `db:"id" pk:"yes" auto:"yes"`
		GoodsId int32   `db:"goods_id"`
		Level   int32   `db:"level"`
		Price   float32 `db:"price"`
		// 限购数量
		MaxQuota int `db:"max_quota"`
		Enabled  int `db:"enabled"`
	}
)

type (
	// 商品
	IGoodsItem interface {
		// 获取聚合根编号
		GetAggregateRootId() int32
		// 设置值
		GetValue() *GoodsItem
		// 获取包装过的商品信息
		GetPackedValue() *valueobject.Goods
		// 设置值
		SetValue(*GoodsItem) error
		// 设置SKU
		SetSku(arr []*Sku) error
		// 保存
		Save() (int32, error)
		// 获取产品
		Product() product.IProduct
		// 商品快照
		Snapshot() *Snapshot
		// 批发
		Wholesale() IWholesaleItem
		// 获取SKU数组
		SkuArray() []*Sku
		// 获取商品的规格
		SpecArray() []*promodel.Spec
		// 获取SKU
		GetSku(skuId int32) *Sku
		// 获取促销信息
		GetPromotions() []promotion.IPromotion
		// 获取促销价
		GetPromotionPrice(level int32) float32
		// 获取会员价销价,返回是否有会原价及价格
		GetLevelPrice(level int32) (bool, float32)
		// 获取促销描述
		GetPromotionDescribe() map[string]string
		// 获取会员价
		GetLevelPrices() []*MemberPrice
		// 保存会员价
		SaveLevelPrice(*MemberPrice) (int32, error)
		// 是否上架
		IsOnShelves() bool
		// 设置上架
		SetShelve(state int32, remark string) error
		// 审核
		Review(pass bool, remark string) error
		// 标记为违规
		Incorrect(remark string) error
		// 更新销售数量,扣减库存
		AddSalesNum(skuId, quantity int32) error
		// 取消销售
		CancelSale(skuId, quantity int32, orderNo string) error
		// 占用库存
		TakeStock(skuId, quantity int32) error
		// 释放库存
		FreeStock(skuId, quantity int32) error
		//// 生成快照
		//GenerateSnapshot() (int64, error)
		//
		//// 获取最新的快照
		//GetLatestSnapshot() *goods.GoodsSnapshot

		// 删除商品
		Destroy() error
	}

	// 商品批发
	IWholesaleItem interface {
		// 获取领域编号
		GetDomainId() int32
		// 开启批发功能
		TurnWholesale(on bool) error
		// 保存
		Save() (int32, error)

		// 根据商品金额获取折扣
		GetWholesaleDiscount(groupId int32, amount int32) float64
		// 获取全部批发折扣
		GetItemDiscount(groupId int32) []*WsItemDiscount
		// 保存批发折扣
		SaveItemDiscount(groupId int32, arr []*WsItemDiscount) error

		// 获取批发价格
		GetWholesalePrice(skuId, quantity int32) float64
		// 根据SKU获取价格设置
		GetSkuPrice(skuId int32) []*WsSkuPrice
		// 保存批发SKU价格设置
		SaveSkuPrice(skuId int32, arr []*WsSkuPrice) error
	}

	IItemWholesaleRepo interface {
		// Get WsItem
		GetWsItem(primary interface{}) *WsItem
		// Select WsItem
		SelectWsItem(where string, v ...interface{}) []*WsItem
		// Save WsItem
		SaveWsItem(v *WsItem, create bool) (int, error)
		// Delete WsItem
		DeleteWsItem(primary interface{}) error
		// Batch Delete WsItem
		BatchDeleteWsItem(where string, v ...interface{}) (int64, error)

		// Get WsSkuPrice
		GetWsSkuPrice(primary interface{}) *WsSkuPrice
		// Select WsSkuPrice
		SelectWsSkuPrice(where string, v ...interface{}) []*WsSkuPrice
		// Save WsSkuPrice
		SaveWsSkuPrice(v *WsSkuPrice) (int, error)
		// Delete WsSkuPrice
		DeleteWsSkuPrice(primary interface{}) error
		// Batch Delete WsSkuPrice
		BatchDeleteWsSkuPrice(where string, v ...interface{}) (int64, error)

		// Get WsItemDiscount
		GetWsItemDiscount(primary interface{}) *WsItemDiscount
		// Select WsItemDiscount
		SelectWsItemDiscount(where string, v ...interface{}) []*WsItemDiscount
		// Save WsItemDiscount
		SaveWsItemDiscount(v *WsItemDiscount) (int, error)
		// Delete WsItemDiscount
		DeleteWsItemDiscount(primary interface{}) error
		// Batch Delete WsItemDiscount
		BatchDeleteWsItemDiscount(where string, v ...interface{}) (int64, error)
	}

	// 批发商品
	WsItem struct {
		// 商品编号
		ItemId int32 `db:"item_id" pk:"yes" auto:"yes"`
		// 是否启用批发
		EnableWholesale int32 `db:"enable_wholesale"`
	}
	// 商品批发价
	WsSkuPrice struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 商品编号
		ItemId int32 `db:"item_id"`
		// SKU编号
		SkuId int32 `db:"sku_id"`
		// 需要数量以上
		RequireQuantity int32 `db:"require_quantity"`
		// 批发价
		WholesalePrice float64 `db:"wholesale_price"`
	}

	// 批发商品折扣
	WsItemDiscount struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 商品编号
		ItemId int32 `db:"item_id"`
		// 客户分组
		BuyerGid int32 `db:"buyer_gid"`
		// 要求金额，默认为0
		RequireAmount int32 `db:"require_amount"`
		// 折扣率
		DiscountRate float64 `db:"discount_rate"`
	}
)
