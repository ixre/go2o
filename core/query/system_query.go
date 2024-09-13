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
)

type SystemQuery struct {
}

func NewSystemQuery() *SystemQuery {
	return &SystemQuery{}
}

// QueryIndustries 查询行业数据
func (q *SystemQuery) QueryIndustries(parentId string) []*sys.GeneralOption {
	if parentId == "0" || parentId == "" {
		dst := make([]*sys.GeneralOption, len(sys.INDUSTRY_DATA))
		// 获取根节点数据
		for i, v := range sys.INDUSTRY_DATA {
			dst[i] = &sys.GeneralOption{
				Label:  v.Label,
				Value:  v.Value,
				IsLeaf: false,
			}
		}
		return dst
	}
	for _, v := range sys.INDUSTRY_DATA {
		if v.Value == parentId {
			// 获取下级数据
			ret := make([]*sys.GeneralOption, len(v.Children))
			for i, v := range v.Children {
				ret[i] = &sys.GeneralOption{
					Label:  v.Label,
					Value:  v.Value,
					IsLeaf: true,
				}
			}
			return ret
		}
	}
	return []*sys.GeneralOption{}
}
