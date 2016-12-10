package promodel

import (
	"database/sql"
	"go2o/core/domain/interface/pro_model"
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
func (s *specServiceImpl) SaveSpec(v *promodel.Spec) (int32, error) {
	id, err := s.rep.SaveSpec(v)
	return int32(id), err
}

// 保存规格项
func (s *specServiceImpl) SaveItem(v *promodel.SpecItem) (int32, error) {
	id, err := s.rep.SaveSpecItem(v)
	return int32(id), err
}

// 删除规格
func (s *specServiceImpl) DeleteSpec(specId int32) error {
	_, err := s.rep.BatchDeleteSpecItem("spec_id=?", specId)
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
	return s.rep.SelectSpecItem("spec_id=?", specId)
}

// 获取产品模型的规格
func (s *specServiceImpl) GetModelSpecs(proModel int32) []*promodel.Spec {
	return s.rep.SelectSpec("pro_model=?", proModel)
}

// 保存模型的规格
func (s *specServiceImpl) SaveModelSpecs(proModel int32, arr []*promodel.Spec) {
	//return nil
}
