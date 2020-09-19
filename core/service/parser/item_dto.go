package parser

import (
	"github.com/ixre/gof/types"
	"go2o/core/domain/interface/item"
	"go2o/core/service/proto"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : item_dto.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-19 13:50
 * description :
 * history :
 */

func ParseSaleLabelDto(src *item.Label) *proto.SItemLabel {
	return &proto.SItemLabel{
		Id:                   src.Id,
		Name:                 src.TagName,
		TagCode:              src.TagCode,
		LabelImage:           src.LabelImage,
		Enabled:              src.Enabled == 1,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
}

func FromSaleLabelDto(src *proto.SItemLabel) *item.Label {
	return &item.Label{
		Id:         src.Id,
		MerchantId: -1,
		TagCode:    src.TagCode,
		TagName:    src.Name,
		LabelImage: src.LabelImage,
		Enabled:    types.IntCond(src.Enabled, 1, 0),
	}
}
