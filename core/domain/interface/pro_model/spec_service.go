package promodel

type (
	// 规格
	Spec struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 产品模型
		ProModel int64 `db:"pro_model"`
		// 规格名称
		Name int64 `db:"name"`
		// 排列序号
		SortNum int `db:"sort_num"`
	}

	// 规格项
	SpecItem struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 规格项名称
		Name string `db:"name"`
		// 规格编号
		SpecId int64 `db:"spec_id"`
		// 产品模型（冗余)
		ProModel int64 `db:"pro_model"`
		// 规格颜色
		SpecColor string `db:"spec_color"`
		// 排列序号
		SortNum int `db:"sort_num"`
	}
)

// 规格服务
type ISpecService interface {
	// 获取规格
	GetSpec(specId int32) *Spec
	// 保存规格
	SaveSpec(*Spec) (int32, error)
	// 保存规格项
	SaveItem(*SpecItem) (int32, error)
	// 删除规格
	DeleteSpec(specId int32) error
	// 删除规格项
	DeleteItem(itemId int32) error
	// 获取规格的规格项
	GetItems(specId int32) []*SpecItem
	// 获取产品模型的规格
	GetModelSpecs(proModel int32) []*Spec
}
