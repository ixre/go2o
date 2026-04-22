package valueobject

// SEO信息表
type OSeo struct {
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 使用者编号
	UseId int64 `db:"use_id"`
	// 使用者类型
	UseType int64 `db:"use_type"`
	// 标题
	Title string `db:"title"`
	// 关键词
	Keywords string `db:"keywords"`
	// 描述
	Description string `db:"description"`
}

type IOtherRepo struct {
}
