package model

// 系统用户
type PermUser struct {
	// ID
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 用户名
	User string `db:"user"`
	// 密码
	Pwd string `db:"pwd"`
	// 标志
	Flag int `db:"flag"`
	// 头像
	Avatar string `db:"avatar"`
	// NickName
	NickName string `db:"nick_name"`
	// Sex
	Sex string `db:"sex"`
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
