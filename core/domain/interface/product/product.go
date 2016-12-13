/**
 * Copyright 2015 @ z3q.net.
 * name : item
 * author : jarryliu
 * date : 2016-06-29 09:31
 * description :
 * history :
 */
package product

import (
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
)

const (
	// 已下架
	ShelvesDown int32 = 1
	// 已上架
	ShelvesOn int32 = 2
	// 已拒绝上架 (不允许上架)
	ShelvesIncorrect int32 = 3
)

var (
	ErrNoSuchProduct *domain.DomainError = domain.NewDomainError(
		"err_product_no_such_product", "产品不存在",
	)
	ErrNoBrand *domain.DomainError = domain.NewDomainError(
		"err_product_no_brand", "未设置商品品牌")

	ErrVendor *domain.DomainError = domain.NewDomainError(
		"err_not_be_review", "商品供应商不正确")

	ErrNotBeReview *domain.DomainError = domain.NewDomainError(
		"err_not_be_review", "商品还未通过审核")

	ErrItemNameLength *domain.DomainError = domain.NewDomainError(
		"err_item_name_length", "商品标题至少10个字")

	ErrItemIncorrect *domain.DomainError = domain.NewDomainError(
		"err_item_incorrect", "商品已被违规下架")

	ErrNotUploadImage *domain.DomainError = domain.NewDomainError(
		"err_goods_not_upload_image", "请上传商品图片")

	ErrDescribeLength *domain.DomainError = domain.NewDomainError(
		"err_item_describe_length", "商品描述至少20个字符")

	ErrNilRejectRemark *domain.DomainError = domain.NewDomainError(
		"err_item_nil_reject_remark", "原因不能为空")
)

type (
	IProduct interface {
		// 获取领域对象编号
		GetDomainId() int32
		// 获取商品的值
		GetValue() Product
		// 设置产品的值
		SetValue(v *Product) error
		// 是否上架
		IsOnShelves() bool

		// 获取销售标签
		//GetSaleLabels() []*Label

		// 保存销售标签
		//SaveSaleLabels([]int) error

		// 设置商品描述
		SetDescribe(describe string) error

		// 设置上架
		SetShelve(state int32, remark string) error

		// 审核
		Review(pass bool, remark string) error

		// 标记为违规
		Incorrect(remark string) error

		// 保存
		Save() (int32, error)

		// 销毁产品
		Destroy() error
	}

	IProductRepo interface {
		// 创建产品
		CreateProduct(*Product) IProduct
		// 根据产品编号获取货品
		GetProduct(id int32) IProduct
		// 获取货品
		GetProductValue(itemId int32) *Product
		// 根据id获取货品
		GetProductsById(ids ...int32) ([]*Product, error)
		SaveProductValue(*Product) (int32, error)
		//todo:  到商品
		// 获取在货架上的商品
		GetPagedOnShelvesProduct(supplierId int32, catIds []int32, start, end int) (total int, goods []*Product)
		//todo:  到商品
		// 获取货品销售总数
		GetProductSaleNum(productId int32) int
		// 删除货品
		DeleteProduct(productId int32) error
	}
)

// 产品
type Product struct {
	// 编号
	Id int32 `db:"id" auto:"yes" pk:"yes"`
	// 分类
	CategoryId int32 `db:"cat_id"`
	// 名称
	Name string `db:"name"`
	//供应商编号(暂时同mch_id)
	VendorId int32 `db:"supplier_id"`
	// 商铺编号
	ShopId int64 `db:"shop_id"`
	// 品牌编号
	BrandId int64 `db:"brand_id"`
	// 商家编码
	Code string `db:"code"`
	// 小标题
	SmallTitle string `db:"small_title"`
	// 图片
	Image string `db:"img"`
	// 成本价
	Cost float32 `db:"cost"`
	// 重量:克(g)
	Weight float32 `db:"weight"`
	// 体积:毫升(ml)
	Bulk int64 `db:"bulk"`
	//定价
	Price float32 `db:"price"`
	//参考销售价
	SalePrice float32 `db:"sale_price"`
	// 运费模板编号
	ExpressTplId int32 `db:"express_tid"`
	// 描述
	Description string `db:"description"`
	// 上架状态
	ShelveState int32 `db:"shelve_state"`
	// 审核状态
	ReviewState int32 `db:"review_state"`
	// 备注
	Remark string `db:"remark"`
	// 状态
	State int32 `db:"state"`
	// 创建时间
	CreateTime int64 `db:"create_time"`
	// 更新时间
	UpdateTime int64 `db:"update_time"`
	// 排序编号
	SortNum int32 `db:"sort_num"`
}

// 转换包含部分数据的产品值对象
func ParseToPartialValueItem(v *valueobject.Goods) *Product {
	return &Product{
		Id:         v.ProductId,
		CategoryId: v.CategoryId,
		Name:       v.Name,
		Code:       v.GoodsNo,
		Image:      v.Image,
		Price:      v.Price,
		SalePrice:  v.SalePrice,
	}
}
