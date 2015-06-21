/**
 * Copyright 2015 @ S1N1 Team.
 * name : parse.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale
import "go2o/src/core/domain/interface/valueobject"

// 转换包含部分数据的产品值对象
func  ParseToPartialValueItem(v *valueobject.Goods)*ValueItem{
    return &ValueItem{
        Id : v.Item_Id,
        CategoryId:v.CategoryId,
        Name :v.Name,
        GoodsNo :v.GoodsNo,
        Image :v.Image,
        Price:v.Price,
        SalePrice:v.SalePrice,
    }
}

// 转换为商品值对象
func ParseToValueGoods(v *valueobject.Goods)*ValueGoods{
    return &ValueGoods{
        Id :v.GoodsId,
        ItemId:v.Item_Id,
        IsPresent:v.IsPresent,
        SkuId:v.SkuId,
        PromotionFlag:v.PromotionFlag,
        StockNum:v.StockNum,
        SaleNum:v.SaleNum,
        SalePrice:v.SalePrice,
        PromPrice:v.PromPrice,
        Price:v.Price,
    }
}