package parser

import (
	"go2o/core/domain/interface/product"
	"go2o/core/service/auto_gen/rpc/ttype"
)

func Category(src *ttype.SCategory) *product.Category {
	s := &product.Category{
		Id:         int(src.ID),
		ParentId:   int(src.ParentId),
		ProdModel:  int(src.ProdModel),
		Priority:   int(src.Priority),
		Name:       src.Name,
		Level:      int(src.Level),
		Icon:       src.Icon,
		IconXy:     src.IconXy,
		VirtualCat: int(src.VirtualCat),
		CatUrl:     src.CatUrl,
		SortNum:    int(src.SortNum),
		Enabled:    int(src.Enabled),
		FloorShow:  int(src.FloorShow),
		CreateTime: int64(src.CreateTime),
	}
	if src.Children != nil {
		s.Children = make([]*product.Category, len(src.Children))
		for i, v := range src.Children {
			s.Children[i] = Category(v)
		}
	}
	return s
}

func CategoryDto(src *product.Category) *ttype.SCategory {
	s := &ttype.SCategory{
		ID:         int32(src.Id),
		ParentId:   int32(src.ParentId),
		ProdModel:  int32(src.ProdModel),
		Priority:   int32(src.Priority),
		Name:       src.Name,
		Level:      int32(src.Level),
		Icon:       src.Icon,
		IconXy:     src.IconXy,
		VirtualCat: int32(src.VirtualCat),
		CatUrl:     src.CatUrl,
		SortNum:    int32(src.SortNum),
		FloorShow:  int32(src.FloorShow),
		Enabled:    int32(src.Enabled),
		CreateTime: src.CreateTime,
	}
	if src.Children != nil {
		s.Children = make([]*ttype.SCategory, len(src.Children))
		for i, v := range src.Children {
			s.Children[i] = CategoryDto(v)
		}
	}
	return s
}
