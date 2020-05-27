package promodel

import (
	"database/sql"
	"go2o/core/domain/interface/pro_model"
	"sort"
)

var _ promodel.ISpecService = new(specServiceImpl)

type specServiceImpl struct {
	rep promodel.IProModelRepo
}

func NewSpecService(rep promodel.IProModelRepo) *specServiceImpl {
	return &specServiceImpl{
		rep: rep,
	}
}

// 获取规格
func (s *specServiceImpl) GetSpec(specId int32) *promodel.Spec {
	return s.rep.GetSpec(specId)
}

// 保存规格
func (s *specServiceImpl) SaveSpec(v *promodel.Spec) (id int32, err error) {
	var i int
	// 如不存在，则新增
	if v.ID <= 0 {
		i, err = s.rep.SaveSpec(v)
		v.ID = int32(i)
		if v == nil {
			return v.ID, err
		}
	}
	// 保存项
	if v.Items != nil {
		v.ItemValues = ""
		for i, iv := range v.Items {
			iv.ProModel = v.ProModel
			iv.SpecId = v.ID
			if i > 0 {
				v.ItemValues += ","
			}
			v.ItemValues += iv.Value
		}
		err = s.saveSpecItems(v.ID, v.Items)
	}
	// 再次保存
	if err == nil {
		_, err = s.rep.SaveSpec(v)
	}
	return v.ID, err
}

// 保存属性项
func (a *specServiceImpl) saveSpecItems(specId int32, items []*promodel.SpecItem) (err error) {
	var i int
	pk := specId
	// 获取存在的项
	old := a.rep.SelectSpecItem("spec_id = $1", pk)
	// 分析当前项目并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*promodel.SpecItem, len(items))
	for _, v := range items {
		currMap[v.ID] = v
	}
	// 筛选出要删除的项
	for _, v := range old {
		if currMap[v.ID] == nil {
			delList = append(delList, v.ID)
		}
	}

	// 删除项
	for _, v := range delList {
		a.rep.DeleteSpecItem(v)
	}
	// 保存项
	for _, v := range items {
		if v.SpecId == 0 {
			v.SpecId = pk
		}
		if v.SpecId == pk {
			if i, err = a.rep.SaveSpecItem(v); err == nil {
				v.ID = int32(i)
			}
		}
	}
	return err
}

// 保存规格项
func (s *specServiceImpl) SaveItem(v *promodel.SpecItem) (int32, error) {
	id, err := s.rep.SaveSpecItem(v)
	return int32(id), err
}

// 删除规格
func (s *specServiceImpl) DeleteSpec(specId int32) error {
	_, err := s.rep.BatchDeleteSpecItem("spec_id= $1", specId)
	if err == nil || err == sql.ErrNoRows {
		err = s.rep.DeleteSpec(specId)
	}
	return err
}

// 删除规格项
func (s *specServiceImpl) DeleteItem(itemId int32) error {
	return s.rep.DeleteSpecItem(itemId)
}

// 获取规格的规格项
func (s *specServiceImpl) GetItems(specId int32) []*promodel.SpecItem {
	return s.rep.SelectSpecItem("spec_id= $1", specId)
}

// 获取产品模型的规格
func (s *specServiceImpl) GetModelSpecs(proModel int32) promodel.SpecList {
	var arr promodel.SpecList
	var items promodel.SpecItemList
	arr = s.rep.SelectSpec("pro_model= $1", proModel)
	for _, v := range arr {
		items = s.GetItems(v.ID)
		sort.Sort(items)
		v.Items = items
	}
	sort.Sort(arr)
	return arr
}
