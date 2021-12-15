/**
 * Copyright 2015 @ 56x.net.
 * name : api_manager.go
 * author : jarryliu
 * date : 2016-05-27 13:28
 * description :
 * history :
 */
package merchant

import (
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/infrastructure/domain"
)

var _ merchant.IApiManager = new(apiManagerImpl)

type apiManagerImpl struct {
	*merchantImpl
	apiInfo *merchant.ApiInfo
}

func newApiManagerImpl(m *merchantImpl) merchant.IApiManager {
	return &apiManagerImpl{
		merchantImpl: m,
	}
}

// 获取API信息
func (a *apiManagerImpl) getApiInfo() *merchant.ApiInfo {
	if a.apiInfo == nil {
		a.apiInfo = a._repo.GetApiInfo(int(a.GetAggregateRootId()))
		//没有API则生成
		if a.apiInfo == nil {
			mchId := int(a.GetAggregateRootId())
			a.apiInfo = &merchant.ApiInfo{
				MerchantId: int(a.GetAggregateRootId()),
				ApiId:      domain.NewApiId(mchId),
				ApiSecret:  domain.NewSecret(mchId),
				WhiteList:  "*",
				Enabled:    0,
			}
			a.SaveApiInfo(a.apiInfo)
		}
	}
	return a.apiInfo
}

// 获取API信息
func (a *apiManagerImpl) GetApiInfo() merchant.ApiInfo {
	return *a.getApiInfo()
}

// 保存API信息
func (a *apiManagerImpl) SaveApiInfo(v *merchant.ApiInfo) error {
	a.apiInfo = v
	a.apiInfo.MerchantId = int(a.GetAggregateRootId())
	return a._repo.SaveApiInfo(a.apiInfo)
}

// 启用API权限
func (a *apiManagerImpl) EnableApiPerm() error {
	v := a.getApiInfo()
	v.Enabled = 1
	return a.SaveApiInfo(v)
}

// 禁用API权限
func (a *apiManagerImpl) DisableApiPerm() error {
	v := a.getApiInfo()
	v.Enabled = 0
	return a.SaveApiInfo(v)
}
