package query

import (
	"github.com/ixre/go2o/pkg/initial/provide"
	"github.com/ixre/gof/db/orm"
)

func getOrm() orm.Orm {
	return provide.GetOrmInstance()
}
