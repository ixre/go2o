package impl

import (
	"context"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/storage"
)

/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2020-09-05 20:14
 * description :
 * history :
 */

var _ proto.CheckServiceServer = new(checkService)

type checkService struct {
	repo         member.IMemberRepo
	registryRepo registry.IRegistryRepo
	store        storage.Interface
	serviceUtil
	proto.UnimplementedCheckServiceServer
}

// CompareCode implements proto.CheckServiceServer.
func (c *checkService) CompareCode(context.Context, *proto.CompareCheckCodeRequest) (*proto.Result, error) {
	panic("unimplemented")
}

// NewCheckService 校验服务
func NewCheckService(repo member.IMemberRepo,
	registryRepo registry.IRegistryRepo,
	store storage.Interface,
) proto.CheckServiceServer {
	s := &checkService{
		repo:         repo,
		registryRepo: registryRepo,
		store:        store,
	}
	return s
}

// SendCode 发送验证码
func (c *checkService) SendCode(_ context.Context, r *proto.SendCheckCodeRequest) (*proto.SendCheckCodeResponse, error) {

	panic("unimplemented")
}
