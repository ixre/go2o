package query

import (
	"github.com/ixre/go2o/core/domain/interface/ad"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type AdvertisementQuery struct {
	fw.BaseRepository[ad.Ad]
	positionRepo fw.BaseRepository[ad.Position]
}

func NewAdvertisementQuery(o fw.ORM) *AdvertisementQuery {
	q := &AdvertisementQuery{}
	q.ORM = o
	q.positionRepo.ORM = o
	return q
}

func (a *AdvertisementQuery) QueryPagingPositions(p *fw.PagingParams) (*fw.PagingResult, error) {
	tables := `ad_position
        LEFT JOIN ad_list ON ad_list.id = ad_position.put_aid`
	fields := `ad_position.*,
        ad_list.name AS ad_title`
	return fw.UnifinedQueryPaging(a.ORM, p, tables, fields)
}

func (a *AdvertisementQuery) QueryPagingAdList(p *fw.PagingParams) (*fw.PagingResult, error) {
	tables := `ad_position
        LEFT JOIN ad_list ON ad_list.id = ad_position.put_aid`
	fields := `ad_position.*,
        ad_list.name AS ad_title`
	return fw.UnifinedQueryPaging(a.ORM, p, tables, fields)
}
