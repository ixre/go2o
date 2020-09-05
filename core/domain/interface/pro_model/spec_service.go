package promodel

import (
	"sort"
)

type (
	// 规格
	Spec struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 产品模型
		ProModel int32 `db:"pro_model"`
		// 规格名称
		Name string `db:"name"`
		// 规格项值
		ItemValues string `db:"item_values"`
		// 排列序号
		SortNum int32 `db:"sort_num"`
		// 规格项
		Items SpecItemList `db:"-"`
	}

	// 规格项
	SpecItem struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 规格编号
		SpecId int32 `db:"spec_id"`
		// 产品模型（冗余)
		ProModel int32 `db:"pro_model"`
		// 规格项值
		Value string `db:"value"`
		// 规格项颜色
		Color string `db:"color"`
		// 排列序号
		SortNum int32 `db:"sort_num"`
	}
)

// 规格服务
type ISpecService interface {
	// 获取规格
	GetSpec(specId int32) *Spec
	// 保存规格
	SaveSpec(*Spec) (int32, error)
	// 保存规格项
	SaveItem(*SpecItem) (int32, error)
	// 删除规格
	DeleteSpec(specId int32) error
	// 删除规格项
	DeleteItem(itemId int32) error
	// 获取规格的规格项
	GetItems(specId int32) []*SpecItem
	// 获取产品模型的规格
	GetModelSpecs(proModel int32) SpecList
	// 保存模型的规格
	// SaveModelSpecs(proModel int32, arr []*Spec)
}

var _ sort.Interface = SpecList{}

type SpecList []*Spec

func (s SpecList) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SpecList) Less(i, j int) bool {
	b := s[i].SortNum < s[j].SortNum
	return b || !b && s[i].ID < s[j].ID
}

// Swap swaps the elements with indexes i and j.
func (s SpecList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var _ sort.Interface = SpecItemList{}

type SpecItemList []*SpecItem

func (s SpecItemList) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SpecItemList) Less(i, j int) bool {
	b := s[i].SortNum < s[j].SortNum
	return b || !b && s[i].ID < s[j].ID
}

// Swap swaps the elements with indexes i and j.
func (s SpecItemList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
