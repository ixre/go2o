/**
 * Copyright 2015 @ z3q.net.
 * name : conf_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

import (
	"errors"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
)

//todo: 这里引用IMemberRep似乎有问题

var _ merchant.ILevelManager = new(LevelManager)

type LevelManager struct {
	_rep        member.IMemberRep
	_merchantId int
	_levelSet   []*merchant.MemberLevel
}

func NewLevelManager(merchantId int, rep member.IMemberRep) merchant.ILevelManager {
	return &LevelManager{
		_merchantId: merchantId,
		_rep:        rep,
	}
}

// 初始化默认等级
func (this *LevelManager) InitDefaultLevels() error {
	if len(this.GetLevelSet()) != 0 {

		return errors.New("已经存在数据，无法初始化!")
	}
	var arr []*merchant.MemberLevel = []*merchant.MemberLevel{
		&merchant.MemberLevel{
			MerchantId: this._merchantId,
			Name:       "普通会员",
			RequireExp: 0,
			Value:      1,
			Enabled:    1,
		},
		&merchant.MemberLevel{
			MerchantId: this._merchantId,
			Name:       "铜牌会员",
			RequireExp: 100,
			Value:      2,
			Enabled:    1,
		},
		&merchant.MemberLevel{
			MerchantId: this._merchantId,
			Name:       "银牌会员",
			RequireExp: 500,
			Value:      3,
			Enabled:    1,
		},
		&merchant.MemberLevel{
			MerchantId: this._merchantId,
			Name:       "金牌会员",
			RequireExp: 1200,
			Value:      4,
			Enabled:    1,
		},
		&merchant.MemberLevel{
			MerchantId: this._merchantId,
			Name:       "白金会员",
			RequireExp: 1500,
			Value:      5,
			Enabled:    1,
		},
	}

	for _, v := range arr {
		v.Id, _ = this.SaveLevel(v)
	}
	return nil
}

// 获取等级设置
func (this *LevelManager) GetLevelSet() []*merchant.MemberLevel {
	if this._levelSet == nil {
		// 已经排好序
		this._levelSet = this._rep.GetMemberLevels(this._merchantId)
	}
	return this._levelSet
}

// 获取等级
func (this *LevelManager) GetLevelById(id int) *merchant.MemberLevel {
	for _, v := range this.GetLevelSet() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 根据等级值获取等级
func (this *LevelManager) GetLevelByValue(value int) *merchant.MemberLevel {
	for _, v := range this.GetLevelSet() {
		if v.Value == value {
			return v
		}
	}
	return nil
}

// 获取下一个等级
func (this *LevelManager) GetNextLevel(value int) *merchant.MemberLevel {
	return this._rep.GetNextLevel(this._merchantId, value)
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
		return this._rep.DeleteMemberLevel(this._merchantId, id)
	}
	return errors.New("no such record")
}

// 保存等级
func (this *LevelManager) SaveLevel(v *merchant.MemberLevel) (int, error) {
	v.MerchantId = this._merchantId
	// 如果新增（非初始化）等级自动设置值
	if v.Id <= 0 && len(this._levelSet) == 0 {
		v.Value = this.getMaxLevelValue() + 1
	}
	this._levelSet = nil
	return this._rep.SaveMemberLevel(this._merchantId, v)
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
	var lv *merchant.MemberLevel
	var levelVal int
	for i := len(this.GetLevelSet()); i > 0; i-- {
		lv = this.GetLevelSet()[i-1]
		if exp >= lv.RequireExp && lv.Value > levelVal {
			levelVal = lv.Value
		}
	}
	return levelVal
}
