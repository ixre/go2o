package query

import (
	"github.com/ixre/go2o/core/domain/interface/ad"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type AdvertisementQuery struct {
	fw.BaseRepository[ad.Ad]
}

func NewAdvertisementQuery(o fw.ORM) *AdvertisementQuery {
	q := &AdvertisementQuery{}
	q.ORM = o
	return q
}
