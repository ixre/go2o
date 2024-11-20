package impl

import (
	"context"
	"errors"

	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
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
	_repo    sys.IApplicationRepository
	_sysRepo sys.ISystemRepo
	s        storage.Interface
	serviceUtil
	proto.UnimplementedAppServiceServer
}

func NewAppService(s storage.Interface, o orm.Orm,
	repo sys.IApplicationRepository,
	sysRepo sys.ISystemRepo,
) proto.AppServiceServer {
	return &appServiceImpl{
		s:        s,
		_repo:    repo,
		_sysRepo: sysRepo,
	}
}

func (a *appServiceImpl) SaveAppDistribution(_ context.Context,
	r *proto.SAppDistribution) (*proto.TxResult, error) {
	im := a._sysRepo.GetSystemAggregateRoot().Application()
	err := im.SaveAppDistribution(&sys.SysAppDistribution{
		Id:             int(r.Id),
		AppName:        r.AppName,
		AppIcon:        r.AppIcon,
		AppDesc:        r.AppDesc,
		UpdateMode:     int(r.UpdateMode),
		DistributeUrl:  r.DistributeUrl,
		DistributeName: r.DistributeName,
		UrlScheme:      r.UrlScheme,
		StableVersion:  r.StableVersion,
		StableDownUrl:  r.StableDownUrl,
		BetaVersion:    r.BetaVersion,
		BetaDownUrl:    r.BetaDownUrl,
	})
	return a.errorV2(err), nil
}

func (a *appServiceImpl) SaveAppVersion(_ context.Context, r *proto.SAppVersion) (*proto.TxResult, error) {
	im := a._sysRepo.GetSystemAggregateRoot().Application()
	err := im.SaveAppVersion(&sys.SysAppVersion{
		Id:              int(r.Id),
		DistributionId:  int(r.DistributionId),
		Version:         r.Version,
		VersionCode:     util.IntVersion(r.Version),
		TerminalOs:      r.TerminalOs,
		TerminalChannel: r.TerminalChannel,
		StartTime:       int(r.StartTime),
		UpdateMode:      int(r.UpdateMode),
		UpdateContent:   r.UpdateContent,
		PackageUrl:      r.PackageUrl,
		IsForce:         int(r.IsForce),
		IsNotified:      int(r.IsNotified),
		CreateTime:      int(r.CreateTime),
		UpdateTime:      int(r.UpdateTime),
	})
	return a.errorV2(err), nil
}

// GetAppDistribution 获取应用分发
func (a *appServiceImpl) GetAppDistribution(_ context.Context, req *proto.SysAppDistributionId) (*proto.SAppDistribution, error) {
	dist := a._sysRepo.GetSystemAggregateRoot().Application().GetAppDistribution(int(req.Value))
	if dist == nil {
		return nil, errors.New("no such app product")
	}
	return &proto.SAppDistribution{
		Id:             int64(dist.Id),
		AppName:        dist.AppName,
		AppIcon:        dist.AppIcon,
		AppDesc:        dist.AppDesc,
		UpdateMode:     int32(dist.UpdateMode),
		DistributeUrl:  dist.DistributeUrl,
		UrlScheme:      dist.UrlScheme,
		DistributeName: dist.DistributeName,
		StableVersion:  dist.StableVersion,
		StableDownUrl:  dist.StableDownUrl,
		BetaVersion:    dist.BetaVersion,
		BetaDownUrl:    dist.BetaDownUrl,
		CreateTime:     int64(dist.CreateTime),
		UpdateTime:     int64(dist.UpdateTime),
	}, nil
}

// GetAppVersion 获取应用版本
func (a *appServiceImpl) GetAppVersion(_ context.Context, id *proto.AppVersionId) (*proto.SAppVersion, error) {
	ia := a._sysRepo.GetSystemAggregateRoot().Application()
	version := ia.GetAppVersion(int(id.Value))
	if version == nil {
		return nil, errors.New("no such version")
	}
	return &proto.SAppVersion{
		Id:              int64(version.Id),
		DistributionId:  int64(version.DistributionId),
		Version:         version.Version,
		VersionCode:     int32(version.VersionCode),
		TerminalOs:      version.TerminalOs,
		TerminalChannel: version.TerminalChannel,
		StartTime:       int64(version.StartTime),
		UpdateMode:      int32(version.UpdateMode),
		UpdateContent:   version.UpdateContent,
		PackageUrl:      version.PackageUrl,
		IsForce:         int32(version.IsForce),
		IsNotified:      int32(version.IsNotified),
		CreateTime:      int64(version.CreateTime),
		UpdateTime:      int64(version.UpdateTime),
	}, nil
}

// DeleteAppDistribution 删除应用分发
func (a *appServiceImpl) DeleteAppDistribution(_ context.Context, req *proto.SysAppDistributionId) (*proto.TxResult, error) {
	err := a._sysRepo.GetSystemAggregateRoot().Application().DeleteAppDistribution(int(req.Value))
	return a.errorV2(err), nil
}

// DeleteAppVersion 删除应用版本
func (a *appServiceImpl) DeleteAppVersion(_ context.Context, id *proto.AppVersionId) (*proto.TxResult, error) {
	err := a._sysRepo.GetSystemAggregateRoot().Application().DeleteAppVersion(int(id.Value))
	return a.errorV2(err), nil
}

func (a *appServiceImpl) CheckAppVersion(_ context.Context, req *proto.CheckAppVersionRequest) (*proto.CheckAppVersionResponse, error) {
	ia := a._sysRepo.GetSystemAggregateRoot().Application()
	prod := ia.GetAppDistributionByName(req.AppName)
	if prod == nil {
		return &proto.CheckAppVersionResponse{
			VersionInfo:   "应用不存在",
			HasNewVersion: false,
		}, nil
	}
	ver := ia.GetLatestVersion(prod.Id, req.TerminalOS, req.TerminalChannel)
	if ver == nil {
		return &proto.CheckAppVersionResponse{
			VersionInfo:   "没有最近发布的版本",
			HasNewVersion: false,
		}, nil
	}
	if v := util.IntVersion(req.Version); ver.VersionCode <= v {
		return &proto.CheckAppVersionResponse{
			LatestVersion: req.Version,
			VersionInfo:   "当前为最新版本",
			HasNewVersion: false,
		}, nil
	}

	return &proto.CheckAppVersionResponse{
		HasNewVersion: true,
		LatestVersion: ver.Version,
		PackageUrl:    ver.PackageUrl,
		VersionInfo:   ver.UpdateContent,
		ForceUpdate:   ver.IsForce == 1,
		UpdateMode:    int32(ver.UpdateMode),
		ReleaseTime:   int64(ver.StartTime),
	}, nil
}
