package testing

import (
	"go2o/core/testing/ti"
	"testing"
)

func TestGetShipmentOrderByOrderId(t *testing.T) {
	repo := ti.ShipmentRepo
	list := repo.GetShipOrders(4, true)
	for _, v := range list {
		t.Logf("%#v", v.Value())
	}
}
