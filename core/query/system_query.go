/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: system_query.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-09-13 19:46:51
 * description:
 * history:
 */

package query

import (
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
)

type SystemQuery struct {
}

func NewSystemQuery() *SystemQuery {
	return &SystemQuery{}
}

// QueryIndustries 查询行业数据
func (q *SystemQuery) QueryIndustries(parentId string) []*sys.GeneralOption {
	data := types.DeepClone(&sys.INDUSTRY_DATA)
	if parentId == "0" || parentId == "" {
		// 获取根节点数据
		for _, v := range *data {
			(*v).IsLeaf = false
			(*v).Children = nil
		}
		return *data
	}
	for _, v := range *data {
		if v.Value == parentId {
			// 获取下级数据
			ret := v.Children
			for _, v := range ret {
				(*v).IsLeaf = true
				(*v).Children = nil
			}
			return ret
		}
	}
	return []*sys.GeneralOption{}
}
