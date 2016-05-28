/**
 * Copyright 2015 @ z3q.net.
 * name : level_manager
 * author : jarryliu
 * date : 2016-05-26 11:33
 * description :
 * history :
 */
package member

const (
	RegisterModeClosed         = 0 // 关闭注册
	RegisterModeNormal         = 1 // 正常注册
	RegisterModeMustInvitation = 2 // 必须邀请注册
	RegisterModeMustRedirect   = 3 // 必须直接注册
)

type (

	// 会员服务
	IMemberManager interface {
		// 等级服务
		LevelManager() ILevelManager

		// 检测注册权限
		RegisterPerm(invitation bool) error

		// IDocManager()IDocManager
	}

	//会员等级
	// todo: value要删除。
	Level struct {
		//编号
		Id int `db:"id" auto:"yes" pk:"yes"`

		//等级名称
		Name string `db:"name"`
		//需要经验值
		RequireExp int `db:"require_exp"`
		// 可编程等级签名,可根据此签名来进行编程
		ProgramSignal string `db:"program_signal"`
		//是否启用
		Enabled int `db:"enabled"`

		//等级值(1,2,4,8,16)
		Value int `db:"-"`
	}

	ILevelManager interface {
		// 获取等级设置
		GetLevelSet() []*Level

		// 获取等级
		GetLevelById(id int) *Level

		// 获取下一个等级
		GetNextLevelById(int int) *Level

		// 删除等级
		DeleteLevel(id int) error

		// 保存等级
		SaveLevel(*Level) (int, error)

		// 根据经验值获取等级值
		GetLevelIdByExp(exp int) int
	}
)
