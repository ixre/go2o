package uams

type (
	Result struct {
		ErrCode int32             `thrift:"ErrCode,1" db:"ErrCode" json:"ErrCode"`
		ErrMsg  string            `thrift:"ErrMsg,2" db:"ErrMsg" json:"ErrMsg"`
		Data    map[string]string `thrift:"Data,3" db:"Data" json:"Data"`
	}

	// 菜单
	Res struct {
		ID       int64
		Title    string
		Code     string
		URI      string
		Data     map[string]string
		Children []*Res
	}

	// 资源
	UamsResource struct {
		// ID
		ID int `db:"id" pk:"yes" auto:"yes"`
		// 上级资源
		ParentId int `db:"parent_id"`
		// 应用编号
		AppId int `db:"app_id"`
		// 资源标题
		ResTitle string `db:"res_title"`
		// 资源类型，如menu/button
		ResType string `db:"res_type"`
		// 资源编码
		ResCode string `db:"res_code"`
		// 资源权限可选值，用"|"分割
		ResOption string `db:"res_option"`
		// 资源是否单独开放给用户
		ResOpen int32 `db:"res_open"`
		// 资源图标
		ResIcon string `db:"-"`
		// 资源URL
		ResUri string `db:"res_uri"`
		// 资源数据
		ResData string `db:"res_data"`
		// 资源状态
		ResState int32 `db:"res_state"`
	}

	// 资源
	ResourceTreeNode struct {
		// ID
		ID int `json:"id"`
		// 资源标题
		Title string `json:"title"`
		// 资源类型，如menu/button
		Type string `json:"type"`
		// 资源编码
		Code string `json:"code"`
		// 图标
		Icon string `json:"icon"`
		// 资源URL
		Uri string `json:"uri"`
		// 子资源
		Children []*ResourceTreeNode `json:"children"`
	}

	// 部门
	UamsDepart struct {
		// 部门ID
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 应用编号
		AppId int64 `db:"app_id"`
		// 上级部门编号
		DeptId int64 `db:"dept_id"`
		// 部门名称
		DeptName string `db:"dept_name"`
		// 部门编码
		DeptCode string `db:"dept_code"`
		// 部门描述
		DeptDesc string `db:"dept_desc"`
		// 子部门
		Children []*UamsDepart `db:"-"`
	}

	// 用户字段
	UamsUserField struct {
		// 字段ID
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 应用编号
		AppId int64 `db:"app_id"`
		// 字段名称
		FieldName string `db:"field_name"`
		// 字段编码
		FieldCode string `db:"field_code"`
		// 字段类型
		FieldType string `db:"field_type"`
		// 字段可选值
		FieldOption string `db:"field_option"`
		// 字段默认值
		FieldDefault string `db:"field_default"`
		// 字段长度
		FieldLen int64 `db:"field_len"`
		// 字段是否必填
		FieldRequire int32 `db:"field_require"`
		// 排序序号
		SortNumber int64 `db:"sort_number"`
		// 字段状态
		FieldState int32 `db:"field_state"`
	}

	// 角色
	UamsRole struct {
		// 角色编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 应用编号
		AppId int64 `db:"app_id"`
		// 角色名称
		RoleName string `db:"role_name"`
		// 角色描述
		RoleDesc string `db:"role_desc"`
		// 默认角色
		IsDefault int `db:"is_default"`
		Enabled   int `db:"enabled"`
	}
)
