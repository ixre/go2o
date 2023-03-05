package promodel

import (
	"fmt"

	"github.com/ixre/go2o/core/domain/interface/pro_model"
	"github.com/ixre/go2o/core/domain/interface/product"
)

var _ promodel.IBrandService = new(brandServiceImpl)

type brandServiceImpl struct {
	repo promodel.IProductModelRepo
}

func NewBrandService(rep promodel.IProductModelRepo) *brandServiceImpl {
	return &brandServiceImpl{
		repo: rep,
	}
}

// 获取品牌
func (b *brandServiceImpl) Get(brandId int) *promodel.ProductBrand {
	return b.repo.GetProBrand(brandId)
}

// 保存品牌
func (b *brandServiceImpl) SaveBrand(v *promodel.ProductBrand) (int, error) {
	return b.repo.SaveProBrand(v)
}

// 删除品牌
func (b *brandServiceImpl) DeleteBrand(id int) error {
	arr := b.repo.SelectProModelBrand("brand_id=$1", id)
	for _, v := range arr {
		if m := b.repo.GetModel(v.ModelId); m != nil {
			return fmt.Errorf(product.ErrBrandIsUsed.Error(), m.Value().Name)
		}
	}
	return b.repo.DeleteProBrand(id)
}

// 获取所有品牌
func (b *brandServiceImpl) AllBrands() []*promodel.ProductBrand {
	return b.repo.SelectProBrand("")
}

// 获取关联的品牌编号
func (b *brandServiceImpl) Brands(proModel int) []*promodel.ProductBrand {
	return b.repo.GetModelBrands(proModel)
}
