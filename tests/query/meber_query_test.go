package query

import (
	"fmt"
	"testing"

	"github.com/ixre/go2o/core/inject"
	_ "github.com/ixre/go2o/tests"
	"github.com/ixre/gof/typeconv"
)

func TestQueryMemberWalletLog(t *testing.T) {
	q := inject.GetMemberQueryService()
	count, rows := q.PagedWalletAccountLog(723, 0, 0, 20, "", "")
	t.Log("count:", count)
	t.Log(fmt.Sprintf("rows:%#v", rows))
	t.Log(typeconv.MustJson(rows))
}
