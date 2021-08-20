package impl

import (
	"context"
	"errors"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types"
	"go2o/core/dao"
	"go2o/core/dao/impl"
	"go2o/core/dao/model"
	"go2o/core/service/proto"
	"strconv"
	"strings"
	"time"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : 2.app_service.go
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

func NewAppService(s storage.Interface, o orm.Orm) *appServiceImpl {
	return &appServiceImpl{
		s:   s,
		dao: impl.NewAppProdDao(o),
	}
}

func IntVersion(s string) int {
	arr := strings.Split(s, ".")
	for i, v := range arr {
		if l := len(v); l < 3 {
			arr[i] = strings.Repeat("0", 3-l) + v
		}
	}
	intVer, err := strconv.Atoi(strings.Join(arr, ""))
	if err != nil {
		panic(err)
	}
	return intVer
}

func (a *appServiceImpl) SaveProd(_ context.Context, r *proto.AppProdRequest) (*proto.Result, error) {
	if c, _ := a.dao.Count("prod_name=$1 AND id <> $2",
		r.ProdName, r.Id); c > 0 {
		return a.error(errors.New("APP已经存在")), nil
	}
	var dst *model.AppProd
	if r.Id > 0 {
		dst = a.dao.Get(r.Id)
	} else {
		dst = &model.AppProd{}
	}
	dst.ProdName = r.ProdName
	dst.ProdDes = r.ProdDes
	dst.PublishUrl = r.PublishUrl
	dst.StableFileUrl = r.StableFileUrl
	dst.AlphaFileUrl = r.AlphaFileUrl
	dst.NightlyFileUrl = r.NightlyFileUrl
	dst.UpdateType = r.UpdateType
	dst.UpdateTime = time.Now().Unix()
	_, err := a.dao.Save(dst)
	return a.error(err), nil
}

func (a *appServiceImpl) SaveVersion(_ context.Context, r *proto.AppVersionRequest) (*proto.Result, error) {
	if c, _ := a.dao.Count("version=$1 AND id <> $2",
		r.Version, r.Id); c > 0 {
		return a.error(errors.New("版本已经存在")), nil
	}
	dst := &model.AppVersion{
		Id:            r.Id,
		ProductId:     r.ProductId,
		Channel:       int16(r.Channel),
		Version:       r.Version,
		VersionCode:   IntVersion(r.Version),
		ForceUpdate:   int16(types.IntCond(r.ForceUpdate, 1, 0)),
		UpdateContent: r.UpdateContent,
	}
	if r.Id <= 0 {
		dst.CreateTime = time.Now().Unix()
	}
	id, err := a.dao.SaveVersion(dst)
	// 如果发布了更新版本,则更新最新的版本
	if err == nil && r.Channel == 0 {
		r.Id = int64(id)
		err = a.updateLatest(r)
	}
	return a.error(err), nil
}

// 更新最新版本
func (a *appServiceImpl) updateLatest(r *proto.AppVersionRequest) error {
	latest := a.getLatest(r.ProductId, int16(r.Channel))
	if latest == nil || IntVersion(r.Version) > IntVersion(latest.Version) {
		prod := a.dao.Get(r.ProductId)
		prod.LatestVid = r.Id
		prod.UpdateTime = time.Now().Unix()
		_, err := a.dao.Save(prod)
		return err
	}
	return nil
}

// 获取最新版本
func (a *appServiceImpl) getLatest(prodId int64, channel int16) *model.AppVersion {
	return a.dao.GetVersionBy(
		"product_id=$1 AND channel=$2 ORDER BY version_code DESC LIMIT 1",
		prodId, channel)
}

func (a *appServiceImpl) GetProd(_ context.Context, id *proto.AppId) (*proto.SAppProd, error) {
	v := a.dao.Get(id.Value)
	if v == nil {
		return nil, nil
	}
	return &proto.SAppProd{
		Id:             v.Id,
		ProdName:       v.ProdName,
		ProdDes:        v.ProdDes,
		LatestVid:      v.LatestVid,
		Md5Hash:        v.Md5Hash,
		PublishUrl:     v.PublishUrl,
		StableFileUrl:  v.StableFileUrl,
		AlphaFileUrl:   v.AlphaFileUrl,
		NightlyFileUrl: v.NightlyFileUrl,
		UpdateType:     v.UpdateType,
		UpdateTime:     v.UpdateTime,
	}, nil
}

func (a *appServiceImpl) GetVersion(_ context.Context, id *proto.AppVersionId) (*proto.SAppVersion, error) {
	v := a.dao.GetVersion(id.Value)
	if v == nil {
		return nil, nil
	}
	return &proto.SAppVersion{
		Id:            v.Id,
		ProductId:     v.ProductId,
		Channel:       int32(v.Channel),
		Version:       v.Version,
		VersionCode:   int32(v.VersionCode),
		ForceUpdate:   v.ForceUpdate == 1,
		UpdateContent: v.UpdateContent,
		CreateTime:    v.CreateTime,
	}, nil
}

func (a *appServiceImpl) DeleteProd(_ context.Context, id *proto.AppId) (*proto.Result, error) {
	err := a.dao.Delete(id.Value)
	return a.error(err), nil
}

func (a *appServiceImpl) DeleteVersion(_ context.Context, id *proto.AppVersionId) (*proto.Result, error) {
	err := a.dao.Delete(id.Value)
	return a.error(err), nil
}

func (a *appServiceImpl) CheckVersion(_ context.Context, r *proto.CheckVersionRequest) (*proto.CheckVersionResponse, error) {
	var v *model.AppVersion
	prod := a.dao.Get(r.AppId)
	switch r.Channel {
	case "stable":
		v = a.dao.GetVersion(prod.LatestVid)
	case "nightly":
		v = a.getLatest(r.AppId, 2)
	case "alpha":
		v = a.getLatest(r.AppId, 1)
	}
	if v == nil {
		return &proto.CheckVersionResponse{
			VersionInfo: "没有最近发布的版本",
			IsLatest:    true,
		}, nil
	}
	if v.VersionCode <= IntVersion(r.Version) {
		return &proto.CheckVersionResponse{
			LatestVersion: r.Version,
			VersionInfo:   "当前为最新版本",
			IsLatest:      true,
		}, nil
	}

	return &proto.CheckVersionResponse{
		LatestVersion: r.Version,
		AppPkgUrl:     a.getVersionPkgURL(prod, r.Channel, r.Version),
		VersionInfo:   v.UpdateContent,
		IsLatest:      false,
		ForceUpdate:   v.ForceUpdate == 1,
		UpdateType:    prod.UpdateType,
		ReleaseTime:   v.CreateTime,
	}, nil
}

// 获取版本更新包地址
func (a *appServiceImpl) getVersionPkgURL(prod *model.AppProd,
	channel string, version string) string {
	s := ""
	switch channel {
	case "stable":
		s = prod.StableFileUrl
	case "nightly":
		s = prod.NightlyFileUrl
	case "alpha":
		s = prod.AlphaFileUrl
	}
	return strings.Replace(s, "{version}", version, -1)
}
