package module

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/log"
	"go2o/core/domain/interface/shipment"
	"go2o/core/module/express/kdniao"
)

var _ Module = new(ExpressModule)

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
func (e *ExpressModule) GetLogisticFlowTrace(shipperCode string, logisticCode string) (*shipment.ShipOrderTrace, error) {
	r, err := kdniao.KdnTraces(shipperCode, logisticCode)
	if err == nil {
		return kdniao.Parse(shipperCode, logisticCode, r), nil
	}
	return nil, err
}
