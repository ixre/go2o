package promodel

import (
	"sort"
)

type (
	// 规格
	Spec struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 产品模型
		ModelId int `db:"prod_model"`
		// 规格名称
		Name string `db:"name"`
		// 规格项值
		ItemValues string `db:"item_values"`
		// 排列序号
		SortNum int `db:"sort_num"`
		// 规格项
		Items SpecItemList `db:"-"`
	}

	// 规格项
	SpecItem struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 规格编号
		SpecId int `db:"spec_id"`
		// 产品模型（冗余)
		ModelId int `db:"prod_model"`
		// 规格项值
		Value string `db:"value"`
		// 规格项颜色
		Color string `db:"color"`
		// 排列序号
		SortNum int `db:"sort_num"`
	}
)

// 规格服务
type ISpecService interface {
	// 获取规格
	GetSpec(specId int) *Spec
	// 保存规格
	SaveSpec(*Spec) (int, error)
	// 保存规格项
	SaveItem(*SpecItem) (int, error)
	// 删除规格
	DeleteSpec(specId int) error
	// 删除规格项
	DeleteItem(itemId int) error
	// 获取规格的规格项
	GetItems(specId int) SpecItemList
	// 获取产品模型的规格
	GetModelSpecs(proModel int) SpecList
	// 保存模型的规格
	// SaveModelSpecs(proModel int, arr []*Spec)
}

var _ sort.Interface = SpecList{}

type SpecList []*Spec

func (s SpecList) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SpecList) Less(i, j int) bool {
	b := s[i].SortNum - s[j].SortNum
	return b < 0 || b == 0 && s[i].Id < s[j].Id
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
	b := s[i].SortNum - s[j].SortNum
	return b < 0 || b == 0 && s[i].Id < s[j].Id
}

// Swap swaps the elements with indexes i and j.
func (s SpecItemList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
