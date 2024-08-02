/**
 * Copyright (C) 2009-2024 56X.NET, All rights reserved.
 *
 * name : model_gen.go
 * author : jarrysix
 * date : 2024/07/15 13:44:47
 * description :
 * history :
 */
package model

// RbacDict 数据字典
type RbacDict struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 字典名称
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 描述
	Remark string `json:"remark" db:"remark" gorm:"column:remark" bson:"remark"`
	// 创建日期
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (r RbacDict) TableName() string {
	return "rbac_dict"
}

// RbacDepart 部门
type RbacDepart struct {
	// ID
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 名称
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 上级部门
	Pid int `json:"pid" db:"pid" gorm:"column:pid" bson:"pid"`
	// 状态
	Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
	// 创建日期
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 部门代码
	Code string `json:"code" db:"code" gorm:"column:code" bson:"code"`
}

func (r RbacDepart) TableName() string {
	return "rbac_depart"
}

// RbacRole 角色
type RbacRole struct {
	// ID
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 名称
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 角色级别
	Level int `json:"level" db:"level" gorm:"column:level" bson:"level"`
	// 数据权限
	DataScope string `json:"dataScope" db:"data_scope" gorm:"column:data_scope" bson:"dataScope"`
	// 备注
	Remark string `json:"remark" db:"remark" gorm:"column:remark" bson:"remark"`
	// 创建日期
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 角色代码
	Code string `json:"code" db:"code" gorm:"column:code" bson:"code"`
}

func (r RbacRole) TableName() string {
	return "rbac_role"
}

