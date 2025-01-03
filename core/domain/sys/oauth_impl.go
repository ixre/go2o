/**
 * Copyright (C) 2007-2025 fze.NET, All rights reserved.
 *
 * name: oauth_impl.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2025-01-03 19:43:15
 * description:
 * history:
 */

package sys

import (
	"errors"
	"strconv"

	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/sp/tencent"
)

type oauthManager struct {
}

func newOAuthManager() sys.IOAuthManager {
	return &oauthManager{}
}

// GetOpenId 获取第三方登录OpenId
func (s *oauthManager) GetOpenId(appId int, clientType, clientCode string) (sys.OAuthOpenIdResponse, error) {
	if len(clientType) == 0 {
		return sys.OAuthOpenIdResponse{}, errors.New("缺少参数: clientType")
	}
	if len(clientCode) == 0 {
		return sys.OAuthOpenIdResponse{}, errors.New("缺少参数: clientCode")
	}
	if clientType == sys.OAuthWechatMiniProgram {
		// 微信小程序
		ret, err := tencent.WECHAT.GetMiniProgramOpenId("", clientCode)
		if err != nil {
			return sys.OAuthOpenIdResponse{}, err
		}
		return sys.OAuthOpenIdResponse{
			AppId:  strconv.Itoa(appId),
			OpenId: ret.OpenID,
			Extra: map[string]string{
				"sessionKey": ret.SessionKey,
				"unionId":    ret.UnionID,
			},
		}, nil
	}
	return sys.OAuthOpenIdResponse{}, errors.New("不支持的第三方登录类型: " + clientType)
}
