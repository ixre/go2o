package sys

/**
 * Copyright (C) 2007-2025 fze.NET, All rights reserved.
 *
 * name: oauth.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2025-01-03 19:37:01
 * description:
 * history:
 */
// 第三方登录客户端类型
const (
	// 微信
	OAuthWechat = "wechat"
	// 微信小程序
	OAuthWechatMiniProgram = "wechat_mini_program"
	// QQ
	OAuthQQ = "qq"
	// 微博
	OAuthWeibo = "weibo"
	// 抖音
	OAuthDouyin = "douyin"
	// 苹果
	OAuthApple = "apple"
	// 谷歌
	OAuthGoogle = "google"
	// 钉钉
	OAuthDingTalk = "dingtalk"
)

// OAuthOpenIdResponse 第三方登录OpenId响应
type OAuthOpenIdResponse struct {
	// 第三方应用Id
	AppId string
	// 用户OpenId
	OpenId string
	// 额外信息
	Extra map[string]string
}

// IOAuthManager 第三方登录管理器
type IOAuthManager interface {
	// 获取第三方登录OpenId
	GetOpenId(appId int, clientType, clientCode string) (OAuthOpenIdResponse, error)
}
