package promodel

import (
	"errors"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/infrastructure/format"
)

var _ promodel.IProductModel = new(modelImpl)

type modelImpl struct {
	rep          promodel.IProductModelRepo
	value        *promodel.ProductModel
	attrService  promodel.IAttrService
	specService  promodel.ISpecService
	brandService promodel.IBrandService
}

func NewModel(v *promodel.ProductModel, rep promodel.IProductModelRepo,
	attrService promodel.IAttrService, specService promodel.ISpecService,
	brandService promodel.IBrandService) promodel.IProductModel {
	return &modelImpl{
		rep:          rep,
		value:        v,
		attrService:  attrService,
		specService:  specService,
		brandService: brandService,
	}
}

// 获取聚合根编号
func (m *modelImpl) GetAggregateRootId() int32 {
	return m.value.ID
}

// 获取值
func (m *modelImpl) Value() *promodel.ProductModel {
	return m.value
}


// 是否启用
func (m *modelImpl) SetValue(model *promodel.ProductModel) error {
	if len(model.Name) == 0{
		return errors.New("model name length")
	}
	 m.value.Enabled = model.Enabled
	 m.value.Name = model.Name
	 return nil
}

// 获取属性
func (m *modelImpl) Attrs() []*promodel.Attr {
	if m.value.Attrs == nil {
		m.value.Attrs = m.attrService.GetModelAttrs(m.GetAggregateRootId())
	}
	return m.value.Attrs
}

// 保存属性
func (m *modelImpl) SetAttrs(a []*promodel.Attr) error {
	if a == nil {
		return promodel.ErrEmptyAttrArray
	}
	m.value.Attrs = a
	return nil
}

// 获取规格
func (m *modelImpl) Specs() []*promodel.Spec {
	if m.value.Specs == nil {
		m.value.Specs = m.specService.GetModelSpecs(m.GetAggregateRootId())
	}
	return m.value.Specs
}

// 保存规格
func (m *modelImpl) SetSpecs(s []*promodel.Spec) error {
	if s == nil {
		return promodel.ErrEmptySpecArray
	}
	m.value.Specs = s
	return nil
}

// 获取关联的品牌编号
func (m *modelImpl) Brands() []*promodel.ProductBrand {
	return m.rep.BrandService().Brands(m.GetAggregateRootId())
}

// 关联品牌
func (m *modelImpl) SetBrands(brandId []int32) error {
	if brandId == nil {
		return promodel.ErrEmptyBrandArray
	}
	m.value.BrandArray = brandId
	return nil
}

// 保存
func (m *modelImpl) Save() (id int32, err error) {
	var i int
	// 新增模型
	if m.GetAggregateRootId() <= 0 {
		i, err = m.rep.SaveProModel(m.value)
		if err == nil {
			m.value.ID = int32(i)
		} else {
			return 0, err
		}
	}
	// 保存品牌
	if m.value.BrandArray != nil {
		err = m.saveModelBrand(m.value.BrandArray)
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
		err = m.saveModelAttrs(m.value.Attrs)
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
		err = m.saveModelSpecs(m.value.Specs)
	}
	// 保存商品模型
	if err == nil {
		i, err = m.rep.SaveProModel(m.value)
		if err == nil {
			m.value.ID = int32(i)
		}
	}
	return m.GetAggregateRootId(), err
}

// 保存规格
func (m *modelImpl) saveModelSpecs(specs []*promodel.Spec) (err error) {
	pk := m.GetAggregateRootId()
	// 获取存在的项
	old := m.rep.SelectSpec("pro_model = $1", pk)
	// 分析当前项目并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*promodel.Spec, len(specs))
	for _, v := range specs {
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
		m.specService.DeleteSpec(v)
	}
	// 保存项
	for _, v := range specs {
		if v.ModelId == 0 {
			v.ModelId = pk
		}
		if v.ModelId == pk {
			v.Id, err = m.specService.SaveSpec(v)
		}
	}
	return err
}

// 保存属性
func (m *modelImpl) saveModelAttrs(attrs []*promodel.Attr) (err error) {
	pk := m.GetAggregateRootId()
	// 获取存在的项
	old := m.rep.SelectAttr("pro_model = $1", pk)
	// 分析当前项目并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*promodel.Attr, len(attrs))
	for _, v := range attrs {
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
		m.attrService.DeleteAttr(v)
	}
	// 保存项
	for _, v := range attrs {
		if v.ModelId == 0 {
			v.ModelId = pk
		}
		if v.ModelId == pk {
			v.Id, err = m.attrService.SaveAttr(v)
		}
	}
	return err
}

// 保存品牌
func (m *modelImpl) saveModelBrand(brandIds []int32) (err error) {
	pk := m.GetAggregateRootId()
	//获取存在的品牌
	old := m.rep.SelectProModelBrand("pro_model = $1", pk)
	//删除不包括的品牌
	idArrStr := format.I32ArrStrJoin(brandIds)
	if len(old) > 0 {
		m.rep.BatchDeleteProModelBrand("pro_model = $1"+
			" AND brand_id NOT IN("+idArrStr+")", pk)
	}
	//写入品牌
	for _, v := range brandIds {
		isExist := false
		for _, vo := range old {
			if vo.BrandId == v {
				isExist = true
				break
			}
		}
		if !isExist {
			e := &promodel.ProModelBrand{
				ID:       0,
				BrandId:  v,
				ProModel: pk,
			}
			_, err = m.rep.SaveProModelBrand(e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

