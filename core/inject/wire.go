//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/impl"
)

func GetStationQueryService() *query.StationQuery {
	panic(wire.Build(
		query.NewStationQuery,
		impl.GetOrmInstance,
	))
}
