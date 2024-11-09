/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: log.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-11-09 09:26:15
 * description: 日志接口定义
 * history:
 */
package sys

import "github.com/ixre/go2o/core/infrastructure/fw"

type LogLevel int

const (
	// 不记录
	LogLevelNone LogLevel = 0
	// 信息
	LogLevelInfo LogLevel = 1
	// 警告
	LogLevelWarn LogLevel = 2
	// 错误
	LogLevelError LogLevel = 3
	// 全部
	LogLevelAll LogLevel = 4
)

type (
	// ILogApp 日志服务
	ILogManager interface {
		// GetApp 获取日志应用
		GetApp(name string) *LogApp
		// AddLog 添加日志
		AddLog(l *SysLog) error
		// CleanLog 清理日志
		CleanLog(appId int, days int) error
	}

	// ILogRepository 日志仓储
	ILogRepository interface {
		// Log 获取日志仓储
		Log() fw.Repository[SysLog]
		// App 获取日志应用仓储
		App() fw.Repository[LogApp]
	}
)

// LogApp 日志应用
type LogApp struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// Name
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 日志级别, 0: 不记录  1: 信息 2:警告  3: 错误 4:全部
	LogLevel int `json:"logLevel" db:"log_level" gorm:"column:log_level" bson:"logLevel"`
}

func (s LogApp) TableName() string {
	return "sys_log_app"
}

// SysLog 日志记录
type SysLog struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 应用编号
	AppId int `json:"appId" db:"app_id" gorm:"column:app_id" bson:"appId"`
	// 用户编号
	UserId int `json:"userId" db:"user_id" gorm:"column:user_id" bson:"userId"`
	// 用户名
	Username string `json:"username" db:"username" gorm:"column:username" bson:"username"`
	// 日志级别, 1:信息  2: 警告  3: 错误 4: 其他
	LogLevel int `json:"logLevel" db:"log_level" gorm:"column:log_level" bson:"logLevel"`
	// Message
	Message string `json:"message" db:"message" gorm:"column:message" bson:"message"`
	// 参数
	Arguments string `json:"arguments" db:"arguments" gorm:"column:arguments" bson:"arguments"`
	// 终端设备型号
	TerminalModel string `json:"terminalModel" db:"terminal_model" gorm:"column:terminal_model" bson:"terminalModel"`
	// 终端名称
	TerminalName string `json:"terminalName" db:"terminal_name" gorm:"column:terminal_name" bson:"terminalName"`
	// 终端应用版本
	TerminalVersion string `json:"terminalVersion" db:"terminal_version" gorm:"column:terminal_version" bson:"terminalVersion"`
	// 额外信息
	ExtraInfo string `json:"extraInfo" db:"extra_info" gorm:"column:extra_info" bson:"extraInfo"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (s SysLog) TableName() string {
	return "sys_log_list"
}
