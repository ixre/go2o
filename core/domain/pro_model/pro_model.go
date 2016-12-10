package promodel

import "go2o/core/domain/interface/pro_model"

var _ promodel.IModel = new(modelImpl)

type modelImpl struct {
    rep          promodel.IProModelRepo
    value        *promodel.ProModel
    attrService  promodel.IAttrService
    specService  promodel.ISpecService
    brandService promodel.IBrandService
    tmpBrandIds  []int32
}

func NewModel(v *promodel.ProModel, rep promodel.IProModelRepo,
attrService promodel.IAttrService, specService promodel.ISpecService,
brandService promodel.IBrandService) promodel.IModel {
    return &modelImpl{
        rep:   rep,
        value: v,
        attrService:attrService,
        specService:specService,
        brandService:brandService,
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

// 获取属性
func (m *modelImpl) Attrs() []*promodel.Attr {
    if m.value.Attrs == nil{
        m.value.Attrs = m.attrService.GetModelAttrs(m.GetAggregateRootId())
    }
    return m.value.Attrs
}

// 保存属性
func (m *modelImpl) SetAttrs(a []*promodel.Attr) error {
    if a == nil || len(a) == 0{
        return promodel.ErrEmptyAttrArray
    }
    m.value.Attrs = a
    return nil
}

// 获取规格
func (m *modelImpl) Specs() []*promodel.Spec {
    if m.value.Specs == nil{
        m.value.Specs = m.specService.GetModelSpecs(m.GetAggregateRootId())
    }
    return m.value.Specs
}

// 保存规格
func (m *modelImpl) SetSpecs(s []*promodel.Spec) error {
    if s == nil || len(s) == 0{
        return promodel.ErrEmptySpecArray
    }
    m.value.Specs = s
    return nil
}

// 获取关联的品牌编号
func (m *modelImpl) Brands() []*promodel.ProBrand {
    return m.rep.BrandService().Brands(m.GetAggregateRootId())
}

// 关联品牌
func (m *modelImpl) SetBrands(brandId []int32) error {
    m.tmpBrandIds = brandId
    return nil
}

// 保存
func (m *modelImpl) Save() (i int32, err error) {
    // 新增模型
    if m.GetAggregateRootId() <= 0 {
        id, err := m.rep.SaveProModel(m.value)
        if err == nil {
            m.value.Id = int32(id)
        } else {
            return 0, err
        }
    }
    // 保存品牌
    if m.tmpBrandIds != nil {
        err = m.rep.BrandService().SetBrands(m.GetAggregateRootId(),
            m.tmpBrandIds)
    }
    // 保存属性
    if m.value.Attrs != nil {
        m.value.AttrStr = ""
        for i, v := range m.value.Attrs {
            if i > 0 {
                m.value.AttrStr += ","
            }
            m.value.AttrStr += v.Name
        }
        //err = m.rep.AttrService().SaveModelAttrs(m.GetAggregateRootId(),
        //    m.value.Attrs)
    }
    // 保存规格
    if m.value.Specs != nil {
        m.value.SpecStr = ""
        for i, v := range m.value.Specs {
            if i > 0 {
                m.value.SpecStr += ","
            }
            m.value.SpecStr += v.Name
        }
        //err = m.rep.SpecService().SaveModelSpecs(m.GetAggregateRootId(),
        //    m.value.Attrs)
    }
    return m.GetAggregateRootId(), err
}

// 是否启用
func (m *modelImpl) Enabled() bool {
    return m.value.Enabled == 1
}