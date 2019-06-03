/**
 * Copyright 2015 @ to2.net.
 * name : api_manager.go
 * author : jarryliu
 * date : 2016-05-27 13:28
 * description :
 * history :
 */
package merchant

import (
	"go2o/core/domain/interface/merchant"
	"go2o/core/infrastructure/domain"
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
		a.apiInfo = a._rep.GetApiInfo(a.GetAggregateRootId())
		//没有API则生成
		if a.apiInfo == nil {
			mchId := int(a.GetAggregateRootId())
			a.apiInfo = &merchant.ApiInfo{
				MerchantId: a.GetAggregateRootId(),
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
	a.apiInfo.MerchantId = a.GetAggregateRootId()
	return a._rep.SaveApiInfo(a.apiInfo)
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
