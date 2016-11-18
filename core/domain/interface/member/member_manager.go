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
	RegisterModeNormal         = 1 // 正常注册
	RegisterModeClosed         = 2 // 关闭注册
	RegisterModeMustRedirect   = 3 // 必须直接注册
	RegisterModeMustInvitation = 4 // 必须邀请注册

)

type (

	// 会员服务
	IMemberManager interface {
		// 等级服务
		LevelManager() ILevelManager

		// 检查手机绑定,同时检查手机格式
		CheckPhoneBind(phone string, memberId int64) error

		// 检查注册信息是否正确
		PrepareRegister(v *Member, pro *Profile,
			invitationCode string) (invitationId int, err error)

		// IDocManager()IDocManager
	}

	//会员等级
	Level struct {
		//编号
		Id int64 `db:"id" auto:"yes" pk:"yes"`
		//等级名称
		Name string `db:"name"`
		//需要经验值
		RequireExp int64 `db:"require_exp"`
		// 可编程等级签名,可根据此签名来进行编程
		ProgramSignal string `db:"program_signal"`
		//是否启用
		Enabled int `db:"enabled"`
		//是否为正式会员/非正式
		IsOfficial int `db:"is_official"`
	}

	ILevelManager interface {
		// 获取等级设置
		GetLevelSet() []*Level

		// 获取最高已启用的等级
		GetHighestLevel() *Level

		// 获取等级
		GetLevelById(id int64) *Level

		// 根据可编程字符获取会员等级
		GetLevelByProgramSign(sign string) *Level

		// 获取下一个等级
		GetNextLevelById(id int64) *Level

		// 删除等级
		DeleteLevel(id int64) error

		// 保存等级
		SaveLevel(*Level) (int64, error)

		// 根据经验值获取等级值
		GetLevelIdByExp(exp int64) int64
	}
)
