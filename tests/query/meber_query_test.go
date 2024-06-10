package query

import (
	"fmt"
	"testing"

	"github.com/ixre/go2o/core/query"
	_ "github.com/ixre/go2o/tests"
	"github.com/ixre/gof/types/typeconv"
)

func TestQueryMemberWalletLog(t *testing.T) {
	q := query.NewMemberQuery(getOrm())
	count, rows := q.PagedWalletAccountLog(723, 0, 0, 20, "", "")
	t.Log("count:", count)
	t.Log(fmt.Sprintf("rows:%#v", rows))
	t.Log(typeconv.MustJson(rows))
}
