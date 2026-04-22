package module

import (
	"errors"

	"github.com/ixre/go2o/internal/module/express/kdniao"
	"github.com/ixre/go2o/pkg/interface/domain/shipment"
	"github.com/ixre/gof"
)

var _ Module = new(ExpressModule)

type ExpressModule struct {
	_app gof.App
}

func (e *ExpressModule) SetApp(app gof.App) {
	e._app = app
}

func (e *ExpressModule) Init() {

}

func (e *ExpressModule) GetLogisticFlowTrack(shipperCode string, logisticCode string, invert bool) (*shipment.ShipOrderTrack, error) {
	r, err := kdniao.KdnTraces(shipperCode, logisticCode)
	if err == nil {
		if r == nil {
			return nil, errors.New("数据获取失败")
		}
		return kdniao.Parse(shipperCode, logisticCode, r, invert), nil
	}
	return nil, err
}
