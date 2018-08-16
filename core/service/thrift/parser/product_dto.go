package parser

import (
	"go2o/core/domain/interface/product"
	"go2o/core/service/auto_gen/rpc/ttype"
)

func Category(src *ttype.SCategory) *product.Category {
	s := &product.Category{
		ID:         src.ID,
		ParentId:   src.ParentId,
		ProModel:   src.ProModel,
		Priority:   src.Priority,
		Name:       src.Name,
		Level:      src.Level,
		Icon:       src.Icon,
		IconXY:     src.IconXY,
		VirtualCat: src.VirtualCat,
		CatUrl:     src.CatUrl,
		SortNum:    src.SortNum,
		Enabled:    src.Enabled,
		FloorShow:  src.FloorShow,
		CreateTime: src.CreateTime,
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
		ID:         src.ID,
		ParentId:   src.ParentId,
		ProModel:   src.ProModel,
		Priority:   src.Priority,
		Name:       src.Name,
		Level:      src.Level,
		Icon:       src.Icon,
		IconXY:     src.IconXY,
		VirtualCat: src.VirtualCat,
		CatUrl:     src.CatUrl,
		SortNum:    src.SortNum,
		FloorShow:  src.FloorShow,
		Enabled:    src.Enabled,
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
