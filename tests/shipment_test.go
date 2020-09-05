package tests

import (
	"go2o/tests/ti"
	"testing"
)

func TestGetShipmentOrderByOrderId(t *testing.T) {
	repo := ti.Factory.GetShipmentRepo()
	list := repo.GetShipOrders(4, true)
	for _, v := range list {
		t.Logf("%#v", v.Value())
	}
}
