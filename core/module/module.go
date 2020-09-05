/**
 * Copyright 2015 @ at3.net.
 * name : module.go
 * author : jarryliu
 * date : 2016-11-25 13:02
 * description :
 * history :
 */
package module

import (
	"errors"
	"github.com/ixre/gof"
	"github.com/ixre/gof/log"
	"sync"
)

var (
	mux       sync.Mutex
	moduleMap map[string]Module
	initOk    = false //是否已初始化
	SSO       = "sso"
	MM        = "member"
	PAY       = "payment"
	B4E       = "bank4e"
	EXPRESS   = "express"
)

// 模块实现
type Module interface {
	// 模块数据
	SetApp(app gof.App)
	// 初始化模块
	Init()
}

// 注册模块
func Register(name string, m Module) error {
	mux.Lock()
	defer mux.Unlock()
	if _, ok := moduleMap[name]; ok {
		return errors.New("已注册相同名称的模块")
	}
	moduleMap[name] = m
	return nil
}

// 初始化模块
func initModule() {
	app := gof.CurrentApp
	moduleMap = map[string]Module{}
	if app != nil {
		registerInternal() //注册内置的模块
		initOk = true      // 标记为已加载
		for k, v := range moduleMap {
			log.Println(" [ Module][ Load]: module => ", k)
			v.SetApp(app)
			v.Init()
		}
	}
}

// 注册内置的模块
func registerInternal() {
	Register(SSO, &SSOModule{})
	Register(MM, &MemberModule{})
	Register(PAY, &PaymentModule{})
	Register(B4E, &Bank4E{})
	Register(EXPRESS, &ExpressModule{})
}

// 获取模块
func Get(name string) Module {
	if !initOk {
		initModule()
	}
	return moduleMap[name]
}
