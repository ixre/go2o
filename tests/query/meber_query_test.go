package query

import (
	"fmt"
	"testing"

	"github.com/ixre/go2o/core/query"
	"github.com/ixre/gof/types/typeconv"
	_ "github.com/ixre/go2o/tests"

)


func TestQueryMemberWalletLog(t *testing.T) {
	q :=  query.NewMemberQuery(getOrm())
	count, rows := q.PagedWalletAccountLog(723, 0,0, 20, "", "")
	t.Log("count:", count)
	t.Log(fmt.Sprintf("rows:%#v", rows))
	t.Log(typeconv.MustJson(rows))
}

