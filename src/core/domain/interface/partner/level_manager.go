/**
 * Copyright 2015 @ S1N1 Team.
 * name : member_manager.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package partner

import "go2o/src/core/domain/interface/valueobject"

type ILevelManager interface {
	// 获取等级设置
	GetLevelSet() []*valueobject.MemberLevel

	// 获取等级
	GetLevelById(id int) *valueobject.MemberLevel

	// 根据等级值获取等级
	GetLevelByValue(value int) *valueobject.MemberLevel

	// 获取下一个等级
	GetNextLevel(value int) *valueobject.MemberLevel

	// 删除等级
	DeleteLevel(id int) error

	// 保存等级
	SaveLevel(*valueobject.MemberLevel) (int, error)

	// 根据经验值获取等级值
	GetLevelValueByExp(exp int) int

	// 初始化默认等级
	InitDefaultLevels() error
}
