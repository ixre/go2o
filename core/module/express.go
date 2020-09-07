package module

import (
	"context"
	"github.com/ixre/gof"
	"github.com/ixre/gof/log"
	"go2o/core/domain/interface/shipment"
	"go2o/core/module/express/kdniao"
	"go2o/core/service"
	"go2o/core/service/proto"
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
	trans, cli, err := service.RegistryServeClient()
	if err == nil {
		defer trans.Close()
		keys := []string{"express_kdn_business_id", "express_kdn_api_key"}
		_, _ = cli.CreateUserRegistry(context.TODO(),
			&proto.UserRegistryCreateRequest{
				Key:          keys[0],
				DefaultValue: "1314567",
				Description:  "快递鸟接口业务ID",
			})
		_, _ = cli.CreateUserRegistry(context.TODO(),
			&proto.UserRegistryCreateRequest{
				Key:          keys[1],
				DefaultValue: "27d809c3-51b6-479c-9b77-6b98d7f3d41",
				Description:  "快递鸟接口KEY",
			})
		data, _ := cli.GetRegistries(context.TODO(),&proto.StringArray{Value: keys})
		kdniao.EBusinessID = data.Value[keys[0]]
		kdniao.AppKey = data.Value[keys[1]]
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
