package kdniao

import (
	"go2o/core/domain/interface/shipment"
)

func Parse(shipperCode, logisticCode string, v *TraceResult, invert bool) *shipment.ShipOrderTrack {
	r := &shipment.ShipOrderTrack{
		LogisticCode: logisticCode,
		ShipperCode:  shipperCode,
		ShipState:    v.State,
		UpdateTime:   0,
		Flows:        []*shipment.ShipFlow{},
	}
	if invert {
		for i := len(v.Traces) - 1; i >= 0; i-- {
			v := v.Traces[i]
			r.Flows = append(r.Flows, &shipment.ShipFlow{
				Subject:    v.AcceptStation,
				CreateTime: v.AcceptTime,
				Remark:     "",
			})
		}
	} else {
		for _, v := range v.Traces {
			r.Flows = append(r.Flows, &shipment.ShipFlow{
				Subject:    v.AcceptStation,
				CreateTime: v.AcceptTime,
				Remark:     "",
			})
		}
	}
	return r
}
