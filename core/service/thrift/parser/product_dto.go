package parser

import (
	"go2o/core/domain/interface/product"
	"go2o/core/service/proto"
	"go2o/core/service/thrift/auto_gen/rpc/ttype"
)

func Category(src *proto.SCategory) *product.Category {
	s := &product.Category{
		Id:         int(src.Id),
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

func CategoryDto(src *product.Category) *proto.SCategory {
	s := &proto.SCategory{
		Id:         int32(src.Id),
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
		s.Children = make([]*proto.SCategory, len(src.Children))
		for i, v := range src.Children {
			s.Children[i] = CategoryDto(v)
		}
	}
	return s
}

func CategoryThrift(src *ttype.SCategory) *product.Category {
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
			s.Children[i] = CategoryThrift(v)
		}
	}
	return s
}

func CategoryDtoThrift(src *product.Category) *ttype.SCategory {
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
			s.Children[i] = CategoryDtoThrift(v)
		}
	}
	return s
}
