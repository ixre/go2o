package rsi

import (
	"go2o/core/domain/interface/pro_model"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/service/thrift/parser"
)

// 产品服务
type productService struct {
	pmRep promodel.IProModelRepo
}

func NewProService(pmRep promodel.IProModelRepo) *productService {
	return &productService{
		pmRep: pmRep,
	}
}

// 获取产品模型
func (p *productService) GetModel(id int32) *promodel.ProModel {
	return p.pmRep.GetProModel(id)
}

// 获取产品模型
func (p *productService) GetModels() []*promodel.ProModel {
	return p.pmRep.SelectProModel("enabled=1")
}

// 获取模型属性
func (p *productService) GetModelAttrs(proModel int32) []*promodel.Attr {
	m := p.pmRep.CreateModel(&promodel.ProModel{Id: proModel})
	return m.Attrs()
}

// 获取模型规格
func (p *productService) GetModelSpecs(proModel int32) []*promodel.Spec {
	m := p.pmRep.CreateModel(&promodel.ProModel{Id: proModel})
	return m.Specs()
}

// 保存产品模型
func (p *productService) SaveModel(v *promodel.ProModel) (*define.Result_, error) {
	var pm promodel.IModel
	var err error
	if v.Id > 0 {
		ev := p.GetModel(v.Id)
		if ev == nil {
			return &define.Result_{Message: "模型不存在"}, nil
		}
		ev.Name = v.Name
		ev.Enabled = v.Enabled
		pm = p.pmRep.CreateModel(ev)
	} else {
		pm = p.pmRep.CreateModel(v)
	}
	// 保存属性
	if err == nil && v.Attrs != nil {
		err = pm.SetAttrs(v.Attrs)
	}
	// 保存规格
	if err == nil && v.Specs != nil {
		err = pm.SetSpecs(v.Specs)
	}
	// 保存品牌
	if err == nil && v.BrandArray != nil {
		err = pm.SetBrands(v.BrandArray)
	}
	// 保存模型
	if err == nil {
		v.Id, err = pm.Save()
	}
	r := parser.Result(err)
	r.ID = v.Id
	return r, nil
}

// 删除产品模型
func (p *productService) DeleteProModel_(id int32) (*define.Result_, error) {
	return &define.Result_{Result_: true}, nil
}

// Get 产品品牌
func (p *productService) GetProBrand_(id int32) *promodel.ProBrand {
	return p.pmRep.BrandService().Get(id)
}

// Save 产品品牌
func (p *productService) SaveProBrand_(v *promodel.ProBrand) (*define.Result_, error) {
	id, err := p.pmRep.BrandService().SaveBrand(v)
	r := parser.Result(err)
	r.ID = id
	return r, nil
}

// Delete 产品品牌
func (p *productService) DeleteProBrand_(id int32) (*define.Result_, error) {
	err := p.pmRep.BrandService().DeleteBrand(id)
	return parser.Result(err), nil
}

// 获取所有产品品牌
func (p *productService) GetBrands() []*promodel.ProBrand {
	return p.pmRep.BrandService().AllBrands()
}

// 获取模型关联的产品品牌
func (p *productService) GetModelBrands(id int32) []*promodel.ProBrand {
	pm := p.pmRep.CreateModel(&promodel.ProModel{Id: id})
	return pm.Brands()
}
