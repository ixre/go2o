package promodel

import "go2o/core/infrastructure/domain"

var (
	ErrEmptyAttrArray *domain.DomainError = domain.NewError(
		"err_empty_attr_array", "请至少包含一个属性")
	ErrEmptySpecArray *domain.DomainError = domain.NewError(
		"err_empty_spec_array", "请至少包含一个规格")
	ErrEmptyBrandArray *domain.DomainError = domain.NewError(
		"err_empty_brand_array", "请至少包含一个品牌")
)

type ProModel struct {
	// 编号
	ID int32 `db:"id" pk:"yes" auto:"yes"`
	// 名称
	Name string `db:"name"`
	// 是否启用
	Enabled int `db:"enabled"`
	// 属性字符
	AttrStr string `db:"attr_str"`
	// 规格字符
	SpecStr string `db:"spec_str"`
	// 属性
	Attrs []*Attr `db:"-"`
	// 规格
	Specs []*Spec `db:"-"`
	// 关联品牌
	BrandArray []int32 `db:"-"`
}

// 产品模型
type IModel interface {
	// 获取聚合根编号
	GetAggregateRootId() int32
	// 获取值
	Value() *ProModel
	// 获取属性
	Attrs() []*Attr
	// 获取规格
	Specs() []*Spec
	// 设置属性
	SetAttrs([]*Attr) error
	// 设置规格
	SetSpecs([]*Spec) error
	// 获取关联的品牌编号
	Brands() []*ProBrand
	// 设置关联品牌
	SetBrands(brandId []int32) error
	// 保存
	Save() (int32, error)
	// 是否启用
	Enabled() bool
}

type IProModelRepo interface {
	// 创建商品模型
	CreateModel(v *ProModel) IModel
	// 获取商品模型
	GetModel(id int32) IModel

	// 属性服务
	AttrService() IAttrService
	// 规格服务
	SpecService() ISpecService

	//获取品牌服务
	BrandService() IBrandService
	// 获取模型的商品品牌
	GetModelBrands(proModel int32) []*ProBrand

	// Get ProModel
	GetProModel(primary interface{}) *ProModel
	// Select ProModel
	SelectProModel(where string, v ...interface{}) []*ProModel
	// Save ProModel
	SaveProModel(v *ProModel) (int, error)
	// Delete ProModel
	DeleteProModel(primary interface{}) error

	// Get Attr
	GetAttr(primary interface{}) *Attr
	// Select Attr
	SelectAttr(where string, v ...interface{}) []*Attr
	// Save Attr
	SaveAttr(v *Attr) (int, error)
	// Delete Attr
	DeleteAttr(primary interface{}) error
	// Batch Delete Attr
	BatchDeleteAttr(where string, v ...interface{}) (int64, error)

	// Get AttrItem
	GetAttrItem(primary interface{}) *AttrItem
	// Select AttrItem
	SelectAttrItem(where string, v ...interface{}) []*AttrItem
	// Save AttrItem
	SaveAttrItem(v *AttrItem) (int, error)
	// Delete AttrItem
	DeleteAttrItem(primary interface{}) error
	// Batch Delete AttrItem
	BatchDeleteAttrItem(where string, v ...interface{}) (int64, error)

	// Get Spec
	GetSpec(primary interface{}) *Spec
	// Select Spec
	SelectSpec(where string, v ...interface{}) []*Spec
	// Save Spec
	SaveSpec(v *Spec) (int, error)
	// Delete Spec
	DeleteSpec(primary interface{}) error
	// Batch Delete Spec
	BatchDeleteSpec(where string, v ...interface{}) (int64, error)

	// Get SpecItem
	GetSpecItem(primary interface{}) *SpecItem
	// Select SpecItem
	SelectSpecItem(where string, v ...interface{}) []*SpecItem
	// Save SpecItem
	SaveSpecItem(v *SpecItem) (int, error)
	// Delete SpecItem
	DeleteSpecItem(primary interface{}) error
	// Batch Delete SpecItem
	BatchDeleteSpecItem(where string, v ...interface{}) (int64, error)

	// Get ProBrand
	GetProBrand(primary interface{}) *ProBrand
	// Save ProBrand
	SaveProBrand(v *ProBrand) (int, error)
	// Delete ProBrand
	DeleteProBrand(primary interface{}) error
	// Select ProBrand
	SelectProBrand(where string, v ...interface{}) []*ProBrand

	// Batch Delete ProBrand
	BatchDeleteProBrand(where string, v ...interface{}) (int64, error)

	// Get ProModelBrand
	GetProModelBrand(primary interface{}) *ProModelBrand
	// Save ProModelBrand
	SaveProModelBrand(v *ProModelBrand) (int, error)
	// Delete ProModelBrand
	DeleteProModelBrand(primary interface{}) error
	// Select ProModelBrand
	SelectProModelBrand(where string, v ...interface{}) []*ProModelBrand
	// Batch Delete ProModelBrand
	BatchDeleteProModelBrand(where string, v ...interface{}) (int64, error)
}
