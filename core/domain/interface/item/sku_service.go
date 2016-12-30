package item

import (
	"go2o/core/domain/interface/pro_model"
)

type (
	ISkuService interface {
		// 将SKU字符串转为字典,如: 1:2;2:3
		SpecDataToMap(specData string) map[int]int
		// 获取规格和项的数组
		GetSpecItemArray(sku []*Sku) ([]int, []int)
		// 合并SKU数组；主要是SKU编号的复制
		Merge(from []*Sku, to *[]*Sku)
		// 重建SKU数组，将信息附加
		RebuildSkuArray(sku *[]*Sku, it *GoodsItem) error
		// 根据SKU更新商品的信息
		UpgradeBySku(it *GoodsItem, arr []*Sku) error
		// 获取SKU的JSON字符串
		GetSkuJson(skuArr []*Sku) []byte
		// 获取商品的规格(从SKU中读取)
		GetSpecArray(skuArr []*Sku) []*promodel.Spec
		// 获取规格选择HTML
		GetSpecHtml(spec []*promodel.Spec) string
	}

	// 商品SKU
	Sku struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 产品编号
		ProductId int32 `db:"product_id"`
		// 商品编号
		ItemId int32 `db:"item_id"`
		// 标题
		Title string `db:"title"`
		// 图片
		Image string `db:"image"`
		// 规格数据
		SpecData string `db:"spec_data"`
		// 规格字符
		SpecWord string `db:"spec_word"`
		// 产品编码
		Code string `db:"code"`
		// 参考价
		RetailPrice float32 `db:"retail_price"`
		// 价格（分)
		Price float32 `db:"price"`
		// 成本（分)
		Cost float32 `db:"cost"`
		// 重量(克)
		Weight int32 `db:"weight"`
		// 体积（毫升)
		Bulk int32 `db:"bulk"`
		// 库存
		Stock int32 `db:"stock"`
		// 已销售数量
		SaleNum int32 `db:"sale_num"`
	}

	// SKU媒介信息，用于与订单数据交换
	SkuMedia struct {
		// SKU编号
		SkuId int32 `db:"sku_id"`
		// 商品标题
		Title string `db:"title"`
		// 产品编号
		ProductId int32 `db:"product_id"`
		// 商品编号
		ItemId int32 `db:"item_id"`
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
		// 图片
		Image string `db:"image"`
		// 规格数据
		SpecData string `db:"spec_data"`
		// 规格字符
		SpecWord string `db:"spec_word"`
		// 产品编码
		ItemCode string `db:"code"`
		// 参考价
		RetailPrice float32 `db:"retail_price"`
		// 价格（分)
		Price float32 `db:"price"`
		// 重量(克)
		Weight int32 `db:"weight"`
		// 体积（毫升)
		Bulk int32 `db:"bulk"`
		// 库存
		Stock int32 `db:"stock"`
		// 已销售数量
		SaleNum int32 `db:"sale_num"`
	}
)

// 转换为SKU媒体
func ParseSkuMedia(it *GoodsItem, sku *Sku) *SkuMedia {
	media := &SkuMedia{
		CatId:      it.CatId,
		Price:      it.Price,
		ExpressTid: it.ExpressTid,
		ItemCode:   it.Code,
		Stock:      it.StockNum,
		Image:      it.Image,
		Title:      it.Title,
	}
	if sku != nil {
		media.SpecData = sku.SpecData
		media.SpecWord = sku.SpecWord
		media.RetailPrice = sku.RetailPrice
		media.Price = sku.Price
		media.Stock = sku.Stock
		if sku.Image != "" {
			media.Image = sku.Image
		}
		if sku.Title != "" {
			media.Title = sku.Title
		}
		if sku.Code != "" {
			media.ItemCode = sku.Code
		}
	}
	return media
}
