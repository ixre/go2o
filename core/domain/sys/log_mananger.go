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
	"time"

	"github.com/ixre/go2o/core/domain/interface/sys"
)

var _ sys.ILogManager = &LogManagerImpl{}

type LogManagerImpl struct {
	_repo sys.ISystemRepo
}

func newLogManager(repo sys.ISystemRepo) sys.ILogManager {
	return &LogManagerImpl{_repo: repo}
}

func (l *LogManagerImpl) GetApp(name string) *sys.LogApp {
	return l._repo.Log().App().FindBy("name = ?", name)
}

// AddLog implements sys.ILogManager.
func (l *LogManagerImpl) AddLog(log *sys.SysLog) error {
	_, err := l._repo.Log().Log().Save(log)
	return err
}

// CleanLog implements sys.ILogManager.
func (l *LogManagerImpl) CleanLog(appId int, days int) error {
	model := sys.SysLog{}
	r := l._repo.Log().Log().Raw()
	lastTime := time.Now().Unix() - int64(days*86400)
	tx := r.Model(&model).Delete("app_id = ? and create_time < ?", appId, lastTime)
	return tx.Error
}
