/**
 * Copyright 2015 @ z3q.net.
 * name : member_manager.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

type (
	MemberLevel struct {
		Id         int64 `db:"id" auto:"yes" pk:"yes"`
		MerchantId int64 `db:"merchant_id"`
		// 等级值(1,2,4,8,16)
		Value      int64  `db:"value" `
		Name       string `db:"name"`
		RequireExp int64  `db:"require_exp"`
		Enabled    int    `db:"enabled"`
	}

	//todo: 先不撤销, 商户也应对应有等级
	ILevelManager interface {
		// 获取等级设置
		GetLevelSet() []*MemberLevel

		// 获取等级
		GetLevelById(id int64) *MemberLevel

		// 根据等级值获取等级
		GetLevelByValue(value int64) *MemberLevel

		// 获取下一个等级
		GetNextLevel(value int64) *MemberLevel

		// 删除等级
		DeleteLevel(id int64) error

		// 保存等级
		SaveLevel(*MemberLevel) (int64, error)

		// 根据经验值获取等级值
		GetLevelValueByExp(exp int64) int64

		// 初始化默认等级
		InitDefaultLevels() error
	}
)
