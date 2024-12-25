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
	// IApplicationManager 应用管理器
	IApplicationManager interface {
		// AddLog 添加日志
		AddLog(l *SysLog) error
		// DeleteLog 删除日志
		DeleteLog(ids []int) error
		// CleanLog 清理日志
		CleanLog(days int) error
		// GetAllAppDistributions 获取所有应用分发
		GetAllAppDistributions() []*SysAppDistribution
		// GetAppDistribution 获取应用分发
		GetAppDistribution(id int) *SysAppDistribution
		// GetAppDistributionByName 获取应用分发
		GetAppDistributionByName(name string) *SysAppDistribution
		// SaveAppDistribution 保存应用分发
		SaveAppDistribution(distribution *SysAppDistribution) error
		// DeleteAppDistribution 删除应用分发
		DeleteAppDistribution(id int) error
		// GetAppVersion 获取应用版本
		GetAppVersion(id int) *SysAppVersion
		// SaveAppVersion 保存应用版本
		SaveAppVersion(version *SysAppVersion) error
		// DeleteAppVersion 删除应用版本
		DeleteAppVersion(id int) error
		// GetLatestVersion 获取最新版本
		GetLatestVersion(distributionId int, terminalOS, terminalChannel string) *SysAppVersion
	}

	// IApplicationRepository 应用仓储
	IApplicationRepository interface {
		// Log 获取日志仓储
		Log() fw.Repository[SysLog]
		// Distribution 获取应用分发仓储
		Distribution() fw.Repository[SysAppDistribution]
		// Version 获取应用版本仓储
		Version() fw.Repository[SysAppVersion]
	}
)

// SysLog 日志记录
type SysLog struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
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
	// 终端入口
	TerminalEntry string `json:"terminalEntry" db:"terminal_entry" gorm:"column:terminal_entry" bson:"terminalEntry"`
	// 终端名称
	TerminalName string `json:"terminalName" db:"terminal_name" gorm:"column:terminal_name" bson:"terminalName"`
	// 终端设备型号
	TerminalModel string `json:"terminalModel" db:"terminal_model" gorm:"column:terminal_model" bson:"terminalModel"`
	// 终端应用版本
	TerminalVersion string `json:"terminalVersion" db:"terminal_version" gorm:"column:terminal_version" bson:"terminalVersion"`
	// 额外信息
	ExtraInfo string `json:"extraInfo" db:"extra_info" gorm:"column:extra_info" bson:"extraInfo"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (s SysLog) TableName() string {
	return "sys_log"
}

// SysAppVersion 应用版本
type SysAppVersion struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" bson:"id"`
	// 分发应用编号
	DistributionId int `json:"distributionId" db:"distribution_id" gorm:"column:distribution_id" bson:"distributionId"`
	// 版本号
	Version string `json:"version" db:"version" gorm:"column:version" bson:"version"`
	// 版本数字代号
	VersionCode int `json:"versionCode" db:"version_code" gorm:"column:version_code" bson:"versionCode"`
	// 终端系统, 如: android / ios
	TerminalOs string `json:"terminalOs" db:"terminal_os" gorm:"column:terminal_os" bson:"terminalOs"`
	// 更新通道, beta: 测试版 nightly:每夜版 stable: 正式版本
	TerminalChannel string `json:"terminalChannel" db:"terminal_channel" gorm:"column:terminal_channel" bson:"terminalChannel"`
	// 开始时间
	StartTime int `json:"startTime" db:"start_time" gorm:"column:start_time" bson:"startTime"`
	// 更新模式, 1:包更新  2: 更新通知
	UpdateMode int `json:"updateMode" db:"update_mode" gorm:"column:update_mode" bson:"updateMode"`
	// 更新内容
	UpdateContent string `json:"updateContent" db:"update_content" gorm:"column:update_content" bson:"updateContent"`
	// 下载包地址
	PackageUrl string `json:"packageUrl" db:"package_url" gorm:"column:package_url" bson:"packageUrl"`
	// 是否强制更新
	IsForce int `json:"isForce" db:"is_force" gorm:"column:is_force" bson:"isForce"`
	// 是否已完成通知,完成后结束更新
	IsNotified int `json:"isNotified" db:"is_notified" gorm:"column:is_notified" bson:"isNotified"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (s SysAppVersion) TableName() string {
	return "sys_app_version"
}

// SysAppDistribution APP产品
type SysAppDistribution struct {
	// 产品编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 英文应用名称,如:mall
	AppName string `json:"appName" db:"app_name" gorm:"column:app_name" bson:"appName"`
	// AppIcon
	AppIcon string `json:"appIcon" db:"app_icon" gorm:"column:app_icon" bson:"appIcon"`
	// AppDesc
	AppDesc string `json:"appDesc" db:"app_desc" gorm:"column:app_desc" bson:"appDesc"`
	// UpdateMode
	UpdateMode int `json:"updateMode" db:"update_mode" gorm:"column:update_mode" bson:"updateMode"`
	// DistributeUrl
	DistributeUrl string `json:"distributeUrl" db:"distribute_url" gorm:"column:distribute_url" bson:"distributeUrl"`
	// UrlScheme
	UrlScheme string `json:"urlScheme" db:"url_scheme" gorm:"column:url_scheme" bson:"urlScheme"`
	// DistributeName
	DistributeName string `json:"distributeName" db:"distribute_name" gorm:"column:distribute_name" bson:"distributeName"`
	// StableVersion
	StableVersion string `json:"stableVersion" db:"stable_version" gorm:"column:stable_version" bson:"stableVersion"`
	// StableDownUrl
	StableDownUrl string `json:"stableDownUrl" db:"stable_down_url" gorm:"column:stable_down_url" bson:"stableDownUrl"`
	// BetaVersion
	BetaVersion string `json:"betaVersion" db:"beta_version" gorm:"column:beta_version" bson:"betaVersion"`
	// BetaDownUrl
	BetaDownUrl string `json:"betaDownUrl" db:"beta_down_url" gorm:"column:beta_down_url" bson:"betaDownUrl"`
	// CreateTime
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// UpdateTime
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (s SysAppDistribution) TableName() string {
	return "sys_app_distribution"
}
