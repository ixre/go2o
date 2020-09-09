/**
 * Copyright 2015 @ to2.net.
 * name : platform_service
 * author : jarryliu
 * date : 2016-05-27 15:30
 * description :
 * history :
 */
package impl

import (
	"context"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/service/proto"
	"strings"
)

var _ proto.RegistryServiceServer = new(registryService)

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


// 获取数据存储
func (s *registryService) GetValue(_ context.Context, key *proto.String) (*proto.RegistryValueResponse, error) {
	v,err := s.registryRepo.GetValue(key.Value)
	rsp := &proto.RegistryValueResponse{Value: v}
	if err != nil {
		rsp.ErrorMsg = err.Error()
	}
	return rsp, nil
}

// 获取键值存储数据
func (s *registryService) GetRegistries(_ context.Context, array *proto.StringArray) (*proto.StringMap, error) {
	mp := make(map[string]string)
	for _, k := range array.Value {
		if ir := s.registryRepo.Get(k); ir != nil {
			mp[k] = ir.StringValue()
		} else {
			mp[k] = ""
		}
	}
	return &proto.StringMap{Value: mp}, nil
}

// 按键前缀获取键数据
func (s *registryService) FindRegistries(_ context.Context, prefix *proto.String) (*proto.StringMap, error) {
	mp := make(map[string]string)
	for _, k := range s.registryRepo.SearchRegistry(prefix.Value) {
		if strings.HasPrefix(k.Key, prefix.Value) {
			mp[k.Key] = k.Value
		}
	}
	return &proto.StringMap{Value: mp}, nil
}

// 搜索注册表
func (s *registryService) SearchRegistry(_ context.Context, key *proto.String) (*proto.RegistriesResponse, error) {
	arr := s.registryRepo.SearchRegistry(key.Value)
	list := make([]*proto.SRegistry, len(arr))
	for i, a := range arr {
		list[i] = &proto.SRegistry{
			Key:         a.Key,
			Value:       a.Value,
			Default:     a.DefaultValue,
			Options:     a.Options,
			Flag:        int32(a.Flag),
			Description: a.Description,
		}
	}
	return &proto.RegistriesResponse{Value: list}, nil
}

// 创建用户自定义注册项
func (s *registryService) CreateUserRegistry(_ context.Context, r *proto.UserRegistryCreateRequest) (*proto.Result, error) {
	if s.registryRepo.Get(r.Key) != nil {
		return s.resultWithCode(-1, "registry is exist"), nil
	}
	rv := &registry.Registry{
		Key:          r.Key,
		Value:        r.DefaultValue,
		DefaultValue: r.DefaultValue,
		Options:      "",
		Flag:         registry.FlagUserDefine,
		Description:  r.Description,
	}
	ir := s.registryRepo.Create(rv)
	err := ir.Save()
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 更新注册表数据
func (s *registryService) UpdateRegistryValues(_ context.Context, registries *proto.StringMap) (*proto.Result, error) {
	for k, v := range registries.Value {
		if ir := s.registryRepo.Get(k); ir != nil {
			if err := ir.Update(v); err != nil {
				return s.error(err), nil
			}
		}
	}
	return s.success(nil), nil
}

// 获取键值存储数据
func (s *registryService) GetRegistryV1(_ context.Context, array *proto.StringArray) (*proto.StringArray, error) {
	a := s._rep.GetsRegistry(array.Value)
	return &proto.StringArray{Value: a}, nil
}
