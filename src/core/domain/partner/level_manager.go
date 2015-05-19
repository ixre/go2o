/**
 * Copyright 2015 @ S1N1 Team.
 * name : conf_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package partner

import (
	"errors"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/valueobject"
)

//todo: 这里引用IMemberRep似乎有问题

var _ partner.ILevelManager = new(LevelManager)

type LevelManager struct {
	_rep       member.IMemberRep
	_partnerId int
	_levelSet  []*valueobject.MemberLevel
}

func NewLevelManager(partnerId int, rep member.IMemberRep) partner.ILevelManager {
	return &LevelManager{
		_partnerId: partnerId,
		_rep:       rep,
	}
}

func (this *LevelManager) insertDefaultLevels() []*valueobject.MemberLevel {
	var arr []*valueobject.MemberLevel = []*valueobject.MemberLevel{
		&valueobject.MemberLevel{
			PartnerId:  this._partnerId,
			Name:       "普通会员",
			RequireExp: 0,
			Value:      1,
			Enabled:    1,
		},
		&valueobject.MemberLevel{
			PartnerId:  this._partnerId,
			Name:       "铜牌会员",
			RequireExp: 100,
			Value:      2,
			Enabled:    1,
		},
		&valueobject.MemberLevel{
			PartnerId:  this._partnerId,
			Name:       "银牌会员",
			RequireExp: 500,
			Value:      3,
			Enabled:    1,
		},
		&valueobject.MemberLevel{
			PartnerId:  this._partnerId,
			Name:       "金牌会员",
			RequireExp: 1200,
			Value:      4,
			Enabled:    1,
		},
		&valueobject.MemberLevel{
			PartnerId:  this._partnerId,
			Name:       "白金会员",
			RequireExp: 1500,
			Value:      5,
			Enabled:    1,
		},
	}

	for _, v := range arr {
		v.Id, _ = this.SaveLevel(v)
	}
	return arr
}

// 获取等级设置
func (this *LevelManager) GetLevelSet() []*valueobject.MemberLevel {
	if this._levelSet == nil {
		this._levelSet = this._rep.GetMemberLevels(this._partnerId)
		if len(this._levelSet) == 0 {
			this._levelSet = this.insertDefaultLevels()
		}
	}
	return this._levelSet
}

// 获取等级
func (this *LevelManager) GetLevelById(id int) *valueobject.MemberLevel {
	for _, v := range this.GetLevelSet() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 根据等级值获取等级
func (this *LevelManager) GetLevelByValue(value int) *valueobject.MemberLevel {
	for _, v := range this.GetLevelSet() {
		if v.Value == value {
			return v
		}
	}
	return nil
}

// 获取下一个等级
func (this *LevelManager) GetNextLevel(value int) *valueobject.MemberLevel {
	return this._rep.GetNextLevel(this._partnerId, value)
}

// 删除等级
func (this *LevelManager) DeleteLevel(id int) error {
	var exists bool = true
	if this._levelSet != nil {
		exists = false
		for i, v := range this._levelSet {
			if v.Id == id {
				exists = true
				this._levelSet = append(this._levelSet[:i], this._levelSet[i+1:]...)
				break
			}
		}
	}
	if exists {
		//todo: 更新会员的等级到下一级
		return this._rep.DeleteMemberLevel(this._partnerId, id)
	}
	return errors.New("no such record")
}

// 保存等级
func (this *LevelManager) SaveLevel(v *valueobject.MemberLevel) (int, error) {
	v.PartnerId = this._partnerId
	// 如果新增（非初始化）等级自动设置值
	if v.Id <= 0 && len(this._levelSet) == 0 {
		v.Value = this.getMaxLevelValue() + 1
	}
	this._levelSet = nil
	return this._rep.SaveMemberLevel(this._partnerId, v)
}

// 获取最大的等级值
func (this *LevelManager) getMaxLevelValue() int {
	var k = 0
	for _, v := range this.GetLevelSet() {
		if v.Value > k {
			k = v.Value
		}
	}
	return k
}

// 根据经验值获取等级
func (this *LevelManager) GetLevelValueByExp(exp int) int {
	return this._rep.GetLevelValueByExp(exp)
}
