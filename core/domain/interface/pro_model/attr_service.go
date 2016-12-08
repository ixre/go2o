package pm

type (
	// 属性
	Attr struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 属性名称
		Name string `db:"name"`
		// 产品模型
		ProModel int64 `db:"pro_model"`
		// 是否作为筛选条件
		IsFilter int `db:"is_filter"`
		// 是否多选
		MultiChk int `db:"multi_chk"`
		// 排列序号
		SortNum int `db:"sort_num"`
	}

	// 属性项
	AttrItem struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 属性编号
		AttrId int64 `db:"attr_id"`
		// 产品模型
		ProModel int64 `db:"pro_model"`
		// 属性值
		Value string `db:"value"`
		// 排列序号
		SortNum int `db:"sort_num"`
	}

	// 产品属性
	ProAttr struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 产品编号
		ProId int64 `db:"pro_id"`
		// 属性编号
		AttrId int64 `db:"attr_id"`
		// 属性值
		AttrData string `db:"attr_data"`
	}
)

// 属性服务
type IAttrService interface {
	// 获取属性
	GetAttr(attrId int32) *Attr
	// 保存属性
	SaveAttr(*Attr) (int32, error)
	// 保存属性项
	SaveItem(*AttrItem) (int32, error)
	// 删除属性
	DeleteAttr(attrId int32) error
	// 删除属性项
	DeleteItem(itemId int32) error
	// 获取属性的属性项
	GetItems(attrId int32) []*AttrItem
	// 获取产品模型的属性
	GetModelAttrs(proModel int32) []*Attr
	// 获取产品的属性
	GetGoodsAttrs(proId int32) []*ProAttr
}

type IProAttrRep interface {
	// Get ProAttr
	GetProAttr(primary interface{}) *ProAttr
	// Save ProAttr
	SaveProAttr(v *ProAttr) (int, error)
	// Delete ProAttr
	DeleteProAttr(primary interface{}) error
	// Select ProAttr
	SelectProAttr(where string, v ...interface{}) []*ProAttr
}
