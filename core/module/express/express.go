package express

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/log"
	"go2o/core/module"
	"go2o/core/module/express/kdniao"
)

var _ module.Module = new(ExpressModule)

type ExpressModule struct {
	_app gof.App
}

func (e *ExpressModule) SetApp(app gof.App) {
	e._app = app
}

func (e *ExpressModule) Init() {
	userId := e._app.Registry().GetString("express:kdn:user_id")
	appKey := e._app.Registry().GetString("express:kdn:api_key")
	kdniao.EBusinessID = userId
	kdniao.AppKey = appKey
	log.Println("--- KDN :", userId, appKey)
}
