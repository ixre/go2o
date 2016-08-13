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
	"go2o/core/domain/interface/merchant"
)

//todo: 这里引用IMemberRep似乎有问题

var _ merchant.ILevelManager = new(LevelManager)

type LevelManager struct {
	_mchRep     merchant.IMerchantRep
	_merchantId int
	_levelSet   []*merchant.MemberLevel
}

func NewLevelManager(merchantId int, rep merchant.IMerchantRep) merchant.ILevelManager {
	return &LevelManager{
		_merchantId: merchantId,
		_mchRep:     rep,
	}
}

// 初始化默认等级
func (l *LevelManager) InitDefaultLevels() error {
	if len(l.GetLevelSet()) != 0 {

		return errors.New("已经存在数据，无法初始化!")
	}
	var arr []*merchant.MemberLevel = []*merchant.MemberLevel{
		&merchant.MemberLevel{
			MerchantId: l._merchantId,
			Name:       "普通会员",
			RequireExp: 0,
			Value:      1,
			Enabled:    1,
		},
		&merchant.MemberLevel{
			MerchantId: l._merchantId,
			Name:       "铜牌会员",
			RequireExp: 100,
			Value:      2,
			Enabled:    1,
		},
		&merchant.MemberLevel{
			MerchantId: l._merchantId,
			Name:       "银牌会员",
			RequireExp: 500,
			Value:      3,
			Enabled:    1,
		},
		&merchant.MemberLevel{
			MerchantId: l._merchantId,
			Name:       "金牌会员",
			RequireExp: 1200,
			Value:      4,
			Enabled:    1,
		},
		&merchant.MemberLevel{
			MerchantId: l._merchantId,
			Name:       "白金会员",
			RequireExp: 1500,
			Value:      5,
			Enabled:    1,
		},
	}

	for _, v := range arr {
		v.Id, _ = l.SaveLevel(v)
	}
	return nil
}

// 获取等级设置
func (l *LevelManager) GetLevelSet() []*merchant.MemberLevel {
	if l._levelSet == nil {
		// 已经排好序
		l._levelSet = l._mchRep.GetMemberLevels(l._merchantId)
	}
	return l._levelSet
}

// 获取等级
func (l *LevelManager) GetLevelById(id int) *merchant.MemberLevel {
	for _, v := range l.GetLevelSet() {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// 根据等级值获取等级
func (l *LevelManager) GetLevelByValue(value int) *merchant.MemberLevel {
	for _, v := range l.GetLevelSet() {
		if v.Value == value {
			return v
		}
	}
	return nil
}

// 获取下一个等级
func (l *LevelManager) GetNextLevel(value int) *merchant.MemberLevel {
	return l._mchRep.GetNextLevel(l._merchantId, value)
}

// 删除等级
func (l *LevelManager) DeleteLevel(id int) error {
	var exists bool = true
	if l._levelSet != nil {
		exists = false
		for i, v := range l._levelSet {
			if v.Id == id {
				exists = true
				l._levelSet = append(l._levelSet[:i], l._levelSet[i+1:]...)
				break
			}
		}
	}
	if exists {
		//todo: 更新会员的等级到下一级
		return l._mchRep.DeleteMemberLevel(l._merchantId, id)
	}
	return errors.New("no such record")
}

// 保存等级
func (l *LevelManager) SaveLevel(v *merchant.MemberLevel) (int, error) {
	v.MerchantId = l._merchantId
	// 如果新增（非初始化）等级自动设置值
	if v.Id <= 0 && len(l._levelSet) == 0 {
		v.Value = l.getMaxLevelValue() + 1
	}
	l._levelSet = nil
	return l._mchRep.SaveMemberLevel(l._merchantId, v)
}

// 获取最大的等级值
func (l *LevelManager) getMaxLevelValue() int {
	var k = 0
	for _, v := range l.GetLevelSet() {
		if v.Value > k {
			k = v.Value
		}
	}
	return k
}

// 根据经验值获取等级
func (l *LevelManager) GetLevelValueByExp(exp int) int {
	var lv *merchant.MemberLevel
	var levelVal int
	for i := len(l.GetLevelSet()); i > 0; i-- {
		lv = l.GetLevelSet()[i-1]
		if exp >= lv.RequireExp && lv.Value > levelVal {
			levelVal = lv.Value
		}
	}
	return levelVal
}
