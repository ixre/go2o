package query

import (
	"testing"

	"github.com/ixre/go2o/pkg/infrastructure/fw"
	"github.com/ixre/go2o/pkg/inject"
	_ "github.com/ixre/go2o/tests"
)

func TestQueryPagingMerchantList(t *testing.T) {
	ms := inject.GetMerchantQueryService()
	p := fw.PagingParams{
		Begin:     0,
		Size:      0,
		Order:     "",
		Arguments: []interface{}{},
	}
	ret, err := ms.QueryPagingMerchantList(&p)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%#v", ret)
}
