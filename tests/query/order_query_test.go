package query

import (
	"fmt"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/tests/ti"
	"testing"
)

func TestQueryMemberNormalOrderList(t *testing.T){
	q := query.NewOrderQuery(ti.GetOrm())
	count, orders := q.QueryPagingNormalOrder(1, 10, 5, false, "", "")
	t.Log("count:",count)
	t.Log(fmt.Sprintf("orders:%#v",orders))
}
