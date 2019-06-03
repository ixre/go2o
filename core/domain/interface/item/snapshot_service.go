/**
 * Copyright 2015 @ to2.net.
 * name : snapshot
 * author : jarryliu
 * date : 2016-06-28 21:41
 * description :
 * history :
 */
package item

type (
	// 快照服务
	ISnapshotService interface {
		// 生成商品快照
		GenerateSnapshot(it *GoodsItem) (int64, error)

		// 获取最新的快照
		GetLatestSnapshot(itemId int64) *Snapshot

		// 获取最新的商品销售快照,如果商品有更新,则更新销售快照
		GetLatestSalesSnapshot(itemId, skuId int64) *TradeSnapshot

		// 根据KEY获取已销售商品的快照
		GetSaleSnapshotByKey(key string) *TradeSnapshot

		// 根据ID获取已销售商品的快照
		GetSalesSnapshot(id int64) *TradeSnapshot
	}

	// 商品快照(针对商品)
	Snapshot struct {
		// 商品编号
		ItemId int64 `db:"item_id" pk:"yes"`
		// 产品编号
		ProductId int64 `db:"product_id"`
		// 快照编码
		Key string `db:"snapshot_key"`
		// 分类编号
		CatId int32 `db:"cat_id"`
		// 供货商编号
		VendorId int32 `db:"vendor_id"`
		// 编号
		BrandId int32 `db:"brand_id"`
		// 商铺编号
		ShopId int32 `db:"shop_id"`
		// 编号分类编号
		ShopCatId int32 `db:"shop_cat_id"`
		// 运费模板
		ExpressTid int32 `db:"express_tid"`
		// 商品标题
		Title string `db:"title"`
		// 短标题
		ShortTitle string `db:"short_title"`
		// 商户编码
		Code string `db:"code"`
		// 商品图片
		Image string `db:"image"`
		// 是否为赠品
		IsPresent int32 `db:"is_present"`
		// 价格区间
		PriceRange string `db:"price_range"`
		// 默认SKU
		SkuId int64 `db:"sku_id"`
		// 成本
		Cost float32 `db:"cost"`
		// 售价
		Price float32 `db:"price"`
		// 零售价
		RetailPrice float32 `db:"retail_price"`
		// 重量(g)
		Weight int32 `db:"weight"`
		// 体积(ml)
		Bulk int32 `db:"bulk"`
		// 会员价
		LevelSales int32 `db:"level_sales"`
		// 上架状态
		ShelveState int32 `db:"shelve_state"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 已销售(交易)商品快照(针对SKU)
	TradeSnapshot struct {
		//快照编号
		Id int64 `db:"id" auto:"yes" pk:"yes"`
		//商品编号
		ItemId int64 `db:"item_id"`
		//商品SKU编号
		SkuId int64 `db:"sku_id"`
		//快照编码: 商户编号+g商品编号+快照时间戳
		SnapshotKey string `db:"snap_key"`
		// 卖家编号
		SellerId int32 `db:"seller_id"`
		// 卖家名称
		//SellerName  string `db:"seller_name"`
		//商品标题
		GoodsTitle string `db:"goods_title"`
		//小标题
		//SmallTitle  string `db:"-"`
		//货号
		GoodsNo string `db:"goods_no"`
		//分类编号
		CategoryId int32 `db:"cat_id"`
		//SKU
		Sku string `db:"sku"`
		//图片
		Image string `db:"img"`
		// 供货价
		Cost float32 `db:"cost"`
		//销售价
		Price float32 `db:"price"`
		// 快照时间
		CreateTime int64 `db:"create_time"`
	}
)
