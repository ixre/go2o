package promodel

// 产品模型
type ProModel struct {
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 名称
	Name string `db:"name"`
	// 是否启用
	Enabled int `db:"enabled"`
	// 关联品牌数
	Brands int64 `db:"brands"`
}

// 产品模型
type IModel interface {
	// 获取聚合根编号
	GetAggregateRootId() int32
	// 获取值
	Value() *ProModel
	// 保存
	Save() (int32, error)
	// 是否启用
	Enabled() bool
	// 获取关联的品牌编号
	Brands() []*ProBrand
	// 关联品牌
	SaveBrands(brandId []int32) error
	// 获取属性
	Attrs() []*Attr
	// 保存属性
	SaveAttrs([]*Attr) error
	// 获取规格
	Specs() []*Spec
	// 保存规格
	SaveSpecs([]*Spec) error
}
