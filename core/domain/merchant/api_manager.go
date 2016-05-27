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
	"errors"
	"fmt"
	"go2o/core/domain/interface/merchant"
)

var _ merchant.IApiManager = new(apiManagerImpl)

type apiManagerImpl struct {
	*MerchantImpl
	_apiInfo *merchant.ApiInfo
}

func newApiManagerImpl(m *MerchantImpl) merchant.IApiManager {
	return &apiManagerImpl{
		MerchantImpl: m,
	}
}

// 获取API信息
func (this *apiManagerImpl) getApiInfo() *merchant.ApiInfo {
	if this._apiInfo == nil {
		this._apiInfo = this._rep.GetApiInfo(this.GetAggregateRootId())
		if this._apiInfo == nil {
			panic(errors.New(fmt.Sprintf("商户:%d-%s 未生成API数据",
				this.GetAggregateRootId(), this.GetValue().Name)))
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
