package promodel

type (
	// 属性
	Attr struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 产品模型
		ModelId int `db:"prod_model"`
		// 属性名称
		Name string `db:"name"`
		// 是否作为筛选条件
		IsFilter int `db:"is_filter"`
		// 是否多选
		MultiCheck int `db:"multi_check"`
		// 属性项值
		ItemValues string `db:"item_values"`
		// 排列序号
		SortNum int `db:"sort_num"`
		// 属性项
		Items []*AttrItem `db:"-"`
	}
	// 属性项
	AttrItem struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 属性编号
		AttrId int `db:"attr_id"`
		// 产品模型
		ModelId int `db:"prod_model"`
		// 属性值
		Value string `db:"value"`
		// 排列序号
		SortNum int `db:"sort_num"`
	}
)

// 属性服务
type IAttrService interface {
	// 获取属性
	GetAttr(attrId int) *Attr
	// 保存属性
	SaveAttr(*Attr) (int, error)
	// 保存属性项
	SaveItem(*AttrItem) (int, error)
	// 删除属性
	DeleteAttr(attrId int) error
	// 删除属性项
	DeleteItem(itemId int) error
	// 获取属性的属性项
	GetItems(attrId int) []*AttrItem
	// 获取产品模型的属性
	GetModelAttrs(proModel int) []*Attr
}
