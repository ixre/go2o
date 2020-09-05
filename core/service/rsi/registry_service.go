/**
 * Copyright 2015 @ to2.net.
 * name : platform_service
 * author : jarryliu
 * date : 2016-05-27 15:30
 * description :
 * history :
 */
package rsi

import (
	"context"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/service/auto_gen/rpc/foundation_service"
	"go2o/core/service/auto_gen/rpc/registry_service"
	"go2o/core/service/auto_gen/rpc/ttype"
	"strings"
)

var _ foundation_service.FoundationService = new(foundationService)

// 基础服务
type registryService struct {
	_rep         valueobject.IValueRepo
	registryRepo registry.IRegistryRepo
	serviceUtil
}

func NewRegistryService(rep valueobject.IValueRepo, registryRepo registry.IRegistryRepo) *registryService {
	return &registryService{
		_rep:         rep,
		registryRepo: registryRepo,
	}
}

// 根据键获取值
func (s *registryService) GetValue(ctx context.Context, key string) (r string, err error) {
	ir := s.registryRepo.Get(key)
	if ir != nil {
		return ir.StringValue(), nil
	}
	return "", nil
}

// 获取键值存储数据
func (s *registryService) GetRegistries(ctx context.Context, keys []string) (map[string]string, error) {
	mp := make(map[string]string)
	for _, k := range keys {
		if ir := s.registryRepo.Get(k); ir != nil {
			mp[k] = ir.StringValue()
		} else {
			mp[k] = ""
		}
	}
	return mp, nil
}

// 按键前缀获取键数据
func (s *registryService) FindRegistries(ctx context.Context, prefix string) (r map[string]string, err error) {
	mp := make(map[string]string)
	for _, k := range s.registryRepo.SearchRegistry(prefix) {
		if strings.HasPrefix(k.Key, prefix) {
			mp[k.Key] = k.Value
		}
	}
	return mp, nil
}

// 搜索注册表
func (s *registryService) SearchRegistry(ctx context.Context, key string) (r []*registry_service.SRegistry, err error) {
	arr := s.registryRepo.SearchRegistry(key)
	list := make([]*registry_service.SRegistry, len(arr))
	for i, a := range arr {
		list[i] = &registry_service.SRegistry{
			Key:         a.Key,
			Value:       a.Value,
			Default:     a.DefaultValue,
			Options:     a.Options,
			Flag:        int32(a.Flag),
			Description: a.Description,
		}
	}
	return list, nil
}

// 获取数据存储
func (s *registryService) GetRegistry(ctx context.Context, key string) (string, error) {
	ir := s.registryRepo.Get(key)
	if ir != nil {
		return ir.StringValue(), nil
	}
	return "", nil
}

// 创建用户自定义注册项
func (s *registryService) CreateUserRegistry(ctx context.Context, key string, defaultValue string, description string) (r *ttype.Result_, err error) {
	if s.registryRepo.Get(key) != nil {
		return s.resultWithCode(-1, "registry is exist"), nil
	}
	rv := &registry.Registry{
		Key:          key,
		Value:        defaultValue,
		DefaultValue: defaultValue,
		Options:      "",
		Flag:         registry.FlagUserDefine,
		Description:  description,
	}
	ir := s.registryRepo.Create(rv)
	err = ir.Save()
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 更新注册表数据
func (s *registryService) UpdateRegistry(ctx context.Context, registries map[string]string) (r *ttype.Result_, err error) {
	for k, v := range registries {
		if ir := s.registryRepo.Get(k); ir != nil {
			if err = ir.Update(v); err != nil {
				return s.error(err), nil
			}
		}
	}
	return s.success(nil), nil
}

// 获取键值存储数据
func (s *registryService) GetRegistryV1(ctx context.Context, keys []string) ([]string, error) {
	return s._rep.GetsRegistry(keys), nil
}
