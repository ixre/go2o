package model

// 部门
type PermDept struct {
	// Id
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 名称
	Name string `db:"name"`
	// 编码
	Code string `db:"code"`
	// 上级部门
	Pid int64 `db:"pid"`
	// 状态
	Enabled int16 `db:"enabled"`
	// 创建日期
	CreateTime int64 `db:"create_time"`
}

// 数据字典详情
type PermDictDetail struct {
	// Id
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 字典标签
	Label string `db:"label"`
	// 字典值
	Value string `db:"value"`
	// 排序
	Sort string `db:"sort"`
	// 字典id
	DictId int64 `db:"dict_id"`
	// 创建日期
	CreateTime int64 `db:"create_time"`
}

// 数据字典
type PermDict struct {
	// Id
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 字典名称
	Name string `db:"name"`
	// 描述
	Remark string `db:"remark"`
	// 创建日期
	CreateTime int64 `db:"create_time"`
}

// 岗位
type PermJob struct {
	// Id
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 岗位名称
	Name string `db:"name"`
	// 岗位状态
	Enabled int16 `db:"enabled"`
	// 岗位排序
	Sort int `db:"sort"`
	// 部门ID
	DeptId int64 `db:"dept_id"`
	// 创建日期
	CreateTime int64 `db:"create_time"`
}

// PermRes
type PermRes struct {
	// 资源ID
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 资源名称
	Name string `db:"name"`
	// 资源类型, 1:页面 2:组件 3:资源
	ResType int16 `db:"res_type"`
	// 上级菜单ID
	Pid int64 `db:"pid"`
	// 深度/层级
	Depth int16 `db:"depth"`
	// 资源键
	Key string `db:"key"`
	// 资源路径
	Path string `db:"path"`
	// 图标
	Icon string `db:"icon"`
	// 排序
	SortNum int `db:"sort_num"`
	// 是否显示到菜单
	IsMenu int16 `db:"is_menu"`
	// 是否启用
	IsEnabled int16 `db:"is_enabled"`
	// 是否禁用
	IsForbidden int64 `db:"is_forbidden"`
	// 创建日期
	CreateTime int64 `db:"create_time"`
	// 组件名称
	ComponentName string `db:"component_name"`
	// 缓存
	Cache string `db:"cache_"`
}

// 角色部门关联
type PermRoleDept struct {
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 角色编号
	RoleId int64 `db:"role_id"`
	// 部门编号
	DeptId int64 `db:"dept_id"`
}

// 角色
type PermRole struct {
	// Id
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 角色代码
	Code string `db:"code"`
	// 名称
	Name string `db:"name"`
	// 角色级别
	Level int `db:"level"`
	// 数据权限
	DataScope string `db:"data_scope"`
	// 备注
	Remark string `db:"remark"`
	// 创建日期
	CreateTime int64 `db:"create_time"`
}

// 角色菜单关联
type PermRoleRes struct {
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 菜单ID
	ResId int64 `db:"res_id"`
	// 角色ID
	RoleId int64 `db:"role_id"`
}

// 系统用户
type PermUser struct {
	// Id
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 用户名
	Username string `db:"username"`
	// 密码
	Password string `db:"password"`
	// 加密盐
	Salt string `db:"salt"`
	// 标志
	Flag int `db:"flag"`
	// 头像
	Avatar string `db:"avatar"`
	// Nickname
	Nickname string `db:"nickname"`
	// Gender
	Gender string `db:"gender"`
	// 邮箱
	Email string `db:"email"`
	// 手机号码
	Phone string `db:"phone"`
	// 部门编号
	DeptId int64 `db:"dept_id"`
	// 岗位编号
	JobId int64 `db:"job_id"`
	// 状态：1启用、0禁用
	Enabled int16 `db:"enabled"`
	// 最后登录的日期
	LastLogin int64 `db:"last_login"`
	// 创建日期
	CreateTime int64 `db:"create_time"`
}

// 用户角色关联
type PermUserRole struct {
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 用户ID
	UserId int64 `db:"user_id"`
	// 角色ID
	RoleId int64 `db:"role_id"`
}
