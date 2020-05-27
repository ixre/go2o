/**
 * Copyright 2015 @ to2.net.
 * name : member_manager.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

type (
	MemberLevel struct {
		Id         int32 `db:"id" auto:"yes" pk:"yes"`
		MerchantId int32 `db:"merchant_id"`
		// 等级值(1,2,4,8,16)
		Value      int32  `db:"value" `
		Name       string `db:"name"`
		RequireExp int32  `db:"require_exp"`
		Enabled    int    `db:"enabled"`
	}

	//todo: 先不撤销, 商户也应对应有等级
	ILevelManager interface {
		// 获取等级设置
		GetLevelSet() []*MemberLevel

		// 获取等级
		GetLevelById(id int32) *MemberLevel

		// 根据等级值获取等级
		GetLevelByValue(value int32) *MemberLevel

		// 获取下一个等级
		GetNextLevel(value int32) *MemberLevel

		// 删除等级
		DeleteLevel(id int32) error

		// 保存等级
		SaveLevel(*MemberLevel) (int32, error)

		// 根据经验值获取等级值
		GetLevelValueByExp(exp int32) int32

		// 初始化默认等级
		InitDefaultLevels() error
	}
)
