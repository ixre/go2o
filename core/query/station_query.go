package query

import (
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

type StationQuery struct {
	db.Connector
	o orm.Orm
}

func NewStationQuery(o orm.Orm) *StationQuery {
	return &StationQuery{
		Connector: o.Connector(),
		o:         o,
	}
}
