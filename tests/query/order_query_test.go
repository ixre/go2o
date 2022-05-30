package query

import (
	"encoding/json"
	"fmt"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/tests/ti"
	"testing"
)

func TestQueryMemberNormalOrderList(t *testing.T){
	q := query.NewOrderQuery(ti.GetOrm())
	count, orders := q.QueryPagingNormalOrder(1, 0, 6, false, "", "")
	t.Log("count:",count)
	t.Log(fmt.Sprintf("orders:%#v",orders))
	bytes,_ := json.Marshal(orders[0])
	t.Log(string(bytes))
}
