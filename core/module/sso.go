/**
 * Copyright 2015 @ at3.net.
 * name : sso.go
 * author : jarryliu
 * date : 2016-11-25 13:02
 * description :
 * history :
 */
package module

import (
    "github.com/jsix/gof"
    "go2o/core/service/thrift/idl/gen-go/define"
    "errors"
    "strings"
    "github.com/jsix/gof/crypto"
)

var _ Module = new(SSOModule)

type SSOModule struct {
    app    gof.App
    apps   map[string]*define.SsoApp
    apiArr []string
}


// 模块数据
func (s *SSOModule) SetApp(app gof.App) {
    app = app
}
// 初始化模块
func (s *SSOModule)Init() {
    s.apps = make(map[string]*define.SsoApp)
}

func (s *SSOModule) Register(app *define.SsoApp) (token string, err error) {
    if app.Name == "" {
        return "", errors.New("-1:app name is null")
    }
    if app.ApiUrl == "" || (!strings.HasPrefix(app.ApiUrl,
        "https//") && !strings.HasPrefix(app.ApiUrl, "http://")) {
        return "", errors.New("-1:api url error")
    }
    if _, ok := s.apps[app.Name]; ok {
        return "", errors.New("-2:app has be registed")
    }
    // 生成TOKEN
    app.Token = crypto.Md5([]byte(app.Name + "#" + app.ApiUrl))
    // 注册
    s.apiArr = nil
    s.apps[app.Name] = app
    return app.Token, nil
}

// 返回同步的应用API地址
func (s *SSOModule) Array() []string {
    if s.apiArr == nil && s.apps != nil && len(s.apps) > 0 {
        s.apiArr = make([]string, len(s.apps))
        i := 0
        for _, v := range s.apps {
            s.apiArr[i] = s.formatApi(v.ApiUrl, v.Token)
            i++
        }
    }
    return s.apiArr
}

// 格式化API地址，加上token参数
func (s *SSOModule) formatApi(api string, token string) string {
    arr := []string{api}
    if strings.Index(api, "?") == -1 {
        arr = append(arr, "?")
    } else {
        arr = append(arr, "&")
    }
    arr = append(arr, "sso_token=")
    arr = append(arr, token)
    return strings.Join(arr, "")
}