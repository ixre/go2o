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
	OAuthWechat = "WECHAT"
	// 微信小程序
	OAuthWechatMiniProgram = "WECHAT_MINI_PROGRAM"
	// QQ
	OAuthQQ = "QQ"
	// 微博
	OAuthWeibo = "WEIBO"
	// 抖音
	OAuthDouyin = "DOUYIN"
	// 苹果
	OAuthApple = "APPLE"
	// 谷歌
	OAuthGoogle = "GOOGLE"
	// 钉钉
	OAuthDingTalk = "DINGTALK"
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
