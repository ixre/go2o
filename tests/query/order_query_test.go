package query

import (
	"testing"

	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/types/typeconv"
)

func TestQueryMemberNormalOrderList(t *testing.T) {
	var memberId int64 = 723
	q := query.NewOrderQuery(ti.GetOrm())
	_, orders := q.QueryPagingNormalOrder(memberId, 0, 50, false, "", "")
	t.Log("count:", len(orders))
	t.Log(typeconv.MustJson(orders[0]))
	//bytes, _ := json.Marshal(orders[0])
	//t.Log(string(bytes))
}
