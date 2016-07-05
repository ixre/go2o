/**
 * Copyright 2015 @ z3q.net.
 * name : snapshot
 * author : jarryliu
 * date : 2016-06-28 21:41
 * description :
 * history :
 */
package goods

type (
	// 快照服务
	ISnapshotManager interface {
		// 获取最新的快照
		GetLatestSnapshot() *Snapshot

		// 生成商品快照
		GenerateSnapshot() (int, error)

		// 生成交易快照
		GenerateSaleSnapshot() (int, error)

		// 根据KEY获取已销售商品的快照
		GetSaleSnapshotByKey(key string) *GoodsSnapshot

		// 根据ID获取已销售商品的快照
		GetSaleSnapshot(id int) *GoodsSnapshot
	}

	// 商品快照
	Snapshot struct {
		//SKU编号
		SkuId int `db:"sku_id" auto:"no" pk:"no"`
		//快照编号: 商户编号+g商品编号+快照时间戳
		Key string `db:"snapshot_key"`
		//供应商编号
		VendorId int `db:"vendor_id"`
		//商品编号
		//GoodsId int `db:"goods_id"`
		//商品标题
		GoodsTitle string `db:"goods_title"`
		//小标题
		SmallTitle string `db:"small_title"`
		//货号
		GoodsNo string `db:"goods_no"`
		//货品编号
		ItemId int `db:"item_id"`
		//分类编号
		CategoryId int `db:"cat_id"`
		//SKU  todo:????
		Sku string `db:"-"`
		// 是否上架
		OnShelves int `db:"on_shelves"`
		//图片
		Image string `db:"img"`
		//定价
		Price float32 `db:"price"`
		//销售价
		SalePrice float32 `db:"sale_price"`
		// 单件重量,单位:克(g)
		Weight int `db:"weight"`
		//是否有会员价
		LevelSales int `db:"level_sales"`
		//销售数量
		SaleNum int `db:"sale_num"`
		//库存
		StockNum int `db:"stock_num"`
		//快照时间
		UpdateTime int64 `db:"update_time"`
	}

	// 商品快照
	GoodsSnapshot struct {
		Id           int    `db:"id" auto:"yes" pk:"yes"`
		Key          string `db:"snapshot_key"`
		ItemId       int    `db:"item_id"`
		GoodsId      int    `db:"goods_id"`
		GoodsName    string `db:"goods_name"`
		GoodsNo      string `db:"goods_no"`
		SmallTitle   string `db:"small_title"`
		CategoryName string `db:"category_name"`
		Image        string `db:"img"`

		//成本价
		Cost float32 `db:"cost"`

		//定价
		Price float32 `db:"price"`

		//销售价
		SalePrice  float32 `db:"sale_price"`
		CreateTime int64   `db:"create_time"`
	}
)
