package item

import (
	"fmt"
	"github.com/ixre/gof/math"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"strconv"
)

var iJsonUtil = &itemJsonUtil{}

// 规格项
type specJdo struct {
	// 规格名称
	Name string `json:"name"`
	// 规格项
	Value []specItemJdo `json:"value"`
}

// 规格值
type specItemJdo struct {
	// 名称
	Name string `json:"name"`
	// 代码
	Code string `json:"code"`
	// 图片地址
	ImageUrl string `json:"imageUrl"`
}

// SKU传输对象
type skuJdo struct {
	// SKU编号
	SkuId string
	// SKU数据
	SpecData string
	// SKU文本
	SpecWord string
	// 商品编码
	Code string
	// 价格
	Price float64
	// 折扣价
	DiscountPrice float64
	// 可售数量
	CanSalesQuantity int32
	// 已售数量
	SalesCount int32
	// 数量与价格字典
	PriceArray []skuPriceJdo
}

type skuPriceJdo struct {
	Quantity int32
	Price    float64
}

type itemJsonUtil struct {
}

// 获取规格JSON数据
func (s *itemJsonUtil) getSpecJdo(spec promodel.SpecList) []specJdo {
	arr := make([]specJdo, len(spec))
	for i, v := range spec {
		arr[i] = specJdo{Name: v.Name}
		arr[i].Value = make([]specItemJdo, len(v.Items))
		for j, v2 := range v.Items {
			arr[i].Value[j] = specItemJdo{
				Name:     v2.Value,
				Code:     fmt.Sprintf("%d:%d", v.ID, v2.ID),
				ImageUrl: v2.Color,
			}
		}
	}
	return arr
}

// 获取SKU的JSON字符串
func (s *itemJsonUtil) getSkuJdo(skuArr []*item.Sku) []skuJdo {
	arr := make([]skuJdo, len(skuArr))
	for i, v := range skuArr {
		arr[i] = skuJdo{
			SkuId:         strconv.Itoa(int(v.ID)),
			SpecData:      v.SpecData,
			SpecWord:      v.SpecWord,
			Price:         math.Round(float64(v.Price), 2),
			DiscountPrice: math.Round(float64(v.Price), 2),
			PriceArray:    []skuPriceJdo{},
		}
	}
	return arr
}
