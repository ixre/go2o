/**
 * Copyright 2015 @ z3q.net.
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
	_apiInfo *merchant.ApiInfo
}

func newApiManagerImpl(m *merchantImpl) merchant.IApiManager {
	return &apiManagerImpl{
		merchantImpl: m,
	}
}

// 获取API信息
func (this *apiManagerImpl) getApiInfo() *merchant.ApiInfo {
	if this._apiInfo == nil {
		this._apiInfo = this._rep.GetApiInfo(this.GetAggregateRootId())
		//没有API则生成
		if this._apiInfo == nil {
			mchId := this.GetAggregateRootId()
			this._apiInfo = &merchant.ApiInfo{
				MerchantId: mchId,
				ApiId:      domain.NewApiId(mchId),
				ApiSecret:  domain.NewSecret(mchId),
				WhiteList:  "*",
				Enabled:    0,
			}
			this.SaveApiInfo(this._apiInfo)
		}
	}
	return this._apiInfo
}

// 获取API信息
func (this *apiManagerImpl) GetApiInfo() merchant.ApiInfo {
	return *this.getApiInfo()
}

// 保存API信息
func (this *apiManagerImpl) SaveApiInfo(v *merchant.ApiInfo) error {
	this._apiInfo = v
	this._apiInfo.MerchantId = this.GetAggregateRootId()
	return this._rep.SaveApiInfo(this._apiInfo)
}

// 启用API权限
func (this *apiManagerImpl) EnableApiPerm() error {
	v := this.getApiInfo()
	v.Enabled = 1
	return this.SaveApiInfo(v)
}

// 禁用API权限
func (this *apiManagerImpl) DisableApiPerm() error {
	v := this.getApiInfo()
	v.Enabled = 0
	return this.SaveApiInfo(v)
}
