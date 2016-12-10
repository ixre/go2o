package promodel

import "go2o/core/domain/interface/pro_model"

var _ promodel.IBrandService = new(brandServiceImpl)

type brandServiceImpl struct {
	rep promodel.IProModelRepo
}

func NewBrandService(rep promodel.IProModelRepo) *brandServiceImpl {
	return &brandServiceImpl{
		rep: rep,
	}
}

// 获取品牌
func (b *brandServiceImpl) Get(brandId int32) *promodel.ProBrand {
	return b.rep.GetProBrand(brandId)
}

// 保存品牌
func (b *brandServiceImpl) SaveBrand(v *promodel.ProBrand) (int32, error) {
	id, err := b.rep.SaveProBrand(v)
	return int32(id), err
}

// 删除品牌
func (b *brandServiceImpl) DeleteBrand(id int32) error {
	return b.rep.DeleteProBrand(id)
}

// 获取所有品牌
func (b *brandServiceImpl) AllBrands() []*promodel.ProBrand {
	return b.rep.SelectProBrand("")
}

// 获取关联的品牌编号
func (b *brandServiceImpl) Brands(proModel int32) []*promodel.ProBrand {
	return b.rep.SelectProBrand("pro_model=?", proModel)
}

// 关联品牌
func (b *brandServiceImpl) SetBrands(proModel int32, brandId []int32) error {
	return b.rep.SetModelBrands(proModel, brandId)
}
