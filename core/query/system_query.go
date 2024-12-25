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
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type SystemQuery struct {
	appRepo             fw.Repository[sys.SysAppVersion]
	appDistributionRepo fw.Repository[sys.SysAppDistribution]
	sysLogRepo          fw.Repository[sys.SysLog]
}

func NewSystemQuery(orm fw.ORM) *SystemQuery {
	return &SystemQuery{
		appRepo:             fw.NewRepository[sys.SysAppVersion](orm),
		appDistributionRepo: fw.NewRepository[sys.SysAppDistribution](orm),
		sysLogRepo:          fw.NewRepository[sys.SysLog](orm),
	}
}

// QueryIndustries 查询行业数据
func (s *SystemQuery) QueryIndustries(parentId string) []*sys.GeneralOption {
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

// QueryAppVersion 查询应用版本
func (s *SystemQuery) QueryAppVersion(params *fw.PagingParams) (*fw.PagingResult, error) {
	return s.appRepo.QueryPaging(params)
}

// QueryAppDistribution 查询应用分发
func (s *SystemQuery) QueryAppDistribution(params *fw.PagingParams) (*fw.PagingResult, error) {
	return s.appDistributionRepo.QueryPaging(params)
}

// QuerySysLog 查询系统日志
func (s *SystemQuery) QueryPagingSysLog(params *fw.PagingParams) (*fw.PagingResult, error) {
	return s.sysLogRepo.QueryPaging(params)
}

// GetSysLog 获取系统日志
func (s *SystemQuery) GetSysLog(id int) *sys.SysLog {
	return s.sysLogRepo.Get(id)
}
