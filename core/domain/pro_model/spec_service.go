package promodel

import (
	"database/sql"
	"sort"

	promodel "github.com/ixre/go2o/core/domain/interface/pro_model"
)

var _ promodel.ISpecService = new(specServiceImpl)

type specServiceImpl struct {
	rep promodel.IProductModelRepo
}

func NewSpecService(rep promodel.IProductModelRepo) *specServiceImpl {
	return &specServiceImpl{
		rep: rep,
	}
}

// 获取规格
func (s *specServiceImpl) GetSpec(specId int) *promodel.Spec {
	return s.rep.GetSpec(specId)
}

// 保存规格
func (s *specServiceImpl) SaveSpec(v *promodel.Spec) (id int, err error) {
	var i int
	// 如不存在，则新增
	if v.Id <= 0 {
		i, err = s.rep.SaveSpec(v)
		v.Id = i
		if v == nil {
			return v.Id, err
		}
	}
	// 保存项
	if v.Items != nil {
		v.ItemValues = ""
		for i, iv := range v.Items {
			iv.ModelId = v.ModelId
			iv.SpecId = v.Id
			if i > 0 {
				v.ItemValues += ","
			}
			v.ItemValues += iv.Value
		}
		err = s.saveSpecItems(v.Id, v.Items)
	}
	// 再次保存
	if err == nil {
		_, err = s.rep.SaveSpec(v)
	}
	return v.Id, err
}

// 保存属性项
func (a *specServiceImpl) saveSpecItems(specId int, items []*promodel.SpecItem) (err error) {
	var i int
	pk := specId
	// 获取存在的项
	old := a.rep.SelectSpecItem("spec_id = $1", pk)
	// 分析当前项目并加入到MAP中
	delList := []int{}
	currMap := make(map[int]*promodel.SpecItem, len(items))
	for _, v := range items {
		currMap[v.Id] = v
	}
	// 筛选出要删除的项
	for _, v := range old {
		if currMap[v.Id] == nil {
			delList = append(delList, v.Id)
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
				v.Id = int(i)
			}
		}
	}
	return err
}

// 保存规格项
func (s *specServiceImpl) SaveItem(v *promodel.SpecItem) (int, error) {
	id, err := s.rep.SaveSpecItem(v)
	return int(id), err
}

// 删除规格
func (s *specServiceImpl) DeleteSpec(specId int) error {
	_, err := s.rep.BatchDeleteSpecItem("spec_id= $1", specId)
	if err == nil || err == sql.ErrNoRows {
		err = s.rep.DeleteSpec(specId)
	}
	return err
}

// 删除规格项
func (s *specServiceImpl) DeleteItem(itemId int) error {
	return s.rep.DeleteSpecItem(itemId)
}

// 获取规格的规格项
func (s *specServiceImpl) GetItems(specId int) promodel.SpecItemList {
	var items promodel.SpecItemList = s.rep.SelectSpecItem("spec_id= $1", specId)
	sort.Sort(items)
	return items

}

// 获取产品模型的规格
func (s *specServiceImpl) GetModelSpecs(proModel int) promodel.SpecList {
	var arr promodel.SpecList = s.rep.SelectSpec("prod_model= $1", proModel)
	for _, v := range arr {
		v.Items = s.GetItems(v.Id)
	}
	sort.Sort(arr)
	return arr
}
