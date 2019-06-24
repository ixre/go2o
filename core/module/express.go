package module

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/log"
	"go2o/core/domain/interface/shipment"
	"go2o/core/module/express/kdniao"
	"go2o/core/service/thrift"
	"os"
)

var _ Module = new(ExpressModule)

type ExpressModule struct {
	_app gof.App
}

func (e *ExpressModule) SetApp(app gof.App) {
	e._app = app
}

func (e *ExpressModule) Init() {
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		keys := []string{"express_kdn_business_id", "express_kdn_api_key"}
		cli.CreateUserRegistry(thrift.Context, keys[0], "1314567", "快递鸟接口业务ID")
		cli.CreateUserRegistry(thrift.Context, keys[1], "27d809c3-51b6-479c-9b77-6b98d7f3d41", "快递鸟接口KEY")
		data, _ := cli.GetRegistries(thrift.Context, keys)
		kdniao.EBusinessID = data[keys[0]]
		kdniao.AppKey = data[keys[1]]
	} else {
		log.Println("intialize express module error:", err.Error())
		os.Exit(1)
	}
}
func (e *ExpressModule) GetLogisticFlowTrack(shipperCode string, logisticCode string, invert bool) (*shipment.ShipOrderTrack, error) {
	r, err := kdniao.KdnTraces(shipperCode, logisticCode)
	if err == nil {
		return kdniao.Parse(shipperCode, logisticCode, r, invert), nil
	}
	return nil, err
}
