package impl

import (
	"context"
	"errors"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"go2o/core/dao"
	"go2o/core/dao/impl"
	"go2o/core/dao/model"
	"go2o/core/service/proto"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : app_service.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-09 15:37
 * description :
 * history :
 */

var _ proto.AppServiceServer = new(appServiceImpl)
type appServiceImpl struct {
	dao dao.IAppProdDao
	s   storage.Interface
	serviceUtil
}


func NewAppService(s storage.Interface,o orm.Orm)*appServiceImpl{
	 d := impl.NewAppProdDao(o)
	c,err := d.Count("prod_name=$1 AND id <> $2","app",1)
	println(c,err)

	return &appServiceImpl{
		s:s,
		dao:impl.NewAppProdDao(o),
	}
}

func (a *appServiceImpl) SaveProd(_ context.Context, r *proto.AppProdRequest) (*proto.Result, error) {
	if c,_ := a.dao.Count("prod_name=$1 AND id <> $2",
		r.ProdName,r.Id);c > 0 {
		return a.error(errors.New("APP已经存在")), nil
	}
	dst := &model.AppProd{
		Id:             r.Id,
		ProdName:       r.ProdName,
		ProdDes:        r.ProdDes,
		PublishUrl:     r.PublishUrl,
		StableFileUrl:  r.StableFileUrl,
		AlphaFileUrl:   r.AlphaFileUrl,
		NightlyFileUrl: r.NightlyFileUrl,
	}
	_, err := a.dao.Save(dst)
	return a.error(err),nil
}

func (a *appServiceImpl) SaveVersion(_ context.Context, request *proto.AppVersionRequest) (*proto.Result, error) {
	panic("implement me")
}

func (a *appServiceImpl) GetProd(_ context.Context, id *proto.AppId) (*proto.SAppProd, error) {
	panic("implement me")
}

func (a *appServiceImpl) GetVersion(_ context.Context, id *proto.AppVersionId) (*proto.SAppVersion, error) {
	panic("implement me")
}

func (a *appServiceImpl) DeleteProd(_ context.Context, id *proto.AppId) (*proto.Result, error) {
	panic("implement me")
}

func (a *appServiceImpl) DeleteVersion(_ context.Context, id *proto.AppVersionId) (*proto.Result, error) {
	panic("implement me")
}

func (a *appServiceImpl) CheckVersion(_ context.Context, request *proto.CheckVersionRequest) (*proto.CheckVersionResponse, error) {
	panic("implement me")
}