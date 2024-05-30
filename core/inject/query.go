//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/ixre/go2o/core/query"
)

var queryProvideSets = wire.NewSet(serviceProvideSets,
	query.NewStationQuery,
	query.NewMerchantQuery,
	query.NewOrderQuery,
)

func GetStationQueryService() *query.StationQuery {
	panic(wire.Build(queryProvideSets))
}
