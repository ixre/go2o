/**
 * Copyright 2015 @ 56x.net.
 * name : platform_service
 * author : jarryliu
 * date : 2016-05-27 15:30
 * description :
 * history :
 */
package impl

import (
	"context"
	"errors"
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
	v, err := s.registryRepo.GetValue(key.Value)
	rsp := &proto.RegistryValueResponse{Value: v}
	if err != nil {
		rsp.ErrMsg = err.Error()
	}
	return rsp, nil
}

func (s *registryService) UpdateValue(_ context.Context, pair *proto.RegistryPair) (r *proto.Result, err error) {
	e := s.registryRepo.Get(pair.Key)
	if e == nil {
		err = errors.New("not exists key")
	} else {
		err = s.registryRepo.UpdateValue(pair.Key, pair.Value)
	}
	if err != nil {
		return s.error(err), nil
	}
	return s.success(nil), nil
}

// 获取键值存储数据
func (s *registryService) GetValues(_ context.Context, array *proto.StringArray) (*proto.StringMap, error) {
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

// Search 搜索键值
func (s *registryService) Search(_ context.Context, r *proto.RegistrySearchRequest) (*proto.StringMap, error) {
	mp := make(map[string]string)
	for _, v := range s.registryRepo.SearchRegistry(r.Key) {
		mp[v.Key] = v.Value
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
			DefaultValue:     a.DefaultValue,
			Options:     a.Options,
			Flag:        int32(a.Flag),
			Description: a.Description,
		}
	}
	return &proto.RegistriesResponse{Value: list}, nil
}

// 创建用户自定义注册项
func (s *registryService) CreateRegistry(_ context.Context, r *proto.RegistryCreateRequest) (*proto.Result, error) {
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

// UpdateValues 更新注册表数据
func (s *registryService) UpdateValues(_ context.Context, registries *proto.StringMap) (*proto.Result, error) {
	for k, v := range registries.Value {
		if ir := s.registryRepo.Get(k); ir != nil {
			err := ir.Update(v)
			if err == nil {
				err = ir.Save()
			}
			if err != nil {
				return s.error(err), nil
			}
		}
	}
	return s.success(nil), nil
}