// RbacJob 岗位
type RbacJob struct {
	// ID
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 岗位名称
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 岗位状态
	Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
	// 岗位排序
	Sort int `json:"sort" db:"sort" gorm:"column:sort" bson:"sort"`
	// 部门ID
	DeptId int `json:"deptId" db:"dept_id" gorm:"column:dept_id" bson:"deptId"`
	// 创建日期
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (r RbacJob) TableName() string {
	return "rbac_job"
}

// RbacLoginLog 用户登录日志
type RbacLoginLog struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 用户编号
	UserId int `json:"userId" db:"user_id" gorm:"column:user_id" bson:"userId"`
	// 登录IP地址
	Ip string `json:"ip" db:"ip" gorm:"column:ip" bson:"ip"`
	// 是否成功
	IsSuccess int `json:"isSuccess" db:"is_success" gorm:"column:is_success" bson:"isSuccess"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (r RbacLoginLog) TableName() string {
	return "rbac_login_log"
}

// RbacRes RbacRes
type RbacRes struct {
	// 资源ID
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 资源名称
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 资源类型, 0: 目录 1: 资源　2: 菜单  3:　 按钮
	ResType int `json:"resType" db:"res_type" gorm:"column:res_type" bson:"resType"`
	// 上级菜单ID
	Pid int `json:"pid" db:"pid" gorm:"column:pid" bson:"pid"`
	// 资源键
	ResKey string `json:"resKey" db:"res_key" gorm:"column:res_key" bson:"resKey"`
	// 资源路径
	Path string `json:"path" db:"path" gorm:"column:path" bson:"path"`
	// 图标
	Icon string `json:"icon" db:"icon" gorm:"column:icon" bson:"icon"`
	// 排序
	SortNum int `json:"sortNum" db:"sort_num" gorm:"column:sort_num" bson:"sortNum"`
	// 是否显示到菜单中
	IsMenu int `json:"isMenu" db:"is_menu" gorm:"column:is_menu" bson:"isMenu"`
	// 是否启用
	IsEnabled int `json:"isEnabled" db:"is_enabled" gorm:"column:is_enabled" bson:"isEnabled"`
	// 创建日期
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 组件路径
	ComponentName string `json:"componentName" db:"component_name" gorm:"column:component_name" bson:"componentName"`
	// 深度/层级
	Depth int `json:"depth" db:"depth" gorm:"column:depth" bson:"depth"`
	// 是否禁止
	IsForbidden int `json:"isForbidden" db:"is_forbidden" gorm:"column:is_forbidden" bson:"isForbidden"`
	// 应用(系统)序号
	AppIndex int `json:"appIndex" db:"app_index" gorm:"column:app_index" bson:"appIndex"`
}

func (r RbacRes) TableName() string {
	return "rbac_res"
}

// RbacRoleDept 角色部门关联
type RbacRoleDept struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 角色编号
	RoleId int `json:"roleId" db:"role_id" gorm:"column:role_id" bson:"roleId"`
	// 部门编号
	DeptId int `json:"deptId" db:"dept_id" gorm:"column:dept_id" bson:"deptId"`
}

func (r RbacRoleDept) TableName() string {
	return "rbac_role_dept"
}

// RbacRoleRes 角色菜单关联
type RbacRoleRes struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 菜单ID
	ResId int `json:"resId" db:"res_id" gorm:"column:res_id" bson:"resId"`
	// 角色ID
	RoleId int `json:"roleId" db:"role_id" gorm:"column:role_id" bson:"roleId"`
	// 权限值, 1:增加  2:删除 4: 更新
	PermFlag int `json:"permFlag" db:"perm_flag" gorm:"column:perm_flag" bson:"permFlag"`
}

func (r RbacRoleRes) TableName() string {
	return "rbac_role_res"
}

// RbacUser 系统用户
type RbacUser struct {
	// ID
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 用户名
	Username string `json:"username" db:"username" gorm:"column:username" bson:"username"`
	// 密码
	Password string `json:"password" db:"password" gorm:"column:password" bson:"password"`
	// 加密盐
	Salt string `json:"salt" db:"salt" gorm:"column:salt" bson:"salt"`
	// 标志
	Flag int `json:"flag" db:"flag" gorm:"column:flag" bson:"flag"`
	// 头像
	ProfilePhoto string `json:"profilePhoto" db:"profile_photo" gorm:"column:profile_photo" bson:"avatar"`
	// 姓名
	Nickname string `json:"nickname" db:"nickname" gorm:"column:nickname" bson:"nickname"`
	// 性别
	Gender int `json:"gender" db:"gender" gorm:"column:gender" bson:"gender"`
	// 邮箱
	Email string `json:"email" db:"email" gorm:"column:email" bson:"email"`
	// 手机号码
	Phone string `json:"phone" db:"phone" gorm:"column:phone" bson:"phone"`
	// 部门编号
	DeptId int `json:"deptId" db:"dept_id" gorm:"column:dept_id" bson:"deptId"`
	// 岗位编号
	JobId int `json:"jobId" db:"job_id" gorm:"column:job_id" bson:"jobId"`
	// 状态：1启用、0禁用
	Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
	// 最后登录的日期
	LastLogin int `json:"lastLogin" db:"last_login" gorm:"column:last_login" bson:"lastLogin"`
	// 创建日期
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (r RbacUser) TableName() string {
	return "rbac_user"
}

// RbacUserRole 用户角色关联
type RbacUserRole struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 用户ID
	UserId int `json:"userId" db:"user_id" gorm:"column:user_id" bson:"userId"`
	// 角色ID
	RoleId int `json:"roleId" db:"role_id" gorm:"column:role_id" bson:"roleId"`
}

func (r RbacUserRole) TableName() string {
	return "rbac_user_role"
}

// RbacDictDetail 数据字典详情
type RbacDictDetail struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 字典标签
	Label string `json:"label" db:"label" gorm:"column:label" bson:"label"`
	// 字典值
	Value string `json:"value" db:"value" gorm:"column:value" bson:"value"`
	// 排序
	Sort string `json:"sort" db:"sort" gorm:"column:sort" bson:"sort"`
	// 字典id
	DictId int `json:"dictId" db:"dict_id" gorm:"column:dict_id" bson:"dictId"`
	// 创建日期
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (r RbacDictDetail) TableName() string {
	return "rbac_dict_detail"
}
