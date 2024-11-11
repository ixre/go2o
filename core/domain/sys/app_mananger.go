/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: log_mananger.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-11-09 09:26:29
 * description:
 * history:
 */

package sys

import (
	"errors"
	"time"

	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/util"
)

var _ sys.IApplicationManager = &LogManagerImpl{}

type LogManagerImpl struct {
	_repo sys.ISystemRepo
}

func newLogManager(repo sys.ISystemRepo) sys.IApplicationManager {
	return &LogManagerImpl{_repo: repo}
}

// AddLog implements sys.ILogManager.
func (l *LogManagerImpl) AddLog(log *sys.SysLog) error {
	log.CreateTime = int(time.Now().Unix())
	if len(log.TerminalName) > 40 {
		log.TerminalName = log.TerminalName[:40]
	}
	if len(log.TerminalModel) > 40 {
		log.TerminalModel = log.TerminalModel[:40]
	}
	_, err := l._repo.App().Log().Save(log)
	return err
}

// CleanLog implements sys.ILogManager.
func (l *LogManagerImpl) CleanLog(days int) error {
	model := sys.SysLog{}
	r := l._repo.App().Log().Raw()
	lastTime := time.Now().Unix() - int64(days*86400)
	tx := r.Model(&model).Delete("create_time < ?", lastTime)
	return tx.Error
}

// DeleteLog 删除日志
func (l *LogManagerImpl) DeleteLog(ids []int) error {
	_, err := l._repo.App().Log().DeleteBy("id in ?", ids)
	return err
}

// GetAppDistribution 获取应用分发
func (l *LogManagerImpl) GetAppDistribution(id int) *sys.SysAppDistribution {
	return l._repo.App().Distribution().Get(id)
}

// GetAppDistributionByName 获取应用分发
func (l *LogManagerImpl) GetAppDistributionByName(name string) *sys.SysAppDistribution {
	return l._repo.App().Distribution().FindBy("app_name = ?", name)
}

// SaveAppDistribution implements sys.IApplicationManager.
func (l *LogManagerImpl) SaveAppDistribution(a *sys.SysAppDistribution) error {
	repo := l._repo.App().Distribution()
	if c, _ := repo.Count("app_name=? AND id <> ?",
		a.AppName, a.Id); c > 0 {
		return errors.New("已经存在分发应用")
	}
	var dst *sys.SysAppDistribution
	if a.Id > 0 {
		dst = repo.Get(a.Id)
	} else {
		if len(a.AppName) == 0 {
			return errors.New("应用名称不能为空")
		}
		dst = &sys.SysAppDistribution{}
		dst.CreateTime = int(time.Now().Unix())
		dst.AppName = a.AppName
		dst.UpdateMode = int(a.UpdateMode)

	}
	dst.AppIcon = a.AppIcon
	dst.AppDesc = a.AppDesc
	dst.DistributeUrl = a.DistributeUrl
	dst.DistributeName = a.DistributeName
	dst.StableVersion = a.StableVersion
	dst.StableDownUrl = a.StableDownUrl
	dst.BetaVersion = a.BetaVersion
	dst.BetaDownUrl = a.BetaDownUrl
	_, err := repo.Save(dst)
	return err
}

// DeleteAppDistribution 删除应用分发
func (l *LogManagerImpl) DeleteAppDistribution(id int) error {
	// 检测是否存在版本
	if len, _ := l._repo.App().Version().Count("distribution_id = ?", id); len > 0 {
		return errors.New("已经存在版本,不能删除")
	}
	return l._repo.App().Distribution().Delete(&sys.SysAppDistribution{Id: id})
}

// GetAppVersion 获取应用版本
func (l *LogManagerImpl) GetAppVersion(id int) *sys.SysAppVersion {
	return l._repo.App().Version().Get(id)
}

// SaveVersion implements sys.IApplicationManager.
func (l *LogManagerImpl) SaveAppVersion(version *sys.SysAppVersion) error {
	app := l.GetAppDistribution(version.DistributionId)
	if app == nil {
		return errors.New("分发应用不存在")
	}
	if len(version.TerminalChannel) == 0 {
		return errors.New("版本类型不能为空")
	}
	repo := l._repo.App().Version()
	var dst *sys.SysAppVersion
	if version.Id > 0 {
		dst = repo.Get(version.Id)
		if dst.TerminalOs != version.TerminalOs {
			return errors.New("终端类型不能修改")
		}
		if dst.TerminalChannel != version.TerminalChannel {
			return errors.New("终端更新渠道不能修改")
		}
	} else {
		version.CreateTime = int(time.Now().Unix())
	}
	if c, _ := repo.Count("version=? AND id <> ?", version.Version, version.Id); c > 0 {
		return errors.New("版本已经存在")
	}
	if version.VersionCode <= 0 {
		// 自动计算版本号
		version.VersionCode = util.IntVersion(version.Version)
	}
	if version.UpdateMode == 1 {
		if len(version.PackageUrl) == 0 {
			// 如果包地址为空,则使用上次版本的包地址
			lastVer := l.GetLatestVersion(version.DistributionId, version.TerminalOs, version.TerminalChannel)
			if lastVer == nil {
				return errors.New("更新/安装包文件地址不能为空")
			}
			version.PackageUrl = lastVer.PackageUrl
		}
		if len(version.TerminalOs) == 0 {
			return errors.New("终端类型不能为空")
		}
	}
	// 更新时间
	if version.StartTime <= 0 {
		version.StartTime = int(time.Now().Unix())
	}
	// 更新模式不能修改
	version.UpdateMode = app.UpdateMode
	_, err := repo.Save(version)
	if err == nil {
		// 如果发布了更新版本,则更新最新的版本
		err = l.updateLatest(version)
	}
	return err
}

// DeleteVersion implements sys.IApplicationManager.
func (l *LogManagerImpl) DeleteAppVersion(id int) error {
	return l._repo.App().Version().Delete(&sys.SysAppVersion{Id: id})
}

// 更新分发应用为最新版本
func (l *LogManagerImpl) updateLatest(version *sys.SysAppVersion) error {
	app := l.GetAppDistribution(version.DistributionId)
	if app == nil {
		return errors.New("应用不存在")
	}
	latest := l.GetLatestVersion(app.Id, version.TerminalOs, version.TerminalChannel)
	if latest == nil || version.Id == latest.Id {
		// 如果是最新版本,则更新最新版本
		if version.TerminalChannel == "stable" {
			app.StableVersion = version.Version
		} else {
			app.BetaVersion = version.Version
		}
		app.UpdateTime = int(time.Now().Unix())
		_, err := l._repo.App().Distribution().Save(app)
		return err
	}
	return nil
}

// 获取最新版本
func (l *LogManagerImpl) GetLatestVersion(distId int, terminalOS, terminalChannel string) *sys.SysAppVersion {
	unix := time.Now().Unix()
	return l._repo.App().Version().FindBy(
		`distribution_id=? 
		AND terminal_os=? 
		AND terminal_channel=? 
		AND start_time BETWEEN 0 AND ?
		ORDER BY version_code DESC`,
		distId, terminalOS, terminalChannel, unix)
}
