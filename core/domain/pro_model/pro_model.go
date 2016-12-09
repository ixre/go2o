package promodel

import "go2o/core/domain/interface/pro_model"

var _ promodel.IModel = new(modelImpl)

type modelImpl struct {
	rep   promodel.IProModelRepo
	value *promodel.ProModel
}

func NewModel(v *promodel.ProModel, rep promodel.IProModelRepo) promodel.IModel {
	return &modelImpl{
		rep:   rep,
		value: v,
	}
}

// 获取聚合根编号
func (m *modelImpl) GetAggregateRootId() int32 {
	return m.value.Id
}

// 获取值
func (m *modelImpl) Value() *promodel.ProModel {
	return m.value
}

// 保存
func (m *modelImpl) Save() (int32, error) {
	return m.GetAggregateRootId(), nil
}

// 是否启用
func (m *modelImpl) Enabled() bool {
	return m.value.Enabled == 1
}

// 获取关联的品牌编号
func (m *modelImpl) Brands() []*promodel.ProBrand {
	return m.rep.BrandService().Brands(m.GetAggregateRootId())
}

// 关联品牌
func (m *modelImpl) SaveBrands(brandId []int32) error {
	return m.rep.BrandService().SetBrands(m.GetAggregateRootId(), brandId)
}

// 获取属性
func (m *modelImpl) Attrs() []*promodel.Attr {
	return nil
}

// 保存属性
func (m *modelImpl) SaveAttrs([]*promodel.Attr) error {
	return nil
}

// 获取规格
func (m *modelImpl) Specs() []*promodel.Spec {
	return nil
}

// 保存规格
func (m *modelImpl) SaveSpecs([]*promodel.Spec) error {
	return nil
}
