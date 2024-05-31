package domain

import (
	"testing"

	"github.com/ixre/go2o/core/inject"
)

func TestGetShipmentOrderByOrderId(t *testing.T) {
	repo := inject.GetShipmentRepo()
	list := repo.GetShipOrders(4, true)
	for _, v := range list {
		t.Logf("%#v", v.Value())
	}
}
