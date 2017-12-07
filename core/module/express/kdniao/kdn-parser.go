package kdniao

import "go2o/core/domain/interface/shipment"

func Parse(shipperCode, logisticCode string, v *TraceResult) *shipment.ShipOrderTrace {
	r := &shipment.ShipOrderTrace{
		LogisticCode: logisticCode,
		ShipperCode:  shipperCode,
		ShipState:    v.State,
		UpdateTime:   0,
		Flows:        []*shipment.ShipFlow{},
	}
	for _, v := range v.Traces {
		r.Flows = append(r.Flows, &shipment.ShipFlow{
			Subject:    v.AcceptStation,
			CreateTime: v.AcceptTime,
			Remark:     "",
		})
	}
	return r
}
