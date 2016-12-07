package pmd

// 产品模型
type ProModel struct {
    // 编号
    Id      int64 `db:"id" pk:"yes" auto:"yes"`
    // 名称
    Name    string `db:"name"`
    // 是否启用
    Enabled int `db:"enabled"`
}
